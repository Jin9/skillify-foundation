---
name: running-business-analysis-workflow
description: >
  Guide an AI agent through a structured business-analysis workflow for software and IT projects: define the problem and goals, analyze the current (as-is) state, elicit functional and non-functional requirements, design the future (to-be) state, assess feasibility, then validate and document. Use when the user asks to "do a business analysis", "gather requirements", "write a BRD or SRS", "create user stories with acceptance criteria", "run a requirements workshop", or "analyze an as-is / to-be process". Detects whether it acts as the BA (autonomous) or as a co-pilot to a human BA, knows which checkpoints the BA owns versus the tech lead, analyst, PM, or QA, and flags items that need a human owner instead of guessing. Do NOT use for writing the application code itself, for non-software business consulting unrelated to a system, or for project scheduling and resourcing.
---

# running-business-analysis-workflow

## Purpose
Drive a software/IT project from a vague problem to a validated, documented requirements set, owning the
business-analysis work while flagging every checkpoint that belongs to a human tech lead, analyst, PM, or QA.

## When to use
- Triggers: *"do a business analysis"*, *"gather requirements"*, *"write a BRD or SRS"*, *"create user stories with acceptance criteria"*, *"run a requirements workshop"*, *"analyze an as-is / to-be process"*.
- **Not this skill:** writing the application code itself → use a coding/implementation workflow; non-software business consulting unrelated to a system → out of scope; project scheduling, staffing, or resourcing → use project-management tooling.

## Mode detection (do this before Stage 1)
Read the context and pick one profile; it changes how every stage behaves.
- **Autonomous BA** — cues: *"you're the BA"*, *"run the analysis"*, or no human BA is mentioned. The agent acts as the analyst.
- **Co-pilot** — cues: *"help me draft…"*, *"review my requirements"*, *"as the BA I want to…"*. A human BA owns the work; the agent assists.
- If genuinely ambiguous, **ask once**, then proceed. State the detected mode at the top of the output.

## Core rule (both modes)
Never invent an answer the agent cannot actually determine — technical feasibility, effort estimates, test results.
**Escalate instead:** insert a labeled placeholder and name the human owner. Fabrication is the cardinal failure here.

## The six-stage workflow
Each stage has the same shape: **Entry → Ask → Capture → Owner → Exit → Output.** Keep this file lean; pull depth from
`references/` on demand (`elicitation-questions.md` for the per-stage question bank, `requirements-techniques.md` for
technique detail, `nfr-catalog.md` for NFR targets, `raci-ownership.md` for role mapping).

### Stage 1 — Define problem & goals
- **Entry:** a request, pain point, or opportunity exists; no agreed problem statement yet.
- **Ask:** What problem, for whom? What measurable business goal defines success? Why now? What is explicitly out of scope? Who are the stakeholders?
- **Capture:** one-paragraph problem statement; ≥1 measurable goal/success metric; stakeholder list; scope boundary; stated assumptions.
- **Owner:** BA frames the problem. PM owns business priority and funding. Tech lead confirms enough technical context exists.
- **Exit:** problem statement + at least one measurable goal agreed.
- **Output:** *Problem & Goals* section of the BRD/SRS (`templates/brd-srs.md`).

### Stage 2 — Analyze current state (as-is)
- **Entry:** problem statement agreed.
- **Ask:** How is this done today — actors, steps, systems? Where are the pain points and workarounds? What data and integrations exist? What are today's volumes?
- **Capture:** as-is process (actors/steps/systems); pain points; current data and integration touchpoints; baseline metrics.
- **Owner:** BA documents the as-is. A domain SME validates the process. Tech lead/architect confirms system facts.
- **Exit:** as-is map validated by an SME; baseline metrics noted.
- **Output:** *As-Is* section of the BRD/SRS.

### Stage 3 — Gather requirements (heaviest — split functional vs non-functional)
- **Entry:** as-is understood.
- **Ask (functional):** What must the system do, per role? What are the acceptance criteria, edge cases, error paths, and business rules?
- **Ask (non-functional):** Performance (latency/throughput numbers)? Availability target? Scale and growth? Security/compliance? Usability? (`references/nfr-catalog.md`)
- **Capture:** functional requirements as user stories with **testable** acceptance criteria (`templates/user-story.md`); NFRs with **concrete numeric targets**; a priority per item; a traceability ID per requirement.
- **Owner:** BA elicits and writes. QA reviews AC for testability. Tech lead/architect supplies feasibility and NFR realism — the BA does not invent these. PM sets priority.
- **Exit:** every story has testable AC; **every NFR has a number**; priorities assigned.
- **Output:** requirements set folded into the BRD/SRS; techniques in `references/requirements-techniques.md`, questions in `references/elicitation-questions.md`.

### Stage 4 — Design future state (to-be)
- **Entry:** requirements set drafted.
- **Ask:** What is the to-be process? Which requirement changes which step? What new data or integrations appear? What is the MVP slice?
- **Capture:** to-be process; the as-is → to-be gap; MVP boundary (core-problem cut test); downstream impacts.
- **Owner:** BA designs the to-be *process*. Architect/tech lead owns the technical solution shape — the BA does not decide architecture. PM confirms MVP scope.
- **Exit:** to-be process mapped; gap and MVP slice agreed.
- **Output:** *To-Be* section + MVP scope in the BRD/SRS.

### Stage 5 — Assess feasibility
- **Entry:** to-be and requirements drafted.
- **Ask (to humans):** Is each requirement technically feasible? Rough effort band? Key risks and dependencies? Anything infeasible or needing a spike?
- **Capture:** a feasibility verdict per requirement (from the tech lead); risks and dependencies; items needing a spike.
- **Owner:** Tech lead/architect owns feasibility and effort — **the BA never fabricates these.** PM owns timeline/resourcing (out of this skill). QA flags test risk.
- **Exit:** every requirement marked feasible / needs-spike / infeasible **by a human owner**; no invented verdicts.
- **Output:** feasibility and risk annotations on the requirements in the BRD/SRS.

### Stage 6 — Validate & document
- **Entry:** requirements and feasibility gathered.
- **Ask:** Do stakeholders agree the requirements are complete and correct? Any conflicts? Is everything traceable goal → requirement → AC? Who signs off?
- **Capture:** validation results; conflict resolutions; a traceability matrix; the RACI (`templates/raci-matrix.md`, `references/raci-ownership.md`); open items.
- **Owner:** BA assembles and runs validation. QA confirms testability and coverage. Tech lead confirms technical sections. PM/sponsor signs off — a human decision.
- **Exit:** BRD/SRS complete, requirements checklist passed, RACI filled, sign-off routed to a human.
- **Output:** completed BRD/SRS (`templates/brd-srs.md`) + requirements checklist (`templates/requirements-checklist.md`) + RACI (`templates/raci-matrix.md`).

## Decision rules
1. **BA-owned items:** autonomous mode produces them directly; co-pilot mode drafts and proposes for the human BA's approval.
2. **Non-BA items (tech lead / architect / PM / QA):** autonomous mode inserts a labeled `[NEEDS INPUT — <role>]` placeholder and asks the user to route it; co-pilot mode notes the typical owner as a suggestion. **Neither fabricates.**
3. **Sign-off:** autonomous mode tracks completeness but flags final sign-off as a human decision; co-pilot mode always defers to the human BA.
4. **NFRs must be numeric.** "Fast" is not a requirement; "p95 < 300 ms at 50 req/s, 99.9% monthly availability" is. No number → not done (`references/nfr-catalog.md`).
5. **Prioritize explicitly** (MoSCoW by default); record the priority on each requirement.
6. **Trace everything**: every requirement links back to a goal and forward to its acceptance criteria.

## Output
```
Artifacts (each maps to a template):
- BRD/SRS document ............ templates/brd-srs.md          # all six stages, assembled
- User stories + AC .......... templates/user-story.md       # one per functional requirement
- Requirements checklist ..... templates/requirements-checklist.md  # six-stage completeness gate
- RACI matrix ................ templates/raci-matrix.md       # who is R/A/C/I per checkpoint
Header on every artifact: detected MODE (autonomous | co-pilot) + current stage.
```

## Checklist
- [ ] Mode detected and stated (autonomous vs co-pilot)
- [ ] Problem statement + ≥1 measurable goal (Stage 1)
- [ ] As-is validated by an SME (Stage 2)
- [ ] Functional requirements as user stories with testable AC (Stage 3)
- [ ] Every NFR has a concrete number (Stage 3)
- [ ] To-be process + MVP slice agreed (Stage 4)
- [ ] Feasibility verdict per requirement from a human owner — none invented (Stage 5)
- [ ] Validation done, traceability + RACI filled, sign-off routed to a human (Stage 6)
- [ ] Every non-BA item carries a labeled owner placeholder, not a guess

## Anti-patterns (never do)
- Fabricate technical feasibility, effort estimates, or test results — escalate with a labeled placeholder instead.
- Decide architecture, commit, deploy, or sign off — this skill produces and recommends only.
- Skip mode detection, or switch modes mid-flow without saying so.
- Accept vague NFRs ("fast", "scalable") with no number.
- Convert ambiguity into a silent assumption instead of an open question for the named owner.
- Cite the research source material or any external source inside the artifacts; write the rule, not its provenance.

## Example
**Input (autonomous):** "You're the BA — run the analysis for adding saved payment cards to checkout."
**Flow (excerpt):** *Mode:* autonomous. *Stage 1:* Goal — cut checkout drop-off; success = +5% conversion. *Stage 3:*
story "As a returning shopper, I want to save a card so checkout is faster," AC in Given/When/Then; NFR "tokenize within
p95 < 200 ms, PCI-DSS scope confined to the vault." *Stage 5:* feasibility of the tokenization vault →
`[NEEDS INPUT — tech lead]`, routed to the user (never guessed). *Exit:* BRD/SRS + checklist assembled; sign-off flagged
as a human decision.

## Human approval gate
**Stop at sign-off.** This skill produces artifacts and recommends; it never decides, commits, or deploys. In autonomous
mode, final sign-off is a human decision and every non-BA verdict (feasibility, effort, test outcomes) must come from its
named owner. In co-pilot mode, defer to the human BA for approval of every BA-owned artifact.
