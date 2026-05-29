---
name: delivery-planning
description: >
  Break a chosen solution into a sequenced, estimated execution plan — tasks, critical path, blockers, and a phased timeline honest enough to tell the business. Use when the user asks to "break this into tasks", "estimate this work", "what's the critical path", or "plan the delivery / timeline". Applies a disciplined estimation process, complexity multipliers, and the no-estimate-on-volatile-requirements rule, ending at a human commit gate. Do NOT use to size raw complexity or unknown-risk before a date is committed — that is `risk-estimation`.
---

# delivery-planning

## Purpose
Break a chosen solution into a sequenced, estimated execution plan — tasks, critical path, blockers, and a phased
timeline honest enough to communicate to the business.

## When to use
- Triggers: *"break this into tasks"*, *"estimate this work"*, *"what's the critical path"*, *"plan the delivery / timeline"*.
- **Not this skill:** negotiating scope politics is stakeholder-communication work; this skill produces the plan + estimate.

## Input
- A chosen architecture + scope (and any known team/calendar constraints).

## Output
A **Delivery Plan** artifact + checklist. Skeleton:

```
# Delivery Plan — <feature>
Tasks (split by interface → handler → business logic):
  - <task> — est. <man-days> — confidence <H/M/L> — depends on …
Sequence / critical path:
Blockers / dependencies (teams, external):
Assumptions:
Complexity multiplier applied:  <why>
Phased timeline for business:  phase 1 … / phase 2 …
```

## Decision rules
1. **Estimation steps**: analyze the business requirement → analyze the problem → assume the solving process → estimate man-days.
2. **Size up for complexity/uncertainty** — increase man-days by the requirement's complexity level.
3. **Don't estimate** when requirements are unclear, likely to change, or coupled to uncertain legacy — these guarantee a wrong estimate.
4. Baseline = a **mid-level developer**; usually the lead estimates for the team.
5. **Round-number estimates** mean the team hasn't thought it through; if a refactor hits 2× its estimate, stop and reassess.
6. Communicate long timelines by listing high-level tasks, time per step, and **discussing scope priorities with the business**.

## Checklist
- [ ] Tasks broken down (interface → handler → business logic)
- [ ] Sequence + critical path identified
- [ ] Blockers/dependencies (incl. external teams) listed
- [ ] Estimates with confidence levels
- [ ] Assumptions explicit
- [ ] Complexity multiplier applied where uncertain
- [ ] Phased timeline prepared for the business

## Anti-patterns (never do)
- Accept a timeline commitment without consulting the team.
- Estimate on unclear or volatile requirements.
- Hand over round-number estimates.
- Promise features without prior design.
- Promise a timeline without consulting the team.
- Answer the business on technical matters without the team knowing.
- Hide technical problems from management.

## Example
**Input:** wallet top-up. **Output (excerpt):** tasks — API contract (1d), provider integration (3d, **critical path**,
external risk), ledger write (1d), tests (2d). Total ~8 man-days ±3 (M confidence). Assumption: provider sandbox
available. Phase 1: top-up; phase 2: refunds.

## Human approval gate
**Stop.** The team and business agree on scope + timeline before it's committed; the lead never commits a date without
the team. The skill produces the plan + estimate; humans commit it.
