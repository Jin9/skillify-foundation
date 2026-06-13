---
name: milestone-jo-check
description: >
  Verify that a vendor's billed milestones stay within the Job Order (JO) limits — that the
  summed Total Manday across all milestones is within the JO Estimated Manday, and the summed
  Total Amount is within the JO Estimated Total Amount Incl. VAT. Use when the user asks to
  "check the milestones against the job order", "verify total manday vs the JO limit", "is the
  vendor billing within the JO estimate", "check JO conditions", or "compare summary manday to
  the JO". Read-only: it extracts the JO .docx limits and reads each Summary_Manday PDF, sums
  them, and prints two comparison tables. Do NOT use to review individual people's daily
  timesheet entries (that is timesheet-individual-review) or to derive replacement tasks from
  Jira (rejected-task-replacement). Never modifies or deletes any source file.
---

# Milestone vs JO Conditions Check

Confirm the billed milestones do not exceed the Job Order's estimated manday and amount limits.
This is a read-only analysis — it writes no files and changes nothing on disk.

## Expected folder layout

```
<JO Folder>/
├── Job Order <xxx>.docx
├── <folder M1>/   (suffix has M + number, e.g. [JO17-M1])
│   └── Summary_Manday_<xxx>.pdf
├── <folder M2>/
│   └── ...
└── <folder MN>/
```

## Steps

1. Ask the user for the JO folder path if it was not given.
2. `ls` the JO folder and every milestone subfolder to locate the `Job Order *.docx` and each
   `Summary_Manday_*.pdf`. Identify a milestone folder by its `M<number>` suffix.
3. Extract the JO limits from the `.docx` with the helper script:
   ```bash
   python3 scripts/extract_jo_docx.py "<path to Job Order .docx>"
   ```
   It prints the document text; read off **Estimated Manday** and **Estimated Total Amount Incl.
   VAT** (and the milestone breakdown if present).
4. Read every `Summary_Manday_*.pdf` (read them in parallel in one call). From each, pull the
   milestone's **Total Manday** and **Total Amount**.
5. Sum across milestones and compare:
   - 4.1 Manday: `Sum(Total Manday) <= Estimated Manday`
   - 4.2 Amount: `Sum(Total Amount) <= Estimated Total Amount Incl. VAT`

## Output (printed only — no file written)

Per-milestone breakdown, then the verdict table:

```
| Milestone | Period      | Total Manday | Total Amount (THB)  |
|-----------|-------------|--------------|---------------------|
| M1        | month–month | xx.x         | x,xxx,xxx.xx        |
| Sum       |             | xxx.x        | x,xxx,xxx.xx        |

| Condition          | Sum          | JO Limit      | Result          |
|--------------------|--------------|---------------|-----------------|
| 4.1 Total Manday   | xxx.x        | <= xxx.x      | PASS / OVER     |
| 4.2 Total Amount   | x,xxx,xxx.xx | <= x,xxx,xxx  | PASS / OVER     |
```

## Guardrails

- NEVER modify, move, or delete any source file (JO `.docx`, PDFs, folders).
- This skill writes no files. If the user wants a saved result, hand off to
  timesheet-individual-review, which owns the written verdict file behind a human-approval gate.
- Report the numbers as extracted; never adjust a sum to make a condition pass.
