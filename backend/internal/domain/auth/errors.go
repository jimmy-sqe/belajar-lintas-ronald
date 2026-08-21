package auth

import "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/pkg/customerror"

// Sentinel errors for the auth domain. Handlers map these to the
// lintas error envelope via internal/http/response.Err.
var (
	ErrInvalidCredentials = customerror.NewUnauthorizedError(
		customerror.CodeUnauthorized, "invalid credentials")
	ErrSessionExpired = customerror.NewUnauthorizedError(
		customerror.CodeUnauthorized, "session expired")
	ErrRefreshInvalid = customerror.NewUnauthorizedError(
		customerror.CodeUnauthorized, "refresh token invalid or expired")
	ErrUserNotFound = customerror.NewNotFoundError(
		customerror.CodeNotFound, "user not found")
)
