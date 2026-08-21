package dto

// LoginRequest is the POST /v1/auth/login body.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

// UserResponse is included inside LoginResponse.
type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email"`
}

// LoginResponse is the payload shape wrapped in SuccessEnvelope.
type LoginResponse struct {
	AccessToken              string       `json:"access_token"`
	RefreshToken             string       `json:"refresh_token"`
	AccessTokenExpiresInSec  int          `json:"access_token_expires_in_sec"`
	RefreshTokenExpiresInSec int          `json:"refresh_token_expires_in_sec"`
	IssuedAt                 int64        `json:"issued_at"`
	User                     UserResponse `json:"user"`
	Permissions              []string     `json:"permissions"`
}

// RenewRequest is the POST /v1/auth/renew body.
type RenewRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RenewResponse is the payload shape wrapped in SuccessEnvelope.
type RenewResponse struct {
	AccessToken             string `json:"access_token"`
	AccessTokenExpiresInSec int    `json:"access_token_expires_in_sec"`
	IssuedAt                int64  `json:"issued_at"`
}

// LogoutRequest is the optional POST /v1/auth/logout body.
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}
