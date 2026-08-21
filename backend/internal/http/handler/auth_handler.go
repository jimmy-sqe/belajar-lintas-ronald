package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	authdomain "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/auth"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/http/handler/dto"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/http/response"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/pkg/customerror"
)

// AuthHandler exposes /v1/auth/* endpoints.
type AuthHandler struct {
	svc authdomain.Authenticator
}

// NewAuthHandler constructs an AuthHandler.
func NewAuthHandler(svc authdomain.Authenticator) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Login handles POST /v1/auth/login.
func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.Err(c, customerror.ErrBadRequest.WithMetadata(
			map[string]any{"detail": err.Error()}))
	}
	if err := c.Validate(&req); err != nil {
		return response.Err(c, customerror.ErrUnprocessable.WithMetadata(
			map[string]any{"detail": err.Error()}))
	}

	result, err := h.svc.SignIn(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return response.Err(c, mapAuthError(err))
	}

	return response.OK(c, http.StatusCreated, "logged in", dto.LoginResponse{
		AccessToken:              result.AccessToken,
		RefreshToken:             result.RefreshToken,
		AccessTokenExpiresInSec:  result.AccessTokenExpiresInSec,
		RefreshTokenExpiresInSec: result.RefreshTokenExpiresInSec,
		IssuedAt:                 result.IssuedAt.UnixMilli(),
		User: dto.UserResponse{
			ID:    result.User.ID.String(),
			Name:  result.User.Name,
			Email: result.User.Email,
		},
		Permissions: []string{},
	})
}

// Logout handles POST /v1/auth/logout.
func (h *AuthHandler) Logout(c echo.Context) error {
	var req dto.LogoutRequest
	_ = c.Bind(&req)
	if req.RefreshToken != "" {
		if err := h.svc.Logout(c.Request().Context(), req.RefreshToken); err != nil {
			return response.Err(c, customerror.ErrInternalServer.WithMetadata(
				map[string]any{"detail": err.Error()}))
		}
	}
	return response.OK(c, http.StatusOK, "logged out", nil)
}

// Renew handles POST /v1/auth/renew.
func (h *AuthHandler) Renew(c echo.Context) error {
	var req dto.RenewRequest
	if err := c.Bind(&req); err != nil {
		return response.Err(c, customerror.ErrBadRequest.WithMetadata(
			map[string]any{"detail": err.Error()}))
	}
	if err := c.Validate(&req); err != nil {
		return response.Err(c, customerror.ErrUnprocessable.WithMetadata(
			map[string]any{"detail": err.Error()}))
	}

	result, err := h.svc.Renew(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return response.Err(c, mapAuthError(err))
	}

	return response.OK(c, http.StatusOK, "token renewed", dto.RenewResponse{
		AccessToken:             result.AccessToken,
		AccessTokenExpiresInSec: result.AccessTokenExpiresInSec,
		IssuedAt:                result.IssuedAt.UnixMilli(),
	})
}

func mapAuthError(err error) customerror.CustomError {
	switch {
	case errors.Is(err, authdomain.ErrInvalidCredentials):
		return customerror.ErrUnauthorized.WithMetadata(
			map[string]any{"detail": "invalid credentials"})
	case errors.Is(err, authdomain.ErrSessionExpired),
		errors.Is(err, authdomain.ErrRefreshInvalid):
		return customerror.ErrUnauthorized.WithMetadata(
			map[string]any{"detail": err.Error()})
	}
	var ce customerror.CustomError
	if errors.As(err, &ce) {
		return ce
	}
	return customerror.ErrInternalServer
}
