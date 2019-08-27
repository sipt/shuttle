package rule

import (
	"context"

	"github.com/sipt/shuttle/dns"
)

const (
	KeyFinal = "FINAL"
)

func init() {
	Register(KeyFinal, finalHandle)
}
func finalHandle(rule *Rule, _ Handle, _ dns.Handle) (Handle, error) {
	return func(ctx context.Context, info RequestInfo) *Rule {
		return rule
	}, nil
}
