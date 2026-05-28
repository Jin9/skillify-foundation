# inventory service

Manages SKU-level stock, reservations, and their lifecycle.

## Implementation status

| File / handler | Status |
|---|---|
| `handler_stock_read.go` | **FULL** |
| `handler_reservation_create.go` | **FULL** (multi-SKU lex-order lock, FOR UPDATE, no-negative, idempotency) |
| `handler_reservation_create.go` — `ReservationReleaseHandler` | **FULL** (RESERVED→release, COMMITTED→sold-restock PR-006) |
| `sweeper.go` | **FULL** (30s ticker, per-row tx, outbox emit for `events.reservation.expired`) |
| `consumer.go` — `handlePaymentCompleted` | **FULL** (state-driven, dedup, PR-003-safe) |
| `outbox_relay.go` | **FULL** (5s poll, SKIP LOCKED, kafka-disabled mode) |
| `handler_stock_bulk_read.go` | **STUB** — 501 NOT_IMPLEMENTED_MVP |
| `handler_stock_adjust.go` | **STUB** — 501 NOT_IMPLEMENTED_MVP |
| `consumer.go` — `handlePaymentFailed` | **STUB** — TODO per td.json |
| `consumer.go` — `handlePaymentExpired` | **STUB** — TODO per td.json |
| `consumer.go` — `handleOrderCancelled` | **STUB** — TODO per td.json |
| `consumer.go` — `handleProductCreated` | **STUB** — TODO per td.json |

## Lock-order rule (MANDATORY)

**ALWAYS lock `stock_levels(sku)` BEFORE `reservations(id)` inside any state-mutating transaction.**

Within a single multi-SKU `reservation.create` call, lock SKUs in **lexicographic order** (ORDER BY sku ASC) before inserting any reservation rows.

This is the single canonical deadlock-prevention pin from `td.json:concurrency_model.lock_ordering`.

Lock-order comments are placed at every lock point in the code.

## No-negative invariant (3-layer defence)

1. **Tx-level guard**: pre-UPDATE check on the locked row (`available_qty >= requested_qty`).
2. **DB CHECK constraints**: `CHECK (available_qty >= 0)`, `CHECK (reserved_qty >= 0)`, `CHECK (sold_qty >= 0)` on `stock_levels`.
3. **Adjust path**: `delta + current available_qty` computed under FOR UPDATE, rejected with `STOCK_NEGATIVE_INVARIANT` before UPDATE if would-go-negative.

## Sweeper

`sweeper.go` fires every `SWEEPER_CADENCE_SECONDS` (default 30).

Per tick:
1. **Candidate scan**: `SELECT ... WHERE status='RESERVED' AND expires_at < NOW() ORDER BY expires_at ASC LIMIT batch FOR UPDATE SKIP LOCKED` in a short tx. Commit releases the scan locks.
2. **Per-row tx**: re-lock `stock_levels(sku)` FIRST (lock-order pin), then `reservations(id)`. If status is no longer RESERVED, skip. Otherwise: flip to EXPIRED, `reserved_qty -= qty`, `available_qty += qty`, insert outbox row for `events.reservation.expired`.

The outbox relay goroutine picks up the outbox row and publishes to `ecom.inventory.events` with `key=orderId` (partition key).

## events.reservation.expired → Order service

The sweeper writes an outbox row; the relay goroutine publishes it to `ecom.inventory.events`. The Order service consumes this event to drive `PENDING_PAYMENT → PAYMENT_EXPIRED` transitions (see BA edge case PR-008).

## TD ambiguity resolutions

| Topic | Decision |
|---|---|
| `reservation.commit` HTTP | NOT exposed. Internal function called from `payment.completed` consumer. |
| `reservation.release` HTTP | EXPOSED. Checkout's compensation path requires it. |
| `product.created` initial qty | Seeded to 0. Admin uses `stock.adjust` to set inventory. |
| Returned `reservationId` | `reservationIds[0]` (first UUID in lex-order; all lookups use `order_id`). |
| `IDEMPOTENCY_KEY_INFLIGHT` | Not returned under default config; FOR UPDATE on `idempotency_keys` serializes waiters. |

## Environment variables

| Var | Default | Description |
|---|---|---|
| `PORT` | (required) | HTTP listen port |
| `DB_URL` | (required) | Full Postgres connection URL |
| `KAFKA_ENABLED` | `false` | Enable Kafka consumer + producer |
| `KAFKA_BROKERS` | | Consumer broker list (CSV) |
| `KAFKA_GROUP_ID` | | Consumer group ID |
| `KAFKA_PRODUCER_BROKERS` | | Producer broker list (CSV) |
| `SWEEPER_CADENCE_SECONDS` | `30` | Expiry sweep interval |
| `SWEEPER_BATCH_SIZE` | `200` | Max rows per sweep tick |
| `RESERVATION_TTL_MINUTES` | `15` | TTL written to `expires_at` on create |
| `JWT_PUBLIC_KEY` | (required) | Base64-encoded ES256 public key |
| `REF_ID_HEADER_KEY` | `X-Ref-Id` | Correlation ID header name |

## Running locally

```bash
docker compose up -d postgres
# apply migrations
psql "$DB_URL" -f migrations/001_init_stock_levels.up.sql
psql "$DB_URL" -f migrations/002_init_reservations.up.sql
psql "$DB_URL" -f migrations/003_init_stock_adjustments.up.sql
psql "$DB_URL" -f migrations/004_init_outbox_events.up.sql
psql "$DB_URL" -f migrations/005_init_consumed_events.up.sql
psql "$DB_URL" -f migrations/006_init_idempotency_keys.up.sql
PORT=8080 DB_URL=postgres://postgres:mypassword@localhost:5432/inventory go run .
```
