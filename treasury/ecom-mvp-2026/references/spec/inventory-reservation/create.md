# POST /api/v1/inventory/reservation/create

## Summary

Multi-SKU stock reservation. Called by Checkout (step 6 of `checkout.commit`) over the private network plane to atomically lock stock for an order. Implements the reservation half of INV-002 / INV-003 / INV-005 / INV-008 and the BA edge case "two concurrent reservations on the last unit".

This is the **load-bearing concurrency endpoint** for the entire system. Every line of the algorithm exists for a specific reason — see `## Business logic steps` and `## Sequence diagram` for the full picture, and ADR-005 for the lock-order pin.

## Story refs

- `STORY_INVENTORY_RESERVE_AND_CONVERT`
- `STORY_CHECKOUT_COMMIT` (caller — checkout's step-6 sync HTTP)

## Contract ref

[`inventory.reservation.create`](../../architecture/contracts.json#inventory.reservation.create)

## Auth

- Tier: **internal_secret**.
- Header: `X-Internal-Secret: <32-byte secret>`.
- Constant-time compare per [ADR-008](../../architecture/ADRs/ADR-008-internal-shared-secret.md).
- Caller per contract: **checkout** (only).
- Idempotency-Key (in body) MUST equal `orderId` per `cross-cutting.idempotency.internal_orderId_rule` — recipient validates and 400s on mismatch.

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["orderId", "items", "expiresAt", "idempotencyKey"],
  "properties": {
    "orderId": {
      "type": "string",
      "format": "uuid",
      "description": "UUID v7 — the new order id minted by Checkout for this commit"
    },
    "items": {
      "type": "array",
      "minItems": 1,
      "maxItems": 50,
      "items": {
        "type": "object",
        "additionalProperties": false,
        "required": ["sku", "qty"],
        "properties": {
          "sku": { "type": "string",  "minLength": 1, "maxLength": 64 },
          "qty": { "type": "integer", "minimum": 1 }
        }
      }
    },
    "expiresAt": {
      "type": "string",
      "format": "date-time",
      "description": "RFC3339; recommended now + RESERVATION_TTL_MINUTES (= 15)"
    },
    "idempotencyKey": {
      "type": "string",
      "format": "uuid",
      "description": "MUST equal orderId — internal_orderId_rule"
    }
  }
}
```

## Response (success)

HTTP 200 (first call) or HTTP 200 (idempotent replay), envelope:

```json
{
  "code": "SUCCESS",
  "message": "reserved",
  "data": {
    "reservationId": "01935b9c-8a9e-7c3f-9d11-22aa33bb44cc",
    "items": [
      { "sku": "SKU-A", "qty": 2, "reserved": true },
      { "sku": "SKU-B", "qty": 1, "reserved": true }
    ],
    "expiresAt": "2026-05-08T15:30:00Z"
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

`reservationId` is the first-row id from the `reservations` table (deterministic by `ORDER BY sku ASC`). All downstream lookups (commit / release / sweeper / event consumers) are by `orderId`, NOT by this id — `reservationId` is informational/audit.

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `VALIDATION_ERROR` | 400 | Empty `items`, any `qty < 1`, `idempotencyKey != orderId`, malformed `expiresAt`, or `items.length > 50` |
| `AUTH_INVALID` | 401 | `X-Internal-Secret` missing or mismatched |
| `INSUFFICIENT_STOCK` | 409 | Any `item.qty > stock.available_qty` (response payload carries per-item details — see below) |
| `STOCK_NEGATIVE_INVARIANT` | 409 | Defensive — concurrent racer slipped past the FOR UPDATE; CHECK constraint fires |
| `IDEMPOTENCY_KEY_REUSED` | 409 | Same `orderId` + DIFFERENT request body hash |
| `DATABASE_UNAVAILABLE` | 503 | Postgres unreachable |

### `INSUFFICIENT_STOCK` payload shape

```json
{
  "code": "INSUFFICIENT_STOCK",
  "message": "one or more SKUs lack available stock",
  "data": {
    "details": [
      { "sku": "SKU-A", "requestedQty": 3, "availableQty": 2 }
    ]
  },
  "traceId": "..."
}
```

Per BA edge-case: when the customer asks for `qty=5` against `availableQty=2`, the response surfaces both numbers so Checkout can render a useful client error.

## Business logic steps

1. `common/middleware.InternalAuth(secret)` validates `X-Internal-Secret` (constant-time compare).
2. Bind + validate body (`common/validator`): `orderId` is UUID, `items.length ∈ [1, 50]`, every `qty >= 1`, `expiresAt` is RFC3339, `idempotencyKey == orderId` (else 400 per internal_orderId_rule).
3. Compute `request_hash = sha256(canonical_json({orderId, items[sorted by sku]}))`.
4. **Application-level sort:** `sort.Slice(items, items[i].Sku < items[j].Sku)` — establishes the lex order the SQL will respect.
5. `BEGIN TX`.
6. **Idempotency check (FOR UPDATE serializes concurrent same-key callers):**
   - `SELECT * FROM idempotency_keys WHERE key = $orderId AND endpoint = 'inventory.reservation.create' FOR UPDATE`.
   - If row found AND `request_hash` matches → `COMMIT`; return persisted `response_envelope` verbatim (idempotent replay; HTTP 200).
   - If row found AND `request_hash` differs → `ROLLBACK`; return 409 `IDEMPOTENCY_KEY_REUSED`.
   - Else proceed.
7. **Lock-order pin (ADR-005):** `SELECT sku, available_qty, reserved_qty, version FROM stock_levels WHERE sku = ANY($sortedSkus) ORDER BY sku ASC FOR UPDATE` — acquires per-row locks in lex sku order.
8. **Per-SKU validation (Layer-1 no-negative guard):** for each `item` in `items`, look up `row = stockRows[item.sku]`:
   - If `row` is `nil` → `insufficientList += {sku, requestedQty, availableQty: 0}` (treat as `INSUFFICIENT_STOCK`; PRODUCT_DELETED is checked upstream by checkout).
   - Else if `row.available_qty < item.qty` → `insufficientList += {sku, requestedQty, availableQty: row.available_qty}`.
   - Else proceed.
9. If `insufficientList` is non-empty → `ROLLBACK`; return 409 `INSUFFICIENT_STOCK` with `data.details = insufficientList`.
10. **Mutate per SKU (still inside the tx; rows are already locked):**
    For each `item` in `items` ordered `ASC by sku`:
    - `UPDATE stock_levels SET available_qty = available_qty - item.qty, reserved_qty = reserved_qty + item.qty, version = version + 1, updated_at = NOW() WHERE sku = item.sku`.
    - The `CHECK (available_qty >= 0)` constraint is the Layer-2 backstop; if Layer-1 missed an edge case, PG raises 23514 → mapped to `STOCK_NEGATIVE_INVARIANT`.
    - `rid = uuidv7()`; `INSERT INTO reservations (id, order_id, sku, qty, status, expires_at) VALUES (rid, $orderId, item.sku, item.qty, 'RESERVED', $expiresAt)`.
    - Append `rid` to `reservationIds`.
11. Build `responseData = {reservationId: reservationIds[0], items: items.map(reserved=true), expiresAt}` and `envelope = {code: 'SUCCESS', message: 'reserved', data: responseData, traceId}`.
12. **Persist idempotency record:** `INSERT INTO idempotency_keys (key=$orderId, endpoint='inventory.reservation.create', request_hash, response_envelope=$envelope, http_status=200, expires_at=NOW()+24h)`.
13. `COMMIT`.
14. Return `envelope`.

## Side effects

- For each item in the request: `UPDATE stock_levels` (`available -= qty`, `reserved += qty`) + `INSERT reservations` (status `RESERVED`).
- One `INSERT idempotency_keys` row keyed on `orderId`.
- NO outbox event (the lifecycle event `events.reservation.expired` is emitted later by the sweeper, not on the create path).

## Idempotency

- Server-side keyed on `(key=orderId, endpoint='inventory.reservation.create')`.
- TTL: 24 hours per `cross-cutting.idempotency.server_side_store.ttl`.
- First call: process normally; persist envelope inside the same tx as the business writes.
- Replay (same `orderId` + same body hash): return persisted envelope verbatim. HTTP 200.
- Same `orderId` + different body hash: 409 `IDEMPOTENCY_KEY_REUSED`.
- Concurrent same-key callers: `SELECT ... FOR UPDATE` on `idempotency_keys` serializes them; the second waiter blocks until the first commit then takes the replay branch. We do NOT return `IDEMPOTENCY_KEY_INFLIGHT` under default config (Checkout rarely double-fires the same `orderId`); the code is reserved for a future config-tunable lock-timeout.

## Performance

- p95 < 200ms for typical 1–5 SKU carts (single tx, FOR UPDATE on each SKU row, INSERTs).
- Hot-SKU contention is the throughput ceiling on a per-SKU basis; concurrent reservations on the same SKU serialize naturally (this is the desired behavior — INV-008).
- Expected QPS at MVP: 5–20 on the spine, bursty under flash-sale conditions.

## Sequence diagram

The lex-lock + per-SKU validation flow, including the IDEMPOTENT REPLAY branch and the INSUFFICIENT_STOCK ROLLBACK branch:

```mermaid
sequenceDiagram
  autonumber
  participant CO as Checkout
  participant INV as Inventory<br/>handler
  participant DB as Postgres<br/>(inventory schema)

  CO->>+INV: POST /reservation/create<br/>Hdr: X-Internal-Secret<br/>Body: {orderId, items=[{A,2},{B,1}], expiresAt, idempotencyKey=orderId}
  INV->>INV: middleware.InternalAuth (constant-time compare)
  INV->>INV: validate body; idempotencyKey == orderId
  INV->>INV: request_hash = sha256(canonical_json)
  INV->>INV: sort.Slice(items, sku ASC)<br/>// [SKU-A, SKU-B]

  INV->>+DB: BEGIN TX
  DB-->>-INV: ok

  Note over INV,DB: Step A — idempotency check (serializes concurrent same-key callers)
  INV->>+DB: SELECT * FROM idempotency_keys<br/>WHERE key=$orderId AND endpoint='...create' FOR UPDATE
  alt existing row + matching hash
    DB-->>INV: row(envelope, hash)
    INV->>+DB: COMMIT
    DB-->>-INV: ok
    INV-->>CO: 200 {SUCCESS, persisted envelope}<br/>(idempotent replay)
  else existing row + DIFFERENT hash
    DB-->>INV: row(envelope, hash')
    INV->>+DB: ROLLBACK
    DB-->>-INV: ok
    INV-->>CO: 409 IDEMPOTENCY_KEY_REUSED
  else no existing row
    DB-->>-INV: nil
    Note over INV,DB: Step B — LOCK ORDER PIN (ADR-005)<br/>stock_levels FIRST, in lex sku order
    INV->>+DB: SELECT sku, available_qty, reserved_qty, version<br/>FROM stock_levels WHERE sku = ANY([SKU-A, SKU-B])<br/>ORDER BY sku ASC FOR UPDATE
    DB-->>-INV: rows[A, B] locked

    Note over INV: Step C — per-SKU validation (Layer-1 guard)
    INV->>INV: for each item: check row.available_qty >= item.qty<br/>build insufficientList
    alt insufficient on any SKU
      INV->>+DB: ROLLBACK
      DB-->>-INV: ok
      INV-->>CO: 409 INSUFFICIENT_STOCK<br/>{details: [{sku, requestedQty, availableQty}]}
    else all sufficient
      Note over INV,DB: Step D — mutate per SKU in lex sku order
      loop for each item ASC by sku
        INV->>+DB: UPDATE stock_levels<br/>SET available_qty -= qty,<br/>    reserved_qty += qty,<br/>    version += 1<br/>WHERE sku = item.sku
        DB-->>-INV: 1 row<br/>(CHECK constraint Layer-2 active)
        INV->>+DB: INSERT reservations<br/>(id=uuidv7, order_id, sku, qty,<br/>  status='RESERVED', expires_at)
        DB-->>-INV: ok
      end

      Note over INV,DB: Step E — persist idempotency record (same tx)
      INV->>+DB: INSERT idempotency_keys<br/>(key=orderId, endpoint, request_hash,<br/>  response_envelope, http_status=200,<br/>  expires_at=NOW()+24h)
      DB-->>-INV: ok
      INV->>+DB: COMMIT
      DB-->>-INV: ok
      INV-->>-CO: 200 {SUCCESS, reservationId, items, expiresAt}
    end
  end
```

## Test cases

### Happy paths

- `reservation_create_single_sku_decrements_available_increments_reserved`
- `reservation_create_multi_sku_atomic_all_or_nothing`
- `reservation_create_idempotent_replay_returns_persisted_envelope`
- `reservation_create_persists_one_idempotency_row_per_orderId`

### Validation

- `reservation_create_idempotencyKey_must_equal_orderId_returns_VALIDATION_ERROR`
- `reservation_create_empty_items_returns_VALIDATION_ERROR`
- `reservation_create_qty_zero_returns_VALIDATION_ERROR`
- `reservation_create_51_items_returns_VALIDATION_ERROR`

### Stock failures

- `reservation_create_one_sku_insufficient_returns_409_with_details_no_partial_decrement`
- `reservation_create_two_skus_one_insufficient_rolls_back_both`

### Auth

- `reservation_create_x_internal_secret_missing_returns_AUTH_INVALID`
- `reservation_create_x_internal_secret_constant_time_compare_audit_log_no_secret_value`

### Concurrency (the load-bearing tests)

- `reservation_create_no_negative_two_goroutines_compete_for_last_unit_one_succeeds_one_INSUFFICIENT_STOCK_never_negative` (BA edge-case)
- `reservation_create_multi_sku_lock_order_two_concurrent_overlapping_sku_sets_in_OPPOSITE_request_order_no_deadlock` (ADR-005 acid test)
- `reservation_create_concurrent_same_orderId_serializes_via_idempotency_FOR_UPDATE_one_succeeds_other_replays`

### Defensive

- `reservation_create_PG_check_violation_maps_to_STOCK_NEGATIVE_INVARIANT_envelope`
- `reservation_create_envelope_carries_traceId_on_every_branch`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Opus 4.7 1M) | Created with mandatory multi-SKU lex-lock sequence diagram |
