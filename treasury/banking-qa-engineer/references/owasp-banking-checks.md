# OWASP Top 10 — banking-flavored checks

Walk every category for every artifact. Mark each as **pass** / **fail**
/ **not-applicable** with reason. A category cannot be skipped.

## A01 — Broken Access Control

- [ ] Every protected handler asserts ownership: `(user_id, resource_owner_id)`.
- [ ] Multi-tenant scoping enforced at SQL `WHERE`, not only in app code.
- [ ] Admin routes are in a separate router group with verb-level RBAC.
- [ ] No mass-assignment from request DTO into domain entity / DB row.
- [ ] Account enumeration prevented (404 not 403 on cross-owner access).

## A02 — Cryptographic Failures

- [ ] PII columns (NIK, PAN, account number, DOB) encrypted at rest with KMS-managed keys.
- [ ] TLS enforced end-to-end; no plaintext channel between services.
- [ ] JWT alg pinned; `kid` validated against the signing-key registry.
- [ ] No password / API-key / token in commits, logs, fixtures, or `.env*` files.
- [ ] PRF / hash choice fits use case (password → argon2id; API token → HMAC-SHA256; not MD5/SHA1).

## A03 — Injection

- [ ] SQL: parameterized queries everywhere. Reject string-concat SQL even in admin tools.
- [ ] NoSQL injection: typed query builders, no `$where` style raw input.
- [ ] OS command injection: avoid `exec`/`system`; if unavoidable, allowlist binary + escape args.
- [ ] LDAP / XPath / template injection considered for any feature that accepts a query expression.
- [ ] Output encoding for any rendered HTML; CSP set at the gateway.

## A04 — Insecure Design

- [ ] Threat model documented for the flow (STRIDE + at least one financial-domain abuse case).
- [ ] Rate limits + lockouts on auth, money-movement, and OTP / 2FA endpoints.
- [ ] Underwriter / admin overrides require dual-control.
- [ ] Disbursement / settlement / refund paths have idempotency keys with persisted dedup state.

## A05 — Security Misconfiguration

- [ ] No default credentials.
- [ ] Error responses do not echo SQL, stack traces, internal hostnames, or version banners.
- [ ] Security headers set (HSTS, X-Content-Type-Options, X-Frame-Options or CSP frame-ancestors, Referrer-Policy).
- [ ] Containers run non-root, read-only root filesystem, dropped capabilities.
- [ ] Per-environment secrets — no SIT/UAT/PRD secret reuse.

## A06 — Vulnerable & Outdated Components

- [ ] SBOM generated on build; signed image with verify-on-pull policy.
- [ ] Dependency scanner (govulncheck / npm audit / pip-audit) clean or with documented exceptions.
- [ ] Direct dependencies pinned; floating versions only where the registry guarantees immutability.
- [ ] No EOL runtimes (e.g., Go 1.18, Node 16, Python 3.7).

## A07 — Identification & Auth Failures

- [ ] Refresh-token rotation on use; old token invalidated server-side.
- [ ] Session fixation prevented; session id regenerated on privilege change.
- [ ] MFA / step-up auth on money-movement, KYC update, profile e-mail change.
- [ ] Brute-force protection on login; CAPTCHA or backoff after N failed attempts.
- [ ] Account-recovery flow does not leak whether an email is registered.

## A08 — Software & Data Integrity Failures

- [ ] CI/CD: pinned Action SHAs, scoped `permissions:` block, OIDC to cloud (no long-lived keys).
- [ ] Branch protection: required reviews, signed commits on `main`, no force-push.
- [ ] Image / artifact signing with cosign; verify-on-pull at runtime.
- [ ] Insecure deserialization avoided (no untrusted Pickle, Java native, etc.).
- [ ] PRD deploy gated on approver distinct from PR author.

## A09 — Security Logging & Monitoring Failures

- [ ] PII fields masked before any log/trace/metric emission.
- [ ] Auth headers (`Authorization`, `Cookie`, `X-API-Key`) and JWTs redacted at the logger.
- [ ] Audit log on credit-decision, KYC state change, disbursement, role change with actor + timestamp + before/after + correlation-id, append-only sink.
- [ ] Alerts wired on auth-failure spikes, IDOR/BOLA-deny rates, and 5xx rates per route.
- [ ] Log retention meets the strictest applicable regulation (PCI-DSS = 1y warm + multi-year archive in many jurisdictions).

## A10 — Server-Side Request Forgery (SSRF)

- [ ] Outbound HTTP allowlist (DNS-resolved before connect) on any user-supplied URL.
- [ ] RFC1918 + link-local + 169.254.169.254 (cloud metadata) blocked.
- [ ] Webhook fetchers re-resolve hostname per request to defeat DNS rebinding.
- [ ] No user-controlled URL → internal service (image processor, PDF renderer, html-to-pdf).

## Output shape per finding

```markdown
| Sev | Category | Description | File:line | Mitigation |
|-----|----------|-------------|-----------|------------|
| P1  | A01 BAC  | `GET /loans/:id` lacks ownership predicate | internal/loan/handler.go:42 | Add `borrower_id` filter at SQL `WHERE`; return 404 not 403 on miss |
```

## Verdict rules

- Any **P1** in any category → **REJECT**.
- Any **P2** without a documented mitigation plan → **REJECT** unless human tech-lead overrides.
- All P3 with mitigations + plan → **APPROVE with conditions**.
- All clean → **APPROVE**.
