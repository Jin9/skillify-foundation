---
artifact_type: qa-signoff-criteria
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
---

# Sign-off Criteria

## Tier: T2

## Go Conditions

- All banking_grade_* tests pass at 100%.
- All compliance tests pass at 100% (PDPA, TRC, CCA, CPA).
- Performance NFR targets met at p95.
- No open P0 or P1 defects.
- BA P1 governance gaps resolved (4 currently blocking).

## Conditional-Go Conditions

- Conditional-go allowed with documented exception per P1 defect (T2 rule).
- Compliance test scope unblocked once BA citation_status moves from 'pending' to 'confirmed' per regulator.

## No-Go Conditions

- Any BA P1 governance gap unresolved.
- Any banking_grade_* test failing.
- Any P0 defect open.
- Tipping-off violation detected in any test description or product copy.
