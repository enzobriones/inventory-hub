package http

import (
	nethttp "net/http"
	"time"
)

// HealthResponse es el payload del endpoint /health.
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

// HealthHandler devuelve un handler que responde con el estado básico del servicio.
// TODO: En fases posteriores, vamos a agregar checks de Postgres, Redis y NATS.
func HealthHandler(version string) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		WriteJSON(w, nethttp.StatusOK, HealthResponse{
			Status:    "ok",
			Timestamp: time.Now().UTC(),
			Version:   version,
		})
	}
}
