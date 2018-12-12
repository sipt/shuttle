package dns

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strings"
)

const (
	DNSTypeStatic = "static"
	DNSTypeDirect = "direct"
	DNSTypeRemote = "remote"
	DNSTypeCache  = "cache"

	MatchTypeDomainSuffix  = "DOMAIN-SUFFIX"
	MatchTypeDomain        = "DOMAIN"
	MatchTypeDomainKeyword = "DOMAIN-KEYWORD"
	MatchNone              = "NONE"
)

type IDNSConfig interface {
	GetDNSServers() []string
	SetDNSServers([]string)
	GetLocalDNS() [][]string
	SetLocalDNS([][]string)

	GetControllerDomain() string
	GetControllerPort() string

	GetGeoIPDBFile() string
}

type DNS struct {
	MatchType string
	Domain    string
	IPs       []string
	DNSs      []string
	Port      string
	Type      string
	Country   string
}

func (d *DNS) String() string {
	buffer := bytes.NewBufferString(d.Domain)
	buffer.WriteString(" IPs:[")
	for i, v := range d.IPs {
		if i > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(v)
	}
	buffer.WriteString("] DNSs:[")
	for i, v := range d.DNSs {
		if i > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(v)
	}
	buffer.WriteString("] Country:[" + d.Country + "]")
	return buffer.String()
}

type DNSConfig struct {
	servers  []string
	localDNS []*DNS
}

var dnsConfig *DNSConfig

func ApplyConfig(config IDNSConfig) error {
	dnsConfig = &DNSConfig{}
	//DNS servers
	servers := config.GetDNSServers()
	if len(servers) == 0 {
		return errors.New("[DNS] [InitDNS] servers is empty")
	}
	dnsConfig.servers = make([]string, len(servers))
	for i, s := range servers {
		dnsConfig.servers[i] = s
		if net.ParseIP(s) == nil {
			return fmt.Errorf("[DNS] [InitDNS] %s is not a IP address", s)
		}
	}
	//Geo IP
	err := InitGeoIP(config.GetGeoIPDBFile())
	if err != nil {
		return err
	}

	//Local DNS
	inputs := config.GetLocalDNS()
	localDNS := make([]*DNS, len(inputs)+1)
	localDNS[0] = &DNS{
		MatchType: MatchTypeDomain,
		Domain:    config.GetControllerDomain(),
		Type:      DNSTypeStatic,
		Port:      config.GetControllerPort(),
		IPs:       []string{"127.0.0.1"},
	}
	var i int
	for j, v := range inputs {
		i = j + 1
		if len(v) != 4 {
			return fmt.Errorf("resolve config file [host] %v length must be 4", v)
		}
		localDNS[i] = &DNS{
			MatchType: v[0],
			Domain:    v[1],
			Type:      v[2],
		}
		if v[0] != MatchTypeDomain && v[0] != MatchTypeDomainSuffix && v[0] != MatchTypeDomainKeyword {
			return fmt.Errorf("resolve config file [host] not support rule type [%v]", v[0])
		}
		switch v[2] {
		case DNSTypeStatic:
			localDNS[i].IPs = strings.Split(v[3], ",")
		case DNSTypeDirect:
			localDNS[i].DNSs = strings.Split(v[3], ",")
		case DNSTypeRemote:
		default:
			return fmt.Errorf("resolve config file [host] not support DNSType [%s]", v[1])
		}
	}
	dnsConfig.localDNS = localDNS
	InitDNSCache()
	return nil
}
