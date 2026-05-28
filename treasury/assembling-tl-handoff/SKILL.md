---
name: assembling-tl-handoff
description: >
  Stage 5 of the BA pipeline: assemble the Handoff Bundle the Tech Lead signs for — the raw requirement plus the Scope Sheet, Story Set, clear Governance Check, and buildable-or-phased Feasibility Note linked by path, plus any explicitly accepted open items — and set state ready-for-tl only when governance is clear and no blocker is open, for the TL acceptance gate (G5). It is composition, not an agent: it references artifacts by path rather than concatenating them, and it carries the raw requirement all the way to the TL. Use when the user is running the BA pipeline and asks to "assemble the handoff", "produce the Handoff Bundle", "is this ready for the TL", or "package this for Tech-Lead acceptance". Do NOT use to design downstream architecture or build anything; do NOT set ready-for-tl while any blocker is open.
compatibility: claude-code, codex, gemini-cli, opencode
---

# assembling-tl-handoff

## Purpose
Stage 5 (composition, not an agent) of the BA pipeline. Assemble the **Handoff Bundle** the Tech Lead signs for — the
raw requirement plus the four stage contracts (linked by path) and any accepted open items — and set
`state: ready-for-tl` only when the invariant holds. The raw requirement travels all the way to the TL, not just the
BA contracts.

## When to use
- Triggers: *"assemble the handoff"*, *"produce the Handoff Bundle"*, *"is this ready for the TL"*, *"package this for Tech-Lead acceptance"*.
- **Not this skill:** designing the downstream architecture or building anything — that is the TL's job after acceptance.

## Model & cost
Runs on a **small** model (e.g. Haiku 4.5 / Gemini 3 Flash), low effort — deterministic assembly plus boolean
invariant checks. Because it links artifacts by path instead of inlining them, it stays within budget even on a small
model: it reads the contract headers + the Story Set INDEX, never every story file.

## Input
- All four upstream artifacts **by path** — **Scope Sheet** (`01-scope-sheet.md`), **Story Set**
  (`02-story-set/INDEX.md` — the index, not every story file), **Governance Check** (`03b-…` re-run, `clear`),
  **Feasibility Note** (`04-…`, `buildable` or `phased`) — and the **raw requirement**. Composition reads the INDEX +
  contract headers; it never loads every story file.

## Output
A **Handoff Bundle** — the shared envelope plus a *manifest* that REFERENCES each upstream artifact by relative path.
The bundle does **not** concatenate full contracts (that would blow the context budget and duplicate the source of
truth); it links them and inlines **only** `open_items` and the invariant assertions.

```
# Handoff Bundle — <task_id>
## envelope
task_id: REQ-<slug>            # identical across all five contracts
stage: S5
intent: assemble the TL handoff
state: ready-for-tl            # ONLY when governance is clear AND no blocker is open
confidence: high | medium | low
provenance: { raw_ref: <raw requirement>, upstream_contract_ref: [ scope_sheet@1.0, story_set@1.0, governance_check@1.0, feasibility_note@1.0 ] }
produced_by: composition
owner: Tech Lead
schemaVersion: 1.0
created_at: <RFC-3339>
## artifact manifest (reference by path — do NOT inline full contracts)
| artifact         | path                            | state / verdict |
|------------------|---------------------------------|-----------------|
| raw_requirement  | <relative path to original ask> | travels all the way here |
| scope_sheet      | ./01-scope-sheet.md             | ready-for-stories |
| story_set        | ./02-story-set/INDEX.md         | drafted (epic + N stories — link the INDEX, not every file) |
| governance_check | ./03b-governance-check-rerun.md | clear (see 03 + 03a for the blocked→resolve arc) |
| feasibility_note | ./04-feasibility-note.md        | buildable | phased |
## open_items (inlined — the only deferred content the TL must see at a glance)
| item                | source_oq | accepted_by  |
|---------------------|-----------|--------------|
| <deferred question> | OQ-<n>    | PM | Sponsor  |     # source_oq MUST be a real Scope-Sheet OQ id; never renumbered
## invariant assertions (all must pass)
- [ ] state: ready-for-tl ⟹ governance_check.verdict == clear AND no open blocker
- [ ] one task_id identical across all five contracts
- [ ] raw_requirement present (not a "BA-contracts-only" handoff)
- [ ] feasibility_note.verdict is buildable | phased (not not-now)
- [ ] every open_items.source_oq exists in the Scope Sheet's open_questions (no invented OQ ids)
- [ ] every open_items entry is PM/Sponsor-accepted
accepted_by: <Tech Lead name + date>   # the G5 sign-off, left unsigned for the TL
```

## Decision rules (invariant assertions — refuse if any fails)
1. `state: ready-for-tl` **only if** `governance_check.verdict == clear` **and** no blocker is open. Otherwise stop and report which invariant failed.
2. **Reference every artifact by relative path** — never paste a full contract body into the bundle; the `story_set` link points at `02-story-set/INDEX.md`, not the individual story files. Only `open_items` and the invariant assertions are inlined.
3. The `task_id` must be **identical** across all five contracts (the one trace id).
4. `raw_requirement` must be **present** — never hand off "BA contracts only".
5. `feasibility_note.verdict` must be `buildable` or `phased` (not `not-now`).
6. A deferred question goes in `open_items` **only if** a PM or Sponsor accepted it; otherwise it is still a blocker.
7. **Every `open_items.source_oq` is a real Scope-Sheet `OQ-n`** — never renumber or invent an id; this keeps each open item traceable to where the question was first raised.
8. **Redact any real PII** to `<PII:REDACTED:CLASS=...>`.
9. Setting `ready-for-tl` and publishing the bundle is a CONFIRM action — it needs the named human at G5; the agent never self-accepts.

## Checklist
- [ ] All four artifacts linked **by path** + the raw requirement present (no full contract bodies inlined)
- [ ] `story_set` link points at `02-story-set/INDEX.md`, not individual story files
- [ ] One `task_id` across all five
- [ ] Governance `clear`, no open blocker
- [ ] Feasibility verdict is `buildable` or `phased`
- [ ] Every `open_items.source_oq` exists in the Scope Sheet; all are PM/Sponsor-accepted
- [ ] `state: ready-for-tl`; `accepted_by` left for the TL to sign; no real PII

## Anti-patterns (never do)
- Set `ready-for-tl` with an open blocker or a non-`clear` governance check.
- Concatenate full upstream contracts into the bundle instead of linking them by path (bloats context, duplicates the source of truth).
- Drop the raw requirement and hand off only the distilled contracts.
- Self-sign `accepted_by` on the agent's behalf.
- Fold an unaccepted open question into `open_items`, or give it an invented / renumbered `source_oq`.

## Example (ShopPilot MVP)
The bundle links `raw_requirement` (the full business spec), `./01-scope-sheet.md`, `./02-story-set/INDEX.md`,
`./03b-governance-check-rerun.md` (`clear`) and `./04-feasibility-note.md` (`buildable`) by path — none inlined — and
inlines three PM/Sponsor-accepted `open_items`, each tagged with its real Scope-Sheet id: cancelled-order review
eligibility (`OQ-1`), partial-refund policy (`OQ-2`), dormant-product archival (`OQ-3`). `state: ready-for-tl`;
`accepted_by` awaits the TL's name + date.

## Human approval gate (G5)
**Stop.** The Tech Lead accepts the bundle (sync, named) and becomes the accountable owner of record for everything
engineering builds from it. The agent assembles and asserts the invariant; it never accepts on the TL's behalf.
