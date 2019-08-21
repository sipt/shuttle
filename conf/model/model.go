package model

type Config struct {
	General struct {
		LoggerLevel string `toml:"logger_level"`
	} `toml:"general"`

	Listener []struct {
		Name string `toml:"name"`
		// Typ: eg. [http, https, socks]
		Typ string `toml:"typ"`
		// Addr: e.g. [":8080", "[::1]:8080", "192.168.1.23:8080"]
		Addr string `toml:"addr"`
		// Params of listener: e.g. {"Password": "password", "UserName": "user name"}
		Params map[string]string `toml:"params"`
	} `toml:"listener"`

	DNS struct {
		IncludeSystem bool     `toml:"include_system"`
		Servers       []string `toml:"servers"`
		TimeoutSec    int      `toml:"timeout_sec"`
		Mapping       []struct {
			Domain string   `toml:"domain"`
			IP     []string `toml:"ip"`
			Server []string `toml:"server"`
		} `toml:"mapping"`
	} `toml:"dns"`

	Server map[string]struct {
		Typ    string            `toml:"typ"`
		Addr   string            `toml:"addr"`
		Port   string            `toml:"port"`
		Params map[string]string `toml:"params"`
	} `toml:"server"`

	ServerGroup map[string]struct {
		// Typ: e.g. ["rtt", "select"]
		Typ string `toml:"typ"`
		// Servers: e.g. in {Server..., ServerGroup...}
		Servers []string          `toml:"servers"`
		Params  map[string]string `toml:"params"`
	} `toml:"server_group"`

	Rule []struct {
		Typ    string            `toml:"typ"`
		Value  string            `toml:"value"`
		Proxy  string            `toml:"proxy"`
		Params map[string]string `toml:"params"`
	} `toml:"rule"`

	Include []struct {
		// e.g. {typ = local, path = "/User/root/config/server.toml"}
		Typ    string            `toml:"typ"`
		Params map[string]string `toml:"params"`
	} `toml:"include"`
}
