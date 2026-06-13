# Output Formats

## On-screen results table (per milestone)

```
| # | Name | Role | Pages | Result          | Reason                 |
|---|------|------|-------|-----------------|------------------------|
| 1 | name | role | x–y   | PASS / REJECT   | — or the reason        |
```

For every REJECT, list below the table: the condition violated + the date(s) to fix.

## Verdict file: ReviewResult-JO<JO-number>.txt

Write only after human confirmation. Example name: `ReviewResult-JO13.txt`. Same folder as the
timesheets. Block format:

```
======M<milestone>======
> p<page> <first name>
- <Date> <short reason>

> p<page> <first name>
Pass
```

Rules:

- A milestone where everyone passes → write just `All pass`.
- Use the **first name only**; keep each reason to half a line at most.
- Duplicate (condition 4) → `Duplicate with <name> p<page>` and include the original date in
  parentheses.
