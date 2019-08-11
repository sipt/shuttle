package conf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type M map[string]string

func TestTomlMarshal_Marshal(t *testing.T) {
	config := &Config{
		General: struct {
			LoggerLevel string
		}{
			"debug",
		},

		Listener: []struct {
			Typ    string
			Addr   string
			Params map[string]string
		}{
			{"https", ":8081", map[string]string{"user": "root", "password": "123123"}},
			{"socks", ":8080", map[string]string{"user": "root", "password": "123123"}},
		},

		ProxyServer: map[string]struct {
			Typ    string
			Addr   string
			Port   string
			Params map[string]string
		}{
			"JP1": {"ss", "jp.remote.com", "8080", M{"user": "root", "password": "123123"}},
			"JP2": {"ss", "jp.remote.com", "8080", M{"user": "root", "password": "123123"}},
			"US1": {"ss", "us.remote.com", "8080", M{"user": "root", "password": "123123"}},
			"US2": {"ss", "us.remote.com", "8080", M{"user": "root", "password": "123123"}},
		},

		ProxyServerGroup: map[string]struct {
			Typ     string
			Servers []string
			Params  map[string]string
		}{
			"Proxy": {"select", []string{"AUTO", "JP", "US"}, nil},
			"AUTO":  {"rtt", []string{"JP1", "JP2", "US1", "US2"}, M{"url": "https://www.google.com"}},
			"JP":    {"select", []string{"JP1", "JP2"}, nil},
			"US":    {"select", []string{"US1", "US2"}, nil},
		},

		Rule: []struct {
			Typ    string
			Value  string
			Params map[string]string
		}{
			{"DOMAIN", "google.com", M{"Proxy": "Proxy", "Comment": "search engine"}},
			{"DOMAIN", "github.com", M{"Proxy": "Proxy", "Comment": "source code"}},
		},
	}
	m, _ := newTomlMarshal(nil)
	data, err := m.Marshal(config)
	assert.NoError(t, err)

	config2, err := m.UnMarshal(data)
	assert.EqualValues(t, config, config2)
}
