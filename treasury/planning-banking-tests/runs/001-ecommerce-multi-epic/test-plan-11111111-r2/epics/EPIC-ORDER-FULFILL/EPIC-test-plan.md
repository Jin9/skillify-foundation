---
artifact_type: qa-epic-test-plan
epic_id: EPIC-ORDER-FULFILL
story_count: 3
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
test_tier: T2
---

# Epic Test Plan — EPIC-ORDER-FULFILL

- Tier: T2
- Stories: 3

## Critical Paths

- Order state machine: paid -> packing -> shipped -> delivered (forward-only, irreversible)
- Admin-cancel before packing with stock restoration and manual-refund OOB notification
- Shipped transition mandatory tracking_number capture and customer-facing tracking display
- Customer order-detail snapshot integrity (price/address immutability per 6.8)
- Cross-tenant authz: customer A cannot read customer B's order
- Audit emission on every status flip with 7 mandatory fields (actor, action, before, after, ts, idem_key, subject_id)

## Stakeholder Owners

- qa-lead@merchant
- ops-admin@merchant
- compliance@merchant
- tl-fulfilment@merchant

## Stories

- EPIC-ORDER-FULFILL-1 — coverage: partial, tests: 12
- EPIC-ORDER-FULFILL-2 — coverage: partial, tests: 8
- EPIC-ORDER-FULFILL-3 — coverage: complete, tests: 6
