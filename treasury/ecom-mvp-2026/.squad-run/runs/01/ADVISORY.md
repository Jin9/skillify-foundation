# ADVISORY.md — Low-Severity Findings

Severity-`low` findings (and a select set of mediums tracked but not repaired in cycle 1) from Reviewer-L1 + QA-L1 accumulate here per `docs/orchestrator.md` § Severity Policy. They never block; they're a backlog of polish/quality items the producers can address opportunistically.

---

## identity

- **REV-IDENT (medium, deferred):** missing `tokenType==refresh` check in refresh handler; same key for access+refresh signers; no email normalization (case-insensitive uniqueness); no body-size cap on auth endpoints; no rate-limiting on `/auth/login` and `/auth/refresh`; logout stub means no kill-switch on 30-day refresh tokens.
- **QA-IDENT-002 (low):** login response carries an undocumented `role` field beyond TD response_schema. Tighten or update spec.
- **QA-IDENT-003 (low):** `phone` field on register lacks pattern validation `^[0-9+\- ]{7,20}$`.
- **QA-IDENT-001 (medium, infrastructure):** `common/wrapper.Response[T]` does not emit `traceId` in the body — every service inherits the gap. Promote to a cross-cutting fix in `common/`.

## catalog

- **CAT-L1-001 (medium):** outbox-relay goroutine has no `defer recover()`; a panic kills the service. `gin.Recovery()` only covers HTTP handlers.
- **CAT-L1-002 (medium):** `extractIsAdmin` in `product.detail` is dead code (claims never set on public route). Recommend separate `/admin/product/detail` route.
- **CAT-L1-003 (medium):** `FOR UPDATE SKIP LOCKED` lock released on `rows.Close()` rather than held through `MarkPublished`. At-least-once still holds; misleading comment.
- **CAT-L1-004 (medium):** ADMIN role check duplicated in every handler — promote to middleware.
- F-001 (low) `inStock` filter not wired to inventory (TODO present); F-002 (low) `stockStatus` hardcoded `OUT_OF_STOCK` fallback; F-003 (info) catalog has 6 stubs not 8.

## inventory

- **INV-L1-05 (medium):** idempotency dedup race — `idemFetch` does `SELECT ... FOR UPDATE` on a non-existent row (acquires no lock). Concurrent same-key first-time calls both proceed; loser hits PK violation → 5xx instead of `IDEMPOTENCY_KEY_INFLIGHT` 409. Fix: `INSERT ... ON CONFLICT (key) DO NOTHING RETURNING ...`.
- G-01 (low) sweeper uses `uuid.New()` (v4) where TD specifies v7. G-02 (low) lock-order comment in `handlePaymentCompleted` conflates batch-scan with per-row pin.
- Stub consumers (4 of 5) and `stock.adjust` 501 stub are **scope decisions** for the buildable-scaffold run, NOT defects. Tracked in KNOWN_ISSUES.md as run-scope caveats.

## cart

- **H-1 (medium, deferred to un-stub time):** `clear-on-checkout` is 501-stubbed; `customerUserId == claims.sub` (AMB-002) MUST land at un-stub time to prevent IDOR.
- **H-2 (medium):** inventory-down silently sets every line `checkoutable=false`/`subtotal=0` with no client-visible reason. Add `reason`/`inventoryStatus` field for clarity.
- **M-1 (medium):** `updateItemRequest.Qty` has only `min=0`, no `max=100` upper bound (mismatch with `addItemRequest`).
- **M-2 (medium):** `UpdateItem(qty=0)` on already-removed line returns 404 — breaks AMB-004 idempotent-retry equivalence.
- **M-3 (medium):** `AddItem` catalog call has no per-call timeout; relies on gateway timeout.
- M-4 (medium); INFO observations on `added_at` reset and `RemoveItem` service ready but handler stubbed.

## checkout

- **M01 (medium):** TD says step-8 compensation is "cancel THEN release"; code runs both unconditionally — pick one and align.
- **M02 (medium):** compensation runs on `context.Background()` with no overall deadline.
- **M03 (medium):** `finalizeAndReturn` always passes `orderID=""` to `saga_log.Write`; migration declares column NOT NULL → silent INSERT failures on pre-step-6 rejection.
- **M04 (medium):** catalog fan-out classifies *every* per-line error as `PRODUCT_INACTIVE` including 5xx/timeout — gets cached for 24h via idempotency envelope.
- 3 low advisories: non-deterministic canonical-hash, switch-case fallthrough on ABANDONED, minor-unit comment alignment.

## order

- **MEDIUM:** stub consumers `payment.failed/expired/reservation.expired` and 4 stub admin/customer-cancel handlers — scope decisions for buildable scaffold; tracked in KNOWN_ISSUES.md.
- **LOW:** `AutoLoggingMiddleware` body-logging behavior cross-cutting concern; admin role enforcement currently per-handler — promote to middleware on `adminGroup` when stubs land.

## payment

- **PAY-L1-001 (low):** amount values logged at WARN (callbacks may surface in audit). Mask if sensitive.
- **PAY-L1-002 (low):** hex-decode short-circuits before MAC compute (cosmetic).

## frontend-web

- **OBS-001 (medium):** `?next=` redirect lacks explicit same-origin allowlist in spec — add `lib/auth/safe-redirect.ts` validator before un-stubbing checkout flow.
- **OBS-002 (low):** `BACKEND_*_BASE_URL` shown as `http://` in `.env.example` (intra-cluster); spec must require `https://` in prod + CI lint.
- **OBS-003 (low):** CSP `script-src 'unsafe-inline'` is a yellow flag — nonce-middleware follow-up; revisit in L2.
- QA gaps: vitest + Playwright test suite absent; `/cart` is fully Client where TD said SC shell + interactive Client; env var name suffix divergence.

---

## Plan-stage advisories (already in PLAN_NOTES.md)

PR-010..PR-015 — 6 advisories from Plan-Reviewer rounds 1+2.
