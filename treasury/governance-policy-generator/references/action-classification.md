# Action classification: reversibility times blast radius

Deep-dive for Workflow steps 1-2. This is done BEFORE any allow rule.

## The mandatory ordering rule

"Before writing any policy rule, classify every action category the agent can
take by two dimensions: reversibility (can this be undone by a Git revert or
equivalent?) and blast radius (how many resources/users are affected if this
action is wrong?). Auto-approve scope should begin with high-reversibility,
low-blast-radius actions only."

This is step zero, not a documentation afterthought: "Start with a taxonomy
of action reversibility and blast radius" is principle 1 of the design
playbook. A rule written before classification is, by definition, ungoverned.

## Axis 1 — reversibility

"Not all agent actions are reversible. A sent email, an external API call, a
deployed cloud resource — these cannot be undone by reverting a Git commit."
Bucket each action:

- **Reversible** — undoable by a Git revert or equivalent (git commit, branch
  creation, local file edit in a tracked tree).
- **Irreversible** — cannot be undone by reverting a commit (deployment,
  external API call, package installation, sent email, DB migration).

If reversibility is unknown, classify as **irreversible until proven
otherwise** — the conservative default the SKILL.md inputs section mandates.

## Axis 2 — blast radius

"Agentic blast radius — the maximum harm an autonomous AI agent can cause if
its credentials or delegated permissions are compromised — grows with the
scope of auto-approve policy." A compromised agent with broad permissions
"can affect many more resources than intended even when its actions appear
locally permitted." Sophos's "lethal trifecta" is identity compromise,
over-privileged tools, and absent blast-radius controls.

- **Low** — one file, one branch, one ephemeral resource.
- **Medium** — a service, a shared dependency, a protected branch.
- **High** — production, many users, organization-wide config, credentials.

## The grid and auto-approve eligibility

| | Low blast | Medium blast | High blast |
|---|---|---|---|
| **Reversible** | auto-approve eligible (start here) | human review | human review |
| **Irreversible** | rollback-scoped review | human review | human review only |

Only the high-reversibility / low-blast-radius cell is the starting
auto-approve envelope. Everything else stays under human review until a
specific, tested, audited justification exists.

## Rollback scope before any irreversible auto-approve

"The rollback scope must be defined and documented *before* an incident
occurs: which effects can be undone, what the target state is for reversible
effects, and what procedure triggers rollback. Auto-approve policies should
therefore be scoped first to reversible actions (Git commits, branch
creation) and extended to irreversible actions (deployment, external API
calls, package installation) only after the rollback story is complete."
Prefer to "keep irreversible actions under human review indefinitely unless
there is a specific, tested, and audited justification."

## Separate observe / recommend / execute

"Best practice separates the permission tiers: agents can always observe
(read code, logs, metrics); agents can recommend (open PRs, file issues)
without triggering immediate action; agents execute only within a defined
auto-approve envelope." This three-tier separation "prevents agents from
acting on their own analysis without a gate."

## Promotion criteria for widening the envelope

Auto-approve "is not a binary switch but a *trust calibration instrument* —
one that should expand gradually as telemetry confirms safe behavior."
Empirically, newer users (under 50 sessions) use full auto-approve roughly
20% of the time, rising to over 40% by 750 sessions; experienced users also
interrupt *more* often (9% of turns vs. 5% for new users) — a shift from
per-action approval to trajectory monitoring.

Phased rollout: "Begin with a sandbox environment and read-only permissions,
log every action, measure human override frequency and reversal rates, and
expand execution privileges only after reversal rates decline to a
pre-defined threshold." Do not widen on intuition: "AI tools caused
experienced open-source developers to be 19% *slower* on average, despite
those developers believing they were 20% faster — a 39-percentage-point
perception gap." Promotion is telemetry-gated, never feel-gated.

## Semantic-versioning baseline

The most defensible first executable rule is dependency updates: "patch and
minor version updates with passing CI tests are the most defensible starting
point for auto-approve" — "the most mature, production-validated auto-approve
pattern in existence" (auto-merge of patch/minor "cutting maintenance time by
approximately 40% on average"). Extend from this baseline rather than from a
blank slate; major version bumps stay under human review.
