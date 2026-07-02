# Monitor

Monitoring is the **Grafana board + the `/healthz` endpoint** only (the exporter aggregates
`audit.jsonl` off the request path). What the board surfaces: requests / errors / tokens (all-time +
today), approvals (✓/✗), denies, rejects, per-tool counters (ok / error / blocked / approved /
declined / deny) for the `mysql_*` tools, per-user activity, chain integrity, bot liveness, and
inferred per-request outcomes.

## /healthz

`GET /healthz` on `HEALTH_PORT` (default 8080), served on the bot's event loop (`bot.py:40`). Cheap by
design — returns `{ok, uptime_s, started_at, model, bot_user}` from cached config; it does NOT walk the
audit chain.

## Audit events (what appears in audit.jsonl)

`received`, `model_call` (tokens), `approval` (approved bool), `deny`, `tool_effect`, `tool_error`,
`guard_block`, `reject`, `replied`, `cap_exceeded`, `error`, `boot_probe`. One `correlation_id` (the
Discord message id) ties a whole request together — grep that id in `audit.jsonl` to see every record
for one request. The exporter's error set is `error`, `cap_exceeded`, `tool_error`, `guard_block`
(`metrics_exporter.py`).

## Chain integrity

`python -c "from audit import AuditLog; print(AuditLog('audit.jsonl').verify())"` → `True` / `False`.
verify() re-walks the SHA-256 `prev_hash`→`hash` chain (`audit.py:63`); a single altered byte returns
False. The Grafana exporter calls the same verifier, so the board reflects the live audit log.

## Grafana stack (Docker)

Bundled in both deploy stacks (`deploy/local/docker-compose.yml` and `deploy/operator/docker-compose.yml`):
exporter `:9108` → Prometheus `:9090` → Grafana `:3000` (user `admin`, password `GF_ADMIN_PASSWORD` or
`admin`). Ports 3000/9090/9108 are shared, so run only one stack at a time. The exporter mounts
`audit.jsonl` read-only and aggregates it directly.

The board's top **🎯 Focus** row foregrounds gate / audit / caps / model; lower rows keep throughput,
per-tool, outcomes, and per-user panels.

## Exported metrics (`:9108/metrics`)

Totals: `harness_requests(_today)`, `harness_errors(_today)`, `harness_tokens(_today)`,
`harness_denies`, `harness_rejects`, `harness_audit_records`. Labeled:
`harness_approvals{decision}`, `harness_tool_calls{tool,status}` (status `blocked` = guard-block),
`harness_outcomes{outcome}`, `harness_user_requests|errors{user_id}`. Health/integrity:
`harness_audit_chain_valid` (1/0/-1), `harness_up`, `harness_uptime_seconds`,
`harness_build_info{model,bot_user}`.

Focus signals (added): `harness_cap_exceeded(_today)` — hard-cap breaches broken out of the error
bucket; `harness_cap_{max_tool_iterations,max_tokens_per_request,rate_limit_per_min,daily_spend_usd}`
— configured ceilings mirrored from `config.yaml` so the board can draw budget lines (a cap emits its
gauge only when set to a non-zero value; a `0`/disabled cap such as unlimited tokens draws no line);
`harness_cost_usd(_today)` + `harness_price_usd_per_1m` — **estimated** LLM spend (total-token ×
blended rate); $0 on the default local qwen3:4b (no per-token price). Tune the rate via
`PRICE_USD_PER_1M` in `metrics_exporter.py` or the `LLM_PRICE_USD_PER_1M` env on the exporter service.
