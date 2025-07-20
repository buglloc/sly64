package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/buglloc/sly64/v2/internal/config"
)

var startCmd = &cobra.Command{
	Use:           "start",
	SilenceUsage:  true,
	SilenceErrors: true,
	Short:         "Starts sly64 in foreground",
	RunE: func(_ *cobra.Command, _ []string) error {
		srv, err := runtime.NewServer()
		if err != nil {
			return fmt.Errorf("create server: %w", err)
		}

		errChan := make(chan error)
		go func() {
			defer close(errChan)

			err := srv.Start()
			if err != nil {
				errChan <- err
			}
		}()

		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-stopChan:
			log.Info().Stringer("signal", sig).Msg("shutting down by signal")

			ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownDeadline)
			defer cancel()

			if err := srv.Shutdown(ctx); err != nil {
				log.Info().Err(err).Msg("shutdown failed")
			}

		case err := <-errChan:
			return fmt.Errorf("unable to start server: %w", err)
		}

		return nil
	},
}
