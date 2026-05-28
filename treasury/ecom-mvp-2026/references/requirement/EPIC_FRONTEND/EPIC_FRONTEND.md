---
template_version: 0.1.0
epic_id: EPIC_FRONTEND
title: Thai-language mobile customer webview (5-tab nav + 8 full-screen routes)
summary: |
  The customer-facing mobile webview shaped by ShopPilot Frontend Spec. Thai language
  (lang=th), 390×844 mobile-first viewport, IBM Plex Sans Thai, emerald #1a8060 brand,
  ฿ Thai-Baht prefix. 5 bottom-tab routes (home/search/cart/orders/profile) + 8
  full-screen routes (detail, login, checkout, payment, orderSuccess, orderResult,
  orderDetail, addresses) per Frontend Spec §3. Admin UI is explicitly out of scope.
business_value: The webview is the entire customer-facing surface for ShopPilot MVP — it is the product the user actually sees. Backend services exist to power this experience.
in_scope_stories:
  - STORY_FRONTEND_HOME
  - STORY_FRONTEND_SEARCH
  - STORY_FRONTEND_DETAIL
  - STORY_FRONTEND_AUTH
  - STORY_FRONTEND_CART
  - STORY_FRONTEND_CHECKOUT
  - STORY_FRONTEND_PAYMENT
  - STORY_FRONTEND_ORDERS
  - STORY_FRONTEND_ORDER_DETAIL
  - STORY_FRONTEND_ADDRESSES
  - STORY_FRONTEND_PROFILE
  - STORY_FRONTEND_AUTH_GUARD
out_of_scope:
  - Admin web UI — Frontend Spec is customer-facing only; admin lives behind backend endpoints (Postman/curl) this run
  - Review route on orderDetail — review module deferred (§8.12)
  - Coupon UX on checkout — coupon module deferred (§8.7); checkout coupon row renders as static disabled placeholder
  - Native iOS/Android app — Frontend Spec ships a mobile webview, not native apps (§4.2 #10)
  - Real product photography — Frontend Spec ships striped diagonal placeholders with category labels
  - Density/font/brand-swatch tweak panel from Frontend Spec §7 — designer-time only, not a runtime customer feature
  - Multi-language UI — Thai-only this run (§4.2 #12); copy is Thai verbatim
success_metrics:
  - All 13 routes (5 tab + 8 full-screen) renderable on a 390×844 viewport without horizontal scroll
  - Auth-required routes redirect to login when no token is present and resume original route on login (Frontend Spec §3 'Auth?' column)
  - The customer journey home → detail → cart → checkout → payment-success → orderDetail completes end-to-end with status=PAID
  - Currency renders as ฿1,290 prefix format throughout (Frontend Spec §1)
  - Primary CTA colour matches the emerald #1a8060 brand token across all routes
dependencies:
  - EPIC_AUTH
  - EPIC_CATALOG
  - EPIC_CART
  - EPIC_CHECKOUT
  - EPIC_PAYMENT
  - EPIC_ORDER
  - EPIC_INVENTORY
compliance_sensitive: false
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created NEW for dry-run #2 to absorb the Frontend Spec; runs/01 had no dedicated frontend epic.
---

# EPIC_FRONTEND — Thai-language mobile customer webview (5-tab nav + 8 full-screen routes)

## Why

The customer mobile webview is the entire surface area the user actually touches. The Frontend Spec (§1–§10) shapes it: Thai-language fashion e-commerce on a 390×844 viewport with IBM Plex Sans Thai, emerald `#1a8060` brand, the Thai-Baht prefix `฿1,290`, and a 5-tab bottom navigation.

The Frontend Spec §3 declares 13 routes:

- **Tab routes (5, with bottom tab bar):** home, search, cart, orders, profile
- **Full-screen routes (8, no tab bar):** detail, login, checkout, payment, orderSuccess, orderResult, orderDetail, addresses

The `Auth?` column distinguishes:
- **Public:** home, search, cart (browsable as guest until checkout), detail, login
- **Required:** orders, profile, checkout, payment, orderSuccess, orderResult, orderDetail, addresses

Each story below owns the *behavioural contract* for one or two routes. The route's visual contract (Thai labels, ฿ prefix, emerald CTAs, tab-bar visibility) is asserted at the AC level — but the AC stops short of prescribing React component structure or class names. That's Tech-Designer's space.

The deliberate scope cut: **no admin UI.** The original requirement §8.13 (Admin Back Office) is partially in scope at the *backend* layer (admin endpoints for product CRUD, category, stock, order transitions exist) but no admin-facing routes are built this run. This is the single biggest interpretive call in this dry-run versus runs/01, which had `frontend-web: Admin pages allow ...` in scope.

## Stories in this epic

| Story ID | Title | Priority |
|---|---|---|
| `STORY_FRONTEND_HOME` | home tab renders featured + categories + new arrivals (Thai labels, ฿ prefix, emerald CTAs, 5-tab bar) | Must |
| `STORY_FRONTEND_SEARCH` | search tab supports query input, category filter chips, and product grid with pagination | Must |
| `STORY_FRONTEND_DETAIL` | detail route shows product info, ฿ price, stock status, add-to-cart CTA | Must |
| `STORY_FRONTEND_AUTH` | login route supports register and login toggle; on success token is stored client-side | Must |
| `STORY_FRONTEND_CART` | cart tab lists lines with quantity stepper, remove, server-computed subtotal | Must |
| `STORY_FRONTEND_CHECKOUT` | checkout route shows address + items + server-computed totals (incl. §11.2) and place-order CTA with Idempotency-Key | Must |
| `STORY_FRONTEND_PAYMENT` | payment route exposes mock-QR + 3 simulate buttons; navigates by service-returned status | Must |
| `STORY_FRONTEND_ORDERS` | orders tab lists customer's orders with all/active/done filter and Thai status badges | Must |
| `STORY_FRONTEND_ORDER_DETAIL` | orderDetail route shows item snapshots, address snapshot, status timeline, tracking when shipped | Must |
| `STORY_FRONTEND_ADDRESSES` | addresses route lists/add addresses; supports picker mode for checkout | Must |
| `STORY_FRONTEND_PROFILE` | profile tab shows name/email/default address and exposes logout | Must |
| `STORY_FRONTEND_AUTH_GUARD` | Auth-required routes redirect to login and resume original route on success | Must |

## Out-of-scope rationale

- **Admin UI** — Frontend Spec is customer-facing only; admin endpoints exist but no admin routes are built this run. This is locked at the BA level so Tech-Lead does not invent admin pages downstream.
- **Review route** — review module deferred (§8.12); the orderDetail route does not surface a "write a review" CTA in MVP.
- **Coupon UX on checkout** — coupon module deferred (§8.7); the Frontend Spec checkout coupon row renders as a static disabled placeholder with copy "เร็ว ๆ นี้" (or equivalent).
- **Native apps** — §4.2 #10 deferred; webview only.
- **Real product imagery** — Frontend Spec §2 explicitly ships striped diagonal placeholders.
- **Tweak panel** — Frontend Spec §7's font/density/brand-swatch tweaker is a designer-time configuration, not a runtime customer toggle.
- **Multi-language** — §4.2 #12 deferred; Thai-only this run with the data model untouched (no i18n schema changes).

## Risks

- **Auth-redirect loop** — STORY_FRONTEND_AUTH_GUARD must remember the originally requested route and resume it on login; otherwise customers who deep-link into orderDetail get bounced back to home and lose context.
- **Token storage** — the access token must be attached to subsequent API calls. Tech-Designer chooses storage (localStorage, sessionStorage, in-memory + cookie) but the BA contract is "subsequent API calls succeed without re-prompting until token expires".
- **§11.2 boundary visual rendering** — the checkout route's totals block must update when subtotal crosses ฿1,500 in either direction. The frontend renders only what the server returns from checkout.preview; it must not compute the shipping fee locally.
- **Payment race** — the payment route's 1.4s processing animation is cosmetic; navigation must depend on the payment service's returned status. A frontend that navigates optimistically to orderSuccess before the service confirms the transition is a bug.
- **Thai-copy verbatim** — the Frontend Spec's prompt template ("Keep all Thai copy verbatim") implies copy regression risk if Tech-Designer paraphrases. BA does not specify the copy strings here; they are owned by the Frontend Spec source HTML.
- **Mobile-3G / narrow viewport** — the layout must reflow at 360-wide without horizontal scroll; the 390×844 figure is the design target, not a hard floor.

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created NEW for dry-run #2 to absorb the Frontend Spec; runs/01 had no dedicated frontend epic. |
