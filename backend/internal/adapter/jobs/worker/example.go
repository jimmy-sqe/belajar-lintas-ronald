package worker

import (
	"context"
	"log"
)

// LogHandler is a trivial Handler that logs each job it receives.
type LogHandler struct{}

// Handle implements Handler.
func (LogHandler) Handle(_ context.Context, j Job) error {
	log.Printf("worker: job=%s payload=%d bytes", j.ID, len(j.Payload))
	return nil
}

// Example shows wiring a Runner against a manual channel.
func Example() {
	in := make(chan Job, 16)
	r := New(in, LogHandler{}, 4)
	r.Start(context.Background())
	in <- Job{ID: "demo", Payload: []byte("hello")}
	close(in)
	_ = r.Stop(context.Background())
}
