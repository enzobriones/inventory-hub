# ADR 0001: Arquitectura hexagonal (ports & adapters)

- Estado: Aceptada
- Fecha: 2026-04-18

## Contexto

El sistema va a tener múltiples adapters externos:
- SII (facturación electrónica chilena)
- Mercado Libre, Shopify (canales de venta)
- WhatsApp Business (canal de ventas y notificaciones)
- Transbank (pagos)

Todas estas APIs son volátiles, tienen rate limits, fallan, cambian versiones,
y no podemos levantarlas en tests. Si el código de dominio conoce detalles de
esas APIs, los cambios externos se propagan a todo el sistema.

## Decisión

Adoptamos arquitectura hexagonal (también conocida como ports & adapters):

- `internal/domain` — entidades y reglas de negocio puras. Sin dependencias
  externas (solo stdlib).
- `internal/ports` — interfaces que el dominio necesita (ej: `ProductRepository`,
  `EventPublisher`). Son contratos que el dominio consume.
- `internal/adapters` — implementaciones concretas de esas interfaces contra
  tecnologías específicas (`postgres`, `redis`, `nats`, `http`).
- `internal/app` — casos de uso / application services que orquestan dominio
  y ports.

## Regla de oro

`domain` **nunca** importa desde `adapters` ni `platform`. El compilador lo
valida implícitamente cuando los paquetes están bien separados. Si aparece un
import prohibido, es bug de diseño, no de estilo.

## Consecuencias

### Positivas
- Tests del dominio sin Docker, sin red, sin fixtures complejos.
- Adapters intercambiables: mock en tests unitarios, real en integración/prod.
- Cambios en APIs externas quedan aislados en un solo archivo por adapter.
- Facilita la migración futura a múltiples servicios si hace falta.

### Negativas
- Más archivos y algo más de "boilerplate" que un approach monolítico clásico.
- Requiere disciplina: es tentador importar `postgres` en un caso de uso
  "solo por esta vez".
- La curva de aprendizaje es notable si no vienes de DDD/clean architecture.

## Alternativas consideradas

- **Layered architecture tradicional (controller → service → repository)**:
  descartada porque acopla service al repositorio concreto y complica tener
  múltiples adapters para el mismo concepto.
- **Microservicios desde el inicio**: descartada por overhead operacional sin
  beneficios reales en esta etapa. Un monolito modular bien separado se puede
  dividir después si hace falta.
