package rule

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"github.com/sipt/shuttle/dns"
)

const (
	KeyIPCidr = "IP-CIDR"
	KeyGeoIP  = "GEOIP"
)

func init() {
	Register(KeyIPCidr, ipCidrHandle)
	Register(KeyGeoIP, geoIPHandle)
}
func ipCidrHandle(rule *Rule, next Handle) (Handle, error) {
	_, cidr, err := net.ParseCIDR(rule.Value)
	if err != nil {
		return nil, errors.Errorf("rule:[%s, %s, %s, %v], ip:[%s] invalid",
			rule.Typ, rule.Value, rule.Proxy, rule.Params, rule.Value)
	}
	return func(ctx context.Context, info Info) *Rule {
		if len(info.IP()) == 0 {
			ip, err := dns.ResolveDomain(ctx, info.Domain())
			if err != nil || len(ip) == 0 || len(ip[0]) == 0 {
				info.SetIP([]byte{})
				return next(ctx, info)
			}
			info.SetIP(ip[0])
		} else if len(info.IP()) > 0 {
			cidr.Contains(info.IP())
		}
		return next(ctx, info)
	}, nil
}

func geoIPHandle(rule *Rule, next Handle) (Handle, error) {
	return func(ctx context.Context, info Info) *Rule {
		if info.IP() == nil {
			ip, err := dns.ResolveDomain(ctx, info.Domain())
			if err != nil || len(ip) == 0 || len(ip[0]) == 0 {
				info.SetIP([]byte{})
				return next(ctx, info)
			}
			info.SetIP(ip[0])
		} else if len(info.IP()) > 0 {
			if dns.GeoLookUp(info.IP()) == rule.Value {
				return rule
			}
		}
		return next(ctx, info)
	}, nil
}
