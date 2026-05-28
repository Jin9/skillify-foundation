# Route spec — `/checkout/payment` (mock payment simulate)

- **Component:** `frontend-web`
- **Story:** [`STORY_FRONTEND_PAYMENT`](../../../requirement/EPIC_FRONTEND/STORY_FRONTEND_PAYMENT.md)
- **Auth:** required
- **Type:** Client Component
- **Tab route:** no

## Purpose

Mock-payment surface — three buttons (`สำเร็จ` / `ไม่สำเร็จ` / `หมดเวลา`) trigger `payment.simulate` with `outcome=SUCCESS|FAILED|TIMEOUT`. Navigation depends on the service's returned status, NEVER an optimistic guess.

## Backend dependencies

| Endpoint | Why |
|---|---|
| `payment.simulate` (POST `/api/payment/simulate`) | Outcome trigger |
| `order.detail` (poll if needed) | Confirm status flip after simulate |

## Key components

- `ShopHeader` (back)
- `MockQrCard` — static visual (no real QR data)
- `OrderRefSummary` — orderNumber + grandTotal in ฿
- `PaymentSimulateButtons` (`'use client'`) — three Btn variants in a row; ~1.4s processing indicator
- `ProcessingOverlay` — spinner + Thai "กำลังประมวลผล…"

## Mobile layout

- QR placeholder card centered
- Order summary below
- Three buttons stacked at bottom; brand-emerald success, amber failure, slate timeout

## Acceptance criteria covered

- **AC1 (3 outcomes drive navigation):**
  - SUCCESS → `/checkout/orderSuccess?orderId=...`
  - FAILED → `/checkout/orderResult?orderId=...&status=PAYMENT_FAILED`
  - TIMEOUT → `/checkout/orderResult?orderId=...&status=PAYMENT_EXPIRED`
- **AC2 (~1.4s processing indicator; navigation by service status, NOT optimistic):** Tap → show `ProcessingOverlay` → POST `/api/payment/simulate` → wait response → if upstream confirms transition (or echoes terminal status), route accordingly. Frontend NEVER pre-decides the destination.

## Edge cases

| Case | Expected | Implementation |
|---|---|---|
| Customer closes browser before tapping | Order remains `PENDING_PAYMENT`; eventually flips to `PAYMENT_EXPIRED` via reservation TTL sweeper (PAY-004 / INV-004); next visit to `/orders` shows it expired | Out of frontend control |
| Double-tap simulate-success | Payment service is idempotent on `(paymentIntentId, providerStatus)`; one transition; one navigation | `InFlightMap` key `payment-simulate` blocks duplicate clicks client-side too |
| `PAYMENT_INTENT_TERMINAL` | Toast + redirect to `/orders/[orderId]` | `error_handling.ui_mapping` |
| Network failure on simulate | Error toast with traceId; buttons re-enabled for retry | Standard `ApiError` |

## Notes

- The 1.4s processing animation is COSMETIC — it gives the customer a moment of feedback while the service responds. Navigation is gated on the response, not on the timer.
- `MockQrCard` is intentionally static visual; no QR generation library imported.
- For the SUCCESS branch: the simulate response carries the new payment status; if the backend's order-state propagation has not yet caught up to PAID by the time the browser reaches `/checkout/orderSuccess`, the success page calls `order.detail` once more to confirm — short retry loop (max 8 polls @ 750ms = 6s) before falling back to the orderResult route.
