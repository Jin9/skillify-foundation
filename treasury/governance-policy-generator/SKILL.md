---
name: governance-policy-generator
description: >
  Emit policy-as-code agent-governance artifacts architecturally separated
  from the agent they govern: an Open Policy Agent (OPA) Rego allowlist with
  default-deny posture (default allow = false), plus a KILLSWITCH.md defining
  TRIGGERS, FORBIDDEN actions, a three-level throttle-pause-stop ESCALATION,
  an append-only JSONL audit, and OVERRIDE conditions. Every action is
  classified by reversibility times blast radius before any allow rule;
  package installation requires registry-allowlist verification. Use when the
  user asks to "write an auto-approve policy", "generate an OPA / Rego
  allowlist for an agent", "create a KILLSWITCH.md", "design a default-deny
  agent governance policy", or "kill switch and audit for an autonomous
  squad". Do NOT use for credential or identity hardening (defer to
  devops-infrastructure-hardener), AGENTS.md prose (defer to
  agent-context-initializer), tool or skill schema CI validation (that is
  universal-spec-validator), or executing or sandboxing the governed commands.
---

# Governance Policy Generator

## Purpose

Produce the policy-as-code and emergency-stop artifacts that govern an
autonomous agent squad, structured so the agent cannot read or rewrite its
own constraints. The deliverables are a default-deny OPA/Rego allowlist
(`default allow = false`, the primary access model), a supplementary
blacklist (an emergency brake, never the primary model), and a KILLSWITCH.md
that lives in infrastructure the agent cannot touch. Every action category is
classified by reversibility times blast radius before a single allow rule is
written. This skill writes governance artifacts only — it never executes,
sandboxes, or merges the actions it governs.

## When to use this skill

- Use when: "write an auto-approve policy for our agent squad" / "generate an OPA / Rego allowlist for a coding agent".
- Use when: "create a KILLSWITCH.md with triggers, forbidden actions, and escalation" / "design a default-deny agent governance policy".
- Use when: "classify our agent actions by reversibility and blast radius before we widen auto-approve" / "policy-as-code for an autonomous squad".
- Use when: "add a registry-allowlist gate so the agent cannot install hallucinated packages".
- Do NOT use when: the ask is credential/identity hardening, AGENTS.md prose, tool/skill schema CI validation, or running/sandboxing the commands. Hand those to the appropriate skill (`devops-infrastructure-hardener`, `agent-context-initializer`, `universal-spec-validator`).

## Inputs

Consumed, never generated here: the agent's tool/command surface (the set of
operations it can attempt); the target environment and what is reversible by
a Git revert there; any existing deny list, sandbox config, or branch
protections; the squad's risk tolerance and the human escalation contacts;
the known-good package registry to allowlist against. If the reversibility of
an action category is unknown, treat it as irreversible until proven otherwise.

## Workflow

1. **Inventory the action surface.** Enumerate every command, tool call, MCP
   operation, and higher-level action (git push, deploy, package install,
   external API call) the agent can attempt. Entry: input list. Exit: a flat
   action inventory with nothing collapsed or assumed. No allow rule may
   reference an action absent from this inventory. See
   `references/action-classification.md`.

2. **Classify by reversibility times blast radius FIRST.** Before writing any
   allow rule, place every inventoried action on the two-axis grid
   (reversibility: undoable by a Git revert or equivalent; blast radius:
   resources/users affected if wrong). Entry: inventory. Exit: a completed
   classification table where every action has an auto-approve eligibility
   verdict. High-reversibility, low-blast-radius only is the starting
   envelope. See `references/action-classification.md`.

3. **Define the three-layer enforcement context.** Record the OS/sandbox
   floor (strongest, the non-negotiable base), the default-deny allowlist
   (primary), and the supplementary blacklist (emergency brake). Note that
   instruction-level rules are the weakest layer and never the boundary.
   Entry: classification. Exit: a stated layer map the Rego sits within. See
   `references/policy-as-code-model.md`.

4. **Generate the default-deny Rego allowlist.** From `templates/policy.rego`:
   set `default allow = false`; write one conditioned allow rule per
   auto-approve-eligible action only; every rule body MUST have at least one
   condition (no unconditioned allows); package installation MUST route
   through a registry-allowlist rule. Mark every rule needing human review.
   Entry: classification + layer map. Exit: a `.rego` that passes
   `scripts/rego_deny_lint.py`. See `references/policy-as-code-model.md`.

5. **Generate the KILLSWITCH.md in separated infrastructure.** From
   `templates/KILLSWITCH.md`: fill TRIGGERS, FORBIDDEN, the three-level
   throttle-pause-stop ESCALATION, the append-only JSONL audit schema, and
   OVERRIDE. Banner that it MUST live where the agent cannot read or modify
   it. Entry: classification. Exit: a complete KILLSWITCH.md. See
   `references/killswitch-standard.md`.

6. **Self-check against the known-bypass catalog, then emit.** Confirm the
   blacklist is supplementary not primary, that no syntactic deny rule is
   relied on as a boundary, and that the six documented bypass classes are
   acknowledged. Run `scripts/rego_deny_lint.py` on the generated `.rego`.
   Entry: drafts. Exit: artifacts emitted only after the lint passes. See
   `references/known-bypass-catalog.md`.

## Output contract

Three artifacts, written to disk, never executed:

- **`policy.rego`** — an OPA/Rego policy opening with `default allow = false`;
  one conditioned allow rule per auto-approve-eligible action; a
  package-installation rule gated on a registry allowlist; explicit
  `# HUMAN REVIEW` markers on every rule that widens autonomy.
- **`KILLSWITCH.md`** — TRIGGERS, FORBIDDEN actions, three-level
  throttle-pause-full-stop ESCALATION, an append-only JSONL audit schema, and
  OVERRIDE conditions, with a banner that it must live in infrastructure the
  agent cannot modify.
- **`action-classification.md`** — the reversibility-times-blast-radius table
  with an auto-approve eligibility column, the provenance for the envelope.

No execution, sandboxing, merging, credential handling, or AGENTS.md prose is
performed. The Rego is a skeleton for human review, not a deployed policy.

## Constraints

- MUST set `default allow = false`; a `default allow = true` or any
  unconditioned allow rule is a hard failure.
- MUST classify every action by reversibility times blast radius BEFORE
  writing any allow rule; never widen the envelope past
  high-reversibility/low-blast-radius without a documented rollback story.
- MUST route package installation through a registry-allowlist rule (roughly
  20% of agent-recommended package names are hallucinated).
- NEVER treat the blacklist as the primary access model — it is a
  supplementary emergency brake; the allowlist is primary, the OS sandbox is
  the floor.
- NEVER place the kill switch, audit, or policy in the agent's own context,
  prompt, or any file the agent can write.
- DO NOT execute, sandbox, or merge governed actions; DO NOT do credential
  hardening, AGENTS.md prose, or schema CI validation.
- DO NOT duplicate the methodology here; it lives one level deep in
  `references/`.

## Validation

- [ ] Frontmatter `name` equals folder, kebab-case, no angle brackets, description 500-1024 chars with trigger language and sibling-skill negatives.
- [ ] Workflow classifies by reversibility times blast radius before any allow rule, mandates `default allow = false`, and ends with a known-bypass self-check plus `scripts/rego_deny_lint.py`.
- [ ] Output contract names `policy.rego`, `KILLSWITCH.md`, `action-classification.md`, the architectural-separation rule, and the no-execution boundary.

## References

- Three-layer enforcement, OPA/Rego rationale, default-deny primacy, allowlist-vs-blacklist trade-offs: `references/policy-as-code-model.md`
- Reversibility-times-blast-radius taxonomy, rollback-before-irreversible rule, promotion criteria: `references/action-classification.md`
- KILLSWITCH.md field semantics, escalation thresholds, JSONL audit schema, OVERRIDE, separation rule: `references/killswitch-standard.md`
- The six bypass techniques, CVEs, and the 36/81/11/20 figures, why blacklist-only is false confidence: `references/known-bypass-catalog.md`
- Skeletons: `templates/policy.rego`, `templates/KILLSWITCH.md`, `templates/action-classification.md`; default-deny linter: `scripts/rego_deny_lint.py`
