# Paper / Source Reading — Action List per Lens

Use this only when source material is complex (long paper, multi-section template, ambiguous domain doc).
For simple "create a skill from this" requests, compress lenses 1–3 into one paragraph and proceed directly to authoring.

---

## 1. Intake
Understand what the material is and why it is being read.

- Identify document type (paper / spec / template / workflow / notes / code)
- Identify objective and target audience
- Identify scope, owner/source, date/version
- Detect missing context or outdated content
- Extract table of contents, key sections, glossary/terminology

**Gate**: If purpose, scope, and key sections are clear → proceed to Extract. Otherwise, ask one clarifying question.

---

## 2. Extract
Pull useful information out of the material.

- Extract: key points, definitions, claims, evidence, process flow
- Extract: requirements, constraints, risks, dependencies, decision points
- Extract: roles, open questions, action items, examples

---

## 3. Understand
Convert the material into clear understanding for the target agent and user.

- Rewrite the core idea in plain, practical terms
- Map concepts to real use cases
- Show example or counter-example if terminology is unclear

**Compression rule**: If the material is under 2 pages or the skill task is straightforward, write one sentence per lens (1–3) and move on.

---

## 4. Analyze
Break down meaning, structure, and implications.

- Analyze structure, logic, trade-offs, pros and cons
- Analyze failure cases, edge cases, complexity
- Analyze security, scalability, maintainability impact (if relevant)

---

## 5. Assumptions
Identify what is not explicitly stated but required for decision-making.

- List assumptions; classify as: safe / risky / needs-confirmation
- Convert risky assumptions into explicit questions or skill constraints
- Identify hidden dependencies

**Rule**: If an assumption is reversible and low-impact → state it briefly and proceed. If it materially changes the skill's scope or reliability → surface it as an open question before continuing.

---

## 6. Validate
Check correctness, feasibility, and completeness.

- Validate logic, requirement completeness, technical feasibility
- Check contradictions, ambiguity, missing acceptance criteria
- Check compliance or security concerns if domain is regulated

---

## 7. Synthesize
Convert reading into skill design.

- Synthesize key ideas into a reusable workflow
- Create decision rules, progressive disclosure structure, validation loop
- Design file layout: `SKILL.md` + optional `references/`, `scripts/`, `assets/`

---

## Minimum Must-Do for Any Source Material

1. Identify purpose and scope
2. Extract key requirements, rules, and constraints
3. Extract assumptions and risks
4. Validate feasibility and gaps
5. Synthesize into an actionable skill design
6. Summarize assumptions made (briefly, inline in the skill draft)

---

## Strong Agent Command Pattern

When asking another model to read source material before skill creation:

> "Read the attached document. Then:
> 1. Intake: identify purpose, scope, audience, and version.
> 2. Extract: list key concepts, requirements, constraints, assumptions, risks, and dependencies.
> 3. Analyze: explain trade-offs, edge cases, and failure cases.
> 4. Validate: check ambiguity, contradiction, feasibility, and missing information.
> 5. Synthesize: convert useful ideas into a reusable skill workflow.
> 6. Summarize: produce a concise markdown draft for review."
