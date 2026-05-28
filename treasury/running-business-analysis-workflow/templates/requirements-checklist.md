# Requirements checklist — six-stage completeness gate

> Mode: <autonomous | co-pilot>. A stage is "done" only when its exit condition is met. Items the BA
> does not own carry a `[NEEDS INPUT — <role>]` placeholder, not a guess.

## Stage 0 — Mode
- [ ] Operating mode detected and stated (autonomous BA vs co-pilot)

## Stage 1 — Define problem & goals
- [ ] Problem statement written in one paragraph
- [ ] At least one measurable goal (metric + target + timeframe)
- [ ] Stakeholders and sponsor identified
- [ ] In-scope and out-of-scope both listed
- [ ] Constraints and assumptions captured
- **Exit:** problem statement + ≥1 measurable goal agreed

## Stage 2 — Analyze current state (as-is)
- [ ] As-is process mapped (actors → steps → systems)
- [ ] Pain points / workarounds captured
- [ ] Current data and integrations listed
- [ ] Baseline volume metrics noted
- [ ] As-is validated by an SME
- **Exit:** SME-validated as-is map + baseline metrics

## Stage 3 — Gather requirements
- [ ] Functional requirements written as user stories
- [ ] Each story has testable acceptance criteria (3–6)
- [ ] Edge cases and error paths covered
- [ ] Every NFR has a concrete number + condition (see nfr-catalog.md)
- [ ] Priority assigned to each requirement (MoSCoW)
- [ ] Traceability ID on each requirement
- **Exit:** testable AC on every story; a number on every NFR; priorities set

## Stage 4 — Design future state (to-be)
- [ ] To-be process mapped
- [ ] As-is → to-be gap identified
- [ ] MVP slice defined (core-problem cut test)
- [ ] Deferred items recorded
- **Exit:** to-be process + MVP slice agreed

## Stage 5 — Assess feasibility
- [ ] Feasibility verdict per requirement (from tech lead — not invented)
- [ ] Rough effort band captured
- [ ] Risks and dependencies recorded
- [ ] Spike/prototype items flagged
- **Exit:** every requirement marked feasible / needs-spike / infeasible by a human owner

## Stage 6 — Validate & document
- [ ] Stakeholders confirm completeness and correctness
- [ ] Conflicts resolved
- [ ] Traceability complete (goal → requirement → AC)
- [ ] RACI filled, one Accountable per row
- [ ] BRD/SRS assembled
- [ ] Sign-off routed to a human decision-maker
- **Exit:** validated BRD/SRS + filled RACI + sign-off routed

## Final
- [ ] No fabricated feasibility, effort, or test results anywhere
- [ ] Every non-BA item carries a labeled owner placeholder
