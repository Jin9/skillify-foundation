# Validation gate

Before emitting `final_report`, `cited_sources`, and `review_notes`, run
every check below. Any failure must be fixed in place — do not emit a
known-broken output and rely on downstream stages to catch it. The gate is
load-bearing because `review-report` is the last research-content stage;
`reporting-research-run` only summarises.

1. Every citation marker in `final_report` resolves to an entry in the
   emitted `cited_sources`.
2. Every entry in `cited_sources` has
   `id ∈ input.cited_sources.ids ∪ {f.source_id for f in findings}`.
3. `review_notes` has the four required sections (Review summary, Applied
   changes, Flags, Unresolvable gaps, Top risks); `claims_verified > 0` if
   the draft contained any citations.
4. If `final_report == draft_report` byte-for-byte, then `Applied changes`
   is empty in `review_notes`.
5. Every `Applied changes` bullet corresponds to a specific observation
   from Procedure steps 2–4; no orphan edits.

## Why these five

Items 1 and 2 enforce the closed-world citation invariant — the single most
important contract this skill maintains. Item 3 catches the rubber-stamp
failure mode (an empty review is meaningless even on a clean draft). Items
4 and 5 catch silent-edit failure modes where the model writes prose changes
into `final_report` without acknowledging them in `Applied changes`, which
would break `reporting-research-run`'s edit accounting downstream.
