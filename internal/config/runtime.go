package config

import (
	"github.com/rs/zerolog"

	"github.com/buglloc/sly64/v2/internal/config/configpb"
)

type Runtime struct {
	Config *configpb.Config
}

func newRuntime(cfg *configpb.Config) (*Runtime, error) {
	r := &Runtime{
		Config: cfg,
	}

	zerolog.SetGlobalLevel(protoLogLevelToZero(r.Config.LogLevel))
	return r, nil
}

func protoLogLevelToZero(lvl configpb.LogLevel) zerolog.Level {
	switch lvl {
	case configpb.LogLevel_LOG_LEVEL_DEBUG:
		return zerolog.DebugLevel

	case configpb.LogLevel_LOG_LEVEL_INFO:
		return zerolog.InfoLevel

	case configpb.LogLevel_LOG_LEVEL_WARN:
		return zerolog.WarnLevel

	case configpb.LogLevel_LOG_LEVEL_ERROR:
		return zerolog.ErrorLevel

	case configpb.LogLevel_LOG_LEVEL_FATAL:
		return zerolog.FatalLevel

	case configpb.LogLevel_LOG_LEVEL_PANIC:
		return zerolog.PanicLevel

	case configpb.LogLevel_LOG_LEVEL_DISABLED:
		return zerolog.Disabled

	default:
		return zerolog.InfoLevel
	}
}
