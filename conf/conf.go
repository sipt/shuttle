package conf

type Config struct {
	General struct {
		LoggerLevel string
	}

	Listener []struct {
		// Typ: eg. [http, https, socks]
		Typ string
		// Addr: e.g. [":8080", "[::1]:8080", "192.168.1.23:8080"]
		Addr string
		// Params of listener: e.g. {"Password": "password", "UserName": "user name"}
		Params map[string]string
	}

	ProxyServer map[string]struct {
		Typ    string
		Addr   string
		Port   string
		Params map[string]string
	}

	ProxyServerGroup map[string]struct {
		// Typ: e.g. ["rtt", "select"]
		Typ string
		// Servers: e.g. in {ProxyServer..., ProxyServerGroup...}
		Servers []string
		Params  map[string]string
	}

	Rule []struct {
		Typ    string
		Value  string
		Params map[string]string
	}
}
