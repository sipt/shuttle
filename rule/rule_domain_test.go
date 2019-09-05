package rule

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainRule(t *testing.T) {
	defaultRule := &Rule{}
	handle := func(ctx context.Context, info RequestInfo) *Rule {
		return defaultRule
	}
	var err error
	ctx := context.Background()
	google := &Rule{Value: "google.com"}
	handle, err = domainHandle(nil, google, handle, nil)
	assert.NoError(t, err)
	facebook := &Rule{Value: "facebook"}
	handle, err = domainKeywordHandle(nil, facebook, handle, nil)
	assert.NoError(t, err)
	github := &Rule{Value: "github.com"}
	handle, err = domainSuffixHandle(nil, github, handle, nil)
	assert.NoError(t, err)

	assert.Equal(t, handle(ctx, &info{domain: "www.google.com"}), defaultRule)
	assert.Equal(t, handle(ctx, &info{domain: "google.com"}), google)
	assert.Equal(t, handle(ctx, &info{domain: "facebook"}), facebook)
	assert.Equal(t, handle(ctx, &info{domain: "www.github.com"}), github)
}

type info struct {
	domain string
	RequestInfo
}

func (i *info) Domain() string {
	return i.domain
}
