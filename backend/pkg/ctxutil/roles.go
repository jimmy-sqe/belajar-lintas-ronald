package ctxutil

import "context"

type rolesKey struct{}

// WithRoles attaches roles to the context.
func WithRoles(ctx context.Context, roles []string) context.Context {
	return context.WithValue(ctx, rolesKey{}, roles)
}

// Roles returns roles or nil.
func Roles(ctx context.Context) []string {
	if v, ok := ctx.Value(rolesKey{}).([]string); ok {
		return v
	}
	return nil
}
