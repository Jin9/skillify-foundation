---
artifact_type: qa-strategy
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
---

# Test Strategy

## Pyramid Allocation

{
  "contract_pct": 15,
  "e2e_pct": 5,
  "integration_pct": 30,
  "unit_pct": 50
}


## Tier Rationale

T2 baseline across all 5 epics per BA initiative.per_epic_tier (all epics tagged T2). Customer-facing PII surfaces and monetary state changes justify T2 pyramid weighting. Integration tests carry the bulk of banking-grade coverage (idempotency replay, audit emission, authz cross-tenant) since unit tests cannot validate cross-component contracts.


## Mock vs Real Strategy

Mock payment provider (mock-PSP) and mock shipping provider for v1.0.0 per BA scope. Real postgres + redis in CI ephemeral env. Contract tests carry the contract-shape to a future real PSP. Race-condition tests require deterministic concurrency primitive with clock control + lock instrumentation.


## High Risk Areas

- Audit log emission on every state-change (cross-cutting, GOVERNANCE-OBS owns)
- Compliance scope blocked by 4 BA P1 governance gaps — cannot ship QA execution until resolved
- Cross-tenant authz on order read paths
- Order state machine irreversibility (paid->shipped, delivered terminal)
- PII redaction in logs and notifications (PDPA P1 governance gap blocks scope)
- Payment idempotency under retry (CART-CHECKOUT STORY-3, STORY-4)
- Stock reservation race conditions on last-unit checkout
