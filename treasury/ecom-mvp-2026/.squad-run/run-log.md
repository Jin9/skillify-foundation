# Run Log — Workflow Squad Dry-Run #2 (v0.7)

**Run ID:** `ecom-mvp-2026-05-08-002`
**Iteration:** v0.7 (per `docs/v0.7-plan.md`)
**Mode:** claude-only (deviation from hybrid-mode plan due to harness GPT-dispatch limit)
**Terminal:** `ShipWithCaveats`

This file is the human narrative; for the structured per-stage run-report see `.squad-run/run-report.md`.

---

## Phase 0 — Archive + prep

- Moved `.squad-run/` → `.squad-run/runs/01/` (preserved dry-run #1 history)
- Pinned `backend/go-template@08eccda` in `docs/templates/dev.md`
- Initialized fresh `state.json` (`run_id: ecom-mvp-2026-05-08-002`, `mode: hybrid`, then `mode: claude-only` post-discovery)

## Phase 1 — Critical v0.7 P0 prereqs (orchestrator-side)

| Item | Result |
|---|---|
| F#P0-2 traceId in `common/wrapper.Response[T]` | shipped — `TraceID string \`json:"traceId,omitempty"\`` + middleware fallback via `wrapper.CtxTraceID` |
| F#P0-8 `identity.profile.read` un-stubbed | shipped — handler + 5-case `_test.go`; service interface extended with `GetUserByID` |
| F#P0-5 `go.sum` regen across 6 services | shipped — also patched 5 services' `replace` directives from broken GitLab path → local `../../common`; **5 pre-existing dry-run #1 build errors recorded** for Wave B Devs |
| F#P0-1 structural validators | shipped — 3 Python scripts; smoke-tested negative inputs; `validate_components.py` correctly catches v0.5→v0.6 schema gaps |
| S#P0-1 `docs/skills.md` | shipped — closes `decisions.md` open #3 |
| S#P0-3 `docs/prompts/<role>.md` × 10 | shipped — binding role prompts for all 10 squad roles |

## Phase 2 — Iteration Planner checkpoint

- Read `docs/v0.7-plan.md` + the new `Frontend Spec.md` + `ShopPilot Mobile (standalone).html`.
- Amended plan: noted hybrid-mode harness limit (GPT dispatch unavailable in this orchestrator session); noted Frontend Spec adds Thai mobile UX as the new shaping factor.
- Marked plan `status: Approved`.

## Phase 3 — Squad pipeline

### Step 1: BA (Opus xHigh)
- Output: `ba.json` (71 in_scope, 44 AC, 17 edges, 23 out_of_scope) + 8 EPIC dirs + 39 stories.
- 5 frontend-shaping changes applied: admin UI removed from frontend-web scope, 13 routes added per Frontend Spec, Thai-mobile non_functional block, EPIC_FRONTEND with 12 stories, address-ownership AC explicit.
- `validate_ba.py` PASS, `validate_story_trace.py` PASS (8 EPICs, 39 stories, 81 story-ACs covering 44 ba ACs).

### Step 2: Tech-Lead (Opus xHigh, 1 retry on overload)
- Output: 57 contracts (vs 55 in run #1), 8 components, 6 architecture docs (infra-summary, infra-topology, connectivity, observability-spec, ADRs/), 10 ADRs.
- All 3 dry-run #1 REV-L2 highs encoded by construction in contracts.
- `validate_components.py` PASS (8 components, 57 contracts, 0 orphan deps).

### Step 3: Plan-Reviewer (Opus xHigh + red-team)
- **block** verdict: 1 medium PR-001 (frontend-web admin-deps) + 3 low. Routes: 1 TL, 3 advisory.
- TL repair (Sonnet, surgical): removed 9 admin contracts from `frontend-web.dependencies`. Validator clean post-repair.
- Per lean rule: skipped Plan-Reviewer round 2 (medium-only repair, no high; validator passes; cycle 1 of 1 used).

### Step 4: Per-component fan-out

**Wave A — 8 Tech-Designers (3 retries needed):**
- Identity, catalog, checkout entered design-mode without writing files on first dispatch (Sonnet failure mode on long structured-output tasks). Re-dispatched with tighter "WRITE FILES NOW" prompts.
- All 8 td.json + 7 erd.md (frontend-web no ERD per spec) + 56 API spec MDs landed.
- Lock-order pin documented at all Inventory lock points; state-driven consumer rule applied across consumers.
- Frontend-web TD extracted Thai mobile visual system from Frontend Spec + HTML mockup; HTML never dropped into `frontend/`.

**Wave B — 8 Devs (clean):**
- All 7 backend services build clean post-Wave B (5 pre-existing dry-run #1 errors fixed: 3 import cycles via Identity-pattern interface boundary; 2 type-mismatch wrapper.Code/Message casts).
- Tests added across services: identity (2 new + 1 from Phase 1), catalog (3), inventory (4), cart (3), checkout (2), order (2), payment (1), frontend-web (visual reshape, no test runner in sandbox).
- Frontend-web reshaped: lang="th", IBM Plex Sans Thai, emerald `#1a8060`, ฿ prefix, mobile 390-wide, 5-tab nav, Thai labels throughout, new TabBar/PaymentSimulator/CheckoutCommitButton components.

**Wave C — 16 reviewers (parallel):**
- QA-L1: 6 of 8 components `code_mismatch` (interface drift TD↔Dev), 1 PASS_WITH_STUBS (order), 1 PASS_WITH_2_LOW (payment).
- Reviewer-L1: 12+ high findings across 6 components — notably checkout's finalize-empty-orderId bug (idempotency totally broken otherwise), identity's broken bcrypt dummy compare (AUTH-007 timing oracle reopened), order using fresh UUID instead of orderId, missing X-Internal-Secret on Checkout's outbound calls.

### Step 5: Layer-2 barrier
- **QA-L2 (Opus xHigh):** `criteria_unmet`. AC trace 17 full / 19 partial / 13 stub / 1 uncovered. 4 routed findings (1 TD, 3 Dev). Top concerns: bulk-read stubbed but called; orders.id ≠ orderId; missing Payment outbox relay (orders never reach PAID); missing dedicated frontend route handlers for Idempotency-Key.
- **Reviewer-L2 (Opus, Max-via-prompt):** `block`. 6 high + 3 medium + 1 low. Closures of dry-run #1: REV-L2-001/002/003 ✓ CLOSED; REV-L2-004 PARTIAL (sync compensation done, async consumers still stubs). Sole-writer rule HELD across all 7 backends.
- Top 5 cross-component highs: REV-L2-101 (Order UUID mismatch), REV-L2-102 (middleware tokenType missing), REV-L2-103 (compound stub partial), REV-L2-104 (Idempotency-Key heterogeneity), REV-L2-105 (Checkout outbound missing X-Internal-Secret).

### L2 repair (Sonnet, surgical multi-service)
- Closed: REV-L2-101, REV-L2-102, REV-L2-105 (3 of 6 highs).
- Order: parses `req.IdempotencyKey` as `orders.id`; idempotent replay returns 200.
- common/middleware/auth.go: rejects non-`access` tokenType.
- Checkout: `client_inventory.go` + `client_payment.go` now carry `internalSecret` and set `X-Internal-Secret` header.
- Build clean: order, checkout, common.

### Layer-2 cycle accounting
**Reviewer-L2 ↔ Tech-Lead cycle 1 of 1 used.** Further block ⇒ ShipWithCaveats per orchestrator policy. Selected ShipWithCaveats with KI-1..KI-13 caveats rather than re-running L2.

## Phase 4 — Post-run report

- Authored `docs/templates/run-report.md` (template; F#P0-9 closure).
- Generated `.squad-run/run-report.md` for run #2 as the seed instance.
- Wrote `KNOWN_ISSUES.md` (13 caveats grouped by spine-breaker / security / interface drift / scope reduction).
- Wrote `ADVISORY.md` (low-severity findings from L1 + L2).
- Wrote `PLAN_NOTES.md` (PR-002 + PR-003 advisories).

## Deviations recorded

3 cross-family deviations from `docs/model-routing.md` (same as dry-run #1 — harness can't dispatch GPT in this orchestrator session):
- Plan-Reviewer: docs say GPT-5.5 xHigh; ran Claude Opus xHigh + red-team prompt
- Tech-Designer: docs say GPT-5.5 xHigh; ran Claude Sonnet/Opus by complexity tier
- QA-L1: docs say GPT-5.5 high; ran Claude Sonnet high

Visible cost: 3 of 8 TDs needed retry on first dispatch (Wave A — design-mode without file writes); cost overrun ~5% above the $60-90 envelope (~$95 actual).

## Run statistics

- ~38 sub-agent dispatches + 5 retries = ~43 total
- ~4.6 M tokens, ~$95 estimated cost
- 7 backend services + 1 Next.js frontend; ~280 source files generated/edited
- 8 EPICs / 39 stories / 56 API spec MDs / 7 ERDs / 10 ADRs
- Total findings: 13 high / 24 medium / 18 low

**Run complete.** Pipeline advances to terminal `ShipWithCaveats` with KNOWN_ISSUES.md as the v0.7.1 punch list.
