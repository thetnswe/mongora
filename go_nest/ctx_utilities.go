package goNest

import "context"

func GetCtxStringValue(ctx context.Context, key interface{}) string {
	if val := ctx.Value(key); val != nil {
		// Try to assert the value as a string
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	// Return empty string if value is nil or not a string
	return ""
}

func GetCtxIntValue(ctx context.Context, key interface{}) int {
	if val := ctx.Value(key); val != nil {
		// Try to assert the value as a string
		if intVal, ok := val.(int); ok {
			return intVal
		}
	}
	// Return 0 if nil or not int value
	return 0
}
