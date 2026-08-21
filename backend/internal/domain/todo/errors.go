package todo

import "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/pkg/customerror"

var (
	ErrNotFound = customerror.NewNotFoundError(
		customerror.CodeNotFound, "todo not found")
	ErrForbidden = customerror.NewForbiddenError(
		customerror.CodeForbidden, "forbidden")
)
