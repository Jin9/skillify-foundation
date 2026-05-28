---
name: running-ba-pipeline
description: >
  Run the five-stage agentic Business-Analysis pipeline that turns a raw requirement into a Tech-Lead-ready handoff: it carries one task_id and a typed envelope across intake, story drafting, governance, feasibility, and TL handoff, placing a human approval gate after each stage and looping back to a named human when governance is blocked. Use when the user asks to "run the BA pipeline", "take this requirement from intake to TL handoff", "which BA stage am I in or what gate is next", "wire these BA stages together", or "enforce the BA pipeline gates and guardrails". Do NOT use to design a brand-new agent pipeline or move its gates (use agentic-workflow-design); do NOT use to produce a single stage's artifact in isolation (use scoping-ba-intake, drafting-ba-stories, checking-ba-governance, assessing-ba-feasibility, or assembling-tl-handoff).
compatibility: claude-code, codex, gemini-cli, opencode
---

# running-ba-pipeline

## Purpose
Run a raw business requirement through five decoupled BA stages to a Tech-Lead-ready handoff, with a human approval
gate after every stage. Autonomy is **assist-level**: each stage agent drafts a typed contract; a named human owns
the decision at each gate. This skill is the **map** — it sequences the stages and holds the shared rules. Each stage
is its own skill and also stands alone; deleting this map would not break any single stage.

## When to use
- Triggers: *"run the BA pipeline"*, *"take this requirement from intake to TL handoff"*, *"which stage or gate am I at"*, *"wire these BA stages together"*, *"enforce the pipeline gates and guardrails"*.
- **Not this skill:** designing a brand-new pipeline or moving its gates → `agentic-workflow-design`; producing one stage's artifact alone → that stage's skill (`scoping-ba-intake`, `drafting-ba-stories`, `checking-ba-governance`, `assessing-ba-feasibility`, `assembling-tl-handoff`).

## Model & cost
Run this map on a **small–mid** model (e.g. Haiku 4.5 / Gemini 3 Flash), low reasoning effort — it only sequences and bookkeeps. Per-node tiers are in **Model & cost routing** below.

## Pipeline at a glance
Each node is `input → agent → output contract → gate → owner`. One agent per stage; the flow is linear with one human loop-back.

| Stage | Agent | Input | Output contract | Gate | Owner |
|-------|-------|-------|-----------------|------|-------|
| S1 Intake & Scope | `intake-agent` | raw requirement | Scope Sheet | G1 confirm scope (sync) | BA / PM |
| S2 Story Drafting | `story-agent` | Scope Sheet | Story Set (`02-story-set/` tree) | G2 accept (async) | BA |
| S3 Governance & Risk | `governance-agent` | Story Set INDEX + selected stories | Governance Check (`clear`/`blocked`) | G3 resolve blocker (sync) | BA + named SME |
| S4 Feasibility & Scope | `feasibility-agent` | Story Set INDEX + selected stories + Governance Check + raw req | Feasibility Note (`buildable`/`not-now`/`phased`) | G4 verdict (sync) | Tech Lead |
| S5 TL Handoff | composition | all of the above, by path | Handoff Bundle | G5 accept (sync) | Tech Lead |

```
raw req -> S1 -> G1 -> S2 -> G2 -> S3 -> G3 --clear--> S4 -> G4 -> S5 -> G5 -> TL owns it
                                  '- blocked -> named human resolves -> re-run the affected stage
```

## Output layout
Artifacts are written under one run directory, numbered by stage. The Story Set is a **directory**, not a single file,
so no node ever loads every story at once:

```
00-run-log.md
01-scope-sheet.md
02-story-set/
  INDEX.md                 # envelope + epic + the enumerated story table — the manifest downstream nodes read first
  ST-01-<slug>.md          # one full user-story-template instance per story
  ST-02-<slug>.md
  ...
03-governance-check.md     # (+ 03a resolution and 03b re-run when G3 fires)
04-feasibility-note.md
05-handoff-bundle.md       # links the artifacts above by path; inlines only open_items + invariant assertions
```

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

## Gate matrix
Gates sit where an action is hard to reverse or has high blast radius — not uniformly.

| Gate | After | Action gated | Reversibility | Blast | Type | Approver |
|------|-------|--------------|---------------|-------|------|----------|
| G1 | S1 | confirm scope | reversible draft | anchors all downstream | sync named | BA / PM |
| G2 | S2 | accept Story Set | reversible | low | async review | BA |
| G3 | S3 | resolve a governance blocker | hard to reverse | high (legal/privacy) | sync named | BA + SME (Legal/DPO/Compliance) |
| G4 | S4 | feasibility verdict | sets scope & cost | medium–high | sync named | Tech Lead |
| G5 | S5 | accept handoff bundle | engineering builds on it | high | sync named | Tech Lead (owner of record) |

An agent never passes an irreversible or control-plane gate on its own confidence; confidence only tunes whether a *reversible* step gets async review.

## Command-safety policy (enforced at the tool layer, not by the prompt)
| Tier | Actions |
|------|---------|
| ALLOW | read the raw requirement; classify/parse it; draft any stage contract; render a diagram; write inside `output/` |
| CONFIRM | publish the Handoff Bundle; mark a blocker resolved; set `state: ready-for-tl`; write to a shared backlog; override an upstream contract's scope |
| DENY | echo real PII; auto-resolve a governance blocker; write outside `output/`; call a non-allowlisted tool; modify another node's contract or its own permissions |

## Never-do guardrails
1. Never echo real PII. Redact to `<PII:REDACTED:CLASS=...>`; personal data stays out of every artifact and log.
2. Never silently repair a defect. Surface ambiguities and gaps as open questions or blockers.
3. Never auto-resolve a blocker. Only a named human clears a Governance Check blocker.
4. Never hand off while a blocker is open. `ready-for-tl` is impossible with an unresolved blocker.
5. Never let an agent relax a gate or approve its own output.
6. Never conflate Compliance (describes the rule) with Legal (interprets the wording) — different owners.

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

## How to run it
1. Assign one `task_id` (the requirement id) — it threads every stage.
2. Run S1 `scoping-ba-intake` → Scope Sheet. **G1:** BA/PM confirm scope and answer blocking open questions.
3. Run S2 `drafting-ba-stories` → Story Set as the `02-story-set/` tree (INDEX + one file per story). **G2:** async BA review (open the INDEX, then individual stories) for INVEST and testable acceptance criteria.
4. Run S3 `checking-ba-governance` (reads the INDEX + the stories it flags) → Governance Check. **G3:** if `blocked`, go to Loop-back; if `clear`, continue.
5. Run S4 `assessing-ba-feasibility` with the Story Set (INDEX + the stories it weighs), the *clear* Governance Check, and the raw requirement → Feasibility Note. **G4:** TL verdict.
6. Run S5 `assembling-tl-handoff` → Handoff Bundle (links artifacts by path; the `story_set` link is `02-story-set/INDEX.md`). **G5:** TL accepts and becomes owner of record.
- Re-running a stage with the same `task_id` reproduces the same contract (idempotent — safe to replay).

## Loop-back & failure
- **Blocked governance:** a named human (Legal / DPO / Compliance / SME / PM) resolves each blocker; the affected upstream stage re-runs with the enriched input. Nothing advances while a blocker is open.
- **Input too thin:** S1 returns `needs-clarification` to the requester — do not invent scope.
- **Transient draft/schema error:** retry with a small hard cap, then stop (do not loop).
- **HITL is the final tier:** when a stage cannot proceed safely it stops and asks a human; it never forces a path through a gate.

## Checklist
- [ ] One `task_id` threads S1–S5
- [ ] Every handoff carries the shared envelope (`schemaVersion` set)
- [ ] The Story Set is the `02-story-set/` tree (INDEX + one file per story); no node loads every story at once
- [ ] Exactly one owner mutates each contract; downstream nodes are read-only
- [ ] No stage advanced past an open blocker
- [ ] The raw requirement travels in the S5 bundle (by path), not just the BA contracts

## Human approval gate
**Stop at every gate.** BA / PM own G1; BA owns G2; BA + named SME own G3; the Tech Lead owns G4 and G5. On G5 the
Tech Lead becomes the accountable owner of record for everything engineering builds from the bundle.

## Example (ShopPilot MVP, `task_id REQ-shoppilot-mvp`)
G1 ✅ scope confirmed → G2 ✅ stories accepted → S3 **`blocked`** (seven governance gaps) → **G3:** named humans
(Legal / DPO / Compliance / Finance / Ops) resolve each → affected stages re-run → S3 re-run **`clear`** → S4
**`buildable`** → **G5:** TL accepts. One trace id throughout; nothing reached the TL while blocked; the raw spec
travelled all the way to the TL.
