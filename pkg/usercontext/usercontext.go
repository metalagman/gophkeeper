package usercontext

import (
	"context"
	"github.com/google/uuid"
)

var EmptyUUID = uuid.NullUUID{
	Valid: false,
}

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

func ReadContextUUID(ctx context.Context, key interface{}) uuid.NullUUID {
	v := ctx.Value(key)
	if v == nil {
		return EmptyUUID
	}
	s, ok := v.(uuid.UUID)
	if !ok {
		return EmptyUUID
	}
	return uuid.NullUUID{
		UUID:  s,
		Valid: true,
	}
}

func ReadUID(ctx context.Context) uuid.NullUUID {
	return ReadContextUUID(ctx, ContextKeyUID{})
}

func WriteUID(ctx context.Context, uid uuid.UUID) context.Context {
	return context.WithValue(ctx, ContextKeyUID{}, uid)
}
