---
name: technical-feasibility
description: >
  Decides whether and how a requirement is buildable on the current stack, surfacing the options, dependencies, and risks before anyone commits to a design or a date. Use when the user asks "can we build this", "is this feasible on our stack", "what are the options and risks", or "do we need a spike/PoC". Produces a Feasibility Assessment with at least two pros/cons options, dependencies, risks with mitigations, a time-boxed spike decision, and a buildable/not-now/phased verdict. Do NOT use for localizing a bug (that belongs to progressive-bug-hunter) or security feasibility of a change (reviewing-software-security).
---

# technical-feasibility

## Purpose
Decide whether — and how — a requirement is buildable on the current stack, and surface the options, dependencies, and
risks before anyone commits to a design or a date. This is about *whether/how + risk*, not the final design.

## When to use
- Triggers: *"can we build this"*, *"is this feasible on our stack"*, *"what are the options and risks"*, *"do we need a spike/PoC"*.
- **Not this skill:** localizing a bug → `progressive-bug-hunter`; security feasibility of a change → `reviewing-software-security`.

## Input
- A clarified requirement (ideally the Requirement Analysis artifact) + current stack, constraints, and deadline.

## Output
A **Feasibility Assessment** artifact + checklist + verdict. Skeleton:

```
# Feasibility — <feature>
Verdict:  BUILDABLE  |  NOT NOW  |  PHASED
Options:
  A) <approach> — pros / cons / cost / maintenance
  B) <approach> — pros / cons / cost / maintenance
Dependencies:  (teams, services, libraries, data)
Risks:  <risk> — likelihood — mitigation
Spike needed?  yes/no — time-box — question it must answer
Timeline impact:  (rough, qualitative)
Recommended option + why:
```

## Decision rules
1. Judge any approach/stack by three lenses: **necessity** (does it make the team's work easier), **maintenance difficulty + community breadth**, and **cost reasonableness** (instance + dev-maintenance cost).
2. **Prefer boring tech.** Only risk new tech when the current stack hits a hard constraint (e.g. fails a security standard), costs more than needed, or its community has stopped supporting it.
3. Before committing to a new tech/protocol, analyze **maturity, community maintenance, and prevention methods** — a passing non-prod test is not proof.
4. **Time-box** a spike/PoC by infra-management sizing and the nature of the feature; give it one clear question to answer.
5. Use **data** when it exists; for early/no-data cases, intuition + a domain-model draft is acceptable.

## Checklist
- [ ] ≥2 options with pros/cons, cost, maintenance
- [ ] Dependencies (teams/services/libs/data) identified
- [ ] Risks listed with likelihood + mitigation
- [ ] Spike decision + time-box (if any)
- [ ] Verdict: buildable / not-now / phased
- [ ] Recommended option with reasoning

## Anti-patterns (never do)
- Recommend tech the team has never used for **critical** work.
- Choose tech for the résumé, not the team's need.
- Treat a non-prod success as production proof.
- Pick something because it's *new*.
- Answer the business on technical matters without the team knowing.

## Example
**Input:** "Stream live order updates to the app." **Output (excerpt):** *Options:* A) gRPC streaming — fast but
prod-risk + ops burden; B) WebSocket — proven for us. *Risk:* gRPC streaming failed in prod historically despite clean
non-prod tests. *Spike:* yes, 2 days, "does gRPC streaming hold at 100 RPS in prod-like net?" *Verdict:* PHASED —
ship polling now, migrate to WebSocket. *Recommended:* B.

## Human approval gate
**Stop.** For one-way-door tech choices, the team discusses the options and a human picks one before any design work
starts. The skill recommends; it does not adopt new tech on its own.
