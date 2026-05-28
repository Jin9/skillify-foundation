# ShopPilot Frontend — LLM Handoff Spec

> Use this document + `ShopPilot Mobile (standalone).html` as the source of truth to rebuild the ShopPilot customer-facing frontend in your stack of choice (Next.js, Vite + React, Remix, etc.).

---

## 1. Product summary

ShopPilot is a Thai-language B2C fashion e-commerce mobile webview. MVP scope is the **customer-facing flow only**: browse → cart → checkout → mock payment → order tracking → reviews. Admin is out of scope.

- **Audience:** Thai shoppers on mobile browsers (iOS Safari + Android Chrome WebView)
- **Language:** Thai (`lang="th"`)
- **Currency:** Thai Baht, prefix format `฿1,290` (toggleable to `1,290 THB`)
- **Target viewport:** 390×844 (iPhone 14) — single column, mobile-first

---

## 2. Visual system

### Brand
- **Primary (`--brand`):** `#1a8060` (emerald)
- **Backgrounds:** `#fff` (cards/surfaces), `#f7f7f5` (page bg), `#f4f4f1` (dividers)
- **Text:** `#111` primary, `rgba(0,0,0,0.55)` secondary, `rgba(0,0,0,0.4)` tertiary
- **Accent palette options:** `#1a8060`, `#0f766e`, `#15803d`, `#ea580c`, `#7c3aed`, `#111`

### Typography
- **Default:** IBM Plex Sans Thai 400/500/600/700
- **Alternates:** Sarabun, Prompt
- Type scale: 10 / 11 / 12 / 13 / 14 / 15 / 17 / 22 / 26 (titles)

### Spacing & shape
- Card radius: `14–16px`; pill radius: `999px`; buttons: `12–14px`
- Page padding: `16px` (regular), `12px` (compact density)
- Card gap: `10–12px`; section spacing: `8–10px` between cards

### Iconography
- Outline style, 1.8 stroke width, 18–22px sizes
- Icons are inlined SVG — see `src/primitives.jsx → Icon` for the full set: home, search, bag, package, user, back, close, chevron, plus, minus, check, heart, pin, tag, truck, card, qr, filter, trash, edit, info

### Imagery
Striped diagonal placeholders with category labels (NOT real photos) — to be replaced with product photos post-launch. 7 tonal variants seeded by `product.tone` index.

---

## 3. Information architecture & routes

| Route | Screen | Tab? | Auth? |
|---|---|---|---|
| `home` | Featured + categories + new arrivals | ✓ | — |
| `search` | Search + filters + product grid | ✓ | — |
| `cart` | Cart line items + summary | ✓ | — |
| `orders` | Order history (filter: all / active / done) | ✓ | required |
| `profile` | Profile + settings + logout | ✓ | required |
| `detail` | Product detail page | — | — |
| `login` | Login / register toggle | — | — |
| `checkout` | Address + items + coupon + summary | — | required |
| `payment` | Mock QR + simulate Success/Fail/Timeout | — | required |
| `orderSuccess` | Confirmation screen | — | required |
| `orderResult` | Failed / expired result | — | required |
| `orderDetail` | Items + status timeline + tracking | — | required |
| `addresses` | List + add new (also picker mode) | — | required |
| `review` | Star rating + comment for an order item | — | required |

Bottom tab bar shows on the 5 tab routes only. Detail/checkout/etc. take over the full screen.

---

## 4. Data models

### Product
```ts
{ id, sku, name, cat, label, tone, price, compareAt?, rating, reviews, stock, desc }
```

### Cart item
```ts
{ id, productId, qty }
```

### Order
```ts
{
  id, orderNumber: 'SP-25-XXXXXX',
  status: 'PENDING_PAYMENT' | 'PAID' | 'PAYMENT_FAILED' | 'PAYMENT_EXPIRED'
        | 'PACKING' | 'SHIPPED' | 'DELIVERED' | 'CANCELLED',
  createdAt: ISO8601,
  address: Address,
  items: OrderItem[],
  subtotal, discount, shippingFee, grandTotal,
  couponCode?, tracking?, carrier?,
  history: { status, at }[]
}
```

### Coupon
```ts
{ code, type: 'fixed' | 'percent' | 'shipping', amount, label, minOrder }
```
Seeded coupons: `WELCOME100` (฿100 off, min ฿500), `SAVE10` (10%, min ฿1000), `FREESHIP`.

---

## 5. Business logic (port verbatim)

### Totals calculation
```js
shippingFee = subtotal >= 1500 ? 0 : 60       // free shipping ≥ ฿1,500
grandTotal  = max(0, subtotal - discount + shippingFee)
```
Coupon validation: requires `subtotal >= coupon.minOrder`. `fixed` subtracts amount, `percent` takes `% of subtotal`, `shipping` zeros the shipping fee.

### Order lifecycle
`PENDING_PAYMENT` → (mock payment) → `PAID` | `PAYMENT_FAILED` | `PAYMENT_EXPIRED`
`PAID` → `PACKING` → `SHIPPED` (gets tracking + carrier) → `DELIVERED` → enables review writing.
`PENDING_PAYMENT` is the only state user can cancel from.

### Payment (mock only — no real gateway)
3 outcomes simulated via buttons after a 1.4s "processing" delay. Real implementation should swap in a payment SDK (Omise, Stripe Thailand, etc.) and consume webhook callbacks.

---

## 6. Component inventory

Reusable atoms (rebuild from `src/primitives.jsx`):
- `<Icon name size color stroke/>`
- `<Btn variant="brand|dark|light|ghost" full disabled/>`
- `<Chip active onClick/>` — pill filter
- `<Stars rating size/>` — 5 star display
- `<ProductImage label tone height radius/>` — striped placeholder

Composites (rebuild from `src/screens-shop.jsx`, `src/screens-flow.jsx`):
- `<TabBar/>` — bottom 5-tab nav
- `<ShopHeader title onBack action/>` — top bar
- `<ProductCard product onClick compact/>`
- `<Section title action onAction/>` — white rounded container
- `<SummaryRow label value bold color/>` — line in totals
- `<StatusBadge status/>` — colored pill with Thai label
- `<Timeline order/>` — 5-step horizontal progress
- `<Sheet title onClose/>` — bottom modal sheet
- `<Field label value onChange placeholder type/>` — form input

---

## 7. Tweakable settings (expose in build)

- Brand color (curated 6 swatches)
- Font (IBM Plex Sans Thai / Sarabun / Prompt)
- Density (compact / regular)
- Currency format (`฿1,290` / `1,290 THB`)

---

## 8. Recommended tech stack

- **Framework:** Next.js 14 App Router (or Vite + React Router)
- **State:** Zustand or React Context (the prototype uses Context + localStorage persistence)
- **Styling:** Tailwind with the brand color as `theme.colors.brand`, OR vanilla CSS with the variables shown above
- **Fonts:** `next/font/google` for IBM Plex Sans Thai
- **Icons:** Lucide React (close substitutes for the inlined SVGs above)
- **API contract:** mirror the data models in §4; the prototype's `localStorage` calls map 1:1 to `GET /api/cart`, `POST /api/orders`, etc.

---

## 9. How to use this spec with an LLM

1. **Open** `ShopPilot Mobile (standalone).html` in a browser to interact with the live prototype.
2. **Paste this file + the standalone HTML** into your LLM of choice with a prompt like:
   > "Rebuild this prototype as a production Next.js 14 app with Tailwind. Match the visual system, routes, and data models exactly. Use server actions for cart/order mutations and replace `localStorage` with a Postgres + Drizzle backend. Keep all Thai copy verbatim."
3. **For incremental work**, attach single source files (`src/screens-shop.jsx` etc.) and ask for one screen at a time.

---

## 10. Files in this export

```
ShopPilot Mobile (standalone).html   ← single-file prototype, drop into any browser
ShopPilot Frontend Spec.md           ← this document
```

Source files in the original project (for reference):
```
ShopPilot Mobile.html       Entry point, loads scripts
src/data.js                 Products, categories, coupons, sample orders, totals fn
src/primitives.jsx          Icon, Btn, Chip, Stars, ProductImage
src/screens-shop.jsx        Home, Search, Detail, Cart, TabBar, Header
src/screens-flow.jsx        Checkout, Payment, Orders, Profile, Auth, Review
src/app.jsx                 Router, store hook, Tweaks panel
src/ios-frame.jsx           iOS device chrome (status bar, dynamic island, home indicator)
src/android-frame.jsx       Android Material 3 chrome
```
