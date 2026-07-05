# Test recipe — prove the policy outcome

`tests/test_harness.py` is a plain script of sequential asserts, run with `python tests/test_harness.py`
(no pytest, no Discord token, no API key). litellm is mocked, so the agent loop is deterministic. To
isolate a case, edit/comment `main()`.

## Mock the DB connection (`tests/test_harness.py:95`)

The MySQL tools never touch a live database in the suite. `db_ops._get_conn` is the single indirection
point — monkeypatch it to a fake connection whose cursor records the SQL + params it was handed:

```python
from tools import db_ops
store = {}

class _FakeCur:
    lastrowid = 1
    def __enter__(self): return self
    def __exit__(self, *a): return False
    def execute(self, sql, params=None):
        store["sql"] = sql; store["params"] = list(params) if params else []
        return 1                      # rows affected
    def fetchone(self): return {"id": 1, "name": "…"}
    def fetchall(self): return [{"id": 1}, {"id": 2}]

class _FakeConn:
    def __enter__(self): return self
    def __exit__(self, *a): return False
    def cursor(self): return _FakeCur()

db_ops._get_conn = lambda: _FakeConn()      # (`tests/test_harness.py:117`)
```

Now you can call a handler directly and assert the exact parameterized SQL (`store["sql"]` /
`store["params"]`) — that is how the db_ops block proves identifiers are allowlist-checked and values go
through `%s` (`tests/test_harness.py:92`).

## The mock LLM helpers (`tests/test_harness.py:161`)

```python
def fmsg(content=None, calls=None): ...   # one assistant message (text OR tool_calls)
def fcall(i, n, a): ...                    # one tool call: id, name, args-dict
def fresp(m): ...                          # wraps a message as a litellm response (+usage)
def script(rs_):                           # returns an async acompletion that yields rs_ in order
    it = iter(rs_)
    async def _ac(**kw): return next(it)
    return _ac
```

Drive the loop by assigning `litellm.acompletion = script([...])` — one `fresp` per model turn. The
LAST turn returns `fmsg(content="…")` with no calls, so the loop replies and exits. The agent-loop
scenarios run with the fake `db_ops._get_conn` from above still installed.

## Fixtures already in main()

- `cfg = policy.load_config()` — the config (`tests/test_harness.py:179`).
- `dblog = AuditLog(tempfile.mktemp(...))` — a throwaway audit log (`tests/test_harness.py:180`).
- `env(intent)` — builds a HandoffEnvelope (`tests/test_harness.py:183`). `approve` / `deny` —
  approval callbacks (`tests/test_harness.py:186`).

## The six existing agent scenarios (the shapes to copy)

| Case | Proves | Anchor |
|------|--------|--------|
| A | `mysql_query_rows` (ALLOW) auto-runs; replies in Thai | `tests/test_harness.py:193` |
| B | `mysql_delete_row` (CONFIRM, approved) executes | `tests/test_harness.py:204` |
| C | `mysql_delete_row` (CONFIRM, declined) never executes | `tests/test_harness.py:213` |
| D | unlisted tool → DENY, loop recovers | `tests/test_harness.py:224` |
| E | identical calls → loop detector terminates (`🛑 Stopped`) | `tests/test_harness.py:233` |
| F | approved op on an UNLISTED table → blocked by handler guard | `tests/test_harness.py:241` |

## Add a scenario for your tool

Pick the assertion that matches the class:
- ALLOW → assert the effect happened (e.g. `store["sql"]` / `store["params"]`) and `out == "<final reply>"`.
- CONFIRM → run once with `approve` (effect happens) and once with `deny` (effect does NOT happen).
- guard → feed worst-case args (e.g. an un-allowlisted table) with `approve`, assert the op was still
  blocked — look for a `guard_block` event in the audit log.

Skeleton: `templates/test-scenario.py`. After adding, keep the final `assert dblog.verify() is True`
(the run's audit chain must stay intact), and end with all `OK` lines.

## Run + verify

```bash
python tests/test_harness.py        # every line OK, ends "ALL HARNESS LOGIC TESTS PASSED"
python -c "from audit import AuditLog; print(AuditLog('audit.jsonl').verify())"   # True
```
