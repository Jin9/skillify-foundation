# Skill Audit Walkthrough

This walkthrough demonstrates how to audit an existing skill using the validation rubric. We'll audit a hypothetical `csv-analyzer` skill.

---

## The Skill Under Audit

```markdown
---
name: CSV Analyzer
description: Analyzes CSV files and produces reports.
---

# CSV Analyzer

This skill helps users analyze CSV data. It can do statistical analysis,
create visualizations, and generate summary reports.

When the user has a CSV file, use this skill to help them understand
their data better. The skill supports various types of analysis including
mean, median, mode, standard deviation, and correlation analysis.

## How to Use

- Load the CSV file
- Look at the columns
- Do the analysis
- Show the results
- If the user wants charts, make charts too

## Notes

- Always use pandas for CSV processing
- Make sure to handle missing values
- Use matplotlib for visualizations
```

---

## Audit Scoring

### 1. Trigger Quality: 2/5 ❌

**Issue**: Description says "Analyzes CSV files" but includes no trigger phrases.
**Fix**:
```yaml
description: >
  Statistical analysis and visualization of CSV data files. Use when the user
  uploads a .csv file, asks to "analyze this data", "show me statistics",
  "create a chart from CSV", or "summarize this spreadsheet".
```

### 2. Scope Focus: 4/5 ✅

**Assessment**: Good — focused on CSV analysis. Could be tighter by separating visualization into its own skill if it grows.

### 3. Workflow Clarity: 1/5 ❌

**Issue**: Steps are vague bullets with no entry or exit conditions, and the load-profile-compute-report order is a real data dependency that the prose leaves implicit. (Numbering is the right verdict here because each step consumes the previous step's output; for judgment work the fix would be goal, constraints, and verification instead.)
**Fix**:
```markdown
## Core workflow
1. **Load data**: Read the CSV file using pandas. Detect encoding and delimiter automatically.
2. **Profile columns**: Identify column types (numeric, categorical, datetime). Report missing value counts.
3. **Compute statistics**: Calculate mean, median, mode, std dev for numeric columns. Frequency counts for categorical.
4. **Generate report**: Output a structured summary with one section per column.
5. **Create visualizations** (if requested): Generate matplotlib charts appropriate to column types.
```

### 4. Output Contract: 1/5 ❌

**Issue**: No specification of what the agent should produce.
**Fix**:
```markdown
## Output format
Produce a markdown summary report with:
- Data overview (rows, columns, types)
- Per-column statistics table
- Missing value report
- Correlation matrix (if >2 numeric columns)
- Charts saved as PNG files (if requested)
```

### 5. Token Efficiency: 4/5 ✅

**Assessment**: Skill is short. But the body contains explanatory text ("This skill helps users...") that wastes tokens. Remove it.

### 6. Conflict Risk: 3/5 ⚠️

**Issue**: Hardcodes "Always use pandas" — what if the repo uses polars?
**Fix**: Change to "Use pandas by default. If the project already uses polars or another library, follow existing patterns."

### 7. Reusability: 4/5 ✅

**Assessment**: Works for any project with CSV files. Minor deduction for the pandas hardcoding.

### 8. Security & Safety: 5/5 ✅

**Assessment**: No destructive commands, no external requests, no bias.

### 9. Frontmatter Correctness: 2/5 ❌

**Issues**:
- Name uses spaces and capitals: `CSV Analyzer` → should be `csv-analyzer`
- Description is too vague
- Missing YAML multi-line syntax

**Fix**:
```yaml
---
name: csv-analyzer
description: >
  Statistical analysis and visualization of CSV data files. Use when the user
  uploads a .csv file, asks to "analyze this data", "show me statistics",
  "create a chart from CSV", or "summarize this spreadsheet".
---
```

### 10. Progressive Disclosure: 3/5 ⚠️

**Assessment**: Skill is short enough that it doesn't need splitting yet. But if visualization instructions grow, they should move to `references/visualization-guide.md`.

---

## Final Score

| Dimension                | Score | Status |
|--------------------------|-------|--------|
| 1. Trigger Quality       | 2/5   | ❌ Fix |
| 2. Scope Focus           | 4/5   | ✅     |
| 3. Workflow Clarity       | 1/5   | ❌ Fix |
| 4. Output Contract        | 1/5   | ❌ Fix |
| 5. Token Efficiency       | 4/5   | ✅     |
| 6. Conflict Risk          | 3/5   | ⚠️ Fix |
| 7. Reusability            | 4/5   | ✅     |
| 8. Security & Safety      | 5/5   | ✅     |
| 9. Frontmatter Correctness| 2/5   | ❌ Fix |
| 10. Progressive Disclosure | 3/5   | ⚠️     |
| **TOTAL**                 | **29/50** | ❌ Fail (needs 40) |

---

## Refactored Skill

After applying all fixes:

```markdown
---
name: csv-analyzer
description: >
  Statistical analysis and visualization of CSV data files. Use when the user
  uploads a .csv file, asks to "analyze this data", "show me statistics",
  "create a chart from CSV", or "summarize this spreadsheet".
  Do NOT use for database queries (use sql-analyzer instead).
---

# CSV Analyzer

## Purpose

Analyze CSV data files to produce statistical summaries and visualizations.

## Core workflow

1. **Load data**: Read the CSV file. Detect encoding and delimiter automatically.
2. **Profile columns**: Identify column types (numeric, categorical, datetime). Report missing value counts.
3. **Compute statistics**: Calculate mean, median, mode, std dev for numeric columns. Frequency counts for categorical.
4. **Generate report**: Output a structured markdown summary with one section per column.
5. **Create visualizations** (if requested): Generate charts appropriate to column types. Save as PNG.

## Output format

Produce a markdown summary report containing:
- Data overview table (rows, columns, types)
- Per-column statistics
- Missing value report
- Correlation matrix (if >2 numeric columns)
- Chart files as PNG (if requested)

## Operating contract

- Instruction priority: the user's request takes precedence over this skill; if a line here would make you pause or diverge from it, follow the user and say which line you set aside.
- Autonomy: "analyze this data" is an instruction; act on the file named. Ask at most one question, only if the target file cannot be found.
- Stop conditions: stop only before overwriting an existing report or chart file, or when the request changes scope. Before ending, check that the last paragraph is not a plan.
- Verification: quote the row and column counts read from the file and the missing-value report before claiming the analysis is complete.
- Delegation: none.
- Progress: open with one line naming the file and the steps; close with a recap of the report sections produced.
- Model-cost tier: small for load and profile; mid for the written summary.

## Constraints

- Use pandas by default and follow the project's existing library if it differs; hardcoding a library breaks portability.
- Handle missing values before computing statistics; unhandled gaps skew every summary number.
- Report data quality issues before the analysis, so the reader can weigh the results.

## Examples

### Example: Basic analysis

**User says**: "Analyze this sales data CSV"

**Action**:
1. Load `sales_data.csv`, detect 5 columns, 1200 rows
2. Profile: 3 numeric, 1 categorical, 1 datetime
3. Compute statistics per column
4. Generate markdown report

**Result**: Structured report with statistics table and data quality summary.
```

**Rescore**: 44/50 ✅ Pass (`scripts/cruft_scan.py`: high=0 medium=0 low=0)
