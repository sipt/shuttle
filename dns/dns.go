package dns

import (
	"context"
	"net"
	"strings"

	"github.com/miekg/dns"
	"github.com/pkg/errors"
	"github.com/sipt/shuttle/conf/model"
	"github.com/sirupsen/logrus"
)

const (
	TypSystem  = "system"
	TypStatic  = "static"
	TypDynamic = "dynamic"
)

type Handle func(ctx context.Context, domain string) *DNS

func ApplyConfig(config *model.Config, fallback Handle) (handle Handle, err error) {
	servers := config.DNS.Servers
	if config.DNS.IncludeSystem {
		// TODO Read File: hosts
	}
	handle = fallback
	handle, _ = newGeneralHandle(servers, handle)
	handle, _ = newCacheHandle(handle)
	for i := len(config.DNS.Mapping) - 1; i >= 0; i-- {
		v := config.DNS.Mapping[i]
		handle, err = newMappingHandle(v.Domain, v.IP, v.Server, handle)
		if err != nil {
			return
		}
	}
	return
}

func newGeneralHandle(servers []string, next Handle) (Handle, error) {
	serverIP := make([]net.IP, len(servers))
	for i, v := range servers {
		serverIP[i] = net.ParseIP(v)
		if len(serverIP[i]) == 0 {
			return nil, errors.Errorf("[DNS.Servers] parse DNS server[%s] failed", v)
		}
	}
	return func(ctx context.Context, domain string) *DNS {
		reply := &DNS{
			Typ:    TypDynamic,
			Domain: domain,
		}
		var err error
		reply.IP, reply.CurrentServer, err = ResolveDomain(ctx, domain, serverIP...)
		if err != nil {
			logrus.WithError(err).WithField("domain", domain).Error("lookup ip failed")
			next(ctx, domain)
		}
		return reply
	}, nil
}

func newMappingHandle(mappingDomain string, server []string, ip []string, next Handle) (Handle, error) {
	if len(server) == 0 && len(ip) == 0 {
		return nil, errors.Errorf("DNS.Mapping[domain:%s, server:%v, ip:%v], server and ip is empty", mappingDomain, server, ip)
	}
	if len(server) > 0 {
		netIP := make([]net.IP, len(server))
		for i, v := range server {
			netIP[i] = net.ParseIP(v)
			if len(netIP[i]) == 0 {
				return nil, errors.Errorf("DNS.Mapping[domain:%s, server:%v, ip:%v], ip[%s] invalid", mappingDomain, server, ip, v)
			}
		}
		return func(ctx context.Context, domain string) *DNS {
			if mappingDomain[0] == '*' && strings.HasSuffix(domain, mappingDomain[1:]) {
			} else if mappingDomain == domain {
			} else {
				return next(ctx, domain)
			}
			reply := &DNS{
				Typ:            TypStatic,
				MappingDomain:  mappingDomain,
				Domain:         domain,
				IP:             netIP,
				CurrentIP:      netIP[0],
				CurrentCountry: GeoLookUp(netIP[0]),
			}
			var err error
			reply.IP, reply.CurrentServer, err = ResolveDomain(ctx, domain, netIP...)
			if err != nil {
				logrus.WithError(err).WithField("domain", domain).Error("lookup ip failed")
			}
			return reply
		}, nil
	} else {
		netIP := make([]net.IP, len(ip))
		for i, v := range ip {
			netIP[i] = net.ParseIP(v)
			if len(netIP[i]) == 0 {
				return nil, errors.Errorf("DNS.Mapping[domain:%s, server:%v, ip:%v], ip[%s] invalid", mappingDomain, server, ip, v)
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
	Server         []net.IP
	CurrentServer  net.IP
	CurrentIP      net.IP
	CurrentCountry string
}

func (d *DNS) IsNil() bool {
	return len(d.Domain) == 0
}

func ResolveDomain(ctx context.Context, domain string, servers ...net.IP) (ips []net.IP, server net.IP, err error) {
	dns.A{}
	return nil, nil, nil
}
