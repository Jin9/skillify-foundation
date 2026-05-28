---
template_version: 0.1.0
epic_id: EPIC_INVENTORY
title: Stock invariant, reservation, sweep, admin adjustment
summary: |
  Each SKU tracks availableQty / reservedQty / soldQty as non-negative integers.
  Add-to-cart does NOT reserve; reservation happens at checkout, with an expiry that
  the sweeper releases on timeout. Payment success converts reserved→sold; payment
  failure or timeout releases reserved back to available. Admin stock adjustments
  write a record with delta, reason, and actor.
business_value: Inventory is the consistency anchor between cart, checkout, payment, and order. Without a strict reservation model, ShopPilot oversells or undersells; with it, every ฿ of stock movement is traceable.
in_scope_stories:
  - STORY_INVENTORY_RESERVE_AND_CONVERT
  - STORY_INVENTORY_RELEASE_ON_FAILURE
  - STORY_INVENTORY_ADMIN_ADJUST
  - STORY_INVENTORY_SWEEP_EXPIRED
out_of_scope:
  - Multi-warehouse / multi-location stock (§4.2 #7) deferred
  - SKU-level reservation queueing (back-order) — not in §8.5
  - External WMS integration deferred
success_metrics:
  - availableQty / reservedQty / soldQty never become negative under any operation (INV-008)
  - Payment success: reserved → sold conversion is exactly equal in size (INV-005)
  - Reservation expiry triggers release within reservation TTL (INV-004)
  - Stock-adjustment record exists with delta, reason, actorUserId for every admin adjust (INV-007)
dependencies:
  - EPIC_AUTH
  - EPIC_CATALOG
compliance_sensitive: false
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.5 + §10; sweeper story added explicitly so the PAYMENT_EXPIRED edge case (browser-walk-away) has an owner.
---

# EPIC_INVENTORY — Stock invariant, reservation, sweep, admin adjustment

## Why

Inventory is the system's consistency anchor. The cart, checkout, payment, and order epics all transitively depend on the reservation model defined here.

The §10 invariant is straightforward: `availableQty >= 0 && reservedQty >= 0 && soldQty >= 0`. The §10.3 example shows the four states a SKU walks through during the happy path:

```
before:    available=10  reserved=0  sold=0
checkout:  available=8   reserved=2  sold=0   (CHK-008)
paid:      available=8   reserved=0  sold=2   (INV-005)
failed:    available=10  reserved=0  sold=0   (INV-006)
```

Two architectural commitments make this work:

1. **Add-to-cart does NOT reserve.** INV-002 / INV-003: a SKU is only locked at checkout, never at cart-add. This keeps cart cheap and avoids zombie reservations.

2. **Reservations have a TTL.** INV-004: if a customer commits checkout and never simulates payment (browser closed, walked away), the sweeper transitions the order to PAYMENT_EXPIRED and releases the reservation. This is what protects the catalog from being silently locked out by abandoned checkouts.

Admin stock adjustment (INV-007) writes a record with delta, reason, and actorUserId — the audit module is deferred (§8.15), but the adjustment record itself is in scope so the data is available when audit ships.

## Stories in this epic

| Story ID | Title | Priority |
|---|---|---|
| `STORY_INVENTORY_RESERVE_AND_CONVERT` | Checkout reserves; payment success converts reserved → sold | Must |
| `STORY_INVENTORY_RELEASE_ON_FAILURE` | Payment failed/timeout/customer-cancel releases reserved → available | Must |
| `STORY_INVENTORY_ADMIN_ADJUST` | Admin adjusts stock; record persists delta, reason, actorUserId | Must |
| `STORY_INVENTORY_SWEEP_EXPIRED` | Sweeper expires reservations past TTL and releases stock | Must |

## Out-of-scope rationale

- **Multi-warehouse** — §4.2 #7 deferred; one logical stock per SKU in MVP.
- **Back-order / preorder** — not mentioned in §8.5; we do not invent it.
- **External WMS integration** — §4.2 #7 deferred.

## Risks

- **Concurrent reservation race** — two checkouts on the last unit must serialise; Tech-Designer must use row-level locking or a CAS pattern. INV-008 means availableQty cannot go negative under ANY interleaving.
- **Sweeper restartability** — the sweeper that expires reservations must be safe to crash and resume; double-release is an INV-008 violation.
- **Admin adjust + active reservation** — adjusting stock down while reservations are live could push availableQty negative. Tech-Designer must decide whether to fail the adjust or to allow it but accept that subsequent checkouts may see INSUFFICIENT_STOCK earlier.
- **Admin-cancel-of-PAID auto-restock** — the MVP policy of soldQty → availableQty on admin cancel of a PAID order must be wired through the events.order.cancelled handler in inventory; otherwise the stock leaks.

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.5 + §10; sweeper story added explicitly. |
