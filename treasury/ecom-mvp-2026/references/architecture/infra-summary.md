# Infra Summary — ShopPilot MVP (B2C E-Commerce Platform)

Generated: 2026-05-08
Tech-Lead: Claude Opus 4.7 (1M context) — dry-run #2 (incremental refactor)

See also: [topology](./infra-topology.md) · [connectivity](./connectivity.md) · [observability](./observability-spec.md) · [contracts](./contracts.json) · [components](./components.json) · [ADRs](./ADRs/)

This run is a refactor of dry-run #1. The eight components (7 Go services + 1 Next.js frontend) are unchanged; the load-bearing edits all sit at the contract layer to close three high-severity REV-L2 findings (REV-L2-001 buyerEmail dataflow, REV-L2-002 INTERNAL_SHARED_SECRET env-var name, REV-L2-003 constant-time secret compare).

---

## Services

| Name | Runtime | Port | Owns | Complexity |
|---|---|---:|---|---|
| identity | Go 1.23 + Gin | 8081 | EPIC_AUTH (register, login, refresh, logout, profile, addresses) | standard |
| catalog | Go 1.23 + Gin | 8082 | EPIC_CATALOG (products, categories, soft-delete) | standard |
| inventory | Go 1.23 + Gin | 8083 | EPIC_INVENTORY (stock_levels, reservations, sweeper, admin adjust) | complex |
| cart | Go 1.23 + Gin | 8084 | EPIC_CART (cart_items, server-computed subtotal, status flagging) | standard |
| checkout | Go 1.23 + Gin | 8085 | EPIC_CHECKOUT (preview + commit orchestration; §11.2 pricing tier) | complex |
| order | Go 1.23 + Gin | 8086 | EPIC_ORDER (state machine, sole writer of order.status, cross-topic event consumers) | complex |
| payment | Go 1.23 + Gin | 8087 | EPIC_PAYMENT (intent + mock simulate + dedup'd callback + outbox) | complex |
| frontend-web | Next.js 14 + TypeScript + Tailwind | 3000 | EPIC_FRONTEND (Thai mobile customer webview) | complex |

> Reference scaffold for all 7 backends: `B2C E-Commerce Platform/backend/go-template/`. Shared libraries: `B2C E-Commerce Platform/backend/common/` (now exposes `wrapper.Response[T].TraceID` + `wrapper.CtxTraceID`; the new `common/middleware/internal_auth.go` is a follow-up landing in dry-run #2 per ADR-008).

## Databases

| Name | Engine | Schema owner | Tables (high level) | Notes |
|---|---|---|---|---|
| identity_db | PostgreSQL 16 | identity service | users, addresses, refresh_tokens (revoked-jti) | bcrypt password hash; immutable email |
| catalog_db | PostgreSQL 16 | catalog service | products, categories, product_images, product_status_history, outbox_events | NUMERIC price; soft-delete via status enum |
| inventory_db | PostgreSQL 16 | inventory service | stock_levels (CHECK >=0 invariant + version), reservations, stock_adjustments, outbox_events, consumed_events, idempotency_keys | FOR UPDATE on stock_levels per reservation |
| cart_db | PostgreSQL 16 | cart service | carts, cart_items (UNIQUE on (cart_id, sku) for CART-006 merge) | informational subtotal only |
| checkout_db | PostgreSQL 16 | checkout service | idempotency_keys (per-(endpoint, customerUserId) for the customer-facing key), saga_log (DEGRADED rows surface to the reconciler) | no Order or Payment rows; pure orchestrator |
| order_db | PostgreSQL 16 | order service | orders, order_items (immutable snapshot incl. buyerEmailSnapshot), order_status_history (append-only), order_addresses (snapshot), outbox_events, consumed_events | sole writer of order.status |
| payment_db | PostgreSQL 16 | payment service | payment_intents, payment_records, payment_callback_dedup (UNIQUE on (paymentIntentId, providerStatus)), outbox_events | dedup-violation = replay |

> Per `cross-cutting.persistence`: PostgreSQL only via `common/database` (pgx). Schema-per-service. NO cross-service joins. NO shared tables.

## Caches

| Name | Engine | Used by | Eviction |
|---|---|---|---|
| (none for MVP) | — | — | — |

> Cache deferred. The only cache-shaped surface is the identity refresh-jti revocation store; in MVP this is a Postgres table inside identity_db (not Redis) to remove a moving part.

## External / MAP APIs

| Name | Vendor | Used by | Endpoint | Auth env var |
|---|---|---|---|---|
| (none — mock-only) | — | — | — | — |

> Real payment, real shipping, social login, and audit/notification log are all out_of_scope per BA. Payment uses an in-process mock; shipping is admin-only and stores tracking number directly on `orders.tracking_number`.

## Message Bus

| Property | Value |
|---|---|
| Broker | Apache Kafka 3.x (Redpanda 23.x acceptable for local Docker Compose) |
| Topics | `ecom.payment.events` (payment.completed/failed/expired), `ecom.order.events` (order.cancelled), `ecom.inventory.events` (reservation.expired), `ecom.catalog.events` (product.created) |
| Partition strategy | per-orderId for payment + order + inventory.reservation events; per-sku for catalog product events |
| Retention | 7 days production; 1 day local |
| Consumer dedup | per-service `consumed_events(eventId PK, consumerName, processedAt)` table; ON CONFLICT DO NOTHING; combined with state-driven SELECT FOR UPDATE on the affected aggregate |

> Cross-topic ordering is NOT preserved by Kafka. Consumers (Order, Inventory) are state-driven — they read current row state inside the consumer transaction and act on actual current state, never on the event's `fromStatus` field. See ADR-009 (state-driven event consumers).

## Observability

| Layer | Component | Notes |
|---|---|---|
| Logs | slog → stdout → Loki via Promtail | structured JSON; required field set in `observability-spec.md`; redaction list enforces no secrets, no JWT bytes, no password material |
| Metrics | Prometheus scraping `/metrics` per service | standard HTTP histogram + per-domain counters (inventory_reservations_total, payment_callback_dedup_hits_total, checkout_orchestration_step_duration_seconds, order_state_transitions_total) |
| Traces | OpenTelemetry → Tempo via OTLP gRPC | W3C `traceparent` propagation; auto-injected into the response envelope's `traceId` field by `wrapper.CtxTraceID` |
| Dashboards | Grafana | "Service Health", "Customer Journey funnel", "Saga Compensation (DEGRADED rows + stuck PENDING_PAYMENT > 15min)", "Auth Anomaly", "Schema Performance" |
| Alerts | Prometheus rules | `HighErrorRate > 5%/5min`, `P95LatencyBreach > 1s/10min`, `SagaInFlightStuck > 15min`, `OutboxBacklog > 1000` |

> Local stack docker-compose snippet: `B2C E-Commerce Platform/observability/docker-compose.yml` — Loki, Tempo, Prometheus, Grafana auto-provisioned.

## Deployment topology

- **Local:** Docker Compose. 7 backend services + 1 frontend + 7 Postgres instances (or 1 instance with 7 schemas) + Kafka (or Redpanda) + observability stack. `docker-compose up`.
- **SIT/UAT/PRD:** Kubernetes; one Deployment per service; Postgres as a managed instance (or one StatefulSet per environment); Kafka as a StatefulSet or managed Confluent / MSK.
- **Auth-token expiry posture:** customer access JWT 15 min, refresh JWT 14 days (suggested defaults; both env-tunable on identity).
- **Internal-secret rotation:** `INTERNAL_SHARED_SECRET=new,old` is honored by all four recipient services (checkout, order, inventory, payment). Rolling restart with the comma-separated value gives a 5-minute overlap window. ADR-008 documents the runbook.
- **Frontend topology constraint:** The Next.js frontend MUST NOT load `INTERNAL_SHARED_SECRET` and MUST NOT forward `X-Internal-Secret`. Customer-originated requests reach backend services as Bearer JWT only. `frontend/.env.example` explicitly omits the var; a regression test asserts the absence of `X-Internal-Secret` on customer-originated requests (REV-L2-013 verified-correct posture).
