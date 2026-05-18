# The deferral meta-principle

## The rule

The ASDLC (Agentic Software Development Life Cycle) framework states the
meta-principle precisely:

> if a constraint can be expressed elsewhere — in a lint rule, a CI check,
> a branch protection rule — it must not be restated in AGENTS.md.

The tool is the constraint; AGENTS.md is for what cannot be enforced
programmatically. The synthesis across every documented trade-off converges
on a principle of **minimum viable context**: write only what agents cannot
discover from code or infer from conventions. The file should be small
*because* most constraints belong in lint configs, CI pipelines, and branch
protection rules, not in natural language.

## The deferral filter (applied in Workflow step 2)

For every candidate rule, ask the three questions in order. A "yes" to any
means the rule is DROPPED from AGENTS.md and logged to the curation
checklist's deferral table (`templates/curation-checklist.md`):

1. **Lint rule?** — formatting, naming, import order, unused symbols,
   complexity ceilings, banned APIs. Belongs in the lint/format config.
   Example dropped rule: "use 4-space indents" → `pyproject.toml`.
2. **CI check?** — "tests must pass", "coverage >= X", "type-check clean",
   "no `TODO` left". Belongs in the CI pipeline as a required job.
   Example dropped rule: "never merge with failing tests" → CI required
   check.
3. **Branch-protection rule?** — "no direct pushes to main", "require N
   approvals", "require status check Y", "no force-push". Belongs in the
   repository's branch-protection settings.
   Example dropped rule: "require one review before merge" → branch
   protection.

What survives the filter is exactly the non-enforceable knowledge: project
structure orientation, the one style snippet, the non-obvious conventions,
and the human-approval boundaries an agent cannot infer.

Machine-enforced *command* allow/deny and auto-approve policy is a separate
concern owned by `governance-policy-generator`; runtime secret and identity
hardening is owned by `devops-infrastructure-hardener`. Defer those too —
do not encode them as prose here.

## Over-specification vs under-specification

Both directions fail:

- **Over-specification.** Files with 15+ sequential "don'ts" and no
  corresponding "dos" cause agents to over-explore, stay conservative, and
  produce incomplete work. Files exceeding roughly 60-200 lines see
  degraded instruction-following as important rules get buried in noise.
- **Under-specification.** Vague statements about code quality ("write
  clean code", "handle errors gracefully") produce agents stuck in
  "exploration mode" that do not complete tasks reliably.

Resolution: pair every prohibition with a concrete alternative, prefer
numbered multi-step workflows over open-ended guidance, and push everything
machine-enforceable out via the deferral filter so the surviving prose
stays short and high-signal.

## Static vs dynamic, and stale-context poisoning

AGENTS.md is static and version-controlled, which is its strength — but
stale content is especially dangerous for agents because it **actively
poisons every request's context**, unlike a stale README that humans
mentally discount. Every token loads on every generation step regardless of
relevance. Treat AGENTS.md as institutional memory: update it when an agent
makes a mistake the file could have prevented, and trim sections the moment
they no longer apply. Distinguish stable context (the AGENTS.md body, which
can be prompt-cached) from volatile task-specific instructions, which do
not belong in the file at all.
