// Package mock provides testify mock implementations for auth domain interfaces.
package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	authdomain "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/auth"
)

// Authenticator is a mock implementation of authdomain.Authenticator.
type Authenticator struct{ mock.Mock }

func (m *Authenticator) SignIn(ctx context.Context, email, password string) (*authdomain.SignInResult, error) {
	args := m.Called(ctx, email, password)
	v, _ := args.Get(0).(*authdomain.SignInResult)
	return v, args.Error(1)
}

func (m *Authenticator) Logout(ctx context.Context, refreshToken string) error {
	return m.Called(ctx, refreshToken).Error(0)
}

func (m *Authenticator) Renew(ctx context.Context, refreshToken string) (*authdomain.RefreshResult, error) {
	args := m.Called(ctx, refreshToken)
	v, _ := args.Get(0).(*authdomain.RefreshResult)
	return v, args.Error(1)
}
