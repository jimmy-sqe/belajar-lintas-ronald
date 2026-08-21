package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/todo"
	todomock "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/todo/mock"
	tehttp "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/http"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/http/handler"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/pkg/ctxutil"
)

func newEcho(t *testing.T) *echo.Echo {
	t.Helper()
	e := echo.New()
	e.Validator = tehttp.NewValidator()
	return e
}

func decode(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	return body
}

func TestTodoHandler_Create_Success(t *testing.T) {
	svc := new(todomock.Service)
	ownerID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	todoID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	created := &todo.Todo{
		ID:        todoID,
		Title:     "Buy milk",
		CreatedBy: ownerID,
		ModifiedBy: ownerID,
		CreatedAt: time.Now(),
		ModifiedAt: time.Now(),
	}
	svc.On("Create", mock.Anything, mock.MatchedBy(func(in todo.CreateInput) bool {
		return in.Title == "Buy milk" && in.CreatedBy == ownerID
	})).Return(created, nil)

	h := handler.NewTodoHandler(svc)
	e := newEcho(t)

	body := `{"title":"Buy milk","description":"from store"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/todos", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req = req.WithContext(ctxutil.WithUserID(context.Background(), ownerID.String()))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, h.Create(c))
	require.Equal(t, http.StatusCreated, rec.Code)

	resp := decode(t, rec)
	assert.Equal(t, true, resp["success"])
	assert.EqualValues(t, 201, resp["code"])
	assert.Equal(t, "todo created", resp["message"])
	data := resp["data"].(map[string]any)
	assert.Equal(t, todoID.String(), data["id"])
	assert.Equal(t, "Buy milk", data["title"])
}

func TestTodoHandler_Create_ValidationFailure(t *testing.T) {
	svc := new(todomock.Service)
	h := handler.NewTodoHandler(svc)
	e := newEcho(t)

	body := `{"title":""}`
	req := httptest.NewRequest(http.MethodPost, "/v1/todos", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ownerID := uuid.New()
	req = req.WithContext(ctxutil.WithUserID(context.Background(), ownerID.String()))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, h.Create(c))
	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	resp := decode(t, rec)
	assert.Equal(t, false, resp["success"])
	assert.EqualValues(t, 42200, resp["code"])
	svc.AssertNotCalled(t, "Create")
}

func TestTodoHandler_Create_MissingUser(t *testing.T) {
	svc := new(todomock.Service)
	h := handler.NewTodoHandler(svc)
	e := newEcho(t)

	body := `{"title":"Buy milk"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/todos", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, h.Create(c))
	require.Equal(t, http.StatusUnauthorized, rec.Code)
	resp := decode(t, rec)
	assert.EqualValues(t, 40100, resp["code"])
	svc.AssertNotCalled(t, "Create")
}

func TestTodoHandler_GetByID_Success(t *testing.T) {
	svc := new(todomock.Service)
	ownerID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	todoID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	stored := &todo.Todo{
		ID:        todoID,
		Title:     "Buy milk",
		CreatedBy: ownerID,
		ModifiedBy: ownerID,
		CreatedAt: time.Now(),
		ModifiedAt: time.Now(),
	}
	svc.On("FindByID", mock.Anything, todoID, ownerID).Return(stored, nil)

	h := handler.NewTodoHandler(svc)
	e := newEcho(t)

	req := httptest.NewRequest(http.MethodGet, "/v1/todos/"+todoID.String(), nil)
	req = req.WithContext(ctxutil.WithUserID(context.Background(), ownerID.String()))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(todoID.String())

	require.NoError(t, h.GetByID(c))
	require.Equal(t, http.StatusOK, rec.Code)
	resp := decode(t, rec)
	assert.Equal(t, true, resp["success"])
	assert.EqualValues(t, 200, resp["code"])
	data := resp["data"].(map[string]any)
	assert.Equal(t, todoID.String(), data["id"])
}

func TestTodoHandler_GetByID_NotFound(t *testing.T) {
	svc := new(todomock.Service)
	ownerID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	todoID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	svc.On("FindByID", mock.Anything, todoID, ownerID).Return(nil, todo.ErrNotFound)

	h := handler.NewTodoHandler(svc)
	e := newEcho(t)

	req := httptest.NewRequest(http.MethodGet, "/v1/todos/"+todoID.String(), nil)
	req = req.WithContext(ctxutil.WithUserID(context.Background(), ownerID.String()))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(todoID.String())

	require.NoError(t, h.GetByID(c))
	require.Equal(t, http.StatusNotFound, rec.Code)
	resp := decode(t, rec)
	assert.Equal(t, false, resp["success"])
	assert.EqualValues(t, 40400, resp["code"])
	assert.Equal(t, "not found", resp["message"])
}

func TestTodoHandler_Update_Success(t *testing.T) {
	svc := new(todomock.Service)
	ownerID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	todoID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	updated := &todo.Todo{
		ID:        todoID,
		Title:     "Buy oat milk",
		CreatedBy: ownerID,
		ModifiedBy: ownerID,
		CreatedAt: time.Now(),
		ModifiedAt: time.Now(),
	}
	svc.On("Update", mock.Anything, todoID, mock.Anything).Return(updated, nil)

	h := handler.NewTodoHandler(svc)
	e := newEcho(t)

	req := httptest.NewRequest(http.MethodPut, "/v1/todos/"+todoID.String(),
		strings.NewReader(`{"title":"Buy oat milk"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req = req.WithContext(ctxutil.WithUserID(context.Background(), ownerID.String()))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(todoID.String())

	require.NoError(t, h.Update(c))
	require.Equal(t, http.StatusOK, rec.Code)
	resp := decode(t, rec)
	assert.Equal(t, "todo updated", resp["message"])
}

func TestTodoHandler_Delete_OKEnvelope(t *testing.T) {
	svc := new(todomock.Service)
	ownerID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	todoID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	svc.On("Delete", mock.Anything, todoID, ownerID).Return(nil)

	h := handler.NewTodoHandler(svc)
	e := newEcho(t)

	req := httptest.NewRequest(http.MethodDelete, "/v1/todos/"+todoID.String(), nil)
	req = req.WithContext(ctxutil.WithUserID(context.Background(), ownerID.String()))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(todoID.String())

	require.NoError(t, h.Delete(c))
	require.Equal(t, http.StatusOK, rec.Code)
	resp := decode(t, rec)
	assert.Equal(t, true, resp["success"])
	assert.EqualValues(t, 200, resp["code"])
	assert.Equal(t, "todo deleted", resp["message"])
	data, ok := resp["data"].(map[string]any)
	require.True(t, ok)
	assert.Empty(t, data)
}

func TestTodoHandler_List_Pagination(t *testing.T) {
	svc := new(todomock.Service)
	ownerID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	todoID := uuid.New()
	svc.On("List", mock.Anything, mock.MatchedBy(func(f todo.ListFilter) bool {
		// page=2, page_size=50 → offset 50, limit 50
		return f.Limit == 50 && f.Offset == 50 && f.OwnerID == ownerID
	})).Return([]todo.Todo{{ID: todoID, CreatedBy: ownerID}}, int64(125), nil)

	h := handler.NewTodoHandler(svc)
	e := newEcho(t)

	req := httptest.NewRequest(http.MethodGet, "/v1/todos?page=2&page_size=50", nil)
	req = req.WithContext(ctxutil.WithUserID(context.Background(), ownerID.String()))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.NoError(t, h.List(c))
	require.Equal(t, http.StatusOK, rec.Code)
	resp := decode(t, rec)

	assert.Equal(t, true, resp["success"])
	pg := resp["pagination"].(map[string]any)
	assert.EqualValues(t, 2, pg["page"])
	assert.EqualValues(t, 50, pg["page_size"])
	assert.EqualValues(t, 125, pg["total_data"])
	assert.EqualValues(t, 3, pg["total_page"])

	data := resp["data"].([]any)
	assert.Len(t, data, 1)
}
