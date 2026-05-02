package logger

import (
	"errors"
	"log/slog"
	"os"

	"github.com/identicalaffiliation/oms-with-events/migrator/internal/config"
)

const (
	DEBUG = "debug"
	ERROR = "error"
	INFO  = "info"
	TEXT  = "text"
	JSON  = "json"
)

type Logger interface {
	Debug(msg string, params ...any)
	Info(msg string, params ...any)
	Error(msg string, params ...any)
}

type slogger struct {
	l *slog.Logger
}

func NewSLogger(cfg *config.OMSMigratorConfig) (Logger, error) {
	levels := map[string]slog.Level{
		ERROR: slog.LevelError,
		DEBUG: slog.LevelDebug,
		INFO:  slog.LevelInfo,
	}

	l, ok := levels[cfg.LoggerConfig.Level]
	if !ok {
		return nil, errors.New("invalid log level")
	}

	handlers := map[string]slog.Handler{
		JSON: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: l}),
		TEXT: slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: l}),
	}

	h, ok := handlers[cfg.LoggerConfig.Format]
	if !ok {
		return nil, errors.New("invalid log format")
	}

	return &slogger{l: slog.New(h)}, nil
}

func (l *slogger) Debug(msg string, params ...any) {
	l.l.Debug(msg, params...)
}

func (l *slogger) Info(msg string, params ...any) {
	l.l.Error(msg, params...)
}

func (l *slogger) Error(msg string, params ...any) {
	l.l.Error(msg, params...)
}
