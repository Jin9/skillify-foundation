# POST /api/v1/checkout/checkout/commit

## Summary

Customer-facing **orchestrator** that converts the caller's cart into a `PENDING_PAYMENT` order, a 15-minute stock reservation, and a `REQUIRES_PAYMENT` payment intent — atomically from the customer's point of view, via 9 sync HTTP steps wrapped by a single Postgres `idempotency_keys` row in the checkout schema. Required `Idempotency-Key` header (CHK-009) makes duplicate POSTs return the originally created order; the four downstream internal hops (steps 6/7/8 + the two compensation calls) all use `X-Internal-Secret` and are idempotent on the checkout-minted `orderId` per `cross-cutting.idempotency.internal_orderId_rule`. Pricing is server-authoritative (CHK-005) and the §11.2 shipping fee comes from the SINGLE source `pricing.go`. The compensation matrix has 3 retries with `[50, 200, 800]ms` exponential backoff per call; on retry-exhaustion the saga row goes `DEGRADED` and the inventory reservation TTL (15 min) is the safety net (ADR-007).

This is the load-bearing endpoint of the entire purchase. Every line of the orchestration exists for a specific reason — see `## Business logic steps`, `## Sequence diagram`, and the compensation matrix in `td.json`. ADR-007 makes the trade-off explicit: sync orchestration with explicit compensations beats a saga choreography for the MVP latency budget.

## Story refs

- `STORY_CHECKOUT_COMMIT`
- `STORY_CHECKOUT_SHIPPING_FEE` (commit re-applies §11.2 from `pricing.go` at step 5; preview is advisory-only)

## Contract ref

[`checkout.commit`](../../architecture/contracts.json#checkout.commit)

## Auth

- Tier: **customer JWT (Bearer)**.
- Header: `Authorization: Bearer <access JWT>`. Validated via `common/middleware/jwt_middleware`; `claims.role` MUST be `CUSTOMER`.
- `AUTH_MISSING` (401) if header absent; `AUTH_FORBIDDEN` (403) if `role != CUSTOMER`.
- The customer JWT is forwarded to four customer-context hops only:
  - Step 1 `cart.read`
  - Step 2 `identity.profile.read` ← REV-L2-001 buyerEmail wire
  - Step 3 `identity.address.list`
  - Step 9 `cart.clear-on-checkout`
- The four internal hops + compensations use `X-Internal-Secret` ONLY (the customer JWT is NOT forwarded into Inventory/Order/Payment):
  - Step 6 `inventory.reservation.create`
  - Step 7 `order.create-from-checkout`
  - Step 8 `payment.intent.create`
  - Compensation: `inventory.reservation.release`
  - Compensation: `order.cancel-on-checkout-failure`
- `INTERNAL_SHARED_SECRET` is loaded from env at boot (ADR-008 LOCKED env-var name); MIN 32 bytes; refuse to start if shorter or absent. Comma-separated rotation overlap supported.

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["shippingAddressId"],
  "properties": {
    "shippingAddressId": {
      "type": "string",
      "format": "uuid",
      "minLength": 1,
      "maxLength": 64,
      "description": "UUID v7; MUST belong to the caller (caller-scoped via identity.address.list)."
    }
  }
}
```

Validation rules beyond the schema (enforced before `TryClaimOrLookup`):
- `Idempotency-Key` header REQUIRED, regex `^[A-Za-z0-9._:-]{1,128}$`. Missing or malformed → 400 `VALIDATION_ERROR` with message `"Idempotency-Key header required"`.
- If body contains ANY of `{total, subtotal, shippingFee, couponDiscount, grandTotal, currentPrice, priceSnapshot}` → 400 `VALIDATION_ERROR` with message `"client-supplied price/total fields are not accepted; server is authoritative (CHK-005)"`.
- If body contains `couponCode` → 400 `VALIDATION_ERROR` with message `"couponCode is not supported in MVP"` (CHK-AMBIG-002).
- The schema's `additionalProperties: false` rejects any unknown field at bind time as 400 `VALIDATION_ERROR`.

Required headers:
- `Authorization: Bearer <access JWT>`
- `Content-Type: application/json`
- `Idempotency-Key: <opaque string>` — UUID v4 recommended; max 128 chars; per-(endpoint, customerUserId) scope; 24-hour TTL.

## Response (success)

HTTP 201, envelope:

```json
{
  "code": "CREATED",
  "message": "order created",
  "data": {
    "orderId": "01935b9c-1a17-7c3d-9e8f-aaaaaaaaaaaa",
    "orderNumber": "ORD-20260508-000123",
    "status": "PENDING_PAYMENT",
    "subtotal": 1198,
    "shippingFee": 60,
    "couponDiscount": 0,
    "total": 1258,
    "paymentIntent": {
      "paymentIntentId": "01935b9c-2a17-7c3d-9e8f-bbbbbbbbbbbb",
      "status": "REQUIRES_PAYMENT",
      "amount": 1258
    }
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

Notes on the success body:
- `orderId` is a UUID v7 minted by Checkout at step 5 — BEFORE any internal hop — so it can double as the `Idempotency-Key` on steps 6/7/8 per the LOCKED `cross-cutting.idempotency.internal_orderId_rule`.
- `orderNumber` is allocated by Order via `order_number_seq` (ORD-008) inside `order.create-from-checkout`'s tx and returned to Checkout in step 7's response.
- `subtotal/shippingFee/total` are computed by `pricing.go` (`ComputeShippingFee` + `ComputeGrandTotal`) — the SINGLE source for §11.2. `couponDiscount` is always 0 in MVP.
- `paymentIntent.amount` MUST equal `total` by construction (step 8 sends `amount=total`).

Replay behavior (per `cross-cutting.idempotency`):
- Same `Idempotency-Key` + same canonicalized body → cached envelope returned verbatim with the cached HTTP status (201 on first success, or whichever decided-error status the original commit produced).
- Same `Idempotency-Key` + DIFFERENT canonicalized body → 409 `IDEMPOTENCY_KEY_REUSED` carrying `data.originalOrderId` from the cached envelope.
- Same `Idempotency-Key` while first call still INFLIGHT → 409 `IDEMPOTENCY_KEY_INFLIGHT` with `Retry-After: 1`.

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `VALIDATION_ERROR` | 400 | empty cart at step 1 (CHK-001); missing/invalid `Idempotency-Key`; client-priced fields present (CHK-005); `couponCode` present (CHK-AMBIG-002); `shippingAddressId` empty/missing; unknown body fields |
| `AUTH_MISSING` | 401 | no `Authorization` header |
| `AUTH_FORBIDDEN` | 403 | `claims.role != CUSTOMER` |
| `ADDRESS_NOT_OWNED` | 403 | `shippingAddressId` not in caller's `identity.address.list` (CHK-002) |
| `ADDRESS_INCOMPLETE` | 400 | address row exists but is missing required fields |
| `PRODUCT_INACTIVE` | 409 | one or more cart items are INACTIVE in catalog (CHK-003); response carries `data.blockers[]` |
| `PRODUCT_DELETED` | 409 | one or more cart items are DELETED (CHK-003); response carries `data.blockers[]` |
| `INSUFFICIENT_STOCK` | 409 | pre-reservation (step 4) OR reservation-race (step 6) detected; response carries `data` with `[{productId, requestedQty, availableQty}]` per item (CHK-004) |
| `IDEMPOTENCY_KEY_REUSED` | 409 | same key + different `request_hash`; `data.originalOrderId` from cached envelope (CHK-009) |
| `IDEMPOTENCY_KEY_INFLIGHT` | 409 | concurrent same-key call still INFLIGHT; `Retry-After: 1` (CHK-009) |
| `DATABASE_UNAVAILABLE` | 503 | checkout DB unreachable, OR a downstream's DB is unreachable and its envelope propagates |
| `UPSTREAM_TIMEOUT` | 504 | a downstream (cart, identity, catalog, inventory, order, payment, or any compensation call) exceeded its budget after its retry — including the `ABANDONED` replay envelope written by the stale-INFLIGHT janitor |

## Business logic steps

The full pseudocode lives in `td.json#endpoints[1].orchestration_pseudocode`. The narrative below mirrors the contracts-locked 9-step shape (with a finalize step 10 that updates `idempotency_keys` to COMPLETED and inserts the `saga_log` row).

**Pre-orchestration (no tx, fast-fail):**

1. **Auth + role gate** via `jwt_middleware`; reject non-CUSTOMER roles.
2. **Idempotency-Key header** validated against `^[A-Za-z0-9._:-]{1,128}$`.
3. **Body validation:** strict bind (`additionalProperties: false`); reject `couponCode` (CHK-AMBIG-002) and reject any of `{total, subtotal, shippingFee, couponDiscount, grandTotal, currentPrice, priceSnapshot}` (CHK-005); require `shippingAddressId`.
4. **Compute `requestHash` = sha256(canonicalize(body))** for the idempotency comparison.

**Step 0 — open the outer tx + atomic idempotency claim:**

5. `BEGIN TX` (ReadCommitted) on the checkout DB.
6. `idempotencyStore.TryClaimOrLookup(tx, idemKey, claims.sub, requestHash)` runs the SINGLE statement `INSERT INTO checkout.idempotency_keys (...) VALUES (...,'INFLIGHT',...) ON CONFLICT (key, customer_user_id) DO NOTHING RETURNING key` — closes the dry-run #1 H02 race window. Branch on result:
   - **winner=true** → proceed (we own the INFLIGHT row).
   - **COMPLETED + hash matches** → ROLLBACK; return cached `response_envelope` verbatim.
   - **COMPLETED + hash mismatch** → ROLLBACK; return 409 `IDEMPOTENCY_KEY_REUSED` with `data.originalOrderId` from the cached envelope.
   - **INFLIGHT** → ROLLBACK; return 409 `IDEMPOTENCY_KEY_INFLIGHT` with `Retry-After: 1`.
   - **ABANDONED** (janitor flipped a stale INFLIGHT > 60s) → DELETE the dead row inside the same tx (we hold the FOR UPDATE lock); re-INSERT INFLIGHT; fall through as winner.

**Orchestration body (tx open; downstream HTTP calls outside tx — checkout's row lock pins INFLIGHT):**

7. **Step 1 — `cart.read`** (Bearer; 800ms; 1 retry). Empty cart → `finalize(VALIDATION_ERROR 'cart is empty (CHK-001)' httpStatus=400)`; COMMIT; return.
8. **Step 2 — `identity.profile.read`** (Bearer; 500ms; 1 retry). REQUIRED — feeds `buyerEmail` to step 7 (REV-L2-001 closure). Failure → `finalize(UPSTREAM_TIMEOUT httpStatus=504)`; COMMIT; return. The customer's own access token is forwarded; the email is sourced server-side and is NEVER read from the request body.
9. **Step 3 — `identity.address.list`** (Bearer; 500ms; 1 retry). Find `req.shippingAddressId`; absent → `finalize(ADDRESS_NOT_OWNED httpStatus=403 CHK-002)`. Incomplete fields → `finalize(ADDRESS_INCOMPLETE httpStatus=400)`. Snapshot the matched address.
10. **Step 4 — catalog + inventory pre-flight:**
    - **`catalog.product.detail` × N (parallel via errgroup; cap 50; 500ms each; 1 retry)** — collect status/price/name/imageUrl. Any line with `status IN {INACTIVE, DELETED}` or `error == NOT_FOUND` becomes a blocker; `len(blockers) > 0` → `finalize(409 PRODUCT_INACTIVE/PRODUCT_DELETED + data.blockers httpStatus=409 CHK-003)`.
    - **`inventory.stock.bulk-read` (X-Internal-Secret; 800ms; 1 retry)** — map `{sku → availableQty}`. Lines where `qty > availableQty` become shortages; `len(shortages) > 0` → `finalize(409 INSUFFICIENT_STOCK + data httpStatus=409 CHK-004)`.
11. **Step 5 — local pricing snapshot + orderId mint:**
    - `orderId = uuidv7()` — minted NOW. This UUID will be the `Idempotency-Key` for steps 6/7/8 per `cross-cutting.idempotency.internal_orderId_rule`.
    - Build `lineItems` with `priceSnapshot = catalog.detail.price` (frozen), `productNameSnapshot`, `productImageUrlSnapshot`, `qty`, `cartItemId`.
    - `subtotalMinor = Σ priceSnapshot * qty` (int64 minor units).
    - `shippingFeeMinor = pricing.ComputeShippingFee(subtotalMinor)` ← **SINGLE SOURCE** — `pricing.go`. Boundary inclusive at 1500 (subtotal=1499→60; subtotal=1500→0; subtotal=9999→0). BA-PR-002.
    - `couponDiscountMinor = 0`.
    - `grandTotalMinor = pricing.ComputeGrandTotal(subtotalMinor, shippingFeeMinor, 0)`.
12. **Step 6 — `inventory.reservation.create`** (X-Internal-Secret + `Idempotency-Key=orderId`; 1.5s; 1 retry). Body: `{orderId, items, expiresAt: now+15min, idempotencyKey: orderId}`. Race → `INSUFFICIENT_STOCK` → `finalize(409 CHK-004)`; other failure → `finalize(UPSTREAM_TIMEOUT or DATABASE_UNAVAILABLE)`. **Compensation on failure here: NONE** (no remote state changed).
13. **Step 7 — `order.create-from-checkout`** (X-Internal-Secret + `Idempotency-Key=orderId`; 1s; 1 retry). Body includes `buyerEmail = profile.email` from step 2 (REV-L2-001), the `addressSnapshot`, `lineItems`, and the int-THB pricing fields. Failure → **COMPENSATION**: 3 retries of `inventory.reservation.release(orderId, reason=ORDER_CANCELLED)` with `[50, 200, 800]ms` exponential backoff; on exhaust mark `saga_log.status=DEGRADED` + log CRITICAL. Either way, `finalize(UPSTREAM_TIMEOUT)`.
14. **Step 8 — `payment.intent.create`** (X-Internal-Secret + `Idempotency-Key=orderId`; 1.5s; 1 retry). Body: `{orderId, amount: grandTotalMinor, ownerUserId: claims.sub}`. Failure → **COMPENSATION** in two phases:
    - **(A)** 3 retries of `order.cancel-on-checkout-failure(orderId, reason='PAYMENT_INTENT_FAILURE')` with `[50, 200, 800]ms` backoff.
    - **(B)** 3 retries of `inventory.reservation.release(orderId, reason=ORDER_CANCELLED)` with `[50, 200, 800]ms` backoff.
    - If either phase exhausts → `saga_log.status=DEGRADED` + log CRITICAL; the reservation TTL (15 min) is the safety net, and `events.reservation.expired` will drive `PAYMENT_EXPIRED` via Order's state-driven consumer. `finalize(UPSTREAM_TIMEOUT)`.
15. **Step 9 — `cart.clear-on-checkout`** (Bearer; 800ms; 2 retries `[50, 250]ms`). Best-effort; **DO NOT compensate the order**. Failure → log WARN + emit `checkout_cart_clear_failed_total`; PROCEED to step 10. Order is `PENDING_PAYMENT` and the customer has the orderId; cart inconsistency is cosmetic and self-healing on the next `cart.read`.
16. **Step 10 — finalize:**
    - `idempotencyStore.Finalize(tx, idemKey, claims.sub, COMPLETED, envelope, httpStatus=201)` — write the cached envelope so retries with the same key + same hash become a verbatim replay.
    - `sagaLogStore.Write(tx, orderId, claims.sub, traceId, COMPLETED, steps_jsonb, compensation=null)`.
    - `COMMIT TX`.
    - Return HTTP 201 with the success envelope.

## Side effects

- **Always (per request):** one row in `checkout.idempotency_keys` (status `INFLIGHT` → `COMPLETED` or `ABANDONED`; the row carries the final response envelope so retries are deterministic).
- **Always (per request):** one row in `checkout.saga_log` written at step 10, regardless of outcome (`COMPLETED` on success or decided client error; `COMPENSATED` if compensation succeeded; `DEGRADED` if compensation exhausted).
- **On success path:**
  - One row in `inventory.reservations` (status `RESERVED`, expires_at = `created_at + 15min`).
  - One row in `order.orders` (status `PENDING_PAYMENT`), N rows in `order.order_items`, one row in `order.order_status_history` (initial), one audit row in `order.outbox_events`.
  - One row in `payment.payment_intents` (status `REQUIRES_PAYMENT`).
  - Cart items moved into the order are removed from `cart.cart_items` (CHK-010); other cart items remain.
- **On step-7 failure path:** reservation row transitions `RESERVED → RELEASED` with `release_reason='ORDER_CANCELLED'` (via the compensation call). Order row never written.
- **On step-8 failure path:** order row transitions `PENDING_PAYMENT → CANCELLED` (via `order.cancel-on-checkout-failure`); reservation row transitions to `RELEASED`. Payment intent never written.
- **On step-9 failure path:** all the success-path effects, MINUS the cart-clear (the cart still contains the just-ordered items; future `cart.read` re-reads `availableQty` and the customer can re-clear). Order, reservation, and payment intent are committed.
- **No event emission from Checkout itself.** Order owns `events.order.cancelled`; Payment owns `events.payment.*`; Inventory owns `events.reservation.expired`. Checkout is a sync orchestrator only (ADR-007).

## Idempotency

- **Key shape:** customer-supplied opaque string in the `Idempotency-Key` HTTP header; max 128 chars; `^[A-Za-z0-9._:-]{1,128}$`; UUID v4 recommended.
- **Scope:** per-(endpoint=`checkout.commit`, `customer_user_id=claims.sub`). The `(key, customer_user_id)` composite PK on `idempotency_keys` enforces this. A single key cannot collide across two different customers.
- **TTL:** 24 hours per `cross-cutting.idempotency.ttl_hours` for `checkout.commit`. The TTL janitor deletes COMPLETED/ABANDONED rows past their `expires_at` in batches of 500 every 5 minutes.
- **Atomic claim:** `INSERT … ON CONFLICT (key, customer_user_id) DO NOTHING RETURNING key` — the SINGLE-statement TryClaimOrLookup pattern. Closes the dry-run #1 H02 race window where a Lookup→InsertInflight pair could allow two concurrent first-time callers to both proceed (two orders for one Idempotency-Key).
- **Replay branches:** see `## Response (success)` notes.
- **Internal hop substitution:** Checkout substitutes `Idempotency-Key=orderId` at steps 6/7/8 + compensations. The customer-facing Idempotency-Key on this commit endpoint is NEVER forwarded downstream — it is a separate scope confined to checkout's local DB.
- **Stale-INFLIGHT janitor:** every 60s, the janitor flips rows with `status='INFLIGHT' AND created_at < NOW() - INTERVAL '60 seconds'` to `ABANDONED` with a cached 504 envelope (the original orchestration crashed mid-flight; replays surface that 504 instead of starting fresh).

## Performance

- p95 < 500 ms (PERF-002 in local Docker Compose).
- Expected QPS at MVP: 1–5.
- Latency budget breakdown (success path, no retries):
  - Pre-tx validation + hash: < 5 ms.
  - Step 0 TryClaimOrLookup (single INSERT): 2–5 ms.
  - Step 1 cart.read: 800ms cap; typical 30–80 ms.
  - Step 2 identity.profile.read: 500ms cap; typical 20–50 ms.
  - Step 3 identity.address.list: 500ms cap; typical 20–50 ms.
  - Step 4 catalog × N parallel + inventory: 800ms cap (parallel wall-clock); typical 50–150 ms.
  - Step 5 local pricing: < 1 ms.
  - Step 6 inventory.reservation.create: 1.5s cap; typical 30–100 ms (FOR UPDATE on stock_levels).
  - Step 7 order.create-from-checkout: 1s cap; typical 30–80 ms.
  - Step 8 payment.intent.create: 1.5s cap; typical 30–80 ms.
  - Step 9 cart.clear-on-checkout: 800ms cap; typical 30–80 ms.
  - Step 10 finalize + COMMIT: < 5 ms.
  - **Total typical p95: 220–470 ms** (within budget). Worst-case-with-retries may approach the 5s global budget; the saga reconciler surfaces those rows.

## Sequence diagram (mandatory — orchestrator endpoint)

```mermaid
sequenceDiagram
  autonumber
  participant FE as Frontend (Next.js route handler)
  participant CO as Checkout
  participant DB as Postgres (checkout schema)
  participant CART as Cart
  participant ID as Identity
  participant CAT as Catalog
  participant INV as Inventory
  participant ORD as Order
  participant PAY as Payment

  FE->>+CO: POST /api/v1/checkout/checkout/commit<br/>Bearer + Idempotency-Key + {shippingAddressId}
  CO->>CO: validate header + body<br/>(reject couponCode CHK-AMBIG-002,<br/>client price fields CHK-005,<br/>missing Idempotency-Key)
  CO->>CO: requestHash = sha256(canonicalize(body))

  CO->>+DB: BEGIN TX (ReadCommitted)
  CO->>DB: TryClaimOrLookup<br/>INSERT idempotency_keys (key, customer_user_id, request_hash, status='INFLIGHT')<br/>ON CONFLICT (key, customer_user_id) DO NOTHING RETURNING key
  alt winner=false (existing row)
    DB-->>CO: existing row (SELECT FOR UPDATE)
    alt COMPLETED + hashMatch
      CO-->>FE: cached envelope verbatim (cached HTTP status)
    else COMPLETED + hashMismatch
      CO-->>FE: 409 IDEMPOTENCY_KEY_REUSED { originalOrderId }
    else INFLIGHT
      CO-->>FE: 409 IDEMPOTENCY_KEY_INFLIGHT (Retry-After: 1)
    else ABANDONED
      CO->>DB: DELETE the dead row + re-INSERT INFLIGHT (we hold the FOR UPDATE lock)
    end
  else winner=true
    note over CO,DB: own the INFLIGHT row; tx stays open through step 10
  end

  rect rgba(220,240,255,0.6)
    note over CO,PAY: 9-step orchestration body — TX open; row lock pins INFLIGHT
    CO->>+CART: 1) POST /cart/cart/read (Bearer)
    CART-->>-CO: items[]
    alt items == []
      CO->>DB: Finalize idempotency=COMPLETED + saga=COMPLETED<br/>envelope = 400 VALIDATION_ERROR 'cart is empty' (CHK-001)
      CO->>DB: COMMIT
      CO-->>FE: 400 VALIDATION_ERROR
    end

    CO->>+ID: 2) POST /identity/profile/read (Bearer)<br/>REV-L2-001 buyerEmail wire
    ID-->>-CO: { userId, email, name, ... }

    CO->>+ID: 3) POST /identity/address/list (Bearer)
    ID-->>-CO: addresses[]
    alt shippingAddressId not in addresses
      CO->>DB: Finalize 403 ADDRESS_NOT_OWNED (CHK-002); COMMIT
      CO-->>FE: 403 ADDRESS_NOT_OWNED
    end

    par 4) catalog fan-out (errgroup; cap 50)
      CO->>+CAT: POST /catalog/product/detail (productId)
      CAT-->>-CO: { status, priceMinor, name, imageUrl }
    and 4) inventory bulk-read
      CO->>+INV: POST /inventory/stock/bulk-read<br/>X-Internal-Secret + { skus }
      INV-->>-CO: stocks[{sku, availableQty}]
    end
    alt blockers (INACTIVE/DELETED) OR shortages
      CO->>DB: Finalize 409 PRODUCT_INACTIVE / PRODUCT_DELETED / INSUFFICIENT_STOCK<br/>(CHK-003 / CHK-004); COMMIT
      CO-->>FE: 409 + data.blockers / data.shortages
    end

    CO->>CO: 5) orderId = uuidv7()<br/>shippingFee = pricing.ComputeShippingFee(subtotal)<br/>(SINGLE SOURCE; 1500 inclusive — BA-PR-002)<br/>total = pricing.ComputeGrandTotal(subtotal, shippingFee, 0)

    CO->>+INV: 6) POST /inventory/reservation/create<br/>X-Internal-Secret + Idempotency-Key=orderId<br/>{orderId, items, expiresAt: now+15min}
    alt 6 OK
      INV-->>CO: { reservationId, items: [{sku, qty, reserved: true}], expiresAt }
    else 6 fails (race INSUFFICIENT_STOCK / 5xx / timeout)
      INV-->>-CO: error
      CO->>DB: Finalize 409 INSUFFICIENT_STOCK or 504 UPSTREAM_TIMEOUT<br/>(NO COMPENSATION — no remote state mutated)
      CO->>DB: COMMIT
      CO-->>FE: 409 / 503 / 504
    end

    CO->>+ORD: 7) POST /order/internal/create-from-checkout<br/>X-Internal-Secret + Idempotency-Key=orderId<br/>{orderId, customerUserId, buyerEmail (REV-L2-001), addressSnapshot, items, subtotal, shippingFee, couponDiscount=0, total}
    alt 7 OK
      ORD-->>CO: { orderId, orderNumber, status: PENDING_PAYMENT, createdAt }
    else 7 fails
      ORD-->>-CO: error
      note over CO,INV: COMPENSATION — 3 retries [50, 200, 800] ms
      loop attempts ∈ {50ms, 200ms, 800ms}
        CO->>+INV: POST /inventory/reservation/release<br/>X-Internal-Secret + {orderId, reason: ORDER_CANCELLED}
        INV-->>-CO: { reservationId, status: RELEASED } OR error
      end
      alt compensation OK
        CO->>DB: Finalize UPSTREAM_TIMEOUT + saga=COMPENSATED; COMMIT
      else compensation EXHAUSTED
        CO->>DB: Finalize UPSTREAM_TIMEOUT + saga=DEGRADED + log CRITICAL<br/>(reservation TTL 15min is safety net)
        CO->>DB: COMMIT
      end
      CO-->>FE: 504 UPSTREAM_TIMEOUT
    end

    CO->>+PAY: 8) POST /payment/intent/create<br/>X-Internal-Secret + Idempotency-Key=orderId<br/>{orderId, amount: total, ownerUserId}
    alt 8 OK
      PAY-->>CO: { paymentIntentId, status: REQUIRES_PAYMENT, amount, expiresAt }
    else 8 fails
      PAY-->>-CO: error
      note over CO,ORD,INV: COMPENSATION — (A) cancel order, (B) release reservation; 3 retries each [50, 200, 800] ms
      loop attempts ∈ {50ms, 200ms, 800ms}
        CO->>+ORD: POST /order/internal/cancel-on-checkout-failure<br/>X-Internal-Secret + {orderId, reason: PAYMENT_INTENT_FAILURE}
        ORD-->>-CO: { orderId, status: CANCELLED } OR error
      end
      loop attempts ∈ {50ms, 200ms, 800ms}
        CO->>+INV: POST /inventory/reservation/release<br/>X-Internal-Secret + {orderId, reason: ORDER_CANCELLED}
        INV-->>-CO: { reservationId, status: RELEASED } OR error
      end
      alt both compensations OK
        CO->>DB: Finalize UPSTREAM_TIMEOUT + saga=COMPENSATED; COMMIT
      else either compensation EXHAUSTED
        CO->>DB: Finalize UPSTREAM_TIMEOUT + saga=DEGRADED + log CRITICAL<br/>(reservation TTL + events.reservation.expired path is safety net)
        CO->>DB: COMMIT
      end
      CO-->>FE: 504 UPSTREAM_TIMEOUT
    end

    CO->>+CART: 9) POST /cart/cart/clear-on-checkout (Bearer)<br/>{customerUserId, cartItemIds, orderId}
    alt 9 OK
      CART-->>CO: { removed: N }
    else 9 fails (best-effort; NO compensation)
      CART-->>-CO: error
      CO->>CO: log WARN + emit checkout_cart_clear_failed_total<br/>(order is committed; cart-clear is self-healing)
    end

    CO->>DB: 10) idempotencyStore.Finalize(COMPLETED, envelope, http=201)
    CO->>DB: sagaLogStore.Write(orderId, customerUserId, traceId, COMPLETED, steps_jsonb, compensation=null)
    CO->>DB: COMMIT TX
    DB-->>-CO: tx committed
    CO-->>-FE: 201 CREATED { orderId, orderNumber, status: PENDING_PAYMENT, subtotal, shippingFee, total, paymentIntent }
  end
```

## Test cases

- `commit_single_call_creates_one_order_and_one_reservation_and_one_payment_intent`
- `commit_two_concurrent_same_idempotency_key_returns_one_inflight_409_and_one_201` (TryClaimOrLookup race; closes dry-run #1 H02)
- `commit_two_serial_same_idempotency_key_returns_cached_201_envelope_verbatim`
- `commit_same_key_different_payload_returns_409_IDEMPOTENCY_KEY_REUSED_with_originalOrderId`
- `commit_inflight_replay_returns_409_IDEMPOTENCY_KEY_INFLIGHT_with_Retry_After_1`
- `commit_abandoned_replay_after_60s_janitor_returns_cached_504_or_re_claims_fresh`
- `commit_subtotal_1499_persists_shippingFee_60_total_1559_BA_PR_002`
- `commit_subtotal_1500_persists_shippingFee_0_total_1500_BA_PR_002_inclusive_boundary`
- `commit_uses_pricing_go_single_source` (grep: `ComputeShippingFee` is the only place subtotal-vs-1500 is decided in checkout)
- `commit_client_supplied_total_field_returns_400_VALIDATION_ERROR` (CHK-005)
- `commit_client_supplied_subtotal_field_returns_400_VALIDATION_ERROR` (CHK-005)
- `commit_couponCode_field_returns_400_VALIDATION_ERROR` (CHK-AMBIG-002)
- `commit_missing_idempotency_key_returns_400_VALIDATION_ERROR`
- `commit_malformed_idempotency_key_returns_400_VALIDATION_ERROR`
- `commit_empty_cart_returns_400_VALIDATION_ERROR_at_step_1` (CHK-001)
- `commit_address_not_owned_returns_403_ADDRESS_NOT_OWNED` (CHK-002)
- `commit_inactive_product_returns_409_PRODUCT_INACTIVE_with_blockers` (CHK-003)
- `commit_deleted_product_returns_409_PRODUCT_DELETED_with_blockers` (CHK-003)
- `commit_qty_exceeds_availableQty_returns_409_INSUFFICIENT_STOCK_with_per_item_details` (CHK-004)
- `commit_step6_race_INSUFFICIENT_STOCK_returns_409_no_order_no_payment_intent` (CHK-004)
- `commit_step7_failure_compensates_inventory_release_3_retries_50_200_800_ms` (compensation matrix row 1)
- `commit_step7_compensation_exhaustion_marks_saga_DEGRADED_logs_CRITICAL` (REV-L2-008)
- `commit_step8_failure_compensates_order_cancel_then_inventory_release_each_3_retries` (compensation matrix row 2)
- `commit_step8_compensation_exhaustion_marks_saga_DEGRADED_reservation_TTL_safety_net`
- `commit_step9_failure_does_NOT_compensate_returns_201_with_warning_log` (compensation matrix row 3)
- `commit_buyerEmail_is_sourced_from_identity_profile_read_step2_NOT_request_body` (REV-L2-001)
- `commit_buyerEmail_is_forwarded_to_order_create_from_checkout_step7` (REV-L2-001)
- `commit_internal_hops_use_X_Internal_Secret_NOT_customer_jwt` (steps 6/7/8 + compensations)
- `commit_customer_jwt_NOT_forwarded_to_inventory_order_payment` (REV-L2-013-style negative smoke)
- `commit_idempotency_key_substitution_internal_orderId_rule_step_6_7_8_use_orderId_as_idempotency_key`
- `commit_admin_token_returns_403_AUTH_FORBIDDEN`
- `commit_no_authorization_header_returns_401_AUTH_MISSING`
- `commit_pgxpool_unreachable_returns_503_DATABASE_UNAVAILABLE`
- `commit_inventory_5xx_with_retry_propagates_504_UPSTREAM_TIMEOUT`
- `commit_payment_5xx_triggers_full_compensation_matrix`
- `commit_does_NOT_write_to_orders_or_payment_intents_or_stock_levels` (negative side-effect smoke; ADR-001 schema isolation)
- `commit_writes_saga_log_row_with_steps_jsonb_per_attempt`
- `commit_inflight_janitor_marks_stale_60s_INFLIGHT_to_ABANDONED_with_504_envelope`
- `commit_ttl_janitor_deletes_COMPLETED_rows_after_24h`
- `commit_saga_reconciler_emits_metric_when_DEGRADED_row_older_than_1h` (REV-L2-008)
- `commit_TryClaimOrLookup_atomic_INSERT_ON_CONFLICT_DO_NOTHING_RETURNING_closes_H02_race`
- `commit_p95_under_500ms_local_docker_compose` (PERF-002 smoke)

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Opus 4.7 1M, dry-run #2) | Created with mandatory mermaid sequence diagram for the 9-step orchestrator. Wires `identity.profile.read` step 2 → `order.create-from-checkout` step 7 buyerEmail (REV-L2-001 closure). Pins TryClaimOrLookup atomic-claim pattern (closes dry-run #1 H02). Pins `pricing.go` as the §11.2 single source. Documents the 3-branch compensation matrix with locked `[50, 200, 800]ms` backoff. References ADR-007, ADR-008, ADR-001, ADR-003. |
