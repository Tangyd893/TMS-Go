package requestid

import "context"

type ctxKey struct{}

const Header = "X-Request-Id"

func FromContext(ctx context.Context) string {
	value, ok := ctx.Value(ctxKey{}).(string)
	if !ok {
		return ""
	}
	return value
}

func WithContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKey{}, id)
}
