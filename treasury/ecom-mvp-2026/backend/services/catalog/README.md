# Catalog Service

Go microservice for the B2C e-commerce platform's product and category catalog.

## Build & Test

```bash
go build ./...   # must be clean
go test ./...    # table-driven handler tests: app/catalog package
```

`go test` covers the 3 FULL handlers with happy + negative paths:
- `handler_product_list_test.go` — CAT-001..005 validation + service stubs
- `handler_product_detail_test.go` — CAT-006 visibility + error paths
- `handler_product_create_test.go` — CAT-007 auth/validation/service error paths

Tests use an in-package `catalogService` interface fake; no DB or Kafka required.

## What's Full vs Stubbed

### FULL handlers (3)

| Endpoint | File | BA ACs |
|----------|------|--------|
| `POST /api/v1/catalog/product/list` | `app/catalog/handler_product_list.go` | CAT-001..005 |
| `POST /api/v1/catalog/product/detail` | `app/catalog/handler_product_detail.go` | CAT-006 |
| `POST /api/v1/catalog/product/create` | `app/catalog/handler_product_create.go` | CAT-007 |

### STUB handlers (8 — return 501 NOT_IMPLEMENTED_MVP)

| Endpoint | File | Notes |
|----------|------|-------|
| `POST /api/v1/catalog/product/update` | `app/catalog/handler_product_update.go` | CAT-008 |
| `POST /api/v1/catalog/product/soft-delete` | `app/catalog/handler_product_soft_delete.go` | CAT-009 |
| `POST /api/v1/catalog/category/create` | `app/catalog/handler_category_create.go` | CATE-001 |
| `POST /api/v1/catalog/category/update` | `app/catalog/handler_category_update.go` | CATE-002 |
| `POST /api/v1/catalog/category/deactivate` | `app/catalog/handler_category_deactivate.go` | CATE-003 |
| `POST /api/v1/catalog/category/list` | `app/catalog/handler_category_list.go` | CATE-003 |

All stub routes are **registered in the router** — they return HTTP 501 with body
`{ "code": "NOT_IMPLEMENTED_MVP", "message": "Endpoint stubbed for MVP — see README" }`.

## How to Run

### Prerequisites

- Go 1.25+
- Docker Compose v2

### Local with Docker Compose

```bash
# From services/catalog/
docker compose up postgres -d
# Wait for postgres to be healthy, then apply migrations manually:
psql -h localhost -U catalog -d catalog -f migrations/001_create_categories.up.sql
psql -h localhost -U catalog -d catalog -f migrations/002_create_products.up.sql
psql -h localhost -U catalog -d catalog -f migrations/003_create_product_images.up.sql
psql -h localhost -U catalog -d catalog -f migrations/004_create_product_status_history.up.sql
psql -h localhost -U catalog -d catalog -f migrations/005_create_outbox_events.up.sql

# Run the service
PORT=8080 \
REF_ID_HEADER_KEY=X-Ref-ID \
POSTGRES_HOST=localhost POSTGRES_USER=catalog SECRET_POSTGRES_PASSWORD=catalogpw POSTGRES_DB=catalog \
JWT_ISSUER=shoppilot JWT_AUDIENCE=shoppilot JWT_PUBLIC_KEY="<your-es256-public-pem>" \
KAFKA_ENABLED=false \
go run .
```

Or use Docker Compose for the full stack (set `JWT_PUBLIC_KEY` in `docker-compose.yml`):

```bash
docker compose up --build
```

### Environment Variables

| Variable | Required | Default | Notes |
|----------|----------|---------|-------|
| `PORT` | yes | — | HTTP listen port |
| `REF_ID_HEADER_KEY` | yes | — | Tracing header name |
| `POSTGRES_HOST` | yes | — | |
| `POSTGRES_USER` | yes | — | |
| `SECRET_POSTGRES_PASSWORD` | no | "" | |
| `POSTGRES_DB` | yes | — | |
| `JWT_ISSUER` | yes | — | |
| `JWT_AUDIENCE` | yes | — | |
| `JWT_PUBLIC_KEY` | yes | — | PEM-encoded ES256 public key |
| `KAFKA_ENABLED` | no | `false` | Set `true` + `KAFKA_PRODUCER_BROKERS` to enable real publishing |
| `KAFKA_PRODUCER_BROKERS` | no | — | CSV broker list |
| `KAFKA_CATALOG_TOPIC` | no | `ecom.catalog.events` | |
| `OUTBOX_RELAY_INTERVAL_SECONDS` | no | `5` | Relay tick interval |
| `MIGRATIONS_AUTO_APPLY` | no | `false` | Apply migrations on startup |
| `INVENTORY_BASE_URL` | no | `http://inventory-svc` | Base URL for inventory client |

## Architecture Notes

### Outbox relay

The relay goroutine starts at service boot (`main.go: startOutboxRelay`). On each tick it:

1. SELECTs up to 100 pending rows: `FOR UPDATE SKIP LOCKED ORDER BY id ASC`
2. Sends each to Kafka topic `ecom.catalog.events` with partition key = `sku`
3. Marks successful rows `published_at = NOW()` in a single batch UPDATE
4. Failed Kafka sends are **not** marked — they are retried on the next tick (at-least-once delivery)

When `KAFKA_ENABLED=false`, step 2 is skipped and all rows are marked processed immediately.
A warning is logged once at boot:

```
WARN KAFKA_ENABLED=false — outbox relay will mark events processed without publishing to Kafka
```

### inStock filter (CAT-005)

`product.list` with `filter.inStock=true` is **partially implemented**. The SQL query
returns ACTIVE+visible products meeting other filters. The inventory cross-service call
(`inventory.stock.bulk-read`) is stubbed with a `// TODO` comment in
`app/catalog/handler_product_list.go` and `app/catalog/service.go`. Wiring requires
the Inventory service to be live.

### Shipping fee tier (§11.2)

Catalog is read-only for pricing. The shipping fee tier calculation (§11.2) lives in
the Checkout service and references product price snapshots at order time. Catalog does
**not** compute or expose shipping fees. No interaction point exists in this service.

### stockStatus on product.detail

The `stockStatus` field in `product.detail` is stubbed to return `OUT_OF_STOCK` until the
inventory client is wired. See `app/catalog/handler_product_detail.go` for the TODO comment.

### Admin visibility

`product.detail` for DELETED/DRAFT products: admin detection uses the JWT claims in context
(set by the middleware group for protected routes). For the public detail route, the service
checks context claims — if no valid token is present, non-ACTIVE/non-INACTIVE products return 404.

## Module and replace directive

```
module github.com/example/shoppilot/catalog
replace gitlab.com/b2c-e-commerce-platform/platform/backend/common => ../../common
```

Build from the monorepo root or ensure `../../common` exists relative to `services/catalog/`.
