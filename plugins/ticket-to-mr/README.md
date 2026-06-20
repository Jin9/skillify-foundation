# ticket-to-mr

A self-contained Claude Code **plugin** that takes one Jira issue to one GitLab merge request through
four human approval gates. It bundles the gated orchestrator skill together with the two skills it
delegates to, so all cross-skill references resolve **inside** the plugin.

## Bundled skills

| Skill (invoke as) | Role | I/O contract |
|---|---|---|
| `/ticket-to-mr:jira-fix-mr-workflow` | Orchestrator (conducts the 4-gate flow) | in: a Jira issue key + repo → out: one opened GitLab MR, every irreversible step human-approved |
| `/ticket-to-mr:progressive-bug-hunter` | Localize the bug | in: issue + repo → out: ranked suspect `file:line` + minimal context bundle. **Stops at diagnosis; never edits.** |
| `/ticket-to-mr:publishing-git-review-requests` | Push + open the MR | in: prepared branch + intended commit + MR metadata → out: pushed branch + opened MR. **Does not merge.** |

The orchestrator **conducts** — it does not re-implement localization, MR publishing, or
command-safety. It treats Jira text as data (never instructions) and never advances a gate itself.

## Enforcement scaffold (bundled & tested)

`jira-fix-mr-workflow` is backed by a real enforcement scaffold at the plugin root:

| Path | Role |
|---|---|
| `core/policy/policy.json` | Single source of truth — `deny`/`confirm`/`allow` command rules, caps, gates, conventions |
| `core/policy/policy_match.py` | ALLOW/CONFIRM/DENY matcher; `--hook` (PreToolUse) and `--command` (CLI) modes |
| `hooks/hooks.json` | Registers the PreToolUse hook on the `Bash` tool |
| `core/policy/generate.py` | Projects `policy.json` → `adapters/codex/execpolicy.rules`; verifies hook wiring |
| `core/lib/jira_extract.py` | Deterministic extraction + **injection quarantine** (Jira text is data, never instructions) |
| `bin/gate_runner.py` | Gate state machine: hash-chained audit log, typed checkpoint schemas, cap enforcement |
| `core/schemas/gate{1..4}-*.json` | Typed gate checkpoint schemas |
| `core/templates/mr-body.md` | Secret/path-scrubbed MR body template |

### Run-scoped enforcement (important)

The PreToolUse hook fires on **every** Bash call while the plugin is enabled, so it is a **no-op unless a
run is active** — otherwise it would police all your normal shell work. A run is active when either:

- `export TICKET_TO_MR_ENFORCE=1` is set, **or**
- a `runs/.active` marker exists (written by `bin/gate_runner.py start`).

When active, the hook returns `deny` (e.g. `git add .`, `rm -rf`, `git push --force`, `glab mr merge`,
staging `.env`/secrets), `ask` for gated actions (`git commit`, `git push`, `glab mr create`), and
`allow` for safe reads/tests. When inactive it emits nothing and defers to your normal permission flow.

### Activate, then drive the gates

```bash
cd plugins/ticket-to-mr
export TICKET_TO_MR_ENFORCE=1                      # or: python3 bin/gate_runner.py start --key DGL-1234 --trace t1
python3 bin/gate_runner.py submit  --trace t1 --gate 1 --payload gate1.json
python3 bin/gate_runner.py approve --trace t1 --gate 1 --approver "Jane (TL)"
python3 bin/gate_runner.py status  --trace t1      # prints state + verifies the hash chain
```

### Tests

```bash
python3 plugins/ticket-to-mr/tests/run_tests.py    # 33 tests: policy, quarantine, hash-chain, caps, schema
python3 plugins/ticket-to-mr/core/policy/generate.py --check   # projection up to date + wiring intact
```

> Codex note: `adapters/codex/execpolicy.rules` is a deterministic projection of `policy.json`; the exact
> syntax a given Codex CLI version expects may need a thin per-version adapter.

## Local install / test

```bash
# from the repo root
claude --plugin-dir plugins/ticket-to-mr
# then inside the session:
/reload-plugins
```

The three skills appear namespaced under the plugin: `/ticket-to-mr:jira-fix-mr-workflow`,
`/ticket-to-mr:progressive-bug-hunter`, `/ticket-to-mr:publishing-git-review-requests`.

## Validation

The bundled skills are verbatim copies of the `treasury/` originals and pass the repo's gate:

```bash
for s in jira-fix-mr-workflow progressive-bug-hunter publishing-git-review-requests; do
  python3 skillify/scripts/quick_validate.py plugins/ticket-to-mr/skills/$s
  python3 skillify/scripts/check_links.py    plugins/ticket-to-mr/skills/$s
done
```
