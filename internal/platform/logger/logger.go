// Package logger provee un logger estructurado compartido por toda la app.
package logger

import (
	"context"
	"log/slog"
	"os"
)

type ctxKey struct{ name string }

var traceIDKey = ctxKey{"trace_id"}

// New construye un *slog.Logger configurado según nivel y entorno
// - development: handler de texto legible
// - cualquier otro: handler JSON estructurado
func New(level, env string) *slog.Logger {
	var lvl slog.Level
	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     lvl,
		AddSource: env == "development",
	}

	var handler slog.Handler
	if env == "development" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

// WithTraceID devuelve un context que lleva el trace ID
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

// TraceIDFrom extrae el trace ID del context, o "" si no hay.
func TraceIDFrom(ctx context.Context) string {
	if v, ok := ctx.Value(traceIDKey).(string); ok {
		return v
	}
	return ""
}
