# Error Codes

This service follows the **lintas team API contract** for error responses. Every error response carries a numeric `code` field in the following format:

```
HTTP_STATUS_PREFIX  ||  2-digit subcode
   (3 digits)       ||    (00 to 99)
```

Example: `40100` = HTTP 401 (Unauthorized), subcode `00` (generic).

## Subcode Convention

| Range  | Meaning                          | Who owns it                          |
|--------|----------------------------------|--------------------------------------|
| `00`   | Generic for the HTTP status      | Boilerplate (`pkg/customerror`)      |
| `01-99`| Service-specific variants        | Each scaffolded service              |

## Generic Codes Shipped With the Boilerplate

Defined in `pkg/customerror/error_code.go`:

| Code    | Constant              | HTTP | Pre-built `Err*` var | Typical use                                |
|---------|-----------------------|------|----------------------|--------------------------------------------|
| `40000` | `CodeBadRequest`      | 400  | `ErrBadRequest`      | Malformed body, missing content-type, bind failure |
| `40100` | `CodeUnauthorized`    | 401  | `ErrUnauthorized`    | Missing / invalid / expired token          |
| `40300` | `CodeForbidden`       | 403  | `ErrForbidden`       | Authenticated but not authorized           |
| `40400` | `CodeNotFound`        | 404  | `ErrNotFound`        | Resource lookup failure                    |
| `40900` | `CodeConflict`        | 409  | `ErrConflict`        | Duplicate / state conflict                 |
| `42200` | `CodeUnprocessable`   | 422  | `ErrUnprocessable`   | Validator failure (semantically invalid)   |
| `42900` | `CodeTooManyRequests` | 429  | `ErrTooManyRequests` | Rate limited                               |
| `50000` | `CodeInternalServer`  | 500  | `ErrInternalServer`  | Unexpected server error (fallback)         |

## Adding Service-Specific Codes

When you scaffold a service from this boilerplate, allocate domain-specific subcodes in the `01-99` range. Add them to `pkg/customerror/error_code.go` (or a sibling file in the same package) following this pattern:

```go
// service-specific catalog
const (
    // Todo domain
    CodeTodoNotFound       ErrorCode = 40410
    CodeTodoInvalidDueDate ErrorCode = 42210

    // User domain
    CodeUserAlreadyExists  ErrorCode = 40910
)

var (
    ErrTodoNotFound        = NewNotFoundError(CodeTodoNotFound, "todo not found")
    ErrTodoInvalidDueDate  = NewUnprocessableEntityError(CodeTodoInvalidDueDate, "invalid todo due date")
    ErrUserAlreadyExists   = NewConflictError(CodeUserAlreadyExists, "user already exists")
)
```

### Guidelines

- **Pick a free subcode in the relevant HTTP range.** `40400` is generic-not-found; `40410` is a service-specific not-found variant.
- **Document what triggers each code** in a comment next to the constant.
- **Reuse the generic `Err*` vars when appropriate** — do not invent a service-specific code just to rename the message; use `customerror.ErrNotFound.WithMetadata(...)` instead.
- **Allocate codes monotonically** within the file (`40410`, `40411`, `40412`, ...) for predictability.
- **Document codes in the service's OpenAPI** — every code your handlers can emit should appear in an example `ErrorResponse` somewhere in `openapi.yaml`.

## Validation vs Body-Parse

- **Body parse failure** (malformed JSON, bind error, missing content-type) → HTTP 400, code `40000`.
- **Validator failure** (field out of range, format mismatch, enum violation) → HTTP 422, code `42200`.

This split follows RFC 7231 / 4918 and matches the convention used by mainstream stacks (FastAPI, Spring `@Valid`).

## References

- Lintas source of truth: `lintas/plugins/sdlc-core/templates/api-contract.yaml`
- Boilerplate envelope renderer: `internal/http/response/envelope.go`
- Custom error primitive: `pkg/customerror/`
