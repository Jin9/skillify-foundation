---
template_version: 0.1.0
epic_id: EPIC_CHECKOUT
title: Server-priced checkout with reservation and idempotency
summary: |
  Customers convert a cart into an order with status PENDING_PAYMENT. Checkout
  validates the address ownership, rejects INACTIVE/DELETED items, validates stock
  per item, recomputes price/shipping server-side under the §11.2 free-shipping ≥ ฿1,500
  rule, reserves stock, and accepts an Idempotency-Key so duplicate POSTs return the
  same order. Coupon application is out of scope this run.
business_value: The checkout endpoint is the consistency point of the entire purchase. Server-priced totals + idempotent commit are what keep ShopPilot from refunds and from negative inventory.
in_scope_stories:
  - STORY_CHECKOUT_PREVIEW
  - STORY_CHECKOUT_COMMIT
  - STORY_CHECKOUT_SHIPPING_FEE
out_of_scope:
  - Coupon application at checkout (CHK-006) — coupon module deferred (§8.7)
  - Multi-address split shipments — single shipping address per order in MVP
  - Saved payment methods — payment intent is created fresh per order
  - Tax calculation — not in §11 pricing rules; THB only, no VAT line in MVP
success_metrics:
  - Two checkout POSTs with the same Idempotency-Key create exactly one order and reserve stock exactly once (CHK-009)
  - Server-computed grandTotal ignores any client-supplied total (CHK-005)
  - §11.2 boundary holds: subtotal=1499 ⇒ shippingFee=60; subtotal=1500 ⇒ shippingFee=0
  - INSUFFICIENT_STOCK returns per-item availableQty/requestedQty without creating the order (CHK-004)
  - p95 checkout latency < 500 ms on local Docker Compose (PERF-002)
dependencies:
  - EPIC_AUTH
  - EPIC_CART
  - EPIC_CATALOG
  - EPIC_INVENTORY
compliance_sensitive: false
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.8 + §11.2; coupon (CHK-006) explicitly deferred and noted on the Frontend Spec checkout coupon row.
---

# EPIC_CHECKOUT — Server-priced checkout with reservation and idempotency

## Why

The checkout endpoint is the point of no return: it locks stock, freezes price, snapshots the address, and creates the order that the payment service will resolve. Every behavior in §8.8 of the requirement exists to keep this transition consistent.

The two non-negotiable rules:

1. **Server-priced everything.** CHK-005: client-supplied totals are ignored; the server recomputes subtotal, shipping fee, and grandTotal from the authoritative product price and the §11.2 free-shipping rule (subtotal ≥ ฿1,500 ⇒ shippingFee=0; otherwise shippingFee=60). The frontend's checkout route renders only what the server returns.

2. **Idempotency on commit.** CHK-009: the customer's mobile webview POSTs the place-order action with an Idempotency-Key header. A duplicate POST with the same key must return the originally created order rather than creating a second one or reserving stock twice. This is what lets the frontend retry safely on flaky network.

Coupon application (CHK-006) is the one CHK-* rule we are NOT shipping this run, because the coupon module (§8.7) is deferred. The Frontend Spec checkout route shows a coupon row; in MVP it renders as a static disabled placeholder.

## Stories in this epic

| Story ID | Title | Priority |
|---|---|---|
| `STORY_CHECKOUT_PREVIEW` | Customer previews a checkout (server-priced totals before commit) | Must |
| `STORY_CHECKOUT_COMMIT` | Customer commits checkout (idempotent, reserves stock, creates order) | Must |
| `STORY_CHECKOUT_SHIPPING_FEE` | §11.2 shipping-fee tier applied server-side at preview and commit | Must |

## Out-of-scope rationale

- **Coupon application** — coupon module deferred per §8.7; CHK-006 not built this run.
- **Tax line** — §11 pricing rules don't define tax for MVP; THB only, no VAT line.
- **Split shipments** — single address per order; multi-address is a v2 problem.
- **Saved payment methods** — every checkout creates a fresh payment intent; nothing carried across orders.

## Risks

- **§11.2 boundary off-by-one** — subtotal exactly = 1500 must apply shippingFee=0 (the inclusive ≥-1500 branch). subtotal=1499 must apply shippingFee=60. Reviewer-L1 should test both sides.
- **Idempotency-Key replay with mismatched payload** — CHK-009 must return the originally-created order *without* silently honoring the new payload. Otherwise an attacker could mutate an order under a stolen idempotency key.
- **Address ownership leak** — CHK-002: checkout must reject an addressId belonging to a different customer with FORBIDDEN/VALIDATION_ERROR. Failing this is a BOLA-class issue.
- **Cart-vs-checkout INSUFFICIENT_STOCK** — checkout's stock validation is authoritative; the cart's CART-002 cap is best-effort. The error envelope on CHK-004 must list every short item, not just the first.

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.8 + §11.2; coupon (CHK-006) explicitly deferred. |
