package shuttle

import (
	"fmt"
	"net"
	"strings"
)

const (
	PolicyReject = "REJECT"
	PolicyDirect = "DIRECT"
	PolicyGlobal = "GLOBAL"
	PolicyMock   = "MOCK"
	PolicyNone   = "NONE"

	RuleDomainSuffix  = "DOMAIN-SUFFIX"
	RuleDomain        = "DOMAIN"
	RuleDomainKeyword = "DOMAIN-KEYWORD"
	RuleGeoIP         = "GEOIP"
	RuleFinal         = "FINAL"
	RuleIPCIDR        = "IP-CIDR"

	ConnModeDirect = "DIRECT"
	ConnModeRemote = "REMOTE"
	ConnModeRule   = "RULE"
	ConnModeReject = "REJECT"
)

var rules []*Rule
var connMode = ConnModeRule

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

func SetConnMode(mode string) error {
	switch connMode {
	case ConnModeDirect, ConnModeRemote, ConnModeRule, ConnModeReject:
		connMode = mode
		return nil
	default:
		return nil
	}
}

func GetConnMode() string {
	return connMode
}

type Rule struct {
	Type    string
	Value   string
	Policy  string
	Options []string
	Comment string
}

func filter(req *Request) (*Rule, error) {
	switch connMode {
	case ConnModeDirect:
		return directRule, nil
	case ConnModeRemote:
		return remoteRule, nil
	case ConnModeReject:
		return rejectRule, nil
	}

	for _, v := range rules {
		switch v.Type {
		case RuleDomainSuffix:
			if req.Addr == v.Value || strings.HasSuffix(req.Addr, "."+v.Value) {
				return v, nil
			}
		case RuleDomain:
			if req.Addr == v.Value {
				return v, nil
			}
		case RuleDomainKeyword:
			if strings.Index(req.Addr, v.Value) >= 0 {
				return v, nil
			}
		case RuleIPCIDR:
			if ipCidrMap[v.Value].Contains(req.IP) {
				return v, nil
			}
		case RuleGeoIP:
			if v.Value == req.DomainHost.Country {
				return v, nil
			}
		case RuleFinal:
			return v, nil
		}
	}
	return nil, nil
}

var directRule = &Rule{
	Type:   "GLOBAL",
	Policy: PolicyDirect,
}
var remoteRule = &Rule{
	Type:   "GLOBAL",
	Policy: PolicyGlobal,
}
var rejectRule = &Rule{
	Type:   "GLOBAL",
	Policy: PolicyReject,
}
var mockRule = &Rule{
	Type:   "MOCK",
	Policy: PolicyMock,
}
var mockServer = &Server{
	Name: "MOCK",
}
