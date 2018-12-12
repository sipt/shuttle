package dns

import (
	"fmt"
	"github.com/miekg/dns"
	"github.com/sipt/shuttle/log"
	"net"
	"strings"
	"time"
)

var dnsServers []string

type Answer struct {
	MatchType string
	Domain    string
	IPs       []string
	Server    string
	Port      string
	Type      string
	Country   string
	Duration  time.Duration
}

func (a *Answer) GetIP() string {
	if a == nil {
		return ""
	} else if len(a.IPs) == 0 {
		return ""
	}
	return a.IPs[0]
}

//resolve domain
func ResolveDomain(domain string) (answer *Answer, err error) {
LOOP:
	for _, v := range dnsConfig.localDNS {
		switch v.MatchType {
		case MatchTypeDomainSuffix:
			if strings.HasSuffix(domain, v.Domain) {
				log.Logger.Debug("[DNS] [Local] ", v.String())
				answer, err = localResolve(v, domain)
				break LOOP
			}
		case MatchTypeDomain:
			if domain == v.Domain {
				log.Logger.Debug("[DNS] [Local] ", v.String())
				answer, err = localResolve(v, domain)
				break LOOP
			}
		case MatchTypeDomainKeyword:
			if strings.Index(domain, v.Domain) >= 0 {
				log.Logger.Debug("[DNS] [Local] ", v.String())
				answer, err = localResolve(v, domain)
				break LOOP
			}
		}
	}
	if answer == nil && err == nil {
		//connect to DNS server
		answer = &Answer{
			MatchType: MatchNone,
			Domain:    domain,
			Type:      DNSTypeDirect,
		}
		start := time.Now()
		var err error
		answer.IPs, answer.Server, err = directResolve(dnsConfig.servers, domain)
		if err != nil {
			log.Logger.Errorf("[DNS] [direct] resolve domain [%s] failed: %s", domain, err.Error())
			return nil, err
		}
		answer.Duration = time.Now().Sub(start)
	}
	if answer != nil {
		if len(answer.IPs) == 0 {
			answer = nil
		} else {
			answer.Country = GeoLookUp(answer.GetIP())
		}
	}
	return
}

//resolve ip
func ResolveIP(ip string) (*Answer, error) {
	return &Answer{
		IPs:     []string{ip},
		Country: GeoLookUp(ip),
	}, nil
}

func localResolve(d *DNS, domain string) (*Answer, error) {
	answer := &Answer{
		MatchType: d.MatchType,
		Domain:    d.Domain,
		Port:      d.Port,
		Type:      d.Type,
		Country:   d.Country,
	}
	switch d.Type {
	case DNSTypeStatic:
		answer.IPs = d.IPs
		return answer, nil
	case DNSTypeDirect:
		//connect to DNS server
		start := time.Now()
		var err error
		answer.IPs, answer.Server, err = directResolve(d.DNSs, domain)
		if err != nil {
			log.Logger.Errorf("[DNS] [direct] resolve domain [%s] failed: %s", domain, err.Error())
			return nil, err
		}
		answer.Duration = time.Now().Sub(start)
	case DNSTypeRemote:
		return nil, nil
	}
	return answer, nil
}

type _Reply struct {
	Addr string
	Msg  *dns.Msg
}

func directResolve(servers []string, domain string) ([]string, string, error) {
	replyChan := make(chan *_Reply, 1)
	for _, s := range servers {
		go resolveDomain(s, "53", domain, replyChan)
	}
	timer := time.NewTimer(2 * time.Second)
	select {
	case reply := <-replyChan:
		var (
			a   *dns.A
			ok  bool
			ips = make([]string, 0, len(reply.Msg.Answer))
		)
		for _, v := range reply.Msg.Answer {
			a, ok = v.(*dns.A)
			if ok {
				ips = append(ips, a.A.String())
			}
		}
		if len(ips) == 0 {
			return nil, "", fmt.Errorf("resolve domain [%s] is empty", domain)
		}
		return ips, reply.Addr, nil
	case <-timer.C:
		log.Logger.Errorf("[DNS] [Local] resolve domain [%s] failed: %s", domain)
		return nil, "", fmt.Errorf("resolve domain [%s] failed: %s", domain)
	}
}

func resolveDomain(addr, port, domain string, c chan *_Reply) {
	m := &dns.Msg{}
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	m.RecursionDesired = true
	r, err := dns.Exchange(m, net.JoinHostPort(addr, port))
	if err != nil {
		log.Logger.Errorf("[DNS] [Local] connect [%s:%s] resolve domain [%s] failed: %s",
			addr, port, domain, err.Error())
		return
	}
	if r == nil || r.Rcode != dns.RcodeSuccess {
		log.Logger.Errorf("[DNS] [Local] connect [%s:%s] resolve domain [%s] failed ",
			addr, port, domain)
		return
	}
	select {
	case c <- &_Reply{Addr: addr, Msg: r}:
	default:
	}
}
