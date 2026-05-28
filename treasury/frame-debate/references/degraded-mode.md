# Reference: degraded-mode procedure (no squad-researcher run dir)

`squad-brainstorm` is designed to debate a *finished squad-researcher run*.
Degraded mode is the explicit, honest fallback when no such run is available
— not a second first-class path. It always costs the run a confidence cap.

## When degraded mode triggers

In Procedure step 1, after `research_run_dir` fails to resolve to a dir
containing all three of `05-final_report.md`, `03-findings.json`,
`01-research_plan.json`:

- `topic` non-empty (after trim) → `grounding: degraded`, proceed below.
- `topic` empty/whitespace → emit the `brief_skipped` shape, stop.

A partial researcher dir (e.g. has `05-final_report.md` but no
`03-findings.json`) is treated as **not resolved** → degraded. Do not
half-read it; the grounding pack must be coherent.

## Building the degraded brief

1. `grounding_pack.summary`:
   - With `source_text`: a ≤400-word neutral digest of it.
   - Without `source_text`: exactly one sentence —
     "Ungrounded debate: no research artifacts or source text supplied;
     panelist positions are reasoning-only and unverified."
2. `grounding_pack.items`:
   - With `source_text`: extract discrete factual passages as `g-n` items;
     `confidence: "low"` for all (no provenance chain), `disputed: false`,
     `origin: "source_text"`.
   - Without `source_text`: `items: []`. Panelists then argue from reasoning
     alone and must say so (their `uncited_assertions` will be large — that
     is expected and honest, not a failure).
3. `grounding_pack.out_of_scope`: `[]` (no research plan to carry).
4. `contested_questions`: inferred from `topic`/`source_text`, all
   `seed: "degraded_inferred"`, `evidence_refs: []` allowed.

## The mandatory banner (`grounding_note`)

Set `debate_brief.grounding_note` to a single line that flows verbatim
through `synthesize-consensus` (which caps `overall_confidence` at `medium`
and prints it under "Confidence & caveats") into `report-debate`'s header:

```
⚠️ degraded grounding — debated from <topic | supplied source text> with no
squad-researcher findings; conclusions are unverified and confidence is
capped at medium.
```

## Cheapest valid run

`rounds: quick` + `grounding: degraded` is explicitly the minimum viable
debate: 3 opening positions, no cross-examination, one open-moderator merge,
one report. Use it for a quick multi-model gut-check on a topic with no
prior research; never present its output as grounded.
