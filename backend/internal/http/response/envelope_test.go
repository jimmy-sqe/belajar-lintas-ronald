package response_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/http/response"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/pkg/customerror"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/pkg/pagination"
)

func newCtx() (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

func TestOK_RendersSuccessEnvelope(t *testing.T) {
	c, rec := newCtx()
	require.NoError(t, response.OK(c, http.StatusOK, "fetched", map[string]any{"id": "abc"}))

	assert.Equal(t, http.StatusOK, rec.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))

	assert.Equal(t, true, body["success"])
	assert.EqualValues(t, 200, body["code"])
	assert.Equal(t, "fetched", body["message"])
	assert.NotEmpty(t, body["timestamp"])
	data, ok := body["data"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "abc", data["id"])
}

func TestOK_NilData_EmitsEmptyObject(t *testing.T) {
	c, rec := newCtx()
	require.NoError(t, response.OK(c, http.StatusOK, "deleted", nil))

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	data, ok := body["data"].(map[string]any)
	require.True(t, ok, "data must be an object")
	assert.Empty(t, data, "data must be empty object {}")
}

func TestPaged_RendersPaginatedEnvelope(t *testing.T) {
	c, rec := newCtx()
	page := pagination.Build(2, 10, 25)
	payload := []map[string]any{{"id": "1"}, {"id": "2"}}
	require.NoError(t, response.Paged(c, "listed", payload, page))

	assert.Equal(t, http.StatusOK, rec.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))

	assert.Equal(t, true, body["success"])
	assert.EqualValues(t, 200, body["code"])
	assert.Equal(t, "listed", body["message"])

	pg, ok := body["pagination"].(map[string]any)
	require.True(t, ok)
	assert.EqualValues(t, 2, pg["page"])
	assert.EqualValues(t, 10, pg["page_size"])
	assert.EqualValues(t, 25, pg["total_data"])
	assert.EqualValues(t, 3, pg["total_page"])

	data, ok := body["data"].([]any)
	require.True(t, ok)
	assert.Len(t, data, 2)
}

func TestErr_CustomError_RendersErrorEnvelope(t *testing.T) {
	c, rec := newCtx()
	err := customerror.ErrNotFound.WithMetadata(map[string]any{"detail": "item missing"})
	require.NoError(t, response.Err(c, err))

	assert.Equal(t, http.StatusNotFound, rec.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))

	assert.Equal(t, false, body["success"])
	assert.EqualValues(t, 40400, body["code"])
	assert.Equal(t, "not found", body["message"])
	assert.NotEmpty(t, body["timestamp"])
	meta, ok := body["metadata"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "item missing", meta["detail"])
}

func TestErr_RawError_FallsBackToInternalServer(t *testing.T) {
	c, rec := newCtx()
	require.NoError(t, response.Err(c, errors.New("kaboom")))

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))

	assert.Equal(t, false, body["success"])
	assert.EqualValues(t, 50000, body["code"])
	assert.Equal(t, "internal server error", body["message"])
}

func TestErr_NoMetadata_OmitsField(t *testing.T) {
	c, rec := newCtx()
	require.NoError(t, response.Err(c, customerror.ErrForbidden))

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	_, present := body["metadata"]
	assert.False(t, present, "metadata must be omitted when empty")
}
