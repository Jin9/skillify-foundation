# Anti-patterns

Sweep this table before emitting `augmented_report`. If any row applies,
fix the issue rather than emitting and relying on `review-report` to catch
it (note: `review-report` runs once during stage 5 — backfill is by
definition past that gate).

| # | Anti-pattern | Why dangerous | Correct alternative |
|---|---|---|---|
| 1 | **Prose touch-up** — "while inserting the diagram I noticed the section's first sentence reads awkwardly, so I tightened it" | Violates the insertion-only contract; downstream tooling diffs input vs. output expecting equality outside the inserted fence | Insert the fence and emit. Prose edits are `review-report`'s job, not this skill's. |
| 2 | **Cited diagram** — `[3]` or `(f12)` markers inside the fenced block | Reviewer cannot audit a diagram for citation grounding; smuggles a claim through the structural side door | Place all `[n]` markers in the surrounding prose. The diagram names entities; the prose carries the cites. |
| 3 | **Boilerplate-section diagram** — a flow chart in `## Executive Summary` or a stack diagram in `## Sources` | Eats summary scannability; breaks the section's contract; reviewer expects boilerplate to stay prose-only | Move the diagram to a themed body section where its grounding findings already live. Skip if no body section fits. |
| 4 | **Padding to 3** — emitting two weak diagrams to "show effort" after one strong one is identified | Lowers signal density; trains the reviewer to scroll past diagrams; reduces the impact of the one good diagram | Emit only diagrams that pass the scan test. Zero is valid; one is common; three is rare. |
| 5 | **Novel-entity diagram** — putting a node in the diagram whose name doesn't appear in the section's cited findings | The diagram becomes a new claim, and an unauditable one | Build every diagram strictly from entities already in the section's grounding set. If the structure is new, propose it as a `synthesize-report` rerun, not a backfill insertion. |
| 6 | **Sources-list edit** — re-sorting, renumbering, or adding to the `## Sources` list (e.g., "I noticed the order didn't match the appearance order in prose") | Breaks the bijection that downstream consumers (review-report, reporting-research-run) rely on | Leave `## Sources` byte-identical. If a real ordering bug exists, file it; do not silently fix it here. |
| 7 | **Decorated heading** — adding box-drawing characters to a section heading to "frame" it | Headings are not diagrams; this confuses parsers and screen readers | Headings stay as `## Title`. Diagrams sit below the first prose paragraph, not on the heading line. |
| 8 | **Language-tagged fence** — opening the fence as ` ```text ` or ` ```diagram ` | Some markdown renderers will apply syntax-highlighting rules to box-drawing characters and corrupt the visual; project rule is no language tag on diagrams | Open every diagram fence with three backticks and an immediate newline. |
| 9 | **Quick-depth override** — adding diagrams to a `quick` report because "structure would help" | Quick reports are budgeted at 300–600 words; a diagram is ~10% of that budget for marginal signal | Honour the depth gate. Quick reports stay diagram-free. |
| 10 | **Section reorder via insertion** — inserting the fence at the wrong byte offset, shifting later headings into a different visual order | Breaks the section parse for downstream rendering; harder to diff than a prose edit | Insertion offset is the newline immediately after the first prose paragraph of the chosen section. Compute it from Step 2's recorded offsets, not by string-matching that may collide. |
