package http

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// ValidatorAdapter wraps go-playground/validator for Echo.
type ValidatorAdapter struct {
	v *validator.Validate
}

// NewValidator returns a fresh ValidatorAdapter.
func NewValidator() *ValidatorAdapter {
	return &ValidatorAdapter{v: validator.New()}
}

// Validate satisfies echo.Validator.
func (a *ValidatorAdapter) Validate(i any) error {
	if err := a.v.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
