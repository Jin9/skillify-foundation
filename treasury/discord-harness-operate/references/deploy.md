# Deploy

Two targets. The agent documents and may run LOCAL docker/compose commands; for AWS it surfaces the
secret/instance commands for a human to run (it never executes `aws ssm put-parameter` or instance
launches itself). Both stacks build from the repo-root `Dockerfile` (the lean MySQL-operator image —
MySQL needs only PyMySQL). Full topology + boot sequence live in `deploy/README.md`.

## 1. Local dev stack — `deploy/local/docker-compose.yml`

Seeded MySQL 8 (`deploy/local/initdb/001-schema.sql`, tables `users` + `orders`, scoped `operator` user)
+ the operator bot + the Grafana board — exercise db ops offline, no AWS. Runs on local qwen3:4b by
default (`config.yaml` → `ollama_chat/qwen3:4b`), reached on the host via `host.docker.internal` — so
keep `ollama serve` running + `ollama pull qwen3:4b` (switch `models.work` to the Gemini line for cloud).
Run from the repo ROOT so the build context resolves:
```bash
cp .env.example .env            # fill DISCORD_TOKEN, OWNER_ID, (CHANNEL_ID/ROLE_ID); GEMINI_API_KEY only on Gemini
touch audit.jsonl               # so Docker bind-mounts a file, not a dir
docker compose -f deploy/local/docker-compose.yml up -d --build
docker compose -f deploy/local/docker-compose.yml logs -f discord-operator
# Grafana http://localhost:3000 (admin / $GF_ADMIN_PASSWORD), Prometheus :9090, exporter :9108, /healthz :8080
```
The compose `env:` overrides the MySQL connection to reach the `mysql` service by name. `config.yaml` is
live-mounted read-only — `docker compose -f deploy/local/docker-compose.yml restart discord-operator`
picks up a model / allowlist swap with no rebuild. See `deploy/local/README.md` for Thai/English test
prompts.

## 2. AWS EC2 operator — `deploy/operator/docker-compose.yml` + `deploy/operator-userdata.sh`

Single g4dn.xlarge spot VM running docker compose: containerized Ollama + the operator bot (host
networking, reaching a real MySQL over the host VPN) + the Grafana board. Discord is outbound-only (no
inbound SG), managed via SSM. One-time setup (human-run):
1. Secrets → SSM SecureString under `/operator/`: `DISCORD_TOKEN`, `OWNER_ID`, `CHANNEL_ID`, `ROLE_ID`,
   `DB_URL` (or `MYSQL_*`).
2. IAM instance role: `AmazonSSMManagedInstanceCore`; `ssm:GetParametersByPath` + `kms:Decrypt` on
   `/operator/*`. Scope the DB user to the one allowlisted database with least privilege.
3. Launch g4dn.xlarge spot, 60 GB gp3, private subnet + NAT, no inbound SG, the NVIDIA Container Toolkit
   (for the ollama container's GPU; CPU works without it), and a host VPN `vpn.service`; user-data =
   `deploy/operator-userdata.sh`.
4. Edit the allowlists in `config.yaml` (the safeguard) — `db.database`, tables, columns.
5. EventBridge Scheduler: start 08:00 / stop 21:00 weekdays (cron is UTC), targeting
   `ec2:StartInstances` / `ec2:StopInstances`.

Boot order each scheduled day: VM start → VPN up on host → `docker compose ... up` → bot retries the DB
on transient failures → Discord ready.

Verify on the box (SSM Session Manager):
```bash
cd /opt/operator
docker compose -f deploy/operator/docker-compose.yml ps          # ollama + discord-operator + grafana Up
curl -s localhost:8080/healthz
python -c "from audit import AuditLog; print(AuditLog('audit.jsonl').verify())"   # True (or AUDIT_PATH)
```

## Adding a safe OPERATION (config-only)

No code change: add a table to `db.allowed_tables` (+ optional `db.allowed_columns`) in `config.yaml`,
then set the tool's class under `tools:` (the model still can only call
`mysql_get_row/query_rows/insert_row/update_row/delete_row` with parameterized values). Restart
`discord-operator`. Adding a NEW tool handler is a code change — that's the `discord-harness-extend`
skill, not this one.
