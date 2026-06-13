#!/usr/bin/env python3
"""Extract the visible text from a Job Order .docx (a zipped OOXML file).

A .docx is a zip archive; the body text lives in word/document.xml as <w:t> runs.
This prints the concatenated run text so the agent can read off Estimated Manday,
Estimated Total Amount Incl. VAT, and the milestone breakdown. Read-only.

Usage:
    python3 extract_jo_docx.py "Job Order JO13.docx"
"""

from __future__ import annotations

import sys
import xml.etree.ElementTree as ET
import zipfile

WORD_NS = "{http://schemas.openxmlformats.org/wordprocessingml/2006/main}"


def extract(path: str) -> str:
    with zipfile.ZipFile(path) as archive:
        xml = archive.read("word/document.xml")
    tree = ET.fromstring(xml)
    return " ".join(node.text for node in tree.iter(f"{WORD_NS}t") if node.text)


def main() -> int:
    if len(sys.argv) != 2:
        print("usage: python3 extract_jo_docx.py <path-to-job-order.docx>", file=sys.stderr)
        return 2
    try:
        print(extract(sys.argv[1]))
    except (FileNotFoundError, zipfile.BadZipFile, KeyError) as err:
        print(f"error: could not read .docx text: {err}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
