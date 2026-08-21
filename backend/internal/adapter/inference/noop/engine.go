// Package noop is the no-op inference adapter. Infer loads no models and
// returns an empty Result echoing the requested ModelType. It is the
// functional default for the inference axis when no local ML runtime is
// wired. Mirrors the java NoopEngine and the domain mock.
package noop

import (
	"context"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/inference"
)

// Engine is the no-op implementation of inference.Engine.
type Engine struct{}

// New returns a no-op inference Engine.
func New() inference.Engine { return &Engine{} }

// compile-time check: Engine must implement inference.Engine.
var _ inference.Engine = (*Engine)(nil)

// Infer returns an empty Result carrying the requested modelType.
func (e *Engine) Infer(_ context.Context, modelType inference.ModelType, _ inference.Input) (*inference.Result, error) {
	return &inference.Result{ModelType: modelType}, nil
}
