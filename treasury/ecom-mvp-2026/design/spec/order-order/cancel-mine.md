# POST /api/v1/order/order/cancel-mine

## Summary

Customer-initiated cancellation of their OWN order while it is still `PENDING_PAYMENT` (ORD-004). Cancels the order, writes the `order_status_history` row, and emits `events.order.cancelled` via the outbox. Inventory's state-driven consumer releases the reservation; Payment's consumer transitions the intent. Re-cancelling a CANCELLED-by-self order is idempotent (200); re-cancelling a CANCELLED-by-admin order via this endpoint returns 409 (the customer no longer owns the action — AMBIG-ORD-3 resolution).

## Story refs

- `STORY_ORDER_CUSTOMER_CANCEL`

## Contract ref

[`order.cancel-mine`](../../architecture/contracts.json#order.cancel-mine)

## Auth

- Tier: **customer_jwt**.
- Server uses combined ownership-check + lock (`SELECT ... WHERE id=$1 AND user_id=claims.sub FOR UPDATE`). Misses → 404 `ORDER_NOT_OWNED` (no enumeration leak).

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["orderId"],
  "properties": {
    "orderId": { "type": "string", "format": "uuid" }
  }
}
```

## Response (success)

HTTP 200, envelope:

```json
{
  "code": "SUCCESS",
  "message": "order cancelled",
  "data": {
    "orderId": "01935b9c-...-uuidv7",
    "status": "CANCELLED"
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `VALIDATION_ERROR` | 400 | `orderId` missing or not UUID-shaped |
| `AUTH_MISSING` | 401 | no `Authorization` header |
| `AUTH_INVALID` | 401 | JWT signature/exp/iss/aud invalid |
| `ORDER_NOT_OWNED` | 404 | id not found OR not owned by `claims.sub` (single 404 — ORD-001 no enumeration) |
| `INVALID_ORDER_STATE` | 409 | current status != `PENDING_PAYMENT` (ORD-004), e.g. trying to cancel a `PAID` order via this endpoint, OR re-cancelling a CANCELLED-by-admin order via this endpoint (AMBIG-ORD-3) |

## Business logic steps

1. Bind + validate body. UUID format check.
2. Read `claims.sub` from context.
3. `service/order_lifecycle.CancelByCustomer(ctx, orderId, claims.sub)`:
   1. `BEGIN;`
   2. `SELECT * FROM orders WHERE id=$1 AND user_id=$2 FOR UPDATE;` (combined ownership + row lock). Miss → ROLLBACK; 404 `ORDER_NOT_OWNED`.
   3. If `status='CANCELLED'`: read most-recent `order_status_history` row. If `actor_user_id = claims.sub` → COMMIT no-op; return 200 (idempotent re-cancel by self per `cancel-mine.idempotency_rules` + AMBIG-ORD-3). Else ROLLBACK; 409 `INVALID_ORDER_STATE`.
   4. If `status != 'PENDING_PAYMENT'` → ROLLBACK; 409 `INVALID_ORDER_STATE`.
   5. `UPDATE orders SET status='CANCELLED', version=version+1, updated_at=NOW() WHERE id=$1 AND version=$old_version;`  *// ADR-003 sole-writer of orders.status*
   6. `INSERT INTO order_status_history (order_id, from_status='PENDING_PAYMENT', to_status='CANCELLED', actor_user_id=claims.sub, actor_role='CUSTOMER', reason=NULL, occurred_at=NOW());`  *// ORD-009*
   7. `INSERT INTO outbox_events (id, aggregate_id=order_id, event_type='order.cancelled', payload_json={eventId:id, occurredAt:NOW(), orderId, cancelActor:'CUSTOMER', fromStatus:'PENDING_PAYMENT'}, created_at=NOW());`  *// ADR-004 outbox + events.order.cancelled.payload_shape*
   8. `COMMIT;`
4. Wrap with `common/wrapper.Success`.

## Side effects

- `UPDATE orders` — status `CANCELLED`, version+1.
- `INSERT order_status_history` — one row, actor=CUSTOMER.
- `INSERT outbox_events` — one row, event_type='order.cancelled', cancelActor='CUSTOMER', fromStatus='PENDING_PAYMENT'. Publisher goroutine (transport/outbox_publisher) ships to Kafka topic `ecom.order.events` (partition_key=orderId).
- Downstream (eventual): Inventory state-driven consumer releases the reservation (RESERVED→available). Payment state-driven consumer transitions the intent to CANCELLED (if still REQUIRES_PAYMENT).

## Idempotency

- Re-cancelling a CANCELLED-by-self order → 200 `SUCCESS` (no extra history row, no extra outbox row).
- Re-cancelling a CANCELLED-by-admin or CANCELLED-by-system order → 409 `INVALID_ORDER_STATE` (terminal state; the customer no longer owns the action — AMBIG-ORD-3 resolution).
- Concurrency: spawn N goroutines that all attempt `cancel-mine` on the same PENDING_PAYMENT order; the FOR UPDATE row lock serializes them; exactly one writes a status_history row and an outbox row; the rest see CANCELLED-by-self and return 200 idempotent.

## Performance

- p95 < 300ms (BA `non_functional.latency`).
- Single transaction with one indexed point lookup + 1-2 INSERTs.

## Test cases

- `cancel_mine_pending_payment_cancels_and_emits_outbox`
- `cancel_mine_paid_order_returns_409_INVALID_ORDER_STATE`
- `cancel_mine_already_cancelled_by_self_returns_200_idempotent`
- `cancel_mine_already_cancelled_by_admin_returns_409`
- `cancel_mine_other_users_order_returns_404_ORDER_NOT_OWNED`
- `cancel_mine_unknown_order_returns_404_ORDER_NOT_OWNED`
- `cancel_mine_concurrent_attempts_one_succeeds_others_idempotent_or_409`
- `cancel_mine_outbox_row_written_in_same_tx_as_status_history`
- `PR003_cancel_mine_loses_to_payment_completed_returns_409` (race-resolution scenario)

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Opus 4.7 1M, dry-run #2) | Created |
