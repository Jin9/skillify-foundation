# POST /api/v1/order/internal/cancel-on-checkout-failure

## Summary

Internal endpoint (Checkout-only) used by the `checkout.commit` compensation matrix (ADR-007) when step 8 (`payment.intent.create`) fails AFTER the order row was created in step 7. Transitions the order from `PENDING_PAYMENT` → `CANCELLED` with `actor_role='SYSTEM'` and the supplied `reason`, and emits `events.order.cancelled` via the outbox so Inventory's state-driven consumer releases the reservation.

This handler MUST be a real implementation, NOT a 501 stub — REV-L2-004 closure (dry-run #1's compound dead-end). ADR-007 makes it explicit: "Reservation TTL is the safety net" — but the safety net is 15 minutes; this endpoint closes the gap immediately within the synchronous checkout-commit window.

## Story refs

- `STORY_CHECKOUT_COMMIT` (Checkout's compensation caller)
- `STORY_ORDER_CUSTOMER_CANCEL` (shares the same outbox + state-history shape)

## Contract ref

[`order.cancel-on-checkout-failure`](../../architecture/contracts.json#order.cancel-on-checkout-failure)

## Auth

- Tier: **internal_secret**.
- Header: `X-Internal-Secret: <secret>`. Validated via `crypto/subtle.ConstantTimeCompare` against `INTERNAL_SHARED_SECRET` (ADR-008). Rotation overlap supported.
- Plain `==`/`!=` is FORBIDDEN.

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["orderId","reason"],
  "properties": {
    "orderId": { "type": "string", "format": "uuid" },
    "reason":  {
      "type": "string",
      "minLength": 1,
      "maxLength": 255,
      "description": "e.g. 'payment intent creation failed: <upstream code>'; persisted on the cancellation status_history row's reason field."
    }
  }
}
```

Required headers:
- `X-Internal-Secret: <secret>` — constant-time-compared (ADR-008).
- `Idempotency-Key: <orderId>` — per `cross-cutting.idempotency.internal_orderId_rule`.

## Response (success)

HTTP 200, envelope:

```json
{
  "code": "SUCCESS",
  "message": "order cancelled",
  "data": {
    "orderId": "01935b9c-...",
    "status": "CANCELLED"
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `VALIDATION_ERROR` | 400 | `orderId` missing/malformed; `reason` missing/empty; `Idempotency-Key != orderId` |
| `AUTH_INVALID` | 401 | missing or wrong `X-Internal-Secret` (constant-time compare; ADR-008) |
| `NOT_FOUND` | 404 | `orderId` does not exist (defensive — should be unreachable in the compensation matrix because step 7 succeeded by the time step 8 runs, but defended) |
| `INVALID_ORDER_STATE` | 409 | current status is not `PENDING_PAYMENT` AND not CANCELLED-by-this-path (e.g. an event-driven transition raced ahead to `PAID` before compensation arrived) — Checkout's caller logic must escalate to admin-cancel |

## Business logic steps

1. **Internal auth middleware:** constant-time-compare. Reject → 401 `AUTH_INVALID`; log `internal_auth_reject`. (ADR-008)
2. **Idempotency-Key validation:** assert `Idempotency-Key == orderId`; mismatch → 400 `VALIDATION_ERROR`.
3. **Body validation:** `orderId` UUID; `reason` 1..255. Else 400 `VALIDATION_ERROR`.
4. `service/order_lifecycle.CancelOnCheckoutFailure(ctx, orderId, reason)`:
   1. `BEGIN;`
   2. `SELECT * FROM orders WHERE id=$1 FOR UPDATE;` Miss → ROLLBACK; 404 `NOT_FOUND`.
   3. **Idempotent replay branch:** if `status='CANCELLED'` AND most-recent `order_status_history` row has `actor_role='SYSTEM'` → COMMIT no-op; return 200 (idempotent compensation replay).
   4. If `status != 'PENDING_PAYMENT'` → ROLLBACK; 409 `INVALID_ORDER_STATE` (Checkout escalates to admin-cancel).
   5. `UPDATE orders SET status='CANCELLED', version=version+1, updated_at=NOW() WHERE id=$1 AND version=$old_version;`  *// ADR-003 sole-writer*
   6. `INSERT INTO order_status_history (order_id, from_status='PENDING_PAYMENT', to_status='CANCELLED', actor_user_id=NULL, actor_role='SYSTEM', reason=$request.reason, occurred_at=NOW());`  *// ORD-009*
   7. `INSERT INTO outbox_events (id=uuidv7(), aggregate_id=order_id, event_type='order.cancelled', payload_json={eventId, occurredAt, orderId, cancelActor:'SYSTEM', fromStatus:'PENDING_PAYMENT', reason:$request.reason});`  *// ADR-004 + events.order.cancelled.payload_shape*
   8. `COMMIT;`
5. Return 200 envelope.

## Side effects

- `UPDATE orders` — status `CANCELLED`, version+1.
- `INSERT order_status_history` — one row, actor=SYSTEM, reason populated.
- `INSERT outbox_events` — one row, `event_type='order.cancelled'`, `cancelActor='SYSTEM'`. Publisher ships to Kafka `ecom.order.events`.
- Downstream (eventual): Inventory state-driven consumer releases reservation. The published event carries `cancelActor=SYSTEM`; Inventory uses current-row read (not event.fromStatus) to decide release branch. Since Checkout's compensation matrix ALSO calls `inventory.reservation.release` directly (in parallel), Inventory's consumer typically observes RELEASED already and applies the idempotent no-op branch.

## Idempotency

- **Key:** `Idempotency-Key` header (which MUST equal `orderId`).
- **Replay with same payload:** 200 `SUCCESS` (no extra status_history row, no extra outbox row). Detected by checking `status='CANCELLED' AND latest_history.actor_role='SYSTEM'`.
- **Replay after order moved to PAID** (race window): 409 `INVALID_ORDER_STATE` — Checkout's caller logic must escalate to admin-cancel since Order is now PAID and the customer's payment succeeded after compensation was triggered. This race is documented in the contract's `idempotency_rules`.
- **TTL:** indefinite (status is terminal).

## Performance

- p95 < 200ms (compensation path; called inside checkout-commit's 500ms PERF-002 budget alongside `inventory.reservation.release`).
- Single tx with one indexed point lookup + 1 UPDATE + 1-2 INSERTs.

## Test cases

- `cancel_on_checkout_failure_pending_payment_to_cancelled_with_system_actor`
- `cancel_on_checkout_failure_idempotent_replay_returns_200_no_extra_history`
- `cancel_on_checkout_failure_status_paid_returns_409_INVALID_ORDER_STATE`
- `cancel_on_checkout_failure_status_already_cancelled_by_admin_returns_409_INVALID_ORDER_STATE`
- `cancel_on_checkout_failure_unknown_orderId_returns_404`
- `cancel_on_checkout_failure_missing_X_Internal_Secret_returns_401`
- `cancel_on_checkout_failure_wrong_X_Internal_Secret_returns_401_constant_time_compare`
- `cancel_on_checkout_failure_idempotency_key_neq_orderId_returns_400`
- `cancel_on_checkout_failure_reason_missing_returns_400`
- `cancel_on_checkout_failure_emits_outbox_with_cancelActor_SYSTEM_fromStatus_PENDING_PAYMENT`
- `cancel_on_checkout_failure_history_row_carries_supplied_reason_string`
- `cancel_on_checkout_failure_with_customer_jwt_returns_401_AUTH_INVALID`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Opus 4.7 1M, dry-run #2) | Created — REV-L2-004 closure; full implementation (no stub). Aligns with ADR-007 + ADR-008 + cross-cutting.idempotency.internal_orderId_rule. |
