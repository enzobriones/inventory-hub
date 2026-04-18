package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config contiene toda la configuración de la aplicación
// Se carga desde variables de entorno al arrancar.
type Config struct {
	HTTPPort     int
	LogLevel     string
	Environment  string
	DatabaseURL  string
	RedisURL     string
	NATSURL      string
	OTLPEndpoint string
	ServiceName  string
}

// Load lee las variables de entorno y construye una Config.
// Si falta una variable obligatoria, retorna error.
func Load() (*Config, error) {
	port, err := strconv.Atoi(getEnv("HTTP_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("Invalid HTTP_PORT: %w", err)
	}

	dbURL, err := requireEnv("DATABASE_URL")
	if err != nil {
		return nil, err
	}
	redisURL, err := requireEnv("REDIS_URL")
	if err != nil {
		return nil, err
	}
	natsURL, err := requireEnv("NATS_URL")
	if err != nil {
		return nil, err
	}

	return &Config{
		HTTPPort:     port,
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		Environment:  getEnv("ENVIRONMENT", "development"),
		DatabaseURL:  dbURL,
		RedisURL:     redisURL,
		NATSURL:      natsURL,
		OTLPEndpoint: getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", ""),
		ServiceName:  getEnv("OTEL_SERVICE_NAME", "inventory-api"),
	}, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func requireEnv(key string) (string, error) {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return "", fmt.Errorf("missing required env var: %s", key)
	}
	return v, nil
}
