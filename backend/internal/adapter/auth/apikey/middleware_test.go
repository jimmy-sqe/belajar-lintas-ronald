package apikey_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/auth/apikey"
)

func setup(t *testing.T) *echo.Echo {
	t.Helper()
	v, err := apikey.NewVerifier(apikey.Config{
		HeaderName: "X-API-Key",
		Keys:       []string{"valid-key-1", "valid-key-2"},
	})
	require.NoError(t, err)
	e := echo.New()
	e.Use(apikey.Middleware(v))
	e.GET("/protected", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	return e
}

func TestMiddleware_ValidKey(t *testing.T) {
	e := setup(t)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("X-API-Key", "valid-key-1")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestMiddleware_InvalidKey(t *testing.T) {
	e := setup(t)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("X-API-Key", "wrong")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestMiddleware_MissingHeader(t *testing.T) {
	e := setup(t)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}
