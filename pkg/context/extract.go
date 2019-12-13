package context

import (
	"context"

	"github.com/sipt/shuttle/global"
	"github.com/sipt/shuttle/rule"
)

type ctxMarker struct{}

var (
	protocolKey    = &ctxMarker{}
	requestInfoKey = &ctxMarker{}
	ruleKey        = &ctxMarker{}
	namespaceKey   = &ctxMarker{}
	profileKey     = &ctxMarker{}
)

func WithRequestInfo(ctx context.Context, v global.RequestInfo) context.Context {
	return context.WithValue(ctx, requestInfoKey, v)
}

func ExtractRequestInfo(ctx context.Context) (global.RequestInfo, bool) {
	v, ok := ctx.Value(requestInfoKey).(global.RequestInfo)
	return v, ok
}

func WithRule(ctx context.Context, v *rule.Rule) context.Context {
	return context.WithValue(ctx, requestInfoKey, v)
}

func ExtractRule(ctx context.Context) (*rule.Rule, bool) {
	v, ok := ctx.Value(requestInfoKey).(*rule.Rule)
	return v, ok
}

func WithProtocol(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, requestInfoKey, v)
}

func ExtractProtocol(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(requestInfoKey).(string)
	return v, ok
}

func WithProfile(ctx context.Context, v *global.Profile) context.Context {
	return context.WithValue(ctx, requestInfoKey, v)
}

func ExtractProfile(ctx context.Context) (*global.Profile, bool) {
	v, ok := ctx.Value(requestInfoKey).(*global.Profile)
	return v, ok
}
