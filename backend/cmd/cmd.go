package cmd

import "github.com/spf13/cobra"

func Root() *cobra.Command {
	root := &cobra.Command{
		Use:   "backend-belajar-lintas-ronald",
		Short: "Template backend service",
	}
	root.AddCommand(serveCmd())
	root.AddCommand(migrateCmd())
	root.AddCommand(seedCmd())
	root.AddCommand(healthcheckCmd())
	return root
}
