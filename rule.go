package shuttle

import (
	"strings"
	"net"
	"fmt"
)

const (
	PolicyReject = "REJECT"
	PolicyDirect = "DIRECT"
	PolicyNone   = "NONE"

	RuleDomainSuffix  = "DOMAIN-SUFFIX"
	RuleDomain        = "DOMAIN"
	RuleDomainKeyword = "DOMAIN-KEYWORD"
	RuleGeoIP         = "GEOIP"
	RuleFinal         = "FINAL"
	RuleIPCIDR        = "IP-CIDR"
)

var rules []*Rule

var ipCidrMap map[string]*net.IPNet

func InitRule(rs []*Rule) error {
	rules = rs
	ipCidrMap = make(map[string]*net.IPNet)
	for _, v := range rs {
		if v.Type == RuleIPCIDR {
			_, ipNet, err := net.ParseCIDR(v.Value)
			if err != nil {
				return fmt.Errorf("[Rule] [IP-CIDR] [%s] error: %v", v.Value, err)
			}
			ipCidrMap[v.Value] = ipNet
		}
	}
	return nil
}

type Rule struct {
	Type    string
	Value   string
	Policy  string
	Options []string
	Comment string
}

func Filter(req *Request) (*Rule, error) {
	for i, v := range rules {
		switch v.Type {
		case RuleDomainSuffix:
			if strings.HasSuffix(req.Addr, v.Value) {
				return rules[i], nil
			}
		case RuleDomain:
			if req.Addr == v.Value {
				return rules[i], nil
			}
		case RuleDomainKeyword:
			if strings.Index(req.Addr, v.Value) >= 0 {
				return rules[i], nil
			}
		case RuleIPCIDR:
			if ipCidrMap[v.Value].Contains(req.IP) {
				return rules[i], nil
			}
		case RuleGeoIP:
			if v.Value == req.DomainHost.Country {
				return rules[i], nil
			}
		case RuleFinal:
			return rules[i], nil
		}
	}
	return nil, nil
}
