# POST /api/v1/order/order/update-status-admin

## Summary

Admin-only state-machine transition endpoint. Accepts `(orderId, toStatus, [trackingNumber], [reason])` and validates the `(currentStatus, toStatus, ADMIN)` triple against the `state_machine.allowed_transitions` table inside a `SELECT ... FOR UPDATE` row lock. Allowed admin transitions: `PAID→PACKING`, `PACKING→SHIPPED` (trackingNumber required — ORD-006), `SHIPPED→DELIVERED`, `PENDING_PAYMENT→CANCELLED` (reason — ORD-005), `PAID→CANCELLED` (reason — ORD-005, restocks per PR-006). All other transitions return 409 `INVALID_ORDER_STATE` (ORD-003 + PR-005 + §9.3).

## Story refs

- `STORY_ORDER_ADMIN_TRANSITIONS`
- `STORY_ORDER_ADMIN_CANCEL`

## Contract ref

[`order.update-status-admin`](../../architecture/contracts.json#order.update-status-admin)

## Auth

- Tier: **admin_jwt**. CUSTOMER tokens → 403 `AUTH_FORBIDDEN` (AUTH-003).

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["orderId","toStatus"],
  "properties": {
    "orderId":        { "type": "string", "format": "uuid" },
    "toStatus":       { "type": "string", "enum": ["PACKING","SHIPPED","DELIVERED","CANCELLED"] },
    "trackingNumber": { "type": "string", "minLength": 1, "maxLength": 64 },
    "reason":         { "type": "string", "minLength": 1, "maxLength": 255 }
  },
  "allOf": [
    { "if": { "properties": { "toStatus": { "const": "SHIPPED" } } },   "then": { "required": ["trackingNumber"] } },
    { "if": { "properties": { "toStatus": { "const": "CANCELLED" } } }, "then": { "required": ["reason"] } }
  ]
}
```

## Response (success)

HTTP 200, envelope:

```json
{
  "code": "SUCCESS",
  "message": "status updated",
  "data": {
    "orderId": "01935b9c-...",
    "status": "SHIPPED"
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `VALIDATION_ERROR` | 400 | missing `trackingNumber` on `toStatus=SHIPPED` (ORD-006); missing `reason` on `toStatus=CANCELLED` (ORD-005); `toStatus` not in admin-allowed enum |
| `AUTH_MISSING` | 401 | no `Authorization` header |
| `AUTH_INVALID` | 401 | JWT signature/exp/iss/aud invalid |
| `AUTH_FORBIDDEN` | 403 | `claims.role != ADMIN` (AUTH-003) |
| `NOT_FOUND` | 404 | `orderId` does not exist |
| `INVALID_ORDER_STATE` | 409 | `(currentStatus, toStatus, ADMIN)` not in `state_machine.allowed_transitions`. Examples: `PACKING→CANCELLED`, `SHIPPED→PAID`, `DELIVERED→CANCELLED`, `PENDING_PAYMENT→PACKING`, `PAID→SHIPPED` (skips PACKING), terminal states (`PAYMENT_FAILED`/`PAYMENT_EXPIRED`/`CANCELLED`/`DELIVERED`) → anything (ORD-003, PR-005, §9.3) |

## Business logic steps

1. `transport/http` admin-decorator: reject 403 `AUTH_FORBIDDEN` if `claims.role != ADMIN`.
2. Bind + validate body. JSON Schema `allOf` enforces conditional required-fields (trackingNumber on SHIPPED, reason on CANCELLED). Else 400 `VALIDATION_ERROR`.
3. `service/order_lifecycle.UpdateStatusByAdmin(ctx, orderId, toStatus, trackingNumber?, reason?, claims.sub)`:
   1. `BEGIN;`
   2. `SELECT * FROM orders WHERE id=$1 FOR UPDATE;` Miss → ROLLBACK; 404 `NOT_FOUND`.
   3. **Idempotent retry guard:** if `status = toStatus` already → INSERT `admin_action_log {outcome_code='SUCCESS', request_summary={idempotent:true}}`; COMMIT no-op (no extra status_history row, no extra outbox row); return 200.
   4. `domain/transitions.ValidateTransition(current=row.status, target=toStatus, actor='ADMIN', provided={trackingNumber?, reason?})`. Miss → INSERT `admin_action_log {outcome_code='INVALID_ORDER_STATE', ...}`; COMMIT; return 409 `INVALID_ORDER_STATE`. (Audit even on rejection.)
   5. `UPDATE orders SET status=$2, [tracking_number=$3 if toStatus='SHIPPED'], version=version+1, updated_at=NOW() WHERE id=$1 AND version=$old_version;`  *// ADR-003 sole-writer of orders.status*
   6. `INSERT INTO order_status_history (order_id, from_status=current, to_status=$2, actor_user_id=claims.sub, actor_role='ADMIN', reason=$reason, occurred_at=NOW());`
   7. If `toStatus='CANCELLED'`: `INSERT INTO outbox_events (id, aggregate_id=order_id, event_type='order.cancelled', payload_json={eventId, occurredAt, orderId, cancelActor:'ADMIN', fromStatus:current, reason:$reason});`  *// ADR-004 + events.order.cancelled*
   8. `INSERT INTO admin_action_log (actor_user_id=claims.sub, endpoint='order.update-status-admin', order_id, request_summary={toStatus, trackingNumber?, reason?}, outcome_code='SUCCESS');`
   9. `COMMIT;`
4. Wrap with `common/wrapper.Success`.

## Side effects

- `UPDATE orders` — `status` (+ `tracking_number` if SHIPPED), version+1.
- `INSERT order_status_history` — one row per successful transition; one row not written for idempotent retries.
- `INSERT outbox_events` — only for `toStatus=CANCELLED` (both fromStatus=PENDING_PAYMENT and fromStatus=PAID origins). Publisher ships to Kafka `ecom.order.events`.
- `INSERT admin_action_log` — one row per call (success, idempotent, AND rejection paths) — full audit ledger.
- Downstream (eventual, for CANCELLED only): Inventory state-driven consumer releases reservation. For `fromStatus=PAID` PR-006 restock policy applies (sold→available). Payment consumer transitions intent if still active.

## Idempotency

- Re-issuing the same `(orderId, toStatus)` when row is already at `toStatus` → 200 `SUCCESS` (audit row written; no extra status_history row, no extra outbox row).
- Optional `Idempotency-Key` header per the contract (admin retry safety) — MVP relies on the FOR UPDATE + status-comparison; the explicit `idempotency_keys` store is not required at this endpoint.

## Performance

- p95 < 300ms.
- Single tx with one indexed point lookup + 2-4 INSERTs (or 1 audit row only on rejection).

## Test cases

- `update_status_admin_paid_to_packing_succeeds`
- `update_status_admin_packing_to_shipped_with_tracking_succeeds`
- `update_status_admin_shipped_to_delivered_succeeds`
- `update_status_admin_pending_payment_to_cancelled_with_reason_succeeds`
- `update_status_admin_paid_to_cancelled_with_reason_succeeds_and_emits_outbox`
- `update_status_admin_packing_to_shipped_without_tracking_returns_400_ORD006`
- `update_status_admin_cancelled_without_reason_returns_400_ORD005`
- `update_status_admin_packing_to_cancelled_returns_409_INVALID_ORDER_STATE_PR005`
- `update_status_admin_shipped_to_cancelled_returns_409`
- `update_status_admin_delivered_to_cancelled_returns_409`
- `update_status_admin_shipped_to_paid_returns_409_backwards`
- `update_status_admin_paid_to_shipped_returns_409_skips_packing`
- `update_status_admin_pending_payment_to_packing_returns_409_skips_paid`
- `update_status_admin_payment_failed_to_anything_returns_409_terminal`
- `update_status_admin_payment_expired_to_anything_returns_409_terminal`
- `update_status_admin_cancelled_to_anything_returns_409_terminal`
- `update_status_admin_idempotent_same_status_returns_200_no_extra_history`
- `update_status_admin_customer_token_returns_403_AUTH_FORBIDDEN`
- `update_status_admin_audit_row_written_on_rejection_path` (INVALID_ORDER_STATE)
- `update_status_admin_concurrent_two_admins_one_succeeds_one_409` (FOR UPDATE serialization)

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Opus 4.7 1M, dry-run #2) | Created |
