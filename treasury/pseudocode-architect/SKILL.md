---
name: pseudocode-architect
description: Analyzes requirements or existing code to generate clean, language-agnostic pseudocode. Use when planning algorithms, breaking down complex logic, or designing systems before actual implementation. Do NOT use for writing production-ready code in specific languages.
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

Produce language-agnostic pseudocode wrapped in `text` code blocks. Precede the pseudocode with a brief problem framing statement.

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