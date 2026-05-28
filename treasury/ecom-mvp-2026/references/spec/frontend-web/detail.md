# Route spec — `/detail/[productId]` (product detail)

- **Component:** `frontend-web`
- **Story:** [`STORY_FRONTEND_DETAIL`](../../../requirement/EPIC_FRONTEND/STORY_FRONTEND_DETAIL.md)
- **Auth:** public (add-to-cart prompts login when unauthenticated)
- **Type:** Server Component (with Client `AddToCartButton` island)
- **Tab route:** no (full-screen)

## Purpose

Conversion surface — the customer reads the product's details and decides to add to cart. Renders name, image placeholder, ฿ price (and `compareAt` if present), stock status badge, description text, and an emerald "เพิ่มลงตะกร้า" CTA.

## Backend dependencies

| Endpoint | Why |
|---|---|
| `catalog.product.detail` (productId) | Full product fields: name, price, compareAt, stockStatus, description, images[], categoryId |
| `cart.add-item` (POST `/api/proxy/cart/add-item`) | Add CTA |

`catalog.product.detail` is public; `cart.add-item` requires auth via `withAuth`.

## Key components

- `ShopHeader` (back button)
- `ProductImage` (large striped placeholder)
- `ProductPrice` (with `compareAt` strikethrough if present)
- `StockBadge` (`พร้อมจำหน่าย` / `สินค้าหมด` / `เหลือน้อย`)
- `ProductDescription`
- `AddToCartButton` (`'use client'`) — disabled when out-of-stock; on tap if unauthenticated → `router.push('/login?next=' + currentPath)`

## Mobile layout

- Full-bleed image at top
- Sticky bottom CTA above iOS home-indicator safe area (px-4 py-3 background `#fff` shadow)

## Acceptance criteria covered

- **AC1 (renders all fields + auth-gated CTA):** name, image placeholder, ฿ price, stock status, description, emerald CTA. Tap when authenticated → `cart.add-item`; success toast.

## Edge cases

| Case | Expected | Implementation |
|---|---|---|
| Unauth tap on CTA | Redirect to `/login?next=/detail/[productId]`; resume add-to-cart on login (caller stores intent in URL or sessionStorage and replays once after `/api/auth/session` confirms auth) | Client island reads cookie via `/api/auth/session` GET; routes accordingly |
| `stockStatus='out_of_stock'` | CTA visually disabled and tap is a no-op | Client guard before fetch |
| Product not found | Server Component throws → segment `not-found.tsx` renders Thai message | Standard Next.js |

## Notes

- `cart.add-item` returns the updated cart envelope; on success, the page can refresh the header `cart-count` badge via a SWR-less `router.refresh()` of the cart-bearing layout, or a simple toast. MVP: toast only; cart count updates on next `/cart` visit.
