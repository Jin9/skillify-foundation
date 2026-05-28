---
name: architecture-decision
description: >
  Chooses a design among options and records it as an ADR with explicit trade-offs, protecting long-term maintainability and setting clear direction the team can follow. Use when the user asks "should we split this into microservices", "write an ADR", "which design should we pick", or "is this over-engineered". Produces an ADR (plus a C4 design doc for larger systems) comparing at least two options, tagging one-way versus two-way doors, and noting consequences plus rollback triggers. Do NOT use for diff-correctness review (use the built-in code-review) or security threat-modeling (reviewing-software-security).
---

# architecture-decision

## Purpose
Choose a design among options and record it as an ADR with explicit trade-offs — protecting long-term maintainability
and setting clear direction the team can follow.

## When to use
- Triggers: *"should we split this into microservices"*, *"write an ADR"*, *"which design should we pick"*, *"is this over-engineered"*.
- **Not this skill:** diff-correctness review → built-in `/code-review`; security threat-modeling → `reviewing-software-security`.

## Input
- Requirement + NFRs + ≥2 candidate designs (or the skill proposes them).

## Output
An **ADR** (+ a **C4 design-doc** for larger systems). Skeletons:

```
# ADR-<n>: <decision>
Context / forces:
Options considered:  A … | B … | C …
Decision:
Trade-offs accepted:
Consequences:
Door type:  ONE-WAY (decide-then-explain)  |  TWO-WAY (reversible)
Rollback / revisit trigger:
```
```
# Design doc (C4)
L2 Network view:        L3 High-level components:        L4 Entity + business flow:
Components sheet:  microservices · database · API mapping
```

## Decision rules
1. **Design order**: business logic + volume → NFRs → scope by business domain (responsibility per boundary) → then trade off stack, number of services, and standards.
2. **Microservices need a reason**: which business scope does each serve? is separation a priority *now*? does the team have capacity to maintain the split? Validate with a perf scenario if unsure. Beware "microservices everywhere" / "3-tier every time" — too many pieces to maintain.
3. **Over-engineering signal**: when the solution process gets more complex than the problem size warrants, stop. If code needs 5 minutes to explain, it's over-engineered.
4. **Tag the door**: for one-way-door decisions (stack choice, data-migration strategy, structural splits), decide-then-explain — set a clear direction, then explain the rationale; for reversible two-way doors, adjust freely. Weight decide-then-explain **more** on one-way doors — under-weighting it leads to no design standard and massive debt. Either way, a human signs off.
5. **Impact sets the process**: whole-team impact → discuss/RFC; the lead writes the high-level (L3) doc, the team owns detail + API spec, and the **lead signs off**.

## Checklist
- [ ] ≥2 options compared with trade-offs
- [ ] NFRs respected; maintainability + bus-factor considered
- [ ] Door type tagged (one-way / two-way)
- [ ] Consequences + rollback/revisit trigger noted
- [ ] Decision recorded as an ADR (and C4 doc if large)
- [ ] Right process used for the impact level (solo / discuss / RFC)

## Anti-patterns (never do)
- Make an important architecture decision **alone**, without the team.
- Recommend tech the team has never used for **critical** work.
- Choose tech for the résumé, not team needs.
- Over-engineer or reach for microservices prematurely.
- Fail to document the decision, or justify it only with "we've always done it this way".
- Delete or deprecate a service without a migration plan.
- Override a team decision without explanation.

## Example
**Input:** 5-person team, moderate traffic, proposal to split into 4 microservices. **Output (excerpt):** Ask the three
questions → only the payments scope justifies its own service now; keep the rest a modular monolith. *Door:*
two-way (can split later). *ADR* records the decision + the revisit trigger (traffic 5×). Team discusses; the lead signs off.

## Human approval gate
**Stop.** The lead signs off on the ADR; one-way-door decisions get explicit team discussion before commit.
The skill drafts and recommends — it does not finalize architecture autonomously.
