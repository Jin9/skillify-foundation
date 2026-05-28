# POST /api/v1/inventory/reservation/release

## Summary

Compensation endpoint. Called by Checkout when a post-step-6 failure of `checkout.commit` requires releasing a previously-created reservation. Per `cross-cutting.persistence.transactional_boundaries` and ADR-007:

- After `inventory.reservation.create` succeeds, if `order.create-from-checkout` (step 7) or `payment.intent.create` (step 8) fails, Checkout MUST release the reservation as part of its compensation matrix.
- The endpoint is also reachable internally from event consumers, BUT consumers call `ReservationService.Release` **directly** (not over HTTP) to avoid loopback overhead. The HTTP route exists exclusively for the cross-process Checkout-compensation path.

The handler is **idempotent on `order_id`** and applies a state-driven branch table per row, mirroring the `events.order.cancelled` consumer's branches. This means the release endpoint correctly handles the rare case where Checkout retries a release after the sweeper / event consumer already released the row.

## Story refs

- `STORY_INVENTORY_RELEASE_ON_FAILURE` (the compensation half of the order lifecycle)
- `STORY_CHECKOUT_COMMIT` (caller — the compensation matrix in step 7/8 failure handling)

## Contract ref

[`inventory.reservation.release`](../../architecture/contracts.json#inventory.reservation.release)

## Auth

- Tier: **internal_secret**.
- Header: `X-Internal-Secret: <32-byte secret>`.
- Constant-time compare per [ADR-008](../../architecture/ADRs/ADR-008-internal-shared-secret.md).
- Caller per contract: **checkout** (compensation path only). Event consumers do NOT use this HTTP route.

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["orderId", "reason"],
  "properties": {
    "orderId":       { "type": "string", "format": "uuid" },
    "reservationId": { "type": "string", "format": "uuid", "description": "Optional; informational. Lookup is by orderId — one orderId may have multiple reservation rows." },
    "reason":        { "type": "string", "enum": ["PAYMENT_FAILED", "PAYMENT_EXPIRED", "ORDER_CANCELLED", "ADMIN_FORCE"] }
  }
}
```

## Response (success)

HTTP 200, envelope:

```json
{
  "code": "SUCCESS",
  "message": "released",
  "data": {
    "reservationId": "01935b9c-8a9e-7c3f-9d11-22aa33bb44cc",
    "status": "RELEASED"
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

`reservationId` echoed back is the first-row id (deterministic by `ORDER BY sku ASC`).

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `VALIDATION_ERROR` | 400 | `orderId` malformed or `reason` not in enum |
| `AUTH_INVALID` | 401 | `X-Internal-Secret` missing or mismatched |
| `RESERVATION_NOT_FOUND` | 404 | No `reservations` row for `order_id` |
| `DATABASE_UNAVAILABLE` | 503 | Postgres unreachable |

## Business logic steps

1. `common/middleware.InternalAuth(secret)` validates `X-Internal-Secret`.
2. Bind + validate body: `orderId` UUID, `reason ∈ {PAYMENT_FAILED, PAYMENT_EXPIRED, ORDER_CANCELLED, ADMIN_FORCE}`.
3. `BEGIN TX`.
4. `SELECT * FROM reservations WHERE order_id = $1 ORDER BY sku ASC FOR UPDATE` — multi-row lock, ordered by sku to honor the lock-order pin.
5. If `rows` is empty → `ROLLBACK`; return `RESERVATION_NOT_FOUND`.
6. **State-driven branch per row** (mirrors the `events.order.cancelled` consumer; the function is shared via `ReservationService.Release`):
   - `RELEASED` or `EXPIRED` → continue (idempotent no-op).
   - `RESERVED`:
     - `SELECT * FROM stock_levels WHERE sku = r.sku FOR UPDATE`.
     - `UPDATE stock_levels SET available_qty += r.qty, reserved_qty -= r.qty, version += 1 WHERE sku = r.sku`.
     - `UPDATE reservations SET status = 'RELEASED', released_at = NOW(), release_reason = $reason WHERE id = r.id`.
   - `COMMITTED` (BA PR-006 soft-cancel-of-PAID branch — also used in the PR-003 race resolution):
     - `SELECT * FROM stock_levels WHERE sku = r.sku FOR UPDATE`.
     - `UPDATE stock_levels SET sold_qty -= r.qty, available_qty += r.qty, version += 1 WHERE sku = r.sku`.
     - `UPDATE reservations SET status = 'RELEASED', released_at = NOW(), release_reason = 'PAID_SOFT_CANCEL' WHERE id = r.id`.
7. `COMMIT`.
8. `wrapper.Respond(c, SUCCESS, {reservationId: rows[0].id, status: 'RELEASED'})`.

## Side effects

For each row matched on `order_id`, depending on prior state:
- RESERVED row: `UPDATE stock_levels` (`available += qty`, `reserved -= qty`) + `UPDATE reservations` (status, released_at, release_reason = caller-supplied reason).
- COMMITTED row: `UPDATE stock_levels` (`sold -= qty`, `available += qty`) + `UPDATE reservations` (status, released_at, `release_reason = 'PAID_SOFT_CANCEL'`).
- RELEASED / EXPIRED row: no change (idempotent no-op).
- NO outbox event. (Inventory does NOT emit a release event; downstream callers either know about the release already (Checkout compensation) or learned via the event they consumed.)

## Idempotency

- **Idempotent on `order_id`** per the contract.
- Per-row idempotency on `reservations.id`: re-invocation on already-RELEASED/EXPIRED rows is a no-op.
- COMMITTED-branch idempotency is naturally safe — once the row is RELEASED, subsequent calls hit the RELEASED no-op branch.
- Two concurrent release calls on the same `order_id` serialize via `FOR UPDATE` on the reservations rows; the second waiter sees status `RELEASED` and no-ops.

## Performance

- p95 < 150ms for typical 1–5 SKU orders (single tx with per-row FOR UPDATE).
- Expected QPS at MVP: << 1 (Checkout's compensation path is the dominant caller; only fires on commit failures).

## Test cases

### Happy paths

- `reservation_release_RESERVED_to_RELEASED_increments_available_decrements_reserved`
- `reservation_release_COMMITTED_to_RELEASED_increments_available_decrements_sold_release_reason_PAID_SOFT_CANCEL` (PR-006)
- `reservation_release_multi_sku_atomic_per_tx`
- `reservation_release_idempotent_on_already_released_no_state_change`
- `reservation_release_idempotent_on_already_expired_no_state_change`

### Validation

- `reservation_release_invalid_reason_returns_VALIDATION_ERROR`
- `reservation_release_malformed_orderId_returns_VALIDATION_ERROR`

### Auth

- `reservation_release_x_internal_secret_missing_returns_AUTH_INVALID`

### Edge cases

- `reservation_release_no_reservation_for_orderId_returns_RESERVATION_NOT_FOUND`
- `reservation_release_concurrent_with_event_consumer_serializes_via_FOR_UPDATE_one_path_releases_other_no_ops`
- `reservation_release_envelope_carries_traceId_on_every_branch`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Opus 4.7 1M) | Created |
