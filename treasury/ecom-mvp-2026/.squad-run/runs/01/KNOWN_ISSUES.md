# Known Issues — Workflow Squad Dry-Run #1

**Run ID:** `ecom-mvp-2026-05-07-001`
**Source workflow:** `B2C E-Commerce Platform/requirement/ecommerce_mvp_agentic_ai_requirement.md`
**Terminal state:** `ShipWithCaveats`
**Cycle accounting:** Layer-2 Reviewer-L2 cycle 1 of cap 1 used.

This file enumerates every accepted defect / scope reduction in the squad's first end-to-end run. Items here did not block delivery; they will not be repaired in this cycle.

For the **full Reviewer-L2 finding objects** (severity, locus, recommendation, routing): see `.squad-run/layer-2/rev-l2.json`.
For **QA-L2 acceptance trace** (per-AC coverage): see `.squad-run/layer-2/qa-l2.json`.
For **low-severity** advisories from Reviewer-L1/L2: see `.squad-run/ADVISORY.md`.
For **plan-stage advisories**: see `.squad-run/PLAN_NOTES.md`.

---

## High-severity caveats (not repaired)

### REV-L2-004 — Compound stub dead end (HIGH, `cross_component_error_propagation_issue`)

A payment failure on a `PENDING_PAYMENT` order has **no recovery path** in this run because the chain of compensating handlers is stubbed:

- `order.cancel-on-checkout-failure` (Order) → 501 stub.
- `events.payment.failed` consumer (Order) → stub.
- `events.payment.expired` consumer (Order) → stub.
- `events.reservation.expired` consumer (Order) → stub.
- `events.payment.failed`, `events.payment.expired`, `events.order.cancelled`, `events.product.created` consumers (Inventory) → stubs.

**Impact:** an order that hits a payment failure stays in `PENDING_PAYMENT` forever. The "self-heal via 15-min reservation TTL" pathway is broken because the consumer that would translate TTL expiry into an `order.status` transition is itself a stub.

**Accepted because:** stub set was a deliberate **buildable-scaffold** scope reduction. Un-stubbing 9 handlers is too large for the L2 repair cycle.

**Action before GA:** un-stub each handler (the FULL `payment.completed` consumer in Order's `consumer.go` is the reference pattern); wire the per-handler state-driven branches; run integration tests on the failure paths.

### REV-L2-001 — `identity.profile.read` still stubbed (HIGH, partial)

The buyerEmail wiring is in place (Checkout calls `identity.profile.read` in step 2b of `handler_commit.go`, then forwards `buyerEmail` to `order.create-from-checkout`), but Identity's `profile.read` handler returns 501. **Every checkout commit fails at step 2b until Identity un-stubs profile.read.**

**Repair file:** `backend/services/checkout/app/checkout/client_identity.go` (already wired).
**Action:** un-stub `POST /api/v1/identity/profile/read` in Identity. ~50 lines: read `users` row by `claims.sub`, return `{userId, email, name, phone, role}` per TD spec.

---

## Medium-severity caveats

### REV-L2-006 — Payment event topic ownership ambiguity (MEDIUM)

`ecom.payment.events` consumer/producer ownership not formally declared in TL contracts. Multiple services consume; canonical producer declaration missing.

### REV-L2-007 — Shared-secret rotation runbook missing (MEDIUM)

`INTERNAL_SHARED_SECRET` rotation procedure is undocumented. A compromise requires manual simultaneous redeploy of all 4 services using the secret.

### REV-L2-008 — Idempotency-key naming inconsistency (MEDIUM)

`idempotencyKey` field in `orderCreateRequest` duplicates `orderId` semantically. Naming pattern divergent vs the `Idempotency-Key` HTTP header used at the Checkout entry point.

### REV-L2-009 — Single secret across 4 services (MEDIUM, `cross_component_security_issue`)

All 4 services using internal-auth (Order, Payment, Inventory, plus Checkout as caller) share **one** secret value. Compromise of any one service compromises the inter-service trust boundary system-wide.

**Mitigation in this run:** secret is env-only (never logged, never written to disk in any service). **Action before GA:** per-edge secrets OR migrate to mTLS / service-to-service JWT.

---

## QA-L2 acceptance gaps (criteria_unmet)

### CU-001 — `traceId` not in response envelope (medium, infrastructure)

`common/wrapper.Response[T]` ships `{code, message, data}` only — `traceId` lives in context but is not serialized to the wire. Affects all 7 backend services. **Single common-lib edit fixes the whole system.** Filed under cross-cutting; was flagged by identity QA-L1 as well (QA-IDENT-001).

### CU-002 — Frontend spine UI partially stubbed

Pages `/checkout`, `/checkout/payment`, `/orders`, `/orders/[id]`, `/register`, `/admin` are "Coming in v2" placeholders. Customer journey AC-35 (UI-driven §7.1 spine) is not executable end-to-end.

### CU-003 — Order failure-branch consumers stubbed (overlaps REV-L2-004)

`onPaymentFailed`, `onPaymentExpired`, `onReservationExpired` are TODO stubs in Order's consumer dispatch. PAY-003, PAY-004, ORD-003 partially un-observable.

### CU-004 — Order admin-transition handlers stubbed

`cancel-mine`, `list-admin`, `update-status-admin`, `cancel-on-checkout-failure` are 501. ORD-004/005/006 not exercisable end-to-end.

### CU-005 — `cart.clear-on-checkout` stubbed

CHK-010 (cart cleared after successful checkout) silently fails — cart is not cleared after a successful checkout commit. Cosmetic for the test path; functional gap for users.

---

## Run-scope reductions (intentional from launch)

These are **not findings** — they were locked in `WORKFLOW_RUN.md` § 1 before launch.

- 4 of 11 bounded contexts deferred to v2: **Shipping, Promotion, Notification, Audit**.
- 4 of 16 modules deferred: §8.7 Promotion/Coupon, §8.11 Shipping Mock, §8.12 Product Review, §8.14 Notification Log, §8.15 Audit Log.
- All 8 components shipped as **buildable scaffolds**, not full implementations: each backend service has 2-4 FULL handlers and the rest as 501 stubs with routes registered. The frontend has 5 FULL pages, 6 stub pages.
- `compliance_sensitive: false` — dual-pass BA + Reviewer-L2 (Opus Max) was not triggered.

---

## Architectural caveats from the single-family deviation

This run used Claude-only sub-agents (no GPT). Per `docs/model-routing.md`, three roles deviated:

- **Plan-Reviewer:** docs say GPT-5.5 xHigh; run used Claude Opus xHigh with red-team-stance prompt. **Outcome:** caught 14 plan-stage findings (4 high), one near-miss self-flagged (PR-004 access-token-in-memory).
- **Tech-Designer:** docs say GPT-5.5 xHigh; run used Claude Sonnet/Opus by complexity tier. **Outcome:** 8 TDs landed cleanly with 5 contract-ambiguities surfaced (mostly around TL contract gaps).
- **QA-L1:** docs say GPT-5.5 high; run used Claude Sonnet. **Outcome:** caught conformance issues per component but two services (catalog, payment) returned `pass` where Reviewer-L1 found mediums; same-family rationalization risk visible but not load-bearing.

**Follow-up natural step:** a hybrid Claude+GPT comparison run on the same input requirement to measure delta on these three roles.

---

## Cycle accounting summary

| Edge | Cycles used | Cap | Status |
|---|--:|--:|---|
| Plan stage (BA + TL) | 1 | 1 | within cap; verdict `advisory_only` after repair |
| Reviewer-L1 ↔ Dev (4 components: identity, inventory, order, checkout) | 1 | 2 | within cap; cycle 2 not run (budget) |
| Reviewer-L1 ↔ Tech-Designer | 0 | 1 | within cap (REV-IDENT-003 fixed at Dev layer; no TD repair needed) |
| Reviewer-L2 ↔ Tech-Lead | 1 | 1 | **at cap**; further block ⇒ ShipWithCaveats |
| Reviewer-L2 ↔ Tech-Designer | 0 | 1 | within cap |
| Tech-Lead reopen from QA-L2 | 0 | 1 | not exercised |
| BA reopen from QA-L2 | 0 | 1 | not exercised |
