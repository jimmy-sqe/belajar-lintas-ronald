package ctxutil

import "context"

type userIDKey struct{}

// WithUserID attaches a user ID to the context.
func WithUserID(ctx context.Context, uid string) context.Context {
	return context.WithValue(ctx, userIDKey{}, uid)
}

// UserID returns the user ID or empty string.
func UserID(ctx context.Context) string {
	if v, ok := ctx.Value(userIDKey{}).(string); ok {
		return v
	}
	return ""
}
