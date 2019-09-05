package rule

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRuleSet(t *testing.T) {
	h, err := ruleSetHandle(context.Background(), &Rule{
		Typ:   KeyRuleSet,
		Value: "https://raw.githubusercontent.com/lhie1/Rules/master/Surge3/Proxy.list",
		Params: map[string]string{
			ParamsKeyInterval: "10s",
		},
	}, func(ctx context.Context, info RequestInfo) *Rule {
		return &Rule{
			Typ:   "h",
			Value: "hehe",
			Proxy: "haha",
		}
	}, nil)
	assert.NoError(t, err)
	ctx := context.Background()
	req := &request{
		network:     "tcp",
		domain:      "www.guqiguqi.com",
		ip:          []byte{13, 32, 1, 1},
		port:        80,
		countryCode: "US",
	}
	r := h(ctx, req)
	fmt.Println(r.String())
	time.Sleep(time.Second * 5)
	r = h(ctx, req)
	fmt.Println(r.String())
	time.Sleep(time.Second * 10)
}

type request struct {
	network     string
	domain      string
	uri         string
	ip          net.IP
	port        int
	countryCode string
}

func (r *request) Network() string {
	return r.network
}
func (r *request) Domain() string {
	return r.domain
}
func (r *request) URI() string {
	return r.uri
}
func (r *request) IP() net.IP {
	return r.ip
}
func (r *request) CountryCode() string {
	return r.countryCode
}
func (r *request) Port() int {
	return r.port
}
func (r *request) SetIP(in net.IP) {
	r.ip = in
}
func (r *request) SetPort(in int) {
	r.port = in
}
func (r *request) SetCountryCode(in string) {
	r.countryCode = in
}
