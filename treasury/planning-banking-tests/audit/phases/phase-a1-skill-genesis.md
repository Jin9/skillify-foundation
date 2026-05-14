---
artifact_type: skill-audit-phase
skill_ref: qa-plan-from-brief
phase: A1
phase_name: skill-genesis
version_documented: 1.0.0
created_at: 2026-05-12
research_session_plan: /Users/IF640063/.claude/plans/let-s-plan-to-research-cryptic-zephyr.md
holdout_used: /Users/IF640063/Desktop/example/qa-bootstrap-kit/e-commerce-v5/output-e5f8b9c2/
---

# Phase A1 — Skill Genesis

## 1. Purpose

Document the design decisions that produced `qa-plan-from-brief` v1.0.0 so that any future auditor, contributor, or v1.1+ engineer can reconstruct *why* the skill has its current shape from the inputs that drove the design.

## 2. Locked decisions

| # | Decision | Rationale | Evidence |
|---|---|---|---|
| D1 | **Architecture: monolithic** (one LLM call per skill invocation) | ADR #8 in `COGNITIVE_OS_PROJECT.md` §2 forbids internal fan-out at v1; reinforced by §1 line 81 and §13 non-goal line 879 | `references/v1.1-role-boundaries.md` §1, §8 |
| D2 | **T2 default tier** | Stage 4c default per `DELIVERY_WORKFLOW_PLAN.md` v2.2 Per-Stage Model Selection Matrix | DELIVERY_WORKFLOW_PLAN line ~132 |
| D3 | **Hybrid input contract** (`oneOf: ba_brief_ref XOR ba_brief_inline`) | Path-based for CLI / test harness; inline for workflow-engine in-memory context per COGNITIVE_OS §5.2 hybrid pattern | `schemas/input.json` |
| D4 | **Role-agnostic canonical output schema** | Decouples canonical artifact from orchestration choice; allows v1.1+ fan-out without breaking schema | `schemas/output.json`, `processing_metadata.role_boundaries[]` |
| D5 | **Cross-document invariants documented as `$comment`** in `allOf` | Draft-07 cannot encode "for every BA array element exists a matching QA array element"; renderer enforces at output time | `schemas/output.json` allOf section |
| D6 | **Pattern over format for IDs** (TC, NFR, COMP) | Draft-07 `format` is annotation-only by default; `pattern` is validator-agnostic | `schemas/output.json` test_cases.items.properties.id |
| D7 | **Deterministic Python renderer**, no LLM at render time | Banking-grade non-negotiable #3 (determinism); enables byte-identical re-runs | `scripts/render_test_plan_tree.py` |
| D8 | **5 conditional roles all fire on holdout** but treated as "facets the single LLM must cover" in v1.0.0 | Security, Performance, Compliance, Accessibility all fire on `output-e5f8b9c2/`; role boundaries preserved for v1.1+ | `references/v1.1-role-boundaries.md` §5 |

## 3. Research evidence

The research phase produced three load-bearing findings:

1. **Architectural compatibility (Step 1)** — HIGH confidence verdict that internal sub-agent fan-out is forbidden in v1 by the cognitive-OS architecture. Source quotes preserved in `references/v1.1-role-boundaries.md` §1 and §8.

2. **Input-locality matrix (Step 2)** — 10 candidate roles mapped to 13 BA-output keys; all 4 conditional roles (Security, Performance, Compliance, Accessibility) fire on the holdout. Preserved in `references/v1.1-role-boundaries.md` §3, §4.

3. **Token-cost measurement (Step 3)** — only Monolithic ($1.15/run) clears T2 envelope of <$1.50/run. Fan-out variants $2.72–$7.55. Preserved in `references/v1.1-role-boundaries.md` §7.

## 4. Re-baselined acceptance numbers

The original design brief's §10.3 #7 acceptance numbers were measured against an incorrect holdout. Re-baselined against `output-e5f8b9c2/`:

| Metric | Original brief | Actual holdout |
|---|---|---|
| Epic plans | 10 | **5** |
| Story plans | 29 | **16** |
| Test cases | ≥100 | **≥80** |
| NFR tests | ≥30 | **≥20** |
| Compliance tests | ≥40 | **≥25** |
| Coverage gaps | ≥20 | **≥15** (tied to 51 OQs + 4 P1 governance gaps) |
| `blocks_qa_execution` | true | true (holdout brief is `blocked_partial_brief`) |

## 5. Open questions deferred to v1.0.0 implementation runtime

1. **"Byte-identical" semantics** — strict text-equal vs canonicalized. The skill's determinism validation (test scenario 1 in `tests/assertions/`) uses canonicalized comparison (sorted keys + masked volatile fields).
2. **Response-cache substrate** — not built; the kit has no LLM response cache yet. Re-runs at T2 temp=0.3 carry provider-side variance. v1.0.0 mitigates by canonicalization in the validation harness, not by caching.
3. **Schema vocabulary completeness** — `compliance_tags` is left open; `nfr_type` enum may need extension; `regulator_code` pattern is open. Tracked as v1.0.x patch candidates.

## 6. Risks tracked into post-ship

- **R1 — temp=0.3 non-determinism without seed pinning.** Test scenario 1 will fail intermittently if the provider returns different tokens. Mitigation: add a response-cache wrapper if observed.
- **R2 — Prompt-caching savings unmodeled.** Irrelevant under monolithic but flips v1.1+ economics in fan-out's favor. Re-measure before v1.1+ planning.
- **R3 — Steps 4–5 manual A/B was skipped under monolithic-locked path.** If v1.1+ re-opens decomposition, that A/B must run for evidence — currently we have cost + ADR arguments only.
- **R4 — Holdout overfitting.** Single T2 e-commerce brief drove all decisions. If applied to T1 banking workloads, re-validate scenario 4 thresholds before sign-off.

## 7. References

- `/Users/IF640063/.claude/plans/let-s-plan-to-research-cryptic-zephyr.md` — full research session record
- `/Users/IF640063/Desktop/example/qa-bootstrap-kit/COGNITIVE_OS_PROJECT.md` §2, §5, §6, §7, §8, §13
- `/Users/IF640063/Desktop/example/qa-bootstrap-kit/DELIVERY_WORKFLOW_PLAN.md` Per-Stage Model Selection Matrix + Banking-Grade Pipeline Properties
- `/Users/IF640063/Desktop/example/qa-bootstrap-kit/e-commerce-v5/output-e5f8b9c2/output.json` — primary holdout
