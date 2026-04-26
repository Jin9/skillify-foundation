# Skill Audit Report Template

Use this template for Audit mode and as a concise structure for Review mode findings.

## Summary

- Skill: `[name]`
- Folder: `[path]`
- Date: `[YYYY-MM-DD]`
- Mode: `Review` or `Audit`
- Result: `Pass`, `Pass with fixes`, or `Fail`

## Deterministic Checks

```text
quick_validate.py: [pass/fail/not run]
check_links.py: [pass/fail/not run]
```

## Rubric Scores

| Dimension | Score | Notes |
|-----------|-------|-------|
| 1. Trigger Quality | /5 | |
| 2. Scope Focus | /5 | |
| 3. Workflow Clarity | /5 | |
| 4. Output Contract | /5 | |
| 5. Token Efficiency | /5 | |
| 6. Conflict Risk | /5 | |
| 7. Reusability | /5 | |
| 8. Security & Safety | /5 | |
| 9. Frontmatter Correctness | /5 | |
| 10. Progressive Disclosure | /5 | |
| **Total** | /50 | Pass threshold: 40/50 and all dimensions at least 4/5 |

## Findings

| Severity | File | Finding | Fix |
|----------|------|---------|-----|
| High/Medium/Low | `[path]` | `[issue]` | `[recommended change]` |

## Anti-Pattern Sweep

| Anti-pattern | Status | Notes |
|--------------|--------|-------|
| Do Everything Skill | Pass/Fail/Mitigated | |
| Weak or Vague Triggers | Pass/Fail/Mitigated | |
| Context Window Bloat | Pass/Fail/Mitigated | |
| Mixing Repo Policy | Pass/Fail/Mitigated | |
| Non-Step Workflows | Pass/Fail/Mitigated | |
| Overriding User Intent | Pass/Fail/Mitigated | |
| Generating Target Output | Pass/Fail/Mitigated | |
| Human Docs in Skill Folder | Pass/Fail/Mitigated | |
| Trigger Phrase Absence | Pass/Fail/Mitigated | |
| Stale Skill | Pass/Fail/Mitigated | |
| Duplicated Cross-Tier Content | Pass/Fail/Mitigated | |
| Hardcoded Platform Assumptions | Pass/Fail/Mitigated | |

## Security Sweep

- External HTTP calls: `[none/found]`
- Secret or environment access: `[none/found]`
- Destructive commands: `[none/found]`
- Broad permissions: `[none/found]`
- Vendor bias: `[none/found]`

## Delta List

1. `[required or recommended change]`
2. `[required or recommended change]`

## Final Recommendation

`[Ship / Refactor before ship / Split / Merge rejected / Retire]`
