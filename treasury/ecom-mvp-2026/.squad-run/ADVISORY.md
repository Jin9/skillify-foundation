# ADVISORY.md — Low-Severity Findings

Severity-`low` findings from Reviewer-L1 and Reviewer-L2. Per `docs/orchestrator.md` § Severity Policy, low findings never block; they're the polish backlog.

For full per-finding detail (file:line, recommendation, CWE refs), see `.squad-run/components/<svc>/rev-l1.json` and `.squad-run/layer-2/rev-l2.json`.

---

## System-level (Reviewer-L2)

- **REV-L2-110 (low, maintainability):** outbox + relay + consumer-dedup patterns are duplicated across 4 services with subtle variations. Promote to `backend/common/outbox` shared lib (STRATEGY P3 item — defer until pattern stabilizes across 5+ services).

## identity (Reviewer-L1)

- Missing rate limiting on `/auth/login`, `/auth/register`, `/auth/refresh` — flagged for v2 (CWE-307).
- Email normalization: case-insensitive uniqueness not enforced at the DB layer.
- No body-size cap on auth endpoints.
- `phone` field validation: pattern present but not all callers honor it.

## catalog (Reviewer-L1)

- `inStock=true` filter pagination under-count documented; `total` is pre-stock-filter (advisory note).
- `category.list` `activeOnly=false` exposes inactive categories without admin gate (already in dry-run #1 advisory).
- `events.product.updated` outbox event has no documented consumer; written for future extensibility.

## inventory (Reviewer-L1)

- Sweeper observability gap: skipped non-RESERVED rows (PR-003 race) not logged or counted. Add `slog.InfoContext` + a `inventory_sweeper_candidates_skipped_total{reason}` counter.
- `reservation.create` does not enforce TD's `items.length ≤ 50` cap (DoS amplification risk; trivial to add binding tag).
- `reservation.create` accepts duplicate SKUs in items (aggregate stock math is correct, per-item check is wrong).
- `qty` field has `min=1` but no `max` upper bound (DoS amplification).

## cart (Reviewer-L1)

- `update-item` `qty` lacks upper bound (`max=100`).
- `update-item` qty=0 path returns 404 on already-removed line (breaks AMB-004 idempotent-retry equivalence).
- `AddItem` catalog call has no per-call deadline (relies on gateway timeout).
- Catalog/inventory upstream error strings reach `wrapper.Respond.Err` (verify wrapper doesn't echo to client body).

## checkout (Reviewer-L1)

- Compensation runs on `context.Background()` with no overall deadline (~4s blocking compensation window; INFLIGHT row stays open).
- Saga-log `orderId` empty-string on pre-step-6 errors; migration declares NOT NULL → silent INSERT failures (related to KI-3).
- Catalog fan-out classifies all per-line errors as `PRODUCT_INACTIVE` including 5xx/timeout (gets cached in idempotency envelope for 24h).

## order (Reviewer-L1)

- `idx_orders_buyer_email` migration drops the `LOWER()` expression that TD spec requires for case-insensitive admin email search.
- Order.payment.completed consumer bypasses canonical `ValidateTransition` gate (uses inline status comparison) — copy-paste regression risk in stub consumers.

## payment (Reviewer-L1)

- amount-mismatch logs leak `got` value (CWE-209/532) — mask if sensitive in PRD.
- HMAC hex-decode short-circuits before MAC compute (cosmetic timing nit).
- `internalSecretMW` does single ConstantTimeCompare; can't honor comma-separated rotation overlap mandated by ADR-008 (real rotation 401s legitimate Checkout traffic). Promote to per-pair secrets in v0.8.

## frontend-web (Reviewer-L1)

- `?next=` redirect has no explicit same-origin allowlist; needs `lib/auth/safe-redirect.ts` validator.
- `BACKEND_*_BASE_URL` shown as `http://` in `.env.example`; spec must require `https://` in PRD + CI lint.
- CSP `script-src 'unsafe-inline'` is a yellow flag (pending nonce-middleware follow-up).
- `withAuth` refresh-then-downstream-throw token-rotation gap (medium edge case).
- Forwarding of expired access cookies on public passthroughs (low risk; public routes don't need access at all).

---

## Plan-stage advisories (already in PLAN_NOTES.md)

PR-001 was a medium repaired in cycle 1 (frontend-web admin-dep cleanup). PR-002 / PR-003 are low advisories, written to PLAN_NOTES.md.
