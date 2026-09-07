# Edge cases and fidelity notes

## Why every value is text

Two damage classes appear when a spreadsheet application type-coerces CSV data, and both are irreversible once saved:

- Long digit IDs collapse to scientific notation: `221673004455667788` becomes `2.21673E+14`, destroying the trailing significant digits.
- Timestamps lose their seconds (`2026-08-01 17:12:34` becomes `1/8/2026 17:12`) and leading zeros vanish (`0042` becomes `42`).

The converter therefore writes every cell as a string with number format `@` (Text). It never prepends an apostrophe to force text - that bakes a literal `'` into the value and corrupts downstream re-exports.

## Pre-damaged input detection

If the source CSV was itself produced by an earlier spreadsheet round-trip, the damage is already in the text. The converter scans cells up to 32 characters long for the pattern `\d\.\d+E\+\d+` and emits a non-fatal `WARN sci_notation` naming the first hit as RxCy (spreadsheet coordinates: header = row 1; record numbers, not physical line numbers). The conversion stays faithful to the CSV as-is; recovering the digits requires re-exporting from the source system. Base64 payload columns cannot false-positive - the base64 alphabet contains no dot.

## Formula neutralization (CSV injection)

A CSV cell like `=HYPERLINK(...)` or `=CMD|...` must never become a live formula in the workbook. Cells whose value starts with `=` are forced to literal text at write time, and the verifier independently asserts the workbook contains zero formula elements.

## Encodings

- Default `utf-8-sig`: reads plain UTF-8 and silently strips a byte-order mark when present.
- Legacy Thai exports: `--encoding cp874` or `--encoding tis-620`.
- A decode error stops the run with `ERROR cannot decode input ...` before any output file is created.

## Delimiters and quoting

- Default comma. Pass `--delimiter ';'` for semicolon or `--delimiter '\t'` for tab (the two-character backslash-t escape is accepted).
- Quoted fields containing delimiters, quotes, or line breaks are handled by the CSV parser; a multi-line field still counts as one record, so record numbers can differ from physical line numbers.

## Ragged rows

- Short rows (fewer fields than the header) are padded with empty cells; common benign causes are omitted trailing commas and an unterminated final line. Reported once as `WARN short_rows=N padded to Y columns`.
- Wide rows (more fields than the header) abort with exit 1: they mean the parse itself is misaligned (wrong delimiter or broken quoting), and continuing would write columns that lie about the data. Silent truncation is never acceptable.
- A completely blank line mid-file becomes one empty, fully padded row.

## Large files and limits

Two-pass design: pass 1 validates and measures without creating any file, so a bad input never leaves a half-written workbook; pass 2 streams rows through the write-only path using inline strings, so no shared-string table is built and memory stays proportional to one row. The CSV field size limit is raised to the platform maximum up front. Hard XLSX-format limits are checked before writing: 1,048,575 data rows and 16,384 columns (the header occupies one row of the 1,048,576 total).

## Sheet naming

The default sheet name is the input file stem with the characters `[]:*?/\` removed, trimmed to 31 characters (the XLSX limit), falling back to `Sheet1` if nothing remains. The same rule is applied to `--sheet` values, and the verifier derives its expected name with the identical rule - pass the same `--sheet` to both scripts.

## Guards

- Lock file `~$NAME.xlsx` beside the output: always refuse, even with `--force` - the workbook is open in a spreadsheet application and writing would conflict.
- Existing output file: refuse without `--force`.
- Output path equal to the input path, or an output directory that does not exist: refuse.

## Reading foreign workbooks

When inspecting a workbook a human has edited, never infer the table width from the sheet dimensions or the maximum column - stray scratch cells outside the table inflate both. The verifier counts rows and columns by iteration and takes the expected shape from the source CSV, and it also rejects saved filter state (`filterMode`, hidden rows) that would make the file open with most rows invisible.
