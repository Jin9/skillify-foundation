# Anti-patterns

Sweep this table before emitting `draft_report`. If any row applies, fix
the offending section in place — do not emit and rely on `review-report` to
catch it. Rows 10–11 cover diagrams; the rest cover prose and citations.

| # | Anti-pattern | Why dangerous | Correct alternative |
|---|---|---|---|
| 1 | **Hallucinated synthesis** — combining two findings into a claim neither finding supports | Looks fluent but is fabrication; downstream review-report may not catch it | Restrict synthesis to a labeled `## Synthesis` section (deep only) and cite all underlying findings; otherwise stay within the literal claim of one finding |
| 2 | **Ghost citations** — `[n]` markers pointing to sources not in the input | 3–13% of unaided LLM-generated citations are fabricated; the worst single failure mode | Build `cited_sources` only by selecting from sources already referenced in `findings` |
| 3 | **Citation drift** — placing cites at paragraph end instead of adjacent to the claim | Recall / precision drop sharply with end-of-paragraph clustering | Place `[n]` immediately after the cited sentence or clause |
| 4 | **Audience default-drift** — writing academic prose for a general audience | Wastes the `audience` input; report fails its primary user | Pass each section through the audience-profile table before drafting |
| 5 | **Padding** — restating topic, methodology, or prior sections to hit length | Lowers signal density; review-report may miss it because prose is "fine" | Cut a section instead of stretching it; under-length is acceptable if findings warrant |
| 6 | **Contradiction laundering** — picking one of two disagreeing findings and silently dropping the other | Misrepresents the source state; biases the report | Name the disagreement in prose; cite both sides; flag in Limitations |
| 7 | **Topic restatement loop** — repeating the topic in every section opening | Reads as filler; insults the reader | Title + executive summary is the topic-restatement budget |
| 8 | **Coverage decay** — later sections shorter and looser than earlier ones | Symptom of single-pass generation | Outline first (Step 4) with per-section finding assignments before drafting any prose |
| 9 | **Source-list / prose disagreement** — `## Sources` lists items the prose never cites, or omits items the prose does cite | Breaks the contract with review-report | Generate `cited_sources` as the source-of-truth in Step 5; render `## Sources` from it in Step 7 |
| 10 | **Decorative diagrams** — ASCII that doesn't add structural information the prose lacks (e.g., a one-box "Memory" sketch, or a flow that just repeats the section's first sentence in arrow form) | Eats vertical space, signals padding, breaks scannability, and dilutes the diagrams that DO carry signal | Skip the diagram. One good paragraph beats one weak diagram. Apply the scan test in Procedure Step 6. |
| 11 | **Diagrams as new claims** — putting an entity, arrow, or layer in a diagram that the prose / findings don't justify | The reviewer cannot audit a diagram for citation grounding; an unsupported diagram smuggles in a claim through the structural side door | Build every diagram strictly from entities already in the section's selected findings; if the structure is new, write it as prose first and only then render it |
