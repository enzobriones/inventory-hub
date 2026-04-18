package http

import (
	"log/slog"
	nethttp "net/http"
	"time"

	"github.com/google/uuid"

	"github.com/enzobriones/inventory-hub/internal/platform/logger"
)

// statusRecorder envuelve un ResponseWriter para capturar el status code escrito.
// net/http no expone el status después de WriteHeader, así que lo guardamos nosotros
type statusRecorder struct {
	nethttp.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// TraceIDMiddleware inyecta un trace ID en el content de cada request.
// Si el cliente mandó uno en el header X-Trace-ID lo respetamos; si no, generamos uno.
// En ambos casos lo devolvemos en la respuesta para que el cliente pueda correlacionar.
func TraceIDMiddleware(next nethttp.Handler) nethttp.Handler {
	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		traceID := r.Header.Get("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.NewString()
		}
		ctx := logger.WithTraceID(r.Context(), traceID)
		w.Header().Set("X-Trace-ID", traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// LoggingMiddleware loggea cada request con método, path, status y duración.
// Se construye con un logger específico para poder testearlo con un buffer.
func LoggingMiddleware(log *slog.Logger) func(nethttp.Handler) nethttp.Handler {
	return func(next nethttp.Handler) nethttp.Handler {
		return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: nethttp.StatusOK}

			next.ServeHTTP(rec, r)

			log.InfoContext(r.Context(), "http_request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", rec.status),
				slog.Duration("duration", time.Since(start)),
				slog.String("trace_id", logger.TraceIDFrom(r.Context())),
			)
		})
	}
}
