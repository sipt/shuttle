package shuttle

import (
	"fmt"
	"github.com/sipt/shuttle/log"
	"github.com/sipt/yaml"
	"io/ioutil"
	"net"
	"regexp"
	"strings"
)

const ControllerDomain = "c.sipt.top"
const ConfigFileVersion = "v1.0.1"

var controllerDomain string
var controllerPort string

type Config struct {
	Ver        string              `yaml:"ver"`
	General    *General            `yaml:"General"`
	Proxy      map[string][]string `yaml:"Proxy,[flow],2quoted"`
	ProxyGroup map[string][]string `yaml:"Proxy-Group,[flow],2quoted"`
	LocalDNSs  [][]string          `yaml:"Local-DNS,[flow],2quoted"`
	Mitm       *Mitm               `yaml:"MITM"`
	Rule       [][]string          `yaml:"Rule,[flow],2quoted"`
	HttpMap    *HttpMap            `yaml:"Http-Map"`
}

type General struct {
	LogLevel            string   `yaml:"loglevel,2quoted"`
	DNSServer           []string `yaml:"dns-server,2quoted"`
	HttpPort            string   `yaml:"http-port,2quoted"`
	HttpInterface       string   `yaml:"http-interface,2quoted"`
	SocksPort           string   `yaml:"socks-port,2quoted"`
	SocksInterface      string   `yaml:"socks-interface,2quoted"`
	ControllerPort      string   `yaml:"controller-port,2quoted"`
	ControllerInterface string   `yaml:"controller-interface,2quoted"`
}

type Mitm struct {
	CA    string   `yaml:"ca,2quoted"`
	Key   string   `yaml:"key,2quoted"`
	Rules []string `yaml:"rules,flow,2quoted"`
}

type HttpMap struct {
	ReqMap  []*ModifyMap `yaml:"Req-Map,2quoted"`
	RespMap []*ModifyMap `yaml:"Resp-Map,2quoted"`
}

type ModifyMap struct {
	Type   string     `yaml:"type,2quoted"`
	UrlRex string     `yaml:"url-rex,2quoted"`
	Items  [][]string `yaml:"items,[flow],2quoted"`
}

func ReloadConfig() (*General, error) {
	DestroyServers()
	ClearDNSCache()
	ClearRecords()
	ClearHttpModify()
	//
	return InitConfig(configFile)
}

func SetMimt(mitm *Mitm) {
	if len(mitm.CA) > 0 {
		conf.Mitm.CA = mitm.CA
	}
	if len(mitm.Key) > 0 {
		conf.Mitm.Key = mitm.Key
	}
	conf.Mitm.Rules = mitm.Rules
	SaveToFile()
}

func SaveToFile() {
	bytes, err := yaml.Marshal(conf)
	if err != nil {
		log.Logger.Errorf("[CONF] yaml marshal config failed : %v", err)
	}
	offset := EmojiDecode(bytes)
	bytes = bytes[:offset]
	err = ioutil.WriteFile(configFile, bytes, 0644)
	if err != nil {
		log.Logger.Errorf("[CONF] save config file failed : %v", err)
	}
}

var configFile string
var conf *Config

func GetGeneralConfig() (general General) {
	general = *conf.General
	return
}

func InitConfig(filePath string) (*General, error) {
	configFile = filePath
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read config file failed: %v", err)
	}
	conf = &Config{}
	err = yaml.Unmarshal(data, conf)
	if err != nil {
		return nil, fmt.Errorf("resolve config file failed: %v", err)
	}
	if conf.Ver != ConfigFileVersion {
		return nil, fmt.Errorf("resolve config file failed: only support ver:%s current:[%s]", ConfigFileVersion, conf.Ver)
	}
	//General
	//logger level
	log.Logger.SetLevel(log.LevelMap[conf.General.LogLevel])

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
			localDNS[i].IPs, err = parseIPs(v[3])
		case DNSTypeDirect:
			localDNS[i].DNSs, err = parseIPs(v[3])
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

	//Servers
	ss := make([]*Server, len(conf.Proxy)+2)
	index := 0
	ss[index] = &Server{Name: PolicyDirect} // 直连
	index ++
	ss[index] = &Server{Name: PolicyReject} // 拒绝
	for k, v := range conf.Proxy {
		index ++
		if len(v) < 2 {
			return nil, fmt.Errorf("resolve config file [proxy] [%s] failed", k)
		}
		ss[index], err = NewServer(k, v)
		if err != nil {
			return nil, err
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

	//Http Map
	var reqMaps, respMaps []*ModifyPolicy
	if conf.HttpMap != nil {
		if len(conf.HttpMap.ReqMap) > 0 {
			reqMaps = make([]*ModifyPolicy, len(conf.HttpMap.ReqMap))
			for i, v := range conf.HttpMap.ReqMap {
				reqMaps[i] = &ModifyPolicy{
					Type:   v.Type,
					UrlRex: v.UrlRex,
				}
				reqMaps[i].rex, err = regexp.Compile(v.UrlRex)
				if err != nil {
					return nil, fmt.Errorf("resolve config file [Http-Map] [%s] failed: %v", err)
				}
				if len(v.Items) > 0 {
					reqMaps[i].MVs = make([]*ModifyValue, len(v.Items))
					for j, e := range v.Items {
						if len(e) != 3 {
							return nil, fmt.Errorf("resolve config file [Http-Map] failed: %v, item's count must be 3", e)
						}
						reqMaps[i].MVs[j] = &ModifyValue{
							Type:  e[0],
							Key:   e[1],
							Value: e[2],
						}
					}
				}
			}
		}
		if len(conf.HttpMap.ReqMap) > 0 {
			respMaps = make([]*ModifyPolicy, len(conf.HttpMap.RespMap))
			for i, v := range conf.HttpMap.RespMap {
				respMaps[i] = &ModifyPolicy{
					Type:   v.Type,
					UrlRex: v.UrlRex,
				}
				respMaps[i].rex, err = regexp.Compile(v.UrlRex)
				if err != nil {
					return nil, fmt.Errorf("resolve config file [Http-Map] [%s] failed: %v", err)
				}
				if len(v.Items) > 0 {
					respMaps[i].MVs = make([]*ModifyValue, len(v.Items))
					for j, e := range v.Items {
						if len(e) != 3 {
							return nil, fmt.Errorf("resolve config file [Http-Map] failed: %v, item's count must be 3", e)
						}
						respMaps[i].MVs[j] = &ModifyValue{
							Type:  e[0],
							Key:   e[1],
							Value: e[2],
						}
					}
				}
			}
		}
	}
	InitHttpModify(reqMaps, respMaps)

	err = InitCert(conf.Mitm)
	if err != nil {
		return nil, fmt.Errorf("mitm init failed: %v", err)
	}
	if conf.Mitm != nil {
		SetMitMRules(conf.Mitm.Rules)
	}
	if len(controllerPort) == 0 {
		controllerPort = conf.General.ControllerPort
	} else {
		conf.General.ControllerPort = controllerPort
	}
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
func EmojiDecode(data []byte) int {
	index, length := 0, len(data)
	offset := 0
	for index < length {
		if data[index] == '\\' && data[index+1] == 'U' {
			index += 2
			decodeEmoji(data[offset:offset], data[index:index+8])
			offset += 4
			index += 8
		} else {
			if index != offset {
				data[offset] = data[index]
			}
			offset ++
			index ++
		}
	}
	return offset
}

func decodeEmoji(dst []byte, src []byte) (err error) {
	const code_length = 8
	var value int
	for k := 0; k < code_length; k++ {
		if !is_hex(src, k) {
			err = fmt.Errorf("is not hex :%v", src[k])
			return
		}
		value = (value << 4) + as_hex(src, k)
	}

	// Check the value and write the character.
	if (value >= 0xD800 && value <= 0xDFFF) || value > 0x10FFFF {
		err = fmt.Errorf("is not hex :%v", value)
		return
	}
	if value <= 0x7F {
		dst = append(dst, byte(value))
	} else if value <= 0x7FF {
		dst = append(dst, byte(0xC0+(value>>6)))
		dst = append(dst, byte(0x80+(value&0x3F)))
	} else if value <= 0xFFFF {
		dst = append(dst, byte(0xE0+(value>>12)))
		dst = append(dst, byte(0x80+((value>>6)&0x3F)))
		dst = append(dst, byte(0x80+(value&0x3F)))
	} else {
		dst = append(dst, byte(0xF0+(value>>18)))
		dst = append(dst, byte(0x80+((value>>12)&0x3F)))
		dst = append(dst, byte(0x80+((value>>6)&0x3F)))
		dst = append(dst, byte(0x80+(value&0x3F)))
	}
	return nil
}
func is_hex(b []byte, i int) bool {
	return b[i] >= '0' && b[i] <= '9' || b[i] >= 'A' && b[i] <= 'F' || b[i] >= 'a' && b[i] <= 'f'
}

// Get the value of a hex-digit.
func as_hex(b []byte, i int) int {
	bi := b[i]
	if bi >= 'A' && bi <= 'F' {
		return int(bi) - 'A' + 10
	}
	if bi >= 'a' && bi <= 'f' {
		return int(bi) - 'a' + 10
	}
	return int(bi) - '0'
}

var ShuttleVersion = "v0.5.0"
