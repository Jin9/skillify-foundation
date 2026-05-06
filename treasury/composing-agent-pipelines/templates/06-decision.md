# 06 — Decision

**Decided**: [ISO timestamp]
**Inputs**: `03-analysis.md`, `04-review.md` (if exists), `05-validation.md` (if exists)

## Question

What decision does this artifact resolve? State in one sentence.

## Options

At least 2. Each row in the table is an option.

| # | Option | Pros | Cons | Risks | Effort | Reversibility |
|---|--------|------|------|-------|--------|---------------|
| 1 | [name] | [bulleted] | [bulleted] | [bulleted] | [S/M/L] | [easy/hard] |
| 2 | [name] | [bulleted] | [bulleted] | [bulleted] | [S/M/L] | [easy/hard] |
| 3 | [name] | [bulleted] | [bulleted] | [bulleted] | [S/M/L] | [easy/hard] |

## Recommendation

**Choose option [N]: [name]**.

**Rationale**: 2–4 paragraphs grounded in the analysis. Cite specific findings (e.g. F2 from `03-analysis.md`). If review or validation exists, address the P1 issues that influenced the choice.

## Why not the alternatives

For each rejected option, one sentence on the deciding factor.

- Option 1: rejected because [reason].
- Option 3: rejected because [reason].

## Residual risks

What could still go wrong even if we pick the recommendation?

- [risk 1] — mitigation: [proposal]
- [risk 2] — mitigation: [proposal]

## Trigger conditions to revisit

When should this decision be reopened?

- If [condition], revisit.
- If [metric] crosses [threshold], revisit.

## Sources

Findings, claims, and validation results that this decision rests on.

- `03-analysis.md#F2`
- `05-validation.md#C1` (verdict: confirmed)
