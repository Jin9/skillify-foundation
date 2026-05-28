# v-planner Template

**Parent:** [`../templates.md`](../templates.md)
**Owner role:** [Iteration Planner](../roles.md) · **`template_version`:** 0.1.0

The execution plan for one **squad-design iteration** (v0.7, v0.8). One per iteration. Drafted from `FEEDBACK.md` + `STRATEGY.md` at iteration start; user approves before workflow runs are accepted against the iteration's docs.

**File location:** `<repo_root>/docs/v0.X-plan.md`.

> Distinct from Plan-Reviewer's `plan-review.json` (workflow-scoped red-team). See [`../roles.md`](../roles.md) for the disambiguation.

---

## Required sections

1. `# v0.X Execution Plan` + frontmatter (see below).
2. `## Iteration objective` — 2-3 sentences.
3. `## In-scope items` — table grouped by P-tier; per row: ID (e.g., `FEEDBACK#P0-1`), title, source, effort, depends-on, DoD.
4. `## Deferred items` — short list with one-line rationale (or just point at FEEDBACK/STRATEGY P-tier).
5. `## Acceptance gates` — concrete pass/fail for "v0.X is done".
6. `## Cost estimate` — single envelope (e.g., `$60-90`) anchored on previous run actuals.
7. `## Risk register` — top 3 named risks with mitigation.
8. `## Decisions resolved` — `decisions.md` open items closed by this iteration.
9. `## Change log`.

Drop the dependency graph if it's <3 edges; drop the risk register if previous iteration's risks were all mitigated.

---

## Frontmatter

```yaml
---
template_version: 0.1.0
iteration: v0.X
status: Drafted | Approved | InProgress | Done | Superseded
inputs: { feedback: FEEDBACK.md, strategy: STRATEGY.md, previous_plan: ..., previous_state: ... }
in_scope_count: { P0: N, P1: N, P2: N, P3: N }
expected_cost_usd: "60-90"
expected_dry_run_terminal: Done | ShipWithCaveats
change_log: [{ date, author, change }]
---
```

---

## Negative example — Plan with no acceptance gates

```markdown
# v0.7 Execution Plan
## In-scope items: …big list…
## Cost estimate: TBD
```

What's wrong:
1. No `Acceptance gates` — operator can't tell when v0.X is "done." Iterations drift indefinitely.
2. `Cost estimate: TBD` — no budget anchor; user can't approve.
3. No risk register, even minimal — known failure modes (deviation accumulation, budget overrun, dep-blocked items) aren't named, so they're not mitigated.

If the dependency graph is non-trivial (P0 unblocks P1, P1 unblocks deferred-P2), name the 2-3 critical edges in prose; don't draw the whole graph.
