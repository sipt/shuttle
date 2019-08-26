package rule

import (
	"context"
)

const (
	KeyFinal = "FINAL"
)

func init() {
	Register(KeyFinal, finalHandle)
}
func finalHandle(rule *Rule, _ Handle) (Handle, error) {
	return func(ctx context.Context, info Info) *Rule {
		return rule
	}, nil
}
