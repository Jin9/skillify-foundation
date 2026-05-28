# Route spec — `/orders/[orderId]` (order detail + timeline)

- **Component:** `frontend-web`
- **Story:** [`STORY_FRONTEND_ORDER_DETAIL`](../../../requirement/EPIC_FRONTEND/STORY_FRONTEND_ORDER_DETAIL.md)
- **Auth:** required
- **Type:** Server Component (with Client `BtnCancel` island)
- **Tab route:** no

## Purpose

End of the customer journey loop — surfaces what was bought, where it ships, what status it's in, and (when shipped) the tracking number + carrier. Renders snapshot fields so a soft-deleted product still displays in a historical order.

## Backend dependencies

| Endpoint | Why |
|---|---|
| `order.detail` (orderId) | Full order tree: status, addressSnapshot, items[] (with priceSnapshot/productNameSnapshot/productImageUrlSnapshot), statusHistory[], trackingNumber, carrier, totals |
| `order.cancel-mine` | Customer cancel (only when status=`PENDING_PAYMENT`) |

## Key components

- `ShopHeader` (back)
- `StatusBadge`
- `OrderStatusTimeline` — 5-step horizontal progress (PENDING_PAYMENT → PAID → PACKING → SHIPPED → DELIVERED); current status highlighted; failure statuses (PAYMENT_FAILED / PAYMENT_EXPIRED / CANCELLED) short-circuit the chain at the failure point with a rose marker
- `AddressSnapshotCard` — receiverName, phone, addressLine, district, province, postalCode (server snapshot; immune to address edits after order)
- `OrderItemRow` — uses snapshot fields (`productNameSnapshot`, `priceSnapshot`, `productImageUrlSnapshot` per CAT-009 / ORD-007)
- `TrackingCard` — visible only when status in {SHIPPED, DELIVERED}; renders `trackingNumber` + `carrier` (mock_express label)
- `CheckoutSummary` (read-only) — final subtotal + shippingFee + grandTotal
- `BtnCancel` (`'use client'`) — visible only when status=`PENDING_PAYMENT` (ORD-004); confirm dialog before POSTing `order.cancel-mine`

## Mobile layout

- Top-to-bottom: status badge + timeline → address snapshot → items list → summary → cancel button (when applicable)
- Sticky-free; the page is read-heavy

## Acceptance criteria covered

- **AC1 (end-to-end completes; PAID + timeline rows):** after the §7.1 spine, the page renders status="ชำระเงินแล้ว" with at least PENDING_PAYMENT and PAID rows in `statusHistory` (ORD-002, ORD-009).
- **AC2 (SHIPPED tracking + carrier; full timeline):** when status=SHIPPED, the page shows trackingNumber, carrier, and statusHistory entries for every transition through PENDING_PAYMENT → PAID → PACKING → SHIPPED with timestamps (ORD-006, ORD-009).

## Edge cases

| Case | Expected | Implementation |
|---|---|---|
| Product was soft-deleted after order | Item row renders snapshot fields; page does not break (CAT-009 / ORD-007) | Snapshot fields are persisted on `order_items`; we render `productNameSnapshot` not the live product name |
| Order belongs to another customer | `order.detail` returns `ORDER_NOT_OWNED` / `NOT_FOUND` → segment `not-found.tsx` Thai message (ORD-001) | Standard error mapping |
| Customer cancels a `PENDING_PAYMENT` order | `BtnCancel` confirm dialog → `order.cancel-mine` → status flips to CANCELLED; `BtnCancel` hides; timeline shows the CANCELLED marker | ORD-004 |
| Timestamps formatting | `Intl.DateTimeFormat('th-TH-u-ca-buddhist', {...})` for Buddhist-era display | `lib/format.ts` |

## Notes

- The Server Component fetches `order.detail` once on render. The cancel button is the only Client island.
- `OrderStatusTimeline` reads `statusHistory` ascending; the timeline component has no business logic of its own — it renders what the backend provides.
