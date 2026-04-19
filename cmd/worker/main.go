// Command worker es el binario de procesamiento de eventos del hub de inventario.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	log := logger.New(cfg.LogLevel, cfg.Environment)
	log.Info("starting worker",
		slog.String("version", version),
		slog.String("env", cfg.Environment),
	)

	// Mismo patrón de graceful shutdown que el API
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	log.Info("worker ready, waiting for ticks (fase 0 stub)")

	for {
		select {
		case <-ctx.Done():
			log.Info("worker stopping cleanly")
			return nil
		case t := <-ticker.C:
			log.Info("worker heartbeat", slog.Time("at", t))
		}
	}
}
