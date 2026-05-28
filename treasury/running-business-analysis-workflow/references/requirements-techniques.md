# Requirements techniques

Field-tested techniques for writing, splitting, prioritizing, and tracing requirements. Use them in Stages 3, 4,
and 6. Pick by fit, not dogma.

## User stories

### Anatomy — Connextra + the 3 Cs
- Format: **"As a `<role>`, I want `<capability>`, so that `<benefit>`."** The benefit clause is the part that
  reveals real intent — never drop it.
- **3 Cs:** *Card* (the short story), *Conversation* (the detail that lives in discussion, not in the card), and
  *Confirmation* (the acceptance criteria and tests). A story is a placeholder for a conversation, not a spec.

### INVEST quality gate
Independent, Negotiable, Valuable, Estimable, Small, Testable. In practice the load-bearing three are **Small,
Valuable, Testable** — a story that is small, clearly valuable, and testable is usually ready; relax Estimable
when needed.

### SPIDR — five ways to split a too-big story
- **Spike** — split off a research/learning task when uncertainty blocks estimation.
- **Path** — separate distinct workflow variants (happy path vs alternate paths).
- **Interface** — split by platform, device, or channel.
- **Data** — slice by data set or data variation.
- **Rules** — temporarily relax a business rule to ship a thinner slice, then add the rule back.
Meta-pattern: reduce variations to one **vertical slice** that still delivers value end-to-end.

## Acceptance criteria

### Format selection
- **Given/When/Then** for behavior with branches and multiple outcomes, especially when it will be automated.
- **Rule-oriented checklist** for validation, UI states, and non-functional checks.
Choose by the story's risk and shape, not habit.

### Quality bar
- Testable, unambiguous (no "fast", "intuitive"), independent (each can be checked alone), and complete
  (includes edge cases and error conditions).
- **Heuristic: 3–6 crisp criteria per story.** More than ~8 usually means the story should be split.
- Calibrated detail: detailed enough to be testable, not so prescriptive that it dictates the solution.

### Example Mapping
A 20–30 min Three Amigos (business, dev, QA) session producing cards for the story, its rules, concrete examples,
and open questions. A growing open-questions pile is the signal the story is not ready.

## Prioritization — and each method's blind spot
- **MoSCoW** (Must / Should / Could / Won't): fast stakeholder alignment. *Blind spot:* "Must" inflation —
  timebox and challenge every Must.
- **RICE** (Reach × Impact × Confidence ÷ Effort): comparable scoring across items. *Blind spot:* subjective
  inputs and undervaluing strategic-but-low-reach work.
- **WSJF** (Cost of Delay ÷ Job Size): sequences by urgency and economics. *Blind spot:* needs a credible
  cost-of-delay estimate.
- **Kano** (Basic / Performance / Excitement): centers satisfaction. *Blind spot:* delighters decay into
  expectations over time.
None of these model dependencies — sequence those separately.

## Scope & MVP
- **Story mapping:** lay out the *backbone* (user activities left-to-right in narrative order) with *ribs*
  (details hanging down in priority order). The top band across the whole backbone is the **walking skeleton** —
  the smallest end-to-end usable system.
- **Slice horizontally, not by layer.** Release slices cut across the whole journey (a thin path end-to-end),
  not "all the UI first". A UI-only or DB-only slice fails the independent-and-valuable test.
- **Core-problem cut test:** "Can users still solve their core problem without this?" If yes, cut it from the MVP.
- **Riskiest Assumption Test:** before building, test the single belief most likely to kill the product — often
  with no code at all.
- **Scope-creep control:** documented scope with explicit exclusions + a formal change process + timeboxing.

## Traceability
- **Directions:** *forward* (requirement → design → code → test, the coverage claim), *backward*
  (test/code → requirement, the rationale), and *bidirectional* (both, full visibility).
- **Requirements Traceability Matrix (RTM):** a living table — requirement ID, description, status, linked
  test-case IDs. Its binding constraint is upkeep; trace as much as the impact-analysis or regulatory payoff
  justifies and no more.
- **Impact analysis:** when a requirement changes, follow its links to find affected design, code, and tests, and
  estimate rework before committing.

## Working with AI on requirements — pitfalls & controls
- **Jagged frontier:** AI helps most on routine drafting (boilerplate stories, first-pass AC) and can *hurt* on
  off-frontier judgment (regulatory edge cases, novel domains). Sort the task before trusting the output.
- **Hallucinated / smoothed requirements:** models invent plausible-but-wrong criteria and quietly resolve real
  ambiguity. Force them to emit definitions, examples, and exceptions — not confident prose — and keep a human
  review gate on semantics and completeness.
- **Lost discovery context:** generated stories drift generic when the persona, epic, and constraints aren't fed
  in. Ground every draft in real artifacts (problem statement, research, current docs).
- **Closed-world prompting:** restrict the model to the authorized source documents; don't let it fill gaps from
  general knowledge. Unknowns become open questions, not invented answers.
