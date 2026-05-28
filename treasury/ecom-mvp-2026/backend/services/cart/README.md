# Cart Service

Manages the B2C shopping cart: add, read, update, and remove items, plus checkout clearing.

## Architecture

```
POST /api/v1/cart/cart/add-item          → handler_add_item.go
POST /api/v1/cart/cart/read              → handler_read.go
POST /api/v1/cart/cart/update-item       → handler_update_item.go
POST /api/v1/cart/cart/remove-item       → handler_remove_item.go       [STUB 501]
POST /api/v1/cart/cart/clear-on-checkout → handler_clear_on_checkout.go [STUB 501]
```

All endpoints require a `Authorization: Bearer <JWT ES256>` customer token. The JWT middleware (from `common/middleware`) validates signature, issuer, audience, and expiry, then injects `token.Claims` into context.

## Fan-out: cart.read

`cart.read` enriches each cart line with live data from two upstream services:

1. **Catalog service** (`CATALOG_BASE_URL/api/v1/catalog/product/detail`) — one call per item, in parallel.
2. **Inventory service** (`INVENTORY_BASE_URL/api/v1/inventory/stock/bulk-read`) — one bulk call for all product IDs, in parallel with catalog calls.

The fan-out is implemented with `golang.org/x/sync/errgroup`:
- A per-request timeout context (`FANOUT_TIMEOUT_MS`, default 3 000 ms) wraps all goroutines.
- A semaphore channel (size `FANOUT_CONCURRENCY`, default 10) caps concurrent catalog calls.
- Items are capped at `FANOUT_MAX_ITEMS` (default 50) before fan-out begins.

### Degradation behavior

| Failure | Effect |
|---------|--------|
| Catalog call times out for a single item | That line: `checkoutable=false`, `productStatus="CATALOG_UNREACHABLE"`. Cart read continues — no error. |
| Inventory bulk-read fails | All lines get `availableQty=0`, `checkoutable=false`. Cart read continues. |
| Catalog call fails hard (non-timeout) | Returns `INTERNAL_ERROR` to client. |

Subtotal is computed server-side: `sum(currentPrice * qty)` for checkoutable lines only. Prices are never stored in cart — always fetched live from catalog.

## Key design decisions

| Decision | Rationale |
|----------|-----------|
| No stock check on add-item | INV-002: reservation happens at checkout, not add-to-cart. Catalog status check only (CART-001). |
| ON CONFLICT merge on add-item | CART-006: same SKU adds merge into one line (`qty = existing + new`). Concurrent adds are safe — UNIQUE constraint + row lock. |
| qty=0 on update-item triggers remove | AMB-004: API accepts `qty >= 0`; `qty=0` delegates to delete path. `qty < 0` is rejected as VALIDATION_ERROR. |
| cart.clear-on-checkout idempotency | CHK-010: keyed on `orderId` in `cart_clear_idempotency` table. Second call returns cached `removed` count without re-executing DELETE. |
| No price snapshot in cart_items | Prices fetched live at read-time. Authoritative snapshot is taken at checkout.commit (stored on order_items). |

## Tunable environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| `FANOUT_TIMEOUT_MS` | `3000` | Per-call timeout for catalog + inventory fan-out (ms). |
| `FANOUT_MAX_ITEMS` | `50` | Maximum cart items processed in fan-out. Items beyond this cap are not returned. |
| `FANOUT_CONCURRENCY` | `10` | Max concurrent catalog HTTP calls per cart.read. |
| `CATALOG_BASE_URL` | — | Base URL of the catalog service (e.g. `http://catalog:8082`). |
| `INVENTORY_BASE_URL` | — | Base URL of the inventory service (e.g. `http://inventory:8083`). |
| `DB_HOST` / `DB_PORT` / `DB_NAME` | — | PostgreSQL connection. |
| `SECRET_DB_USERNAME` / `SECRET_DB_PASSWORD` | — | PostgreSQL credentials. |
| `JWT_ISSUER` / `JWT_AUDIENCE` / `SECRET_JWT_PUBLIC_KEY` | — | ES256 JWT verification. |

## What is stubbed (returns 501)

- `cart.remove-item` — `handler_remove_item.go`. Full logic: validate UUID, DELETE cart_items idempotently, touch updated_at, return enriched cart.
- `cart.clear-on-checkout` — `handler_clear_on_checkout.go`. Full logic: validate `customerUserId == claims.sub`, idempotency check on `orderId`, `DELETE WHERE id=ANY($2::uuid[])`, persist idempotency record, return `{removed}`.

Both have complete service-layer implementations in `service.go` (`RemoveItem`, `ClearOnCheckout`), storage implementations in `access/`, and the DB schema in `migrations/003`. Wiring the handlers is a one-commit task.

## Database schema

Schema: `cart` (PostgreSQL 15+).

| Table | Purpose |
|-------|---------|
| `cart.carts` | One row per customer. `UNIQUE(user_id)`. |
| `cart.cart_items` | Line items. `UNIQUE(cart_id, product_id)` enforces SKU-merge invariant (CART-006). |
| `cart.cart_clear_idempotency` | Replay guard for `cart.clear-on-checkout` keyed on `order_id`. |

## Running locally

```bash
docker-compose up
```

Runs the cart service on port 8084 with a local Postgres on 5432. Migrations are applied automatically via `docker-entrypoint-initdb.d`.
