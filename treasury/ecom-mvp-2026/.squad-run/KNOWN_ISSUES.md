# Known Issues — Dry-Run #2 (v0.7)

**Run ID:** `ecom-mvp-2026-05-08-002`
**Terminal state:** `ShipWithCaveats`
**Source of truth for full findings:** `.squad-run/components/<svc>/{qa-l1.json, rev-l1.json}` per service + `.squad-run/layer-2/{qa-l2.json, rev-l2.json}` system-level. This file is the operator's punch list.

Each item below is a **P0 or P1 caveat** carried into v0.7.1 because we hit Layer-2 cycle 1 of cap 1 and chose `ShipWithCaveats` over a second L2 repair. Items are grouped by impact on the §7.1 customer-journey spine.

---

## §7.1 spine breakers (must close before "demo-ready")

These three items mean the happy-path browse → cart → checkout → simulate-payment → see-paid-order journey **does not work end-to-end** even with all current code in place.

### KI-1 — Payment service has no outbox relay goroutine

**Source:** QA-L2 `l2-fail-003`; QA-L1 catalog of payment finds the same.
**Where:** `backend/services/payment/main.go` — relay goroutine never started.
**Effect:** `events.payment.completed/failed/expired` rows are written to `payment.outbox_events` but never produced to Kafka. Order's consumer never receives them. **Orders stuck in `PENDING_PAYMENT` forever** on simulate-success.
**Fix:** mirror the outbox-relay pattern already in Order/Catalog/Inventory. ~50 lines + a goroutine in main.

### KI-2 — Order's 3 async consumers still TODO stubs

**Source:** REV-L2-103 (high). Partial closure of dry-run #1 REV-L2-004.
**Where:** `backend/services/order/app/order/consumer.go` — `onPaymentFailed`, `onPaymentExpired`, `onReservationExpired` return nil immediately.
**Effect:** abandoned-checkout / declined-simulate / TTL-expiry paths leave orders frozen.
**Fix:** copy the `onPaymentCompleted` state-driven pattern (FOR UPDATE on the order row, ValidateTransition, status update, status_history insert, optional outbox emit). Each consumer ~60 lines.

### KI-3 — Checkout's `finalizeAndReturn` writes empty `orderID=""` to a NOT NULL column

**Source:** REV-CHECKOUT-001 (high).
**Where:** `backend/services/checkout/app/checkout/handler_commit.go` — early-return paths pre-step-6 call `finalizeAndReturn` with empty orderId; `saga_log.order_id UUID NOT NULL` cast fails.
**Effect:** **Idempotency contract broken on every successful commit.** The wrapping pgx tx aborts in the cleanup write, `idempotency_keys.Finalize` is rolled back, row stays `INFLIGHT`. The 60s janitor flips it to `ABANDONED`, the next retry re-executes from scratch — duplicate orders + reservations + payment intents.
**Fix:** make `saga_log.order_id` NULLABLE OR skip `saga_log.Write` on the empty-orderId paths. ~5-line surgical fix.

---

## Security highs (must close before any non-MVP usage)

### KI-4 — AUTH-007 timing oracle reopened (broken dummy bcrypt)

**Source:** REV-IDENT-101 (high, code_security_issue).
**Where:** `backend/services/identity/app/identity/service.go:124` — dummy hash is `$2a$12$dummyhashfortimingnulltarget....` which is **not** a valid bcrypt hash.
**Effect:** `bcrypt.CompareHashAndPassword` returns `ErrHashTooShort` in ~40 ns vs ~175 ms for a real compare. **~4 million× wall-clock skew on the unknown-email branch.** Username enumeration via timing is fully back.
**Fix:** at init time, run `bcrypt.GenerateFromPassword([]byte("sentinel"), 12)` once and reuse the result on the dummy-compare path.

### KI-5 — `inventory.stock.adjust` admin role guard missing

**Source:** REV-INV-SEV1 (high).
**Where:** `backend/services/inventory/router/router.go:75-83` — same JWT-only middleware group as `stock.read`.
**Effect:** Once the stub un-stubs (TD spec mandates ADMIN-only), customer JWTs will mutate stock + audit log records will pin breaches on the wrong actor.
**Fix:** add `middleware.RequireRole("ADMIN")` to the route or implement an inline `claims.role == "ADMIN"` check in the handler.

### KI-6 — Idempotency-Key carriage heterogeneity across services

**Source:** REV-L2-104 (high, cross_component_data_flow_issue).
**Where:** Frontend sends header at `/api/proxy/checkout/commit` → caught by generic `[...path]` proxy that drops headers; Inventory + Order read from body field; Payment reads from header. Three different conventions in one system.
**Effect:** Customer double-tapping Place-Order can create duplicate orders because the dedup chain doesn't carry through.
**Fix:** standardize on **HTTP header** end-to-end. Update generic proxy to forward `Idempotency-Key`; update Inventory + Order to also accept the header; add dedicated frontend route handlers for checkout.commit + payment.simulate.

### KI-7 — Inventory outbox relay drops events when `KAFKA_ENABLED=false`

**Source:** REV-INV-SEV3 (high, error propagation).
**Where:** `backend/services/inventory/app/inventory/outbox_relay.go` — disabled-Kafka path marks rows `published_at=NOW()` without publishing.
**Effect:** Once Kafka is re-enabled, those rows are never re-emitted. Silent data loss; at-least-once contract broken.
**Fix:** in the disabled-Kafka branch, log a warning and `return` without advancing rows. Outbox accumulates until Kafka comes back.

### KI-8 — Inventory idempotency-key INFLIGHT race

**Source:** REV-INV-SEV2 (high, race condition CWE-362).
**Where:** `backend/services/inventory/app/inventory/handler_reservation_create.go` — `SELECT … FOR UPDATE` on a non-existent row acquires no lock; concurrent same-key callers both pass the existence check.
**Effect:** Concurrent same-orderId reservation.create calls produce a winner with 200 + a loser with 503 (PK violation surfaced as `DATABASE_UNAVAILABLE`) instead of the loser replaying the winner's envelope.
**Fix:** mirror the Checkout pattern — `INSERT … ON CONFLICT DO NOTHING RETURNING` for atomic-claim. Same as dry-run #1's H02 fix in Checkout.

---

## Interface drift (Dev ↔ TD spec mismatches; close before integration testing)

These are the QA-L1 `code_mismatch` findings. None are functional bugs; they are **interface drift** between TD specs and Dev code that will fail integration tests.

- **KI-9** Identity: response envelope `code` is `"0000"` literal vs TD's `"SUCCESS"`; profile.read response missing `defaultAddressId/status/createdAt/updatedAt` per TD; profile.read route is `GET /api/v1/identity/profile/me` while TD specifies `POST /api/v1/identity/profile/read` (Checkout's `client_identity.go` calls the TD path).
- **KI-10** Catalog: `price` is JSON string in response vs TD's JSON number; `search > 200 chars` is silently truncated vs TD's VALIDATION_ERROR 400; ADMIN detection on public detail route is dead code (claims never set on public routes).
- **KI-11** Cart: `CartItemView` uses `float64` for prices vs cross-cutting int64-minor-units rule; missing `available` + `unavailReason` fields per TD; STOCK_INSUFFICIENT mapped to 409 vs TD 400 VALIDATION_ERROR.
- **KI-12** Checkout: `client_inventory.go` and `client_payment.go` previously didn't send `X-Internal-Secret` (closed by L2 repair as REV-L2-105). Test count drift (Dev claimed 41; actual 34).
- **KI-13** Frontend: `OrdersPage` calls `/api/proxy/order/list` vs TD's `/api/proxy/order/list-mine`; Search tab uses `/products` vs TD's `/search` (route doesn't exist); `clearAccessCookie/clearRefreshCookie` omit `SameSite` and `Secure` attributes that the build counterparts include — production browsers may not overwrite the original session cookie on AUTH_REVOKED.

---

## Run-scope reductions (intentional, recorded)

These are NOT findings — they were locked in `docs/v0.7-plan.md` § Out-of-scope before dispatch.

- 4 of 11 bounded contexts deferred (Shipping, Promotion, Notification, Audit) — same as dry-run #1.
- Address handlers in Identity: `address.set-default` un-stubbed in Wave B; the other 4 (`create`, `list`, `update`, `delete`) remain 501 stubs.
- Catalog admin handlers: `product.update`, `product.soft-delete`, all 4 `category.*` admin handlers still 501.
- Cart `remove-item` and `clear-on-checkout` still 501 (service-layer ready; un-stub trivially).
- Order admin handlers `cancel-mine`, `list-admin`, `update-status-admin` still 501.

---

## Architectural caveats from the single-family deviation (3 of 10 roles)

This run's Plan-Reviewer / Tech-Designer / QA-L1 ran on Claude instead of GPT-5.5 because the orchestrator session has no GPT dispatch path. **Risk:** same-family rationalization. **Visible cost** in this run:

- 3 of 8 Tech-Designers entered design-mode without writing files (identity, catalog, checkout) on first dispatch; needed retry. Hypothesis: Claude's verbose-reasoning bias overpowered the "write files" forcing function. Reflects in cost overrun (~$95 vs $60-90 estimate).
- Catalog and Payment QA-L1 returned `pass`/`pass-with-low` while Reviewer-L1 found mediums on the same code — same-family conformance check missed quality issues a cross-family check might have caught.
- Plan-Reviewer caught 1 medium + 3 low; deferred to v0.8 to compare against a hybrid run when GPT dispatch is available.

---

## Cycle accounting

| Edge | Cycles used | Cap | Status |
|---|--:|--:|---|
| Plan stage (BA + TL routes) | 1 | 1 | within cap; Plan-Reviewer round 2 skipped per lean rule (validator clean post-repair) |
| Reviewer-L1 ↔ Dev | 0 | 2 | not used; per-component repairs deferred to v0.7.1 |
| Reviewer-L1 ↔ TD | 0 | 1 | not used |
| Reviewer-L2 ↔ TL | 1 | 1 | **at cap**; further block ⇒ ShipWithCaveats |
| Reviewer-L2 ↔ TD | 0 | 1 | not used |
| QA-L2 → Dev | 0 | 2 | not used; QA-L2's 4 routed findings carried as caveats |
| QA-L2 → TL | 0 | 1 | not used |
| QA-L2 → BA | 0 | 1 | not used |

**Layer-2 ↔ Tech-Lead** is the only edge at cap. Items KI-1..KI-8 above are the consequence of choosing `ShipWithCaveats` over a second L2 cycle.

---

## What "Closed by construction" looked like vs dry-run #1

For the record, here's how dry-run #2 closed the four dry-run #1 Reviewer-L2 highs:

| Dry-run #1 high | Closure mechanism (run #2) |
|---|---|
| REV-L2-001 buyerEmail | TL added `identity.profile.read` as customer-JWT contract; Phase 1 un-stubbed it; Checkout TD wired step 2 → step 7 explicitly; Order TD made `buyerEmail` `binding:"required,email"`. |
| REV-L2-002 env-var split | TL encoded `INTERNAL_SHARED_SECRET` mandate in `cross-cutting.auth` + ADR-008. All 4 services use this single name. |
| REV-L2-003 constant-time compare | TL mandate in `cross-cutting.auth.transport` + ADR-008. Reviewer-L1 fires `code_security_issue` on any `==`/`!=`. All services compliant. |
| REV-L2-004 compound stub dead-end | Partial. Sync `cancel-on-checkout-failure` un-stubbed (Wave B). Async consumer un-stubs deferred to v0.7.1 (KI-2). |
