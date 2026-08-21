package customerror

// Generic error codes shipped with the boilerplate. The catalog follows the
// lintas API contract: HTTP_PREFIX + 2-digit subcode (subcode 00 = generic).
//
// Services scaffolded from this boilerplate SHOULD extend this catalog
// with their own domain-specific subcodes in the 01-99 range per HTTP status.
// See docs/ERROR_CODES.md for the full convention.
const (
	CodeBadRequest      ErrorCode = 40000
	CodeUnauthorized    ErrorCode = 40100
	CodeForbidden       ErrorCode = 40300
	CodeNotFound        ErrorCode = 40400
	CodeConflict        ErrorCode = 40900
	CodeUnprocessable   ErrorCode = 42200
	CodeTooManyRequests ErrorCode = 42900
	CodeInternalServer  ErrorCode = 50000
)

// Pre-built generic errors for direct use in handlers and helpers.
var (
	ErrBadRequest      = NewClientError(CodeBadRequest, "bad request")
	ErrUnauthorized    = NewUnauthorizedError(CodeUnauthorized, "unauthorized")
	ErrForbidden       = NewForbiddenError(CodeForbidden, "forbidden")
	ErrNotFound        = NewNotFoundError(CodeNotFound, "not found")
	ErrConflict        = NewConflictError(CodeConflict, "conflict")
	ErrUnprocessable   = NewUnprocessableEntityError(CodeUnprocessable, "unprocessable entity")
	ErrTooManyRequests = NewTooManyRequestError(CodeTooManyRequests, "too many requests")
	ErrInternalServer  = NewInternalServerError(CodeInternalServer, "internal server error")
)
