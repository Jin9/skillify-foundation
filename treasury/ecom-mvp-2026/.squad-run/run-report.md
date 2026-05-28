---
template_version: 0.1.0
run_id: ecom-mvp-2026-05-08-002
iteration: v0.7
started_at: 2026-05-08
ended_at: 2026-05-08
mode: claude-only
terminal_state: ShipWithCaveats
total_subagents_dispatched: 38
expected_cost_usd: "60-90"
actual_cost_usd: "~95"
deviations_from_docs: 3
findings_total: { high: 13, medium: 24, low: 18 }
---

# Run Report — v0.7 Dry-Run #2 (B2C E-Commerce Platform Refactor)

## Summary

Dry-run #2 ran the v0.7 squad design (10 LLM roles, 18 templates, hybrid Claude+GPT model map) against an **enhanced requirement** that added `ShopPilot Frontend Spec.md` (Thai mobile webview, 390×844 viewport, emerald `#1a8060`, IBM Plex Sans Thai, 5-tab nav, 13 routes) plus a designer's standalone HTML mockup as visual oracle. Mode was incremental — kept the dry-run #1 service tree, archived `.squad-run/` to `runs/01/`, refactored over the existing artifacts.

**Phase 1 prereqs landed cleanly:** traceId in `common/wrapper.Response[T]`, `identity.profile.read` un-stubbed with a 5-case test, `go.sum` regenerated across all 7 services (also patched 5 services' `replace` directives from a broken GitLab path to local `../../common`), 3 Python orchestrator validators (`validate_ba.py`, `validate_story_trace.py`, `validate_components.py`), `docs/skills.md` registry, and 10 sub-agent prompt templates under `docs/prompts/`. Pre-existing build errors from dry-run #1 (3 import cycles in cart/inventory/payment, 2 type-mismatch in catalog/checkout) recorded for Wave B Devs to clear during refactor — they did, all 7 services now build clean.

**Phase 3 squad pipeline executed end-to-end** with one plan-stage repair cycle (1 medium, frontend-web admin-deps cleanup) and one Layer-2 repair cycle (3 cross-component highs closed: order using its own UUID instead of orderId, common/middleware missing tokenType check, Checkout outbound calls missing X-Internal-Secret). Reviewer-L2 still flagged 3 unresolved highs that require larger work (3 async consumer un-stubs in Order; Idempotency-Key carriage heterogeneity; broken bcrypt dummy compare for AUTH-007 timing equivalence) — all captured in `KNOWN_ISSUES.md` for v0.7.1.

**Hybrid mode was promised but harness can't dispatch GPT** — recorded as a known orchestrator-session limit; 3 cross-family deviations (Plan-Reviewer, Tech-Designer, QA-L1) reproduced from dry-run #1. Hybrid run is deferred to v0.8 when a GPT-capable orchestrator lands.

## Acceptance gate status

| Gate (from `docs/v0.7-plan.md`) | Status | Notes |
|---|---|---|
| All 13 P0 items ☑ | **partial** | 6 critical P0 prereqs landed in Phase 1; remaining 7 P0 items (post-run report template, _test.go siblings on every FULL handler — done partially per service) covered in this run |
| ≥ 12 of 16 P1 items ☑ | **passed** | 16/16 P1 items have artifacts emitted (EPICs/Stories ✓, infra docs ✓, per-API spec MDs ✓, ERDs ✓, observability spec ✓, ADRs ✓, QA CSVs ✓, etc.) |
| Structural validator refuses fan-out on malformed `ba.json` | **passed** | Manually tested — removing a `then` clause causes `validate_ba.py` exit-2 with clear error |
| Terminal `Done` or `ShipWithCaveats` with KNOWN_ISSUES citing P2/P3 only | **partial fail** | Terminal IS `ShipWithCaveats`, but KNOWN_ISSUES cites P0/P1-class items (broken bcrypt dummy, stub consumers). Strict acceptance fails; promoted as v0.7.1 work |
| Dry-run produces `run-report.md` per template | **passed** | This file. `docs/templates/run-report.md` was authored alongside as the seed |
| Cost actuals within ±25% of $60-90 envelope | **partial** | ~$95 actual (slightly over; 5 sub-agents over the 33 estimated due to 3 retry dispatches) |
| Hybrid deviation count: zero | **failed** | 3 deviations recorded (harness can't dispatch GPT this session) |

## Stage outcomes

| Stage | Verdict | Sub-agents | Output paths | Repair cycles |
|---|---|--:|---|--:|
| Iteration Planner (Phase 2) | Approved | 1 (orchestrator self) | `docs/v0.7-plan.md` | 0 |
| BA | pass (validators clean) | 1 | `.squad-run/ba.json`, `requirement/EPIC_*/` | 0 |
| Tech-Lead | pass (validators clean) | 1 (+1 Opus retry on overload) | `design/architecture/*` + 10 ADRs | 0 |
| Plan-Reviewer | block → repair → pass | 1 + 1 repair | `.squad-run/plan-review.json` | 1 (TL surgical edit on PR-001 medium) |
| Tech-Designer fan-out (Wave A) | partial then complete | 8 (+ 3 retries: identity, catalog, checkout) | 8 td.json, 7 erd.md, 56 API spec MDs | 1 |
| Dev fan-out (Wave B) | clean builds | 8 | 7 service trees + frontend reshape | 0 |
| QA-L1 + Reviewer-L1 (Wave C) | mixed | 16 | `.squad-run/components/<svc>/qa-l1.{json,csv}` + `rev-l1.json` | 0 |
| Layer-2 | block → repair | 2 + 1 repair | `.squad-run/layer-2/qa-l2.{json,csv}`, `rev-l2.json` | 1 (3 cross-component highs closed) |
| Terminal | ShipWithCaveats | — | this file + `KNOWN_ISSUES.md` (forthcoming) | — |

## Closure of prior-run highs

Dry-run #1 surfaced 4 Reviewer-L2 highs. Status this run:

| Prior finding | Status | Evidence |
|---|---|---|
| REV-L2-001 buyerEmail wire | **Closed** | `identity.profile.read` FULL handler (Phase 1 un-stub) + Checkout step 2 calls `client_identity.ReadProfile()` + Order's `create-from-checkout.buyerEmail` is `binding:"required,email"` |
| REV-L2-002 env-var split | **Closed** | All 4 services (order, payment, inventory, checkout) reference exactly `INTERNAL_SHARED_SECRET`; `INTERNAL_SECRET` is gone |
| REV-L2-003 constant-time compare | **Closed** | `crypto/subtle.ConstantTimeCompare` in middleware on order, inventory, payment, checkout; no plain `==` / `!=` survives |
| REV-L2-004 compound stub dead-end | **Partial** | `order.cancel-on-checkout-failure` un-stubbed (sync compensation works); but `onPaymentFailed` / `onPaymentExpired` / `onReservationExpired` async consumers in Order are still TODO stubs (re-fired as REV-L2-103 high this run) |

## New findings (top 10 by severity, distinct from prior-run highs)

Cross-component (Reviewer-L2):

1. **REV-L2-101 (high, data flow)** — Order generated its own UUID for `orders.id` instead of using Checkout's `orderId`. Closed by L2 repair: handler now parses `req.IdempotencyKey` as UUID, idempotent replay returns 200.
2. **REV-L2-102 (high, security)** — `common/middleware/auth.go` accepted refresh tokens on Bearer routes (no `tokenType == "access"` check). Closed by L2 repair: middleware now rejects non-access tokens with 401.
3. **REV-L2-105 (high, security)** — Checkout's outbound calls to Inventory and Payment didn't include `X-Internal-Secret`. Closed by L2 repair: clients now carry the secret and set the header.
4. **REV-L2-103 (high, error propagation)** — REV-L2-004 partial closure. **Open** — async consumer un-stubs deferred to v0.7.1.
5. **REV-L2-104 (high, data flow)** — Idempotency-Key carriage heterogeneous: header at frontend (dropped at proxy catch-all), body at inventory/order, missing at payment. **Open** — needs cross-cutting cleanup in v0.7.1.

Per-component (Reviewer-L1, samples):

6. **REV-IDENT-101 (high)** — AUTH-007 dummy bcrypt compare uses an invalid hash (`$2a$12$dummy...`); the timing oracle for unknown-email vs wrong-password is fully reopened. Pre-compute one valid bcrypt sentinel hash at init. **Open** — Dev repair candidate.
7. **REV-INV-SEV1 (high)** — `stock.adjust` is admin-only per TD but the route only validates JWT signature; no role gate. Once the stub un-stubs, customer JWTs will mutate inventory.
8. **REV-CART-001 (high)** — `cart.read` fan-out swallows upstream errors silently (catalog timeout → silent fallback) instead of strict UPSTREAM_TIMEOUT 504.
9. **REV-CHECKOUT-001 (high)** — `finalizeAndReturn` passes empty `orderID=""` to `saga_log.Write` whose column is `UUID NOT NULL`; cast fails, tx aborts, idempotency_keys row stays INFLIGHT, next retry creates duplicate orders.
10. **REV-FRONTEND-WEB-001 (high)** — `/api/checkout/commit` and `/api/payment/simulate` route handlers don't exist; both endpoints fall through to the generic `[...path]` catch-all which doesn't forward `Idempotency-Key`. Replay-safety broken end-to-end despite client correctly sending the header.

QA-L1 / QA-L2 most-impactful:

11. **QA-L2 l2-fail-003 (criteria_unmet)** — Payment service is missing the `outbox_relay` goroutine. `events.payment.completed/failed/expired` are written to the outbox table but never produced to Kafka. Order's consumer never receives them; orders stuck in `PENDING_PAYMENT`. **§7.1 customer-journey end-state `status=PAID` is structurally unreachable** even on a happy-path simulate-success.

Full per-component finding files in `.squad-run/components/<svc>/{qa-l1.json, rev-l1.json}` and `.squad-run/layer-2/{qa-l2.json, rev-l2.json}`.

## Deviations

| Role | Docs assignment | Run assignment | Reason | Impact |
|---|---|---|---|---|
| Plan-Reviewer | GPT-5.5 xHigh (cross-family) | Claude Opus xHigh + red-team prompt | Harness can't dispatch GPT in this orchestrator session | Same-family rationalization risk; partially compensated by red-team prompt; 1 medium found, 3 low |
| Tech-Designer | GPT-5.5 xHigh | Claude Sonnet 4.6 (Opus on complex) | Same | 3 of 8 TDs needed retry (identity, catalog, checkout entered design-mode without writing files); successful on second pass |
| QA-L1 | GPT-5.5 high (cross-family vs Claude Dev) | Claude Sonnet 4.6 high | Same | Catalog and Payment QA-L1 returned conditionally `pass` while Reviewer-L1 found mediums on the same code — same-family rationalization visible |

## Cost ledger (best-effort estimate)

Cost ledger (STRATEGY#P1-3) is not yet wired; numbers below are token-count summations from each sub-agent's reported `total_tokens`:

| Stage | Sub-agents | ~Tokens | Est. cost |
|---|--:|--:|--:|
| Plan stage (BA + TL + Plan-Reviewer + 1 TL repair) | 5 | ~570k | ~$11 |
| Wave A (8 TDs + 3 retries) | 11 | ~1130k | ~$23 |
| Wave B (8 Devs) | 8 | ~880k | ~$18 |
| Wave C (16 reviewers) | 16 | ~1640k | ~$33 |
| Layer-2 (QA-L2 + Rev-L2 + 1 repair) | 3 | ~380k | ~$8 |
| **Total** | **~38 dispatches + 5 retries = ~43** | **~4.6M** | **~$95** |

**~$95 actual vs $60-90 estimated.** Slightly over. Sources of overrun: 3 TD retries (Wave A); 1 plan-stage repair; 1 L2 repair. The cost estimate model in `docs/v0.7-plan.md` assumed 0 retries; needs calibration in v0.8.

## Recommendations for next iteration (v0.7.1 / v0.8)

Tactical (v0.7.1 — close v0.7 ShipWithCaveats):

1. Un-stub Order's 3 async consumers (`onPaymentFailed`, `onPaymentExpired`, `onReservationExpired`) — closes REV-L2-103 + l2-fail-003 indirectly when paired with Payment's outbox relay.
2. Add Payment's `outbox_relay` goroutine (mirror Order's pattern). Closes l2-fail-003 — orders can finally reach `status=PAID` on simulate-success.
3. Fix REV-CHECKOUT-001 finalize-empty-orderId bug. Single-line fix; high impact (idempotency totally broken otherwise).
4. Pre-compute valid bcrypt sentinel hash for AUTH-007 timing-equivalence. Closes REV-IDENT-101.
5. Add ADMIN role gate on `inventory.stock.adjust`. Closes REV-INV-SEV1.
6. Standardize Idempotency-Key carriage cross-system (header end-to-end, no body-field-vs-header heterogeneity).
7. Frontend: add dedicated `/api/checkout/commit` and `/api/payment/simulate` route handlers that explicitly forward Idempotency-Key.

Strategic (v0.8 candidates from STRATEGY.md):

- Wire `agent-spend-review` skill for a structured cost ledger so the run-report cost section is real, not estimated.
- Hybrid run when GPT dispatch becomes available — measure the deviation cost we currently document but can't quantify.
- `docs/prompts/<role>.md` self-check loop: TD retries (3 of 8) suggest the prompt template needs a stronger "WRITE FILES NOW" forcing function for output-heavy roles.

## Manual verification notes

- `validate_ba.py` smoke-tested with deliberately malformed input → correctly rejects.
- `go build ./...` passes for all 7 backend services after Wave B repairs.
- Frontend visual conformance NOT manually tested (sandbox can't run `next dev`); ShopPilot Mobile HTML mockup not opened side-by-side. Defer to v0.7.1 demo.
- Reviewer-L2 explicitly noted: **sole-writer rule HELD** — only `order/storage_order.go::UpdateStatus` mutates `orders.status` across all 7 backends.

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Orchestrator (Claude Opus 4.7 1M, this session) | Generated as the v0.7 dry-run #2 seed instance of `docs/templates/run-report.md` |
