# Inventory Hub

Hub de sincronización de inventario multicanal para PyMEs chilenas.

Una fuente única de verdad para el stock, con sincronización automática hacia
múltiples canales de venta (POS, Mercado Libre, Shopify, WhatsApp Business) y
emisión de documentos tributarios electrónicos (DTE) vía integración con SII.

> 🚧 **Estado**: en desarrollo activo — Fase 0 completada, Fase 1 en curso (~55%).
> Contexto completo en [`docs/HANDOFF.md`](./docs/HANDOFF.md).

## Stack

- **Lenguaje**: Go 1.22+
- **Base de datos**: PostgreSQL 16
- **Cache y locks**: Redis 7
- **Broker de eventos**: NATS con JetStream
- **Observabilidad**: OpenTelemetry + Jaeger
- **Arquitectura**: hexagonal (ports & adapters)

## Requisitos

- Go 1.22 o superior
- Docker y Docker Compose
- `make`

## Setup inicial

```bash
# 1. Clonar el repo
git clone https://github.com/enzobriones/inventory-hub.git
cd inventory-hub

# 2. Copiar el archivo de variables de entorno
cp .env.example .env

# 3. Levantar la infra local (postgres, redis, nats, jaeger)
make up

# 4. Verificar que todo esté sano
make ps

# 5. Correr el API
make run-api
```

El API queda escuchando en `http://localhost:8080`.

Probá el health check:

```bash
curl -i http://localhost:8080/health
```

## Comandos útiles

```bash
make help         # lista todos los comandos disponibles
make up           # levanta la infra
make down         # detiene la infra
make run-api      # corre el API
make run-worker   # corre el worker
make build        # compila ambos binarios a ./bin/
make test         # corre tests con race detector
make lint         # corre linters
make fmt          # formatea código
```

## Estructura del proyecto
```
├── cmd/
│   ├── api/              # entry point del servidor HTTP
│   └── worker/           # entry point del worker de eventos
├── internal/
│   ├── domain/           # entidades y reglas de negocio puras
│   ├── ports/            # interfaces (puertos hexagonales)
│   ├── adapters/         # implementaciones concretas
│   │   ├── http/         # handlers HTTP
│   │   ├── postgres/     # repositorios SQL
│   │   ├── redis/        # cache y locks
│   │   └── nats/         # pub/sub de eventos
│   ├── app/              # casos de uso
│   ├── config/           # carga de configuración
│   └── platform/         # cross-cutting: logger, telemetría
├── migrations/           # migraciones SQL versionadas
├── db/queries/           # queries para sqlc
├── docs/adr/             # Architecture Decision Records
└── deployments/          # docker-compose y affines
```

## Documentación

Las decisiones arquitectónicas importantes están documentadas en
[`docs/adr/`](./docs/adr/).

## Licencia

Por definir.
