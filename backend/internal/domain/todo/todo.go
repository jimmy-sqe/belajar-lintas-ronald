// Package todo manages todo items owned by authenticated users.
package todo

import (
	"time"

	"github.com/google/uuid"
)

// Todo is the canonical todo entity. CreatedBy is the ownership marker;
// queries filter on it to enforce per-user isolation.
type Todo struct {
	ID          uuid.UUID
	Title       string
	Description string
	DueDate     *time.Time
	CreatedAt   time.Time
	CreatedBy   uuid.UUID
	ModifiedAt  time.Time
	ModifiedBy  uuid.UUID
}

// CreateInput is what handlers pass to Service.Create.
type CreateInput struct {
	Title       string
	Description string
	DueDate     *time.Time
	CreatedBy   uuid.UUID
}

// UpdateInput is what handlers pass to Service.Update. Title and
// Description are nil-pointer to mean "no change". DueDate uses
// **time.Time to distinguish "no change" (nil) from "clear value"
// (*time.Time = nil inside outer pointer) — simplest distinction.
type UpdateInput struct {
	Title       *string
	Description *string
	DueDate     **time.Time
	ModifiedBy  uuid.UUID
}

// ListFilter is what handlers pass to Service.List.
type ListFilter struct {
	OwnerID uuid.UUID
	Limit   int
	Offset  int
}
