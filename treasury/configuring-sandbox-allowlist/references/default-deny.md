# Why default-deny matters

The sandboxed implement runner has only one network: `proxy-net`, which
the compose file marks `internal: true`. The runner literally cannot
route to the host network — its only path to the outside is through
tinyproxy, and tinyproxy denies by default unless the requested host
matches `filter`.

This is the property the scaffold relies on:

> If a host is not in `filter`, the runner cannot reach it.

Anything that breaks that property breaks the sandbox.

## What breaks default-deny

- Wildcards that match too much (see `pattern-rules.md`).
- Internal-host redirects that resolve to private space.
- Adding `127.0.0.1`-resolving hostnames that the runner can hit
  directly via `host.docker.internal` (already allowlisted for the
  LiteLLM gateway use case — be careful).
- Editing `tinyproxy.conf` to widen `Allow`, change
  `FilterDefaultDeny`, or add `ConnectPort` entries.

This skill refuses each of those.

## Verifying enforcement

Before each sandboxed implement run, the dispatcher attempts to fetch a
non-allowlisted host (`example.com` is the canonical test). If that
fetch *succeeds*, the implement runner aborts the workflow rather than
ship code that ran outside the sandbox.

```bash
just doctor --full
```

Exercises this path end-to-end. After any filter change, the user must
re-run this command before declaring the change complete.

## What "rebuild the proxy" actually does

```bash
docker compose -f docker/sandbox-compose.yml build proxy --no-cache
```

Tinyproxy reads `filter` at startup. Without `--no-cache`, Docker may
serve a layer where the old `filter` is baked in. Without restarting
the proxy container, the running tinyproxy instance keeps its old
filter loaded.

```bash
docker compose -f docker/sandbox-compose.yml up -d proxy
```

…restarts the container with the new image.

## Why this skill never builds

Two reasons:

1. The user is the one who decides when to interrupt running workflows
   that depend on the proxy. A surprise rebuild can kill an active run.
2. Building involves Docker context with credentials potentially in
   environment variables. Surface the command; let the user run it in
   their shell.
