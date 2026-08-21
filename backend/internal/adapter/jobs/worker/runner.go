// Package worker is a queue-driven worker pool adapter. Pair with a
// messaging adapter to feed the in channel.
package worker

import (
	"context"
	"sync"
)

// Job is a unit of work consumed by workers.
type Job struct {
	ID      string
	Payload []byte
}

// Handler processes a single Job. Implementations should be idempotent.
type Handler interface {
	Handle(ctx context.Context, job Job) error
}

// Runner runs N goroutines pulling Job values from a channel.
type Runner struct {
	in      <-chan Job
	handler Handler
	workers int
	wg      sync.WaitGroup
	cancel  context.CancelFunc
}

// New constructs a Runner. workers <= 0 falls back to 1.
func New(in <-chan Job, handler Handler, workers int) *Runner {
	if workers <= 0 {
		workers = 1
	}
	return &Runner{in: in, handler: handler, workers: workers}
}

// Start launches the worker goroutines. Returns immediately.
func (r *Runner) Start(parent context.Context) {
	ctx, cancel := context.WithCancel(parent)
	r.cancel = cancel
	for i := 0; i < r.workers; i++ {
		r.wg.Add(1)
		go r.loop(ctx)
	}
}

func (r *Runner) loop(ctx context.Context) {
	defer r.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-r.in:
			if !ok {
				return
			}
			_ = r.handler.Handle(ctx, job)
		}
	}
}

// Stop signals shutdown and waits for in-flight handlers.
func (r *Runner) Stop(_ context.Context) error {
	if r.cancel != nil {
		r.cancel()
	}
	r.wg.Wait()
	return nil
}
