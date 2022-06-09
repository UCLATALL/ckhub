package app

import (
	"github.com/spf13/cobra"

	"github.com/uclatall/ckhub/cmd/ckhub/app/server"
)

// NewCommand creates a new root command for ckhub application.
func NewCommand(version string) *cobra.Command {
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "play [command] [flags]",
		Short:   "Playground service for ckhub.",
		Version: version,
	}

	cmd.AddCommand(
		server.NewCommand(version),
	)

	return cmd
}
