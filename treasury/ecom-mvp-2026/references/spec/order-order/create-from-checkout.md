# POST /api/v1/order/internal/create-from-checkout

## Summary

Internal endpoint (Checkout-only). Creates the `orders` row in `PENDING_PAYMENT` plus the immutable `order_items` snapshot block, the initial `order_status_history` row, and an audit-only `outbox_events` row — ALL in a SINGLE transaction (atomic). Allocates the human-readable `order_number` via `order_number_seq` (ORD-008). Validates `Idempotency-Key == orderId` per `cross-cutting.idempotency.internal_orderId_rule`. Replay with the same body returns the original envelope (de-facto idempotency via `orders.id` PK). Replay with a different body returns 409 `IDEMPOTENCY_KEY_REUSED`.

This endpoint exists because Order is the SOLE WRITER of `orders.status` (ADR-003). Checkout cannot INSERT directly into `order.orders`; it must go through this handler — which is also why path (a) of `chosen_order_creation_path` was selected over event-driven creation.

## Story refs

- `STORY_CHECKOUT_COMMIT` (Checkout's caller side)
- `STORY_ORDER_LIST_DETAIL` (the order row this endpoint creates is what list-mine + detail subsequently read)

## Contract ref

[`order.create-from-checkout`](../../architecture/contracts.json#order.create-from-checkout)

## Auth

- Tier: **internal_secret**.
- Header: `X-Internal-Secret: <secret>`.
- Validation: `crypto/subtle.ConstantTimeCompare([]byte(supplied), []byte(secret)) == 1` against `INTERNAL_SHARED_SECRET` (env var name LOCKED — ADR-008). Comma-separated rotation overlap supported (try each segment; first match wins).
- Plain `==` / `!=` is FORBIDDEN (CWE-208 / OWASP API2:2023).
- On reject: `slog.Error("internal_auth_reject", remoteAddr, x_forwarded_for, path, requestId)` — NEVER log the supplied/expected secret value. Returns 401 `AUTH_INVALID`.
- Frontend MUST NOT load `INTERNAL_SHARED_SECRET` nor forward `X-Internal-Secret` (REV-L2-013).

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["orderId","customerUserId","buyerEmail","addressSnapshot","items","subtotal","shippingFee","couponDiscount","total"],
  "properties": {
    "orderId":         { "type": "string", "format": "uuid", "description": "UUID v7; client-supplied by Checkout; recipient stores as orders.id." },
    "customerUserId":  { "type": "string", "format": "uuid" },
    "buyerEmail":      { "type": "string", "format": "email", "maxLength": 254 },
    "addressSnapshot": {
      "type": "object",
      "required": ["addressId","receiverName","phone","addressLine","province","district","postalCode"],
      "properties": {
        "addressId":    { "type": "string", "format": "uuid" },
        "receiverName": { "type": "string", "minLength": 1, "maxLength": 120 },
        "phone":        { "type": "string", "minLength": 1, "maxLength": 32 },
        "addressLine":  { "type": "string", "minLength": 1 },
        "province":     { "type": "string", "minLength": 1 },
        "district":     { "type": "string", "minLength": 1 },
        "postalCode":   { "type": "string", "minLength": 1, "maxLength": 16 }
      }
    },
    "items": {
      "type": "array",
      "minItems": 1,
      "items": {
        "type": "object",
        "required": ["productId","sku","productNameSnapshot","priceSnapshot","qty"],
        "properties": {
          "productId":               { "type": "string", "format": "uuid" },
          "sku":                     { "type": "string", "minLength": 1 },
          "productNameSnapshot":     { "type": "string", "minLength": 1 },
          "productImageUrlSnapshot": { "type": ["string","null"] },
          "priceSnapshot":           { "type": "integer", "minimum": 0 },
          "qty":                     { "type": "integer", "minimum": 1 }
        }
      }
    },
    "subtotal":       { "type": "integer", "minimum": 0 },
    "shippingFee":    { "type": "integer", "minimum": 0 },
    "couponDiscount": { "type": "integer", "const": 0 },
    "total":          { "type": "integer", "minimum": 0 }
  }
}
```

Required headers:
- `X-Internal-Secret: <secret>` — validated via constant-time compare (ADR-008).
- `Idempotency-Key: <orderId UUID v7>` — MUST equal `request.orderId` (cross-cutting.idempotency.internal_orderId_rule).

## Response (success)

HTTP 201, envelope:

```json
{
  "code": "CREATED",
  "message": "order created",
  "data": {
    "orderId": "01935b9c-...-uuidv7",
    "orderNumber": "ORD-20260508-000123",
    "status": "PENDING_PAYMENT",
    "createdAt": "2026-05-08T07:12:33Z"
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `VALIDATION_ERROR` | 400 | `buyerEmail` missing/malformed; `Idempotency-Key != orderId` (REV-L2-006); `subtotal != Σ priceSnapshot*qty`; `total != subtotal + shippingFee - couponDiscount`; required field missing |
| `AUTH_INVALID` | 401 | missing or wrong `X-Internal-Secret` (constant-time compare; ADR-008) |
| `IDEMPOTENCY_KEY_REUSED` | 409 | same `orderId` + DIFFERENT canonicalized body — error data carries the ORIGINAL `orderNumber` so caller can recover |
| `DATABASE_UNAVAILABLE` | 503 | Postgres unreachable |

## Business logic steps

1. **Internal auth middleware:** constant-time-compare `X-Internal-Secret` against `INTERNAL_SHARED_SECRET` (or each comma-segment for rotation overlap). Reject → 401 `AUTH_INVALID`; log `internal_auth_reject`. (ADR-008)
2. **Idempotency-Key validation:** assert header `Idempotency-Key == request.orderId`; mismatch → 400 `VALIDATION_ERROR` (REV-L2-006).
3. **Body validation:** schema bind; `buyerEmail` RFC5322; defensive recompute `Σ priceSnapshot*qty == subtotal` and `total == subtotal + shippingFee - couponDiscount`; mismatch → 400 `VALIDATION_ERROR`. `couponDiscount` MUST equal 0 in MVP.
4. **Single transaction:**
   1. `BEGIN;`
   2. Try `INSERT INTO orders (id=$orderId, user_id=$customerUserId, status='PENDING_PAYMENT', subtotal, shipping_fee, coupon_discount, grand_total=$total, currency='THB', address_snapshot=$jsonb, buyer_email_snapshot=$buyerEmail, idempotency_key=$orderId, version=0) ON CONFLICT (id) DO NOTHING RETURNING id;`  *// ADR-003 sole-writer (initial INSERT)*
   3. **Idempotent replay branch:** if no row returned (duplicate `orders.id`):
      - `SELECT * FROM orders WHERE id=$1; SELECT * FROM order_items WHERE order_id=$1`.
      - Compute `request_hash = sha256(canonicalize(request))` and compare to a stored hash (we use `idempotency_key` column as the dedup discriminator — same orderId means same intent, but if items differ we still detect by row content).
      - If body matches stored row content → COMMIT; return original `{orderId, orderNumber, status:'PENDING_PAYMENT', createdAt}` envelope as 200 `SUCCESS` (per cross-cutting.idempotency: duplicate same payload returns persisted envelope).
      - If body differs → ROLLBACK; return 409 `IDEMPOTENCY_KEY_REUSED` with `data.originalOrderNumber` populated.
   4. **Fresh-create branch:** allocate `seq = nextval('order_number_seq')`. Format `orderNumber = fmt.Sprintf("ORD-%s-%06d", time.Now().UTC().Format("20060102"), seq)` (ORD-008). `UPDATE orders SET order_number=$1 WHERE id=$2;`
   5. `INSERT INTO order_items (one row per request.items[i]; price_snapshot/name_snapshot/image_url_snapshot frozen — ORD-007); line_subtotal = priceSnapshot * qty.` Batched insert.
   6. `INSERT INTO order_status_history (order_id, from_status=NULL, to_status='PENDING_PAYMENT', actor_user_id=NULL, actor_role='SYSTEM', reason='checkout.commit step 7', occurred_at=NOW());`  *// ORD-009 initial row*
   7. `INSERT INTO outbox_events (id=uuidv7(), aggregate_id=order_id, event_type='order.created', payload_json={...audit-only PENDING_PAYMENT signal...}, created_at=NOW());`  *// per contracts.json line 1067 — audit-only; no MVP consumer*
   8. `COMMIT;`
5. Return 201 envelope.

## Side effects

- `INSERT orders` (1 row, status `PENDING_PAYMENT`).
- `INSERT order_items` (1+ rows, immutable snapshots — ORD-007).
- `INSERT order_status_history` (1 row, initial — ORD-009).
- `INSERT outbox_events` (1 row, `order.created` audit-only — no MVP consumer; future analytics).
- `SELECT nextval('order_number_seq')` (sequence side effect).

## Idempotency

- **Key:** `Idempotency-Key` header (which MUST equal `orderId` per the internal_orderId_rule).
- **Scope:** per-orderId; the recipient uses `orders.id` PK as the de-facto dedup mechanism — `INSERT ... ON CONFLICT (id) DO NOTHING` short-circuits replays.
- **TTL:** indefinite (the `orders` row lives forever).
- **First call:** processes normally; commits the trio + audit outbox.
- **Replay with same body:** returns the original `{orderId, orderNumber, status, createdAt}` envelope verbatim (HTTP 200 — note the contract says 201 on first create; replays are 200 because no new row was created).
- **Replay with different body:** 409 `IDEMPOTENCY_KEY_REUSED` carrying the original orderNumber so the caller (Checkout) can surface a clean error to the customer.
- **Mismatched header:** `Idempotency-Key != orderId` → 400 `VALIDATION_ERROR` (the caller violated `internal_orderId_rule`).

## Performance

- p95 < 250ms (cross-component PERF-002 checkout < 500ms budget; this endpoint is one of three downstream calls, so its slice is ~150-250ms).
- Single transaction with 1 INSERT + 1 sequence call + 1 UPDATE (order_number) + N INSERTs (items) + 1 INSERT (history) + 1 INSERT (outbox) = O(N) where N = items count, typically ≤ 10.

## Sequence diagram (mandatory)

```mermaid
sequenceDiagram
  autonumber
  participant CO as Checkout
  participant ORD as Order Service
  participant MW  as internal_auth middleware
  participant DB  as Postgres (order schema)
  participant SEQ as order_number_seq
  participant OB  as outbox_events
  participant PUB as outbox publisher (goroutine)
  participant K   as Kafka (ecom.order.events)

  CO->>+ORD: POST /api/v1/order/internal/create-from-checkout<br/>Headers: X-Internal-Secret, Idempotency-Key=<orderId><br/>Body: {orderId, customerUserId, buyerEmail, addressSnapshot, items[], subtotal, shippingFee, couponDiscount=0, total}
  ORD->>+MW: subtle.ConstantTimeCompare(supplied, INTERNAL_SHARED_SECRET) [ADR-008]
  alt mismatch
    MW-->>CO: 401 AUTH_INVALID (slog internal_auth_reject; never logs secret)
  else match
    MW-->>-ORD: ok
    ORD->>ORD: assert Idempotency-Key == orderId<br/>(REV-L2-006; mismatch -> 400 VALIDATION_ERROR)
    ORD->>ORD: validate body, recompute subtotal + total
    ORD->>+DB: BEGIN
    ORD->>DB: INSERT orders (id=orderId, user_id, status='PENDING_PAYMENT',<br/>subtotal, shipping_fee, coupon_discount=0, grand_total=total,<br/>address_snapshot, buyer_email_snapshot, idempotency_key=orderId, version=0)<br/>ON CONFLICT (id) DO NOTHING RETURNING id<br/>// ADR-003 sole-writer initial INSERT
    alt duplicate orderId (idempotent replay)
      DB-->>ORD: no row
      ORD->>DB: SELECT orders + order_items WHERE id=orderId
      DB-->>ORD: existing row + items
      alt body matches stored
        ORD->>DB: COMMIT
        ORD-->>CO: 200 {orderId, orderNumber, status:'PENDING_PAYMENT', createdAt}
      else body differs
        ORD->>DB: ROLLBACK
        ORD-->>CO: 409 IDEMPOTENCY_KEY_REUSED {originalOrderNumber}
      end
    else fresh INSERT
      DB-->>ORD: new id
      ORD->>+SEQ: SELECT nextval('order_number_seq')
      SEQ-->>-ORD: seq=N
      ORD->>DB: UPDATE orders SET order_number = 'ORD-YYYYMMDD-NNNNNN' (ORD-008)
      ORD->>DB: INSERT order_items (1+ rows, price/name/image snapshots — ORD-007)
      ORD->>DB: INSERT order_status_history (from=NULL, to='PENDING_PAYMENT',<br/>actor_role='SYSTEM', reason='checkout.commit step 7') — ORD-009
      ORD->>+OB: INSERT outbox_events (event_type='order.created', payload audit-only)
      OB-->>-ORD: ok
      ORD->>DB: COMMIT  // atomic: orders + items + history + outbox
      DB-->>-ORD: tx committed
      ORD-->>-CO: 201 CREATED {orderId, orderNumber, status:'PENDING_PAYMENT', createdAt}
    end
  end

  note over PUB,K: outbox publisher goroutine (ADR-004)<br/>polls every 1s; ORDER BY id ASC LIMIT 100
  PUB->>+OB: SELECT id, payload_json FROM outbox_events WHERE published_at IS NULL ORDER BY id LIMIT 100
  OB-->>-PUB: rows including order.created audit row
  PUB->>+K: produce(topic='ecom.order.events', key=orderId, value=payload)
  K-->>-PUB: ack
  PUB->>OB: UPDATE outbox_events SET published_at=NOW() WHERE id=$id
```

## Test cases

- `create_from_checkout_writes_orders_items_history_and_audit_outbox_atomically`
- `create_from_checkout_allocates_order_number_in_format_ORD_YYYYMMDD_NNNNNN`
- `create_from_checkout_initial_status_history_row_actor_role_SYSTEM_from_status_NULL`
- `create_from_checkout_idempotent_replay_with_same_body_returns_original_orderNumber`
- `create_from_checkout_replay_with_different_body_returns_409_IDEMPOTENCY_KEY_REUSED`
- `create_from_checkout_idempotency_key_neq_orderId_returns_400_VALIDATION_ERROR`
- `create_from_checkout_buyerEmail_missing_returns_400_VALIDATION_ERROR`
- `create_from_checkout_total_inconsistent_returns_400_VALIDATION_ERROR`
- `create_from_checkout_couponDiscount_nonzero_returns_400_VALIDATION_ERROR`
- `create_from_checkout_missing_X_Internal_Secret_returns_401_AUTH_INVALID`
- `create_from_checkout_wrong_X_Internal_Secret_returns_401_constant_time_compare`
- `create_from_checkout_constant_time_compare_no_secret_in_logs` (security smoke)
- `create_from_checkout_rotation_overlap_accepts_either_segment`
- `create_from_checkout_rollback_on_items_insert_failure_leaves_no_partial_state`
- `create_from_checkout_outbox_audit_row_present_for_observability`
- `create_from_checkout_with_admin_jwt_returns_401_AUTH_INVALID` (admin tokens are not valid here; only X-Internal-Secret)
- `create_from_checkout_with_customer_jwt_returns_401_AUTH_INVALID` (customer must NEVER reach this endpoint — REV-L2-013 frontend posture)

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Opus 4.7 1M, dry-run #2) | Created with mandatory mermaid sequence diagram; aligns with cross-cutting.idempotency.internal_orderId_rule, ADR-003, ADR-004, ADR-007, ADR-008. |
