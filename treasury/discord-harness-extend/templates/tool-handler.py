"""Skeleton for adding a MySQL operator tool to the harness. Copy a block into tools/db_ops.py,
rename, and wire the real query. Then classify it in config.yaml under `tools:` — a handler with no
config class is silently DENY and is never offered to the model.

Keep the in-handler re-validation: the safety gate approving a CONFIRM call does NOT mean the
arguments are safe (defense in depth). The model supplies only identifiers/values — NEVER raw SQL;
identifiers are checked against the allowlist, values always go through %s placeholders, so model
output cannot inject SQL. This file is a template, not an importable module.
"""
from __future__ import annotations

from contracts import ToolContext
# Re-validate identifiers in every handler (db_ops also imports table_projection / allowed_database):
from policy import allowed_column, allowed_table  # noqa: F401
# These helpers already live in tools/db_ops.py — reuse them when you paste a block in:
#   _ident(name)      → shape-check + backtick-quote an identifier (blocks injection chars)
#   _projection(table)→ the allowlisted SELECT column list ('*' only when the table permits all)
#   _get_conn()       → opens the connection; enforces allowed_database(db.database) inside
#   _rows_to_str(...) → JSON-encode + truncate rows for the model
from . import tool


# ── ALLOW variant: a read, auto-runs (still re-validates identifiers) ──────────
@tool(
    "my_read_tool",                       # MUST match the key under `tools:` in config.yaml
    "One line: what it reads and any limits (shown to the model).",
    {
        "type": "object",
        "properties": {
            "table": {"type": "string", "description": "Allowlisted table name"},
            "key_column": {"type": "string", "description": "Column to match on"},
            "key_value": {"description": "Value to match (string or number)"},
        },
        "required": ["table", "key_column", "key_value"],
    },
)
def my_read_tool(args: dict, ctx: ToolContext) -> str:
    table = allowed_table(str(args["table"]))             # raises PermissionError if off-allowlist
    col = allowed_column(table, str(args["key_column"]))  # re-validate even for ALLOW
    sql = f"SELECT {_projection(table)} FROM {_ident(table)} WHERE {_ident(col)} = %s LIMIT 1"
    with _get_conn() as conn, conn.cursor() as cur:       # _get_conn enforces allowed_database
        cur.execute(sql, (args["key_value"],))            # value via %s — never interpolated
        row = cur.fetchone()
    return _rows_to_str([row] if row else [])


# ── CONFIRM variant: a destructive write (gated by the approval button) ───────
@tool(
    "my_write_tool",
    "One line: the irreversible effect this performs.",
    {
        "type": "object",
        "properties": {
            "table": {"type": "string", "description": "Allowlisted table name"},
            "key_column": {"type": "string", "description": "Column identifying the row(s)"},
            "key_value": {"description": "Value identifying the row(s) to affect"},
        },
        "required": ["table", "key_column", "key_value"],
    },
)
def my_write_tool(args: dict, ctx: ToolContext) -> str:
    table = allowed_table(str(args["table"]))             # defense in depth: re-check identifiers
    key_col = allowed_column(table, str(args["key_column"]))
    sql = f"DELETE FROM {_ident(table)} WHERE {_ident(key_col)} = %s"   # parameterized — no SQL from model
    with _get_conn() as conn, conn.cursor() as cur:
        n = cur.execute(sql, (args["key_value"],))
    return f"affected {n} row(s) in {table}"


# config.yaml — add under `tools:` (omit = silently DENY):
#   my_read_tool: ALLOW
#   my_write_tool: CONFIRM
# and allowlist any new identifiers under `db:` (db.allowed_tables / db.allowed_columns).
#
# If this handler lives in a NEW module tools/<x>.py (not db_ops.py), also add `from . import <x>`
# to the import block in tools/__init__.py.
