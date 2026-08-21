package cron

import (
	"context"
	"log"
)

// Example shows how to wire a cron Runner. Copy/adapt this in your service.
func Example() {
	r := New()
	if err := r.Schedule("@every 1m", func(_ context.Context) {
		log.Println("cron: tick")
	}); err != nil {
		log.Fatalf("cron schedule: %v", err)
	}
	r.Start()
	// later: r.Stop(ctx)
}
