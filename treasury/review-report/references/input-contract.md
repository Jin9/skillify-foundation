# Input contract

Canonical field catalog for the five inputs the workflow engine passes to
`review-report`. The SKILL.md body points here. No input may be inferred or
fabricated; missing/malformed/empty inputs trigger the `input_incomplete`
recovery (see `failure-modes.md`).

| Field | Type | Source | Notes |
|-------|------|--------|-------|
| `topic` | string | workflow input | The subject the report is about. The reviewer checks the report stays on this topic. |
| `draft_report` | string (markdown) | `stages.synthesize-report.draft_report` | The unreviewed report. Citation markers refer to entries in `cited_sources`. |
| `cited_sources` | array of `{id, title, url, …}` | `stages.synthesize-report.cited_sources` | The sources the draft actually cited. Each has a stable `id`. |
| `findings` | array of `{claim, evidence, source_id, confidence, …}` | `stages.extract-findings.findings` | The full evidence pool from extract-findings. `source_id` references entries in the upstream `sources` set. The reviewer's closed-world citation universe. |
| `audience` | string | workflow input | Who the report is written for (e.g. "general", "engineering manager", "researchers"). The reviewer checks fit, not invents new audience. |
