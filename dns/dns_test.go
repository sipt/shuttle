package dns

import (
	"context"
	"net"
	"testing"

	"github.com/sipt/shuttle/conf/model"

	"github.com/stretchr/testify/assert"
)

func TestParseDnsServer(t *testing.T) {
	_, err := ParseDnsServer("://8.8.8.8:53")
	assert.Error(t, err)
	_, err = ParseDnsServer("tcp://8.8.8.8:53")
	assert.Error(t, err)

	s, err := ParseDnsServer("8.8.8.8")
	assert.NoError(t, err)
	assert.EqualValues(t, s.String(), "udp://8.8.8.8:53")

	s, err = ParseDnsServer("udp://8.8.8.8:88")
	assert.NoError(t, err)
	assert.EqualValues(t, s.String(), "udp://8.8.8.8:88")
}

func TestResolveDomain(t *testing.T) {
	ctx := context.Background()
	ips, server, err := ResolveDomain(ctx, "www.baidu.com", []*DnsServer{{"udp", net.IP{8, 8, 8, 8}, 53}, {"udp", net.IP{114, 114, 114, 114}, 53}}...)
	assert.NoError(t, err)
	assert.NotEqual(t, len(ips), 0)
	assert.True(t, len(server.IP.String()) > 0)
}

func TestApplyConfig(t *testing.T) {
	fileName := "../GeoLite2-Country.mmdb"
	h, err := ApplyConfig(&model.Config{
		DNS: struct {
			IncludeSystem bool     `json:"include_system"`
			Servers       []string `toml:"servers"`
			Mapping       []struct {
				Domain string   `json:"domain"`
				IP     []string `json:"ip"`
				Server []string `json:"server"`
			} `json:"mapping"`
		}{
			Servers: []string{
				"udp://8.8.8.8:53",
				"114.114.114.114",
			},
		},
	}, func(ctx context.Context, domain string) *DNS {
		t.Error("failed")
		return nil
	})
	assert.NoError(t, err)
	ctx := context.Background()
	answer := h(ctx, "www.baidu.com")
	assert.NotNil(t, answer)
	assert.True(t, len(answer.CurrentIP.String()) > 0)
	assert.EqualValues(t, answer.CurrentCountry, "CN")
}
