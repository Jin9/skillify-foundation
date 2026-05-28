---
name: refactor-decision
description: >
  Decide whether a refactor is necessary, bound its scope, protect existing behavior with tests first,
  and avoid cosmetic refactors that have no business reason. Use when the user asks "should we refactor
  this", "how big should this refactor be", "is this legacy worth fixing", or "refactor or leave it".
  Produces a Refactor Decision artifact with a scoped boundary, a 2x stop-rule, and a checklist. Do NOT
  use for a from-scratch redesign decision and its ADR (that is a separate architecture-decision step).
---

# refactor-decision

## Purpose
Decide whether a piece of code earns a refactor, **bound its scope tightly**, protect existing behavior with tests before touching it, and avoid cosmetic refactors with no business reason — landing on DO-NOW / SCHEDULE / LEAVE+ticket.

## When to use
- Triggers: *"should we refactor this"*, *"how big should this refactor be"*, *"is this legacy worth fixing"*, *"refactor or leave it"*.
- **Not this skill:** a from-scratch redesign decision plus an ADR for it — that is an architecture-decision activity, not a scoped refactor of existing code.

## Input
- The code/legacy in question, the trigger (bug, debt, over-engineering), and any deadline.
- (Optional) deployment runbook, existing test coverage, PM time constraints.

## Output
A **Refactor Decision** artifact + checklist + decision. Skeleton:
```
# Refactor Decision — <area>
Trigger:           bug | debt | over-engineering
Impact class:      must-change | consider | nit
Scope boundary:    in:  …    out:  …          # peel off at the edges; keep it bounded
Behavior protection: tests written FIRST covering current behavior  [ ] yes
Effort:            <n> man-days   |  2x stop-rule: if it crosses 2x estimate → STOP & reassess
Can remove feature from deployment runbook?  yes → pull & plan  |  no → decide the risk explicitly
Decision:          DO NOW  |  SCHEDULE  |  LEAVE + tracking ticket   # never silently drop
Rollback / guard:  <how we revert + the fitness check that locks the invariant>
```

## Decision rules
1. A refactor earns priority **only by debt risk** — debt = risk of a problem; triage **pay-now** (blocks customers) / **tolerate** (works but sub-standard) / **live-with** (nit).
2. If the solution process is **more complex than the problem warrants**, the real fix may be removing over-engineering, not adding more.
3. Pre-launch critical bug needing a big refactor: understand the bug → size the refactor → discuss effort/timeline with PM → **ask whether the feature can be pulled from the deployment runbook**; if yes, pull it and plan the fix; if no, decide the risk explicitly, try to fix, else accept the risk on the record.
4. Legacy "broken but still running": estimate man-days → classify **must-change / consider / nit** → discuss a time slot with PM → check it fits the timeline → **if it can't be committed, leave it but record a tracking ticket** — make it visible, never silently drop.
5. **Simple first / YAGNI** — don't refactor for elegance with no business reason.
6. **Refactor crossing 2x the estimate → stop and reassess**.
7. Keep the refactor **bounded** — peel off at the edges along the coupling profile — and **lock the invariant with an automated fitness check** so it cannot silently regress.

## Checklist
- [ ] Trigger identified (bug / debt / over-engineering)
- [ ] Impact class set (must-change / consider / nit)
- [ ] Scope boundary written (in / out)
- [ ] Behavior-protecting tests written FIRST
- [ ] Effort man-days + 2x stop-rule stated
- [ ] Runbook-removal question answered
- [ ] Decision recorded; if LEAVE → tracking ticket filed

## Anti-patterns (never do)
- Cosmetic refactor with no business reason — violates simple-first/YAGNI.
- Refactor without a behavior-protecting test net, or push past 2x estimate without reassessing.
- Silently drop a refactor you can't commit instead of filing a visible ticket.
- "Solve" complexity by adding more complexity.
- Recommend a large refactor without an impact assessment.
- Delete or deprecate a service without a migration plan.

## Example
**Input:** "Our fintech payment-callback handler is tangled; a duplicate-charge bug is showing up days before launch — refactor or leave it?"
**Output (excerpt):** *Trigger:* bug. *Impact class:* **must-change** (blocks customers). *Scope boundary:* in = callback idempotency + ledger write; out = provider SDK swap. *Behavior protection:* characterization tests on current callback paths first. *Effort:* ~3 man-days; *2x stop-rule:* >6 → stop. *Runbook:* the new provider can be pulled from the deploy → recommend pull-and-plan. → **DO NOW** (scoped), invariant locked by a CI fitness check.

## Human approval gate
**Stop.** This skill **recommends a decision; it does not commit one.** The PM owns the timeline/runbook call and whether to accept residual risk; a large refactor never proceeds without an impact assessment, and the lead signs off the artifact rather than deciding the team's schedule alone.
