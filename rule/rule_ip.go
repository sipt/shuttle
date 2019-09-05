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
func ipCidrHandle(_ context.Context, rule *Rule, next Handle, dnsHandle dns.Handle) (Handle, error) {
	_, cidr, err := net.ParseCIDR(rule.Value)
	if err != nil {
		return nil, errors.Errorf("rule:[%s, %s, %s, %v], ip:[%s] invalid",
			rule.Typ, rule.Value, rule.Proxy, rule.Params, rule.Value)
	}
	return func(ctx context.Context, info RequestInfo) *Rule {
		if len(info.IP()) == 0 {
			dns := dnsHandle(ctx, info.Domain())
			if dns == nil || len(dns.CurrentIP) == 0 {
				info.SetIP([]byte{})
				return next(ctx, info)
			}
			info.SetIP(dns.CurrentIP)
		}
		if len(info.IP()) > 0 && cidr.Contains(info.IP()) {
			return rule
		}
		return next(ctx, info)
	}, nil
}

func geoIPHandle(_ context.Context, rule *Rule, next Handle, dnsHandle dns.Handle) (Handle, error) {
	return func(ctx context.Context, info RequestInfo) *Rule {
		if len(info.CountryCode()) > 0 {
			if info.CountryCode() == rule.Value {
				return rule
			}
		} else {
			if info.IP() == nil {
				answer := dnsHandle(ctx, info.Domain())
				if answer == nil || len(answer.CurrentIP) == 0 {
					info.SetIP([]byte{})
					return next(ctx, info)
				}
				info.SetIP(answer.CurrentIP)
			}
			if len(info.IP()) > 0 {
				info.SetCountryCode(dns.GeoLookUp(info.IP()))
				if info.CountryCode() == rule.Value {
					return rule
				}
			}
			if len(info.CountryCode()) > 0 && info.CountryCode() == rule.Value {
				return rule
			}
		}
		return next(ctx, info)
	}, nil
}
