# ADR-008 — `INTERNAL_SHARED_SECRET` env var name + `crypto/subtle.ConstantTimeCompare` for internal service-to-service auth

- **Status:** Accepted (LOCKED — closes REV-L2-002 and REV-L2-003)
- **Date:** 2026-05-08
- **Deciders:** Tech-Lead

## Context

Inter-service HTTP calls between the four internal services (Checkout->Inventory, Checkout->Order, Checkout->Payment, Order->Payment cancel-on-failure) carry an `X-Internal-Secret` header. The recipient validates the header against a configured secret.

REV-L2 (dry-run #1) found two breakages:
- **REV-L2-002:** three of four services read env `INTERNAL_SHARED_SECRET`; Payment read `INTERNAL_SECRET`. In any environment that set only the canonical name, every Checkout->Payment call rejected with 401 -> Checkout's compensation matrix triggered on every commit -> compound dead-end with REV-L2-004. The split was a silent partial-outage waiting to happen.
- **REV-L2-003:** Order used `crypto/subtle.ConstantTimeCompare` (correct). Inventory and Payment used Go's plain `!=` / `==` string comparison, which short-circuits at the first differing byte. CWE-208 / OWASP API2:2023 / ASVS V2.4.7. An attacker on the inter-pod network could probe the secret byte-by-byte by measuring response-time variance.

## Decision

- **Env var name:** `INTERNAL_SHARED_SECRET` is the ONLY accepted name. All four recipient services (checkout, order, inventory, payment) MUST read this exact variable. Alternative names — `INTERNAL_SECRET`, `SHARED_SECRET`, `SVC_SECRET` — are rejected at startup with a fatal config error.
- **Value format:** ASCII string, MIN 32 bytes (256 bits of entropy). Recipient services validate length at startup; refuse to boot if shorter.
- **Comparison rule:** Recipient services MUST use:
  ```go
  if subtle.ConstantTimeCompare([]byte(supplied), []byte(secret)) != 1 {
      // reject
  }
  ```
  Plain Go `==` / `!=` is FORBIDDEN. Reviewer-L1 / Reviewer-L2 grep for `!=.*Secret\|Secret.*!=\|==.*Secret\|Secret.*==` patterns and flag any hit as a high-severity finding.
- **Shared helper:** `backend/common/middleware/internal_auth.go` is the canonical implementation. The four services call:
  ```go
  router.Use(common_middleware.InternalAuth(secret))
  ```
  In MVP, the four services MAY ship with their own copy of this middleware ONLY if each one provably uses subtle.ConstantTimeCompare; the shared helper is the right end-state per REV-L2-010.
- **Rotation overlap:** the env var MAY accept a comma-separated value `INTERNAL_SHARED_SECRET=new,old`. Recipients try ConstantTimeCompare against each comma-segment; first match wins. This lets ops roll the secret across all four services without a synchronized restart (5-minute overlap window typical).
- **Frontend posture:** the Next.js frontend MUST NOT load `INTERNAL_SHARED_SECRET` and MUST NOT forward `X-Internal-Secret` (REV-L2-013 verified-correct). `frontend/.env.example` explicitly omits the variable; a frontend regression test asserts the absence of `X-Internal-Secret` on customer-originated requests.
- **Audit logging:** on startup every recipient logs `slog.Info("internal_auth_configured", "service", svc, "sha256_prefix", hex(sha256(secret))[:8])` — ops can confirm at a glance that all four services hash the same secret without ever logging the secret itself. On every reject, recipient logs `slog.Error("internal_auth_reject", "remoteAddr", ..., "x_forwarded_for", ..., "path", ..., "requestId", ...)` — NEVER logs the supplied or expected secret value.
- **Metrics:** `internal_auth_rejects_total{service, path}` counter; alert at >10/sec for 5 minutes (possible secret rotation gap or attack).

## Consequences

- **Positive:** dry-run #2 cannot produce REV-L2-002 by construction (single accepted env-var name, startup-fatal on the wrong one).
- **Positive:** dry-run #2 cannot produce REV-L2-003 in code that uses the shared helper; per-service implementations are explicitly grep-checked at review time.
- **Positive:** rotation has a documented runbook and overlap window — no full-fleet downtime.
- **Negative:** still a shared-secret model — leakage from any one service exposes all four. Acceptable for MVP per REV-L2-009; mTLS is the v2 hardening path.
- **Negative:** the comma-separated parsing is non-trivial enough that the shared helper is the right place to land it; per-service drifted implementations risk subtle splits.
