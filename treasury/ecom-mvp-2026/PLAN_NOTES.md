# PLAN_NOTES.md — Plan-Stage Advisories (Dry-Run #2)

Low-severity Plan-Reviewer findings (`plan_polish` / advisory). Per `docs/orchestrator.md` § Severity Policy, never block; polish backlog.

Source: `.squad-run/plan-review.json` (verdict: `block` round 1; PR-001 medium repaired in cycle 1; PR-002 + PR-003 advisory below).

---

## PR-002 (low / advisory) — `buyerEmailSnapshot` denormalization shape

**Locus:** `design/architecture/contracts.json` `order.create-from-checkout` payload + Order TD `erd.md`.

The buyerEmail snapshot lives on `orders.buyer_email_snapshot` (per-order) rather than denormalized to `order_items` rows. Plan-Reviewer flagged this as wider than necessary on dry-run #1 (PR-002 in run #1 too). Tech-Lead chose per-order; defensible — keep as-is unless admin email-search performance becomes an issue.

## PR-003 (low / advisory) — `STORY_CHECKOUT_COMMIT` AC trace tightening

**Locus:** `requirement/EPIC_CHECKOUT/STORY_CHECKOUT_COMMIT.md`.

The contract-level `binding:"required,email"` on `buyerEmail` in `order.create-from-checkout` doesn't have an explicit AC in `STORY_CHECKOUT_COMMIT`. Add an AC like:

> Given a customer with a registered profile, when they POST `/api/v1/checkout/checkout/commit`, then Order receives `buyerEmail` populated from the customer's profile and persists it on `orders.buyer_email_snapshot`.

Pure trace-tightening; behavior is already implemented.

---

For the plan-stage repaired finding (PR-001 medium — admin-dep cleanup on frontend-web), see the run-log entry under `plan_repair_tl_completed`. Repair was a 9-contract removal from `frontend-web.dependencies[]`; validator passes clean post-repair.
