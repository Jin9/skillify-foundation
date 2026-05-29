# Taxonomy: Stack-Specific Review Patterns

Read the section that matches the artifact in front of you. Each area lists concrete patterns to look for, the typical CWE/API/ASVS identifier, and the stack-specific safer pattern.

## Area index — headline focus + standards

Each finding is tagged with one **primary** area (A–K) + optional secondary tags. Jump to the matching section below for the full check list and safer-pattern rewrites.

| # | Area | Headline focus | Standards |
|---|---|---|---|
| A | Application security (Go/Gin) | input binding + size limits, money math via `shopspring/decimal`, middleware order, goroutine context safety | ASVS V1–V14, CWE Top 25 |
| B | API security (Gin + Kong/APISIX) | BOLA / BOPLA / BFLA, JWT alg pinning, rate-limit + size-limit on every public route, route-shadow prevention | OWASP API Top 10 (2023) |
| C | Architecture (DDD / CQRS / event-driven) | aggregate invariants, commands carry identity, bounded-context contracts not shared joins, idempotent replay | — |
| D | AuthN / AuthZ | ownership predicate per handler, mTLS for s2s, tenant scoping at SQL `WHERE`, policy-as-code RBAC | ASVS V4 |
| E | Secrets & configuration | no hardcoded creds, Vault / Secrets Manager via IRSA / ESO, per-environment secrets with ≤ 90-day rotation | NIST SSDF PS.2 |
| F | Logging & observability | PII masked, auth headers / JWTs redacted, append-only audit on credit-decision / KYC / disbursement, no internal-hostname echoes | — |
| G | Database (MySQL / RDS) | parameterized queries, app user lacks `DROP`/`GRANT`/`FILE`, KMS-encrypted PII columns, TLS-enforced connections | CWE-89, CIS MySQL |
| H | Kubernetes & container | digest-pinned base, `runAsNonRoot` + `readOnlyRootFilesystem`, default-deny NetworkPolicy, scoped RBAC verbs, IRSA over keys | CIS Kubernetes, NSA Hardening |
| I | CI/CD & supply chain | SBOM + cosign, secret scanning, branch protection, pinned Action SHAs, scoped `permissions:`, OIDC to cloud | SLSA, NIST SSDF PW.4 |
| J | Event-driven (Kafka) | per-principal topic ACLs, schema-registry compatibility, idempotent producers, dedup-key consumers, header re-validation | — |
| K | Financial / lending data | regulatory inventory (PDPA/GDPR/PCI-DSS/BOT/OJK/MAS), KYC docs via short-lived pre-signed URLs, signed credit-decision events, dual-control on overrides, masked PRD-to-lower-env copies | PCI-DSS v4, BOT/OJK/MAS |

---

## A. Application security (Go / Gin)

| Pattern to flag | Standard | Safer pattern |
|---|---|---|
| `c.ShouldBindJSON(&v)` with no struct tags or no `binding:"required"` | CWE-20, ASVS V5.1 | Tag every field with `binding`; reject early with `c.AbortWithStatusJSON(400, ...)` |
| Reading body via `io.ReadAll(c.Request.Body)` without `MaxBytesReader` | CWE-770 | Wrap once in middleware: `c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)` |
| `interface{}` / `map[string]any` decoding of untrusted JSON | CWE-915 | Decode into a typed struct; reject unknown fields with `dec.DisallowUnknownFields()` |
| `float64` for monetary amounts | CWE-682 | `github.com/shopspring/decimal` or fixed-point integer minor units |
| Custom auth middleware not registered on a route group | CWE-862, API5:2023 | Group protected routes: `r.Group("/v1", auth.Required(), authz.Loan())` |
| `c.Copy()` not used inside a goroutine spawned from a handler | CWE-362 | Always pass `c.Copy()` into goroutines; never mutate the original `*gin.Context` after `Next()` |
| `recover()` middleware after auth middleware | CWE-755 | Order: `Recovery() → RequestID() → Auth() → Authz() → RateLimit() → Handler` |
| `time.Now()` used for audit-window comparisons without monotonic guarantee | CWE-367 | Use `time.Since(start)` (monotonic) for durations; UTC `time.Now().UTC()` for absolute timestamps |
| Panics that leak state via Gin's default recovery | CWE-209 | Custom recovery that emits sanitized 500 with correlation-id only |

**Money-math reminder:** any code involving `loan.Amount * rate` or `installment / N` must use `decimal.Decimal` with explicit rounding mode (`ROUND_HALF_UP` for consumer credit, `ROUND_DOWN` for fees that must not be over-charged); float arithmetic on money is **High** severity because rounding errors compound across the loan lifecycle.

---

## B. API security (Gin + Kong / APISIX) — OWASP API Top 10 (2023)

### B.1 In Gin handlers
| Pattern | API# | Action |
|---|---|---|
| Path/query id used directly in `repo.GetByID` with no owner check | API1 (BOLA) | Pass authenticated subject into repo; filter at SQL `WHERE` |
| Mass-assignment via `c.ShouldBindJSON(&domainEntity)` | API3 (BOPLA) | Bind into a request DTO, map explicitly to domain entity |
| Admin endpoints in the same router group as user endpoints | API5 (BFLA) | Separate group with explicit role check; verb-level RBAC |
| Outbound HTTP using user-supplied URL (e.g., webhook fetcher) | API7 (SSRF) | Resolve hostname, deny RFC1918 + link-local + 169.254.169.254 (metadata) before connect |
| `?fields=...` mass-projection that includes internal columns | API3 (BOPLA) | Allowlist response fields per role |

### B.2 In Kong declarative config (`kong.yml`)
- **Route precedence pitfall:** a route with `paths: ["/v1/public"]` and **no plugins** declared after a more specific authenticated route can shadow auth if `regex_priority` or path overlap is wrong. Always pin auth plugins **per-route**, not only globally — global plugins do not protect against a route declared without them.
- **Required plugins on every public route:** `key-auth` or `jwt`, `rate-limiting`, `request-size-limiting`, `cors` (if needed), `correlation-id`.
- **Consumer credentials:** never share a `consumer.id` across SIT/UAT/PRD; never check the literal `key_auth.key` into git (Critical).
- **`proxy-rewrite`:** verify it does not strip `Authorization` or inject internal headers (`X-Internal-User`) that downstream services then trust.
- **Custom Lua plugins:** review for `string.format` injection, `os.execute` calls (banned), and untrusted `loadstring`/`load` use.

### B.3 In APISIX route YAML
- `consumers` and `consumer-credentials` are stored in etcd; verify etcd is encrypted in transit and at rest.
- Plugin order matters: `key-auth` / `jwt-auth` must run before `proxy-rewrite` so rewritten paths cannot bypass auth.
- `serverless-pre-function` / `serverless-post-function`: review Lua for the same banned calls as Kong.
- `uri-blocker` and `consumer-restriction` can give a false sense of security if combined with a permissive `routes`-level allowlist.

### B.4 OpenAPI hygiene
- `additionalProperties: true` on a payload containing PII → BOPLA risk.
- `string` without `maxLength` → DoS vector.
- Response schema includes fields the client should never see (`internal_score`, `risk_band`, `db_id`).
- Missing `required` arrays → optional fields silently accepted as null and bypass validation.

---

## C. Architecture (DDD / CQRS / event-driven)

### C.1 Aggregate boundary checks
- An aggregate must own its invariants. **Red flag:** a service method that begins a transaction, mutates `LoanAggregate`, then mutates `BorrowerAggregate` in the same transaction. Two aggregates in one transaction is an aggregate-boundary smell **and** a deadlock risk on busy MySQL rows.
- Cross-aggregate consistency must be eventual via domain events, not via direct DB joins.

### C.2 CQRS read-side authorization
- The read model is not authoritative for authorization. **Red flag:** a query handler computes `if user.Role == "underwriter" { return projection }` — authorization must be enforced *before* the projection is queried, ideally at the application service.
- Projections must not include fields the requester is not authorized to see, even if the UI hides them; a curl client can still read them.
- Replica lag (RDS read replica) can let a stale role grant access. For role changes, force a read against the primary for the next N seconds, or invalidate the user's cached principal.

### C.3 Event sourcing pitfalls
- **Replay attacks on event stream:** consumers must be idempotent (see J).
- **PII in immutable event log:** if events contain NIK or bank account numbers and you must honor a GDPR/PDPA erasure request, you need a *crypto-shredding* plan: events store ciphertext; deletion is by destroying the per-subject key in KMS. Plan this **before** writing PII to events; retrofitting is painful.
- **Event versioning:** schema changes must be additive (`BACKWARD` compatibility minimum, `FULL` preferred). Breaking changes are integrity threats.

### C.4 Bounded context isolation
- Shared MySQL schema between two contexts is an anti-pattern. The skill should flag any `JOIN` across schemas owned by different services as a cross-context coupling and a security boundary erosion (one team's DB user can read another team's PII).

---

## D. AuthN / AuthZ

### D.1 Authentication
- **JWT:** verify the algorithm against an allowlist (`HS256` or `RS256`, never both), pin the issuer, validate `kid` against a key set served from a stable URL (jwks endpoint cached with TTL). Reject `alg: none` even if the library claims to.
- **Refresh tokens:** rotate on each use; detect re-use of a previously rotated token as compromise (revoke the family). Persist the active token id (`jti`) in MySQL with `(user_id, jti, issued_at, replaced_by)`.
- **Session fixation:** Gin sessions must regenerate the session id on auth-state change (login, role escalation).
- **mTLS for service-to-service:** SPIFFE/SPIRE or a service mesh (Istio/Linkerd). Service identity is `spiffe://td/ns/<ns>/sa/<sa>`, not an IP.

### D.2 Authorization
- **BOLA on every `:id` route.** The check must happen *in the same query* that fetches the resource (`WHERE id = ? AND owner_id = ?`), not as a separate "fetch then check" pattern (race window).
- **BFLA / verb-level RBAC.** Use casbin or oso. Policy lives in a versioned file checked into git; runtime loads from configmap. Avoid `if user.Role == "admin"` scattered through handlers.
- **Tenant scoping.** In multi-tenant, every query carries `tenant_id` in `WHERE`. App-layer filtering alone is **High** severity (a single missing predicate exposes cross-tenant data).
- **404 vs 403 leakage.** When the resource exists but the caller does not own it, return 404 — returning 403 confirms existence (information disclosure).

### D.3 Service-to-service
- Per-service Kong/APISIX consumer (no shared "internal" consumer).
- Internal mesh policy: deny-by-default; explicit `AuthorizationPolicy` per (source service, destination route).
- No implicit-trust internal networks. The phrase "it's behind the VPN" is not a control.

---

## E. Secrets & configuration

### E.1 Detection patterns (always Critical when found in repo)
- Go literals: `apiKey := "AKIA..."`, `password := "..."`.
- Test fixtures: `testdata/.env` containing real production secrets (look for non-`example`/non-`fake` values).
- `_test.go` files using real bearer tokens (string starts with `eyJ` and decodes to a real issuer).
- Helm `values.yaml` with embedded passwords (`db.password: ...`).
- Kong `kong.yml` with literal `keyauth_credentials.key`.

### E.2 Storage hierarchy (preferred → least preferred)
1. **External Secrets Operator + AWS Secrets Manager via IRSA.** Service account assumes a role; ESO syncs Secret Manager values into K8s `Secret`s; pod mounts them. No long-lived keys.
2. **HashiCorp Vault sidecar (Agent Injector).** Vault writes secrets into a tmpfs mount; app reads file. Lease-based.
3. **K8s `Secret` API directly.** Acceptable for low-sensitivity values when `EncryptionConfiguration` is enabled on etcd. Still better than env literals.
4. **Plain `env:` block in `Deployment`.** **Avoid.** Visible in any pod inspection, kubectl describe, audit logs.

### E.3 Rotation
- Signing keys (JWT RS256): rotate every 90 days, support 2-key window during rotation (jwks publishes both).
- DB passwords: rotate every 90 days; ESO + Secrets Manager handles this automatically when the secret is configured for rotation.
- Kong consumer credentials: rotate on personnel change and every 180 days.
- **Same secret across SIT/UAT/PRD is Critical** — a SIT compromise leaks PRD.

---

## F. Logging & observability

### F.1 PII fields that must be masked or omitted
- Borrower name, NIK / national ID, passport number.
- Bank account number, debit/credit card PAN (PCI), CVV (never log under any circumstance).
- Email, phone, address, date of birth.
- KYC document URLs (signed URLs leak the document if logged with the signature).
- Income, credit score, internal risk band (compliance-sensitive).

### F.2 Token / credential fields that must be redacted
- `Authorization`, `Cookie`, `Set-Cookie`, `X-API-Key`, `X-Auth-Token`.
- JWTs anywhere in body.
- Kong `apikey` header.
- Internal mTLS client cert subjects.

### F.3 Implementation pattern (Go `log/slog`)
- Custom `slog.Handler` that walks attrs and replaces values for keys in a redaction set.
- Test the handler: a unit test that logs a struct containing a known-PII field and asserts the formatted output does not contain the value.
- For structured logs, prefer typed attrs (`slog.String("user.id", id)`) over `slog.Any` so the redaction set can match by key, not value pattern.

### F.4 Audit log requirements (lending-specific)
Every credit-decision, KYC-state-change, disbursement, role-change, and underwriter override **must** record:
- Actor identity (subject id + auth method).
- Timestamp (UTC, monotonic-stable).
- Before/after state (or a hash if size is a concern, with the full state stored elsewhere).
- Correlation id linking to the originating request.
- Reason code (machine-readable enum, not free text).

Sink must be append-only (S3 with object lock, or a separate audit DB with no `DELETE` grant). Retention bounded by regulation (typically 5-10 years for lending).

### F.5 Error responses
- Never echo SQL, stack traces, internal hostnames, file paths, or library versions to the client.
- Gin recovery emits a sanitized 500 with correlation-id only; the full stack is logged server-side.
- 4xx errors: machine-readable code + human message; do not leak whether a resource exists (use 404 not 403 for ownership-denied).

---

## G. Database (MySQL / RDS)

### G.1 Query construction
- `database/sql`: always `?` placeholders; never `fmt.Sprintf` into the query.
- `sqlx`: `sqlx.In("WHERE id IN (?)", ids)` is safe; rebind correctly.
- `gorm`: `db.Raw(...).Scan(...)` with `?` placeholders is safe; `db.Exec("...")` with concatenation is **Critical**.
- `sqlboiler`: generated code is parameterized; flag any hand-written `qm.Where(fmt.Sprintf(...))`.

### G.2 Privilege separation
- App user: `SELECT, INSERT, UPDATE, DELETE` on data tables only. **No** `DROP`, `CREATE`, `ALTER`, `GRANT`, `FILE`, `SUPER`.
- Migrator user: `DDL` privileges, used only by the migration job. Different password; rotated separately.
- Read-only reporter user: `SELECT` only on a subset of tables; used by analytics.
- Root: never used by any application. Console-only, MFA-gated.

### G.3 PII at rest
- Column-level encryption for: NIK, passport, bank account number, PAN.
- Two patterns:
  - **Application-layer**: encrypt in Go before insert (`crypto/aes-gcm` with KMS-derived data key). Search-by-encrypted-value requires deterministic encryption (HMAC-derived key) with care.
  - **MySQL function**: `AES_ENCRYPT(value, KEY())` with key from KMS. Less control over key rotation.
- Tokenization is preferred for PAN (PCI-DSS requirement): replace PAN with a token; vault holds the mapping.

### G.4 Transport and at-rest
- `require_secure_transport=ON` on the RDS parameter group.
- Application connection string includes `?tls=true&sslmode=require` (Go driver-specific).
- RDS storage encryption enabled at instance create time (cannot be enabled later without snapshot-restore).
- Backups: automated backups encrypted with the same KMS key; cross-account snapshot copy reviewed against PII export controls.

### G.5 Read replicas (RDS) and CQRS
- Replica lag can let a stale auth state grant access. For role changes, mark a read-after-write window per user, route reads to primary during the window, or invalidate the user's principal cache.
- Replicas in different VPCs/accounts are PII export points; verify the IAM and security-group posture matches the primary.

---

## H. Kubernetes / container

### H.1 `Dockerfile` review (line by line)
- `FROM image:tag` with floating tag (`:latest`, `:1.21`) → flag, recommend digest pin (`@sha256:...`).
- Build secrets via `--build-arg SECRET=...` → flag, recommend BuildKit `--mount=type=secret`.
- `RUN curl https://... | sh` → flag, this fetches code at build time without integrity check.
- `USER root` (or no `USER` directive) → flag, set `USER 65532:65532` or similar high non-root UID.
- `COPY . /app` without `.dockerignore` → flag, will copy `.git`, `.env`, `node_modules`.
- Multi-stage build that copies CA bundles incorrectly → TLS verification fails silently.

### H.2 `Deployment` `securityContext`
Required pod-level:
```yaml
securityContext:
  runAsNonRoot: true
  runAsUser: 65532
  fsGroup: 65532
  seccompProfile: { type: RuntimeDefault }
```
Required container-level:
```yaml
securityContext:
  allowPrivilegeEscalation: false
  readOnlyRootFilesystem: true
  capabilities: { drop: [ALL] }
  runAsNonRoot: true
```
If `readOnlyRootFilesystem: true` causes failures, the fix is an `emptyDir` volume mounted at the writable path, **not** disabling the flag.

### H.3 NetworkPolicy
- Default-deny per namespace:
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata: { name: default-deny }
spec:
  podSelector: {}
  policyTypes: [Ingress, Egress]
```
- Then explicit allow rules per workload (egress to RDS, Kafka, Vault/KMS, DNS).
- Missing default-deny in PRD namespace is **High** by escalation floor.

### H.4 RBAC
- No `verbs: ["*"]` on resources.
- No `ClusterRoleBinding` to a service account unless absolutely necessary; prefer namespaced `RoleBinding`.
- IRSA for AWS access (annotation `eks.amazonaws.com/role-arn`); never long-lived `AWS_ACCESS_KEY_ID` in `Secret`.

### H.5 Resource limits and PDB
- `requests` and `limits` on every container; missing `limits` → noisy-neighbor DoS.
- `PodDisruptionBudget` for any workload that participates in security-critical paths so security-driven restarts (image patches, NetworkPolicy updates) do not destabilize the system.

### H.6 Admission control
- Kyverno or OPA Gatekeeper enforcing the above as cluster policy. The skill prefers a policy that **fails CI** over a recommendation that relies on memory.

---

## I. CI/CD & supply chain

### I.1 GitHub Actions specifics
- **Action SHA pinning:** `uses: actions/checkout@v3` is unsafe (tag can be moved); pin to the commit SHA `actions/checkout@a12a3943b...`. Dependabot can update SHAs.
- **`permissions:` block:** at workflow or job level, set the minimum (`contents: read`, add `id-token: write` only if using OIDC). Default token is too permissive.
- **`pull_request_target` + checkout of head ref:** runs untrusted PR code with repo write secrets — **Critical**.
- **Long-lived cloud credentials in secrets** (`AWS_ACCESS_KEY_ID`): replace with OIDC trust policy (`role-to-assume`). Eliminates rotation burden.
- **Self-hosted runners** processing PRs from forks: **Critical**; use ephemeral or restrict to internal events.

### I.2 SBOM, signing, scanning
- Generate SBOM at build time (Syft, Trivy). Store with the artifact.
- Sign images with cosign (keyless via OIDC preferred). Verify on pull (cosign policy controller, Kyverno `verifyImages`).
- Secret scanning: gitleaks/trufflehog as pre-commit hook **and** CI gate. Pre-commit alone is insufficient (developers can `--no-verify`).
- Dependency scanning: govulncheck, trivy, dependabot. Track unfixed criticals with an SLA.

### I.3 Branch protection and separation of duties
- `main`: required reviews ≥ 1 (≥ 2 for production-impact paths), signed commits, no force-push, status checks required.
- PRD deploy gated on an approver distinct from the PR author (CODEOWNERS-driven approval rule).
- Deploy keys scoped per environment; no shared "deploy" key across SIT/UAT/PRD.

### I.4 SLSA build levels
- Provenance attestation (`slsa-github-generator`) at minimum SLSA Build L2.
- Hermetic builds where feasible (no internet during build, all deps from a vetted proxy).

---

## J. Event-driven (Kafka)

### J.1 Topic ACLs
- Per-service producer/consumer principal (service account). No shared "app" principal.
- `Read`/`Write` ACLs scoped to specific topics and consumer groups.
- `IdempotentWrite` ACL required for idempotent producers.

### J.2 Schema registry
- Compatibility set to `BACKWARD` minimum, `FULL` preferred. Breaking changes require a migration plan (dual-write, consumer upgrade, deprecation window).
- Schemas reviewed for over-collection: an event payload should contain only what consumers need; do not publish full borrower record because "someone might want it later."

### J.3 Idempotency
- Producer: `enable.idempotence=true`, `acks=all`, `max.in.flight.requests.per.connection<=5`.
- Consumer: persist `(topic, partition, offset)` or a domain-level `event_id` in MySQL with a unique constraint. On restart/replay, the constraint causes a duplicate-insert error which the consumer treats as already-processed.
- Critical for **disbursement, credit-decision, KYC-state-change** consumers — without idempotency, at-least-once delivery causes double-disbursement.

### J.4 Poison-pill handling
- DLQ topic with retention bounded by PII policy (do not retain failed events with PII forever).
- Retry budget: max N retries with backoff, then DLQ; not infinite retry that amplifies into a self-DoS.
- Alerting on DLQ growth — a silent DLQ is a hidden integrity failure.

### J.5 Header trust and replay safety
- Headers like `actor-id`, `tenant-id`, `correlation-id` are *informational* unless signed. The consumer must re-derive trust (e.g., look up the actor's permissions in the authoritative service), not trust the header.
- For high-impact events (disbursement), sign the event payload at the producer (HMAC with a per-service key) and verify at the consumer. Replay protection via a `(event_id, ttl)` cache.

### J.6 Encryption
- TLS in transit between brokers and clients (mTLS preferred).
- Broker-side at-rest encryption (KMS).
- Payload-level encryption for PII fields if the cluster is multi-tenant or shared with less-trusted services.

---

## K. Financial / lending data

### K.1 Regulatory inventory (surface based on data classes present)
- **PDPA (Thailand/Singapore/Malaysia variants):** PII rights, consent, breach notification.
- **GDPR (EU borrowers):** lawful basis, right to erasure (drives crypto-shredding plan for events).
- **PCI-DSS:** if PAN is present, scope is the entire system unless tokenized.
- **BOT (Bank of Thailand) / OJK (Indonesia) / MAS (Singapore) / BNM (Malaysia):** lending-specific recordkeeping, decision auditability, fair-lending non-discrimination.
- **AML/KYC (FATF-aligned):** identity verification, sanctions screening, suspicious-activity reporting.

### K.2 KYC artifact handling
- Documents (ID photos, selfies, bank statements) stored in encrypted object store (S3 SSE-KMS) with bucket policy denying public access.
- Access via short-lived (≤ 5 min) pre-signed URLs; never serve documents from app server.
- Pre-signed URL generation logged (actor + document id + expiry) — without logging, a leaked URL is undetectable.
- Document destruction policy: separate from event log (KYC docs may be deletable while the decision record remains).

### K.3 Credit-decision integrity
- Decision events signed (HMAC) at the decisioning service; downstream services verify before acting.
- Bureau-pull responses retained with cryptographic hash for dispute resolution (a borrower disputing an adverse decision must be able to see what data drove it).
- Underwriter manual override requires dual-control: a second approver from a different role recorded in the audit log.

### K.4 Disbursement
- Account-number change on a loan requires re-verification (penny-test or bureau lookup); **never** disburse to a newly-set account without it.
- Payout instructions signed end-to-end (originating system → payment rail).
- No plaintext bank account in logs, Kafka headers, or response bodies.
- Disbursement consumer **must** be idempotent (see J.3).

### K.5 Data subject rights
- For GDPR/PDPA erasure requests where events are immutable: implement crypto-shredding (per-subject KMS key; deletion = key destruction). Plan **before** writing PII to events.
- A "delete user" operation that only removes the row from the user table while events still contain the PII is a regulatory failure, not a fix.

### K.6 Environment hygiene
- PRD data **must not** be copied to SIT/UAT without masking. A `prod_dump.sql` in a non-PRD environment is **Critical**.
- Synthetic data generators (`go-faker`, `mockaroo` exports) are the supported path.
- Detect: any seed file or fixture containing data shapes consistent with PRD (real-looking NIK ranges, bank-account checksums that validate, sequential IDs starting from PRD's id range).
