package cmd

import (
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/app"
	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/config"
)

func serveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the HTTP server",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			ctx, cancel := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
			defer cancel()

			c, err := app.NewContainer(ctx, cfg)
			if err != nil {
				return err
			}
			defer c.Close() //nolint:errcheck // best-effort container teardown

			return app.NewServer(c).Start(ctx)
		},
	}
}
