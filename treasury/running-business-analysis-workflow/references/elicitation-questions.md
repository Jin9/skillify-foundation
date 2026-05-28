# Elicitation question bank (per stage)

Use these to drive interviews, workshops, and document review. Ask open questions first, then close with
specifics. Capture answers verbatim where wording matters (rules, thresholds, exclusions). When an answer
is a guess, mark it as an assumption and route it to the owner — never promote a guess to a fact.

## Elicitation patterns
- **Conversational, iterative interviewing.** Ask one topic at a time, reflect the answer back, then probe the
  gap ("you said usually — when is it *not* usual?"). Follow-ups surface far more requirements than a flat checklist.
- **Mine existing inputs first.** Tickets, support logs, emails, chat threads, current docs, and screen recordings
  are pre-existing requirement sources; harvest them before scheduling new interviews.
- **Example Mapping / Three Amigos.** Bring business + dev + QA together around a story. Use four card colors:
  the story, its business rules, concrete examples, and **open questions**. A pile of open-question cards means the
  story is not ready — that is a signal, not a failure.
- **Workshop hygiene.** One facilitator per ~10 participants; smaller for remote (~6 for deep engagement). Timebox,
  make the timer visible, name a single decision-maker, and end each block with "what did we just decide / still owe?"

## Stage 1 — Define problem & goals
- What problem are we solving, and for whom (which user/role/segment)?
- How do we know it's a real problem — what evidence or pain exists today?
- What does success look like as a **measurable** outcome (metric + target + timeframe)?
- Why now? What changes if we do nothing?
- What is explicitly **out** of scope for this effort?
- Who are the stakeholders, sponsors, and impacted teams?
- What hard constraints exist (budget ceiling, deadline, regulation, platform)?

## Stage 2 — Analyze current state (as-is)
- Walk me through how this is done today, step by step — who does what, in which system?
- Where do people slow down, make errors, or invent workarounds?
- What data is created, read, or moved? Where does it live?
- Which systems or third parties are already integrated into this flow?
- What are today's volumes (users, transactions, peak vs average)?
- What works well today that we must **not** break?

## Stage 3 — Gather requirements
**Functional**
- What must the system let each role *do*? (Frame as "As a <role>, I want <capability> so that <benefit>.")
- For each capability: what are the acceptance criteria — how do we know it's done and correct?
- What are the edge cases, error paths, and "what if it fails halfway" scenarios?
- What business rules, validations, or policies constrain it?
- What must happen on success? On failure? On retry?

**Non-functional** (demand numbers — see `nfr-catalog.md`)
- Performance: expected and peak load? acceptable latency (p95/p99)? throughput?
- Availability: target uptime? acceptable downtime window? RTO/RPO?
- Scale: data growth rate? user growth over 12–24 months?
- Security & compliance: what data is sensitive? which regulations apply? authn/authz needs?
- Usability, accessibility, localization, observability, maintainability targets?

## Stage 4 — Design future state (to-be)
- In the improved flow, what does each step look like — who/what does it now?
- Which requirement changes which as-is step (the gap)?
- What new data, integrations, or interfaces appear?
- What is the smallest slice that still solves the core problem (the MVP)?
- What can be deferred to a later phase without blocking the core?

## Stage 5 — Assess feasibility (ask the technical owners)
- Is each requirement technically feasible with the current stack and constraints?
- What is the rough effort band, and what drives it?
- What are the key technical risks and external dependencies?
- Is anything infeasible as stated, or does it need a spike/prototype first?
- Which NFR targets are realistic, and which need negotiation?

## Stage 6 — Validate & document
- Do stakeholders agree the requirements are complete and correct? What's missing?
- Are there conflicts or contradictions between requirements or stakeholders?
- Can every requirement be traced back to a goal and forward to its acceptance criteria?
- Who is Responsible / Accountable / Consulted / Informed for each checkpoint?
- Who has the authority to sign off, and what do they need to see first?
