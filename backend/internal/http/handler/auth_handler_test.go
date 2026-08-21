package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	authdomain "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/auth"
	authmock "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/auth/mock"
	tehttp "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/http"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/http/handler"
)

func newAuthEcho(t *testing.T) *echo.Echo {
	t.Helper()
	e := echo.New()
	e.Validator = tehttp.NewValidator()
	return e
}

func TestAuthHandler_Login_Success(t *testing.T) {
	svc := new(authmock.Authenticator)
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	now := time.Now().UTC()
	result := &authdomain.SignInResult{
		AccessToken:              "acc.tok.en",
		RefreshToken:             "ref.tok.en",
		AccessTokenExpiresInSec:  3600,
		RefreshTokenExpiresInSec: 86400,
		IssuedAt:                 now,
		User: &authdomain.User{
			ID:    userID,
			Email: "demo@example.com",
			Name:  "Demo User",
		},
	}
	svc.On("SignIn", mock.Anything, "demo@example.com", "Demo1234!").Return(result, nil)

	h := handler.NewAuthHandler(svc)
	e := newAuthEcho(t)

	body := `{"email":"demo@example.com","password":"Demo1234!"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, h.Login(c))
	require.Equal(t, http.StatusCreated, rec.Code)

	resp := decode(t, rec)
	assert.Equal(t, true, resp["success"])
	assert.EqualValues(t, 201, resp["code"])
	assert.Equal(t, "logged in", resp["message"])
	data := resp["data"].(map[string]any)
	assert.Equal(t, "acc.tok.en", data["access_token"])
	assert.Equal(t, "ref.tok.en", data["refresh_token"])
	user := data["user"].(map[string]any)
	assert.Equal(t, "demo@example.com", user["email"])
	assert.Equal(t, "Demo User", user["name"])
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	svc := new(authmock.Authenticator)
	svc.On("SignIn", mock.Anything, "bad@example.com", "BadPass1!").
		Return(nil, authdomain.ErrInvalidCredentials)

	h := handler.NewAuthHandler(svc)
	e := newAuthEcho(t)

	body := `{"email":"bad@example.com","password":"BadPass1!"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, h.Login(c))
	require.Equal(t, http.StatusUnauthorized, rec.Code)

	resp := decode(t, rec)
	assert.Equal(t, false, resp["success"])
	assert.EqualValues(t, 40100, resp["code"])
}

func TestAuthHandler_Login_ValidationFailure(t *testing.T) {
	svc := new(authmock.Authenticator)
	h := handler.NewAuthHandler(svc)
	e := newAuthEcho(t)

	// Empty body — missing required fields
	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, h.Login(c))
	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	resp := decode(t, rec)
	assert.Equal(t, false, resp["success"])
	assert.EqualValues(t, 42200, resp["code"])
	svc.AssertNotCalled(t, "SignIn")
}

func TestAuthHandler_Renew_Success(t *testing.T) {
	svc := new(authmock.Authenticator)
	now := time.Now().UTC()
	result := &authdomain.RefreshResult{
		AccessToken:             "new.acc.tok",
		AccessTokenExpiresInSec: 3600,
		IssuedAt:                now,
	}
	svc.On("Renew", mock.Anything, "valid.refresh.token").Return(result, nil)

	h := handler.NewAuthHandler(svc)
	e := newAuthEcho(t)

	body := `{"refresh_token":"valid.refresh.token"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/renew", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, h.Renew(c))
	require.Equal(t, http.StatusOK, rec.Code)

	resp := decode(t, rec)
	assert.Equal(t, true, resp["success"])
	assert.EqualValues(t, 200, resp["code"])
	assert.Equal(t, "token renewed", resp["message"])
	data := resp["data"].(map[string]any)
	assert.Equal(t, "new.acc.tok", data["access_token"])
}

func TestAuthHandler_Renew_InvalidToken(t *testing.T) {
	svc := new(authmock.Authenticator)
	svc.On("Renew", mock.Anything, "bad.token").Return(nil, authdomain.ErrRefreshInvalid)

	h := handler.NewAuthHandler(svc)
	e := newAuthEcho(t)

	body := `{"refresh_token":"bad.token"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/renew", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, h.Renew(c))
	require.Equal(t, http.StatusUnauthorized, rec.Code)

	resp := decode(t, rec)
	assert.Equal(t, false, resp["success"])
	assert.EqualValues(t, 40100, resp["code"])
}

func TestAuthHandler_Logout_NoBody(t *testing.T) {
	svc := new(authmock.Authenticator)
	h := handler.NewAuthHandler(svc)
	e := newAuthEcho(t)

	req := httptest.NewRequest(http.MethodPost, "/v1/auth/logout", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, h.Logout(c))
	require.Equal(t, http.StatusOK, rec.Code)

	resp := decode(t, rec)
	assert.Equal(t, true, resp["success"])
	assert.EqualValues(t, 200, resp["code"])
	assert.Equal(t, "logged out", resp["message"])
	svc.AssertNotCalled(t, "Logout")
}

func TestAuthHandler_Logout_WithToken(t *testing.T) {
	svc := new(authmock.Authenticator)
	svc.On("Logout", mock.Anything, "my.refresh.token").Return(nil)

	h := handler.NewAuthHandler(svc)
	e := newAuthEcho(t)

	body := `{"refresh_token":"my.refresh.token"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/auth/logout", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, h.Logout(c))
	require.Equal(t, http.StatusOK, rec.Code)

	resp := decode(t, rec)
	assert.Equal(t, true, resp["success"])
	svc.AssertCalled(t, "Logout", mock.Anything, "my.refresh.token")
}
