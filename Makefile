.PHONY: help up down logs ps restart build run-api run-worker test lint fmt tidy clean

# Carga variables del .env si existe (necesario para run-api y run-worker)
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

help: ## Muestra esta ayuda
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

## --- Infra ---

up: ## Levanta toda la infra local (postgres, redis, nats, jaeger)
	docker compose -f deployments/docker-compose.yml up -d

down: ## Detiene la infra local
	docker compose -f deployments/docker-compose.yml down

restart: ## Reinicia la infra local
	$(MAKE) down
	$(MAKE) up

logs: ## Sigue los logs de la infra
	docker compose -f deployments/docker-compose.yml logs -f

ps: ## Muestra estado de los contenedores
	docker compose -f deployments/docker-compose.yml ps

## --- Go ---

build: ## Compila api y worker a ./bin/
	@mkdir -p bin
	go build -o bin/api ./cmd/api
	go build -o bin/worker ./cmd/worker
	@echo "binaries built in ./bin/"

run-api: ## Corre el API (lee .env)
	go run ./cmd/api

run-worker: ## Corre el worker (lee .env)
	go run ./cmd/worker

test: ## Corre todos los tests con race detector
	go test -race -count=1 ./...

lint: ## Corre golangci-lint (requiere instalarlo)
	golangci-lint run ./...

fmt: ## Formatea el código
	gofmt -s -w .
	go mod tidy

tidy: ## go mod tidy
	go mod tidy

clean: ## Elimina binarios compilados
	rm -rf bin/
