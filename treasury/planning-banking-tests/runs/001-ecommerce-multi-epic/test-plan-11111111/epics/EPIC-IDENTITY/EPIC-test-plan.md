---
artifact_type: qa-epic-test-plan
epic_id: EPIC-IDENTITY
story_count: 3
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
test_tier: T2
---

# Epic Test Plan — EPIC-IDENTITY

- Tier: T2
- Stories: 3

## Critical Paths

- signup -> password-hash storage -> auto-login or redirect (EPIC-IDENTITY-1)
- login with valid creds -> session token issued -> account home (EPIC-IDENTITY-2)
- login with wrong creds / unknown email -> identical generic error (EPIC-IDENTITY-2)
- address create/edit/default/delete -> default invariant + active-order block (EPIC-IDENTITY-3)
- cross-customer authz denial on address edit (EPIC-IDENTITY-3)

## Stakeholder Owners

- Customer (account-holder)
- Merchant Admin
- Data Protection Officer (DPO) / Privacy [absent - P1 governance gap]
- Legal Counsel [absent - P1 governance gap]
- Security Reviewer [absent - P1 governance gap]

## Stories

- EPIC-IDENTITY-1 — coverage: partial, tests: 4
- EPIC-IDENTITY-2 — coverage: partial, tests: 5
- EPIC-IDENTITY-3 — coverage: partial, tests: 5
