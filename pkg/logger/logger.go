package logger

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

const componentKey = "component"

func init() {
	// setup global logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

type Logger struct {
	zerolog.Logger
}

func NewGlobal(cfg Config) {
	zl := New(cfg).Logger
	log.Logger = zl
	zerolog.DefaultContextLogger = &zl
	zerolog.SetGlobalLevel(zl.GetLevel())
	log.Logger.Debug().Msg("Running in verbose mode")
}

// New constructor
func New(cfg Config) Logger {
	logLevel := zerolog.InfoLevel
	if cfg.Verbose {
		logLevel = zerolog.DebugLevel
	}
	zl := log.Logger.Level(logLevel)
	if cfg.Pretty {
		zl = zl.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}

	return Logger{zl}
}

// WithComponent creates child logger for named component
func (l Logger) WithComponent(name string) Logger {
	return Logger{Logger: l.With().Str(componentKey, name).Logger()}
}

// Ctx gets or creates context logger
func Ctx(ctx context.Context) Logger {
	logger := zerolog.Ctx(ctx)
	return Logger{Logger: *logger}
}

// Global returns current global logger
func Global() *Logger {
	return &Logger{log.Logger}
}
