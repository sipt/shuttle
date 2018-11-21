package config

import (
	"fmt"
	"github.com/sipt/yaml"
	"io/ioutil"
)

const ConfigFileVersion = "v1.0.1"

// load config file
func LoadConfig(filePath string) (*Config, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read config file failed: %v", err)
	}

	conf := &Config{}
	err = yaml.Unmarshal(data, conf)
	if err != nil {
		return nil, fmt.Errorf("resolve config file failed: %v", err)
	}
	if conf.Ver != ConfigFileVersion {
		return nil, fmt.Errorf("resolve config file failed: only support ver:%s current:[%s]", ConfigFileVersion, conf.Ver)
	}
	return conf, nil
}

// save config file
func SaveConfig(filePath string, config *Config) error {
	return nil
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
	ReqMap  []*ModifyMap `yaml:"Req-Map,2quoted"`
	RespMap []*ModifyMap `yaml:"Resp-Map,2quoted"`
}

type ModifyMap struct {
	Type   string     `yaml:"type,2quoted"`
	UrlRex string     `yaml:"url-rex,2quoted"`
	Items  [][]string `yaml:"items,[flow],2quoted"`
}
