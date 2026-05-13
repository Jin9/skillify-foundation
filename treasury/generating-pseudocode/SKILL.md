---
name: generating-pseudocode
description: >
  Analyzes requirements or existing code to generate clean, language-agnostic
  pseudocode that bridges high-level intent and implementation, readable by a
  Python, Go, or TypeScript developer without translation. Use when the user
  says "write pseudocode for X", "plan the algorithm for X", "language-agnostic
  logic for X", "break this logic into steps before I code it", "draft the
  flow before implementation", "logic breakdown for X", or "design the
  algorithm first". Also use when stakeholders need to agree on business
  logic before a language is chosen, or when an existing function is too
  dense to review and needs to be re-expressed as steps. Outputs a
  problem-framing block (Inputs, Outputs, Constraints / edge cases), a
  pseudocode block in a fenced text segment, and a short verification note.
  Do NOT use for writing production-ready code in specific languages, for
  translating pseudocode back into a specific language, or for high-level
  system architecture.
---

# Pseudocode Generation

## Purpose

Analyze complex logic, requirements, or existing code to generate clean, language-agnostic pseudocode. This acts as a bridge between high-level intent and actual implementation.

## When to use this skill

- Planning algorithms or business logic flows.
- Breaking down complex logic before actual implementation.
- Designing systems where the logic needs to be agreed upon by stakeholders before coding.
- Trigger terms: "pseudocode", "logic breakdown", "algorithm plan", "language-agnostic".

## Modes

### `draft` - Initial pseudocode generation
Draft high-level logic steps and detailed pseudocode from requirements.

### `verify` - Logic verification against constraints
Review generated pseudocode against edge cases and constraints to ensure robustness.

## Core workflow

1. **Problem Framing**: Briefly state inputs (data structures), outputs (expected result), and constraints/edge cases (empty lists, timeouts, null values).
2. **Drafting**: Write the high-level logic steps.
3. **Generation**: Output structured pseudocode using the reference rules.
4. **Verification**: Review the generated pseudocode against the constraints identified in step 1.

## Output format

Produce three sections in this order:

1. **Problem framing** — three labeled lines:
   - `Inputs:` data types and shapes the logic consumes.
   - `Outputs:` what the logic returns or emits.
   - `Constraints / edge cases:` empty collections, nulls, timeouts, concurrent access, ordering, or any failure mode that must be handled.
2. **Pseudocode** — a fenced code block tagged as `text`, following the keyword and formatting rules in `references/pseudocode-rules.md`. Use capitalized control-flow keywords (FUNCTION, IF/THEN, FOR EACH, TRY/CATCH). No language-specific syntax.
3. **Verification note** — one or two sentences naming the constraints from step 1 that were checked against the pseudocode, and any constraint deliberately deferred.

Do not surround the response with explanatory prose beyond what is needed to label these three sections.

## Constraints

- DO NOT write actual production code (e.g., Python, Go, TypeScript) when generating pseudocode.
- DO NOT use language-specific boilerplate or syntax.
- DO NOT proceed to code implementation until the pseudocode logically satisfies all edge cases.

## Troubleshooting

| Signal | Action |
|--------|--------|
| Unclear requirements | Ask for exact inputs, outputs, and edge cases before writing. |
| Logic is too complex | Break down into smaller helper functions in the pseudocode. |
| Overly specific syntax | Remind the agent to remain language-agnostic using generic keywords. |

## References

- **Pseudocode Rules**: See [references/pseudocode-rules.md](references/pseudocode-rules.md) for core principles and formatting examples.

## Validation gate

Before sending the output, confirm:
1. Pseudocode is language-agnostic.
2. Problem framing (Inputs, Outputs, Constraints) is present.
3. Constraints and edge cases have been logically handled in the pseudocode.