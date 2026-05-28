---
template_version: 0.1.0
epic_id: EPIC_CART
title: Customer cart management
summary: |
  Authenticated customers add active products to a cart, edit quantities respecting
  availableQty, remove lines, and have duplicate adds merged into one line with a
  summed quantity. The cart computes a server-side subtotal and exposes updatedAt for
  stale-cart detection. Adding to cart does NOT reserve stock — reservation happens at
  checkout.
business_value: Bridge between catalog browsing and checkout. Without a working cart, the customer cannot accumulate purchase intent, and the checkout route has no input.
in_scope_stories:
  - STORY_CART_ADD
  - STORY_CART_UPDATE_QTY
  - STORY_CART_REMOVE
out_of_scope:
  - Saved-for-later / wishlist surface — not in original requirement
  - Cross-device cart sync beyond authenticated single-cart-per-user — single ACTIVE cart per customer per §12.1
  - Cart-level promotion preview — coupon module deferred (PROMO-001..008 per §8.7)
success_metrics:
  - Adding the same SKU twice produces ONE line with merged quantity (CART-006)
  - Quantity edits exceeding availableQty are rejected with a validation error citing availableQty (CART-002)
  - INACTIVE/DELETED products in cart are flagged non-checkoutable (CART-005)
  - Cart subtotal recomputes server-side on every read; client subtotals are never trusted (CART-004)
dependencies:
  - EPIC_AUTH
  - EPIC_CATALOG
  - EPIC_INVENTORY
compliance_sensitive: false
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.6; tightened "INACTIVE-flagging" AC to the cart-tab visual contract from Frontend Spec §3.
---

# EPIC_CART — Customer cart management

## Why

The cart tab (Frontend Spec §3) is the staging ground between catalog browse and checkout. The customer accumulates purchase intent here and reviews the server-priced subtotal before committing to a checkout flow.

The original requirement (§8.6) specifies seven cart rules:
- CART-001: add an active product with quantity > 0 (auth required)
- CART-002: edit quantity capped at current availableQty
- CART-003: remove a line and recalc subtotal
- CART-004: subtotal computed server-side from current price (or stored snapshot per design)
- CART-005: cart flags items whose product became INACTIVE/DELETED
- CART-006: duplicate SKU adds merge into one line with summed quantity
- CART-007: cart row carries updatedAt for stale-cart detection

The architectural commitment that **add-to-cart does NOT reserve stock** (INV-002 / INV-003) is what allows the cart to be cheap and the checkout to be the consistency point. Reservations are owned by the checkout/inventory boundary, not the cart.

## Stories in this epic

| Story ID | Title | Priority |
|---|---|---|
| `STORY_CART_ADD` | Customer adds an active product to cart (quantity merge included) | Must |
| `STORY_CART_UPDATE_QTY` | Customer updates quantity of a cart line | Must |
| `STORY_CART_REMOVE` | Customer removes a cart line | Must |

## Out-of-scope rationale

- **Wishlist / saved-for-later** — not in §8.6 or §4.1; would be invented behaviour.
- **Coupon preview on cart** — coupon module deferred (§8.7 deferred); checkout-side coupon (CHK-006) also deferred this run.
- **Cross-device cart merge** — §12.1 declares one ACTIVE cart per user; no requirement to merge anonymous carts on login.

## Risks

- **Quantity-cap race** — CART-002 caps quantity at availableQty *at the moment of edit*; the cap can be stale by the time checkout reserves. The cart-side check is best-effort; the authoritative check lives in checkout (CHK-004).
- **Subtotal price source** — CART-004 allows either "current price" or "stored snapshot" depending on design. Tech-Designer chooses; BA does not prescribe. Whatever is chosen, ORD-007 (snapshot at order creation) must still hold so historical orders are immutable.
- **Stale-cart detection** — `updatedAt` (CART-007) is a hint, not a lock. Concurrent edits from two tabs need an ordering rule the Tech-Designer owns.

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.6; tightened "INACTIVE-flagging" AC to the cart-tab visual contract. |
