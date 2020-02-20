package dns

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/sipt/shuttle/constant/typ"

	"github.com/miekg/dns"
	"github.com/pkg/errors"
	"github.com/sipt/shuttle/conf/model"
	"github.com/sirupsen/logrus"
)

const (
	TypSystem  = "system"
	TypStatic  = "static"
	TypDynamic = "dynamic"

	dnsTimeout = time.Minute * 10
)

type Handle func(ctx context.Context, domain string) *DNS

func ApplyConfig(config *model.Config, fallback Handle) (handle Handle, cache ICache, err error) {
	timeout := dnsTimeout
	if config.DNS.TimeoutSec != 0 {
		timeout = time.Second * time.Duration(config.DNS.TimeoutSec)
	}
	err = InitGeoIP()
	if err != nil {
		return
	}
	servers := config.DNS.Servers
	if config.DNS.IncludeSystem {
		// TODO Read File: hosts
	}
	handle = fallback
	handle, _ = newGeneralHandle(servers, timeout, handle)
	handle, cache, _ = newCacheHandle(handle)
	for i := len(config.DNS.Mapping) - 1; i >= 0; i-- {
		v := config.DNS.Mapping[i]
		handle, err = newMappingHandle(v.Domain, v.Server, v.IP, timeout, handle)
		if err != nil {
			return
		}
	}
	return
}

func ApplyRuntime(_ context.Context, _ typ.Runtime) error {
	return nil
}

func newGeneralHandle(servers []string, timeout time.Duration, next Handle) (Handle, error) {
	serverAddrs := make([]*DnsServer, len(servers))
	var err error
	for i, v := range servers {
		serverAddrs[i], err = ParseDnsServer(v)
		if err != nil {
			return nil, err
		}
	}
	return func(ctx context.Context, domain string) *DNS {
		reply := &DNS{
			Typ:    TypDynamic,
			Domain: domain,
		}
		var err error
		reply.IP, reply.CurrentServer, err = ResolveDomain(ctx, domain, serverAddrs...)
		reply.CurrentIP = SelectIP(reply.IP)
		reply.CurrentCountry = GeoLookUp(reply.CurrentIP)
		reply.ExpireAt = time.Now().Add(timeout)
		if err != nil {
			logrus.WithError(err).WithField("domain", domain).Error("lookup ip failed")
			next(ctx, domain)
		}
		return reply
	}, nil
}

// can override SelectIP
var SelectIP = func(ips []net.IP) net.IP {
	if len(ips) > 0 {
		return ips[0]
	}
	return nil
}

func newMappingHandle(mappingDomain string, servers []string, ips []string, timeout time.Duration, next Handle) (Handle, error) {
	if len(servers) == 0 && len(ips) == 0 {
		return nil, errors.Errorf("DNS.Mapping[domain:%s, server:%v, ip:%v], server and ip is empty", mappingDomain, servers, ips)
	}
	if len(servers) > 0 {
		serverAddrs := make([]*DnsServer, len(servers))
		var err error
		for i, v := range servers {
			serverAddrs[i], err = ParseDnsServer(v)
			if err != nil {
				return nil, err
			}
		}
		return func(ctx context.Context, domain string) *DNS {
			if mappingDomain[0] == '*' && strings.HasSuffix(domain, mappingDomain[1:]) {
			} else if mappingDomain == domain {
			} else {
				return next(ctx, domain)
			}
			reply := &DNS{
				Typ:           TypStatic,
				MappingDomain: mappingDomain,
				Domain:        domain,
				Server:        serverAddrs,
			}
			var err error
			reply.IP, reply.CurrentServer, err = ResolveDomain(ctx, domain, serverAddrs...)
			if err != nil {
				logrus.WithError(err).WithField("domain", domain).Error("lookup ip failed")
				return next(ctx, domain)
			}
			reply.ExpireAt = time.Now().Add(timeout)
			return reply
		}, nil
	} else {
		netIP := make([]net.IP, len(ips))
		for i, v := range ips {
			netIP[i] = net.ParseIP(v)
			if len(netIP[i]) == 0 {
				return nil, errors.Errorf("DNS.Mapping[domain:%s, server:%v, ip:%v], ip[%s] invalid", mappingDomain, servers, ips, v)
			}
		}

		return func(ctx context.Context, domain string) *DNS {
			if mappingDomain[0] == '*' && strings.HasSuffix(domain, mappingDomain[1:]) {
			} else if mappingDomain == domain {
			} else {
				return next(ctx, domain)
			}
			return &DNS{
				Typ:            TypStatic,
				MappingDomain:  mappingDomain,
				Domain:         domain,
				IP:             netIP,
				CurrentIP:      netIP[0],
				CurrentCountry: GeoLookUp(netIP[0]),
			}
		}, nil
	}
}

type DNS struct {
	Typ            string
	MappingDomain  string
	Domain         string
	IP             []net.IP
	Server         []*DnsServer
	CurrentServer  DnsServer
	CurrentIP      net.IP
	CurrentCountry string
	ExpireAt       time.Time
}

func ParseDnsServer(value string) (*DnsServer, error) {
	if strings.Index(value, "://") < 0 {
		value = "udp://" + value
	}
	u, err := url.Parse(value)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "udp" {
		return nil, errors.Errorf("[DNS.Servers] [%s] not support [%s]", value, u.Scheme)
	}
	s := &DnsServer{
		Network: u.Scheme,
	}
	s.IP = net.ParseIP(u.Hostname())
	if len(s.IP) == 0 {
		return nil, errors.Errorf("[DNS.Servers] [%s] [%s] ip invalid", value, u.Hostname())
	}
	if len(u.Port()) == 0 {
		s.Port = 53
	} else {
		s.Port, err = strconv.Atoi(u.Port())
		if err != nil {
			return nil, errors.Errorf("[DNS.Servers] [%s] [%s] port invalid", value, u.Port())
		}
	}
	return s, nil
}

type DnsServer struct {
	Network string
	IP      net.IP
	Port    int
}

func (d *DnsServer) Addr() string {
	return net.JoinHostPort(d.IP.String(), strconv.Itoa(d.Port))
}

func (d *DnsServer) String() string {
	return fmt.Sprintf("%s://%s", d.Network, d.Addr())
}

func (d *DNS) IsNil() bool {
	return len(d.Domain) == 0
}

func ResolveDomain(ctx context.Context, domain string, servers ...*DnsServer) (ips []net.IP, server DnsServer, err error) {
	type _reply struct {
		ips    []net.IP
		server *DnsServer
	}
	c := make(chan *_reply, 1)
	for _, v := range servers {
		go func(s *DnsServer) {
			m := &dns.Msg{}
			m.SetQuestion(dns.Fqdn(domain), dns.TypeA).
				RecursionDesired = true
			r, err := dns.Exchange(m, s.Addr())
			if err != nil {
				logrus.WithError(err).WithField("domain", domain).
					WithField("dns_server", s.String()).
					Error("resolve dns failed")
				return
			}
			ips := make([]net.IP, 0, len(r.Answer))
			for _, v := range r.Answer {
				a, ok := v.(*dns.A)
				if ok {
					ips = append(ips, a.A)
				}
			}
			select {
			case c <- &_reply{ips: ips, server: s}:
			default:
			}
		}(v)
	}
	select {
	case reply := <-c:
		return reply.ips, *reply.server, nil
	case <-ctx.Done():
		return nil, DnsServer{}, errors.Errorf("[ResolveDomain] context was canceled")
	}
}
