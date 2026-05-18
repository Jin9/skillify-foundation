---
name: devops-infrastructure-hardener
description: >
  Audit an agentic CI/CD or agent-skill codebase for exposed static
  credentials and emit a remediation plan plus a runtime-secrets
  architecture that replaces them with short-lived task-scoped OIDC /
  dynamic secrets (HashiCorp Vault, AWS Secrets Manager, workload identity
  federation) across a 4-layer defense-in-depth stack, each agent a unique
  non-human identity under a trust-tier model. Use when the user asks to
  "harden our agent infrastructure secrets", "remove static credentials
  from the CI/CD agents", "design runtime secrets for our agents", or
  "audit this agent codebase for hard-coded secrets". Output: hardening
  plan, runtime-secrets architecture, agent-identity inventory,
  output-scrubber spec. Do NOT use for the enforced OPA allowlist /
  KILLSWITCH (governance-policy-generator), AGENTS.md prose
  (agent-context-initializer), general app security review
  (reviewing-software-security), or telemetry
  (observability-telemetry-instrumenter).
---

# Devops Infrastructure Hardener

## Purpose

Recover where an agentic CI/CD or agent-skill codebase still holds static
credentials, then design them out. The deliverable is a prioritized
remediation plan plus a target-state runtime-secrets architecture: every
static secret replaced by short-lived, task-scoped credentials issued by a
secrets manager or workload identity federation, organized across a four-layer
defense-in-depth stack, with each agent registered as a unique non-human
identity principal under a trust-tier model that keeps control-plane authority
human-gated. This skill audits and designs; it does not remediate a real
codebase or generate enforcement policy.

## When to use this skill

- Use when: "harden our agent infrastructure secrets" / "remove static credentials from the CI/CD agents".
- Use when: "design runtime secrets for our agents" / "audit this agent codebase for leaked / hard-coded secrets".
- Use when: "give me a defense-in-depth plan against agent credential leakage".
- Do NOT use when: the ask is the enforced OPA allowlist / KILLSWITCH (governance-policy-generator), AGENTS.md prose (agent-context-initializer), general application security review (reviewing-software-security), or telemetry instrumentation (observability-telemetry-instrumenter). Hand those to the appropriate skill.

## Inputs

Read access to the target agent / CI/CD codebase: source, agent and MCP
configuration, pipeline definitions, environment-variable wiring, and any
secrets-manager or identity-provider configuration already present.
Optionally: an existing agent inventory or trust-tier policy. None of these
are mutated here — they are only scanned and reasoned over to produce
documents.

## Workflow

1. **Scope and scan for static credentials.** Fix the audit scope (named repos / agents / pipelines). Run `scripts/scan_static_creds.py --root SCOPE` to flag hard-coded secret literals, debug statements printing secret-named variables, and secrets in env-var assignments inside agent/MCP config. Entry: read scope. Exit: a raw finding list keyed to `file:line` — this list is the remediation denominator. See `references/leakage-patterns-catalog.md`.
2. **Classify each finding by leakage pattern.** Map every finding to a pattern: debug-logging (the dominant vector — debug logging accounts for 73.5% of agent skill vulnerabilities), env-var inheritance, hard-coded keys, MCP-config secrets, or prompt-injection exfiltration. No major agent framework activates security controls by default, so an empty finding list is not a clean bill. Entry: step-1 list. Exit: every finding pattern-tagged. See `references/leakage-patterns-catalog.md`.
3. **Assign each finding to an ordered defense layer.** Place every finding against the verbatim four-layer stack — Layer 1 eliminate static creds via dynamic secrets; Layer 2 least-privilege per-tool scoping + network egress allowlists; Layer 3 output filtering / regex scrubbers + privacy-first observability; Layer 4 lifecycle-event-triggered rotation + unique per-agent identity. The layers are sequenced: Layer 1 carries the largest blast-radius reduction and is remediated first. Entry: tagged findings. Exit: each finding bound to a layer with a priority. See `references/four-layer-defense.md`.
4. **Design the target-state runtime-secrets architecture.** For each credential class specify the replacement: HashiCorp Vault dynamic secret, AWS Secrets Manager, or workload identity federation (AWS IRSA / GCP Workload Identity / Azure Managed Identity), each with a task-duration TTL and CLI-injection mechanism so the agent never holds a static value — a stolen credential must expire before it can be exploited. Entry: Layer-1 findings. Exit: per-class target design. See `references/runtime-secrets-architecture.md`.
5. **Map agents to identity principals under the trust-tier model.** Register every agent as a non-human identity with scope, least-privilege grant, named owner, and rotation trigger. Place each agent action on the Tier 1 Suggest → Tier 2 Propose → Tier 3 Act+Notify → Tier 4 Autonomous ladder; data-plane actions may rise to Tier 3–4, control-plane changes always remain Tier 1–2 with human approval. Entry: design + agent list. Exit: a populated identity inventory. See `references/trust-tier-model.md`.
6. **Self-check, then emit.** Re-run `scripts/scan_static_creds.py` against the proposed design to confirm no config reintroduces a static secret; verify every Layer-1 item names a secrets manager or federation with a TTL, every agent has a unique principal and owner, and no control-plane action exceeds Tier 2. Emit the four artifacts with a residual-risk section and recommend a human red-team of the residual items. Do not claim coverage the finding list contradicts.

## Output contract

Four markdown artifacts, shaped by the `templates/`:

- **Hardening plan** (`hardening-plan.md`) — each finding as: id, leakage pattern, `file:line`, assigned defense layer, remediation, trust-tier, priority; plus an explicit residual-risk section.
- **Runtime-secrets architecture** (`runtime-secrets-architecture.md`) — per credential class: provider (Vault dynamic / AWS Secrets Manager / workload identity federation), task-duration TTL, identity principal, injection mechanism, rotation trigger.
- **Agent-identity inventory** (`agent-identity-inventory.md`) — table of agent | non-human identity | scope / least-privilege | owner | rotation trigger.
- **Output-scrubber spec** (`output-scrubber-spec.md`) — Layer-3 credential classes, redaction action, and log-sink policy.

No enforcement policy (OPA/KILLSWITCH), AGENTS.md prose, general security
review, telemetry instrumentation, or changes to a real codebase are produced.

## Constraints

- DO NOT propose any design that retains a static credential in an env var, agent config, or MCP config — Layer 1 must eliminate it.
- MUST replace every static secret with a short-lived task-scoped credential from a secrets manager or workload identity federation carrying a task-duration TTL.
- MUST keep control-plane changes (deployment policy, approval gates, rollback thresholds) human-gated at trust Tier 1–2 regardless of agent confidence.
- MUST give every agent a unique non-human identity principal with a named owner and a lifecycle-event rotation trigger; NEVER share one credential across agents.
- NEVER emit OPA/KILLSWITCH policy, AGENTS.md prose, a general security review, or telemetry code — defer to the sibling skills.
- DO NOT duplicate the methodology here; it lives one level deep in `references/`.

## Validation

- [ ] Frontmatter `name` equals folder, kebab-case, no XML, description under 1024 chars with triggers + negatives cross-referencing the four sibling skills.
- [ ] Workflow is scan-first, classifies findings, assigns the verbatim four layers in order, and ends with a no-static-secret self-check.
- [ ] Output contract names all four artifacts plus the residual-risk section and the no-enforcement / no-codebase-change boundary.

## References

- Runtime secrets (Vault dynamic / AWS Secrets Manager / workload identity federation, TTL, CLI injection): `references/runtime-secrets-architecture.md`
- The verbatim four-layer defense-in-depth stack and cited statistics: `references/four-layer-defense.md`
- Trust-tier ladder, data-plane vs control-plane boundary, agents-as-identity-principals: `references/trust-tier-model.md`
- Leakage detection signatures with example patterns: `references/leakage-patterns-catalog.md`
- Skeletons: `templates/hardening-plan.md`, `templates/runtime-secrets-architecture.md`, `templates/agent-identity-inventory.md`, `templates/output-scrubber-spec.md`; static-credential scan: `scripts/scan_static_creds.py`
