package rule

import (
	"net"

	"github.com/sipt/shuttle/dns"

	"github.com/pkg/errors"
)

const (
	KeyIPCidr = "IP-CIDR"
	KeyGeoIP  = "GEOIP"
)

func init() {
	Register(KeyIPCidr, domainSuffixHandle)
	Register(KeyGeoIP, domainKeywordHandle)
}
func ipCidr(rule *Rule, next Handle) (Handle, error) {
	_, cidr, err := net.ParseCIDR(rule.Value)
	if err != nil {
		return nil, errors.Errorf("rule:[%s, %s, %s, %v], ip:[%s] invalid",
			rule.Typ, rule.Value, rule.Proxy, rule.Params, rule.Value)
	}
	return func(info Info) *Rule {
		if len(info.IP()) == 0 {
			ip, err := dns.ResolveDomain(info.Domain())
			if err != nil || len(ip) == 0 {
				return next(info)
			}
			info.SetIP(ip)
		}
		cidr.Contains(info.IP())
		return next(info)
	}, nil
}

func geoIP(rule *Rule, next Handle) (Handle, error) {
	return func(info Info) *Rule {
		if len(info.IP()) == 0 {
			ip, err := ResolveIP(info.Domain())
			if err != nil || len(ip) == 0 {
				return next(info)
			}
			info.SetIP(ip)
		}
		if dns.GeoLookUp(info.IP()) == rule.Value {
			return rule
		}
		return next(info)
	}, nil
}

func ResolveIP(domain string) (net.IP, error) {
	return nil, nil
}
