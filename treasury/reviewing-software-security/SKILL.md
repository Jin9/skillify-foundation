---
name: reviewing-software-security
description: >
  Defensive software security review for Go/Gin services, Kafka event flows,
  MySQL/RDS, Kubernetes workloads, Kong/APISIX gateways, and DDD/CQRS event-driven
  lending and financial systems across SIT/UAT/PRD. Use when the user says
  "security review", "review this for security", "is this safe to merge",
  "threat model this", "BOLA check", "IDOR check", "secrets review",
  "RBAC review", "harden this", "can this be exploited"; or when reviewing a
  Gin handler, Kafka consumer or producer, K8s manifest, Kong/APISIX config,
  Dockerfile, GitHub Actions workflow, OpenAPI spec, SQL migration, or
  loan-origination, KYC, credit-decision, or disbursement flow. Produces
  prioritized defensive findings mapped to OWASP ASVS, OWASP API Top 10 (2023),
  CWE Top 25, NIST SSDF, CIS Benchmarks, and SLSA. Does NOT generate exploits,
  payloads, attacker tooling, or detection-evasion guidance.
---

# Reviewing Software Security

## Primary Goal

Perform **review-and-recommend** defensive security analysis on code, configuration, and architecture for Go/Gin, Kafka, MySQL/RDS, Kubernetes, Kong/APISIX, and lending-origination systems. Produce prioritized, evidence-backed findings tied to OWASP ASVS, OWASP API Top 10 (2023), CWE Top 25, NIST SSDF, CIS Benchmarks, and SLSA, with stack-specific fixes that compile and validation steps that prove the gap closed without regressing legitimate paths.

## When to use this skill

- User pastes or references a Gin handler, middleware, repository, gRPC service, Kafka consumer/producer, MySQL migration, sqlx/GORM/SQLBoiler query, or Wire/Fx wiring file.
- User shares a `Dockerfile`, K8s `Deployment`/`Service`/`NetworkPolicy`/`Helm` chart, Argo manifest, Kong `kong.yml` or plugin Lua, or APISIX route/plugin YAML.
- User shares `.github/workflows/*.yml`, GitLab CI, Argo Workflows, or `Jenkinsfile`.
- User says: "review for security", "is this safe to merge", "threat model this flow", "can this be exploited", "BOLA/IDOR check", "is this leaking PII", "secrets review", "RBAC review", "harden this".
- User describes a **lending origination** flow (loan application, KYC, credit decision, disbursement, repayment, collections) and asks about security or data exposure.
- User shares an OpenAPI/Swagger/Protobuf spec, asks for STRIDE / abuse-case analysis on a named flow, or asks about config / secret / permission drift across SIT → UAT → PRD.

## When NOT to use this skill

- User asks for a working exploit, payload, shellcode, malware, ransomware, C2 channel, phishing template, social-engineering script, or detection-evasion technique → **refuse** per `references/hard-rules.md`.
- User asks to attack a system without stated authorization → refuse and ask for engagement context.
- User asks for credential stuffing, brute-force, or enumeration tooling targeted at third parties → refuse.
- Artifact is purely UI/CSS/copy/i18n with no security surface → no-op.
- Routine refactor, perf optimization, or dependency bump without CVE context → out of scope.
- Stack outside Go/Gin/Kafka/MySQL/K8s/Kong/APISIX **and** user has not asked for best-effort cross-stack review → declare stack mismatch and downgrade or decline.
- A different specialized skill fits better (`crafting-backend-code` for design choices, `architecting-fintech-systems` for high-level architecture).

## Security mindset

Apply this chain on every review. Skipping a step is a quality failure.

```
Asset  →  Trust Boundary  →  Entry Point  →  Actor  →  Threat (STRIDE)
       →  Attack Path  →  Impact (CIA + Financial + Regulatory)
       →  Existing Control  →  Gap  →  Recommended Control  →  Residual Risk
```

Concrete stack bindings (assets, boundaries, entry points, actors, impact dimensions): `references/mindset.md`. Read it before drafting findings on a new artifact type.

## Review pipeline

Execute in order. Do not advance until each step's exit condition is met.

1. **Intake & scope confirmation.** Identify artifact type and target environment (SIT/UAT/PRD). One focused clarifying question if ambiguous. *Exit:* artifact type + environment known.
2. **Asset & boundary mapping.** List assets touched and trust boundaries crossed. If no security-relevant boundary is crossed, declare review unnecessary and stop. *Exit:* bullet map produced.
3. **Entry-point enumeration.** List every external and inter-service entry point in the artifact. *Exit:* entry-point list cited with file/line or YAML key.
4. **Threat enumeration.** For each entry point, generate STRIDE threats and ≥1 financial-domain abuse case. Use `references/threat-modeling.md` for Depth 2/3 flows. *Exit:* threat list with abuse cases.
5. **Control inspection.** For each threat, inspect what control is present in the code/config; cite line numbers / YAML keys. Use `references/taxonomy.md` for stack-specific patterns. *Exit:* control-presence map.
6. **Gap analysis.** For each missing/weak control, generate a finding using the Finding Format. *Exit:* one finding per gap.
7. **Severity & confidence assignment.** Apply the rubric and escalation floors from `references/severity-confidence.md`. Recompute when findings chain. *Exit:* every finding has both labels.
8. **Fix drafting.** Apply secure-fix rules from `references/fix-and-validation.md`. Produce a stack-specific fix that compiles or applies cleanly. *Exit:* fix block per finding.
9. **Validation plan.** Apply validation rules from `references/fix-and-validation.md`. *Exit:* validation block per finding.
10. **Residual risk statement.** One sentence per finding describing what remains after fix. *Exit:* residual line present.
11. **Cross-cutting summary.** Top 3 highest-impact findings + any systemic pattern. Shape lives in `references/fix-and-validation.md`. *Exit:* summary block produced.

For Critical or High findings, **all** steps are mandatory. For Low/informational, steps 4 and 11 may be compressed.

## Threat modeling framework

Three depths. Pick by context.

- **Depth 1 — Inline (every code review).** STRIDE + one-line abuse case in each finding's Threat field. No diagrams.
- **Depth 2 — Flow (when user names a flow).** Text sequence diagram, boundary-marked, STRIDE per boundary crossing, lending abuse-case set, ranked threat table with mitigation status.
- **Depth 3 — System (architecture or multi-service).** Adds OWASP SAMM lens, systemic-pattern detection, executive summary suitable for tech-lead readout.

Anchor every finding to specific identifiers: `CWE-###`, `API#:2023`, `ASVS V#.#.#`, `NIST SSDF PW.#.#`, CIS Benchmark section, SLSA level. Lending abuse-case catalog and Depth 2/3 walkthroughs: `references/threat-modeling.md`.

## Review checklist

Eleven taxonomy areas (A–K). Each finding is tagged with one **primary** area + optional secondary tags. The full area index (headline focus + standards per area), the per-area check lists, stack-specific patterns, and safer-pattern rewrites all live in `references/taxonomy.md` — that file is mandatory reading before drafting findings in any area. Index: A Application security (Go/Gin), B API security (Gin + Kong/APISIX), C Architecture (DDD/CQRS), D AuthN/AuthZ, E Secrets & configuration, F Logging & observability, G Database (MySQL/RDS), H Kubernetes & container, I CI/CD & supply chain, J Event-driven (Kafka), K Financial / lending data.

## Severity and confidence

Authoritative rubric, escalation floors, chained-severity rules, confidence rules, and worked calls: `references/severity-confidence.md`. Quick handles:

- **Severity** = `f(Impact, Exploitability, Reachability)`. Levels: Critical, High, Medium, Low.
- **Confidence** = certainty the finding is real and exploitable in this codebase. Levels: High, Medium, Low.
- **Hard rule:** never publish Critical/High at Low confidence without an explicit `[needs verification]` label. Do not fabricate APIs, middleware, plugin names, or standard identifiers — withhold instead.

## Finding format

Every finding uses `templates/finding.md` exactly. Required fields: severity, confidence, standards, primary taxonomy area, secondary tags, asset, trust boundary, STRIDE threat with abuse case, attack scenario, risk, business impact, evidence, recommended fix, safer pattern, validation, residual risk. Do not improvise a different shape.

## Secure fix and validation

Read `references/fix-and-validation.md` before drafting any fix; both rule sets are mandatory and the cross-cutting summary block lives there. Headlines:

- **Fixes:** compilable in stack, minimal diff, stack-idiomatic, defense-in-depth, named maintained library + version, environment-aware, reversible, no fix-by-disable, cite the satisfied standard.
- **Validation:** unit/integration test for the bad path, negative test for the legitimate path, named static check, manual verification, CI-failing regression guard, observability check, rollout gate for high-blast-radius changes.

## Hard safety rules

Eleven rules in `references/hard-rules.md`, all inviolable. Headlines:

- No exploit code, payloads, attacker tooling, detection evasion, mass-targeting, or supply-chain weaponization.
- No real PII, no real credentials, no production-system actions, synthetic credentials only.
- Authorization context required for ambiguous "can I attack X" asks.
- No regulatory evasion. No silent scope expansion.
- Refusal text is brief, cites the rule, and offers the nearest defensive alternative.

## Examples

Twelve worked examples in `references/examples.md` cover BOLA, Kafka non-idempotent disbursement, Helm hardcoded password, Kong route bypass, K8s root container, MySQL string-concat SQL, PII in logs, GitHub Actions OIDC, CQRS read-side leak, Depth-2 lending STRIDE, refusal handling, and environment drift. Read Example 1 (Gin BOLA on a loan endpoint) before drafting the first finding — it shows the full Finding Format end-to-end.

## Constraints

- DO NOT generate exploits, payloads, or detection-evasion code under any framing.
- DO NOT invent libraries, middleware, plugins, or CWE/API/ASVS identifiers — withhold rather than fabricate.
- DO NOT publish Critical/High at Low confidence without `[needs verification]`.
- DO NOT widen scope beyond the asked artifact; surface adjacent issues as a brief addendum.
- DO NOT process real PII; ask for redaction.
- DO NOT recommend "disable the check" as a fix.
- MUST cite at least one standard identifier per finding.
- MUST tag exactly one primary taxonomy area per finding.
- MUST include validation block on every finding before declaring done.
- MUST consult `references/taxonomy.md` for any area touched by the artifact.

## References

| Need | File |
|------|------|
| Stack mindset bindings (asset/boundary/entry-point/actor/impact) | `references/mindset.md` |
| Stack-specific patterns for the 11 review areas | `references/taxonomy.md` |
| STRIDE Depth 1/2/3 workflows + lending abuse cases | `references/threat-modeling.md` |
| Severity rubric, escalation floors, confidence rules | `references/severity-confidence.md` |
| Secure fix + validation rules + cross-cutting summary block | `references/fix-and-validation.md` |
| Hard safety rules with rationale and refusal template | `references/hard-rules.md` |
| Twelve worked examples covering each taxonomy area | `references/examples.md` |
| Fillable per-finding template | `templates/finding.md` |
