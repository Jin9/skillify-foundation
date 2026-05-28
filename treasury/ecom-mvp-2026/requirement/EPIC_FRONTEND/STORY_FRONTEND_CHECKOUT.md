---
template_version: 0.1.0
story_id: STORY_FRONTEND_CHECKOUT
epic_id: EPIC_FRONTEND
title: checkout route shows address + items + server-computed totals (incl. §11.2) and place-order CTA with Idempotency-Key
as_a: Customer
i_want: a Thai-language checkout route that shows my address, items, and server-priced totals
so_that: I can review and confirm the purchase with confidence that totals are server-priced
acceptance_criteria:
  - given: An authenticated customer whose cart subtotal is 1500 THB on the checkout route
    when: The checkout route renders the totals block
    then: |
      shippingFee shows ฿0 and grandTotal shows ฿1,500 (Frontend Spec §5 + §11.2).
  - given: An authenticated customer whose cart subtotal is 1499 THB on the checkout route
    when: The checkout route renders the totals block
    then: |
      shippingFee shows ฿60 and grandTotal shows ฿1,559 (Frontend Spec §5 + §11.2).
  - given: An authenticated customer on the checkout route with a non-empty cart and a default address
    when: The customer taps the 'place order' CTA
    then: |
      The frontend POSTs to checkout commit with an Idempotency-Key header generated
      client-side; on success the customer is routed to the payment route with the
      created order's id; on a network retry the same Idempotency-Key is reused so
      the server returns the original order (CHK-009).
priority: Must
size: L
edge_cases:
  - case: Customer changes the cart between checkout preview and commit such that subtotal crosses the §11.2 threshold
    expected: The commit endpoint recomputes shippingFee from the cart at commit time; the previously-shown preview is advisory only
  - case: Place-order CTA is double-tapped (rapid re-tap before the response returns)
    expected: The frontend reuses the same Idempotency-Key for both attempts; the server returns the same order id; the customer is routed to the payment route exactly once (CHK-009)
  - case: Place-order returns 409 INSUFFICIENT_STOCK
    expected: The frontend surfaces the Thai error message and lists the short items per the response detail (CHK-004)
requirement_refs: [CHK-001, CHK-002, CHK-005, CHK-009, "§11.2", "Frontend Spec §3", "Frontend Spec §5"]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created NEW for dry-run #2; explicit Idempotency-Key on the place-order CTA so client retries are safe.
---

# STORY_FRONTEND_CHECKOUT — checkout route shows address + items + server-computed totals (incl. §11.2) and place-order CTA with Idempotency-Key

## User narrative

As a **Customer**, I want **a Thai-language checkout route that shows my address, items, and server-priced totals** so that **I can review and confirm the purchase with confidence that totals are server-priced**.

## Why this story exists in EPIC_FRONTEND

The checkout route is the customer's last review surface before payment. The §11.2 visual rendering and the Idempotency-Key generation on the place-order CTA are the two non-negotiable behaviours.

## Edge cases

| Case | Expected |
|---|---|
| Cart changes between preview and commit | Commit recomputes from current cart |
| Double-tap place-order | Same Idempotency-Key; same order returned |
| INSUFFICIENT_STOCK on commit | Thai error with per-item details |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created NEW for dry-run #2; explicit Idempotency-Key behaviour on the place-order CTA. |
