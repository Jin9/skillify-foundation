---
name: defining-engineering-standards
description: >
  Define and own the team's reusable standard — tech-stack baseline, DDD domain boundaries, a C4 L3/L4 design skeleton, API naming and folder/layer conventions, and the project scaffold — so the team builds consistently and can run independently. Use when the user asks "define our team's standard template", "what's our tech-stack baseline", "set our API naming and folder conventions", or "set up a new project's standards from scratch". Produces a standard-template artifact plus a checklist and a lead sign-off gate, with the stack chosen on necessity, maintenance, and cost as a decide-then-explain one-way door, and the conventions that conformance reviews later check against. Do NOT use to generate the actual project skeleton code from the standard (use your scaffolding tool or framework generator) or for security standards and threat-modeling (reviewing-software-security).
---

# defining-engineering-standards

## Purpose
Define and own the team's reusable **standard template** — tech-stack baseline, DDD domain boundaries, a C4 L3/L4
design skeleton, API naming and folder/layer conventions, and the project scaffold — so the team builds consistently
and can run independently. This produces the standard the team develops against and that conformance reviews check
against; it does not make one-off, single-system design decisions.

## When to use
- Triggers: *"define our team's standard template"*, *"what's our tech-stack baseline"*, *"set our API naming and folder conventions"*, *"set up a new project's standards from scratch"*.
- **Not this skill:** generating the actual project skeleton code from the standard → your scaffolding tool / framework generator (e.g. cookiecutter, a `create-*` CLI); security standards / threat-modeling → `reviewing-software-security`.

## Input
- A new or existing team/project context: business domains + expected volume (for NFRs), the current or candidate stack,
  and the conventions to standardize (or the skill proposes a baseline).

## Output
A **Standard Template** artifact + checklist + sign-off. Skeleton:

```
# Standard Template — <team / project>
Tech-stack baseline:   <language · framework · datastore · messaging · infra — chosen on necessity, maintenance/community, cost; door: ONE-WAY (decide-then-explain)>
Domain & boundaries:   <DDD: business domains -> aggregate boundaries -> responsibility per boundary>
C4 design skeleton:    L2 network · L3 high-level components (lead writes) · L4 entity + business-flow (team, in scope)
Components sheet:      <services · datastores · API mapping>
API conventions:       naming /api/v1/<domain>/<aggregate>/<action> · request/response/error/versioning shape · change-log
Folder / layer layout: <project scaffold — layers, module structure, where business logic lives>
Coding standards:      <lint/format · test layout · naming · the "matches template" review bar>
Ownership process:     <lead writes L3 · team owns detail + API spec (owner-driven) · squad-lead + team review · lead signs off>
Revisit trigger:       <when the standard re-opens — stack hits a hard constraint · costs too much · community abandoned>
Verdict:  TEMPLATE READY for lead sign-off  |  GAPS: <…>
```

## Decision rules
1. **Get domain design right before standardizing anything else.** Order: business logic + expected volume → NFRs → scope by business domain (DDD: responsibility per aggregate boundary) → only then set the stack, service count, and coding standards. Wrong domain boundaries propagate debt through the whole system.
2. **The stack baseline is a one-way door — decide-then-explain.** Choose on: necessity for the team (does it make work easier), maintenance difficulty + community breadth, and cost (instance + dev maintenance). Adopt new or cutting-edge tech only when the existing stack hits a hard constraint, costs more than needed, or its community has stopped supporting it.
3. **Carry an owner-driven C4 skeleton.** The standard includes L2 (network), L3 (high-level components), L4 (entity + business-flow), and a components sheet (services, datastore, API mapping). The lead writes the high-level L3; the team owns the detail (L4) and the API spec.
4. **One naming and layout convention.** API `/api/v1/<domain>/<aggregate>/<action>` with a fixed request/response/error/versioning shape and a change-log; a fixed folder/layer layout so every module looks the same and business logic sits in a predictable place. This is the bar conformance reviews check against.
5. **Make the team self-sufficient, not just documented.** The standard exists so the team builds consistently and runs independently (DDD → aggregate boundary → standard template → monitored autonomy) — pair the doc with the project scaffold and a learning framework, not a wiki page alone.
6. **Review and sign off.** The squad-lead + team review the standard (critical issues the lead discusses directly, minor ones the team decides); the **lead signs off**. Whole-team-impact changes (stack, structure) are discussed openly, never changed silently.

## Checklist
- [ ] NFRs derived from business logic + expected volume
- [ ] Domains scoped by DDD; aggregate boundaries + responsibilities defined
- [ ] Stack baseline chosen on necessity / maintenance + community / cost; door tagged one-way
- [ ] C4 skeleton present (L2/L3/L4 + components sheet); lead owns L3, team owns L4 + API spec
- [ ] API naming `/api/v1/<domain>/<aggregate>/<action>` + request/response/error/versioning/change-log shape
- [ ] Folder/layer layout + coding standards defined (the "matches template" bar)
- [ ] Project scaffold + learning framework pair the doc (team self-sufficient)
- [ ] Revisit trigger stated (stack constraint / cost / abandoned community)
- [ ] Reviewed (squad-lead + team); lead signed off

## Anti-patterns (never do)
- Ship code conventions as a standard before domain design is right — wrong boundaries propagate debt system-wide.
- Standardize on tech the team has never used for critical work, or pick the stack for the résumé over team need and cost.
- Over-compromise to consensus without weighting decide-then-explain on one-way-door choices (stack, structure) — the result is no design standard, unclear domain authority, and massive debt.
- Leave the standard as a doc with no scaffold or learning framework — the team won't actually run independently.
- Over-engineer the template (the "problem size vs. solution complexity" test).
- Change a whole-team standard (stack, structure) silently, without team discussion and a lead sign-off.
- Let API naming or folder layout drift per developer — the conformance bar becomes meaningless.

## Example
**Input:** new project, 6-person squad, payments domain, moderate volume. **Output (excerpt):** NFRs from volume →
DDD aggregates `wallet / ledger / payment`; stack baseline chosen on maintenance + cost (one-way → decide-then-explain);
C4 L3 by the lead, L4 per aggregate by the team; API shape `POST /api/v1/payment/wallet/top-up` + change-log; folder
layout `.../<domain>/<aggregate>/{handler,domain,repo}`; scaffold repo + onboarding checklist pair the doc. Reviewed by
squad-lead + team; **lead signs off.** **Verdict: TEMPLATE READY** once the message-queue choice names a maintenance owner.

## Human approval gate
**Stop.** The squad-lead + team review the standard and the **lead signs off** before it becomes the team's template;
one-way-door choices (stack, structure) get explicit team discussion first. The skill drafts and recommends the
standard — it does not finalize or enforce it autonomously.
