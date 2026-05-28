---
template_version: 0.1.0
story_id: STORY_CHECKOUT_SHIPPING_FEE
epic_id: EPIC_CHECKOUT
title: §11.2 shipping-fee tier applied server-side at preview and commit
as_a: Customer
i_want: free shipping when my subtotal is at least ฿1,500 and a flat ฿60 fee otherwise
so_that: I see a transparent and consistent shipping fee at preview and commit
acceptance_criteria:
  - given: A customer's cart with subtotal of 1499 THB
    when: checkout preview is called
    then: |
      The response sets shippingFee=60 and grandTotal=1559 (§11.2).
  - given: A customer's cart with subtotal of 1500 THB
    when: checkout preview is called
    then: |
      The response sets shippingFee=0 and grandTotal=1500 (§11.2 threshold; subtotal
      ≥ 1500 ⇒ free shipping).
priority: Must
size: S
edge_cases:
  - case: Subtotal lands exactly on the threshold (1500 THB)
    expected: Server applies shippingFee=0 (the ≥ 1500 branch is inclusive); subtotal=1499 instead applies shippingFee=60 (§11.2)
  - case: Subtotal is much greater than 1500 (e.g. 9999 THB)
    expected: shippingFee=0 regardless of how large; the rule is binary, not progressive
  - case: Customer changes the cart between preview and commit such that subtotal crosses the threshold
    expected: Commit recomputes shippingFee from the cart at commit time; the preview's shippingFee is advisory and not binding
requirement_refs: ["§11.2"]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - date: 2026-05-08
    author: BA (Claude Opus 4.7 1M, dry-run #2)
    change: Created from requirement §11.2 MVP shipping fee rule; carries forward from runs/01 with explicit boundary edge cases.
---

# STORY_CHECKOUT_SHIPPING_FEE — §11.2 shipping-fee tier applied server-side at preview and commit

## User narrative

As a **Customer**, I want **free shipping when my subtotal is at least ฿1,500 and a flat ฿60 fee otherwise** so that **I see a transparent and consistent shipping fee at preview and commit**.

## Why this story exists in EPIC_CHECKOUT

The §11.2 free-shipping rule is the only pricing rule in scope this run (coupon module deferred). Both the preview AND the commit must apply it server-side; the frontend renders only what the server returns.

## Edge cases

| Case | Expected |
|---|---|
| Subtotal = 1500 (boundary inclusive) | shippingFee=0 |
| Subtotal much greater than 1500 | shippingFee=0 (binary, not progressive) |
| Cart changes between preview and commit | Commit recomputes from current cart; preview is advisory |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | BA (Claude Opus 4.7 1M, dry-run #2) | Created from requirement §11.2. |
