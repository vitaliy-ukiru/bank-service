package logging

import (
	"context"
	"log/slog"
)

type loggerKey struct{}

func FromContext(ctx context.Context) Logger {
	v, ok := ctx.Value(loggerKey{}).(Logger)
	if !ok {
		return Logger{slog.Default()}
	}
	return v
}

func Context(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}
