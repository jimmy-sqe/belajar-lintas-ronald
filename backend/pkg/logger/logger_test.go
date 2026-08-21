package logger

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	ctx := t.Context()
	logger := New(true)

	errCtx := logger.ErrorContext(ctx)
	assert.IsType(t, &zerolog.Event{}, errCtx)

	infoCtx := logger.InfoContext(ctx)
	assert.IsType(t, &zerolog.Event{}, infoCtx)

	debugCtx := logger.DebugContext(ctx)
	assert.IsType(t, &zerolog.Event{}, debugCtx)

	fatalCtx := logger.FatalContext(ctx)
	assert.IsType(t, &zerolog.Event{}, fatalCtx)

	warnCtx := logger.WarnContext(ctx)
	assert.IsType(t, &zerolog.Event{}, warnCtx)
}
