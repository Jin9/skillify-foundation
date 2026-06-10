# Severity, exit-code & risk-tiered gate policy

The gate is risk-tiered, not a blunt kill: an over-aggressive auto-fail takes a
healthy agent offline just as surely as drift does (per *Strict Schema
Evolution* and *Agent command safety*). Map severity → gate; only `block`
findings fail the build.

## Severity ladder
| severity | meaning |
|----------|---------|
| `critical` | exploitable / destructive command surface (C1) |
| `high` | breaking evolution, lock-in, broad scope, injection, missing HITL |
| `medium` | portability friction, weak versioning/identity/creds |
| `low` | minor, advisory |
| `info` | context, skipped check, recommendation |

## Default severity → gate
```
critical → block
high     → block
medium   → warn
low      → warn
info     → info
```
`block` ⇒ counts toward non-zero exit. `warn`/`info` ⇒ printed, exit-neutral.

## Exit codes (the gate signal)
- `0` — zero `block` findings (warns allowed). Pipeline proceeds.
- `1` — ≥1 `block` finding. Pipeline must fail.
- `2` — usage / parse / unreadable-spec error (treat as failure).

## Config override (`.spec-validator.yaml`)
```yaml
compatibility_mode: backward      # backward | forward | full
transitive: false                 # parsed but not yet read by validate_spec.py
target_models: [openai, anthropic, gemini]   # informs P1/P2/P5 (not yet read by validate_spec.py)
gate_overrides:                    # per-rule severity→gate remap (risk tiering)
  P2: info                         # e.g. accept strict-mode drift, just notify
  C5: block                        # e.g. enforce pinning hard in this repo
ignore:                            # suppress by rule id + locator glob
  - rule: P6
    locator: "schemas/legacy/**"
fail_on: [critical, high]          # which severities count as block (default)
```
Precedence: `ignore` → `gate_overrides` → `fail_on` default.

## Report schema (`spec-validation.json`)
```json
{
  "summary": {"specs": 0, "block": 0, "warn": 0, "info": 0,
              "by_axis": {"portability": 0, "command-safety": 0, "schema-evolution": 0}},
  "exit_code": 0,
  "findings": [
    {"axis": "command-safety", "id": "C1", "severity": "critical",
     "gate": "block", "spec": "tools/deploy.json",
     "locator": "$.tools[2].command", "message": "...",
     "rule_ref": "command-safety-rules.md#C1"}
  ]
}
```
`spec-validation.md` renders the same data via
`templates/spec-validation-report-template.md`.

Recommended operating pattern (state in the report, never auto-apply): contract
as single source of truth, risk-tiered automation (autonomous only for
low-risk, notify-human for medium, block for high/critical), evolution-safe
registry changes, exact version pins as the default with invalidation as the
exception path.
