# Severity & Confidence: Full Rubric

This file is the authoritative rubric. SKILL.md carries a summary; when a finding is borderline, consult this file.

---

## Severity rubric

Severity = `f(Impact, Exploitability, Reachability)`. Assigned per finding, never per file.

### Critical

**Definition.** Direct, immediate financial loss, regulatory breach, or unauthenticated PII exposure of more than one record. The exploit is straightforward (no chained vulnerabilities, no insider access required), and the affected entry point is reachable from the internet or from a low-privilege actor.

**Concrete examples in this stack:**
- Disbursement endpoint missing ownership check (BOLA on disbursement) → attacker authenticated as borrower A triggers payout to borrower B's account.
- Hardcoded RDS root password committed to a public or shared repo.
- Kong route bypassing JWT plugin (route precedence error or plugin attached only globally).
- KYC document URL not signed and accessible by id enumeration.
- Same secret reused across SIT/UAT/PRD (a SIT compromise leaks PRD).
- `pull_request_target` GitHub Actions workflow checking out untrusted PR head with repo-write secrets.
- Public S3 bucket containing PII or KYC documents.
- SQL injection on any internet-reachable endpoint.

### High

**Definition.** Sensitive data exposure, privilege escalation, or integrity loss with preconditions (authenticated user, specific timing, partial knowledge). Limited or single-actor blast radius, but the exploit is realistic.

**Concrete examples:**
- BOLA on `/loans/:id` with predictable sequential ids — authenticated borrower can read other borrowers' loans.
- PII (NIK, bank account) appearing in Kafka headers or in structured logs in PRD.
- Secrets in K8s `Deployment.spec.template.spec.containers[].env[]` plain values.
- CSRF on underwriter override action.
- Replay of `loan.approved` event causes double-disburse (no consumer-side idempotency).
- Refresh-token reuse not detected (no token-family compromise tracking).
- Default-deny `NetworkPolicy` missing in PRD namespace.
- Disbursement / credit-decision / KYC-state-change without an audit log entry.
- mTLS not enforced for internal service-to-service traffic in PRD.
- Cross-tenant query missing `tenant_id` filter at SQL level (relying on app-layer filter only).

### Medium

**Definition.** Defense-in-depth gap, hardening miss, or exploit requiring chaining of multiple weaknesses. Limited blast radius, often partial mitigation already exists.

**Concrete examples:**
- Missing rate-limit on auth endpoint (brute force still requires valid usernames).
- Verbose error responses (stack trace, library version) on a non-PII endpoint.
- `runAsNonRoot` / `readOnlyRootFilesystem` not set on a workload.
- Missing audit log on a non-financial state change (profile email update).
- Floating image tags (`:1.21`) instead of digest pins.
- Container without `requests`/`limits` (DoS vector against neighbors).
- GitHub Actions floating tag (`@v3`) instead of SHA pin.
- API response includes a benign internal field (e.g., `created_by_service: "loan-api"`).
- Custom Lua plugin in Kong with `string.format` of user input (not exploitable to RCE in the observed call sites, but fragile).

### Low

**Definition.** Best-practice deviation, code-quality issue with weak security adjacency, or informational. Not directly exploitable; closes a future hardening gap.

**Concrete examples:**
- Missing security headers (`X-Content-Type-Options`, `Strict-Transport-Security`) on a non-sensitive route.
- Dependency pinned to minor version instead of patch version (no known CVE).
- Comment leaks an internal Jira ID or service URL pattern.
- Inconsistent audit-log format across services (machine-parseable but mixed schema).
- Helpful but unauthenticated `/health/detailed` endpoint exposing build info.

---

## Escalation floors

These rules **override the matrix upward**. If a finding qualifies for both a base severity and a floor, the floor wins.

| Pattern | Floor |
|---|---|
| Hardcoded production credential of any kind in repo, Helm chart, ConfigMap, or `kong.yml` | **Critical** |
| PII (NIK, PAN, bank account, KYC doc URL with signature) in PRD logs | **High** |
| Default-deny `NetworkPolicy` missing in PRD namespace | **High** |
| Disbursement / credit-decision / KYC-state-change path missing audit log | **High** |
| Same secret across SIT/UAT/PRD (any service) | **Critical** |
| PRD data copied to SIT/UAT without masking | **Critical** |
| Public bucket / unauthenticated endpoint serving PII | **Critical** |
| RBAC with `verbs: ["*"]` on a `Secret` or `ConfigMap` resource in PRD | **High** |
| `pull_request_target` checking out PR head with secrets exposure | **Critical** |

---

## Severity recomputation under chaining

When two findings combine to enable an exploit neither enables alone, the **chained severity** is the higher of:
- The product of the individual severities (Medium + Medium can chain to High).
- The severity that the chained outcome warrants on its own.

Document the chain explicitly: list the chained findings by id, the combined attack path, and the chained severity.

**Example:** A Medium "unauthenticated `/health/detailed` exposes service version" + a Medium "outdated library X with public CVE-YYYY-NNNN" can chain to High because the disclosed version enables targeted exploitation. Both individual findings remain at Medium; an additional High chain finding is created with cross-references.

---

## Confidence rubric

Confidence = the skill's certainty that the finding is **real and exploitable in this codebase**, not theoretical.

### High

- The vulnerable pattern is present in the artifact and cited with file:line or YAML key.
- The exploit path is reasoned end-to-end (actor → entry point → control gap → impact).
- No plausible compensating control is visible elsewhere in the artifact.
- The skill can name what the existing code does that is wrong **and** what a correct version would look like.

### Medium

- The pattern is present but a compensating control may exist outside the artifact.
- Examples:
  - "BOLA-shaped handler in Gin, but there may be a Kong consumer-restriction plugin not shown."
  - "Plain `env:` secret in `Deployment`, but External Secrets Operator may be syncing it from Secrets Manager — config not provided."
- The skill must explicitly name **what would need to be true to invalidate the finding** (e.g., "If `kong.yml` shows `acl` plugin scoped to the borrower's group, the finding does not apply.").

### Low

- The skill suspects a problem from shape, but lacks enough context.
- Examples:
  - User pasted a single handler; the surrounding middleware is unknown.
  - K8s Deployment without the corresponding NetworkPolicy file.
  - OpenAPI spec without the implementing service.
- The skill states the **exact artifact needed to upgrade** confidence (e.g., "Provide `internal/auth/middleware.go` and the route group registration to confirm.").

---

## Confidence rules

1. **Never publish Critical or High at Low confidence without an explicit `[needs verification]` label** in the finding header. Example:
   ```
   **Confidence:** Low [needs verification — provide router registration]
   ```

2. **Do not fabricate to raise confidence.** Hallucinated APIs, fictional Gin middleware, invented Kong plugin names, or made-up CWE/API/ASVS identifiers are a quality failure, not a Low-confidence finding. The skill must withhold rather than fabricate.

3. **Confidence requires a citation.** If the skill cannot point to a file:line, YAML key, or specific config block, confidence is **Medium maximum**. A purely descriptive finding with no anchor is at Medium.

4. **Confidence is independent of severity.** A High-severity finding can be Low confidence (suspicious shape, need more context). A Critical finding at Low confidence is published with `[needs verification]` because the worst-case warrants surfacing even when uncertain.

5. **Confidence can be downgraded by the user.** If the user provides context that contradicts the finding ("we have a Kong ACL plugin you didn't see"), update the finding to Low confidence or withdraw it; do not argue.

---

## Worked severity / confidence calls

| Scenario | Severity | Confidence | Reason |
|---|---|---|---|
| Gin handler `loan.GetByID(c.Param("id"))` with no owner check, surrounding code shows authentication middleware but no authz | High | High | BOLA pattern clear; auth ≠ authz; no compensating control visible |
| Same handler, but only the handler was pasted; no context | High | Medium | Pattern is the same; cannot rule out a casbin policy enforced upstream |
| `kong.yml` with `key-auth` plugin attached globally and a public route declared without per-route plugin pinning, but `acl` plugin presence unknown | Critical | Medium | Route precedence can shadow auth; need to see the full plugin chain |
| Hardcoded password `"hunter2"` in `_test.go` fixture | Critical | High | Floor: hardcoded credential. Even in test, it leaks via repo history |
| `env: { name: DB_PASSWORD, value: "..." }` in `Deployment` | Critical | High | Floor: hardcoded credential in PRD config |
| `runAsNonRoot` not set on a service that processes loan applications | Medium | High | Hardening miss; container escape less likely than direct app vulns but increases blast radius |
| `pull_request_target` workflow checking out PR head and running `npm install` with secrets | Critical | High | Direct path to repo compromise via fork PR |
| Suspected race condition on credit-limit check, but only one file shown | High | Low | Pattern fits race; need transaction boundary view to confirm; mark `[needs verification]` |
