# Run & configure

## Environment variables (`.env` — secrets only)

| Var | Purpose | Required |
|-----|---------|----------|
| `DISCORD_TOKEN` | bot token; `bot.py` exits without it (`bot.py:211`) | yes |
| `OWNER_ID` | your Discord user id; auto-added to allowlist + approvers | recommended |
| `CHANNEL_ID` | channel(s) the bot answers in (comma/space sep). Set → answers ONLY there, DMs blocked | optional |
| `ROLE_ID` | role(s) allowed to invoke (enables the privileged members intent) | optional |
| `GEMINI_API_KEY` | provider key for the `models.work: gemini/…` fallback route (not needed on the default local Ollama) | conditional |
| `OLLAMA_API_BASE` | local model endpoint; default `http://localhost:11434` (Docker overrides to host) | optional |
| `DB_URL` or `MYSQL_HOST/PORT/USER/PASSWORD/DB` | the ONE MySQL database the operator manages (`DB_URL` wins if both set) | yes |
| `HEALTH_PORT` | `/healthz` port (default 8080) | optional |
| `AUDIT_PATH` | point the audit chain at a durable volume | optional |
| `AUDIT_STDOUT` | mirror each audit line to stdout for the container log driver | optional |
| `HARNESS_CONFIG` | point at an alternate policy file (rarely needed — `config.yaml` is the only profile) | optional |
| `GF_ADMIN_PASSWORD` | Grafana admin password (default `admin`) | optional |

## config.yaml (policy, committed)

`config.yaml` is the single operator policy file (it IS the operator profile). Sections: `system_prompt`;
`models.work` / `models.params`; `allowlist.users` / `channels` / `roles`; `approvers`; `caps`
(max_tool_iterations, max_tokens_per_request [0 = unlimited], max_tool_result_chars, rate_limit_per_min, daily_spend_usd);
`tools:` classes (the 5 `mysql_*` tools); and the MySQL operation allowlist `db.database` /
`db.allowed_tables` / `db.allowed_columns`. Identifiers are validated against those lists; values always
go through parameterized queries. There are no fs/shell/web/EKS sections — those capabilities were
removed, so any unlisted tool fails closed to DENY.

## Swap the work model

Edit one line — `models.work`. This harness wires exactly two routes:
- `ollama_chat/qwen3:4b` — the default; local, no API key; needs `ollama serve` + `ollama pull qwen3:4b`
  (with `models.params.num_ctx: 8192` / `think: true` set in `config.yaml` — thinking ON keeps the
  reply clean by routing qwen3's reasoning to a separate field the loop drops).
- `gemini/gemini-2.5-flash-lite` — cloud fallback; needs `GEMINI_API_KEY`.

Then relaunch `python bot.py`, or `docker compose -f deploy/local/docker-compose.yml restart
discord-operator` (config is live-mounted read-only).

## Single config (no profile swap)

There is one policy file — `config.yaml` — and it is the operator profile, so there is nothing to switch
between. `python bot.py` loads it with no `HARNESS_CONFIG` set. `HARNESS_CONFIG` still lets you point at
an alternate policy file (e.g. an isolated test config); `policy.py:20` resolves a relative name next to
the repo. The config IS the allowlist — it lists only the `mysql_*` tools, so unlisted tools fail closed
to DENY and the operator can run only validated MySQL ops.

## Provider-neutral routing

The loop sends one OpenAI-format request; LiteLLM adapts it to whatever provider `models.work` names,
and `litellm.drop_params = True` (`agent.py:27`) drops params a provider doesn't support (e.g. Ollama
ignores OpenAI-only fields; `models.params.num_ctx` is sent through to Ollama).

## Intents

`intents.message_content = True` (`bot.py:26`) — also enable Message Content Intent in the Developer
Portal. The members intent is requested only when a role allowlist or `ROLE_ID` is set (`bot.py:30`).
