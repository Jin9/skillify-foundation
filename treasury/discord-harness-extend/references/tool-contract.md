# Tool contract — how a tool is wired

The registry knows *what* a tool does; it never knows *whether* it's allowed. Authorization is
`policy.classify()` + config (see `safety-classes.md`). Source: `tools/__init__.py`.

## The decorator

`tools/__init__.py:30`
```python
def tool(name: str, description: str, parameters: dict):
    def deco(fn: ToolHandler) -> ToolHandler:
        REGISTRY[name] = ToolSpec(name, description, parameters, fn)
        return fn
    return deco
```

- `name` — the tool name the model calls; MUST match the key you add under `tools:` in config.
- `description` — shown to the model; state what it does and any limits.
- `parameters` — a JSON Schema object (`type: object`, `properties`, `required`).
- `ToolHandler = Callable[[dict, ToolContext], Union[Awaitable[str], str]]` (`tools/__init__.py:16`).
  Keep the `Union[...]` form — do not rewrite to `|` (Python ≤3.9 back-compat).

## Handler signature and context

A handler is `def fn(args: dict, ctx: ToolContext) -> str` (or `async def`). It returns a string the
loop feeds back to the model. `ToolContext` (`contracts.py:72`) is read-only provenance + config:

```python
@dataclass
class ToolContext:
    config: dict[str, Any]      # the loaded config.yaml
    correlation_id: str         # = the Discord message id
    user_id: int
    workspace: str = ""         # scratch dir, created if missing (unused by the MySQL tools)
    notes: list[str] = field(default_factory=list)
```

Read config through `ctx.config.get(...)` (e.g. a handler reads `ctx.config["db"]["allowed_tables"]`).

## How tools reach the model (and how DENY hides them)

`openai_schema()` (`tools/__init__.py:37`) builds the `tools=` payload; `dispatch()`
(`tools/__init__.py:46`) routes a call to the handler. The loop only OFFERS non-DENY tools
(`agent.py:65`):

```python
offered_tools = [s for s in tools.openai_schema()
                 if policy.classify(s["function"]["name"]) is not SafetyClass.DENY]
```

So an unlisted (DENY) tool is never advertised — the model can't even propose it.

## Registration side effects

`tools/__init__.py:58` imports the handler module(s) so their `@tool` decorators run:

```python
from . import db_ops     # MySQL operator tools
```

A NEW module must be added here, or its handlers never register.

## Example — read tool (ALLOW)

`tools/db_ops.py:109` `mysql_get_row` — the model supplies `{table, key_column, key_value}`; the
handler validates identifiers with `allowed_table` / `allowed_column`, projects only allowlisted
columns, and runs only **parameterized** SQL (`%s`), so model output cannot inject SQL or leak an
un-allowlisted column:
```python
@tool("mysql_get_row",
      "Fetch a single row by a key column from an allowlisted table.",
      {"type": "object",
       "properties": {"table": {"type": "string"}, "key_column": {"type": "string"},
                      "key_value": {"description": "Value to match"}},
       "required": ["table", "key_column", "key_value"]})
def mysql_get_row(args: dict, ctx: ToolContext) -> str:
    table = allowed_table(str(args["table"]))            # re-validate identifiers (defense in depth)
    col = allowed_column(table, str(args["key_column"]))
    sql = f"SELECT {_projection(table)} FROM {_ident(table)} WHERE {_ident(col)} = %s LIMIT 1"
    with _get_conn() as conn, conn.cursor() as cur:      # _get_conn enforces allowed_database
        cur.execute(sql, (args["key_value"],))           # value via %s — never interpolated
    ...
```

## Example — write tool (CONFIRM)

`tools/db_ops.py:223` `mysql_delete_row` is the destructive op (CONFIRM in config). Same shape: the
model supplies only `{table, key_column, key_value}`, the handler re-validates the identifiers and
builds a parameterized `DELETE`. Mirror this for any new tool: allowlist the identifiers, parameterize
the values, never let model output reach the SQL string.

## Files you touch to add one tool

| File | What |
|------|------|
| `tools/db_ops.py` (or a new `tools/<x>.py`) | the `@tool` handler + in-handler re-validation |
| `tools/__init__.py` | add `from . import <x>` ONLY if you created a new module |
| `config.yaml` | the `tools:` safety class (+ any new `db.allowed_tables` / `db.allowed_columns`) |
| `tests/test_harness.py` | a scenario asserting the policy outcome |
