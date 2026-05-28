# Topology — ShopPilot MVP (B2C E-Commerce Platform)

## Layers

- **Edge / UI** — Customer mobile browser (390x844, Thai-language webview) → Next.js (Server Components + route handlers acting as the gateway; HttpOnly cookies for `access` + `refresh`; 6-step `withAuth` per `cross-cutting.auth.token_acquisition_flow`).
- **Processor** — 7 Go 1.23 + Gin microservices behind no API gateway in MVP (the Next.js route handlers play that role). Internal calls between services use `X-Internal-Secret` (env: `INTERNAL_SHARED_SECRET`; constant-time compare).
- **Data** — PostgreSQL per service (schema-per-service via `common/database`); Apache Kafka 3.x (or Redpanda for local) bus with 4 topics — `ecom.payment.events`, `ecom.order.events`, `ecom.inventory.events`, `ecom.catalog.events`.
- **Observability** — slog → Loki via Promtail; Prometheus (15s scrape) → Tempo via OTLP gRPC; Grafana dashboards. `wrapper.CtxTraceID` injects W3C trace id into the response envelope.

## Diagram

```
                         +------------------------------+
                         |   Customer mobile browser    |
                         |   (Thai, 390x844 viewport)   |
                         +--------------+---------------+
                                        | HTTPS
                                        v
                         +------------------------------+
                         |  Next.js  (frontend-web)     |
                         |  - SSR + RSC pages           |
                         |  - /api/proxy/[...path]      |
                         |  - withAuth (6-step flow)    |
                         |  - Idempotency-Key gen       |
                         |  HttpOnly: access, refresh   |
                         |  NEVER forwards X-Internal-* |
                         +--------------+---------------+
                                        | POST /api/v1/<svc>/<agg>/<action>
                                        | Authorization: Bearer <jwt>
        +------------+-----------+------+-----+----------+-----------+----------+
        v            v           v            v          v           v          v
   +--------+   +--------+   +--------+   +--------+ +--------+  +--------+ +--------+
   |identity|   |catalog |   |inventory|  |  cart  | |checkout|  | order  | |payment |
   | :8081  |   | :8082  |   | :8083   |  | :8084  | | :8085  |  | :8086  | | :8087  |
   |        |   |        |   |        |   |        | |  ORCH  |  | SOLE   | | DEDUP  |
   |  JWT   |   | public |   |  Bearer|   | Bearer | | +X-Int |  | WRITER | | +OUTBOX|
   | issuer |   |+admin  |   | + X-Int|   |        | | -Secret|  |status  | |        |
   +---+----+   +---+----+   +----+---+   +---+----+ +---+----+  +---+----+ +---+----+
       |            |             |           |          |           |          |
       v            v             v           v          v           v          v
   +--------+   +--------+   +--------+   +--------+ +--------+  +--------+ +--------+
   |identity|   |catalog |   |inventory|  | cart   | |checkout|  | order  | |payment |
   |  _db   |   |  _db   |   |  _db    |  |  _db   | |  _db   |  |  _db   | |  _db   |
   | (PG16) |   | (PG16) |   | (PG16)  |  | (PG16) | | (PG16) |  | (PG16) | | (PG16) |
   +--------+   +--------+   +--------+   +--------+ +--------+  +--------+ +--------+
                                              ^                                  ^
                       checkout->cart---------+                                  |
                       checkout->catalog (per-item, errgroup)                    |
                       checkout->identity.profile.read (Bearer; buyerEmail)      |
                       checkout->inventory.reservation.create (X-Int-Secret)     |
                       checkout->order.create-from-checkout (X-Int-Secret)-------|
                       checkout->payment.intent.create (X-Int-Secret)------------+
                       checkout->cart.clear-on-checkout (best-effort)

       ============== Kafka topic: ecom.payment.events (key=orderId) ==============
                                       ||
              +========================++=========================+
              ||  payment ==> payment.completed/failed/expired   ||
              vv                                                 vv
        +----------+                                       +-----------+
        |  order   |                                       | inventory |
        | consumer |                                       |  consumer |
        | (state-  |                                       |  (state-  |
        |  driven) |                                       |   driven) |
        +----------+                                       +-----------+

       =============== Kafka topic: ecom.order.events (key=orderId) ===============
                                       ||
                       order ==> order.cancelled (cancelActor=CUSTOMER|ADMIN|SYSTEM)
                                       vv
                                 +-----------+        +-----------+
                                 | inventory |        |  payment  |
                                 |  consumer |        |  consumer |
                                 | (state-   |        |  (state-  |
                                 |  driven)  |        |   driven) |
                                 +-----------+        +-----------+

       ============ Kafka topic: ecom.inventory.events (key=orderId) =============
                inventory ==> reservation.expired (sweeper-driven)
                                       vv
                                 +-----------+
                                 |   order   |
                                 |  consumer |
                                 | -> PAYMENT_EXPIRED
                                 +-----------+

       ============ Kafka topic: ecom.catalog.events (key=sku) ===================
                catalog ==> product.created
                                       vv
                                 +-----------+
                                 | inventory |
                                 | seeds stock_levels(sku) all-zero

       ----- system-internal sweepers (in-process goroutines, no HTTP route) -----
         ::> inventory.reservation.sweep-expired (every 30s; FOR UPDATE SKIP LOCKED; batch 200)
         ::> payment.intent.expiry-sweeper (every 60s; FOR UPDATE SKIP LOCKED)
         ::> outbox-relay (every 5s per service that owns an outbox: catalog, inventory, order, payment)
         ::> checkout saga reconciler (every 60s; surfaces saga_log WHERE status='DEGRADED' AND age > 1h)
```

## Legend

- `-->` synchronous HTTP request (same row indicates same source/target tier)
- `==>` async Kafka event (per-orderId or per-sku partition key)
- `::>` system-internal cron / sweeper (no HTTP route, no caller, no JWT)

## Notes on the diagram

- The frontend MUST NOT forward `X-Internal-Secret` and MUST NOT load `INTERNAL_SHARED_SECRET` (REV-L2-013 verified-correct posture). Every customer-originated arrow at the top of the diagram carries `Authorization: Bearer <jwt>` only.
- All four `+X-Int-Secret` arrows from `checkout` are constant-time compared on the recipient via `subtle.ConstantTimeCompare` (REV-L2-003 closure; ADR-008).
- `checkout -> identity.profile.read` is the buyerEmail dataflow path that closes REV-L2-001. It uses Bearer JWT (the same access token the customer's request arrived with), NOT internal secret — Identity is on the customer auth tier.
- `order` is the SOLE writer of `orders.status` (verified by grep across 7 backends in dry-run #1 REV-L2). Payment never writes Order rows; Inventory never writes Order rows; only Order's storage layer.
- All four Kafka consumers (Order x3, Inventory x4) are state-driven — they read current row state via FOR UPDATE inside the consumer transaction and act on actual current state, never on the event's `fromStatus` field. ADR-009 documents this.
