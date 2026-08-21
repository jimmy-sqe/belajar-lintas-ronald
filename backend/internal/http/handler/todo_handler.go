package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/todo"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/http/handler/dto"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/http/response"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/pkg/ctxutil"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/pkg/customerror"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/pkg/pagination"
)

// TodoHandler is the HTTP-facing todo API.
type TodoHandler struct {
	svc todo.Service
}

// NewTodoHandler constructs a TodoHandler.
func NewTodoHandler(svc todo.Service) *TodoHandler {
	return &TodoHandler{svc: svc}
}

func (h *TodoHandler) currentUserID(c echo.Context) (uuid.UUID, error) {
	uid := ctxutil.UserID(c.Request().Context())
	if uid == "" {
		return uuid.Nil, customerror.ErrUnauthorized
	}
	return uuid.Parse(uid)
}

// Create handles POST /v1/todos.
func (h *TodoHandler) Create(c echo.Context) error {
	var req dto.CreateTodoRequest
	if err := c.Bind(&req); err != nil {
		return response.Err(c, customerror.ErrBadRequest.WithMetadata(
			map[string]any{"detail": err.Error()}))
	}
	if err := c.Validate(&req); err != nil {
		return response.Err(c, customerror.ErrUnprocessable.WithMetadata(
			map[string]any{"detail": err.Error()}))
	}
	userID, err := h.currentUserID(c)
	if err != nil {
		return response.Err(c, customerror.ErrUnauthorized)
	}
	t, err := h.svc.Create(c.Request().Context(), todo.CreateInput{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		CreatedBy:   userID,
	})
	if err != nil {
		return response.Err(c, todoCustomError(err))
	}
	return response.OK(c, http.StatusCreated, "todo created", todoToResponse(t))
}

// GetByID handles GET /v1/todos/:id.
func (h *TodoHandler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Err(c, customerror.ErrBadRequest.WithMetadata(
			map[string]any{"detail": "invalid id"}))
	}
	userID, err := h.currentUserID(c)
	if err != nil {
		return response.Err(c, customerror.ErrUnauthorized)
	}
	t, err := h.svc.FindByID(c.Request().Context(), id, userID)
	if err != nil {
		return response.Err(c, todoCustomError(err))
	}
	return response.OK(c, http.StatusOK, "todo fetched", todoToResponse(t))
}

// maxPageSize bounds the list page size. It matches the OpenAPI contract
// (`page_size.maximum`); enforcing it here keeps the implementation in sync
// with the advertised limit and caps the SQL LIMIT so a crafted large
// page_size can't exhaust memory or overflow int on the uint64→int cast.
const maxPageSize = 100

// List handles GET /v1/todos.
// Query params: page (default 1), page_size (default 10, max 100).
func (h *TodoHandler) List(c echo.Context) error {
	userID, err := h.currentUserID(c)
	if err != nil {
		return response.Err(c, customerror.ErrUnauthorized)
	}
	page, _ := strconv.ParseUint(c.QueryParam("page"), 10, 64)
	pageSize, _ := strconv.ParseUint(c.QueryParam("page_size"), 10, 64)
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	f := todo.ListFilter{
		OwnerID: userID,
		Limit:   int(pageSize),
		Offset:  int(pagination.SQLOffset(page, pageSize)),
	}
	todos, total, err := h.svc.List(c.Request().Context(), f)
	if err != nil {
		return response.Err(c, todoCustomError(err))
	}
	payload := make([]dto.TodoResponse, 0, len(todos))
	for i := range todos {
		payload = append(payload, *todoToResponse(&todos[i]))
	}
	pg := pagination.Build(page, pageSize, uint64(total))
	return response.Paged(c, "todos fetched", payload, pg)
}

// Update handles PUT /v1/todos/:id.
func (h *TodoHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Err(c, customerror.ErrBadRequest.WithMetadata(
			map[string]any{"detail": "invalid id"}))
	}
	var req dto.UpdateTodoRequest
	if err := c.Bind(&req); err != nil {
		return response.Err(c, customerror.ErrBadRequest.WithMetadata(
			map[string]any{"detail": err.Error()}))
	}
	if err := c.Validate(&req); err != nil {
		return response.Err(c, customerror.ErrUnprocessable.WithMetadata(
			map[string]any{"detail": err.Error()}))
	}
	userID, err := h.currentUserID(c)
	if err != nil {
		return response.Err(c, customerror.ErrUnauthorized)
	}
	var dueDatePtr **time.Time
	if req.DueDate != nil {
		dueDatePtr = &req.DueDate
	}
	in := todo.UpdateInput{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     dueDatePtr,
		ModifiedBy:  userID,
	}
	t, err := h.svc.Update(c.Request().Context(), id, in)
	if err != nil {
		return response.Err(c, todoCustomError(err))
	}
	return response.OK(c, http.StatusOK, "todo updated", todoToResponse(t))
}

// Delete handles DELETE /v1/todos/:id.
func (h *TodoHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Err(c, customerror.ErrBadRequest.WithMetadata(
			map[string]any{"detail": "invalid id"}))
	}
	userID, err := h.currentUserID(c)
	if err != nil {
		return response.Err(c, customerror.ErrUnauthorized)
	}
	if err := h.svc.Delete(c.Request().Context(), id, userID); err != nil {
		return response.Err(c, todoCustomError(err))
	}
	return response.OK(c, http.StatusOK, "todo deleted", nil)
}

func todoToResponse(t *todo.Todo) *dto.TodoResponse {
	return &dto.TodoResponse{
		ID:          t.ID.String(),
		Title:       t.Title,
		Description: t.Description,
		DueDate:     t.DueDate,
		CreatedAt:   t.CreatedAt,
		CreatedBy:   t.CreatedBy.String(),
		ModifiedAt:  t.ModifiedAt,
		ModifiedBy:  t.ModifiedBy.String(),
	}
}

// todoCustomError maps domain sentinels to generic CustomErrors.
func todoCustomError(err error) customerror.CustomError {
	switch {
	case errors.Is(err, todo.ErrNotFound):
		return customerror.ErrNotFound.WithMetadata(
			map[string]any{"detail": err.Error()})
	case errors.Is(err, todo.ErrForbidden):
		return customerror.ErrForbidden
	}
	var ce customerror.CustomError
	if errors.As(err, &ce) {
		return ce
	}
	return customerror.ErrInternalServer
}
