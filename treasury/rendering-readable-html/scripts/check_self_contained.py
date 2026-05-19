#!/usr/bin/env python3
"""Deterministic gate: prove one HTML file is self-contained and JS-free.

Exit 0 only if the file: starts with a doctype, declares a utf-8 charset,
contains no script element, no inline event handler, no javascript: URL, and no
remote (http/https/protocol-relative) resource load. Hyperlinks (<a href>) and
data: URIs are allowed. Python 3 standard library only.
"""

from __future__ import annotations

import argparse
import re
import sys
from html.parser import HTMLParser
from pathlib import Path

# tag -> attributes that LOAD an external resource on render
RESOURCE_ATTRS = {
    "img": ("src", "srcset"),
    "source": ("src", "srcset"),
    "video": ("src", "poster"),
    "audio": ("src",),
    "track": ("src",),
    "iframe": ("src",),
    "embed": ("src",),
    "object": ("data",),
    "input": ("src",),
    "link": ("href",),
    "use": ("href", "xlink:href"),
    "image": ("href", "xlink:href"),
}

URL_FUNC_RE = re.compile(r"url\(\s*['\"]?\s*([^'\")]+)", re.IGNORECASE)
IMPORT_RE = re.compile(r"@import\s+(?:url\(\s*)?['\"]?\s*([^'\");]+)", re.IGNORECASE)


def is_remote(value: str) -> bool:
    v = value.strip().lower()
    return v.startswith(("http://", "https://", "//"))


def is_js_url(value: str) -> bool:
    return value.strip().lower().startswith("javascript:")


def srcset_urls(value: str) -> list[str]:
    out: list[str] = []
    for part in value.split(","):
        token = part.strip().split()
        if token:
            out.append(token[0])
    return out


class Checker(HTMLParser):
    def __init__(self) -> None:
        super().__init__(convert_charrefs=True)
        self.errors: list[str] = []
        self.saw_doctype = False
        self.saw_charset = False
        self._in_style = False

    def _loc(self) -> str:
        line, col = self.getpos()
        return f"line {line}:{col}"

    def handle_decl(self, decl: str) -> None:
        if decl.strip().lower().startswith("doctype"):
            self.saw_doctype = True
            if "html" not in decl.lower():
                self.errors.append(f"{self._loc()}: doctype is not <!doctype html>")

    def handle_starttag(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        self._inspect(tag, attrs)

    def handle_startendtag(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        self._inspect(tag, attrs)

    def _inspect(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        amap = {k.lower(): (v or "") for k, v in attrs}

        if tag == "script":
            self.errors.append(f"{self._loc()}: <script> element is not allowed (zero JavaScript)")
        if tag == "style":
            self._in_style = True

        if tag == "meta":
            if amap.get("charset", "").strip().lower() == "utf-8":
                self.saw_charset = True
            if amap.get("http-equiv", "").strip().lower() == "content-type" and \
                    "charset=utf-8" in amap.get("content", "").lower():
                self.saw_charset = True

        for name, value in amap.items():
            if name.startswith("on"):
                self.errors.append(f"{self._loc()}: inline event handler '{name}' is JavaScript")
            if is_js_url(value):
                self.errors.append(f"{self._loc()}: javascript: URL in '{name}'")

        # <a href> is navigation, not a load — allowed (only the js: check above applies).
        for attr in RESOURCE_ATTRS.get(tag, ()):  # noqa: SIM118
            raw = amap.get(attr, "")
            if not raw:
                continue
            candidates = srcset_urls(raw) if attr == "srcset" else [raw]
            for c in candidates:
                if is_remote(c):
                    self.errors.append(
                        f"{self._loc()}: <{tag} {attr}> loads a remote resource: {c.strip()[:80]}"
                    )

        style_attr = amap.get("style", "")
        if style_attr:
            self._scan_css(style_attr)

    def handle_endtag(self, tag: str) -> None:
        if tag == "style":
            self._in_style = False

    def handle_data(self, data: str) -> None:
        if self._in_style:
            self._scan_css(data)

    def _scan_css(self, css: str) -> None:
        for m in URL_FUNC_RE.finditer(css):
            if is_remote(m.group(1)):
                self.errors.append(f"{self._loc()}: CSS url() loads a remote resource: {m.group(1).strip()[:80]}")
        for m in IMPORT_RE.finditer(css):
            self.errors.append(f"{self._loc()}: CSS @import is not allowed (inline styles only)")


def check(path: Path) -> list[str]:
    text = path.read_text(encoding="utf-8", errors="replace")
    parser = Checker()
    parser.feed(text)
    parser.close()
    errors = list(parser.errors)
    if not parser.saw_doctype:
        errors.insert(0, "missing <!doctype html> at the top of the file")
    if not parser.saw_charset:
        errors.insert(0, "missing <meta charset=\"utf-8\">")
    return errors


def main() -> int:
    ap = argparse.ArgumentParser(description="Check that an HTML file is self-contained and JS-free.")
    ap.add_argument("html_file", help="Path to the generated .html file")
    args = ap.parse_args()

    path = Path(args.html_file).expanduser().resolve()
    if not path.is_file():
        print(f"error: not a file: {path}", file=sys.stderr)
        return 1

    errors = check(path)
    if errors:
        for e in errors:
            print(f"error: {e}", file=sys.stderr)
        return 1

    print(f"ok: {path}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
