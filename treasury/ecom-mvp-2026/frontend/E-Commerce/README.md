# ShopPilot — Next.js Frontend

Next.js 14 App Router · TypeScript 5 · Tailwind CSS 3 · pnpm

## Quick Start

```bash
# 1. Copy env
cp .env.example .env.local
# Edit .env.local with your backend service URLs

# 2. Install dependencies
pnpm install

# 3. Run dev server
pnpm dev
```

Open http://localhost:3000

## Environment Variables

| Variable | Description | Required |
|---|---|---|
| `BACKEND_IDENTITY_URL` | Identity service base URL (e.g. `http://identity:8080`) | Yes |
| `BACKEND_CATALOG_URL` | Catalog service base URL | Yes |
| `BACKEND_CART_URL` | Cart service base URL | Yes |
| `BACKEND_CHECKOUT_URL` | Checkout service base URL | Yes |
| `BACKEND_ORDER_URL` | Order service base URL | Yes |
| `BACKEND_PAYMENT_URL` | Payment service base URL | Yes |
| `BACKEND_INVENTORY_URL` | Inventory service base URL | Yes |
| `NEXT_PUBLIC_APP_NAME` | App name shown in UI | No (default: ShopPilot) |
| `NODE_ENV` | `development` or `production` | Auto-set by Next.js |

**Security note:** All `BACKEND_*` URLs are server-only (no `NEXT_PUBLIC_` prefix). They are never bundled into client JavaScript.

## Journey Demo (Happy Path)

```
/ (Home)
  → Browse categories and new arrivals
  → Click "Shop Now" or a category card

/products
  → Filter by category, price range, in-stock, search term
  → Pagination links at the bottom

/products/[id]
  → Product image, description, price, stock status
  → Click "Add to Cart" (requires login; redirects to /login?next=... if not authenticated)

/login
  → Email + password form with zod validation
  → On success: HttpOnly cookies set, redirect to previous page or /

/cart
  → View items, adjust quantity, remove items
  → Items with unavailable products flagged as non-checkoutable
  → "Proceed to Checkout" (coming in v2)
```

## Auth Flow Summary (6-Step Token Acquisition)

Lives in `lib/auth.ts → withAuth()`. Called by every authenticated route handler.

1. Read `access` HttpOnly cookie from the incoming request.
2. If present and not expired (30 s clock-skew tolerance) → forward `Authorization: Bearer <access>` to backend.
3. If absent or expired → read `refresh` cookie. If missing → return 401 `AUTH_MISSING`.
4. Call `identity.refresh` with the refresh token.
5. On success → set-cookie both rotated tokens, forward downstream with new access token.
6. On `AUTH_REVOKED`/`AUTH_INVALID` → clear both cookies (Max-Age=0), return 401 `AUTH_REVOKED`; browser redirects to `/login`.

Cookies are `HttpOnly; SameSite=Lax; Path=/` (plus `Secure` in production). No JWT ever reaches JavaScript.

## What Is FULL vs Stubbed

### FULL Pages (5)
| Route | Notes |
|---|---|
| `/` | RSC — hero banner, category grid, top 12 products |
| `/products` | RSC + Client filter pane; URL params drive server fetch |
| `/products/[id]` | RSC product detail + Client AddToCartButton |
| `/login` | Client form (react-hook-form + zod); sets cookies on success |
| `/cart` | Client — fetches cart, qty stepper, remove, subtotal |

### STUB Pages ("Coming in v2")
`/register`, `/checkout`, `/checkout/payment`, `/orders`, `/orders/[id]`, `/admin`

### FULL Route Handlers (5)
| Handler | Description |
|---|---|
| `GET /api/auth/session` | Session probe — decodes access cookie, no refresh side-effect |
| `POST /api/proxy/identity/auth/login` | Login → sets both HttpOnly cookies, strips tokens from client response |
| `POST /api/proxy/identity/auth/refresh` | Explicit refresh → rotates both cookies |
| `POST /api/proxy/catalog/product/list` | Public passthrough (forwards Bearer if present) |
| `POST /api/proxy/catalog/product/detail` | Public passthrough |

### STUB Route Handler (catch-all)
`POST /api/proxy/[...path]` — forwards all other routes to the correct backend with the 6-step auth flow. Does NOT implement per-route role guards or idempotency-key forwarding (flagged for v2 dedicated handlers).

## Security Headers (CSP)

Set in `next.config.js` via `headers()`:

```
default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; connect-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'
```

`X-Frame-Options: DENY`, `X-Content-Type-Options: nosniff`, `Referrer-Policy: strict-origin-when-cross-origin`, `Permissions-Policy: geolocation=(), microphone=(), camera=()`.

HSTS is added in production only.

**`unsafe-inline` for `script-src`** is a known Next.js App Router constraint pending nonce-middleware (tracked for v2).

## Directory Structure

```
app/
  layout.tsx           Root layout (fonts, globals.css, Header)
  globals.css          Tailwind directives
  page.tsx             Home (FULL)
  products/page.tsx    Product listing (FULL)
  products/[id]/       Product detail (FULL)
  login/page.tsx       Login (FULL)
  cart/page.tsx        Cart (FULL)
  register/page.tsx    STUB
  checkout/            STUB
  orders/              STUB
  admin/               STUB
  api/auth/session/    FULL GET handler
  api/proxy/identity/auth/login/   FULL
  api/proxy/identity/auth/refresh/ FULL
  api/proxy/catalog/product/list/  FULL
  api/proxy/catalog/product/detail/FULL
  api/proxy/[...path]/ Catch-all STUB proxy
lib/
  auth.ts              6-step token acquisition + cookie helpers
  http.ts              Backend fetch wrapper + service URL map
  idempotency.ts       UUID generator + InFlightMap
components/
  Header.tsx           Client nav with login state
  ProductCard.tsx      Server product tile
  ProductFilters.tsx   Client filter pane
  CartLineItem.tsx     Client cart row
  AddToCartButton.tsx  Client add-to-cart
```
