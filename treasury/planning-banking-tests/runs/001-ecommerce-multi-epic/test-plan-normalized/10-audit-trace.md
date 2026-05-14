---
artifact_type: qa-audit-trace
test_plan_id: TPLAN-EPIC-SHOPPILOT-MVP
---

# Audit Trace

## Processing Metadata

```json
{
  "ba_brief_idempotency_key": "e5f8b9c2-3d4a-4e6f-8b9c-1d2e3f4a5b6c",
  "ba_brief_version": "1.2.1",
  "environment_inventory_consumed": false,
  "pyramid_allocation_rules_version": "1.0",
  "role_boundaries": [
    {
      "notes": "v1.0.0 monolithic; role doc in references/v1.1-role-boundaries.md",
      "owns_sections": [
        "strategy",
        "signoff_criteria",
        "high_risk_areas"
      ],
      "role": "qa-strategist"
    },
    {
      "notes": "Per-Story SDET in v1.1+ fan-out model",
      "owns_sections": [
        "test_cases",
        "execution_dag"
      ],
      "role": "test-designer"
    },
    {
      "notes": "v1.1+ Performance + Security + Accessibility leads",
      "owns_sections": [
        "nfr_tests"
      ],
      "role": "nfr-engineer"
    },
    {
      "notes": "v1.1+ Compliance/Regulatory lead per regulator",
      "owns_sections": [
        "compliance_tests"
      ],
      "role": "compliance-mapper"
    },
    {
      "notes": "v1.1+ Test Data Engineer",
      "owns_sections": [
        "test_data_specs"
      ],
      "role": "fixture-architect"
    },
    {
      "notes": "v1.1+ Release Manager",
      "owns_sections": [
        "signoff_criteria",
        "qa_readiness_checklist"
      ],
      "role": "signoff-officer"
    },
    {
      "notes": "v1.1+ runs last; merges fragments",
      "owns_sections": [
        "coverage_matrix",
        "coverage_gaps"
      ],
      "role": "coverage-auditor"
    }
  ],
  "scenario_mapping_table_version": "1.0",
  "tier_decisions": [
    {
      "rationale": "Inherited from BA initiative.per_epic_tier",
      "target": "EPIC-CART-CHECKOUT",
      "tier": "T2"
    },
    {
      "rationale": "Inherited from BA initiative.per_epic_tier",
      "target": "EPIC-CATALOG",
      "tier": "T2"
    },
    {
      "rationale": "Inherited from BA initiative.per_epic_tier",
      "target": "EPIC-GOVERNANCE-OBS",
      "tier": "T2"
    },
    {
      "rationale": "Inherited from BA initiative.per_epic_tier",
      "target": "EPIC-IDENTITY",
      "tier": "T2"
    },
    {
      "rationale": "Inherited from BA initiative.per_epic_tier",
      "target": "EPIC-ORDER-FULFILL",
      "tier": "T2"
    }
  ],
  "tl_design_consumed": false
}
```

## QA Readiness Checklist

- all_banking_grade_applies_have_tests: True
- all_compliance_regulators_covered: False
- all_dependencies_have_test_envs: True
- all_nfr_targets_resolved: False
- all_stories_have_test_cases: True
- no_blocking_ba_governance_gaps: False
