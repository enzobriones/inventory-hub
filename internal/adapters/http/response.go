package http

import (
	"encoding/json"
	"log/slog"
	nethttp "net/http"
)

// WriteJSON serializa v como JSON y lo escribe en w con el status dado.
// Si el encode falla, lo loggea — a esa altura ya mandamos el header,
// no podemos hacer mucho más.
func WriteJSON(w nethttp.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("failed to encode json response", slog.String("error", err.Error()))
	}
}

// ErrorResponse es el shape estándar de errores de la API.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// WriteError responde con un ErrorResponse en formato JSON.
func WriteError(w nethttp.ResponseWriter, status int, code, message string) {
	WriteJSON(w, status, ErrorResponse{Error: code, Message: message})
}
