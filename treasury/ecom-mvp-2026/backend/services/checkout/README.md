# checkout service

Sync HTTP orchestrator that turns a customer's cart + chosen address into a server-priced preview or a transactionally-safe `Order + Reservation + PaymentIntent` triple.

## What is fully implemented

| File | Status |
|---|---|
| `app/checkout/pricing.go` | FULL — sole §11.2 shipping-fee tier source |
| `app/checkout/handler_preview.go` | FULL — all 4 downstream steps + server-side pricing |
| `app/checkout/handler_commit.go` | FULL — all 9 orchestration steps + full compensation matrix |
| `app/checkout/access/storage_idempotency.go` | FULL — Lookup (FOR UPDATE), InsertInflight, Finalize, CleanupExpired, MarkStaleInflightAbandoned |
| `app/checkout/access/storage_saga_log.go` | Sketch (write-only insert) |
| All `client_*.go` | FULL — per-call timeouts, retry logic, compensation backoff |

## Orchestration steps (`checkout.commit`)

```
Step 0  Idempotency-Key lookup / INSERT INFLIGHT (Postgres TX opened here)
Step 1  cart.read                          timeout 800ms, 1 retry
Step 2  identity.address.list              timeout 500ms, 1 retry
Step 3  catalog.product.detail (parallel)  timeout 500ms each, errgroup
Step 4  inventory.stock.bulk-read          timeout 800ms, 1 retry (fail-fast)
Step 5  orderId = uuid.NewV7() + server-side pricing snapshot
Step 6  inventory.reservation.create       timeout 1.5s, 1 retry, idempotent on orderId
Step 7  order.create-from-checkout         timeout 1s, 1 retry, idempotent on orderId
Step 8  payment.intent.create              timeout 1.5s, 1 retry, idempotent on orderId
Step 9  cart.clear-on-checkout             best-effort; 2 retries; failure is non-fatal
Step 10 Finalize idempotency_keys COMPLETED + saga_log INSERT + TX COMMIT
```

## Compensation matrix

| Step succeeded | Next step failed | Compensation |
|---|---|---|
| Step 5 (orderId generated) | Step 6 reservation.create | None — no remote state to undo |
| Step 6 reservation.create | Step 7 order.create-from-checkout | `inventory.reservation.release` — 3 retries, 50/200/800ms backoff |
| Step 7 order.create-from-checkout | Step 8 payment.intent.create | (A) `order.cancel-on-checkout-failure` 3 retries; then (B) `inventory.reservation.release` 3 retries |
| Step 8 payment.intent.create | Step 9 cart.clear-on-checkout | **None** — order is PENDING_PAYMENT; log warning + metric only |

On compensation failure (all retries exhausted): log `CRITICAL`, write `saga_log.status=DEGRADED`, rely on inventory reservation TTL (15 min) as the absolute safety net.

## Idempotency-Key semantics

- Header: `Idempotency-Key: <opaque string, max 128 chars, UUID v4 recommended>`.
- Scoped per `(endpoint, customerUserId)` — cross-tenant collision is impossible.
- TTL: 24 hours.
- States: `INFLIGHT` → `COMPLETED` | `ABANDONED`.
- Duplicate same-payload: return cached response, zero side-effects.
- Duplicate different-payload: `409 IDEMPOTENCY_KEY_REUSED`.
- Concurrent same-key: `409 IDEMPOTENCY_KEY_INFLIGHT` + `Retry-After: 1`.
- Stale INFLIGHT > 60s: janitor background goroutine marks `ABANDONED` (HTTP 504 envelope).

## Pricing (§11.2 BA-PR-002)

Single source of truth: `app/checkout/pricing.go`.

```
shippingFee = 60 THB  if subtotal < 1500 THB
shippingFee = 0 THB   if subtotal >= 1500 THB  (inclusive at 1500)
grandTotal  = subtotal + shippingFee - couponDiscount (always 0 in MVP)
```

## TD ambiguity handling

**CHK-AMBIG-001** (`order.create-from-checkout` missing contract): implemented as `POST /api/v1/order/order/create-from-checkout` and `POST /api/v1/order/order/cancel-on-checkout-failure` per the TD's chosen path. Internal auth via `X-Internal-Secret` header. TL must add these as first-class contracts. See `client_order.go`.

**CHK-AMBIG-002** (`couponCode` field): requests containing `couponCode` receive `400 VALIDATION_ERROR` with message "couponCode is not supported in MVP" — NOT silently stripped.

## Environment variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `SERVER_PORT` | no | `8080` | HTTP listener port |
| `POSTGRES_DSN` | **yes** | — | Postgres connection string |
| `JWT_PUBLIC_KEY` | **yes** | — | ES256 public key (PEM) for JWT verification |
| `JWT_ISSUER` | no | `shoppilot-identity` | Expected JWT issuer |
| `JWT_AUDIENCE` | no | `shoppilot-api` | Expected JWT audience |
| `INTERNAL_SHARED_SECRET` | **yes** | — | Shared secret for `X-Internal-Secret` header sent to order service |
| `CART_BASE_URL` | **yes** | — | Cart service base URL |
| `CATALOG_BASE_URL` | **yes** | — | Catalog service base URL |
| `INVENTORY_BASE_URL` | **yes** | — | Inventory service base URL |
| `PAYMENT_BASE_URL` | **yes** | — | Payment service base URL |
| `ORDER_BASE_URL` | **yes** | — | Order service base URL |
| `IDENTITY_BASE_URL` | **yes** | — | Identity service base URL |
| `CART_READ_TIMEOUT` | no | `800ms` | Per TD timeouts_and_retries_summary |
| `CATALOG_DETAIL_TIMEOUT` | no | `500ms` | |
| `INVENTORY_BULK_READ_TIMEOUT` | no | `800ms` | |
| `INVENTORY_RESERVATION_TIMEOUT` | no | `1500ms` | |
| `INVENTORY_RELEASE_TIMEOUT` | no | `1000ms` | |
| `PAYMENT_INTENT_TIMEOUT` | no | `1500ms` | |
| `ORDER_CREATE_TIMEOUT` | no | `1000ms` | |
| `ORDER_CANCEL_TIMEOUT` | no | `1000ms` | |
| `IDENTITY_ADDRESS_TIMEOUT` | no | `500ms` | |
| `HTTP_CLIENT_ENABLE_LOG_DEBUG` | no | `false` | Verbose HTTP request/response logging |

## Running locally

```bash
docker compose up
```

Runs on `http://localhost:8083`. PostgreSQL on `localhost:5433`.

Apply migrations manually or mount the `migrations/` directory as shown in `docker-compose.yml`.
