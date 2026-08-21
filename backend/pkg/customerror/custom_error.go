package customerror

import (
	"fmt"
	"net/http"
)

// ErrorCode is a 5-digit numeric code: HTTP_PREFIX + 2-digit subcode.
// Convention documented in docs/ERROR_CODES.md.
type ErrorCode int

type CustomError struct {
	code     ErrorCode
	message  string
	httpCode int
	metadata map[string]any
}

func (ce CustomError) Error() string            { return ce.message }
func (ce CustomError) HTTPCode() int            { return ce.httpCode }
func (ce CustomError) Code() ErrorCode          { return ce.code }
func (ce CustomError) Message() string          { return ce.message }
func (ce CustomError) Metadata() map[string]any { return ce.metadata }

func (ce CustomError) WithMetadata(metadata map[string]any) CustomError {
	ce.metadata = metadata
	return ce
}

// Equal compares CustomErrors by code, httpCode, and message. Metadata is
// intentionally ignored so enriched errors (via WithMetadata) remain equal
// to their base error.
func (ce CustomError) Equal(err error) bool {
	if e, ok := err.(CustomError); ok {
		return ce.code == e.code && ce.httpCode == e.httpCode && ce.message == e.message
	}
	return false
}

// Is implements errors.Is support so that errors.Is(enrichedErr, sentinel)
// returns true when both are CustomError values with the same code, httpCode,
// and message — regardless of metadata.
func (ce CustomError) Is(target error) bool {
	return ce.Equal(target)
}

func (ce CustomError) Sprintf(a ...any) CustomError {
	newCe := ce
	newCe.message = fmt.Sprintf(ce.message, a...)
	return newCe
}

func NewErrorWithCode(code ErrorCode, message string, httpCode int) CustomError {
	return CustomError{code, message, httpCode, nil}
}

func NewClientError(code ErrorCode, message string) CustomError {
	return NewErrorWithCode(code, message, http.StatusBadRequest)
}

func NewUnprocessableEntityError(code ErrorCode, message string) CustomError {
	return NewErrorWithCode(code, message, http.StatusUnprocessableEntity)
}

func NewForbiddenError(code ErrorCode, message string) CustomError {
	return NewErrorWithCode(code, message, http.StatusForbidden)
}

func NewUnauthorizedError(code ErrorCode, message string) CustomError {
	return NewErrorWithCode(code, message, http.StatusUnauthorized)
}

func NewTooManyRequestError(code ErrorCode, message string) CustomError {
	return NewErrorWithCode(code, message, http.StatusTooManyRequests)
}

func NewNotFoundError(code ErrorCode, message string) CustomError {
	return NewErrorWithCode(code, message, http.StatusNotFound)
}

func NewConflictError(code ErrorCode, message string) CustomError {
	return NewErrorWithCode(code, message, http.StatusConflict)
}

func NewInternalServerError(code ErrorCode, message string) CustomError {
	return NewErrorWithCode(code, message, http.StatusInternalServerError)
}

func NewUnsupportedMediaType(code ErrorCode, message string) CustomError {
	return NewErrorWithCode(code, message, http.StatusUnsupportedMediaType)
}
