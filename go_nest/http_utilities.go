package goNest

import "context"

// ContextWithRequestID :  Store request ID in context
func ContextWithRequestID(ctx context.Context, key string, val any) context.Context {
	return context.WithValue(ctx, key, val)
}
