// Package mock provides a test double for the inference.Engine port.
package mock

import (
	"context"
	"sync"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/inference"
)

// Call records one Infer invocation for later assertion.
type Call struct {
	ModelType inference.ModelType
	Input     inference.Input
}

// Engine is a thread-safe test double for inference.Engine.
// It records every Infer call and optionally delegates to a stub function.
type Engine struct {
	mu    sync.Mutex
	calls []Call
	stub  func(inference.ModelType, inference.Input) (*inference.Result, error)
}

// compile-time check: Engine must implement inference.Engine.
var _ inference.Engine = (*Engine)(nil)

// Infer records the call and invokes the stub if set, else returns an empty Result.
func (m *Engine) Infer(_ context.Context, modelType inference.ModelType, input inference.Input) (*inference.Result, error) {
	m.mu.Lock()
	m.calls = append(m.calls, Call{ModelType: modelType, Input: input})
	m.mu.Unlock()
	if m.stub != nil {
		return m.stub(modelType, input)
	}
	return &inference.Result{ModelType: modelType}, nil
}

// SetStub replaces the default no-op response with a custom function.
func (m *Engine) SetStub(fn func(inference.ModelType, inference.Input) (*inference.Result, error)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stub = fn
}

// Calls returns a snapshot of all recorded Infer invocations.
func (m *Engine) Calls() []Call {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]Call(nil), m.calls...)
}

// Reset clears the recorded call log. Use between sub-tests sharing a mock instance.
func (m *Engine) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = m.calls[:0]
}
