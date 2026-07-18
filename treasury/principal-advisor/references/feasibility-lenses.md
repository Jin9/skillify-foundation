# Feasibility lenses — judging whether an idea can and should be built

How to run an `assess` judgment. Verdict first, reasoning after, riskiest unknown always named.

## The three lenses

Judge any approach, stack, or idea through:

1. **Necessity** — does it make the team's work easier or the product possible, or is it complexity without a payer?
2. **Maintenance + community** — how hard is it to operate over time, and how broad and alive is the community that keeps it viable?
3. **Cost reasonableness** — build cost plus run cost plus developer-maintenance cost, weighed against the value it unlocks.

## Boring-tech rule

Prefer boring, proven technology. Risking new tech is justified only when the current approach hits a hard constraint (fails a security or scale requirement), costs clearly more than needed, or its community has stopped maintaining it. Before adopting anything new, examine maturity, community maintenance, and failure-prevention methods — a passing demo or non-prod test is not production proof.

## Spike discipline

When the riskiest unknown is empirical, propose a time-boxed spike: one clear question it must answer, a box sized to the feature's importance, and a decision rule for both outcomes. A spike without a single question is a hobby.

## Verdicts

- **BUILDABLE** — feasible on the current stack and team; state the recommended path.
- **NOT NOW** — feasible in principle, but a named blocker (cost, team capability, dependency, licensing) makes it wrong today; name what would change the answer.
- **PHASED** — buildable via a staged path; name phase 1 and the trigger for phase 2 (ship the boring version now, migrate when the measured need appears).

## Evidence posture

Use data when it exists. For early ideas with no data, informed intuition plus a rough domain sketch is acceptable — and say that is what it is. Always name the riskiest unknown; a feasibility answer that names no unknown is overconfident. State a confidence level on the verdict and name the cheapest evidence that would change it — a spike, a benchmark, one reference customer, a licensing answer.

## Ownership check

Ask who will own, operate, debug, upgrade, and eventually retire the choice. Feasibility without a named owner is a NOT NOW in disguise — the build is possible, the operation is not.

## Anti-patterns

- Recommending tech the team has never used for critical-path work.
- Choosing for the resume, or because it is new.
- Treating a non-prod success as production proof.
- A verdict with no named unknown and no condition that would flip it.
