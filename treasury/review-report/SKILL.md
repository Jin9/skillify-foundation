---
name: review-report
description: >
  Self-review pass over a draft research report. Verifies every cited claim is
  grounded in the input findings, checks audience and topic fit, surfaces
  internal-consistency issues, and may perform bounded citation surgery (drop,
  add, swap) using ONLY sources already present in the input findings. Emits
  final_report, updated cited_sources, and review_notes. Single-pass, no
  re-search. Use when the workflow stage is `review-report`, or when the user
  says "review this draft report", "self-review my report", "audit citations
  against findings". Do NOT use to write the draft (use `synthesize-report`)
  or to re-run search.
compatibility: claude-code, codex, copilot, gemini, antigravity
allowed-tools: ""  # no external tools — closed-world citation surgery against in-context findings
version: 0.2.2
stage: review-report
workflow: researcher
inputs:
  - { name: topic,         type: string, required: true, source: workflow.inputs.topic,                  description: "the subject" }
  - { name: draft_report,  type: string, required: true, source: stages.synthesize-report.draft_report,  description: "the unreviewed draft" }
  - { name: cited_sources, type: array,  required: true, source: stages.synthesize-report.cited_sources, description: "draft's citation set; may be modified" }
  - { name: findings,      type: array,  required: true, source: stages.extract-findings.findings,       description: "ground truth for the review" }
  - { name: audience,      type: string, required: true, source: workflow.inputs.audience,               description: "for fit checks" }
outputs:
  - { name: final_report,  type: string, description: "post-review markdown; may equal draft if no changes" }
  - { name: cited_sources, type: array,  description: "post-review citation set; may differ from input" }
  - { name: review_notes,  type: string, description: "what was flagged, changed, or left unresolvable" }
---

# Review Report

## Purpose

Run one disciplined self-review pass over the draft research report produced by
`synthesize-report`. Verify every cited claim against the input findings, check
the prose fits the topic and the stated audience, surface internal-consistency
issues, and apply only **bounded** revisions — minimum-edit, never destructive.
Emit the post-review report, an updated citation set (which may differ from the
input), and review notes that downstream readers can audit.

This is **stage 5 of 6** in the `researcher` workflow (the meta `reporting-research-run`
stage runs after). It is **single-pass**: no re-search, no recursive
re-review. Unresolvable gaps are recorded, not fixed.

## When to use this skill

- Use when: the workflow's prior stage is `synthesize-report` and the draft is
  ready for the final review pass before user delivery.
- Use when: the user says "review this draft", "self-review my report",
  "audit citations against findings".
- Do NOT use when: the request is to *write* the report — that's
  `synthesize-report`.
- Do NOT use when: the reviewer wants to *re-search* — that is explicitly
  outside this skill's scope; the only sources available are those in the
  input `findings` array.
- Do NOT use when: the request is a security or code review (different skills).
- Do NOT use when: the report needs to be re-targeted to a different audience —
  emit a `review_notes` entry and stop; the workflow operator may re-run from
  `synthesize-report` with a new audience.

## Input contract

The skill MUST receive five inputs verbatim from the workflow engine — `topic`,
`draft_report`, `cited_sources`, `findings`, `audience`. No input may be
inferred or fabricated. If any required input is missing, malformed, or empty,
emit `final_report = draft_report`, `cited_sources = input cited_sources`, and a
`review_notes` entry of category `input_incomplete` describing the blocker.

Field-by-field catalog (types, sources, notes):
[`references/input-contract.md`](references/input-contract.md).

**Closed-world rule:** any source the reviewer keeps, swaps in, or adds MUST
trace to either `cited_sources` (already present) or a `findings[*].source_id`
(promotable). The reviewer never fabricates new sources.

## Procedure

Run all six steps in order. Steps 1–4 collect observations. Step 5 decides
what to apply vs. flag. Step 6 emits the structured output. Do not start
rewriting prose before step 5 — collecting first prevents destructive,
mid-stream revision.

### Step 1 — Build the claim → citation map

Walk the `draft_report`. For every citation marker, record the sentence (or
smallest claim unit) it grounds, the cited `source_id`, and the section heading
it sits under. Also list **uncited claims** that read as factual assertions
(numbers, named results, quoted positions, recommendations) — candidates for
grounding in Step 3. Soft framings ("in general," "broadly speaking,"
introductions, transitions) are not assertions and need no grounding.

**Fenced code blocks are excluded from this scan.** ASCII diagrams emitted by
`synthesize-report` (v0.2.0+) are structural renderings of claims the prose
already makes, not new assertions. Treat each block as opaque — do not extract
claims from it, do not flag missing `[n]` markers inside it. Diagram-vs-prose
consistency is audited in Step 4.

### Step 2 — Verify each cited claim against findings

For each cited claim from step 1, look up the `findings` entry whose `source_id`
matches the citation:

- **Match, evidence supports the claim** → `verified`.
- **Match, evidence does NOT support** → `unsupported_by_cited_source`
  (candidate *swap* or *drop*).
- **No matching finding for that source_id** → `citation_broken` (candidate
  *swap* or *drop*).

For `unsupported_by_cited_source` and `citation_broken`, scan the rest of
`findings` for a better source: exactly one higher-or-equal-confidence finding
supporting the claim → candidate *swap*; none → candidate *drop* (reword/remove
the claim). Record every check, including verified ones — a non-empty verified
ledger is the antidote to rubber-stamping.

### Step 3 — Find uncited claims that should be cited

For each uncited factual assertion from step 1, scan `findings` for a claim
whose `evidence` supports it:

- **Exactly one supporting finding** → candidate *add* (cite that source).
- **Multiple findings agree** → candidate *add* using the highest-confidence
  finding's source.
- **No supporting finding** → `unresolvable_gap`: asserted but ungroundable
  with available evidence. Flag for the Step 5 decision (hedge, remove, or
  accept as a flagged gap).

### Step 4 — Topic, audience, and internal consistency

Four lightweight scans, each yielding zero or more observations:

- **Topic fit:** sections that drift wildly off `topic` → `off_topic`.
- **Audience fit:** check terminology and assumed background against `audience`.
  *Minor* (one place of jargon for a "general" audience) → flag-only, no
  rewrite. *Pervasive* (whole report assumes deep domain knowledge) →
  top-level `review_notes` entry of category `audience_mismatch_pervasive`; do
  NOT attempt a full audience rewrite — flag and let the operator decide.
- **Internal consistency:** check numeric, date, and entity-name agreement plus
  claim→conclusion alignment. Contradictions → `internal_contradiction`; prefer
  the higher-confidence side when resolvable, flag-only when not.
- **Diagram consistency:** for each fenced block, verify its structure matches
  the surrounding prose (same entities, arrows/hierarchy, flow direction), sits
  in a permitted section (themed body or `## Synthesis`), and contains no `[n]`
  marker. Mismatches → `diagram_inconsistent`. Do NOT fix a diagram's contents
  — see Step 5.

### Step 5 — Decide: revise vs. flag-only vs. record unresolvable

With every observation from steps 1–4 collected, map each one to its action via
the decision table, then enforce the hard constraints over every applied
revision. Do not start until collection is complete — deciding mid-collection
risks destructive, mid-stream rewriting.

Canonical Step 5 decision table (observation → action) and Hard constraints on
revision: [`references/revision-rules.md`](references/revision-rules.md). Those
rules are authoritative; do not improvise.

### Step 6 — Emit output

Produce three fields. See **Output contract** below.

## Output contract

The skill emits exactly three fields. All three are required.

### `final_report` (string, markdown)

The post-review report. If no revisions were applied in step 5, this MUST
equal `draft_report` byte-for-byte. If revisions were applied, this is the
draft with those revisions and only those revisions. Every citation marker
must reference an `id` present in the emitted `cited_sources`.

### `cited_sources` (array)

The set of sources actually cited in `final_report`. May differ from the
input `cited_sources`:

- Sources whose citations were all dropped → removed.
- Sources newly added during step 3 → appended (with the same shape as input
  entries; the `id` comes from the finding's `source_id`).
- Sources swapped in → present; sources swapped out → removed if no other
  citation still uses them.

**Invariant:** every `id` in the emitted `cited_sources` must come from
`input.cited_sources ∪ {f.source_id for f in findings}`. No fabricated
sources.

### `review_notes` (string, markdown)

A structured report of what the reviewer did. Required sections:

```markdown
## Review summary
- claims_verified: <n>
- citations_swapped: <n>
- citations_added: <n>
- citations_dropped: <n>
- sentences_revised: <n>
- flags_raised: <n>
- unresolvable_gaps: <n>

## Applied changes
For each change, one bullet:
- [<change_type>] §<section>: "<short quote>" — <why> [old→new source if citation change]

## Flags (not changed)
For each flag, one bullet:
- [<flag_category>] §<section>: "<short quote>" — <why flagged, why not fixed>

## Unresolvable gaps
For each gap, one bullet:
- §<section>: "<short quote>" — claim asserted, no supporting finding in input. Recommend: <re-search this sub-topic | accept as acknowledged limitation | remove>.

## Top risks
At most 3 bullets summarizing the highest-impact issues the operator should know about. Empty list is allowed only if the review was clean.
```

`change_type` ∈ `{swap_citation, add_citation, drop_citation, revise_sentence,
delete_sentence, resolve_contradiction, delete_diagram}`.

`flag_category` ∈ `{off_topic, audience_mismatch_minor,
audience_mismatch_pervasive, internal_contradiction_unresolved,
unsupported_claim_softened, citation_unverifiable, diagram_inconsistent}`.

**Anti-rubber-stamp rule:** `review_notes` is never empty. Even a clean
draft produces a non-empty "Review summary" with `claims_verified > 0` and
the verified-ledger evidence that the checks were actually run.

## Failure modes

See [`references/failure-modes.md`](references/failure-modes.md) for the
detection-and-recovery table covering missing inputs, empty findings,
broken citation markers, oversized drafts, rubber-stamp output, style-edit
attempts, new-claim attempts, and pervasive audience mismatch. The
recoveries there are canonical; do not improvise.

## Anti-patterns

See [`references/revision-rules.md`](references/revision-rules.md) → Hard
constraints on revision for the canonical rule list. This skill's
anti-patterns are the inverse of those constraints (no loop, no re-search, no
fabrication, no paraphrase of passing prose, no hedging-to-placate, no
scope-widening, no empty `review_notes`, no section deletion, no silent
citation drops). If a behaviour is not in Hard constraints, this skill does
not enforce against it.

## Worked example

See `examples/remote-work.md` for a complete worked illustration. It traces
each bounded edit type (verified, swap_citation, revise_sentence + add_citation,
resolve_contradiction) on a single small draft and shows the full
`review_notes` shape including the verified-ledger and Top risks sections.

## Constraints

See [`references/revision-rules.md`](references/revision-rules.md) → Hard
constraints on revision for the canonical constraint list (single-pass,
closed-world citations, minimum-edit, audience-retargeting out of scope, etc.).

## Validation gate

See [`references/validation-gate.md`](references/validation-gate.md) for the
5-check pre-emit gate (citation-marker resolution, closed-world `id` set,
required `review_notes` sections, byte-for-byte equality invariant, and the
no-orphan-edits rule). Any failure must be fixed in place — do not emit a
known-broken output.
