# Threat Modeling: Depths, Workflows, and Lending Abuse Cases

Pick the depth from context. Apply the workflow exactly. Cite identifiers per finding.

---

## Depth 1 — Inline (every code review)

**When:** every finding produced by the Review Pipeline.

**Workflow:**
1. For the entry point in scope, generate STRIDE threats:
   - **S**poofing: who could impersonate whom? (forged JWT, replayed cookie, spoofed Kafka header)
   - **T**ampering: what can be modified in flight or at rest? (request body, event payload, DB row)
   - **R**epudiation: can an actor deny an action? (missing audit log on disbursement)
   - **I**nformation disclosure: what data is exposed? (PII in response, error message, log)
   - **D**enial of service: what causes the system to stop serving? (unbounded payload, no rate-limit, no resource limits)
   - **E**levation of privilege: how does a low-privilege actor become high-privilege? (BOLA, BFLA, role injection)
2. Pick at least one **financial-domain abuse case** from the catalog below. Write it as a one-liner ("Borrower replays `loan.approved` to double-disburse").
3. Embed in the finding's **Threat (STRIDE + abuse case)** field.

**No diagrams at Depth 1.** Keep it to two lines.

---

## Depth 2 — Flow (when user names a flow)

**When:** user asks "threat model this flow", "STRIDE on the disbursement path", or shares a sequence of services.

**Workflow:**

### Step 1 — Render the flow as a text diagram

```
Borrower
   │ HTTPS
   ▼
Kong/APISIX  ──[trust boundary: internet → DMZ]──
   │ HTTP (mTLS preferred)
   ▼
loan-api (Gin)  ──[trust boundary: DMZ → app tier]──
   │
   ├──▶ MySQL (RDS) ──[trust boundary: app → data tier]──
   │
   └──▶ Kafka(loan.applied) ──[trust boundary: app → event bus]──
                │
                ▼
         decision-svc  ──▶  bureau-api (external)
                │
                ▼
         Kafka(loan.approved)
                │
                ▼
         disbursement-svc  ──▶  payment-rail (external)
```

### Step 2 — Mark trust boundaries explicitly

A trust boundary is any place where the *caller's authority changes*: internet ↔ DMZ, DMZ ↔ app tier, app ↔ data tier, app ↔ event bus, internal ↔ external.

### Step 3 — STRIDE per boundary crossing (not per component)

Boundary-centric is faster and catches more — most exploitable bugs live at boundaries, not inside components.

For each boundary:
- Who crosses it? (actor)
- With what data? (payload)
- What authentication / authorization is asserted at the crossing?
- What integrity / confidentiality controls protect the data in transit?
- What logging / audit captures the crossing?

### Step 4 — Apply lending abuse cases (catalog below)

Pick the cases relevant to the flow. For a loan-application flow: synthetic-ID, replay, decision-tampering, fee-evasion. For a disbursement flow: redirection, double-disburse, race-on-credit-limit. For a KYC flow: substitution, document-tampering.

### Step 5 — Produce a ranked threat table

```markdown
| # | Threat | STRIDE | Affected boundary | Severity | Mitigation status | Reference |
|---|--------|--------|-------------------|----------|-------------------|-----------|
| 1 | Borrower iterates loan ids → reads other borrowers | I, E | Kong → loan-api | Critical | Unmitigated | [SEV-1] |
| 2 | Replay of `loan.approved` → double-disburse | T, E | Kafka → disbursement-svc | High | Partial (DLQ only, no idempotency key) | [SEV-3] |
| 3 | KYC document URL leaked in logs | I | loan-api → log sink | High | Unmitigated | [SEV-7] |
| 4 | DoS via unbounded loan-application body | D | Kong → loan-api | Medium | Mitigated (`request-size-limiting` plugin) | — |
| 5 | Underwriter denies action ("I didn't approve that") | R | loan-api → MySQL audit | High | Mitigated (audit log + dual-control) | — |
```

**Mitigation status values:**
- `Mitigated` — control is present and verified.
- `Partial` — control exists but is incomplete (DLQ without idempotency, rate-limit without per-tenant quota).
- `Unmitigated` — no control. Generates a finding.
- `Accepted` — risk explicitly accepted by the org with rationale and expiry.

### Step 6 — Output

Embed the table at the top of the review output, before individual findings. Each finding cross-references the table row.

---

## Depth 3 — System (architecture or multi-service review)

**When:** user shares architecture diagram, multiple services at once, or asks for a system-level security read.

**Adds to Depth 2:**

### S.1 — OWASP SAMM lens

For each unmitigated threat, classify where in the SDLC it should be caught:

- **Design** — should have been caught at architecture review (e.g., "events with PII published before crypto-shredding plan").
- **Implementation** — should have been caught in code review (e.g., "BOLA pattern not flagged").
- **Verification** — should have been caught in testing (e.g., "no negative test for cross-tenant access").
- **Operations** — should have been caught at runtime (e.g., "default-deny NetworkPolicy missing").

This signals to the architect *which control plane* needs strengthening, not just which finding to fix.

### S.2 — Systemic-pattern detection

Look for a single root cause behind multiple findings:
- "Auth middleware missing on 4 of 7 routes" → router-group misconfiguration, not a per-route fix.
- "PII leaks in 3 services" → no shared structured-logging library with redaction, not a per-service fix.
- "Three services each have their own retry-loop" → no shared idempotency convention.

Surface the root cause prominently; recommend a platform-level fix.

### S.3 — Executive summary (for tech-lead readout)

```markdown
## Executive summary

**Posture:** [one paragraph: where the system stands]
**Top 3 risks:** [bullet list with severity]
**Systemic patterns:** [bullet list]
**Recommended platform investments:** [e.g., "shared auth middleware lib", "central log-redaction handler", "OPA admission policy bundle"]
**Quick wins (≤ 1 week):** [bullet list]
**Medium-term (≤ 1 quarter):** [bullet list]
```

---

## Lending-domain abuse-case catalog

Use these as starting points. Combine with STRIDE to produce concrete threats.

### Origination
- **Synthetic identity application.** Attacker creates a borrower from fragments of real PII (NIK from one source, address from another) to obtain credit they will not repay. Detection: bureau cross-check, document liveness, behavioral signals.
- **Replay of a successful application.** Attacker re-submits an approved application to obtain a second loan or to evade a hold. Defense: idempotency at application receipt, application-id replay window.
- **Application tampering in transit.** Attacker modifies amount or term mid-flight. Defense: signed application payload (HMAC), TLS, integrity check at decisioning.
- **Document substitution.** Attacker submits another person's documents with their own face/signature swapped. Defense: liveness detection, document authenticity checks, separate verification of identity vs. document.

### Decisioning
- **Decision tampering.** Attacker modifies a `loan.approved` event before the disbursement service consumes it. Defense: signed events, consumer-side verification.
- **Bureau-pull manipulation.** Attacker influences which bureau record is fetched (e.g., by submitting a similar NIK). Defense: exact-match policy, name+DOB cross-check.
- **Race on credit limit.** Attacker submits N parallel applications under one borrower to exceed the per-borrower limit before any single application commits. Defense: per-borrower lock or credit-reservation pattern at application receipt.
- **Underwriter coercion / collusion.** Insider with override authority approves bad applications. Defense: dual-control, override audit, anomaly detection on per-underwriter approval rates.

### Disbursement
- **Account-number redirection.** Attacker modifies the destination bank account between approval and disbursement. Defense: account-change re-verification, payout signed end-to-end.
- **Double-disburse via replay.** At-least-once delivery causes the disbursement consumer to fire twice. Defense: consumer idempotency keyed by `loan_id` or event id.
- **Cross-currency abuse.** Attacker exploits FX rounding between the loan currency and the disbursement currency. Defense: `decimal` arithmetic, explicit rounding rule.
- **Disbursement to sanctioned party.** Defense: sanctions-list screening at disbursement, not only at onboarding.

### Servicing & repayment
- **Repayment misallocation.** Repayment applied to the wrong loan or the wrong principal/interest split. Defense: deterministic allocation rule, audit log of every allocation.
- **Refund to attacker-controlled account.** A borrower-initiated refund routes funds to an account the attacker controls. Defense: refund-account verification, dual-control on refunds above threshold.
- **Charge-off manipulation.** Insider charges off an active loan to write off a confederate's debt. Defense: dual-control, anomaly detection.

### Cross-cutting
- **PII scrape via BOLA / BFLA.** Defense: object-level authz on every `:id` route, BFLA on admin verbs.
- **Log-aggregator breach exposes PII.** Defense: redaction at source, log sink access controls, retention bounded.
- **Backup compromise.** Defense: encrypt backups with same KMS posture as primary, cross-account snapshot reviewed.

---

## Cross-references

Every threat the skill flags must cite **at least one** of:
- `CWE-###` (Common Weakness Enumeration)
- `API#:2023` (OWASP API Top 10 2023)
- `ASVS V#.#.#` (OWASP ASVS V1–V14)
- `NIST SSDF PW.#.#` or `PS.#.#`
- `CIS K8s/MySQL/Docker §#.#`
- `SLSA L#`

If no identifier fits cleanly, the finding belongs in a domain-specific category (lending abuse case) and should cite the lending regulatory anchor (PDPA, GDPR, BOT, OJK, MAS, PCI-DSS) instead.
