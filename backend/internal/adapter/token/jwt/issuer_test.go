package jwt_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/adapter/token/jwt"
	authdomain "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/auth"
)

func TestIssuer_SignAndParse_RoundTrip(t *testing.T) {
	iss := jwt.NewIssuer("test-secret-32-chars-or-more-pad...")
	userID := uuid.New()
	now := time.Now().UTC().Truncate(time.Second)

	claims := authdomain.Claims{
		UserID:    userID,
		Email:     "demo@example.com",
		Kind:      authdomain.KindAccess,
		IssuedAt:  now,
		ExpiresAt: now.Add(1 * time.Hour),
	}

	token, err := iss.Sign(claims)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	parsed, err := iss.Parse(token)
	require.NoError(t, err)
	assert.Equal(t, userID, parsed.UserID)
	assert.Equal(t, "demo@example.com", parsed.Email)
	assert.Equal(t, authdomain.KindAccess, parsed.Kind)
}

func TestIssuer_Parse_RejectsTamperedToken(t *testing.T) {
	iss := jwt.NewIssuer("test-secret-32-chars-or-more-pad...")
	token, err := iss.Sign(authdomain.Claims{
		UserID: uuid.New(), Email: "demo@example.com",
		Kind: authdomain.KindAccess, IssuedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour),
	})
	require.NoError(t, err)

	tampered := token + "x"
	_, err = iss.Parse(tampered)
	assert.Error(t, err)
}

func TestIssuer_Parse_RejectsExpiredToken(t *testing.T) {
	iss := jwt.NewIssuer("test-secret-32-chars-or-more-pad...")
	now := time.Now().UTC()
	token, err := iss.Sign(authdomain.Claims{
		UserID: uuid.New(), Email: "demo@example.com",
		Kind: authdomain.KindAccess, IssuedAt: now.Add(-2 * time.Hour),
		ExpiresAt: now.Add(-1 * time.Hour),
	})
	require.NoError(t, err)

	_, err = iss.Parse(token)
	assert.Error(t, err)
}
