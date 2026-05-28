# Route spec — `/checkout` (checkout review + place-order)

- **Component:** `frontend-web`
- **Story:** [`STORY_FRONTEND_CHECKOUT`](../../../requirement/EPIC_FRONTEND/STORY_FRONTEND_CHECKOUT.md)
- **Auth:** required
- **Type:** Client Component
- **Tab route:** no

## Purpose

Customer's last review surface before payment. Displays the picked address, order items (snapshot of cart at preview time), the §11.2 server-priced totals (subtotal + shippingFee + grandTotal), and a "ยืนยันคำสั่งซื้อ" CTA that POSTs to `/api/checkout/commit` with a client-generated `idempotencyKey` (UUID v4) per click.

## Backend dependencies

| Endpoint | Why |
|---|---|
| `identity.address.list` | Address picker (default selected) |
| `checkout.preview` | Server-priced totals; re-runs whenever shippingAddressId changes |
| `checkout.commit` | Place-order action; idempotent on customer-supplied key |

## Key components

- `ShopHeader` (back)
- `AddressPicker` — selectable list; tap → "เปลี่ยนที่อยู่" navigates to `/addresses?mode=pick&returnTo=/checkout`
- `OrderItemList` — snapshot from preview (`productNameSnapshot`, `priceSnapshot`, `qty`, `lineSubtotal`)
- `CouponRow` — disabled placeholder "เร็ว ๆ นี้" per EPIC_FRONTEND.out_of_scope
- `CheckoutSummary` — subtotal + shippingFee + couponDiscount(=0) + grandTotal verbatim from server
- `ShippingFeeExplainer` — small text "ส่งฟรีเมื่อสั่งครบ ฿1,500" under the shippingFee row
- `CheckoutCommitButton` (`'use client'`) — generates `idempotencyKey` on click, debounces via `InFlightMap`

## Mobile layout

- Sticky bottom action bar with grandTotal + commit button above iOS safe area
- Body: address card → items list → summary

## Acceptance criteria covered

- **AC1 (1500 boundary):** subtotal=1500 → shippingFee=฿0, grandTotal=฿1,500 (Frontend Spec §5 + §11.2)
- **AC2 (1499 boundary):** subtotal=1499 → shippingFee=฿60, grandTotal=฿1,559 (Frontend Spec §5 + §11.2)
- **AC3 (idempotent commit):** Place-order POSTs `idempotencyKey` (UUID v4); on success → `router.push('/checkout/payment?orderId=...')`; on retry the SAME key returns the SAME order id (CHK-009).

## Edge cases

| Case | Expected | Implementation |
|---|---|---|
| Cart changes between preview and commit | Commit recomputes shippingFee from cart at commit time; preview is advisory | Backend handles; frontend just renders commit response |
| Double-tap place-order | `InFlightMap` key `checkout-commit` blocks the second click; the same key is sent on a forced retry; server returns the same order id | `idempotency_strategy.in_flight_dedup` |
| `INSUFFICIENT_STOCK` 409 | Toast Thai message + per-item availableQty inline; commit blocked until cart adjusted | `error_handling.ui_mapping.INSUFFICIENT_STOCK` |
| `IDEMPOTENCY_KEY_REUSED` | Read `envelope.data.orderId`; toast "คำสั่งซื้อนี้ถูกสร้างไว้แล้ว"; `router.push('/orders/[orderId]')`; do NOT retry | `idempotency_strategy.409_handling` |
| `IDEMPOTENCY_KEY_INFLIGHT` | Honor `Retry-After` (1s); retry once with same key; if still inflight, soft toast and stop | `idempotency_strategy.409_handling` |
| Address picked from `/addresses?mode=pick` | URL query `?addressId=...` triggers preview re-fetch | `useEffect([addressId])` |

## Sequence — checkout commit (with idempotency)

```mermaid
sequenceDiagram
  participant U as User
  participant B as Browser (CheckoutCommitButton)
  participant H as Route Handler (/api/checkout/commit)
  participant C as Backend Checkout

  U->>B: tap "ยืนยันคำสั่งซื้อ"
  B->>B: const key = crypto.randomUUID()
  B->>B: InFlightMap.has('checkout-commit')? skip : continue
  B->>H: POST {idempotencyKey: key, shippingAddressId, ...}
  H->>H: withAuth: read access cookie / refresh if needed
  H->>C: POST /api/v1/checkout/checkout/commit\nIdempotency-Key: <key>\nAuthorization: Bearer <access>
  alt First call (no inflight, no reuse)
    C-->>H: 201 CREATED {orderId, paymentIntentId, ...}
    H-->>B: 201 envelope
    B->>B: router.push('/checkout/payment?orderId=...')
  else Same payload, replay
    C-->>H: 200 SUCCESS (cached envelope)
    H-->>B: 200 envelope (same orderId)
    B->>B: router.push('/checkout/payment?orderId=...')
  else Same key, different payload
    C-->>H: 409 IDEMPOTENCY_KEY_REUSED {orderId}
    H-->>B: 409 envelope
    B->>B: router.push('/orders/[orderId]'); toast
  else Concurrent inflight
    C-->>H: 409 IDEMPOTENCY_KEY_INFLIGHT (Retry-After: 1)
    H-->>B: 409 envelope
    B->>B: wait 1s, retry once with SAME key
  end
```

## Notes

- `CheckoutSummary` NEVER recomputes pricing client-side. It renders exactly what the server returns from `checkout.preview` (initial render) and `checkout.commit` (response). This is the §11.2 boundary contract from ADR-010.
- `idempotencyKey` is held in `useState` only — not persisted. A page refresh discards it; the user re-clicks; a fresh key is generated; the previous orderId is what survives via the URL the user was redirected to.
