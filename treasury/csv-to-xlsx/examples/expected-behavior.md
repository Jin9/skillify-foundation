# Expected behavior on the fixture

Run from the skill root; write outputs to a scratch directory outside the skill folder. Use these exact lines as the oracle after any script change.

## Minimal mode (default)

    python3 scripts/csv_to_xlsx.py examples/sample-input.csv -o /tmp/sample-minimal.xlsx
    python3 scripts/verify_xlsx.py /tmp/sample-minimal.xlsx examples/sample-input.csv

Expected converter stderr:

    WARN sci_notation cells=1 first=R8C7 sample=2.21673E+14 (input already Excel-damaged; conversion stays faithful to the CSV as-is)

Expected converter stdout (absolute paths vary by machine):

    OK input=.../examples/sample-input.csv output=/tmp/sample-minimal.xlsx sheet=sample-input rows=7 cols=7 filter=A1:G8 mode=minimal warnings=1

Expected verifier stdout:

    PASS file=/tmp/sample-minimal.xlsx rows=7 cols=7 filter=A1:G8 mode=minimal checks=8/8

## Styled mode

    python3 scripts/csv_to_xlsx.py examples/sample-input.csv -o /tmp/sample-styled.xlsx --styled
    python3 scripts/verify_xlsx.py /tmp/sample-styled.xlsx examples/sample-input.csv --styled

Same `WARN` line; the `OK` line ends `mode=styled warnings=1`; the `PASS` line ends `mode=styled checks=8/8`.

## What the fixture exercises

Row numbers are spreadsheet rows (header = row 1).

| Row | Case | Guarantee proven |
|---|---|---|
| 2 | 18-digit `citizen_ref`, Thai text, `1500.50` | no scientific notation, UTF-8 intact, trailing zero kept |
| 3 | quoted `"Somsri, Jr."`, `0042` | embedded comma handled, leading zeros kept |
| 4 | `=SUM(A1:A2)` in `note` | formula-looking value neutralized to literal text |
| 5 | `007` | leading zeros kept |
| 6 | quoted multi-line `full_name` | one record despite the line break |
| 7 | empty `amount_thb` | empty cell round-trips |
| 8 | `2.21673E+14` in `legacy_ref` | pre-damage detected as `WARN sci_notation` at R8C7 |

## Guard behaviors (exact messages)

| Trigger | Result |
|---|---|
| Output exists, no `--force` | exit 1, `ERROR output exists: ... (pass --force to overwrite)` |
| `~$NAME.xlsx` lock beside output | exit 1, `ERROR Excel holds a lock on the output (...); close the workbook first` |
| Row wider than header | exit 1, `ERROR row N has X fields, expected Y (delimiter or quoting mismatch?)`, no file written |
| Short rows | non-fatal `WARN short_rows=N padded to Y columns`, conversion proceeds |
