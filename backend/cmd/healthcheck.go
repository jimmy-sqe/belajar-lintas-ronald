package cmd

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// healthcheckCmd performs an in-process liveness probe against the running
// server's /healthz endpoint. It exists so container healthchecks work on the
// distroless final image, which ships no shell, wget, or curl. Reads HTTP_PORT
// (default 8000) directly rather than the full config to stay dependency-free.
func healthcheckCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "healthcheck",
		Short: "Probe the local /healthz endpoint (exit 0 healthy, 1 unhealthy)",
		RunE: func(_ *cobra.Command, _ []string) error {
			port := os.Getenv("HTTP_PORT")
			if port == "" {
				port = "8000"
			}
			client := &http.Client{Timeout: 3 * time.Second}
			resp, err := client.Get(fmt.Sprintf("http://127.0.0.1:%s/healthz", port))
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("healthcheck: unexpected status %d", resp.StatusCode)
			}
			return nil
		},
	}
}
