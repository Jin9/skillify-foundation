# Run Report Template

**Parent:** [`../templates.md`](../templates.md)
**Owner role:** Orchestrator (deterministic; not an LLM role) · **`template_version`:** 0.1.0

The post-run summary emitted by the orchestrator at terminal state. One per workflow run. Distinct from per-iteration `docs/v0.X-plan.md` (input) and from per-component QA/Reviewer outputs (mechanical state).

**File location:** `<workflow_root>/.squad-run/run-report.md` (per run, generated at terminal state).

**Why this template exists:** v0.6 dry-run #1 ended with `KNOWN_ISSUES.md` + `ADVISORY.md` + `run-log.md` as ad-hoc post-run prose. v0.7 standardizes the shape so `agent-workflow-postmortem` can be wired in later, and so cross-run regression detection (FEEDBACK#P3) has stable structure to diff.

---

## Required sections (frontmatter + 8 H2 sections)

### Frontmatter (YAML)

```yaml
---
template_version: 0.1.0
run_id: <e.g., ecom-mvp-2026-05-08-002>
iteration: v0.X
started_at: YYYY-MM-DD
ended_at: YYYY-MM-DD
mode: hybrid | claude-only | gpt-only
terminal_state: Done | ShipWithCaveats | Failed | HardFail
total_subagents_dispatched: N
expected_cost_usd: "<low>-<high>"
actual_cost_usd: "~<n>"
deviations_from_docs: N
findings_total: { high: N, medium: N, low: N }
---
```

### H2 sections (in order)

1. **`## Summary`** — 2-3 paragraph human narrative. What ran, what landed, what didn't, why the terminal state.
2. **`## Acceptance gate status`** — table: each gate from the iteration's `v0.X-plan.md` § Acceptance gates with pass/fail status.
3. **`## Stage outcomes`** — per-stage table: stage / verdict / sub-agent count / output paths / repair-cycle count.
4. **`## Closure of prior-run highs`** — for each high finding from the previous run, status: Closed | Partial | Open. Cite finding ID.
5. **`## New findings`** — top 10 by severity. Cross-link to the per-component / layer-2 finding files. Don't dump all findings here — that's what the per-component files are for.
6. **`## Deviations`** — list every entry from `state.json.deviations_from_docs` with reason + impact observed during the run.
7. **`## Cost ledger`** — per-stage / per-role token totals if the orchestrator's cost ledger (STRATEGY#P1-3) is wired; otherwise: best-effort estimate vs the iteration's `expected_cost_usd` envelope.
8. **`## Recommendations for next iteration`** — 3-7 bullets. What goes into the next iteration's `v0.X+1-plan.md`. This section feeds back into FEEDBACK.md / STRATEGY.md.

Optional 9th: `## Manual verification notes` — operator-run smoke-tests, visual conformance checks, etc.

Trailing: `## Change log` (table).

---

## Frontmatter schema

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "run-report/v0.1.0",
  "type": "object",
  "additionalProperties": false,
  "required": [
    "template_version", "run_id", "iteration", "started_at", "ended_at",
    "mode", "terminal_state", "total_subagents_dispatched",
    "expected_cost_usd", "actual_cost_usd",
    "deviations_from_docs", "findings_total"
  ],
  "properties": {
    "template_version":          { "type": "string", "pattern": "^\\d+\\.\\d+\\.\\d+$" },
    "run_id":                    { "type": "string", "minLength": 1 },
    "iteration":                 { "type": "string", "pattern": "^v\\d+\\.\\d+$" },
    "started_at":                { "type": "string", "pattern": "^\\d{4}-\\d{2}-\\d{2}$" },
    "ended_at":                  { "type": "string", "pattern": "^\\d{4}-\\d{2}-\\d{2}$" },
    "mode":                      { "enum": ["hybrid", "claude-only", "gpt-only"] },
    "terminal_state":            { "enum": ["Done", "ShipWithCaveats", "Failed", "HardFail"] },
    "total_subagents_dispatched":{ "type": "integer", "minimum": 1 },
    "expected_cost_usd":         { "type": "string" },
    "actual_cost_usd":           { "type": "string" },
    "deviations_from_docs":      { "type": "integer", "minimum": 0 },
    "findings_total": {
      "type": "object", "additionalProperties": false,
      "required": ["high", "medium", "low"],
      "properties": {
        "high":   { "type": "integer", "minimum": 0 },
        "medium": { "type": "integer", "minimum": 0 },
        "low":    { "type": "integer", "minimum": 0 }
      }
    }
  }
}
```

---

## Negative example — content-thin report

```markdown
---
run_id: foo
terminal_state: Done
---
# Run report
The run finished. KNOWN_ISSUES has stuff. Done.
```

What's wrong:

1. Missing required frontmatter fields (`template_version`, `iteration`, `mode`, `started_at`, `ended_at`, `total_subagents_dispatched`, `expected_cost_usd`, `actual_cost_usd`, `deviations_from_docs`, `findings_total`).
2. No H2 sections; reader can't tell what gate passed/failed, what the stages did, what the next iteration should pick up.
3. Body is content-free; `agent-workflow-postmortem` skill (when wired in) couldn't ingest this.

A good run report fits in 200-400 lines and is the single artifact a future Claude session reads to understand what happened in the run without trawling `run-log.md` + 16 per-component finding files.
