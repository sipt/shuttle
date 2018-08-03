package shuttle

import (
	"net"
	"io/ioutil"
	"fmt"
	"gopkg.in/yaml.v2"
	"strings"
)

type Config struct {
	Ver        string              `yaml:"ver"`
	General    *General            `yaml:"general"`
	Proxy      map[string][]string `yaml:"proxy"`
	ProxyGroup map[string][]string `yaml:"proxy-group"`
	LocalDNSs  [][]string          `yaml:"local-dns"`
	Rule       [][]string          `yaml:"rule"`
}

type General struct {
	LogLevel       string   `yaml:"loglevel"`
	DNSServer      []string `yaml:"dns-server"`
	SocksPort      string   `yaml:"socks-port"`
	HttpPort       string   `yaml:"http-port"`
	HttpInterface  string   `yaml:"http-interface"`
	SocksInterface string   `yaml:"socks-interface"`
}

func InitConfig(filepath string) (*General, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("read config file failed: %v", err)
	}
	conf := &Config{}
	err = yaml.Unmarshal(data, conf)
	if err != nil {
		return nil, fmt.Errorf("resolve config file failed: %v", err)
	}
	if conf.Ver != "v1.0.0" {
		return nil, fmt.Errorf("resolve config file failed: only support ver:v1.0.0 current:[%s]", conf.Ver)
	}
	//Servers
	ss := make([]*Server, len(conf.Proxy)+2)
	index := 0
	ss[index] = &Server{Name: PolicyDirect} // 直连
	index ++
	ss[index] = &Server{Name: PolicyReject} // 拒绝
	for k, v := range conf.Proxy {
		index ++
		if len(v) != 4 {
			return nil, fmt.Errorf("resolve config file [proxy] [%s] failed", k)
		}
		ss[index] = &Server{
			Name:     k,
			Host:     v[0],
			Port:     v[1],
			Method:   v[2],
			Password: v[3],
		}
	}
	gs := make([]*ServerGroup, len(conf.ProxyGroup))
	index = 0
	for k := range conf.ProxyGroup {
		gs[index] = &ServerGroup{Name: k}
		index ++
	}
	getServer := func(name string) interface{} {
		for i := range ss {
			if ss[i].Name == name {
				return ss[i]
			}
		}
		for i := range gs {
			if gs[i].Name == name {
				return gs[i]
			}
		}
		return nil
	}
	var cs []string
	for _, v := range gs {
		cs = conf.ProxyGroup[v.Name]
		if len(cs) < 2 {
			return nil, fmt.Errorf("resolve config file [proxy_group] [%s] failed", v.Name)
		}
		v.SelectType = cs[0]
		v.Servers = make([]interface{}, len(cs)-1)
		for i := range v.Servers {
			v.Servers[i] = getServer(cs[i+1])
			if v.Servers[i] == nil {
				return nil, fmt.Errorf("resolve config file [proxy_group] [%s] [%s] not found", v.Name, cs[i+1])
			}
		}
	}
	err = InitServers(gs, ss)
	if err != nil {
		return nil, fmt.Errorf("init server failed: %v", err)
	}

	//Rule
	rules := make([]*Rule, len(conf.Rule))
	for i, v := range conf.Rule {
		if len(v) != 4 {
			return nil, fmt.Errorf("resolve config file [rule] %v length must be 4", v)
		}
		rules[i] = &Rule{
			Type:    v[0],
			Value:   v[1],
			Policy:  v[2],
			Comment: v[3],
		}
		//switch v[0] {
		//case RuleDomainSuffix, RuleDomain, RuleDomainKeyword, RuleGeoIP, RuleFinal:
		//default:
		//	return nil, fmt.Errorf("resolve config file [rule] not support rule type [%v]", v[0])
		//}
		if getServer(v[2]) == nil {
			return nil, fmt.Errorf("resolve config file [rule] not support policy[%s]", v[2])
		}
	}
	err = InitRule(rules)
	if err != nil {
		return nil, fmt.Errorf("init rule failed: %v", err)
	}

	//DNS
	dns := make([]net.IP, len(conf.General.DNSServer))
	for i, v := range conf.General.DNSServer {
		dns[i] = net.ParseIP(v)
		if dns[i] == nil {
			return nil, fmt.Errorf("resolve config file [general.dns-server] not support [%s]", v)
		}
	}
	localDNS := make([]*DNS, len(conf.LocalDNSs))
	for i, v := range conf.LocalDNSs {
		if len(v) != 4 {
			return nil, fmt.Errorf("resolve config file [host] %v length must be 4", v)
		}
		localDNS[i] = &DNS{
			MatchType: v[0],
			Domain:    v[1],
			Type:      v[2],
		}
		if v[0] != RuleDomain && v[0] != RuleDomainSuffix && v[0] != RuleDomainKeyword {
			return nil, fmt.Errorf("resolve config file [host] not support rule type [%v]", v[0])
		}
		switch v[2] {
		case DNSTypeStatic:
			localDNS[0].IPs, err = parseIPs(v[3])
		case DNSTypeDirect:
			localDNS[0].DNSs, err = parseIPs(v[3])
		case DNSTypeRemote:
		default:
			return nil, fmt.Errorf("resolve config file [host] not support DNSType [%s]", v[1])
		}
		if err != nil {
			return nil, fmt.Errorf("resolve config file [host] [%v] [%v]", v, err)
		}
	}

	err = InitDNS(dns, localDNS)
	if err != nil {
		return nil, fmt.Errorf("init rule failed: %v", err)
	}

	//logger level
	SetLeve(conf.General.LogLevel)
	return conf.General, nil
}

func parseIPs(line string) ([]net.IP, error) {
	ips := strings.Split(line, ",")
	if len(ips) == 0 {
		return nil, nil
	}
	reply := make([]net.IP, len(ips))
	for i, v := range ips {
		reply[i] = net.ParseIP(v)
		if reply[i] == nil {
			return nil, fmt.Errorf("not resove ips [%s]", line)
		}
	}
	return reply, nil
}
