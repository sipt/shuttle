package rule

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuleIPCiDR(t *testing.T) {
	h, err := ipCidrHandle(nil, &Rule{
		Typ:   "IP-CIDR",
		Value: "13.32.0.0/16",
		Proxy: "hehe",
	}, func(ctx context.Context, info RequestInfo) *Rule {
		return &Rule{
			Typ:   "IP-CIDR",
			Value: "13.32.0.0/16",
			Proxy: "haha",
		}
	}, nil)
	assert.NoError(t, err)
	req := &request{
		network:     "tcp",
		domain:      "www.guqiguqi.com",
		ip:          []byte{13, 32, 1, 1},
		port:        80,
		countryCode: "US",
	}
	r := h(nil, req)
	assert.EqualValues(t, r.Proxy, "hehe")
}
