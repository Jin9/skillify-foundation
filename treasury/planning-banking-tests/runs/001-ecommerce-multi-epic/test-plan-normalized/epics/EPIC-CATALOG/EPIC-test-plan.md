---
artifact_type: qa-epic-test-plan
epic_id: EPIC-CATALOG
story_count: 3
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
test_tier: T2
---

# Epic Test Plan — EPIC-CATALOG

- Tier: T2
- Stories: 3

## Critical Paths

- Guest listing returns only active products with valid pagination
- PDP suppresses add-to-cart affordance for non-sellable products
- Admin product edits emit audit events and remain idempotent under replay

## Stakeholder Owners

- Merchant Admin (back-office operator)
- Product Manager (business-side spec author)
- Guest / Customer (read-side consumer)

## Stories

- EPIC-CATALOG-1 — coverage: partial, tests: 5
- EPIC-CATALOG-2 — coverage: complete, tests: 4
- EPIC-CATALOG-3 — coverage: partial, tests: 7
