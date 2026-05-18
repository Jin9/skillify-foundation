# Topology decision record — &lt;task / system name&gt;

> Skeleton. Replace every bracketed slot. Delete guidance lines before delivery.

## Task summary

[One paragraph: what the task is, the agents considered, the shared state.]

## Four-step gate checklist

- [ ] **Step 1 — Deterministic 90%+ flowchart?** [yes -> recommend a
      workflow, not agents, and STOP | no -> continue]. Evidence: [why]
- [ ] **Step 2 — Single agent first.** [Is it tightly coupled, mostly
      sequential, global-context, under 20 tools? If yes, recommend ONE
      agent and STOP.] Evidence: [why]
- [ ] **Step 3 — Multi-agent justified?** [Only if parallelizable,
      read-heavy, with independent sub-problems and sub-agent isolation
      prevents context pollution.] Evidence: [why]
- [ ] **Step 4 — Hard quality gate present?** [A schema + a verifier
      that REJECTS rather than blends outputs. If absent, do not build
      an independent multi-agent system.] Evidence: [verifier identity]

## Decision

- Recommendation: [workflow | single agent | multi-agent squad]
- Chosen topology (if squad): [orchestrator-worker (default) | other]
- Single-vs-multi rationale: [tie the decision to the step that gated it]

## Justified boundary

[State the specific boundary that justifies any handoff: genuine
parallelism, context-window limit, or separated security domains. If no
such boundary exists, the recommendation must be single agent.]

## Cost acknowledgement

[Token/latency cost accepted: e.g. ~15x tokens vs single chat for
orchestrator-worker; justified because task value is high enough.]
