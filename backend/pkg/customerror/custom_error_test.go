package customerror

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstructors_SetHTTPCode(t *testing.T) {
	cases := []struct {
		name     string
		err      CustomError
		wantHTTP int
		wantCode ErrorCode
	}{
		{"client", NewClientError(40001, "bad"), http.StatusBadRequest, 40001},
		{"unauthorized", NewUnauthorizedError(40101, "no token"), http.StatusUnauthorized, 40101},
		{"forbidden", NewForbiddenError(40301, "denied"), http.StatusForbidden, 40301},
		{"notfound", NewNotFoundError(40401, "missing"), http.StatusNotFound, 40401},
		{"conflict", NewConflictError(40901, "dup"), http.StatusConflict, 40901},
		{"unprocessable", NewUnprocessableEntityError(42201, "invalid"), http.StatusUnprocessableEntity, 42201},
		{"too-many", NewTooManyRequestError(42901, "rate"), http.StatusTooManyRequests, 42901},
		{"unsupported", NewUnsupportedMediaType(41501, "bad type"), http.StatusUnsupportedMediaType, 41501},
		{"internal", NewInternalServerError(50001, "boom"), http.StatusInternalServerError, 50001},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.wantHTTP, tc.err.HTTPCode())
			assert.Equal(t, tc.wantCode, tc.err.Code())
		})
	}
}

func TestCustomError_MessageAndError(t *testing.T) {
	e := NewInternalServerError(CodeInternalServer, "boom")
	assert.Equal(t, "boom", e.Error())
	assert.Equal(t, "boom", e.Message())
	assert.Equal(t, CodeInternalServer, e.Code())
	assert.Equal(t, http.StatusInternalServerError, e.HTTPCode())
}

func TestCustomError_WithMetadata_ReturnsCopy(t *testing.T) {
	base := NewClientError(CodeBadRequest, "bad")
	enriched := base.WithMetadata(map[string]any{"field": "name"})

	assert.Nil(t, base.Metadata(), "base must not mutate")
	assert.Equal(t, "name", enriched.Metadata()["field"])
}

func TestCustomError_Sprintf_FormatsMessage(t *testing.T) {
	tmpl := NewClientError(CodeBadRequest, "field %s missing")
	got := tmpl.Sprintf("name")
	assert.Equal(t, "field name missing", got.Error())
	assert.Equal(t, "field %s missing", tmpl.Error(), "template must not mutate")
}

func TestCustomError_Equal(t *testing.T) {
	a := NewNotFoundError(CodeNotFound, "x")
	b := NewNotFoundError(CodeNotFound, "x")
	c := NewNotFoundError(CodeNotFound, "y")

	assert.True(t, a.Equal(b))
	assert.False(t, a.Equal(c))
	assert.False(t, a.Equal(errors.New("not custom")))
}

func TestPrebuiltGenericErrors_HaveExpectedShape(t *testing.T) {
	assert.Equal(t, CodeBadRequest, ErrBadRequest.Code())
	assert.Equal(t, http.StatusBadRequest, ErrBadRequest.HTTPCode())

	assert.Equal(t, CodeNotFound, ErrNotFound.Code())
	assert.Equal(t, http.StatusNotFound, ErrNotFound.HTTPCode())

	assert.Equal(t, CodeConflict, ErrConflict.Code())
	assert.Equal(t, http.StatusConflict, ErrConflict.HTTPCode())

	assert.Equal(t, CodeInternalServer, ErrInternalServer.Code())
	assert.Equal(t, http.StatusInternalServerError, ErrInternalServer.HTTPCode())
}
