# Order Service

Owns the order lifecycle state machine. **Sole writer of `orders.status`** — no other service updates that column.

## What is FULL vs STUB

| Endpoint / Consumer | Status |
|---|---|
| `POST /api/v1/order/order/create-from-checkout` | FULL |
| `POST /api/v1/order/order/detail` | FULL |
| `POST /api/v1/order/order/list-mine` | FULL |
| `POST /api/v1/order/order/cancel-mine` | STUB (501) |
| `POST /api/v1/order/order/list-admin` | STUB (501) |
| `POST /api/v1/order/order/update-status-admin` | STUB (501) |
| `POST /api/v1/order/order/cancel-on-checkout-failure` | STUB (501) |
| Consumer `events.payment.completed` | FULL |
| Consumer `events.payment.failed` | STUB (log + ack) |
| Consumer `events.payment.expired` | STUB (log + ack) |
| Consumer `events.reservation.expired` | STUB (log + ack) |

## State Machine

Encoded as a 2-D map in `app/order/state_machine.go` keyed by `(from, to, actor)` struct.
O(1) lookup via `ValidateTransition(from, to, actor)` — single source of truth for ALL transition paths.

| From | To | Actor |
|---|---|---|
| PENDING_PAYMENT | PAID | SYSTEM |
| PENDING_PAYMENT | PAYMENT_FAILED | SYSTEM |
| PENDING_PAYMENT | PAYMENT_EXPIRED | SYSTEM |
| PENDING_PAYMENT | CANCELLED | CUSTOMER |
| PENDING_PAYMENT | CANCELLED | ADMIN |
| PENDING_PAYMENT | CANCELLED | SYSTEM |
| PAID | PACKING | ADMIN |
| PACKING | SHIPPED | ADMIN |
| SHIPPED | DELIVERED | ADMIN |
| PAID | CANCELLED | ADMIN |

Terminal states (no outbound transitions): `PAYMENT_FAILED`, `PAYMENT_EXPIRED`, `DELIVERED`, `CANCELLED`.

## Internal Auth Path

`create-from-checkout` and `cancel-on-checkout-failure` are **internal-only** — accessible only by the Checkout service.

Authentication: `X-Internal-Secret` header must equal `$INTERNAL_SHARED_SECRET`. Reject with `401 ORD401` otherwise.

Hardening path (deferred, AMBIG-ORD-1): replace with mTLS or SERVICE-role JWT per cross-cutting.route-convention.

## Order Number Sequence

Global Postgres sequence `order.order_number_seq` (not per-day). Format: `ORD-YYYYMMDD-NNNNNN`.

- `NNNNNN` does NOT reset at midnight — ORD-008 only requires the shape, not daily reset.
- Allocation is race-free by sequence semantics; `UNIQUE` constraint on `order_number` is the safety net.
- See `migrations/002_create_order_number_seq.up.sql` and `access/storage_order.go#NextOrderNumber`.

## Kafka Modes

| Variable | Description |
|---|---|
| `KAFKA_ENABLED=false` | Consumer disabled; service runs HTTP-only |
| `KAFKA_BROKERS` | Consumer broker list (CSV) |
| `KAFKA_PRODUCER_BROKERS` | Producer broker list (empty = outbox publisher disabled) |
| `KAFKA_GROUP_ID` | Consumer group ID |
| `KAFKA_TOPICS` | Topics to subscribe (CSV): `ecom.payment.events,ecom.inventory.events` |

## Consumer Dispatch Summary

`consumer.go#RegisterConsumerRoutes` returns a `map[string]kafka.KafkaHandler` keyed by `EventName`:

- `payment.completed` → `onPaymentCompleted` (FULL: two-layer idempotency, FOR UPDATE, state-driven guard)
- `payment.failed` → `onPaymentFailed` (STUB: logs TODO)
- `payment.expired` → `onPaymentExpired` (STUB: logs TODO)
- `reservation.expired` → `onReservationExpired` (STUB: logs TODO)

The router wires these via `kafka.NewEventRouter` (common/kafka), which dispatches by `envelope.EventName`.

## TD Ambiguity Handling

**AMBIG-ORD-1 (internal auth):** Implemented shared-secret header (`X-Internal-Secret`) — a minor scope addition over pure network-plane trust. Reviewer-L2 to ratify or downgrade.

**AMBIG-ORD-2 (buyer_email_snapshot):** This service assumes Checkout passes `buyerEmail` in the `create-from-checkout` request body. Checkout must load it from `identity.profile.read` before calling this endpoint. Contracts.json needs clarification.

**AMBIG-ORD-3 (idempotent re-cancel):** Implementation note in `handler_cancel_mine.go` stub: idempotent 200 only when last history row actor_user_id matches claims.sub; otherwise 409 for CANCELLED-by-other.

**Path (a) decision:** Confirmed — `order.create-from-checkout` is a synchronous internal HTTP endpoint (not Kafka event), as required by contracts.json line 892 (checkout.commit.internal_orchestration step 6) for the frontend to receive orderId in the commit response.

## Running Locally

```bash
cp .env.example .env   # fill in secrets
docker compose up -d postgres kafka
go run .
```

## Migrations

Apply with any golang-migrate compatible tool:

```bash
migrate -path ./migrations -database "postgres://order:order@localhost:5434/orderdb?search_path=order&sslmode=disable" up
```
