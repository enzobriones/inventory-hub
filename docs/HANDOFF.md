# Handoff — Inventory Hub

Documento de contexto completo para retomar el desarrollo en una nueva sesión.
Última actualización: 2026-04-18.

---

## 1. Visión del proyecto

**Qué es**: Hub de sincronización de inventario multicanal para PyMEs chilenas.
Una fuente única de verdad para el stock, que se sincroniza automáticamente hacia
múltiples canales de venta (POS propio, Mercado Libre, Shopify, WhatsApp Business)
y emite documentos tributarios electrónicos (DTE) vía integración con el SII
como capstone final.

**Enfoque declarado**: el objetivo primario es **aprender** (no vender).
Construido para subir de nivel técnico con Go y sistemas distribuidos, y para
servir como pieza central de portafolio profesional.

**Nivel del dueño**: primer proyecto serio en Go. Viene de otros lenguajes.
Chileno, habla español naturalmente. Prefiere estilo guiado paso a paso ("tú me
dices qué hacer, lo hago, avanzamos"), no volcados gigantes de código.

---

## 2. Roadmap completo

El proyecto está estructurado en 11 fases. Cada fase tiene entregables concretos
y desbloquea aprendizaje específico.

| Fase | Tema | Estado |
|------|------|--------|
| 0 | Fundación del proyecto: estructura, infra local, CI, logging, graceful shutdown | ✅ **Completada** |
| 1 | Modelo de dominio: entidades, invariantes, tests unitarios puros | 🟡 **En curso** (paso 11 de ~20 completado) |
| 2 | Persistencia e idempotencia: Postgres, sqlc, patrón inbox/outbox | Pendiente |
| 3 | API REST + auth: JWT, refresh tokens, multi-tenancy, rate limiting | Pendiente |
| 4 | Primer adapter: POS propio con ventas que descuentan stock atómicamente | Pendiente |
| 5 | Workers y eventos: NATS JetStream, reintentos, DLQ | Pendiente |
| 6 | Adapter externo real: Mercado Libre (OAuth2, webhooks, rate limits) | Pendiente |
| 7 | Conflictos y conciliación: reservas, oversell, workflows de resolución | Pendiente |
| 8 | Observabilidad seria: OpenTelemetry, Prometheus, Grafana | Pendiente |
| 9 | Segundo adapter: Shopify o WhatsApp Business | Pendiente |
| 10 | SII (capstone): firma XML, certificado digital, DTE, timbre PDF417 | Pendiente |

---

## 3. Decisiones de diseño confirmadas (Fase 1)

El modelo de dominio ya fue discutido y consensuado. Estas decisiones son
**vinculantes** — cualquier cambio requiere revisar consecuencias.

### Estructura física del inventario

| # | Decisión | Motivo |
|---|----------|--------|
| 1.1 | Multi-ubicación desde día 1 | Evita migración dolorosa cuando el cliente crece a 2+ locales |
| 1.2 | Stock compartido + reservas por canal (híbrido, configurable) | Cubre todos los casos de uso reales |

### Catálogo

| # | Decisión | Motivo |
|---|----------|--------|
| 2.1 | `Product` padre + `ProductVariant` hijas | Estándar de la industria (Shopify, ML) |
| 2.2 | SKU vive en la variante, no en el producto padre | La variante es la unidad real de venta |

### Reglas de stock

| # | Decisión | Motivo |
|---|----------|--------|
| 3.1 | Flag `allow_backorder` configurable **por producto** | Flexible para pre-ventas, servicios, productos bajo pedido |
| 3.2 | Event sourcing ligero: `StockMovement` como fuente de verdad, `Inventory` como proyección materializada | Auditoría completa + queries rápidas |

### Canales y ventas

| # | Decisión | Motivo |
|---|----------|--------|
| 4.1 | `Channel` como entidad en BD (con credenciales OAuth propias) | Una PyME puede tener 2+ cuentas del mismo tipo |
| 4.2 | `Sale` como entidad propia en el hub (no solo delegada al canal) | Habilita reportería consolidada cross-canal |
| 4.3 | Precio override por `(variante × canal)` mediante tabla `ChannelPrice` | Marketplaces suelen requerir precios distintos |

### Dinero e impuestos

| # | Decisión | Motivo |
|---|----------|--------|
| 5.1 | Montos como `amount int64 + currency char(3)` (patrón Stripe) | Multimoneda sin bugs de float |
| 5.2 | Precios con IVA incluido en BD (neto se calcula al emitir DTE) | Cómo piensa el usuario del sistema |
| 5.3 | Flag `tax_exempt` por producto | Productos exentos (libros, ciertos alimentos) |

### Flujo venta → stock

- `Sale` confirmada genera automáticamente los `StockMovement` correspondientes.
- `StockReservation` es **flexible** (abierta): apunta a `(variant, location)`,
  pero puede vincularse opcionalmente a `Sale` pendiente o usarse de forma
  standalone. Configurable por tenant: política estricta (descuenta disponible)
  vs laxa (solo informativa).

### Transversal

- **Multi-tenant desde la migración #1**: toda fila de toda tabla lleva
  `tenant_id`. No-negociable.
- **Arquitectura hexagonal**: `domain` nunca importa `adapters` ni `platform`.

---

## 4. Stack técnico

| Pieza | Tecnología | Versión | Rol |
|-------|------------|---------|-----|
| Lenguaje | Go | 1.22+ (directiva `go 1.22` en `go.mod`, local 1.26.2) | Backend |
| BD principal | PostgreSQL | 16 | Fuente de verdad |
| Cache / locks | Redis | 7 | Rate limiting, locks distribuidos |
| Broker de eventos | NATS | 2.10 + JetStream | Pub/sub persistente |
| Tracing | Jaeger | 1.57 (local) | OpenTelemetry backend |
| Migraciones | golang-migrate o goose | por decidir en Fase 2 | |
| Queries | sqlc | por decidir en Fase 2 | SQL tipado |
| Tests de integración | testcontainers-go | por decidir en Fase 2 | |
| Linter | golangci-lint | v2.11.4 | `revive`, `staticcheck`, `gosec`, etc. |
| CI | GitHub Actions | Node 20 runners | Build + test + lint |
| Orquestación local | Docker Compose | | Todo el stack con un comando |

**Decisiones de librerías ya tomadas**:
- Logging: `log/slog` (stdlib, Go 1.21+). No usar Zap/Logrus.
- Config: `os.Getenv` directo. No usar Viper/envconfig por ahora.
- Router HTTP: `net/http` con sintaxis `"GET /path"` (Go 1.22+). No usar Gin/Echo/Chi.
- UUIDs: `github.com/google/uuid`.

---

## 5. Estado del repositorio

**URL**: https://github.com/enzobriones/inventory-hub (público)
**Rama principal**: `main`
**CI**: verde ✓

### Estructura actual

```
.
├── .github/workflows/ci.yml           # build+test+lint en cada push/PR
├── .golangci.yml                      # config del linter (v2, locale ES)
├── .env.example                       # template de variables
├── .env                               # local, ignorado en git
├── .gitignore
├── Makefile                           # comandos cómodos (up, run-api, test, lint)
├── README.md                          # con badge de CI
├── go.mod                             # module github.com/enzobriones/inventory-hub, go 1.22
├── go.sum
├── cmd/
│   ├── api/main.go                    # binario del API HTTP
│   └── worker/main.go                 # binario del worker (stub por ahora)
├── deployments/
│   └── docker-compose.yml             # postgres + redis + nats + jaeger
├── docs/
│   ├── HANDOFF.md                     # ESTE archivo
│   └── adr/
│       └── 0001-hexagonal-architecture.md
├── internal/
│   ├── adapters/http/
│   │   ├── health.go                  # GET /health
│   │   ├── middleware.go              # TraceID + Logging
│   │   └── response.go                # WriteJSON / WriteError
│   ├── config/config.go               # Load() con env vars
│   ├── domain/
│   │   ├── errors.go                  # ErrInvalidInput, ErrNotFound, etc.
│   │   ├── ids.go                     # tipos tipados + generadores + parsers
│   │   ├── money.go                   # value object Money
│   │   └── money_test.go              # table-driven tests
│   └── platform/logger/logger.go      # slog + trace ID en context
├── migrations/                        # vacío, para Fase 2
├── db/queries/                        # vacío, para Fase 2 (sqlc)
└── scripts/                           # vacío
```

### Lo que ya funciona end-to-end

- `make up` levanta postgres + redis + nats + jaeger sanos.
- `make run-api` arranca el servidor en `:8080` con graceful shutdown.
- `curl http://localhost:8080/health` responde JSON con `X-Trace-Id` en header.
- El trace ID del header coincide con el del log estructurado (correlación OK).
- `Ctrl+C` cierra limpio con log `server stopped cleanly`.
- `make run-worker` emite heartbeat cada 5s y también cierra limpio.
- `make test` corre tests con race detector, todos verdes.
- `make lint` pasa sin issues.
- CI en GitHub Actions corre build+test+lint en cada push.

### Comandos útiles de referencia

```bash
make help           # lista todo lo disponible
make up             # levanta infra local
make down           # baja infra local
make logs           # tail de logs de infra
make run-api        # API en :8080
make run-worker     # worker con heartbeats
make build          # compila bin/api y bin/worker
make test           # tests con -race
make lint           # golangci-lint v2.11.4
make fmt            # gofmt + go mod tidy
gh run watch        # ver CI en vivo
gh run view --web   # abrir CI en navegador
```

### Variables de entorno (`.env`)

```
HTTP_PORT=8080
LOG_LEVEL=debug
ENVIRONMENT=development
DATABASE_URL=postgres://inventory:inventory_dev@localhost:5432/inventory?sslmode=disable
REDIS_URL=redis://localhost:6379/0
NATS_URL=nats://localhost:4222
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318
OTEL_SERVICE_NAME=inventory-api
```

---

## 6. Convenciones del proyecto

Estas convenciones se establecieron y aplicaron en las fases 0 y 1. **Respetar
en todo código nuevo**.

### Código Go

- **Package comments obligatorios** en un archivo por paquete: `// Package X ...`.
- **Todo símbolo exportado** necesita comentario que empiece con su nombre
  (convención de godoc). `revive` lo verifica.
- **Errores empiezan con minúscula**: `fmt.Errorf("invalid HTTP_PORT: %w", err)`,
  no `"Invalid..."`. Convención Go, `staticcheck ST1005` la verifica.
- **Envolver errores con `%w`** siempre, para habilitar `errors.Is`/`errors.As`.
- **Acrónimos en mayúsculas**: `ID`, `URL`, `HTTP`, `API`, `JSON`. Nunca
  `Id`, `Url`, `Http`. `revive var-naming` lo verifica.
- **Parámetros no usados**: `_` en vez de un nombre. `revive unused-parameter`
  lo verifica.
- **Idioma**: nombres de símbolos en inglés, comentarios/documentación en
  español. `misspell` está configurado con lista de palabras en español ignoradas.

### Arquitectura

- **`domain` es puro**: no importa `adapters`, `platform`, ni nada fuera de
  stdlib (excepto `github.com/google/uuid` por practicidad). Si un test de
  dominio necesita red, Docker o fixtures complejos, es bug de diseño.
- **Value objects inmutables**: constructor valida, campos privados, todas
  las operaciones devuelven nueva instancia (ver `Money`).
- **IDs tipados con type safety**: el compilador impide pasar `ProductID`
  donde se espera `TenantID`. Ver `internal/domain/ids.go`.
- **Sentinel errors en `domain/errors.go`**: los errores específicos se
  crean envolviendo los base con `fmt.Errorf("...: %w", ErrInvalidInput)`.

### Tests

- **Paquete `_test`**: `package domain_test` (no `package domain`). Fuerza
  testear solo la API pública.
- **Table-driven tests** con `t.Run(name, ...)` para subtests.
- **`errors.Is`** para verificar tipo de error, nunca comparar mensajes.
- **Convención de mensajes**: "esperaba X, obtuve Y".

### Git / commits

- Mensajes en inglés con prefijo semántico: `feat`, `fix`, `chore`, `docs`,
  `ci`, `style`, `refactor`, `test`.
- Scope opcional entre paréntesis: `feat(domain): add typed ids ...`.
- CI debe quedar verde antes de seguir al próximo paso.

---

## 7. Cómo se trabaja en este proyecto (estilo de sesión)

El dueño del proyecto pidió explícitamente este estilo:

1. **Claude guía paso a paso**, un entregable a la vez.
2. El dueño confirma "listo" antes de avanzar al siguiente paso.
3. Si algo falla, pega el error completo y se resuelve antes de seguir.
4. Las decisiones importantes de diseño se discuten con opciones + trade-offs
   + recomendación razonada. El dueño elige.
5. Preferir **explicar el porqué** antes que tirar código. El código sin
   contexto se olvida; las razones se retienen.
6. Priorizar aprendizaje sobre velocidad. Si una pieza enseña algo valioso,
   vale hacerla bien aunque tarde más.
7. Escribir los commits y documentación **en inglés**, pero la comunicación
   y los comentarios en español.

**No hacer**:
- Volcados gigantes de código sin contexto.
- Asumir conocimiento previo de DDD, hexagonal, event sourcing, etc. — explicar
  cuando aparecen por primera vez.
- Avanzar si el dueño no confirmó el paso anterior.

---

## 8. Qué pasó en cada fase (resumen cronológico)

### Fase 0 — Fundación

- Inicializado módulo Go (`github.com/enzobriones/inventory-hub`, directiva
  `go 1.22` para compatibilidad con herramientas).
- Creada estructura hexagonal completa de carpetas.
- Docker Compose con postgres, redis, nats (JetStream + monitor en `:8222`)
  y jaeger.
- `internal/config/config.go` con `Load()` que lee env vars, diferencia entre
  required (DATABASE_URL, etc.) y optional (LOG_LEVEL default "info").
- `internal/platform/logger/logger.go` con `slog`, handler de texto en dev y
  JSON en prod, helpers `WithTraceID`/`TraceIDFrom` usando `context.Context`
  con key de tipo privado (patrón idiomático Go).
- `internal/adapters/http/`: `HealthHandler` devuelve status/timestamp/version,
  `TraceIDMiddleware` genera UUID si el cliente no manda uno, `LoggingMiddleware`
  registra método/path/status/duración con `statusRecorder` (patrón para
  capturar status después del `WriteHeader`).
- `cmd/api/main.go` con graceful shutdown vía `signal.NotifyContext`, timeouts
  defensivos (`ReadHeaderTimeout` para mitigar Slowloris), ListenAndServe
  en goroutine para no bloquear el `select`.
- `cmd/worker/main.go` como stub con ticker cada 5s, misma estructura de
  graceful shutdown.
- Makefile con `include .env` condicional (`ifneq .wildcard .env`), comandos
  coloreados en `make help`.
- CI con dos jobs paralelos (`build-and-test`, `lint`), verifica `go mod tidy`
  limpio, tests con `-race`.
- `.golangci.yml` v2 (formato nuevo), `misspell` con lista de palabras en
  español ignoradas.
- ADR `0001-hexagonal-architecture.md` documenta decisión base.
- README con badge de CI, setup en 3 comandos.

### Fase 1 — Modelo de dominio (en curso)

Completado:
- Discusión exhaustiva de decisiones de diseño (11 preguntas agrupadas en
  5 bloques, todas resueltas).
- Diagrama visual del modelo completo mostrado y validado con el dueño.
- `internal/domain/errors.go`: 4 errores sentinela
  (`ErrInvalidInput`, `ErrNotFound`, `ErrConflict`, `ErrUnauthorized`).
- `internal/domain/ids.go`: 10 tipos de ID tipados + generadores (`NewXxxID`)
  + 6 parsers (`ParseXxxID` para IDs que vienen de HTTP).
- `internal/domain/money.go`: value object `Money` con `amount int64` +
  `currency string`, constructor validante (`NewMoney`), constructor para
  tests (`MustMoney`), accessors de solo lectura, operaciones
  `Add`/`Sub`/`Mul`/`Equal`, predicados `IsZero`/`IsNegative`/`IsPositive`,
  `String()` para logs. Monedas soportadas: CLP, USD, EUR.
- `internal/domain/money_test.go`: tests table-driven exhaustivos que
  validan la API pública (construcción, operaciones, predicados, errores
  por monedas distintas).

Pendiente en Fase 1:
- **Paso 12**: modelar `Product` y `ProductVariant` con invariantes:
  handle slug válido, al menos 1 variante, SKU único por tenant,
  flag `allow_backorder`, flag `tax_exempt`.
- **Paso 13**: modelar `Location`.
- **Paso 14**: modelar `StockMovement` (inmutable, tipos: sale, purchase,
  adjustment, transfer_in, transfer_out, return).
- **Paso 15**: modelar `Inventory` como proyección con `on_hand`, `reserved`,
  `available = on_hand - reserved`.
- **Paso 16**: modelar `StockReservation` flexible (política configurable).
- **Paso 17**: modelar `Channel` con credenciales opacas.
- **Paso 18**: modelar `Sale` y `SaleItem`.
- **Paso 19**: modelar `ChannelPrice` (override precio por variante × canal).
- **Paso 20**: tests de invariantes transversales + primer caso de uso
  mínimo (probablemente "registrar venta POS") en `internal/app/` para
  validar que el modelo funciona compuesto.

---

## 9. Issues técnicos ya resueltos (para no repetirlos)

1. **NATS monitoring port**: el flag `-m 8222` es obligatorio; solo exponer
   el puerto en Docker no activa el endpoint HTTP.
2. **Go version en `go.mod`**: `go 1.26.2` rompe herramientas del ecosistema
   compiladas con versiones previas. Bajado a `go 1.22` (la mínima que
   necesitamos por el router con verbos).
3. **golangci-lint v1 no soporta Go 1.26**: migrado a v2.11.4.
4. **golangci-lint v2 requiere action v7**: `golangci/golangci-lint-action@v7`,
   no `@v6`.
5. **Config de golangci-lint v2**: formato nuevo, `version: "2"`,
   `linters.default: none` (no `disable-all: true`), formatters separados.
6. **misspell en español**: configurar `locale: US` pero agregar las
   palabras en español a `ignore-rules`.
7. **Package comments**: `revive` v2 los exige en un archivo por paquete.
   Poner en uno solo — duplicar da warning.

---

## 10. Cómo arrancar la próxima sesión

1. El dueño abre el proyecto, corre `make up` para levantar infra.
2. Verifica que CI esté verde (`gh run list`).
3. Le comparte este documento a Claude (o lo pega al inicio de la conversación).
4. Dice: "Sigamos con el paso 12" (o el paso que corresponda).
5. Claude responde con el paso concreto, el dueño lo ejecuta, confirma, siguen.

### Punto exacto donde retomar

**Paso 12: Modelar `Product` y `ProductVariant`**.

Crear `internal/domain/product.go` con:
- Value object `Handle` o `Slug` (ej: `polera-manga-corta-2026`), validado
  (solo `[a-z0-9-]`, sin espacios, longitud razonable).
- Value object `SKU` (ej: `POL-2026-ROJ-M`), validado.
- Entidad `Product`:
  - Campos: ID, TenantID, Handle, Name, Description, AllowBackorder,
    TaxExempt, CreatedAt, UpdatedAt, Variants (slice).
  - Invariante: debe tener al menos 1 variante.
  - Invariante: handle único por tenant (validado en repositorio después,
    pero el dominio tiene que exponer el método para detectar).
  - Constructor `NewProduct(...)` que valide todo.
  - Métodos: `AddVariant`, `RemoveVariant` (no si es la última), getters.
- Entidad `ProductVariant`:
  - Campos: ID, ProductID, SKU, Attributes (map string→string para
    talla/color/etc.), Price (Money), CreatedAt, UpdatedAt.
  - Constructor `NewProductVariant(...)` que valide.
  - Invariante: precio positivo o cero, nunca negativo.
- Tests unitarios table-driven de ambos.

Antes de escribir código, plantear al dueño:
- ¿Handle se genera automático desde el name o se pide explícito?
- ¿Atributos de variante son map libre o estructura fija (talla+color)?
- ¿Precio base vive en la variante (como quedó en diagrama) o en el producto?

---

## 11. Referencias internas

- **ADR 0001**: `docs/adr/0001-hexagonal-architecture.md` (arquitectura base).
- **Diagrama del dominio**: fue mostrado visualmente, no está guardado
  como archivo. Si se necesita recrear, está en la conversación de
  la sesión 1 y en la sección 3 de este handoff.
- **Repositorio**: https://github.com/enzobriones/inventory-hub
