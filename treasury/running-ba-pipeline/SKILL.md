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
Run this map on a **small–mid** model (e.g. Haiku 4.5 / Gemini 3 Flash), low reasoning effort — it only sequences and bookkeeps. Per-node model & cost routing: `references/pipeline-contracts.md`.

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

## Context budget, shared envelope, gates & safety (shared contracts)
Each node declares the files it reads and reads no more, and every handoff carries one typed envelope; the gate matrix
and command-safety tiers govern what an agent may do unattended. These cross-stage tables are shared by all five child
skills — the authoritative copy lives in the references below, so a node and the map never drift.
- Shared envelope + per-node context-budget reads + model/cost routing: `references/pipeline-contracts.md`
- Gate matrix + command-safety policy (ALLOW/CONFIRM/DENY): `references/gates-and-safety.md`

## Never-do guardrails
1. Never echo real PII. Redact to `<PII:REDACTED:CLASS=...>`; personal data stays out of every artifact and log.
2. Never silently repair a defect. Surface ambiguities and gaps as open questions or blockers.
3. Never auto-resolve a blocker. Only a named human clears a Governance Check blocker.
4. Never hand off while a blocker is open. `ready-for-tl` is impossible with an unresolved blocker.
5. Never let an agent relax a gate or approve its own output.
6. Never conflate Compliance (describes the rule) with Legal (interprets the wording) — different owners.

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
