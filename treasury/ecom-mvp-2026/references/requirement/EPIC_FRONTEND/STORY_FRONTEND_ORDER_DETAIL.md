---
template_version: 0.1.0
story_id: STORY_FRONTEND_ORDER_DETAIL
epic_id: EPIC_FRONTEND
title: orderDetail route shows item snapshots, address snapshot, status timeline, tracking when shipped
as_a: Customer
i_want: a Thai-language orderDetail route showing items, address, status timeline, and tracking
so_that: I can see exactly what I bought, where it ships, and what status it's in
acceptance_criteria:
  - given: An authenticated customer with at least one cart line and one default address
    when: The customer drives the mobile webview through home → search → detail → add-to-cart → cart → checkout → payment-success → orderSuccess → orderDetail
    then: |
      The end-to-end journey completes without error, the orderDetail route shows
      status=PAID, item snapshots, address snapshot, and a status-history timeline
      with at least PENDING_PAYMENT and PAID entries (Frontend Spec §3 + ORD-002,
      ORD-009).
  - given: An order in SHIPPED state with a tracking number
    when: The customer opens orderDetail for that order
    then: |
      The page shows the tracking number, the carrier (or 'mock_express' label), and
      the status badge in Thai for SHIPPED; the status-history timeline shows every
      transition through PENDING_PAYMENT → PAID → PACKING → SHIPPED with timestamps
      (ORD-006, ORD-009).
priority: Must
size: M
edge_cases:
  - case: Customer opens orderDetail for an order whose product was soft-deleted
    expected: Item lines render the snapshot fields (name, price snapshot, image placeholder); the page does not break (CAT-009, ORD-007)
  - case: Customer opens orderDetail for an order belonging to another customer
    expected: NOT_FOUND or FORBIDDEN; the page surfaces a Thai error message via the standard envelope (ORD-001)
requirement_refs: [ORD-001, ORD-002, ORD-006, ORD-007, ORD-009, "Frontend Spec §3"]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created NEW for dry-run #2 from Frontend Spec §3 orderDetail row + §6 timeline.
---

# STORY_FRONTEND_ORDER_DETAIL — orderDetail route shows item snapshots, address snapshot, status timeline, tracking when shipped

## User narrative

As a **Customer**, I want **a Thai-language orderDetail route showing items, address, status timeline, and tracking** so that **I can see exactly what I bought, where it ships, and what status it's in**.

## Why this story exists in EPIC_FRONTEND

The orderDetail route is the end of the customer journey loop. It surfaces the snapshot guarantees (priceSnapshot, addressSnapshot, status history) that ORD-002 + ORD-007 + ORD-009 mandate.

## Edge cases

| Case | Expected |
|---|---|
| Product soft-deleted | Snapshot rendering still works |
| Other customer's order | NOT_FOUND/FORBIDDEN with Thai message |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created NEW for dry-run #2 from Frontend Spec §3 + §6. |
