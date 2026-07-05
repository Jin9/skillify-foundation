---
name: discord-harness-operate
description: >-
  Run, configure, monitor, and deploy this Discord MySQL-operator bot without changing its tool code.
  Use when the user says "run the harness bot", "swap the work model", "check healthz / monitor the
  bot", "set up the local MySQL test stack", or "deploy the operator to Docker or AWS". Covers the
  .env-secrets vs config.yaml-policy split, the message_content intent, model swap via models.work,
  the Grafana/Prometheus stack, the /healthz endpoint, and the local and AWS EC2 deploy
  targets. Do NOT use to add or modify tools or change the policy gate (use discord-harness-extend), or
  to edit the audit, state, or contracts internals.
---

# discord-harness-operate

Run, configure, monitor, and deploy the Discord MySQL-operator bot — without changing its tool code.

## When to use

The agent has been asked to operate the running system: start the bot, swap the model, watch it, or
deploy it. For adding/modifying tools or touching the policy gate, use `discord-harness-extend`.

This harness is now a single-purpose **MySQL operator** — only the 5 `mysql_*` tools exist (reads +
validated upserts auto-run; `mysql_delete_row` gates on a human). There is no research assistant and no
fs/shell/web/EKS tools; anything unlisted fails closed to DENY. There is ONE config file (`config.yaml`)
— no profile to switch between.

## First: the secrets / policy split

Two files, never mixed (`.env.example`, `config.yaml`):
- **`.env` = secrets only** — `DISCORD_TOKEN`, `OWNER_ID`, optional `CHANNEL_ID` / `ROLE_ID`, the
  provider API key matching `models.work`, the MySQL connection (`DB_URL` or `MYSQL_*`), optional
  `OLLAMA_API_BASE` / `HEALTH_PORT` / `AUDIT_PATH` / `AUDIT_STDOUT` / `GF_ADMIN_PASSWORD`.
- **`config.yaml` = policy** — `system_prompt`, model routing, allowlists, caps, tool classes, and the
  `db.database` / `allowed_tables` / `allowed_columns` operation allowlist. Never put a secret here;
  never put policy in `.env`.

Auth fails closed: empty `allowlist.users` + no `OWNER_ID` = nobody authorized (`policy.py:50`). Run
`scripts/preflight.py` before starting to catch missing keys and a locked allowlist.

## Tasks

### Run
- Local: `python bot.py` (in the `.venv`, with `.env` filled). No `HARNESS_CONFIG` needed — `config.yaml`
  IS the operator profile. `bot.py` exits if `DISCORD_TOKEN` is unset (`bot.py:211`).
- Docker (dev): from the repo root, `docker compose -f deploy/local/docker-compose.yml up -d --build`
  (seeded MySQL 8 + bot + Grafana). There is no root `docker-compose.yml` anymore.
- The `message_content` intent must be ON in the Discord Developer Portal AND set in `bot.py:26`, or
  message bodies arrive empty.
- Liveness: `GET /healthz` on `HEALTH_PORT` (default 8080) → `{ok, uptime_s, started_at, model, bot_user}`
  (`bot.py:51`). Full env table: `references/run-and-config.md`.

### Configure
- **Swap the work model:** edit `models.work` in `config.yaml` to one of the two wired routes —
  local `ollama_chat/qwen3:4b` (default; no key; needs `ollama serve` + `ollama pull qwen3:4b`, and
  `models.params.num_ctx`/`think` set) or cloud `gemini/gemini-2.5-flash-lite` (needs `GEMINI_API_KEY`
  in `.env`). Relaunch, or `docker compose -f deploy/local/docker-compose.yml restart discord-operator`
  (config is live-mounted read-only).
- **Single config:** `config.yaml` is the only policy file — there is no second profile to switch
  between. `HARNESS_CONFIG` still points at an alternate policy file (e.g. an isolated test config) but
  is not needed for normal operation (`policy.py:20`). `config.yaml` lists only the `mysql_*` tools, so
  everything else fails closed to DENY.
- Allowlists, caps, tool classes, and the `db.*` table/column allowlist are all config — restart to
  apply. LiteLLM routes any provider with `litellm.drop_params = True` (`agent.py:27`). More in
  `references/run-and-config.md`.

### Monitor
- The Grafana board at `:3000` (Prometheus `:9090`, exporter `:9108`) from the compose stack is the
  monitoring UI. The exporter is read-only — it reads `audit.jsonl` + `config.yaml` only and reuses the
  bot's own `AuditLog.verify()`, so its integrity result cannot drift (`metrics_exporter.py`).
- Liveness any time: `curl -s localhost:8080/healthz`.
- Audit integrity any time:
  `python -c "from audit import AuditLog; print(AuditLog('audit.jsonl').verify())"`.
- Metrics + event glossary: `references/monitor.md`.

### Deploy
Two targets (`references/deploy.md` has commands + topology); both build from the repo-root `Dockerfile`:
- **Local dev stack** — `deploy/local/docker-compose.yml`: seeded MySQL 8 + bot + Grafana, fully
  offline. From the repo root: `docker compose -f deploy/local/docker-compose.yml up -d --build`.
- **AWS EC2 prod** — `deploy/operator/docker-compose.yml` + `deploy/operator-userdata.sh`: g4dn spot,
  containerized Ollama + bot reaching a real MySQL over the host VPN + Grafana, secrets from SSM
  SecureString, EventBridge start/stop schedule.
- The agent surfaces cloud/secret commands (`aws ssm put-parameter`, instance launch) for the human to
  run; it does not execute them.

## Verify after any change

```bash
curl -s localhost:8080/healthz            # liveness + active model
curl -s localhost:9108/metrics | head     # exporter scrape (chain status, tools, outcomes as gauges)
python -c "from audit import AuditLog; print(AuditLog('audit.jsonl').verify())"   # True
```

## Hard rules

- Secrets never enter `config.yaml`; policy never enters `.env`.
- Don't disable the `message_content` intent or the bot goes deaf.
- Changing a tool's behavior or its safety class is NOT an operate task — that's
  `discord-harness-extend` (it needs a test + the registration cross-check).

## References

| Need | File |
|------|------|
| Env var table, intents, model swap, LiteLLM routing | `references/run-and-config.md` |
| /healthz fields, audit events, exported metrics, Grafana stack | `references/monitor.md` |
| Local dev / AWS EC2 operator deploy targets | `references/deploy.md` |
| Pre-flight check (keys present, config parses, auth not locked) | `scripts/preflight.py` |
