package usercontext

import "context"

type ContextKeyUID struct{}

func ReadContextString(ctx context.Context, key interface{}) string {
	v := ctx.Value(key)
	if v == nil {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}

func ReadUID(ctx context.Context) string {
	return ReadContextString(ctx, ContextKeyUID{})
}

func WriteUID(ctx context.Context, uid string) context.Context {
	return context.WithValue(ctx, ContextKeyUID{}, uid)
}
