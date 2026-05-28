---
template_version: 0.1.0
story_id: STORY_ORDER_LIST_DETAIL
epic_id: EPIC_ORDER
title: Customer lists own orders and reads order detail with status timeline
as_a: Customer
i_want: to list my own orders and read the detail of any one of them
so_that: I can track what I've bought, see status, and confirm price snapshots
acceptance_criteria:
  - given: Customer A's order O_A and Customer B's access token
    when: Customer B GETs orders listing or GETs order O_A by id
    then: |
      Listing does not include O_A and direct GET returns NOT_FOUND or FORBIDDEN;
      never the order body (ORD-001).
  - given: An order moved through PENDING_PAYMENT → PAID → PACKING
    when: Customer GETs order detail
    then: |
      Response includes item lines with priceSnapshot, addressSnapshot, currentStatus
      = PACKING, and status history with three rows in chronological order each
      carrying actor and timestamp (ORD-002, ORD-007, ORD-009).
  - given: An order whose product was soft-deleted by admin AFTER order creation
    when: Customer GETs order detail
    then: |
      The order detail still resolves the product's name, image, and priceSnapshot
      from the order item snapshot rather than from the live catalog (CAT-009,
      ORD-007).
priority: Must
size: M
edge_cases:
  - case: Customer GETs orders with status filter (active vs done)
    expected: Listing respects the filter; status-to-bucket mapping is documented
  - case: Customer GETs an order id that does not exist
    expected: NOT_FOUND
requirement_refs: [ORD-001, ORD-002, ORD-007, ORD-009, CAT-009]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §8.9 ORD-001..002, ORD-007, ORD-009; combined list+detail since they share auth scope.
---

# STORY_ORDER_LIST_DETAIL — Customer lists own orders and reads order detail with status timeline

## User narrative

As a **Customer**, I want **to list my own orders and read the detail of any one of them** so that **I can track what I've bought, see status, and confirm price snapshots**.

## Why this story exists in EPIC_ORDER

The mobile webview's orders tab and orderDetail route (Frontend Spec §3) are built on these endpoints. ORD-001 (customer-scope enforcement) and CAT-009 + ORD-007 (snapshot resolution after soft-delete) are the two critical correctness rules.

## Edge cases

| Case | Expected |
|---|---|
| Status-bucket filter (active/done) | Documented mapping respected |
| Unknown order id | NOT_FOUND |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §8.9. |
