package marshal

import (
	"fmt"
	"testing"

	"github.com/sipt/shuttle/conf/model"
	"github.com/stretchr/testify/assert"
)

type M map[string]string

var config = &model.Config{
	General: struct {
		LoggerLevel    string `toml:"logger-level" json:"logger-level" yaml:"logger-level"`
		DefaultTestURI string `toml:"default-test-uri" json:"default-test-uri" yaml:"default-test-uri"`
	}{
		"debug", "https://www.bing.com",
	},

	Listener: []struct {
		Name   string            `toml:"name" json:"name" yaml:"name"`
		Typ    string            `toml:"typ" json:"typ" yaml:"typ"`
		Addr   string            `toml:"addr" json:"addr" yaml:"addr"`
		Params map[string]string `toml:"params" json:"params" yaml:"params"`
	}{
		{"http-inbound", "https", ":8081", map[string]string{"user": "root", "password": "123123"}},
		{"socks-inbound", "socks", ":8080", map[string]string{"user": "root", "password": "123123"}},
	},

	Server: map[string]struct {
		Typ    string            `toml:"typ" json:"typ" yaml:"typ"`
		Host   string            `toml:"host" json:"host" yaml:"host"`
		Port   int               `toml:"port" json:"port" yaml:"port"`
		Params map[string]string `toml:"params" json:"params" yaml:"params"`
	}{
		"JP1": {"ss", "jp.remote.com", 8080, M{"user": "root", "password": "123123"}},
		"JP2": {"ss", "jp.remote.com", 8080, M{"user": "root", "password": "123123"}},
		"US1": {"ss", "us.remote.com", 8080, M{"user": "root", "password": "123123"}},
		"US2": {"ss", "us.remote.com", 8080, M{"user": "root", "password": "123123"}},
	},

	ServerGroup: map[string]struct {
		// Typ: e.g. ["rtt", "select"]
		Typ string `toml:"typ" json:"typ" yaml:"typ"`
		// Servers: e.g. in {Server..., ServerGroup...}
		Servers []string          `toml:"servers" json:"servers" yaml:"servers"`
		Params  map[string]string `toml:"params" json:"params" yaml:"params"`
	}{
		"Proxy": {"select", []string{"AUTO", "JP", "US"}, nil},
		"AUTO":  {"rtt", []string{"JP1", "JP2", "US1", "US2"}, M{"url": "https://www.google.com"}},
		"JP":    {"select", []string{"JP1", "JP2"}, nil},
		"US":    {"select", []string{"US1", "US2"}, nil},
	},

	Rule: []model.Rule{
		{"DOMAIN", "google.com", "Proxy", M{"Comment": "search engine"}},
		{"DOMAIN", "github.com", "Proxy", M{"Comment": "source code"}},
	},
}

func TestTomlMarshal_Marshal(t *testing.T) {

	m, _ := newTomlMarshal(nil)
	data, err := m.Marshal(config)
	assert.NoError(t, err)

	fmt.Println(string(data))

	str := `
[general]
  logger_level = "debug"

[[listener]]
  typ = "https"
  addr = ":8081"
  [listener.params]
    password = "123123"
    user = "root"

[[listener]]
  typ = "socks"
  addr = ":8080"
  [listener.params]
    password = "123123"
    user = "root"

[server]
  JP1 = {typ = "ss", addr = "jp.remote.com", port = "8080", params = { password = "123123", user = "root" }}
  JP2 = {typ = "ss", addr = "jp.remote.com", port = "8080", params = { password = "123123", user = "root" }}
  US1 = {typ = "ss", addr = "us.remote.com", port = "8080", params = { password = "123123", user = "root" }}
  US2 = {typ = "ss", addr = "us.remote.com", port = "8080", params = { password = "123123", user = "root" }}

[server_group]
  [proxy_server_group.AUTO]
    typ = "rtt"
    servers = ["JP1", "JP2", "US1", "US2"]
    [proxy_server_group.AUTO.params]
      url = "https://www.google.com"
  [proxy_server_group.JP]
    typ = "select"
    servers = ["JP1", "JP2"]
  [proxy_server_group.Proxy]
    typ = "select"
    servers = ["AUTO", "JP", "US"]
  [proxy_server_group.US]
    typ = "select"
    servers = ["US1", "US2"]

[[rule]]
  typ = "DOMAIN"
  value = "google.com"
  proxy = "Proxy"
  [rule.params]
    Comment = "search engine"

[[rule]]
  typ = "DOMAIN"
  value = "github.com"
  proxy = "Proxy"
  [rule.params]
    Comment = "source code"
`

	config2, err := m.UnMarshal([]byte(str))
	assert.EqualValues(t, config, config2)
}
