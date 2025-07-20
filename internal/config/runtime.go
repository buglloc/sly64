package config

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Runtime struct {
	Config     *Config
	toShutdown []shutdowner
}

func newRuntime(c *Config) (*Runtime, error) {
	r := &Runtime{
		Config: c,
	}

	zerolog.SetGlobalLevel(r.Config.LogLevel)
	return r, nil
}

func (r *Runtime) Shutdown(ctx context.Context) {
	for _, s := range r.toShutdown {
		if err := s.Shutdown(ctx); err != nil {
			log.Warn().
				Err(err).
				Str("shutdowner", fmt.Sprintf("%T", s)).
				Msg("unable to shutdown")
		}
	}
}

type shutdowner interface {
	Shutdown(ctx context.Context) error
}
