package rule

import (
	"context"
	"regexp"

	"github.com/pkg/errors"
)

const (
	KeyUrlRegex = "URL-REGEX"
)

func init() {
	Register(KeyUrlRegex, urlRegexHandle)
}
func urlRegexHandle(rule *Rule, next Handle) (Handle, error) {
	reg, err := regexp.Compile(rule.Value)
	if err != nil {
		return nil, errors.Errorf("rule:[%s, %s, %s, %v], regex:[%s] invalid",
			rule.Typ, rule.Value, rule.Proxy, rule.Params, rule.Value)
	}
	return func(ctx context.Context, info Info) *Rule {
		if reg.MatchString(info.URI()) {
			return rule
		}
		return next(ctx, info)
	}, nil
}
