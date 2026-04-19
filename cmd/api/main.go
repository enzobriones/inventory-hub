// Command api es el servidor HTTP del hub de inventario.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	nethttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apphttp "github.com/enzobriones/inventory-hub/internal/adapters/http"
	"github.com/enzobriones/inventory-hub/internal/config"
	"github.com/enzobriones/inventory-hub/internal/platform/logger"
)

var version = "dev"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// 1. Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// 2. Construir logger
	log := logger.New(cfg.LogLevel, cfg.Environment)
	log.Info("starting api",
		slog.String("version", version),
		slog.String("env", cfg.Environment),
		slog.Int("port", cfg.HTTPPort),
	)

	// 3. Armar el router
	mux := nethttp.NewServeMux()
	mux.HandleFunc("GET /health", apphttp.HealthHandler(version))

	// 4. Aplicar middleware (orden: el más externo envuelve al resto)
	// 		TraceID corre primero -> inyecta el ID en el context
	// 		Logging corre después -> lee el ID del context
	handler := apphttp.TraceIDMiddleware(
		apphttp.LoggingMiddleware(log)(mux),
	)

	// 5. Configurar servidor con timeouts explícitos
	srv := &nethttp.Server{
		Addr:              fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// 6. Preparar graceful shutdown: escuchamos SIGINT (Ctrl+C) y SIGTERM (docker/k8s stop)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 7. Arrancar el servidor en una goroutine (ListenAndServe bloquea)
	errCh := make(chan error, 1)
	go func() {
		log.Info("http server listening", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, nethttp.ErrServerClosed) {
			errCh <- err
		}
	}()

	// 8. Esperar: o bien llega una señal, o bien el servidor falla
	select {
	case <-ctx.Done():
		log.Info("shutdown signal received")
	case err := <-errCh:
		return fmt.Errorf("server error: %w", err)
	}

	// 9. Shutdown con timeout: dejamos 10s para que terminen las respuestas en vuelo
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	log.Info("server stopped cleanly")
	return nil
}
