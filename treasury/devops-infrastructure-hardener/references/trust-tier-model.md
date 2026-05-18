# Trust-tier model and agents-as-identity-principals

## The trust-tier ladder

Agents are given progressively greater autonomy as they demonstrate
reliability, bounded in all tiers by policy-as-code guardrails.

```
Governance tiers for agent actions
─────────────────────────────────────────────
Tier 1: Suggest    → human decides & executes
Tier 2: Propose    → human approves, agent executes
Tier 3: Act+Notify → agent executes, human informed
Tier 4: Autonomous → agent executes; audit trail only
   (control-plane changes always remain Tier 1–2)
```

The agent's default output is a recommendation, not a mutation. Autonomy is
earned per action type and can be promoted only after the agent demonstrates
reliability at the tier below, instrumented with quality metrics (bug rate,
review time, incident rate per change), never granted up front.

## Data-plane vs control-plane boundary

This boundary is the hard rule of the model.

- **Data-plane actions** — patch a specific vulnerability, rerun a failed
  test, roll back a canary deployment, generate a PR. These act on content,
  not pipeline structure. They MAY rise to Tier 3–4 once reliability is
  demonstrated.
- **Control-plane actions** — changing deployment policies, approval gates,
  rollback thresholds, or production system / configuration. No production
  system delegates control-plane authority without human sign-off.
  Control-plane changes ALWAYS remain Tier 1–2 with prior human approval,
  regardless of the agent's confidence level. Irreversible actions (data
  deletion, production modification, config changes) always require prior
  human approval.

When mapping a finding's remediation to a tier, classify the action first: if
it touches the control plane it is capped at Tier 2 and the hardening plan
must record it as human-gated.

## Agents as identity principals

Every agent should have a documented non-human identity (service account or
workload identity), a defined scope of permissions, and a named owner. This
enables prompt revocation and scoped audit. Without an inventory, incident
response is reactive and slow — 64% of leaked secrets from 2022 are still
active in 2026, primarily a governance gap.

Each agent registers as:

- **Non-human identity** — a unique service account / workload identity, never
  shared with another agent or a human engineer, so a compromised credential
  can be revoked in isolation.
- **Scope / least-privilege** — the minimum permission set for the agent's
  tasks; read-only vs write, specific resource identifiers, time-bounded.
- **Owner** — a named human accountable for the identity's lifecycle.
- **Rotation trigger** — a lifecycle event (deployment change, config change,
  anomaly alert, scope change), never a calendar interval.

The agent-identity inventory (`templates/agent-identity-inventory.md`)
captures one row per agent with these columns.

## Phased autonomy promotion

1. **Phase 0 — Foundation.** Do not promote autonomy without mature delivery
   foundations (test coverage, reliable pipelines, practiced incident
   response). AI amplifies existing capability; weak foundations get worse.
2. **Phase 1 — Scoped + HITL.** Bounded, human-supervised tasks at Tier 1–2
   with approval for every action; instrument quality metrics; define
   escalation triggers before deployment.
3. **Phase 2 — Governance before scale.** Build audit trails, pre-validation,
   spend caps, and an incident playbook before widening scope.
4. **Phase 3 — HOTL + broader scope.** Promote selected data-plane actions to
   Tier 3–4 under human-on-the-loop monitoring; control-plane actions never
   leave Tier 1–2.

Promotion is per action type and reversible: a regression in the tracked
metrics demotes the action.
