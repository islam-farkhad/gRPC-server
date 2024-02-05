package logger

import (
	"context"

	"go.uber.org/zap"
)

var defaultLogger *zap.Logger

type ctxKey struct{}

// SetGlobal sets the provided logger as the defaultLogger.
func SetGlobal(logger *zap.Logger) {
	defaultLogger = logger
}

// FromContext fetches a logger from the provided context or returns the default logger.
func FromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return logger
	}
	return defaultLogger
}

// ToContext sets a logger in the context and returns a new context with the logger.
func ToContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

// Infof logs message with Info level using the logger from the context.
func Infof(ctx context.Context, format string, args ...any) {
	FromContext(ctx).Sugar().Infof(format, args...)
}

// Errorf logs message with Error level using the logger from the context.
func Errorf(ctx context.Context, format string, args ...any) {
	FromContext(ctx).Sugar().Errorf(format, args...)
}
