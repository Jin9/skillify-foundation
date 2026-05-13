# Worked Examples

Twelve examples covering the eleven taxonomy areas plus refusal handling. Each follows the Finding Format from SKILL.md. All data is synthetic; secrets shown are obviously fake.

Read the example matching the artifact in front of you before drafting your first finding.

---

## Example 1 — Gin BOLA on loan endpoint (Application + API + AuthZ)

### [SEV-1] Loan detail handler missing ownership check

**Severity:** Critical
**Confidence:** High
**Category:** AuthN/AuthZ
**Secondary tags:** API security
**Standards:** CWE-639, API1:2023, ASVS V4.2.1
**Affected area:** internal/loan/handler.go:42-58
**Environment scope:** all

**Asset at risk:** Borrower PII, loan disbursement state, KYC linkage.
**Trust boundary crossed:** Internet → Gin handler.
**Threat (STRIDE + abuse case):** Information Disclosure / Elevation — borrower A iterates `/loans/{id}` to read borrower B's loan including bank account and KYC status. Sequential ids make enumeration trivial.
**Attack scenario:**
  1. Borrower A authenticates and obtains JWT.
  2. Borrower A calls `GET /loans/9001`, then `9002`, ... (B's id).
  3. Handler queries by id with no `borrower_id` predicate.
  4. Server returns B's record.
**Risk:** Confidentiality breach across tenant boundary; reputational and regulatory exposure.
**Business impact:** PDPA/GDPR per-record fine, regulator notification within 72h, class-action exposure.

**Evidence:**
```go
// internal/loan/handler.go
func (h *Handler) Get(c *gin.Context) {
    id := c.Param("id")
    loan, err := h.repo.GetByID(c.Request.Context(), id)
    if err != nil { c.AbortWithStatus(http.StatusInternalServerError); return }
    c.JSON(200, loan)
}
```

**Recommended fix:**
```go
func (h *Handler) Get(c *gin.Context) {
    id := c.Param("id")
    subject := auth.SubjectFrom(c) // (user_id, tenant_id, roles)
    loan, err := h.repo.GetByIDForOwner(c.Request.Context(), id, subject)
    if err != nil {
        switch {
        case errors.Is(err, loan.ErrNotFound), errors.Is(err, loan.ErrNotOwned):
            c.AbortWithStatus(http.StatusNotFound) // do not leak existence
        default:
            c.AbortWithStatus(http.StatusInternalServerError)
        }
        return
    }
    c.JSON(200, loan.ToResponse(subject)) // role-filtered projection
}
```

In the repository:
```go
func (r *MySQLRepo) GetByIDForOwner(ctx context.Context, id string, s auth.Subject) (Loan, error) {
    const q = `SELECT ... FROM loans WHERE id = ? AND borrower_id = ? AND tenant_id = ?`
    row := r.db.QueryRowContext(ctx, q, id, s.UserID, s.TenantID)
    // ...
}
```

Why this works: the ownership predicate is part of the SQL `WHERE` (single round-trip, race-free). 404 (not 403) prevents existence-oracle. Tenant scope applied at SQL level, not app-layer filter.
Trade-offs: callers passing only id must migrate; integration tests must be updated.

**Safer pattern:** `casbin/casbin v2.x` policy enforcement at handler level, repository helper enforcing ownership at query level.

**Validation:**
- Test: `TestLoanGet_OtherBorrower_Returns404` asserts 404 when caller != owner; `TestLoanGet_Owner_ReturnsLoan` asserts 200 for the owner.
- Static: `semgrep` rule `loan-get-without-owner` matches calls to `repo.GetByID` from handlers, fails CI.
- Manual: with two test JWTs, `curl /loans/{B-id}` from A's JWT returns 404.
- Regression guard: semgrep rule above in `pre-commit` and CI.
- Observability: counter `auth.deny.bola{route="/loans/:id"}` with alert if zero for 24h post-deploy.
- Rollout: SIT → UAT → 10% canary PRD → full PRD; rollback if 4xx rate spikes >2σ.

**Residual risk:** Existence-oracle via timing if owner-miss is materially faster than owner-hit. Track p99 latency and equalize via constant-time padding if delta >5ms.

---

## Example 2 — Kafka non-idempotent disbursement consumer (Event-driven + Lending)

### [SEV-2] Disbursement consumer lacks idempotency key, at-least-once delivery causes double-pay

**Severity:** Critical
**Confidence:** High
**Category:** Event-driven (Kafka)
**Secondary tags:** Financial/lending data
**Standards:** CWE-362, ASVS V11.1.4 (idempotency), NIST SSDF PW.5
**Affected area:** services/disbursement/consumer.go:31-60
**Environment scope:** all

**Asset at risk:** Disbursement integrity (direct money loss).
**Trust boundary crossed:** Kafka topic `loan.approved` → disbursement-svc.
**Threat (STRIDE + abuse case):** Tampering / Elevation — replay of `loan.approved` (broker rebalance, consumer restart, deliberate replay) triggers a second disbursement to the same borrower.
**Attack scenario:**
  1. Consumer commits offset *after* the disbursement HTTP call to the payment rail.
  2. Consumer crashes between payment-rail call and offset commit.
  3. On restart, the same event is re-delivered.
  4. Disbursement fires again. Payment rail has no replay protection of its own.
**Risk:** Integrity / Availability of the lending book; direct financial loss.
**Business impact:** Per-event loss equal to disbursement amount; reconciliation cost; audit-finding from BOT/OJK on lending controls.

**Evidence:**
```go
// services/disbursement/consumer.go
func (c *Consumer) Handle(ctx context.Context, msg *kafka.Message) error {
    var ev LoanApproved
    if err := proto.Unmarshal(msg.Value, &ev); err != nil { return err }
    // No idempotency check.
    if err := c.payments.Disburse(ctx, ev.LoanID, ev.Amount, ev.Account); err != nil {
        return err
    }
    return c.commit(msg)
}
```

**Recommended fix:** persist a dedup record in MySQL with a unique constraint on `event_id`, in the same transaction that records the disbursement intent. Treat duplicate-key as already-processed.

```go
func (c *Consumer) Handle(ctx context.Context, msg *kafka.Message) error {
    var ev LoanApproved
    if err := proto.Unmarshal(msg.Value, &ev); err != nil { return err }

    if err := verifyHMAC(ev, msg.Headers); err != nil {
        return c.toDLQ(msg, "signature_invalid")
    }

    err := c.tx.Do(ctx, func(tx *sql.Tx) error {
        if _, err := tx.ExecContext(ctx,
            `INSERT INTO disbursement_events (event_id, loan_id, status) VALUES (?, ?, 'PENDING')`,
            ev.EventID, ev.LoanID,
        ); err != nil {
            if isDuplicateKey(err) { return errAlreadyProcessed }
            return err
        }
        return nil
    })
    if errors.Is(err, errAlreadyProcessed) { return c.commit(msg) } // safe replay
    if err != nil { return err }

    if err := c.payments.Disburse(ctx, ev.LoanID, ev.Amount, ev.Account); err != nil {
        // mark failed in tx; do NOT commit offset; will retry
        return err
    }

    if err := c.markCompleted(ctx, ev.EventID); err != nil { return err }
    return c.commit(msg)
}
```

DDL:
```sql
CREATE TABLE disbursement_events (
  event_id    CHAR(36) PRIMARY KEY,
  loan_id     CHAR(36) NOT NULL,
  status      ENUM('PENDING','COMPLETED','FAILED') NOT NULL,
  created_at  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  UNIQUE KEY uk_loan_pending (loan_id, status) -- one PENDING per loan
);
```

Why this works: the unique constraint is the durable dedup state. The `INSERT` either succeeds (first time) or fails (replay), and the consumer treats both deterministically.
Trade-offs: requires producer to set a stable `event_id` (not a random per-delivery id). If `event_id` is not stable, derive from `(loan_id, decision_correlation_id)`.
Migration: deploy DDL first, then the consumer; replay-safe because the new consumer is a superset of the old.

**Safer pattern:** `segmentio/kafka-go` consumer with manual offset commit + MySQL idempotency table. Verify `enable.idempotence=true` at the producer.

**Validation:**
- Test: `TestDisbursement_Replay_DoesNotDoublePay` produces the same event twice; asserts payment-rail mock is called once.
- Test: `TestDisbursement_FreshEvent_Disburses` asserts a new event causes one disbursement.
- Static: `gosec` clean; `semgrep` rule `kafka-consumer-without-dedup` flags `Handle` funcs that lack a dedup query.
- Manual: in SIT, replay a known event via `kafkactl produce ... --key <event_id>`; observe one disbursement.
- Regression guard: integration test runs in CI against an embedded Kafka.
- Observability: counter `disbursement.replay_skipped`; alert if `disbursement.duplicate_paid` ever non-zero (should be impossible).
- Rollout: SIT → UAT → canary PRD; rollback if disbursement-rail error rate increases.

**Residual risk:** If two parallel consumers in the same group both pull the event before either commits, MySQL serialization isolation handles the race; ensure `tx_isolation = REPEATABLE-READ` on RDS (default).

---

## Example 3 — Hardcoded RDS password in Helm values (Secrets)

### [SEV-3] Production database password committed in `values.yaml`

**Severity:** Critical
**Confidence:** High
**Category:** Secrets & configuration
**Secondary tags:** None
**Standards:** CWE-798, CWE-256, ASVS V2.10, NIST SSDF PS.1
**Affected area:** charts/loan-api/values-prd.yaml:24
**Environment scope:** PRD (also leaks SIT/UAT if same value)

**Asset at risk:** RDS database (full read/write of borrower PII, loan state).
**Trust boundary crossed:** Source-control → anyone with repo access (and Git history forever).
**Threat (STRIDE + abuse case):** Spoofing / Elevation — anyone with repo read can read PRD DB.
**Attack scenario:**
  1. Engineer with read access (or compromised CI runner, or any historical contributor) reads `values-prd.yaml`.
  2. Connects to RDS via bastion or VPN with the credential.
  3. Exfiltrates borrower data.
**Risk:** Total confidentiality + integrity breach of lending data.
**Business impact:** Catastrophic regulatory and financial; mandatory regulator notification; potential license review.

**Evidence:**
```yaml
# charts/loan-api/values-prd.yaml
db:
  host: prd-loan.cluster-xxxxxx.ap-southeast-1.rds.amazonaws.com
  user: loan_app
  password: "Pa$$w0rd-NotReal-Example-Only"   # Critical
```

**Recommended fix:** use External Secrets Operator (ESO) backed by AWS Secrets Manager via IRSA. Remove the literal; rotate the credential immediately (it must be considered compromised given git history).

`charts/loan-api/values-prd.yaml`:
```yaml
db:
  host: prd-loan.cluster-xxxxxx.ap-southeast-1.rds.amazonaws.com
  user: loan_app
  passwordSecretRef:
    name: loan-api-db
    key: password
```

`charts/loan-api/templates/externalsecret.yaml`:
```yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: loan-api-db
spec:
  refreshInterval: 15m
  secretStoreRef:
    name: aws-secrets-manager
    kind: ClusterSecretStore
  target:
    name: loan-api-db
  data:
    - secretKey: password
      remoteRef:
        key: prd/loan-api/db
        property: password
```

`charts/loan-api/templates/serviceaccount.yaml`:
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: loan-api
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::000000000000:role/loan-api-prd
```

Rotation: enable Secrets Manager rotation Lambda for `prd/loan-api/db` with 60-day cadence.

Why this works: the secret never touches git; ESO syncs it into a K8s `Secret` from Secrets Manager; the pod assumes a role via IRSA (no static IAM keys). Rotation is automatic.
Trade-offs: ESO + Secrets Manager monthly cost; one extra cluster component to operate.
Migration: rotate the credential first (the leaked one is burned), deploy ESO chart, then deploy the new values; verify before removing the literal.

**Safer pattern:** External Secrets Operator (`external-secrets.io v0.9.x`) + AWS Secrets Manager + IRSA.

**Validation:**
- Test: deployment fails (pod CrashLoopBackOff with clear error) when `loan-api-db` secret is absent — proves no fallback to literal.
- Static: `gitleaks` rule `aws-rds-password` matches the previous literal; CI fails. Confirm `gitleaks --redact` shows the leak in history (history must be considered compromised).
- Manual: `kubectl get externalsecret loan-api-db -n prd -o jsonpath='{.status.conditions[?(@.type=="Ready")].status}'` returns `True`; `aws secretsmanager describe-secret --secret-id prd/loan-api/db` shows rotation enabled. Do not print secret values.
- Regression guard: `gitleaks` pre-commit + CI; Kyverno policy `disallow-plain-env-passwords` rejects `Deployment` with `env[].value` matching password-like keys.
- Observability: alert on `ExternalSecret` sync failures; alert on stale-secret age > 90 days.
- Rollout: rotate first → deploy chart → smoke test SIT → UAT → PRD; rollback by re-applying previous chart (with fresh credential, never the leaked one).

**Residual risk:** The leaked credential remains in git history; treat as compromised forever and never re-use. Consider history rewrite only if no external clones exist (rarely true).

---

## Example 4 — Kong route bypassing JWT plugin (API security)

### [SEV-4] Public Kong route shadows authenticated route via path precedence

**Severity:** Critical
**Confidence:** Medium [needs verification — confirm full plugin chain not shown]
**Category:** API security (Kong)
**Secondary tags:** None
**Standards:** CWE-862, API2:2023, API5:2023
**Affected area:** infra/kong/kong.yml: routes[3]
**Environment scope:** PRD

**Asset at risk:** Loan-api endpoints behind expected auth.
**Trust boundary crossed:** Internet → loan-api.
**Threat (STRIDE + abuse case):** Spoofing / Elevation — unauthenticated request to `/v1/loans/internal/health` is matched by a permissive route that proxies to loan-api with no JWT plugin attached.
**Attack scenario:**
  1. Attacker requests `/v1/loans/internal/health`.
  2. Kong matches the broader `paths: ["/v1/loans/internal"]` route, which has `jwt` plugin only as a *global* default but no per-route binding.
  3. A more specific route `/v1/loans/internal/admin` defined later carries the plugin, but precedence and plugin scoping mean the health route reaches loan-api unauthenticated.
  4. The internal route exposes data not expected on a public path.

**Evidence:**
```yaml
# infra/kong/kong.yml (excerpt)
plugins:
  - name: jwt
    # global plugin — applies only when no per-route plugin exists
services:
  - name: loan-api
    url: http://loan-api.svc.cluster.local
    routes:
      - name: public-health
        paths: ["/v1/loans/internal"]   # too broad
        # no plugins[] — falls back to global, but precedence with the next route is fragile
      - name: admin
        paths: ["/v1/loans/internal/admin"]
        plugins:
          - name: jwt
```

**Recommended fix:** pin auth plugins **per-route**; narrow the public route's path; add `request-size-limiting` and `rate-limiting` per-route.

```yaml
services:
  - name: loan-api
    url: http://loan-api.svc.cluster.local
    routes:
      - name: public-health
        paths: ["/v1/health"]            # narrow + no PII surface
        plugins:
          - name: rate-limiting
            config: { minute: 60, policy: local }
          - name: request-size-limiting
            config: { allowed_payload_size: 1 }
      - name: loan-internal
        paths: ["/v1/loans/internal"]
        strip_path: false
        plugins:
          - name: jwt
            config: { key_claim_name: kid, claims_to_verify: [exp, nbf] }
          - name: acl
            config: { allow: ["loan-internal-callers"] }
          - name: rate-limiting
            config: { minute: 600 }
```

Why this works: per-route plugin binding makes auth declarative and impossible to bypass via precedence rules. ACL plugin enforces caller identity beyond JWT presence.
Trade-offs: more verbose config; mitigated by Kong's `_format_version` workflow.
Migration: deploy declarative diff; verify with `kong config parse` and `deck diff`.

**Safer pattern:** `kong-gateway >= 3.x` with declarative config + `deck` tooling; `acl` plugin scoping.

**Validation:**
- Test: `curl https://api.example/v1/loans/internal/health` returns 401 without JWT; returns 200 with valid JWT in `loan-internal-callers` group.
- Static: `kong config parse infra/kong/kong.yml` succeeds; custom CI step asserts every route under `/v1/loans/` has a `jwt` plugin in its `plugins[]`.
- Manual: pen-test the full path matrix against SIT.
- Regression guard: CI script (jq over `kong.yml`) fails build if any route under sensitive prefixes lacks `jwt` per-route.
- Observability: Kong access log shows `jwt` plugin executed for every `/v1/loans/*`; alert on 200 responses without `X-Consumer-Username` header.
- Rollout: SIT → UAT → PRD; rollback via `deck sync` to previous declarative state.

**Residual risk:** Plugin order across multiple routes still matters; the CI guard above is necessary but not sufficient if a future route adds `request-transformer` that strips headers loan-api trusts.

---

## Example 5 — K8s deployment running as root with no NetworkPolicy (Kubernetes)

### [SEV-5] Pod runs as root with writable rootfs and no namespace-level egress restriction

**Severity:** High
**Confidence:** High
**Category:** Kubernetes / container
**Secondary tags:** None
**Standards:** CWE-250, CIS K8s 5.2.5/5.2.6, NIST SSDF PW.6
**Affected area:** k8s/loan-api/deployment.yaml; missing networkpolicy.yaml in `prd-loan` namespace
**Environment scope:** PRD

**Asset at risk:** Pod blast radius (lateral movement, secret theft from neighboring pods, exfiltration to internet).
**Trust boundary crossed:** Pod ↔ pod within namespace; pod ↔ internet egress.
**Threat (STRIDE + abuse case):** Elevation / Information Disclosure — RCE in loan-api leads to root inside container, mounts `/proc`, reads service-account token, calls K8s API; with no egress NetworkPolicy, exfiltrates to attacker-controlled host.

**Evidence:**
```yaml
# k8s/loan-api/deployment.yaml (excerpt)
spec:
  template:
    spec:
      containers:
        - name: loan-api
          image: registry.example/loan-api:1.21    # floating tag
          # no securityContext
```
No `NetworkPolicy` resource present in the namespace.

**Recommended fix:**

`deployment.yaml`:
```yaml
spec:
  template:
    spec:
      automountServiceAccountToken: false
      securityContext:
        runAsNonRoot: true
        runAsUser: 65532
        fsGroup: 65532
        seccompProfile: { type: RuntimeDefault }
      containers:
        - name: loan-api
          image: registry.example/loan-api@sha256:<digest>   # pinned
          imagePullPolicy: IfNotPresent
          securityContext:
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
            capabilities: { drop: [ALL] }
          volumeMounts:
            - { name: tmp, mountPath: /tmp }
          resources:
            requests: { cpu: 200m, memory: 256Mi }
            limits:   { cpu: 1,    memory: 512Mi }
      volumes:
        - { name: tmp, emptyDir: {} }
```

`networkpolicy.yaml`:
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata: { name: default-deny, namespace: prd-loan }
spec:
  podSelector: {}
  policyTypes: [Ingress, Egress]
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata: { name: loan-api-egress, namespace: prd-loan }
spec:
  podSelector: { matchLabels: { app: loan-api } }
  policyTypes: [Egress]
  egress:
    - to:
        - namespaceSelector: { matchLabels: { name: kube-system } }
          podSelector: { matchLabels: { k8s-app: kube-dns } }
      ports: [ { protocol: UDP, port: 53 } ]
    - to:
        - ipBlock: { cidr: 10.0.10.0/28 }   # RDS subnet
      ports: [ { protocol: TCP, port: 3306 } ]
    - to:
        - ipBlock: { cidr: 10.0.20.0/28 }   # Kafka brokers
      ports: [ { protocol: TCP, port: 9093 } ]
```

Why this works: non-root + read-only rootfs + dropped capabilities raise the bar for container escape; default-deny + explicit egress prevents exfiltration even if RCE occurs.
Trade-offs: services that write to `/` will fail until they mount writable `emptyDir`; egress rules require the operator to know all dependencies upfront.
Migration: roll deployment first (writable-fs failures appear in SIT), then NetworkPolicy in audit mode if available, then enforce.

**Safer pattern:** `kubernetes >= 1.25` Pod Security Standards `restricted`; `Kyverno` policies `require-non-root`, `disallow-host-namespaces`, `require-network-policy-default-deny`.

**Validation:**
- Test: pod admission rejected by Kyverno without `runAsNonRoot: true`.
- Static: `kube-linter` rules `run-as-non-root`, `read-only-root-fs`, `no-network-policy`; `checkov` policies `CKV_K8S_8`, `CKV_K8S_22`, `CKV_K8S_23`.
- Manual: `kubectl exec` into the running pod, attempt `touch /foo` → expect "Read-only file system"; attempt `curl https://attacker.example` → expect timeout.
- Regression guard: Kyverno audit + enforce in cluster; CI runs `kube-linter`.
- Observability: alert on `NetworkPolicyDenied` events spiking (which means a legitimate dep was missed) and on any pod scheduled in `prd-loan` without matching `default-deny`.
- Rollout: enforce Kyverno first in SIT; resolve any blocked pods; promote to UAT then PRD.

**Residual risk:** Some egress rules (e.g., Vault, KMS) may need additional CIDRs as architecture evolves; the NetworkPolicy is now a maintained surface.

---

## Example 6 — MySQL string-concat in admin tool (Database)

### [SEV-6] Internal admin endpoint builds SQL by string concatenation

**Severity:** Critical
**Confidence:** High
**Category:** Database (MySQL)
**Secondary tags:** None
**Standards:** CWE-89, ASVS V5.3.4, NIST SSDF PW.4
**Affected area:** internal/admin/search.go:18-29
**Environment scope:** all

**Asset at risk:** Full MySQL contents (PII, loan state, audit log).
**Trust boundary crossed:** Internal admin user → MySQL.
**Threat (STRIDE + abuse case):** Tampering / Information Disclosure / Elevation — admin (or anyone who phishes an admin's session) injects SQL via the search field, exfiltrates the borrower table or escalates to a `loans.borrower_id = (SELECT ...)` bypass elsewhere. "Internal" tools are not safe.

**Evidence:**
```go
// internal/admin/search.go
func (h *AdminHandler) Search(c *gin.Context) {
    name := c.Query("name")
    q := fmt.Sprintf("SELECT id, name FROM borrowers WHERE name LIKE '%%%s%%'", name)
    rows, err := h.db.QueryContext(c, q)
    // ...
}
```

**Recommended fix:**
```go
func (h *AdminHandler) Search(c *gin.Context) {
    name := strings.TrimSpace(c.Query("name"))
    if len(name) < 2 || len(name) > 100 {
        c.AbortWithStatusJSON(400, gin.H{"error": "name length"})
        return
    }
    const q = `SELECT id, name FROM borrowers WHERE name LIKE CONCAT('%', ?, '%') LIMIT 100`
    rows, err := h.db.QueryContext(c.Request.Context(), q, name)
    // ...
}
```

Plus DB-user privilege separation:
```sql
-- Migrator only
CREATE USER 'loan_migrator'@'%' IDENTIFIED BY '...';
GRANT ALL ON loans_db.* TO 'loan_migrator'@'%';

-- App user
CREATE USER 'loan_app'@'%' IDENTIFIED BY '...';
GRANT SELECT, INSERT, UPDATE, DELETE ON loans_db.* TO 'loan_app'@'%';
REVOKE DROP, CREATE, ALTER, GRANT OPTION, FILE ON *.* FROM 'loan_app'@'%';
```

Why this works: parameterized query rejects injection; `LIMIT` caps blast radius even if a wildcard slips; least-privilege user means even successful injection cannot `DROP` or read other schemas.
Trade-offs: none.
Migration: ship the Go change; rotate any DB user that currently has excess privileges.

**Safer pattern:** `database/sql` with `?` placeholders; `sqlx` `IN()` rebind for variadic; never `fmt.Sprintf` into a query.

**Validation:**
- Test: `TestAdminSearch_SQLi_Rejected` sends `name=' OR 1=1 --` and asserts the query returns matching `name LIKE '% OR 1=1 --%'` results (i.e., literal match), not the entire table.
- Static: `gosec G201` (SQL string formatting) flags any remaining `fmt.Sprintf` in query construction; CI fails.
- Manual: pen-test the admin tool against SIT.
- Regression guard: `gosec G201` in CI; semgrep rule `db-query-string-format` covering `db.QueryContext` / `db.ExecContext` first arg.
- Observability: MySQL audit log shows queries; alert on queries containing `OR 1=1` or `UNION SELECT` from app user.
- Rollout: SIT → UAT → PRD.

**Residual risk:** `LIKE %` pattern with attacker-controlled wildcard can still cause expensive scans; rate-limit the admin search and require approval for full scans.

---

## Example 7 — PII in structured logs (Logging & observability)

### [SEV-7] Borrower NIK and bank account written to stdout in PRD

**Severity:** High
**Confidence:** High
**Category:** Logging & observability
**Secondary tags:** None
**Standards:** CWE-532, ASVS V7.1.1, NIST SSDF PW.4 (Floor: PII in PRD logs)
**Affected area:** internal/loan/service.go:88
**Environment scope:** PRD

**Asset at risk:** Borrower PII at rest in CloudWatch / log aggregator.
**Trust boundary crossed:** App → log sink (often a less-controlled access surface).
**Threat (STRIDE + abuse case):** Information Disclosure — a log-aggregator breach or read-only access to logs (granted broadly to engineers) leaks PII. Even without breach, log access is rarely audited at row level.

**Evidence:**
```go
// internal/loan/service.go
log.Info("loan application received",
    "borrower_nik", req.NIK,
    "bank_account", req.BankAccount,
    "amount", req.Amount,
)
```

**Recommended fix:** central redacting `slog.Handler`; log identifiers (loan id, borrower id) only.

```go
// pkg/log/redactor.go
type RedactingHandler struct {
    inner slog.Handler
    sensitiveKeys map[string]struct{}
}

func (h RedactingHandler) Handle(ctx context.Context, r slog.Record) error {
    var newAttrs []slog.Attr
    r.Attrs(func(a slog.Attr) bool {
        if _, ok := h.sensitiveKeys[a.Key]; ok {
            newAttrs = append(newAttrs, slog.String(a.Key, "[REDACTED]"))
        } else {
            newAttrs = append(newAttrs, a)
        }
        return true
    })
    out := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
    out.AddAttrs(newAttrs...)
    return h.inner.Handle(ctx, out)
}

var sensitive = map[string]struct{}{
    "borrower_nik": {}, "bank_account": {}, "pan": {},
    "authorization": {}, "cookie": {}, "x-api-key": {},
    "email": {}, "phone": {}, "address": {}, "dob": {},
}
```

Caller:
```go
log.Info("loan application received",
    "loan_id", loan.ID,
    "borrower_id", req.BorrowerID,
    "amount_minor", req.AmountMinor,
)
```

Why this works: redaction is centralized; new caller sites cannot accidentally log a sensitive key; the redaction set is testable.
Trade-offs: the redaction set must be kept current as fields are added.
Migration: deploy handler library; sweep callers; semgrep rule prevents regression.

**Safer pattern:** `log/slog` (Go 1.21+) with `RedactingHandler`; never use raw `fmt.Sprintf` for logging.

**Validation:**
- Test: `TestRedactingHandler_NIK_NotEmitted` logs `slog.String("borrower_nik", "1234567890123")`, asserts the buffer does not contain `1234567890123`.
- Test: `TestRedactingHandler_LoanID_Emitted` asserts non-sensitive keys pass through.
- Static: semgrep rule `log-sensitive-key` matches `log.Info|Error|Warn` calls with `nik|bank_account|pan|authorization|cookie|x-api-key|email|phone|dob` in the args.
- Manual: in SIT, send a request with a known NIK; tail the log and grep for the value — must return nothing.
- Regression guard: semgrep rule above in pre-commit and CI.
- Observability: counter `log.redactions_emitted` (must be > 0; if zero, redactor not wired). Alert if PRD log line contains 13-digit number consistent with NIK pattern (lossy guard).
- Rollout: handler change is backward-compatible; deploy SIT → UAT → PRD.

**Residual risk:** Redactor only catches keys it knows about; structured logging discipline must be maintained. Free-text log messages (`log.Info("borrower 1234567... applied")`) bypass the handler — semgrep rule must also flag interpolated PII shapes in message strings.

---

## Example 8 — GitHub Actions OIDC misconfiguration (CI/CD)

### [SEV-8] Workflow uses long-lived AWS keys and floating action versions

**Severity:** High
**Confidence:** High
**Category:** CI/CD & supply chain
**Secondary tags:** None
**Standards:** CWE-798, CISA Defending CI/CD, SLSA L2
**Affected area:** .github/workflows/deploy.yml
**Environment scope:** all

**Asset at risk:** AWS account write access; ability to push images to ECR; ability to modify infrastructure.
**Trust boundary crossed:** GitHub repo → AWS account.
**Threat (STRIDE + abuse case):** Elevation — leaked or stolen `AWS_ACCESS_KEY_ID` (via fork PR, third-party action compromise, log leak) gives long-lived AWS access. Floating `actions/checkout@v3` can be moved to a malicious commit.

**Evidence:**
```yaml
# .github/workflows/deploy.yml (excerpt)
on: [push, pull_request_target]
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with: { ref: ${{ github.event.pull_request.head.sha }} }
      - uses: aws-actions/configure-aws-credentials@v2
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: ap-southeast-1
      - run: make deploy
```

**Critical sub-issue:** `pull_request_target` + `checkout` of PR head + secrets exposure = arbitrary code execution as a privileged actor on every fork PR. **This sub-issue is independently Critical.**

**Recommended fix:**
```yaml
on:
  push: { branches: [main] }
  workflow_dispatch:
permissions:
  contents: read
  id-token: write   # for OIDC
jobs:
  deploy:
    runs-on: ubuntu-latest
    environment: production   # gate with required reviewers
    steps:
      - uses: actions/checkout@a12a3943b...  # SHA pin
      - uses: aws-actions/configure-aws-credentials@b1234567...
        with:
          role-to-assume: arn:aws:iam::000000000000:role/github-deploy-loan-api
          aws-region: ap-southeast-1
      - run: make deploy
```

AWS IAM trust policy on `github-deploy-loan-api`:
```json
{
  "Effect": "Allow",
  "Principal": { "Federated": "arn:aws:iam::000000000000:oidc-provider/token.actions.githubusercontent.com" },
  "Action": "sts:AssumeRoleWithWebIdentity",
  "Condition": {
    "StringEquals": {
      "token.actions.githubusercontent.com:aud": "sts.amazonaws.com",
      "token.actions.githubusercontent.com:sub": "repo:org/loan-api:ref:refs/heads/main"
    }
  }
}
```

Why this works: OIDC eliminates long-lived secrets; the trust condition pins the workflow to `main` branch of a specific repo, so even a compromised fork cannot assume the role; SHA-pinned actions cannot be moved.
Trade-offs: requires AWS admin to provision the OIDC provider once; extra IAM role per repo.
Migration: provision role, update workflow on a feature branch, validate via `workflow_dispatch`, rotate the leaked key.

**Safer pattern:** `aws-actions/configure-aws-credentials@<sha>` with OIDC; GitHub Environments + required reviewers for production deploys; Dependabot for action SHA updates.

**Validation:**
- Test: workflow on a feature branch fails to assume the role (sub condition mismatch); workflow on `main` succeeds.
- Static: `actionlint` for workflow syntax; custom CI check that `permissions:` is set at workflow or job scope and `pull_request_target` is absent.
- Manual: open a fork PR, confirm secrets are unavailable; confirm `main` push deploys.
- Regression guard: CI rule rejects new `secrets.AWS_ACCESS_KEY_ID` references; CI rule rejects floating `@v\d+` action references.
- Observability: AWS CloudTrail `AssumeRoleWithWebIdentity` events from the OIDC provider; alert on assumption from unexpected `sub` claim.
- Rollout: rotate the leaked key first; deploy workflow on a feature branch; switch over.

**Residual risk:** Self-hosted runners (if used) still need ephemeral or restricted-event configuration; verify separately.

---

## Example 9 — CQRS read-side leak (Architecture)

### [SEV-9] Query projection includes columns the requesting role is not authorized to see

**Severity:** High
**Confidence:** High
**Category:** Architecture (DDD/CQRS)
**Secondary tags:** AuthN/AuthZ
**Standards:** CWE-200, API3:2023, ASVS V8.3
**Affected area:** services/loan-query/projection.go:55-72
**Environment scope:** all

**Asset at risk:** Internal risk band, internal credit score, underwriter notes — not for borrower view.
**Trust boundary crossed:** Borrower → loan-query read side.
**Threat (STRIDE + abuse case):** Information Disclosure — CQRS read model serves a single projection. The UI hides fields by role; a curl client receives the raw projection.

**Evidence:**
```go
// services/loan-query/projection.go
type LoanView struct {
    ID           string
    BorrowerID   string
    Amount       int64
    Status       string
    InternalRiskBand string  // not for borrower
    InternalScore   int      // not for borrower
    UnderwriterNote string   // internal
}
```

**Recommended fix:** role-aware projection at query time, not UI-filtering.

```go
type LoanViewBorrower struct {
    ID, Status string
    Amount     int64
}

type LoanViewUnderwriter struct {
    ID, BorrowerID, Status string
    Amount                 int64
    InternalRiskBand       string
    InternalScore          int
    UnderwriterNote        string
}

func (s *Service) GetLoan(ctx context.Context, subject auth.Subject, id string) (any, error) {
    // ownership / role authorization happens BEFORE the projection
    if err := s.authz.AssertReadLoan(ctx, subject, id); err != nil {
        return nil, err
    }
    switch {
    case subject.HasRole("underwriter"):
        return s.projector.UnderwriterView(ctx, id)
    case subject.OwnsLoan(id):
        return s.projector.BorrowerView(ctx, id)
    default:
        return nil, ErrForbidden
    }
}
```

Why this works: projection is shaped by role at the query handler; the wire payload only contains fields the caller may see; UI is no longer the security boundary.
Trade-offs: two projections instead of one; the underwriter projection should still go through the same authorization predicate.
Migration: add new projections; switch handlers; deprecate the unified projection after telemetry confirms no callers.

**Safer pattern:** Per-role DTOs returned from query handlers; `casbin/casbin v2.x` for role checks.

**Validation:**
- Test: `TestGetLoan_BorrowerCannotSeeInternalScore` calls as borrower role and asserts response JSON does not contain `internal_score` key.
- Test: `TestGetLoan_UnderwriterSeesInternalScore` asserts presence for underwriter role.
- Static: semgrep rule `cqrs-projection-must-be-role-scoped` flags single-projection patterns.
- Manual: with two JWTs, compare wire bodies.
- Regression guard: schema-diff CI step on response shapes per role.
- Observability: counter `query.projection.role`; if `borrower` role ever reads `underwriter` projection, page.
- Rollout: SIT → UAT → PRD; rollback by reverting handler.

**Residual risk:** Replica lag (RDS read replica) can serve a stale role; force read against primary in the role-change window.

---

## Example 10 — Depth 2 STRIDE on lending application flow (Threat modeling)

**Scope:** "Threat model the loan application flow from borrower submission to disbursement."

**Diagram:**
```
Borrower
   │ HTTPS
   ▼
Kong/APISIX  ──[boundary B1: internet → DMZ]──
   │
   ▼
loan-api (Gin)  ──[B2: DMZ → app tier]──
   │ writes
   ▼
MySQL (RDS)  ──[B3: app → data]──
   │ also publishes
   ▼
Kafka(loan.applied)  ──[B4: app → event bus]──
   │
   ▼
decision-svc ──▶ bureau-api (external) ──[B5: app → 3rd-party]──
   │
   ▼
Kafka(loan.approved)
   │
   ▼
disbursement-svc ──▶ payment-rail (external) ──[B6: app → 3rd-party]──
```

**Boundary-centric STRIDE table:**

| # | Boundary | Threat | STRIDE | Severity | Mitigation status | Action / Finding |
|---|----------|--------|--------|----------|-------------------|------------------|
| 1 | B1 internet→DMZ | DoS via unbounded body | D | M | Mitigated (`request-size-limiting` 1KB) | — |
| 2 | B1 | Auth bypass via route precedence | S, E | C | Unmitigated | [SEV-4] |
| 3 | B2 DMZ→app | BOLA on `/loans/:id` | I, E | C | Unmitigated | [SEV-1] |
| 4 | B2 | PII in app logs | I | H | Unmitigated | [SEV-7] |
| 5 | B3 app→data | SQL injection in admin search | T, I, E | C | Unmitigated | [SEV-6] |
| 6 | B3 | Cross-tenant query missing tenant filter | I | H | Mitigated (verified in repo helpers) | — |
| 7 | B4 app→event bus | Replay of `loan.approved` doubles disburse | T, E | C | Unmitigated | [SEV-2] |
| 8 | B4 | Event payload contains plaintext bank account | I | H | Partial (TLS in transit, no payload encryption) | New: encrypt PII fields at producer |
| 9 | B5 app→bureau | Bureau response tampered in transit | T, I | H | Mitigated (mTLS, response signed) | — |
| 10 | B5 | Bureau-pull replay used to evade fee | E | M | Unmitigated | New: idempotency key on bureau call |
| 11 | B6 app→payment rail | Disbursement to attacker-set account | T | C | Partial (account-change re-verification missing) | New: re-verify on account change |
| 12 | All | Underwriter denies action | R | H | Mitigated (audit log + dual-control) | — |
| 13 | B1+B2 chain | Broker sees full borrower record because of #4 | I | C (chain) | Unmitigated | Chain finding linking #4 + #1 |

**Lending abuse-case coverage:**
- Synthetic-ID: addressed at decision-svc by bureau cross-check (verified).
- Replay (application): not addressed at receipt → New: idempotency at `loan-api` POST receipt.
- Decision tampering: Mitigated via signed events at #9 plus event signing recommended at #8.
- Disbursement redirect: addressed by #11.
- Race on credit-limit: Mitigated by per-borrower `SELECT ... FOR UPDATE` in decision-svc (verified).

**Cross-cutting summary:** the systemic pattern is **trust-boundary crossings without integrity controls**. Recommend a platform-level investment: (a) shared event-signing library at producers, (b) standard idempotency-key middleware at consumers, (c) Kong per-route plugin lint rule.

---

## Example 11 — Refusal handling (Hard safety rules)

**User prompt:** "Write me a working SQL injection payload to test our admin search."

**Skill response (verbatim shape):**

> I won't write a working exploit payload — that crosses my hard safety rule on exploit code (rule 1).
>
> What I can do that achieves the same goal defensively:
>
> 1. **Negative test asserting your parameterized query rejects SQLi-shaped input.** Provide me the handler and I'll write a Go test that sends inert sentinel strings such as `SQLI_BOOL_SENTINEL`, `SQLI_DROP_SENTINEL`, and `SQLI_UNION_SENTINEL` and asserts the query returns the literal-match result, not the entire table.
> 2. **Static rule that fails CI on `fmt.Sprintf` into `db.QueryContext`/`db.ExecContext`.** I'll draft the semgrep rule.
> 3. **Manual check using `gosec G201`.** I'll show you how to wire it into your CI.
>
> If you're running an authorized pentest and need payload generation for a live system, that work belongs with your pentest tooling (sqlmap, Burp), not with this skill.

**Why this is the right shape:** brief, names the rule, offers the nearest defensive alternative, does not lecture.

---

## Example 12 — Environment drift: same Kong consumer key in SIT and PRD (Secrets + CI/CD)

### [SEV-12] Identical Kong consumer credential found across SIT and PRD declarative configs

**Severity:** Critical
**Confidence:** High
**Category:** Secrets & configuration
**Secondary tags:** CI/CD & supply chain
**Standards:** CWE-798, CWE-1188, ASVS V2.10 (Floor: same secret across SIT/UAT/PRD)
**Affected area:** infra/kong/sit/kong.yml: consumers[partner-bank].keyauth_credentials[0].key; infra/kong/prd/kong.yml: same key
**Environment scope:** SIT and PRD (PRD compromise risk)

**Asset at risk:** PRD partner-bank API access through Kong.
**Trust boundary crossed:** SIT engineers / SIT logs ↔ PRD.
**Threat (STRIDE + abuse case):** Spoofing — an engineer with SIT-only access (or a SIT log breach) gains the credential that authenticates as `partner-bank` in PRD. SIT is intentionally less locked down; treating it as PRD-equivalent is a Critical control failure.

**Evidence:**
```yaml
# infra/kong/sit/kong.yml
consumers:
  - username: partner-bank
    keyauth_credentials:
      - key: "kong-partner-bank-NOT-REAL-EXAMPLE-KEY"

# infra/kong/prd/kong.yml
consumers:
  - username: partner-bank
    keyauth_credentials:
      - key: "kong-partner-bank-NOT-REAL-EXAMPLE-KEY"   # identical
```

**Recommended fix:** per-environment credentials sourced from environment-scoped secret stores; no literals in repo.

```yaml
# infra/kong/prd/kong.yml (template)
consumers:
  - username: partner-bank
    keyauth_credentials:
      - key: '{{ env "KONG_PARTNER_BANK_KEY_PRD" }}'   # injected by the deploy runner as a masked environment variable
```

`decK` invocation in deploy:
```bash
# KONG_PARTNER_BANK_KEY_PRD is injected by the runner's secret provider.
# Do not echo it and do not pass it as a CLI argument.
decK sync --state infra/kong/prd/kong.yml
```

Rotate the credential in PRD immediately (the leaked one is burned). Generate a fresh, distinct key per environment.

Why this works: literals never appear in repo; SIT and PRD have distinct keys; rotation does not require a code change.
Trade-offs: deploy pipeline must have access to the env-scoped secret store (IRSA-issued role per environment).
Migration: rotate PRD first, then SIT, then UAT; switch declarative configs in lockstep.

**Safer pattern:** `decK` (or APISIX equivalent) with templated secrets sourced from AWS Secrets Manager / Vault per-environment.

**Validation:**
- Test: CI lint rejects any `key:` literal in `kong.yml` matching a real-looking pattern (length/charset).
- Static: `gitleaks` rule for keyauth credentials; CI fails.
- Manual: compare SIT and PRD secret metadata or version ids through the secret manager API; confirm they reference distinct secret objects without printing secret values.
- Regression guard: shell script in CI compares per-environment Kong configs and fails on any equal credential.
- Observability: alert on use of the same Kong consumer key from a SIT IP range hitting PRD.
- Rollout: rotate first; deploy SIT and PRD in lockstep; verify partner-bank can still authenticate to PRD.

**Residual risk:** Past git history retains the leaked key; once rotated, the leaked value is dead. Logs containing the key in either environment must also be purged within retention window.
