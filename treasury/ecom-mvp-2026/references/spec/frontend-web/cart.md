# Route spec — `/cart` (cart tab)

- **Component:** `frontend-web`
- **Story:** [`STORY_FRONTEND_CART`](../../../requirement/EPIC_FRONTEND/STORY_FRONTEND_CART.md)
- **Auth:** public read; write actions trigger AUTH_MISSING → /login flow
- **Type:** Server Component shell + Client line controls
- **Tab route:** yes

## Purpose

Staging ground for checkout — lists every cart line with image, name, ฿ unit price, quantity stepper, remove icon, line subtotal, and a server-computed cart subtotal. The CTA "ดำเนินการชำระเงิน" navigates to `/checkout` when every line is `checkoutable=true`.

## Backend dependencies

| Endpoint | Why |
|---|---|
| `cart.read` | Initial server-rendered cart contents (subtotal, items[], checkoutable flags per CART-005) |
| `cart.update-item` | Quantity stepper +/- |
| `cart.remove-item` | Remove icon |

All three go through `withAuth`. An unauthenticated cart read returns AUTH_MISSING; the layout redirects to `/login?next=/cart`.

## Key components

- `ShopHeader`
- `CartLineItem` (`'use client'`) — image + name + ฿ price + qty stepper + remove + line subtotal + non-checkoutable badge
- `CartSummary` — server-supplied subtotal block
- `BtnCheckout` — primary CTA, disabled when any line is `checkoutable=false`
- `EmptyCart` — placeholder with a CTA back to `/`
- `TabBar`

## Mobile layout

- Single column; each line is a card with 14px radius
- Stepper buttons sized 44x44 for tap targets
- Sticky `CartSummary` at the bottom above `TabBar`

## Acceptance criteria covered

- **AC1 (lines + server subtotal):** Each line shows the required atoms; the cart subtotal equals the sum of `lineSubtotal` returned by `cart.read` (frontend NEVER recomputes).
- **AC2 (stepper updates qty + subtotal):** Tap `+` → `cart.update-item` qty=current+1 → re-renders with new server subtotal.

## Edge cases

| Case | Expected | Implementation |
|---|---|---|
| Step over `availableQty` | `cart.update-item` returns `INSUFFICIENT_STOCK`; toast in Thai with `availableQty` hint; qty unchanged | `error_handling.ui_mapping.INSUFFICIENT_STOCK` |
| `checkoutable=false` line (PRODUCT_INACTIVE / DELETED) | Visual badge "ไม่สามารถสั่งซื้อได้"; checkout CTA disabled until removed | CART-005 + CHK-003 |
| Empty cart | Render `EmptyCart` with CTA → `/` | Server Component conditional |
| Auth missing on a write action | Redirect to `/login?next=/cart`; on success, return to `/cart` | `withAuth` returns `AUTH_MISSING` envelope; client interprets and routes |

## Notes

- Subtotal is always server-truth from `cart.read`; the route does NOT compute totals locally — important for boundary correctness with §11.2.
- Quantity mutations re-fetch `cart.read` to keep state coherent; optimistic update is allowed but the next render reconciles to server.
