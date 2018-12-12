package rule

import (
	"fmt"
	"github.com/sipt/shuttle/dns"
	"github.com/sipt/shuttle/proxy"
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

	OptionTunMode = "tun-mode"
)

var (
	rules    []*Rule
	connMode = ConnModeRule

	FailedRule = &Rule{Type: "FAILED", Policy: "FAILED"}
	DirectRule = &Rule{Type: "GLOBAL", Policy: PolicyDirect}
	RemoteRule = &Rule{Type: "GLOBAL", Policy: PolicyGlobal}
	RejectRule = &Rule{Type: "GLOBAL", Policy: PolicyReject}
	MockRule   = &Rule{Type: "MOCK", Policy: PolicyMock}
)

var ipCidrMap map[string]*net.IPNet

type IRuleConfig interface {
	GetRule() [][]string
	SetRule([][]string)
}

func ApplyConfig(config IRuleConfig) error {
	rs := make([]*Rule, len(config.GetRule()))
	ipCidrMap = make(map[string]*net.IPNet, 16)
	for i, v := range config.GetRule() {
		if len(v) != 4 {
			return fmt.Errorf("resolve config file [rule] %v length must be 4", v)
		}
		rs[i] = &Rule{
			Type:    v[0],
			Value:   v[1],
			Policy:  v[2],
			Comment: v[3],
		}
		if _, err := proxy.GetServer(v[2]); err != nil {
			return fmt.Errorf("resolve config file [rule] not support policy[%s]", v[2])
		}
		if v[0] == RuleIPCIDR {
			_, ipNet, err := net.ParseCIDR(v[1])
			if err != nil {
				return fmt.Errorf("[Rule] [IP-CIDR] [%s] error: %v", v[1], err)
			}
			ipCidrMap[v[1]] = ipNet
		}
	}
	rules = rs
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

type IRequest interface {
	Network() string
	Domain() string
	IP() string
	Port() string
	Answer() *dns.Answer
}

type Rule struct {
	Type    string
	Value   string
	Policy  string
	Options []string
	Comment string
}

func RuleFilter(req IRequest) (*Rule, error) {
	switch connMode {
	case ConnModeDirect:
		return DirectRule, nil
	case ConnModeRemote:
		return RemoteRule, nil
	case ConnModeReject:
		return RejectRule, nil
	}

	for _, v := range rules {
		switch v.Type {
		case RuleDomainSuffix:
			if req.Domain() == v.Value || strings.HasSuffix(req.Domain(), "."+v.Value) {
				return v, nil
			}
		case RuleDomain:
			if req.Domain() == v.Value {
				return v, nil
			}
		case RuleDomainKeyword:
			if strings.Index(req.Domain(), v.Value) >= 0 {
				return v, nil
			}
		case RuleIPCIDR:
			if len(req.IP()) > 0 && ipCidrMap[v.Value].Contains(net.ParseIP(req.IP())) {
				fmt.Println(v.Value, ":", req.IP(), ipCidrMap[v.Value].Contains(net.ParseIP(req.IP())))
				return v, nil
			}
		case RuleGeoIP:
			if req.Answer() != nil && v.Value == req.Answer().Country {
				return v, nil
			}
		case RuleFinal:
			return v, nil
		}
	}
	return nil, nil
}
