package rule

import (
	"context"
	"strings"
)

const (
	KeyDomainSuffix  = "DOMAIN-SUFFIX"
	KeyDomain        = "DOMAIN"
	KeyDomainKeyword = "DOMAIN-KEYWORD"
)

func init() {
	Register(KeyDomainSuffix, domainSuffixHandle)
	Register(KeyDomain, domainHandle)
	Register(KeyDomainKeyword, domainKeywordHandle)
}
func domainSuffixHandle(rule *Rule, next Handle) (Handle, error) {
	return func(ctx context.Context, info Info) *Rule {
		if strings.HasSuffix(info.Domain(), rule.Value) {
			return rule
		}
		return next(ctx, info)
	}, nil
}
func domainHandle(rule *Rule, next Handle) (Handle, error) {
	return func(ctx context.Context, info Info) *Rule {
		if len(info.Domain()) == len(rule.Value) && info.Domain() == rule.Value {
			return rule
		}
		return next(ctx, info)
	}, nil
}
func domainKeywordHandle(rule *Rule, next Handle) (Handle, error) {
	return func(ctx context.Context, info Info) *Rule {
		if len(info.Domain()) >= len(rule.Value) && strings.Index(info.Domain(), rule.Value) > -1 {
			return rule
		}
		return next(ctx, info)
	}, nil
}
