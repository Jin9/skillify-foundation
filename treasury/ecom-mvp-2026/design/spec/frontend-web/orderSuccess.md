# Route spec — `/checkout/orderSuccess` (paid confirmation)

- **Component:** `frontend-web`
- **Story:** [`STORY_FRONTEND_PAYMENT`](../../../requirement/EPIC_FRONTEND/STORY_FRONTEND_PAYMENT.md) (post-success branch)
- **Auth:** required
- **Type:** Server Component
- **Tab route:** no

## Purpose

Final confirmation screen after a successful payment simulate. Surfaces the order number, grand total, and two CTAs ("ดูคำสั่งซื้อ" → `/orders/[orderId]`, "กลับสู่หน้าหลัก" → `/`).

## Backend dependencies

| Endpoint | Why |
|---|---|
| `order.detail` (orderId from query) | Confirm status=PAID; show order number + total |

## Key components

- `SuccessIcon` (Lucide `CheckCircle2` in brand emerald, 64px)
- `OrderRefSummary` — orderNumber + grandTotal
- `Btn` "ดูคำสั่งซื้อ" → `/orders/[orderId]`
- `Btn` (light variant) "กลับสู่หน้าหลัก" → `/`

## Mobile layout

- Centered icon at top
- Summary card middle
- Two stacked CTAs at bottom

## Acceptance criteria covered

- **AC (post-success transition from /checkout/payment):** customer arrives here ONLY when `payment.simulate` (or `order.detail` confirmation poll) reports status `PAID`. The page asserts `order.status==='PAID'` server-side; if not, redirects to `/checkout/orderResult` with the actual status.

## Edge cases

| Case | Expected | Implementation |
|---|---|---|
| Direct deep-link to `/checkout/orderSuccess?orderId=X` while order is still PENDING_PAYMENT | Server Component reads `order.detail` → status mismatch → redirect to `/orders/[orderId]` | Server-side `redirect()` |
| Order belongs to another customer | `order.detail` returns `ORDER_NOT_OWNED` / `NOT_FOUND` → segment `not-found.tsx` Thai message | Standard error mapping |

## Notes

- This route is short-lived in the user's mind — it's a one-screen receipt. Don't overload with details; the full timeline lives at `/orders/[orderId]`.
- `SuccessIcon` color = `--brand` emerald; reinforces the brand tokens established at `/`.
