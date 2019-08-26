package conf

import (
	"context"
	"testing"

	"github.com/sipt/shuttle/conf/storage"

	"github.com/stretchr/testify/assert"
)

var dataMap = map[string][]byte{
	"shuttle.toml": []byte(`
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
  AUTO = {typ = "rtt", servers = ["JP1", "JP2", "US1", "US2"], params={url="https://www.google.com"}}
  JP = {typ = "select", servers = ["JP1", "JP2"]}
  Proxy = {typ = "select", servers = ["AUTO", "JP", "US"]}
  US = {typ = "select", servers = ["US1", "US2"]}

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

[[include]]
  typ = "map"
  [include.params]
    path = "include_1.toml"
[[include]]
  typ = "map"
  [include.params]
    path = "include_2.toml"`),
	"include_1.toml": []byte(`
[server_group.Other]
    typ = "select"
    servers = ["DIRECT"]
[server.UK1]
  typ = "ss"
  addr = "us.remote.com"
  port = "8080"
  params = { password = "123123", user = "root" }
[server.UK2]
  typ = "ss"
  addr = "us.remote.com"
  port = "8080"
  params = { password = "123123", user = "root" }`),
	"include_2.toml": []byte(`
[[rule]]
  typ = "DOMAIN-SUFFIX"
  value = "bing.com"
  proxy = "Proxy"
  [rule.params]
    Comment = "search engine"

[[rule]]
  typ = "DOMAIN-SUFFIX"
  value = "qq.com"
  proxy = "Proxy"
  [rule.params]
    Comment = "source code"`),
}

func init() {
	storage.Register("map", NewLocalInclude)
}

func NewLocalInclude(params map[string]string) (storage.IStorage, error) {
	return &localInclude{
		data: dataMap[params["path"]],
		name: params["path"],
	}, nil
}

type localInclude struct {
	data []byte
	name string
}

func (l *localInclude) Name() string {
	return l.name
}

func (l *localInclude) Load() ([]byte, error) {
	return l.data, nil
}
func (l *localInclude) RegisterNotify(ctx context.Context, notify func()) error {
	return nil
}
func (l *localInclude) Save(data []byte) error {
	return nil
}

func TestLoadConfig(t *testing.T) {
	ctx := context.Background()
	config, err := LoadConfig(ctx, "map", "toml", map[string]string{
		"path": "shuttle.toml",
	}, nil)
	assert.NoError(t, err)
	assert.EqualValues(t, len(config.Rule), 4)
	assert.EqualValues(t, len(config.ServerGroup), 5)
	assert.EqualValues(t, len(config.Server), 6)
}
