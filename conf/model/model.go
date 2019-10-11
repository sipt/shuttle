package model

type Config struct {
	Info struct {
		Name string
		URI  string
	} `toml:"-" json:"-" yaml:"-"`

	General struct {
		LoggerLevel    string `toml:"logger-level" json:"logger-level" yaml:"logger-level"`
		DefaultTestURI string `toml:"default-test-uri" json:"default-test-uri" yaml:"default-test-uri"`
	} `toml:"general" json:"general" yaml:"general"`

	Listener []struct {
		Name string `toml:"name" json:"name" yaml:"name"`
		// Typ: eg. [http, https, socks]
		Typ string `toml:"typ" json:"typ" yaml:"typ"`
		// Addr: e.g. [":8080", "[::1]:8080", "192.168.1.23:8080"]
		Addr string `toml:"addr" json:"addr" yaml:"addr"`
		// Params of listener: e.g. {"auth-type": "basic", "password": "password", "user": "user name"}
		Params map[string]string `toml:"params" json:"params" yaml:"params"`
	} `toml:"listener" json:"listener" yaml:"listener"`

	Controller struct {
		// Addr: e.g. [":8080", "[::1]:8080", "192.168.1.23:8080"]
		Addr string `toml:"addr" json:"addr" yaml:"addr"`
		// Params of listener: e.g. {"auth-type": "basic", "password": "password", "user": "user name"}
		Params map[string]string `toml:"params" json:"params" yaml:"params"`
	} `toml:"controller" json:"controller" yaml:"controller"`

	DNS struct {
		IncludeSystem bool     `toml:"include-system" json:"include-system" yaml:"include-system"`
		Servers       []string `toml:"servers" json:"servers" yaml:"servers"`
		TimeoutSec    int      `toml:"timeout-sec" json:"timeout-sec" yaml:"timeout-sec"`
		Mapping       []struct {
			Domain string   `toml:"domain" json:"domain" yaml:"domain"`
			IP     []string `toml:"ip" json:"ip" yaml:"ip"`
			Server []string `toml:"server" json:"server" yaml:"server"`
		} `toml:"mapping" json:"mapping" yaml:"mapping"`
	} `toml:"dns" json:"dns" yaml:"dns"`

	Server map[string]struct {
		Typ    string            `toml:"typ" json:"typ" yaml:"typ"`
		Host   string            `toml:"host" json:"host" yaml:"host"`
		Port   int               `toml:"port" json:"port" yaml:"port"`
		Params map[string]string `toml:"params" json:"params" yaml:"params"`
	} `toml:"server" json:"server" yaml:"server"`

	ServerGroup map[string]struct {
		// Typ: e.g. ["rtt", "select"]
		Typ string `toml:"typ" json:"typ" yaml:"typ"`
		// Servers: e.g. in {Server..., ServerGroup...}
		Servers []string          `toml:"servers" json:"servers" yaml:"servers"`
		Params  map[string]string `toml:"params" json:"params" yaml:"params"`
	} `toml:"server-group" json:"server-group" yaml:"server-group"`

	Rule []Rule `toml:"rule" json:"rule" yaml:"rule"`

	UDPRule []Rule `toml:"udp-rule" json:"udp-rule" yaml:"udp-rule"`

	Include []struct {
		// e.g. {typ = local, path = "/User/root/config/server.toml"}
		Typ    string            `toml:"typ" json:"typ" yaml:"typ"`
		Params map[string]string `toml:"params" json:"params" yaml:"params"`
	} `toml:"include" json:"include" yaml:"include"`

	Filter []struct {
		Typ    string            `toml:"typ" json:"typ" yaml:"typ"`
		Params map[string]string `toml:"params" json:"params" yaml:"params"`
	} `toml:"filter" json:"filter" yaml:"filter"`

	Stream struct {
		Before []struct {
			Typ    string            `toml:"typ" json:"typ" yaml:"typ"`
			Params map[string]string `toml:"params" json:"params" yaml:"params"`
		} `toml:"before" json:"before" yaml:"before"`
		After []struct {
			Typ    string            `toml:"typ" json:"typ" yaml:"typ"`
			Params map[string]string `toml:"params" json:"params" yaml:"params"`
		} `toml:"after" json:"after" yaml:"after"`
	} `toml:"stream" json:"stream" yaml:"stream"`

	Plugins map[string]map[string]string `json:"plugins"`
}

type Rule struct {
	Typ    string            `toml:"typ" json:"typ" yaml:"typ"`
	Value  string            `toml:"value" json:"value" yaml:"value"`
	Proxy  string            `toml:"proxy" json:"proxy" yaml:"proxy"`
	Params map[string]string `toml:"params" json:"params" yaml:"params"`
}
