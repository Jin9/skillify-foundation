# Pipeline contracts: shared envelope, context-budget reads, model & cost routing

These cross-stage contract tables are shared by all five child skills (`scoping-ba-intake`,
`drafting-ba-stories`, `checking-ba-governance`, `assessing-ba-feasibility`,
`assembling-tl-handoff`). They are the authoritative copy — the `SKILL.md` body points here.

## Shared envelope (every handoff carries it)
| Field | Type | Meaning |
|-------|------|---------|
| `task_id` | string | the requirement id — single trace id across S1–S5 |
| `stage` | `S1..S5` | which stage produced this contract |
| `intent` | string | what the stage was asked to produce |
| `state` | string | the stage verdict / output type (see each contract) |
| `confidence` | `high`/`medium`/`low` | agent self-rating; `low` forces async human review |
| `provenance` | object | `{ raw_ref, upstream_contract_ref }` — audit back-pointers |
| `produced_by` | string | distinct agent identity (e.g. `story-agent`) |
| `owner` | string | named human owner of record for this stage |
| `schemaVersion` | string | contract version, for drift detection |
| `created_at` | RFC-3339 | timestamp |

## Context-budget rule (load only what you need)
Each node declares the files it reads and reads no more. The Story Set INDEX is always cheap; individual stories are
pulled only when a node needs that story's depth — so a 200k window (and a mid-tier model) is never blown by the whole
set.

| Node | Reads | Never loads |
|------|-------|-------------|
| `drafting-ba-stories` | the Scope Sheet | — (it writes the `02-story-set/` tree) |
| `checking-ba-governance` | `02-story-set/INDEX.md` + the specific `ST-NN-*.md` it flags | every story file at once |
| `assessing-ba-feasibility` | `INDEX.md` + the stories it weighs + clear Governance Check + raw requirement | every story file at once |
| `assembling-tl-handoff` | contract headers + `INDEX.md` | every story file; full contract bodies |

## Model & cost routing
Nodes are decoupled, so assign a model per node instead of running the heavy model everywhere. Tiers are capability levels (vendor-neutral) with dated examples; for a full per-role policy see `model-selection`.

| Node | Tier (example, 2026) | Effort | Why |
|------|----------------------|--------|-----|
| `running-ba-pipeline` | small–mid (Haiku 4.5 / Gemini 3 Flash) | low | orchestration only |
| `scoping-ba-intake` | mid (Sonnet 4.6 / Gemini 3 Pro) | medium | structuring, no deep trade-offs |
| `drafting-ba-stories` | mid → frontier | medium–high | mid is fine; frontier sharpens acceptance criteria |
| `checking-ba-governance` | mid | medium | checklist sweep; escalate to frontier only if compliance wording is ambiguous |
| `assessing-ba-feasibility` | frontier (Opus 4.7 / top Gemini) | high | the one node that genuinely wants the heavy model |
| `assembling-tl-handoff` | small (Haiku 4.5 / Gemini 3 Flash) | low | deterministic assembly + boolean checks |

Model tier never changes who approves: a bigger model never earns the skipping of a human gate.
