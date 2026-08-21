// Package cron is the in-process cron job adapter using robfig/cron/v3.
package cron

import (
	"context"
	"fmt"

	rcron "github.com/robfig/cron/v3"
)

// Runner schedules and runs cron jobs in-process.
type Runner struct {
	c      *rcron.Cron
	ctx    context.Context
	cancel context.CancelFunc
}

// New constructs a Runner. Time zone defaults to UTC.
func New() *Runner {
	ctx, cancel := context.WithCancel(context.Background())
	return &Runner{c: rcron.New(rcron.WithSeconds()), ctx: ctx, cancel: cancel}
}

// Schedule registers fn under spec (e.g., "0 * * * * *" for every minute).
// fn receives a derived context that is canceled when Stop() is called, so
// long-running jobs can observe shutdown and wind down promptly.
func (r *Runner) Schedule(spec string, fn func(ctx context.Context)) error {
	if _, err := r.c.AddFunc(spec, func() { fn(r.ctx) }); err != nil {
		return fmt.Errorf("cron: schedule %q: %w", spec, err)
	}
	return nil
}

// Start begins running scheduled jobs. Non-blocking.
func (r *Runner) Start() { r.c.Start() }

// Stop halts all jobs and waits for in-flight ones to finish. The job context
// is canceled first so in-flight handlers can abort, then we drain.
func (r *Runner) Stop(_ context.Context) error {
	r.cancel()
	stopCtx := r.c.Stop()
	<-stopCtx.Done()
	return nil
}
