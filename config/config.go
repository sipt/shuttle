package config

import (
	"fmt"
	"github.com/sipt/shuttle/util"
	"github.com/sipt/yaml"
	"io/ioutil"
)

const ConfigFileVersion = "v1.0.1"
const SetAsSystemProxyAuto = "auto"

var configFile string
var conf *Config

func CurrentConfig() *Config {
	return conf
}

func CurrentConfigFile() string {
	return configFile
}

// load config file
func LoadConfig(filePath string) (*Config, error) {
	util.RLock(filePath)
	defer util.RUnLock(filePath)
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
	configFile = filePath
	return conf, nil
}

// save config file
func SaveConfig(configFile string, config *Config) error {
	util.Lock(configFile)
	defer util.UnLock(configFile)
	bytes, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("[CONF] yaml marshal config failed : %v", err)
	}
	offset := EmojiDecode(bytes)
	bytes = bytes[:offset]
	err = ioutil.WriteFile(configFile, bytes, 0644)
	if err != nil {
		return fmt.Errorf("[CONF] save config file failed : %v", err)
	}
	return nil
}

func ReloadConfig() (*Config, error) {
	if configFile == "" {
		return nil, fmt.Errorf("config file not found")
	}
	return LoadConfig(configFile)
}

// download config file
func downloadConfig(url string) error {
	return nil
}

type Config struct {
	Ver        string              `yaml:"ver"`
	General    *General            `yaml:"General"`
	Proxy      map[string][]string `yaml:"Proxy,[flow],2quoted"`
	ProxyGroup map[string][]string `yaml:"Proxy-Group,[flow],2quoted"`
	LocalDNSs  [][]string          `yaml:"Local-DNS,[flow],2quoted"`
	Mitm       *Mitm               `yaml:"MITM"`
	Rule       [][]string          `yaml:"Rule,[flow],2quoted"`
	HttpMap    *HttpMap            `yaml:"Http-Map"`
	RttUrl     string              `yaml:"rtt-url"`
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
	SetAsSystemProxy    string   `yaml:"set-as-system-proxy,2quoted"`
}

type Mitm struct {
	CA    string   `yaml:"ca,2quoted"`
	Key   string   `yaml:"key,2quoted"`
	Rules []string `yaml:"rules,flow,2quoted"`
}

type HttpMap struct {
	ReqMap  []*ModifyMap `yaml:"Req-Map,2quoted" json:"req_map"`
	RespMap []*ModifyMap `yaml:"Resp-Map,2quoted" json:"resp_map"`
}

type ModifyMap struct {
	Type   string     `yaml:"type,2quoted" json:"type"`
	UrlRex string     `yaml:"url-rex,2quoted" json:"url_rex"`
	Items  [][]string `yaml:"items,[flow],2quoted" json:"items"`
}

//dns
//func GetControllerDomain() string ==> controller
//func GetControllerPort() string ==> controller
func (c *Config) GetDNSServers() []string {
	return c.General.DNSServer
}
func (c *Config) SetDNSServers(servers []string) {
	c.General.DNSServer = servers
}
func (c *Config) GetLocalDNS() [][]string {
	return c.LocalDNSs
}
func (c *Config) SetLocalDNS(localDNSs [][]string) {
	c.LocalDNSs = localDNSs
}
func (c *Config) GetGeoIPDBFile() string {
	return "GeoLite2-Country.mmdb"
}

//logger
func (c *Config) GetLogLevel() string {
	return c.General.LogLevel
}
func (c *Config) SetLogLevel(l string) {
	c.General.LogLevel = l
}

//controller
func (c *Config) GetControllerDomain() string {
	return "c.sipt.top"
}
func (c *Config) GetControllerInterface() string {
	return c.General.ControllerInterface
}
func (c *Config) SetControllerInterface(inter string) {
	c.General.ControllerInterface = inter
}
func (c *Config) GetControllerPort() string {
	return c.General.ControllerPort
}
func (c *Config) SetControllerPort(port string) {
	c.General.ControllerPort = port
}

//HTTP Proxy
func (c *Config) GetHTTPInterface() string {
	return c.General.HttpInterface
}
func (c *Config) SetHTTPInterface(inter string) {
	c.General.HttpInterface = inter
}
func (c *Config) GetHTTPPort() string {
	return c.General.HttpPort
}
func (c *Config) SetHTTPPort(port string) {
	c.General.HttpPort = port
}

//SOCKS Proxy
func (c *Config) GetSOCKSInterface() string {
	return c.General.SocksInterface
}
func (c *Config) SetSOCKSInterface(inter string) {
	c.General.SocksInterface = inter
}
func (c *Config) GetSOCKSPort() string {
	return c.General.SocksPort
}
func (c *Config) SetSOCKSPort(port string) {
	c.General.SocksPort = port
}

//Proxy & Proxy Group
func (c *Config) GetProxy() map[string][]string {
	return c.Proxy
}
func (c *Config) SetProxy(proxy map[string][]string) {
	c.Proxy = proxy
}
func (c *Config) GetProxyGroup() map[string][]string {
	return c.ProxyGroup
}
func (c *Config) SetProxyGroup(group map[string][]string) {
	c.ProxyGroup = group
}
func (c *Config) SetRttUrl(rttUrl string) {
	c.RttUrl = rttUrl
}
func (c *Config) GetRttUrl() string {
	return c.RttUrl
}

//Rule
func (c *Config) GetRule() [][]string {
	return c.Rule
}
func (c *Config) SetRule(rule [][]string) {
	c.Rule = rule
}

//HttpMap
func (c *Config) GetHTTPMap() *HttpMap {
	return c.HttpMap
}
func (c *Config) SetHTTPMap(httpMap *HttpMap) {
	c.HttpMap = httpMap
}

//MITM
func (c *Config) GetMITM() *Mitm {
	return c.Mitm
}
func (c *Config) SetMITM(mitm *Mitm) {
	c.Mitm = mitm
}
