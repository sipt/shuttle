package model

type Config struct {
	Info struct {
		Name string
		URI  string
	} `toml:"-" yaml:"-" json:"-" yaml:"-"`

	General struct {
		LoggerLevel    string `toml:"logger-level" yaml:"logger-level"`
		DefaultTestURI string `toml:"default-test-uri" yaml:"default-test-uri"`
	} `toml:"general" yaml:"general"`

	Listener []struct {
		// Typ: eg. [http, https, socks]
		Typ string `toml:"typ" yaml:"typ"`
		// Addr: e.g. [":8080", "[::1]:8080", "192.168.1.23:8080"]
		Addr string `toml:"addr" yaml:"addr"`
		// Params of listener: e.g. {"auth-type": "basic", "password": "password", "user": "user name"}
		Params map[string]string `toml:"params" yaml:"params"`
	} `toml:"listener" yaml:"listener"`

	Controller struct {
		// Addr: e.g. [":8080", "[::1]:8080", "192.168.1.23:8080"]
		Addr string `toml:"addr" yaml:"addr"`
		// Params of listener: e.g. {"auth-type": "basic", "password": "password", "user": "user name"}
		Params map[string]string `toml:"params" yaml:"params"`
	} `toml:"controller" yaml:"controller"`

	DNS struct {
		IncludeSystem bool     `toml:"include-system" yaml:"include-system"`
		Servers       []string `toml:"servers" yaml:"servers"`
		TimeoutSec    int      `toml:"timeout-sec" yaml:"timeout-sec"`
		Mapping       []struct {
			Domain string   `toml:"domain" yaml:"domain"`
			IP     []string `toml:"ip" yaml:"ip"`
			Server []string `toml:"server" yaml:"server"`
		} `toml:"mapping" yaml:"mapping"`
	} `toml:"dns" yaml:"dns"`

	Server map[string]struct {
		Typ    string            `toml:"typ" yaml:"typ"`
		Host   string            `toml:"host" yaml:"host"`
		Port   int               `toml:"port" yaml:"port"`
		Params map[string]string `toml:"params" yaml:"params"`
	} `toml:"server" yaml:"server"`

	ServerGroup map[string]struct {
		// Typ: e.g. ["rtt", "select"]
		Typ string `toml:"typ" yaml:"typ"`
		// Servers: e.g. in {Server..., ServerGroup...}
		Servers []string          `toml:"servers" yaml:"servers"`
		Params  map[string]string `toml:"params" yaml:"params"`
	} `toml:"server-group" yaml:"server-group"`

	Rule []Rule `toml:"rule" yaml:"rule"`

	UDPRule []Rule `toml:"udp-rule" yaml:"udp-rule"`

	Include []struct {
		// e.g. {typ = local, path = "/User/root/config/server.toml"}
		Typ    string            `toml:"typ" yaml:"typ"`
		Params map[string]string `toml:"params" yaml:"params"`
	} `toml:"include" yaml:"include"`

	Filter []struct {
		Typ    string            `toml:"typ" yaml:"typ"`
		Params map[string]string `toml:"params" yaml:"params"`
	} `toml:"filter" yaml:"filter"`

	Stream struct {
		Before []struct {
			Typ    string            `toml:"typ" yaml:"typ"`
			Params map[string]string `toml:"params" yaml:"params"`
		} `toml:"before" yaml:"before"`
		After []struct {
			Typ    string            `toml:"typ" yaml:"typ"`
			Params map[string]string `toml:"params" yaml:"params"`
		} `toml:"after" yaml:"after"`
	} `toml:"stream" yaml:"stream"`

	Plugins map[string]map[string]string `json:"plugins"`
}

type Rule struct {
	Typ    string            `toml:"typ" yaml:"typ"`
	Value  string            `toml:"value" yaml:"value"`
	Proxy  string            `toml:"proxy" yaml:"proxy"`
	Params map[string]string `toml:"params" yaml:"params"`
}
