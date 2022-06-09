package server

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/uclatall/ckhub/pkg/logging"
	"github.com/uclatall/ckhub/pkg/runtime"
	"github.com/uclatall/ckhub/sandbox"
	"github.com/uclatall/ckhub/sandbox/server"
)

// NewCommand creates a new server management command.
func NewCommand(version string) *cobra.Command {
	flags := NewFlags()

	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "server [flags]",
		Short: "Start a playground server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			cmd.SilenceErrors = true

			log := logging.NewLogger(flags.Logger())
			log.Info("starting server", logging.String("version", version))

			cfg, err := NewConfig(flags.Config())
			if err != nil {
				log.Error("server interrupted", logging.Error(err))
				return err
			}

			mgr, err := sandbox.NewManager(
				sandbox.Logger(log.Name("sandbox")),
				cfg.Sandbox,
			)

			srv, err := server.NewServer(
				mgr,
				server.Logger(log.Name("sandbox")),
				cfg.Server,
			)
			if err != nil {
				log.Error("server interrupted", logging.Error(err))
				return err
			}

			run, err := runtime.NewRuntime(
				runtime.Logger(log),
				runtime.Group{mgr, srv},
			)
			if err != nil {
				log.Error("server interrupted", logging.Error(err))
				return fmt.Errorf("server interrupted: %w", err)
			}

			err = run.Run(cmd.Context())
			if err != nil {
				log.Error("server interrupted", logging.Error(err))
				return fmt.Errorf("server interrupted: %w", err)
			}

			log.Info("server stopped")
			return nil
		},
	}

	flags.Register(cmd.Flags())
	return cmd
}
