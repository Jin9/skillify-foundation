---
name: business-logic-extractor
description: >
  Extract the implemented business logic and rules FROM a codebase —
  cross-referenced with available requirements and agent/execution traces —
  into a faithful, traceable specification that BOUNDS information loss:
  salient rules, edge cases, and decision paths preserved with code provenance
  (file:line), not a lossy summary. Direction is code → spec (reverse). Use
  when the user asks to "extract the business rules from this module /
  service", "what business logic does this code implement",
  "reverse-engineer the spec / requirements from this codebase", "document
  the decision logic without dropping edge cases", or "produce a traceable
  rules spec from this implementation". Output: a structured rules spec plus
  an explicit information-loss / coverage ledger. Do NOT use to generate new
  code or new requirements (the forward requirement-to-code direction is out
  of scope), for general or lossy code summarization / docstring generation,
  for architecture diagrams, or for test generation.
---

# Business Logic Extractor

## Purpose

Recover what the code actually decides — every rule, branch, and edge case —
into a specification a human or downstream agent can audit, with each rule
traced to `file:line` (and, where available, to a requirement and an
execution-trace span). Compression is bounded and accounted for: the
deliverable includes a ledger of exactly what was preserved verbatim, what was
abstracted, and what is known to be missing. This is reverse extraction, never
forward generation.

## When to use this skill

- Use when: "extract the business rules from this module / service" / "what business logic does this code implement".
- Use when: "reverse-engineer the spec / requirements from this codebase" / "produce a traceable rules spec from this implementation".
- Use when: "document the decision logic without dropping edge cases".
- Do NOT use when: the ask is to generate new code or requirements (forward direction), to write a general/lossy summary or docstrings, to draw architecture diagrams, or to generate tests. Hand those to the appropriate skill.

## Inputs

Read access to the target code (module/service/repo scope). Optionally: a
requirements/PRD/issue source and agent or execution traces (OTel GenAI /
OpenInference spans, logs) for cross-referencing. None of these are generated
here — only consumed.

## Workflow

1. **Scope and inventory.** Fix the extraction scope (named modules/services). Enumerate the decision sites: conditionals, validations, state transitions, calculations, guards, error/edge handling. This inventory is the coverage denominator — record it.
2. **Salience-first pass (before any prose).** Identify the rule/entity set: the conditions, thresholds, named entities, and decision paths that carry business meaning. Carry this set forward as an explicit retention constraint — do not let an abstractive pass silently drop it. See `references/loss-bounding-method.md`.
3. **Extract rules with provenance.** For each rule, write one testable statement (EARS-style: "When TRIGGER, the SYSTEM shall RESPONSE" / Given–When–Then) bound to exact `file:line`. Capture every branch and edge case as its own rule or a decision-table row. Quote the load-bearing code span; never infer an API or rule the code does not contain. See `references/traceability-method.md`.
4. **Cross-reference.** Link each rule to a requirement id (if a requirements source exists) and to execution-trace evidence that the path actually runs (distinguish coded-but-dead from exercised logic). Cite `trace_id`/`span_id` alongside `file:line`. AI proposes the link; mark confidence; a human confirms. See `references/trace-cross-referencing.md`.
5. **Build the loss ledger.** Record what was preserved verbatim vs. abstracted, coverage against the step-1 inventory (which decision sites have a rule, which do not), enumerated omissions/uncertainties, and a confidence band. Gate on coverage/omission, not faithfulness alone.
6. **Self-check, then emit.** Verify every rule resolves to real `file:line` (run `scripts/check_provenance.py <spec-file>`), the ledger is present and non-empty, and no rule asserts behavior absent from the code. Emit the two artifacts. Recommend a human spot-check of the flagged gaps; do not claim completeness the ledger contradicts.

Do not recursively re-summarize an already-extracted spec (compounding loss);
re-extract from source instead. Treat the spec as durable externalized state,
not a transient summary.

## Output contract

Two markdown artifacts (default `business-logic-spec.md` and
`loss-ledger.md`, shaped by the `templates/`):

- **Business-logic spec** — enumerated rules; each with: id, EARS-style statement, code provenance `path:line(-line)`, edge cases, a decision table where logic branches, requirement xref (or `none`), execution-trace xref (or `none`), and confidence.
- **Information-loss / coverage ledger** — preserved-verbatim vs. abstracted inventory; coverage of step-1 decision sites (covered / partial / uncovered); enumerated omissions and uncertainties; overall confidence band; human-spot-check checklist.

No code, requirements, tests, or diagrams are produced. The spec is descriptive of existing behavior only.

## Constraints

- DO NOT generate code, requirements, tests, or architecture diagrams — extraction only, code → spec.
- DO NOT invent a rule, API, or identifier the code does not contain (intrinsic/extrinsic hallucination); every rule cites real `file:line`.
- DO NOT drop edge cases or rare branches to make the summary shorter — they are the first casualties and the highest-value content.
- DO NOT recursively re-summarize the spec; re-extract from source to avoid compounding loss.
- DO NOT report completeness the loss ledger does not support; gate on omission, not just faithfulness.
- DO NOT duplicate the methodology here; it lives one level deep in `references/`.

## Validation

- [ ] Frontmatter `name` equals folder, kebab-case, no XML, description under 1024 chars with triggers + negatives.
- [ ] Workflow is salience-first, provenance-bound, and ends with a loss ledger + provenance self-check.
- [ ] Output contract names both artifacts and the no-generation boundary.

## References

- Loss-bounding / anti-lossy-compression extraction method: `references/loss-bounding-method.md`
- Requirement↔code traceability (reverse) technique: `references/traceability-method.md`
- Agent/execution-trace cross-referencing format: `references/trace-cross-referencing.md`
- Failure modes & anti-extrapolation: `references/extraction-pitfalls.md`
- Skeletons: `templates/business-logic-spec.md`, `templates/loss-ledger.md`; provenance check: `scripts/check_provenance.py`
