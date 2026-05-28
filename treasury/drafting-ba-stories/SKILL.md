---
name: drafting-ba-stories
description: >
  Stage 2 of the BA pipeline: consume a confirmed Scope Sheet and emit a typed Story Set — one epic plus INVEST user stories, each a full user-story-template instance (Description, Business Logic with decision tables, Gherkin acceptance criteria, Out of Scope), emitted as one file per story under an enumerated INDEX so a mid-tier model loads one story at a time — wrapped in the shared envelope, for the async story-review gate (G2). Use when the user is running the BA pipeline and asks to "draft the stories", "produce the Story Set", "write the epic and stories for the pipeline", or "turn the scope into INVEST stories". Do NOT use for heavy banking-grade elicitation with force-filled governance fields (use eliciting-banking-brief); do NOT scope the requirement here (use scoping-ba-intake).
compatibility: claude-code, codex, gemini-cli, opencode
---

# drafting-ba-stories

## Purpose
Stage 2 (`story-agent`) of the BA pipeline. Turn a confirmed **Scope Sheet** into a typed **Story Set** — one epic and
a set of INVEST user stories, **each a full instance of the `user-story-template`** (the Jira-style ticket shape:
Description → Business Logic → Acceptance Criteria → Out of Scope) — wrapped in the shared envelope, ready for Gate G2.
The Story Set is emitted as a **directory** (one file per story plus a thin INDEX), so no reader ever loads every story
at once. This is a black-box node: it consumes only the Scope Sheet contract and emits one contract.

## When to use
- Triggers: *"draft the stories"*, *"produce the Story Set"*, *"write the epic and stories for the pipeline"*, *"turn the scope into INVEST stories"*.
- **Not this skill:** heavy banking-grade elicitation with force-filled governance fields → `eliciting-banking-brief`; scoping the requirement → `scoping-ba-intake`; standalone PII / compliance flagging → `checking-ba-governance` (S3).

## Model & cost
Runs on a **mid** model (e.g. Sonnet 4.6 / Gemini 3 Pro), medium–high effort; a **frontier** model sharpens acceptance
criteria. Either tier produces the same contract shape. Because the output is one-file-per-story, a mid model drafts —
and a reviewer reads — a single story within budget; it never needs the whole set in context at once.

## Input
- A confirmed **Scope Sheet** (S1 output, `state: ready-for-stories`).

## Output
A **Story Set** emitted as a *directory*, not one file, so downstream nodes and a mid-tier model load selectively and
no node ever holds every story at once:

```
02-story-set/
  INDEX.md                 # shared envelope + epic + the enumerated story table — load me first
  ST-01-<slug>.md          # one full user-story-template instance per story
  ST-02-<slug>.md
  ...
```

The **shared envelope lives once, in `INDEX.md`** (one envelope per Story Set, not per story). Each `ST-NN-<slug>.md`
is a full instance of `templates/user-story-template.md` (or `templates/user-story-template-th.md`; see *Language*),
prefixed by a 4-line HTML-comment locator stub so a story file is self-locating when opened on its own:

```
<!-- ba-pipeline: story-set REQ-<slug> · ST-NN · state: drafted | needs-scope-rework
     scope_ref: ../01-scope-sheet.md · open_questions: [OQ-n, ...]  (carried from the Scope Sheet, never invented) -->
```

### INDEX.md
```
# Story Set — <task_id>
## envelope
task_id: REQ-<slug>            # same trace id as the Scope Sheet
stage: S2
intent: draft epic + stories
state: drafted | needs-scope-rework
confidence: high | medium | low
provenance: { raw_ref: <source>, upstream_contract_ref: scope_sheet@1.0 }
produced_by: story-agent
owner: BA lead
schemaVersion: 1.0
created_at: <RFC-3339>
## epic
title:   <epic title>
summary: <2–3 sentence epic summary>
## stories
| ST-ID | issue-key  | title                        | summary (1 line)      | state   | file              |
|-------|------------|------------------------------|-----------------------|---------|-------------------|
| ST-01 | <KEY/none> | [<TYPE>] <imperative title>  | <one-line capability> | drafted | ST-01-<slug>.md   |
| ...   | ...        | ...                          | ...                   | ...     | ...               |
```
`issue-key` is the tracker key if one exists, else `<none>`. This table is the manifest downstream nodes read first.

### Per-story file = a full `user-story-template` instance
Each `ST-NN-<slug>.md` follows the vendored template exactly. Where each section's content comes from:

| Template section | Required? | Source in the Scope Sheet (S1) |
|---|---|---|
| breadcrumb `< Back · <ISSUE-KEY> · <RELATED-KEY>` | optional | tracker key(s) if any, else drop the line |
| `# [<TYPE>] <title>` | **required** | one capability from `in_scope`; tag `[Feature]` / `[BAU]` / `[Spike]` |
| `## Description` (As a / I want / So that) | **required** | the capability + its value (`business_goal`) |
| `## Product` | optional | product / feature area named in scope |
| `## Background` | optional | `business_goal` + the trigger for this capability |
| `## Business Logic` (numbered rules, a./b. sub-items) | optional | rules the scope says the system must enforce |
| `### <Decision matrix>` table | optional | **only when the scope gives conditional logic** (rate / limit / state tables) — never invent one |
| `> ⚠️ Note:` callout | optional | a cross-cutting rule / effective date / ownership |
| `## Acceptance Criteria` | **kept** (pipeline override) | ≥1 testable Gherkin Given/When/Then with concrete values — the G2 gate requires it |
| `## Out of Scope` | optional | this story's slice of `out_of_scope`; name the owning team |
| `## Mock-up references` (AS-IS / TO-BE, Spec API) | optional | only if the scope references screens / APIs; else delete the section |
| METADATA (Attachments / Subtasks / Linked / Confluence / Activity) | optional | `<none>` placeholders unless mirroring a real ticket export |

Per the template, **only Title and Description are required; keep the sections you need and delete the rest** — except
Acceptance Criteria, which this pipeline keeps (the G2 gate is about testable criteria). Banking / compliance rigor is
**not** force-filled here: a story carries rigor through Business-Logic rules + decision tables + Acceptance Criteria;
standalone PII / compliance flags live in `checking-ba-governance` (S3). When the scope supports it a story reaches the
richness of `references/example_TierRateStory.md` (numbered Business-Logic rules, a `### <name>` decision table, a
`> ⚠️ Note`). When the scope gives no conditional logic, omit the table — do not invent one.

### Deterministic naming (re-run-idempotent)
1. **ST-NN** — number stories densely from `01`, zero-padded, in the order they are derived from `in_scope` (top to
   bottom). Same Scope Sheet ⇒ same order ⇒ same numbers.
2. **slug** — the title lowercased, with the `[<TYPE>]` tag and punctuation dropped, spaces → `-`, repeats collapsed,
   **ASCII word-characters only**, capped at 8 words (e.g. `[Feature] Checkout with server-side totals` →
   `ST-03-checkout-with-server-side-totals.md`). If a title has no ASCII word-characters, the file is just `ST-NN.md`.
3. Filenames are pure functions of `(ST-NN, title)` — no timestamps, no counters, no randomness — so re-running on the
   same Scope Sheet reproduces identical filenames (idempotent; safe to replay).

### Consistency invariant (check before emitting)
`epic`'s story count **==** number of `ST-NN-*.md` files **==** number of INDEX `## stories` rows. If they disagree,
the Story Set is malformed — fix it before advancing to G2.

### Language
Story **content** follows the source language (the template states content may be written in any language). Pick the
scaffold to match: English content → `templates/user-story-template.md`; Thai content →
`templates/user-story-template-th.md`. The exemplar is Thai content on the English scaffold — that is allowed; match
the source. Do **not** emit parallel EN+TH copies of a story.

## Decision rules
1. **One capability per story** (INVEST-Small) — split a story that hides two capabilities.
2. **Every story is a full `user-story-template` instance** — never hand-write a divergent shape; fill Title +
   Description always, and the other sections only when the Scope Sheet supports them.
3. **Acceptance criteria are concrete** Gherkin given/when/then with real values (amounts, states, messages), so they
   are testable; keep ≥1 per story.
4. **Emit one file per story plus a thin INDEX** — never dump all stories into a single file (it blows the context
   budget for a mid-tier model). Name files deterministically (see *Deterministic naming*).
5. **Add a `### <decision matrix>` table only when the scope gives conditional logic** — never invent one to look rich.
6. **Carry unresolved items forward** in the per-story `open_questions` stub using the Scope Sheet's exact `OQ-n` id;
   never invent an answer the Scope Sheet did not settle, and never renumber an open question.
7. Each story's **Out of Scope** names the owning team when the boundary lands elsewhere.
8. Keep every story **traceable to the Scope Sheet** — do not add capability the scope did not authorize.
9. If the scope is too thin to split into stories, set `state: needs-scope-rework` and loop back to `scoping-ba-intake`.
10. **Redact any real PII** to `<PII:REDACTED:CLASS=...>`; use synthetic values in examples.

## Checklist
- [ ] `INDEX.md` carries the envelope + epic (title + summary) + the story table
- [ ] One file per story (`ST-NN-<slug>.md`), each a full `user-story-template` instance
- [ ] Count invariant holds: epic story count == ST-NN files == INDEX rows
- [ ] Each story is INVEST-shaped (one capability) with Title + Description and ≥1 Gherkin acceptance criterion
- [ ] Decision tables appear only where the scope has conditional logic (none invented)
- [ ] Unresolved items carried in the per-story stub with the Scope Sheet's exact `OQ-n` id (not invented, not renumbered)
- [ ] Deterministic filenames (no timestamps/randomness); correct language scaffold; no real PII

## Anti-patterns (never do)
- Hand-write a terse, divergent story shape instead of a full `user-story-template` instance.
- Dump every story into one file — breaks the per-story context budget and the enumerated layout.
- Invent a decision-matrix table (or any section) the Scope Sheet does not support, to look thorough.
- Vague acceptance criteria ("works correctly") with no observable outcome; a "god story" bundling capabilities.
- Invent answers to open questions, or renumber a Scope Sheet `OQ-n`.

## Example (ShopPilot MVP)
`INDEX.md` carries the envelope + epic, then a story table whose first rows are:

```
| ST-ID | issue-key | title                                 | summary (1 line)                     | state   | file                                       |
|-------|-----------|---------------------------------------|--------------------------------------|---------|--------------------------------------------|
| ST-03 | <none>    | [Feature] Checkout with server-side totals | Pay exactly the server-computed total | drafted | ST-03-checkout-with-server-side-totals.md  |
| ST-05 | <none>    | [Feature] Stock reservation & last-item race | One winner on the last unit, at confirm | drafted | ST-05-stock-reservation-last-item-race.md |
```

`ST-03-checkout-with-server-side-totals.md` (abbreviated — a full `user-story-template` instance):

```
<!-- ba-pipeline: story-set REQ-shoppilot-mvp · ST-03 · state: drafted
     scope_ref: ../01-scope-sheet.md · open_questions: [] -->

`< Back`  ·  `<none>`

# [Feature] Checkout with server-side totals

## Description
*As a* logged-in customer
*I want* to review cart, address, coupon, shipping and net total before paying
*So that* I pay exactly what the system computed, not a tampered client value.

## Business Logic
1. Totals are computed **server-side** (§4.6); the client value is never trusted.
2. Shipping is free at ≥ 1,500 THB, else 60 THB (§6.3).
3. One coupon per order; the total is floored at 0 and never negative (§6.4).
4. The coupon is **re-validated at confirm** (§4.5).

## Acceptance Criteria
1. Given a 1,200 THB cart with coupon `WELCOME100`, when I confirm, then the total is 1,160 THB (1,200 − 100 + 60 shipping) and an order is created in `awaiting-payment`.
2. Given a 1,600 THB cart, when I confirm, then shipping is 0 (free ≥ 1,500).
3. Given a coupon that expired while I browsed, when I confirm, then the order is rejected with "coupon no longer valid" and no order is created.

## Out of Scope
1. Coupon stacking (Phase 2).
```

`ST-05` carries `open_questions: [OQ-6]` (the reservation-TTL question) in its locator stub — carried verbatim from
the Scope Sheet, not invented here.

## Human approval gate (G2)
**Stop.** A BA reviews the Story Set asynchronously — opening `INDEX.md`, then individual story files — for INVEST-ness
and testable acceptance criteria before it goes to `checking-ba-governance`. Low `confidence` forces this review; the
artifact is reversible and low-blast, so the gate stays async.
