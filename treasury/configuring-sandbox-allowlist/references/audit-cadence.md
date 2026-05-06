# Audit cadence

Allowlists drift. Hosts get added for one workflow and never removed.
Patterns that were precise become stale. This is how to audit.

## Cadence

| Trigger | Audit scope |
|---|---|
| Quarterly | Full filter review. |
| After any postmortem citing the sandbox | Targeted: every host added since the last incident. |
| After off-boarding a researcher | Targeted: every host added by that researcher in the last 90 days. |
| Before raising the team-wide cap default | Full filter review. |

## Audit data sources

### Proxy logs

```bash
docker compose -f docker/sandbox-compose.yml logs proxy --since 720h \
  | grep -E 'CONNECT|ALLOW|DENY'
```

Per-host hit counts:

```bash
docker compose -f docker/sandbox-compose.yml logs proxy --since 720h \
  | awk '/CONNECT/ {print $NF}' \
  | sort | uniq -c | sort -rn
```

Hosts in `filter` with zero hits in the last 30 days are candidates for
removal.

### State / runlog cross-reference

For each candidate-for-removal host, check whether any recent workflow
referenced it (search `runlog.md` and `stages/*.md`):

```bash
grep -ril '<host>' .agent/runlog.md .agent/stages/*.md \
  $(find . -name 'runlog.md' -path '*/.agent/*')
```

If found, the host is in active use even if recent traffic is low.

## Audit table shape

```markdown
| Section | Pattern | Last used | Hits/30d | Verdict |
|---------|---------|-----------|----------|---------|
| Model APIs | `^api\.openai\.com$` | 2026-05-05 | 142 | keep |
| Package managers | `^pypi\.org$` | 2026-05-04 | 38 | keep |
| Public source hosts | `^codeload\.github\.com$` | 2026-04-08 | 0 | recommend remove |
| (custom) | `^foo\.partner\.example\.com$` | 2026-02-01 | 0 | recommend remove |
```

The audit produces this table. Removal happens via a separate Remove
operation per host the user agrees to drop.

## Patterns that always need scrutiny

- Any pattern added more than 90 days ago whose use case is no longer
  active.
- Any pattern with `|` joining ≥ 4 alternatives — usually a sign the
  list is growing organically and could be split into a section.
- Any pattern matching a CDN / object-store host (often covers more
  than the user realized).
- Any pattern added by a researcher who has off-boarded.

## What the audit does NOT do

- Recommend new entries. The audit shrinks; new entries enter through
  the Add operation with a stated reason.
- Modify `tinyproxy.conf`, `sandbox-compose.yml`, or `Dockerfile`.
- Trigger a proxy rebuild. The user rebuilds after any agreed removal.
