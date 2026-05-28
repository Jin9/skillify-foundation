---
template_version: 0.1.0
epic_id: EPIC_ORDER
title: Order lifecycle, customer cancel, admin fulfillment transitions
summary: |
  Customers list and read their own orders, cancel a PENDING_PAYMENT order, and view
  status history. Admins drive the fulfillment state machine PAID → PACKING → SHIPPED
  → DELIVERED via backend endpoints (no admin UI). The state machine enforces §9.2
  allowed transitions; forbidden transitions return INVALID_ORDER_STATE. Admin-cancel
  of a PAID order auto-restocks per MVP policy.
business_value: The order is the durable record of the sale. Strict state-machine semantics protect the business from sold-out-then-shipped, double-cancel, and out-of-band status changes.
in_scope_stories:
  - STORY_ORDER_LIST_DETAIL
  - STORY_ORDER_CUSTOMER_CANCEL
  - STORY_ORDER_ADMIN_TRANSITIONS
  - STORY_ORDER_ADMIN_CANCEL
out_of_scope:
  - Shipment record creation (SHIP-001..006) — deferred (§8.11); tracking number stored on the order itself
  - Admin order UI — backend endpoints only this run per Frontend Spec scope
  - Return / refund flow (§4.2 #6) deferred
  - Customer review post-DELIVERED (§8.12) deferred
success_metrics:
  - Forbidden transitions per §9.3 always return 409 INVALID_ORDER_STATE (ORD-003)
  - Customer cannot view another customer's order (ORD-001)
  - Status history captures every transition with from-status, to-status, actor, timestamp (ORD-009)
  - Admin-cancel of PAID auto-restocks (sold→available) per MVP policy
  - p95 order detail latency < 300 ms (PERF-003)
dependencies:
  - EPIC_AUTH
  - EPIC_CHECKOUT
  - EPIC_PAYMENT
  - EPIC_INVENTORY
compliance_sensitive: false
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.9 + §9.2 + §9.3; admin endpoints kept but admin UI removed per Frontend Spec scope.
---

# EPIC_ORDER — Order lifecycle, customer cancel, admin fulfillment transitions

## Why

The order entity is the durable record of every sale and the audit ground-truth for refunds, support, and history. The §9 state machine is what keeps it consistent under concurrent admin/customer actions and replayed payment callbacks.

Three rules anchor this epic:

1. **Allowed transitions only.** §9.2 declares the legal arrows (PENDING_PAYMENT → PAID / PAYMENT_FAILED / PAYMENT_EXPIRED / CANCELLED; PAID → PACKING → SHIPPED → DELIVERED; PAID → CANCELLED with reason). §9.3 hard-bans DELIVERED → CANCELLED, SHIPPED → PAID, etc. ORD-003 maps every forbidden transition to INVALID_ORDER_STATE.

2. **Status history is append-only.** ORD-009: every transition writes a row with from-status, to-status, actor, timestamp. The frontend orderDetail route renders this as a chronological timeline (Frontend Spec §6 `<Timeline>`).

3. **Customer scope is enforced server-side.** ORD-001: customer A's token cannot see customer B's order — neither in listing nor in direct GET. NOT_FOUND or FORBIDDEN, never the order body.

The MVP policy choice: admin-cancel of a PAID order auto-restocks (soldQty → availableQty). Original requirement leaves cancel-after-paid open; we lock this to "release stock" per dry-run #1's decision. Tech-Designer should wire this through the events.order.cancelled handler in inventory.

## Stories in this epic

| Story ID | Title | Priority |
|---|---|---|
| `STORY_ORDER_LIST_DETAIL` | Customer lists own orders and reads order detail with status timeline | Must |
| `STORY_ORDER_CUSTOMER_CANCEL` | Customer cancels own order while PENDING_PAYMENT | Must |
| `STORY_ORDER_ADMIN_TRANSITIONS` | Admin moves PAID → PACKING → SHIPPED (tracking required) → DELIVERED | Must |
| `STORY_ORDER_ADMIN_CANCEL` | Admin cancels PENDING_PAYMENT or PAID with reason; refused after PACKING | Should |

## Out-of-scope rationale

- **Shipment record** — SHIP-001..006 deferred per §8.11; the tracking number and carrier are stored on the order itself rather than a separate shipment aggregate.
- **Admin UI** — Frontend Spec is customer-only; admin endpoints exist but no UI is built this run.
- **Customer reviews** — REV-001..005 deferred per §8.12; the orderDetail route does not surface a "write a review" CTA in MVP.
- **Return / refund** — §4.2 #6 deferred entirely.

## Risks

- **State-machine concurrency** — two admins cancelling the same PAID order concurrently must not double-restock. Tech-Designer must use a row-level guard (status comparison or version) on the transition.
- **§9.3 forbidden transitions** — Reviewer-L1 should diff the implemented transition table against §9.3 line by line; missing one (e.g. PAYMENT_EXPIRED → PAID) is a P1 correctness bug.
- **Customer-vs-admin cancel boundary** — customer can only cancel from PENDING_PAYMENT; admin can additionally cancel from PAID. ORD-004 vs ORD-005 must produce distinct error codes when used wrong.
- **Tracking number requirement on SHIPPED** — ORD-006 requires the admin SHIP transition to carry a tracking string; missing tracking returns VALIDATION_ERROR.

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.9 + §9.2 + §9.3; admin UI removed per Frontend Spec scope. |
