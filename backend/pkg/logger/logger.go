package logger

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

type Logger struct{}

func New(isEnableDebug bool) *Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	log.Logger = zerolog.New(os.Stderr).
		With().
		Caller().
		Timestamp().
		Stack().
		Logger()

	if isEnableDebug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	return &Logger{}
}

func (l *Logger) ErrorContext(_ context.Context) *zerolog.Event {
	return log.Error()
}

func (l *Logger) InfoContext(_ context.Context) *zerolog.Event {
	return log.Info()
}

func (l *Logger) DebugContext(_ context.Context) *zerolog.Event {
	return log.Debug()
}

func (l *Logger) FatalContext(_ context.Context) *zerolog.Event {
	return log.Fatal()
}

func (l *Logger) WarnContext(_ context.Context) *zerolog.Event {
	return log.Warn()
}
