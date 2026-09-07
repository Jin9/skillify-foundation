#!/usr/bin/env python3
"""Independently verify an XLSX produced from a CSV by csv_to_xlsx.py.

Re-reads both files (no code shared with the converter) and asserts:
  1 sheet        exactly one worksheet, expected name
  2 counts       row count by iteration == CSV records (never trusts dimensions)
  3 header       header row equality vs the CSV
  4 spot_rows    first / middle / last data rows: full string equality
  5 autofilter   autoFilter ref spans header + data exactly
  6 clean_state  no saved filter state, no hidden rows, no live formulas
  7 freeze_style frozen pane + header style present iff --styled, absent otherwise
  8 no_bloat     no sharedStrings part (inline strings only)

stdout: one final line, either
  PASS file=... rows=N cols=N filter=... mode=... checks=8/8    (exit 0)
  FAIL check=NAME expected=... actual=...                        (exit 1)
"""

import argparse
import csv
import sys
import xml.etree.ElementTree as ET
import zipfile
from pathlib import Path

try:
    import openpyxl
    from openpyxl.utils import get_column_letter
except ImportError:
    sys.stderr.write(
        "ERROR openpyxl is required. Install: python3 -m pip install --user openpyxl\n"
    )
    sys.exit(1)

NS = "{http://schemas.openxmlformats.org/spreadsheetml/2006/main}"
SHEET_FORBIDDEN = "[]:*?/\\"
SHAREDSTRINGS_MAX_BYTES = 10240


def fail(check, expected, actual):
    print("FAIL check=%s expected=%s actual=%s" % (check, expected, actual))
    sys.exit(1)


def raise_field_limit():
    try:
        csv.field_size_limit(sys.maxsize)
    except OverflowError:
        csv.field_size_limit(2 ** 31 - 1)


def sanitize_sheet_name(raw):
    name = "".join(ch for ch in raw if ch not in SHEET_FORBIDDEN)
    name = name.strip()[:31].strip()
    return name or "Sheet1"


def norm(values, width):
    out = ["" if v is None else str(v) for v in values]
    if len(out) < width:
        out = out + [""] * (width - len(out))
    return out[:width]


def csv_facts(path, delimiter, encoding):
    """One streaming pass: header, record count, first and last data rows."""
    header = None
    n_rows = 0
    first_row = None
    last_row = None
    with open(str(path), "r", newline="", encoding=encoding) as fh:
        for row in csv.reader(fh, delimiter=delimiter):
            if header is None:
                header = row
                continue
            n_rows += 1
            if first_row is None:
                first_row = row
            last_row = row
    return header, n_rows, first_row, last_row


def csv_row_at(path, delimiter, encoding, data_index):
    """Fetch 1-based data row `data_index` in a second streaming pass."""
    with open(str(path), "r", newline="", encoding=encoding) as fh:
        seen_header = False
        i = 0
        for row in csv.reader(fh, delimiter=delimiter):
            if not seen_header:
                seen_header = True
                continue
            i += 1
            if i == data_index:
                return row
    return None


def sheet_xml_facts(xlsx_path):
    """Structural facts straight from the zip parts (read-only mode hides these)."""
    facts = {
        "sheet_members": 0,
        "autofilter_ref": None,
        "filter_mode": None,
        "hidden_rows": 0,
        "formulas": 0,
        "pane": None,
        "row1_ht": None,
        "sharedstrings_size": None,
        "styles_xml": "",
    }
    with zipfile.ZipFile(str(xlsx_path)) as zf:
        names = zf.namelist()
        sheets = [
            n for n in names
            if n.startswith("xl/worksheets/sheet") and n.endswith(".xml")
        ]
        facts["sheet_members"] = len(sheets)
        if "xl/sharedStrings.xml" in names:
            facts["sharedstrings_size"] = zf.getinfo("xl/sharedStrings.xml").file_size
        if "xl/styles.xml" in names:
            with zf.open("xl/styles.xml") as fh:
                facts["styles_xml"] = fh.read().decode("utf-8", "replace")
        if len(sheets) != 1:
            return facts
        with zf.open(sheets[0]) as fh:
            for _, elem in ET.iterparse(fh):
                tag = elem.tag
                if tag == NS + "autoFilter":
                    facts["autofilter_ref"] = elem.get("ref")
                elif tag == NS + "sheetPr":
                    facts["filter_mode"] = elem.get("filterMode")
                elif tag == NS + "pane":
                    facts["pane"] = dict(elem.attrib)
                elif tag == NS + "f":
                    facts["formulas"] += 1
                elif tag == NS + "row":
                    if elem.get("hidden") in ("1", "true"):
                        facts["hidden_rows"] += 1
                    if elem.get("r") == "1":
                        facts["row1_ht"] = elem.get("ht")
                    elem.clear()  # keep memory flat on large sheets
    return facts


def main():
    parser = argparse.ArgumentParser(
        description="Verify an XLSX against its source CSV: content equality,"
        " AutoFilter range, clean filter state, mode-dependent styling."
    )
    parser.add_argument("xlsx", help="workbook to verify")
    parser.add_argument("csv_input", metavar="csv", help="source CSV file")
    parser.add_argument("--styled", action="store_true", help="expect the styled header/freeze")
    parser.add_argument("--sheet", help="expected worksheet name (default: derived from the CSV name)")
    parser.add_argument("--delimiter", default=",", help="CSV delimiter (default ','; use \\t for tab)")
    parser.add_argument("--encoding", default="utf-8-sig", help="CSV encoding (default utf-8-sig)")
    args = parser.parse_args()

    raise_field_limit()
    delimiter = "\t" if args.delimiter == "\\t" else args.delimiter
    mode = "styled" if args.styled else "minimal"

    xlsx_path = Path(args.xlsx).resolve()
    csv_path = Path(args.csv_input).resolve()
    if not xlsx_path.is_file():
        fail("input", "existing xlsx", "missing %s" % xlsx_path)
    if not csv_path.is_file():
        fail("input", "existing csv", "missing %s" % csv_path)

    try:
        header, n_rows, first_row, last_row = csv_facts(csv_path, delimiter, args.encoding)
    except (UnicodeDecodeError, csv.Error) as exc:
        fail("input", "readable CSV", repr(exc))
    if header is None:
        fail("input", "non-empty CSV", "empty file")
    n_cols = len(header)
    expected_sheet = sanitize_sheet_name(args.sheet if args.sheet else csv_path.stem)
    expected_ref = "A1:%s%d" % (get_column_letter(n_cols), n_rows + 1)

    # Checks 1-4: values, via openpyxl read-only streaming
    try:
        wb = openpyxl.load_workbook(str(xlsx_path), read_only=True)
    except Exception as exc:  # zip/format corruption surfaces here
        fail("input", "loadable xlsx", repr(exc))
    try:
        if wb.sheetnames != [expected_sheet]:
            fail("sheet", "[%s]" % expected_sheet, str(wb.sheetnames))
        ws = wb[expected_sheet]
        mid_index = (n_rows + 1) // 2 if n_rows else 0  # 1-based data index
        targets = {1: None, mid_index: None, n_rows: None} if n_rows else {}
        xl_count = 0
        xl_header = None
        for row in ws.rows:
            xl_count += 1
            if xl_count == 1:
                xl_header = norm([c.value for c in row], n_cols)
                continue
            data_idx = xl_count - 1
            if data_idx in targets:
                targets[data_idx] = norm([c.value for c in row], n_cols)
        if xl_count != n_rows + 1:
            fail("counts", "%d rows (header+%d)" % (n_rows + 1, n_rows), str(xl_count))
        if xl_header != norm(header, n_cols):
            fail("header", str(norm(header, n_cols))[:200], str(xl_header)[:200])
        if n_rows:
            expected_rows = {
                1: first_row,
                n_rows: last_row,
                mid_index: csv_row_at(csv_path, delimiter, args.encoding, mid_index),
            }
            for idx in sorted(targets):
                if targets[idx] != norm(expected_rows[idx], n_cols):
                    fail("spot_rows", "data row %d equal to CSV" % idx, "mismatch at data row %d" % idx)
    finally:
        wb.close()

    # Checks 5-8: structure, straight from the XML parts
    facts = sheet_xml_facts(xlsx_path)
    if facts["sheet_members"] != 1:
        fail("clean_state", "1 worksheet part", str(facts["sheet_members"]))
    if facts["autofilter_ref"] != expected_ref:
        fail("autofilter", expected_ref, str(facts["autofilter_ref"]))
    if facts["filter_mode"] in ("1", "true"):
        fail("clean_state", "no filterMode", "filterMode=%s" % facts["filter_mode"])
    if facts["hidden_rows"]:
        fail("clean_state", "0 hidden rows", "%d hidden" % facts["hidden_rows"])
    if facts["formulas"]:
        fail("clean_state", "0 formula cells", "%d formulas" % facts["formulas"])

    pane = facts["pane"]
    if args.styled:
        pane_ok = (
            pane is not None
            and pane.get("ySplit") == "1"
            and pane.get("topLeftCell") == "A2"
            and pane.get("state") in ("frozen", "frozenSplit")
        )
        if not pane_ok:
            fail("freeze_style", "frozen pane ySplit=1 topLeftCell=A2", str(pane))
        if facts["row1_ht"] is None:
            fail("freeze_style", "header row height set", "none")
        if "1F4E78" not in facts["styles_xml"]:
            fail("freeze_style", "header fill 1F4E78 in styles", "absent")
    else:
        if pane is not None:
            fail("freeze_style", "no frozen pane in minimal mode", str(pane))
        if "1F4E78" in facts["styles_xml"]:
            fail("freeze_style", "no header fill in minimal mode", "1F4E78 present")

    ss_size = facts["sharedstrings_size"]
    if ss_size is not None and ss_size > SHAREDSTRINGS_MAX_BYTES:
        fail("no_bloat", "no sharedStrings part (or under 10 KB)", "%d bytes" % ss_size)

    print(
        "PASS file=%s rows=%d cols=%d filter=%s mode=%s checks=8/8"
        % (xlsx_path, n_rows, n_cols, expected_ref, mode)
    )


if __name__ == "__main__":
    main()
