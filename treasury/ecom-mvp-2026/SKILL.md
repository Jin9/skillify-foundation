---
name: ecom-mvp-2026
description: Design and build 7 backend services plus a Thai-language mobile customer web frontend that fulfill ShopPilot MVP's brow…
when_to_use:
  - EPIC_AUTH
  - EPIC_CATALOG
  - EPIC_INVENTORY
  - EPIC_CART
  - EPIC_CHECKOUT
  - EPIC_ORDER
  - EPIC_PAYMENT
  - EPIC_FRONTEND
not_for:
  - Admin web UI — Frontend Spec is customer-facing only; Admin remains backend-only via Postman/curl in MVP. Admin product create/edit/soft-delete, category create/edit, stock adjustment, and order stat…
  - Promotion/Coupon module (requirement section 8.7, PROMO-001..008) deferred this run; checkout coupon application (CHK-006) therefore deferred; the Frontend Spec coupon row on checkout displays a stat…
  - Shipping Mock module (requirement section 8.11, SHIP-001..006) deferred; no shipment records are created. The order state machine keeps PACKING/SHIPPED/DELIVERED but those transitions are admin-manua…
  - Product Review module (requirement section 8.12, REV-001..005) deferred; product detail's review summary returns empty/placeholder data and the Frontend Spec review route is not built this run.
  - Notification Log module (requirement section 8.14, NOTI-001..004) deferred.
  - Audit Log module (requirement section 8.15, AUD-001..005) deferred; admin actions are not written to a persisted audit log this run.
  - Admin Back Office dashboards (subset of section 8.13: ADM-002 dashboard summary metrics, ADM-005 customer list page, ADM-006 audit log viewer) deferred.
  - Real payment gateway integration (requirement section 4.2 #1).
  - Real shipping provider integration (requirement section 4.2 #2).
  - Multi-vendor marketplace (section 4.2 #3).
  - Loyalty point system (section 4.2 #4).
  - Installment / BNPL (section 4.2 #5).
  - Return / refund flow (section 4.2 #6).
  - Warehouse management (section 4.2 #7).
  - Recommendation engine and AI product recommendation (section 4.2 #8 and #13).
  - Real-time chat (section 4.2 #9).
  - Native mobile app (section 4.2 #10) — the Frontend Spec ships a mobile webview, not a native app.
  - Multi-currency (section 4.2 #11) — Thai Baht only this run.
  - Full multi-language UI (section 4.2 #12) — Thai-only this run; data model must not preclude later i18n.
  - Accounting integration (section 4.2 #14).
  - Density/font/brand-swatch tweak panel from Frontend Spec §7 — exposed only as a designer-time setting, not as a runtime customer feature.
  - Real product photography — Frontend Spec ships striped diagonal placeholders with category labels; replacement is post-launch.
domain: b2c-retail
template_version: 0.1.0
generated_at: 2026-05-08
generator: package_skill.py@v0.1.0
source_run_id: ecom-mvp-2026-05-08-002
terminal_state: ShipWithCaveats
---

# shoppilot-mvp-core-commerce

## Inputs

- identity: Guest registers a CUSTOMER account with email, password, name; email uniqueness enforced server-side; password stored only as a one-way hash (AUTH-001, AUTH-006).
- identity: User logs in with email/password and receives an access token plus refresh token, both with expiry; access token expiry is shorter than refresh token expiry (AUTH-002, AUTH-005).
- identity: Login failure returns a single generic error message and code that does not disclose whether email or password was wrong (AUTH-007).
- identity: Role separation between CUSTOMER and ADMIN; ADMIN-only endpoints reject CUSTOMER tokens with FORBIDDEN (AUTH-003).
- identity: User logout revokes/invalidates the refresh token so subsequent refresh attempts are rejected (AUTH-004).
- identity: Customer can read own profile (name, email, phone, default address) and edit name and phone; email is immutable in MVP (CUST-001, CUST-002).
- identity: Customer can add a shipping address with receiver name, phone, address line, province, district, postal code, and select a default address used by checkout (CUST-003, CUST-004).
- catalog: Guest and Customer can list products, seeing only ACTIVE and visible products (CAT-001).
- catalog: Product listing supports pagination (page, limit, sort) and filters by categoryId, minPrice/maxPrice, and inStock (CAT-002, CAT-003, CAT-004, CAT-005).
- catalog: Product detail returns name, description, images, price, stock status, category, review summary placeholder (CAT-006).
- catalog: Admin creates a product with unique SKU, price greater than zero, and a valid category (CAT-007).
- catalog: Admin edits an existing product; price/status/stock-policy edits are recorded server-side for later audit hookup (CAT-008).
- catalog: Admin soft-deletes a product so that the product still resolves inside historical orders (CAT-009).
- catalog: Product status enum is exactly DRAFT, ACTIVE, INACTIVE, DELETED and listing only exposes ACTIVE (CAT-010).
- catalog: Admin creates a category with name unique within its level and slug unique globally; admin can edit a category and mark a category inactive so it stops appearing in filters (CATE-001, CATE-0…
- inventory: Each SKU tracks availableQty, reservedQty, soldQty as non-negative integers (INV-001, INV-008).
- inventory: Adding to cart does NOT reserve stock; reservation happens only at checkout (INV-002, INV-003).
- inventory: Stock reservations carry an expiry; on expiry the stock is released back to available (INV-004).
- inventory: Payment success converts the reservation: reservedQty decreases and soldQty increases by the same amount; availableQty is unchanged from its post-reserve value (INV-005).
- inventory: Payment failure or timeout releases the reservation: reservedQty decreases, availableQty increases by the same amount (INV-006).
- inventory: Admin adjusts stock and the system writes a stock-adjustment record with delta, reason, and the admin's actorUserId (INV-007).
- inventory: Every stock operation enforces the invariant that no field becomes negative (INV-008).
- cart: Authenticated Customer adds an active product with quantity greater than zero to cart (CART-001).
- cart: Customer changes quantity of a cart item; new quantity must not exceed the SKU's current availableQty (CART-002).
- cart: Customer removes a cart item and the cart subtotal updates accordingly (CART-003, CART-004).
- cart: Adding the same SKU twice merges into one cart line with summed quantity (CART-006).
- cart: Cart flags items whose product became INACTIVE/DELETED so they are visibly non-checkoutable (CART-005).
- cart: Cart subtotal is computed server-side from current product price (or stored snapshot per design) (CART-004).
- cart: Cart row carries updatedAt to allow stale-cart detection downstream (CART-007).
- checkout: Customer submits checkout from a cart that has at least one item (CHK-001).
- (+41 more — see `references/ba-acceptance-criteria.md`)

## Steps

1. [EPIC_AUTH](references/requirement/EPIC_AUTH/EPIC_AUTH.md)
2. [EPIC_CATALOG](references/requirement/EPIC_CATALOG/EPIC_CATALOG.md)
3. [EPIC_INVENTORY](references/requirement/EPIC_INVENTORY/EPIC_INVENTORY.md)
4. [EPIC_CART](references/requirement/EPIC_CART/EPIC_CART.md)
5. [EPIC_CHECKOUT](references/requirement/EPIC_CHECKOUT/EPIC_CHECKOUT.md)
6. [EPIC_ORDER](references/requirement/EPIC_ORDER/EPIC_ORDER.md)
7. [EPIC_PAYMENT](references/requirement/EPIC_PAYMENT/EPIC_PAYMENT.md)
8. [EPIC_FRONTEND](references/requirement/EPIC_FRONTEND/EPIC_FRONTEND.md)

## Acceptance Criteria

See [`references/ba-acceptance-criteria.md`](references/ba-acceptance-criteria.md).

## Architecture

- [Component list](references/architecture/components.json)
- [Integration contracts](references/architecture/contracts.json)
- [Infra summary](references/architecture/infra-summary.md)
- [Infra topology](references/architecture/infra-topology.md)
- [Connectivity](references/architecture/connectivity.md)
- [Observability spec](references/architecture/observability-spec.md)
- [ADRs](references/architecture/ADRs/)

## Per-Component Specs

- **identity** — [td.json](references/components/identity/td.json) · [erd.md](references/components/identity/erd.md)
- **catalog** — [td.json](references/components/catalog/td.json) · [erd.md](references/components/catalog/erd.md)
- **inventory** — [td.json](references/components/inventory/td.json) · [erd.md](references/components/inventory/erd.md)
- **cart** — [td.json](references/components/cart/td.json) · [erd.md](references/components/cart/erd.md)
- **checkout** — [td.json](references/components/checkout/td.json) · [erd.md](references/components/checkout/erd.md)
- **order** — [td.json](references/components/order/td.json) · [erd.md](references/components/order/erd.md)
- **payment** — [td.json](references/components/payment/td.json) · [erd.md](references/components/payment/erd.md)
- **frontend-web** — [td.json](references/components/frontend-web/td.json) · [erd.md](references/components/frontend-web/erd.md)

## API Reference

See [`references/spec/`](references/spec/) for one MD per endpoint, grouped by `<service>-<aggregate>/`.

## Tests

- identity: [qa-l1 csv](tests/qa-l1-identity.csv) · [qa-l1 json](tests/qa-l1-identity.json)
- catalog: [qa-l1 csv](tests/qa-l1-catalog.csv) · [qa-l1 json](tests/qa-l1-catalog.json)
- inventory: [qa-l1 csv](tests/qa-l1-inventory.csv) · [qa-l1 json](tests/qa-l1-inventory.json)
- cart: [qa-l1 csv](tests/qa-l1-cart.csv) · [qa-l1 json](tests/qa-l1-cart.json)
- checkout: [qa-l1 csv](tests/qa-l1-checkout.csv) · [qa-l1 json](tests/qa-l1-checkout.json)
- order: [qa-l1 csv](tests/qa-l1-order.csv) · [qa-l1 json](tests/qa-l1-order.json)
- payment: [qa-l1 csv](tests/qa-l1-payment.csv) · [qa-l1 json](tests/qa-l1-payment.json)
- frontend-web: [qa-l1 csv](tests/qa-l1-frontend-web.csv) · [qa-l1 json](tests/qa-l1-frontend-web.json)
- system: [qa-l2 csv](tests/qa-l2.csv) · [qa-l2 json](tests/qa-l2.json)

## Known Issues

See [`KNOWN_ISSUES.md`](KNOWN_ISSUES.md) for the full caveat list (this run terminated as `ShipWithCaveats`).

## Advisories

See [`ADVISORY.md`](ADVISORY.md).

## Provenance

Generated by squad-flow workflow squad. See `run-log.md` and `run-report.md` under `.squad-run/runs/<run_id>/` for the full execution narrative.
