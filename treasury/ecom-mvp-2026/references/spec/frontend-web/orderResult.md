# Route spec — `/checkout/orderResult` (payment failure / expired)

- **Component:** `frontend-web`
- **Story:** [`STORY_FRONTEND_PAYMENT`](../../../requirement/EPIC_FRONTEND/STORY_FRONTEND_PAYMENT.md) (failure branch)
- **Auth:** required
- **Type:** Server Component
- **Tab route:** no

## Purpose

Failure / expiry screen after `payment.simulate` returns `FAILED` or `TIMEOUT`. Shows the reason (status-derived Thai message) and offers a retry path or a go-home action.

## Backend dependencies

| Endpoint | Why |
|---|---|
| `order.detail` (orderId from query) | Read final status (PAYMENT_FAILED / PAYMENT_EXPIRED / CANCELLED) |

## Key components

- `FailIcon` (Lucide `XCircle` for FAILED, `AlertCircle` for EXPIRED) in `text-rose-600`
- `OrderRefSummary` — orderNumber + grandTotal
- `ReasonText` — Thai message keyed by status
  - `PAYMENT_FAILED` → "การชำระเงินไม่สำเร็จ กรุณาลองอีกครั้ง"
  - `PAYMENT_EXPIRED` → "หมดเวลาชำระเงิน คำสั่งซื้อถูกยกเลิก"
  - `CANCELLED` → "คำสั่งซื้อถูกยกเลิก"
- `BtnRetryPayment` — visible only when status=`PAYMENT_FAILED` AND payment intent is still retryable; navigates to `/checkout/payment?orderId=...`
- `BtnHome` — light variant → `/`
- `BtnViewOrder` (link) → `/orders/[orderId]`

## Acceptance criteria covered

- **AC (failure branch):** customer arrives here when `payment.simulate` returns `FAILED` or `TIMEOUT`. Status-driven copy + CTAs. The retry button is gated by the contract (only when retry is legal) — for MVP, retry only when status=`PAYMENT_FAILED`.

## Edge cases

| Case | Expected |
|---|---|
| Order status is PAID (mismatched arrival) | Redirect to `/checkout/orderSuccess` |
| Retry button tapped on `PAYMENT_EXPIRED` | Button is hidden; if URL-tampered, server-side guard rejects with toast and stays on the result page |

## Notes

- This page is one of two terminal screens for the checkout flow (the other being `/checkout/orderSuccess`). Both link out to `/orders/[orderId]` for the full timeline view.
- The `ReasonText` keys live in `lib/i18n.ts` next to `STATUS_LABEL_TH`.
