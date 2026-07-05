"""Skeleton test scenario for a new MySQL operator tool. Paste into tests/test_harness.py main() after
the existing agent scenarios (A–F), using the helpers already defined there: fmsg / fcall / fresp /
script / env / approve / deny / cfg / dblog, plus the fake DB connection (db_ops._get_conn is already
monkeypatched to _FakeConn and `store` records the SQL + params). Assert BOTH the side effect (the
parameterized SQL / params, or that nothing ran) AND the policy outcome. Snippet, not a runnable module.
"""

# ── ALLOW tool: auto-runs, no approval ────────────────────────────────────────
store.clear()
litellm.acompletion = script([
    fresp(fmsg(calls=[fcall("t1", "my_read_tool", {"table": "users", "key_column": "id", "key_value": 1})])),
    fresp(fmsg(content="Done.")),
])
out = asyncio.run(agent.run_agent(env("use my_read_tool"), config=cfg, audit=dblog, approval_cb=approve))
assert out == "Done."
assert store["params"] == [1]                 # value went through %s, not interpolated into the SQL
print("my_read_tool: ALLOW auto-ran  OK")

# ── CONFIRM tool: runs only when approved ─────────────────────────────────────
store.clear()
litellm.acompletion = script([
    fresp(fmsg(calls=[fcall("t1", "my_write_tool", {"table": "users", "key_column": "id", "key_value": 9})])),
    fresp(fmsg(content="Done it.")),
])
asyncio.run(agent.run_agent(env("use my_write_tool"), config=cfg, audit=dblog, approval_cb=approve))
assert store.get("sql", "").startswith("DELETE FROM `users`")   # the effect's SQL was built + executed
print("my_write_tool: CONFIRM(approved) ran  OK")

# declined → effect must NOT happen (no SQL reaches the fake cursor)
store.clear()
litellm.acompletion = script([
    fresp(fmsg(calls=[fcall("t1", "my_write_tool", {"table": "users", "key_column": "id", "key_value": 9})])),
    fresp(fmsg(content="Skipped.")),
])
asyncio.run(agent.run_agent(env("use my_write_tool"), config=cfg, audit=dblog, approval_cb=deny))
assert "sql" not in store                      # declined → handler never ran, no SQL executed
print("my_write_tool: CONFIRM(declined) did not run  OK")

# ── guard: even an APPROVED op on an UN-allowlisted table is blocked ───────────
store.clear()
litellm.acompletion = script([
    fresp(fmsg(calls=[fcall("t1", "my_write_tool", {"table": "secrets", "key_column": "id", "key_value": 1})])),
    fresp(fmsg(content="Blocked.")),
])
asyncio.run(agent.run_agent(env("write to secrets"), config=cfg, audit=dblog, approval_cb=approve))
assert "sql" not in store                      # allowed_table raised → handler aborted before any SQL
print("my_write_tool: guard blocked an off-allowlist table  OK")
