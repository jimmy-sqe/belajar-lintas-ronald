// Package dto holds request/response shapes for HTTP handlers. Validator
// tags are checked by the validator middleware.
package dto

import "time"

type CreateTodoRequest struct {
	Title       string     `json:"title" validate:"required,min=1,max=200"`
	Description string     `json:"description" validate:"max=2000"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

type UpdateTodoRequest struct {
	Title       *string    `json:"title,omitempty" validate:"omitempty,min=1,max=200"`
	Description *string    `json:"description,omitempty" validate:"omitempty,max=2000"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

type TodoResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	CreatedBy   string     `json:"created_by"`
	ModifiedAt  time.Time  `json:"modified_at"`
	ModifiedBy  string     `json:"modified_by"`
}
