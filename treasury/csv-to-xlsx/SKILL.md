---
name: csv-to-xlsx
description: Convert a CSV file into an XLSX workbook with an Excel AutoFilter on the header row, preserving every value exactly as text - long digit IDs never collapse into scientific notation, leading zeros and timestamp seconds survive. Default output is data plus the header filter only; a --styled flag adds a bold white-on-navy header, frozen top row, and readable column widths. Use when the user asks to "convert this csv to xlsx", "csv to excel with filter", "add filter at the header", or "export this data as excel". Streams large files, handles UTF-8 with or without BOM, custom delimiters, and quoted fields, and verifies output by independent read-back. Do NOT use for xlsx-to-csv export, rendering tables as HTML (rendering-readable-html), charts (dataviz), or pivot-table data analysis.
---

# CSV to XLSX

Convert a CSV file into an XLSX workbook with filter arrows on the header row and every value preserved as literal text. The fragile logic lives in two deterministic scripts; the agent's job is to run them in order and relay their contract lines. Every step is model-cost tier: small.

## When to use this skill

- The user wants a CSV (or CSV-like export) turned into an Excel/XLSX file.
- The user asks for "the filter at the header" - an AutoFilter over the header row.
- A database or report export must reach Excel without corrupting IDs, leading zeros, or timestamps.

## When NOT to use this skill

- XLSX to CSV (reverse direction), or editing an existing workbook.
- Rendering tables as HTML (rendering-readable-html) or charts (dataviz).
- Data analysis, pivot tables, or summarizing the data itself.

## Requirements

`python3` with the `openpyxl` package. Entry check: `python3 -c "import openpyxl"`. If it fails, install with `python3 -m pip install --user openpyxl` (or `uv pip install openpyxl`), then re-check before proceeding.

## Workflow

1. [small] Resolve the input CSV path. Choose the output path (default: same folder, same name, `.xlsx` extension) and the mode: minimal (default) or styled (`--styled` adds a frozen header row, bold white-on-navy header cells, and sized columns). If the output file already exists, confirm with the user before using `--force`. Exit condition: input path exists, output path and mode decided.
2. [small] Convert:
   `python3 scripts/csv_to_xlsx.py INPUT.csv [-o OUTPUT.xlsx] [--styled] [--sheet NAME] [--delimiter D] [--encoding E] [--force]`
   Relay the final `OK` line and every `WARN` line verbatim. `WARN sci_notation` means the source CSV was already damaged by an earlier Excel round-trip; the conversion stays faithful to the CSV as-is, but the lost digits are unrecoverable - suggest re-exporting from the source system. `WARN short_rows` means ragged short rows were padded with empty cells. Exit condition: exit code 0 and one `OK` line captured.
3. [small] Verify - never skip this step and never hand over an unverified artifact:
   `python3 scripts/verify_xlsx.py OUTPUT.xlsx INPUT.csv [--styled] [--sheet NAME] [--delimiter D] [--encoding E]`
   Use the same mode, sheet, delimiter, and encoding flags as step 2. Exit condition: a `PASS` line. On `FAIL`, report the line to the user and stop.
4. [small] Report to the user: output path, sheet name, rows, columns, filter range, mode, warnings, and the verifier's `PASS` line. Never claim success from the converter's own output alone.

## Output contract

One machine-readable line each, space-separated key=value:

- Converter success (stdout): `OK input=PATH output=PATH sheet=NAME rows=N cols=N filter=A1:CN mode=minimal|styled warnings=N`
- Converter warnings (stderr, non-fatal): `WARN short_rows=N padded to N columns` and `WARN sci_notation cells=N first=RxCy sample=VALUE (...)`
- Converter errors (stderr, fatal): `ERROR message`
- Verifier (stdout): `PASS file=PATH rows=N cols=N filter=REF mode=MODE checks=8/8` or `FAIL check=NAME expected=X actual=Y`

## Constraints and guarantees

- Every cell is written as text (number format `@`); values are never coerced, truncated, or prefixed with an apostrophe.
- Values starting with `=` are neutralized to literal text - the workbook never contains live formulas.
- The AutoFilter range spans header plus data exactly, computed from actual counts, never inferred from sheet dimensions.
- No saved filter state and no hidden rows: every row is visible when the file opens.
- Streaming both ways (write-only and read-only modes): large files convert in O(row) memory with no sharedStrings bloat.
- Rows shorter than the header are padded; rows wider than the header abort the run (a misaligned parse must never produce misaligned columns).
- Refuses to run while Excel holds a lock (`~$NAME.xlsx`) on the output, and refuses to overwrite without `--force`.
- XLSX limits enforced before writing: 1,048,575 data rows, 16,384 columns.

## Failure modes

| Exit | Meaning | Agent action |
|---|---|---|
| 0 | Success (warnings possible) | Relay `OK`/`WARN` lines, then run step 3 |
| 1 | Input, validation, or environment error (`ERROR` line) | Report the line; fix the cause (encoding, delimiter, lock, `--force`) and retry once |
| 2 | Usage error (bad flags) | Fix the command line and retry |
| verifier 1 | `FAIL check=...` | Report; do not deliver the file; re-convert only if the cause is fixable |

## References

- `references/edge-cases-and-fidelity.md` - encodings (legacy Thai included), delimiters, ragged rows, damage classes, limits, sheet naming, guards.
- `examples/sample-input.csv` - synthetic fixture exercising the edge cases.
- `examples/expected-behavior.md` - exact expected `OK`/`WARN`/`PASS` lines for the fixture; use as the oracle after any script change.
