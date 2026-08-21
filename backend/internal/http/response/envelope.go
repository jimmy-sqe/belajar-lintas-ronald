// Package response renders HTTP responses in the lintas-standard envelope format.
//
// Three envelopes are exposed:
//   - SuccessEnvelope     — single-resource and no-content responses
//   - PaginatedEnvelope   — list responses with pagination metadata
//   - ErrorEnvelope       — all error responses, with optional metadata
//
// Handlers MUST use the helpers (OK / Paged / Err); manual c.JSON of
// envelope-shaped maps defeats the single-source-of-truth guarantee.
package response

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/pkg/customerror"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/pkg/pagination"
)

// SuccessEnvelope matches lintas SuccessResponse.
type SuccessEnvelope struct {
	Success   bool      `json:"success"`
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Data      any       `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

// PaginatedEnvelope matches lintas PaginatedResponse.
type PaginatedEnvelope struct {
	Success    bool            `json:"success"`
	Code       int             `json:"code"`
	Message    string          `json:"message"`
	Data       any             `json:"data"`
	Pagination pagination.Page `json:"pagination"`
	Timestamp  time.Time       `json:"timestamp"`
}

// ErrorEnvelope matches lintas ErrorResponse, with an optional metadata
// extension (omitempty so the wire format stays contract-compliant when empty).
type ErrorEnvelope struct {
	Success   bool                  `json:"success"`
	Code      customerror.ErrorCode `json:"code"`
	Message   string                `json:"message"`
	Timestamp time.Time             `json:"timestamp"`
	Metadata  map[string]any        `json:"metadata,omitempty"`
}

// OK renders SuccessEnvelope. data=nil is encoded as `{}` to satisfy the
// lintas requirement that `data` is always a present object.
func OK(c echo.Context, status int, message string, data any) error {
	if data == nil {
		data = struct{}{}
	}
	return c.JSON(status, SuccessEnvelope{
		Success:   true,
		Code:      status,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC(),
	})
}

// Paged renders PaginatedEnvelope. data should be a slice; the helper does
// not enforce slice-ness so callers may pass any JSON-marshalable array.
func Paged(c echo.Context, message string, data any, page pagination.Page) error {
	return c.JSON(http.StatusOK, PaginatedEnvelope{
		Success:    true,
		Code:       http.StatusOK,
		Message:    message,
		Data:       data,
		Pagination: page,
		Timestamp:  time.Now().UTC(),
	})
}

// Err renders ErrorEnvelope. A customerror.CustomError is rendered as-is;
// any other error falls back to customerror.ErrInternalServer (caller is
// expected to log the original error upstream via middleware).
func Err(c echo.Context, err error) error {
	var ce customerror.CustomError
	if !errors.As(err, &ce) {
		ce = customerror.ErrInternalServer
	}
	return c.JSON(ce.HTTPCode(), ErrorEnvelope{
		Success:   false,
		Code:      ce.Code(),
		Message:   ce.Message(),
		Timestamp: time.Now().UTC(),
		Metadata:  ce.Metadata(),
	})
}
