# Troubleshooting

Lookup table for common signals encountered while augmenting a report.
Default for any unlisted edge case: emit `augmented_report` equal to
`final_report` (zero diagrams) and log the reason in `augment_notes`. Zero
is always a valid answer for this skill.

| Signal | Action |
|---|---|
| `research_plan.depth` is `quick` | Short-circuit: emit input verbatim and notes saying "skipped — quick depth". Do not parse sections. |
| `research_plan.depth` is missing | Default to `standard`. Note nothing in `augmented_report`; record the default in `augment_notes`. |
| A candidate section has zero `[n]` markers | Skip it. The section has no grounding set; any diagram would be a new claim. Log `skip — no grounded findings`. |
| All candidate sections fail the scan test | Emit `augmented_report` equal to `final_report` (zero diagrams). Zero is a valid answer. Log per-section reasons in notes. |
| Two candidate sections discuss the same structure (e.g., both describe a request → router → model pipeline) | Pick the section whose grounding set is larger; skip the other to avoid duplicate signal. Note the dedup in `augment_notes`. |
| The first "prose paragraph" of a section is actually a bullet list | Treat the bullet list as the first paragraph for placement purposes — insert the fence after the list's terminating blank line. Do not insert between bullet items. |
| Final report has no body sections (only Executive Summary + Limitations) | Emit input verbatim. Zero diagrams. Log `skip — no eligible body sections`. |
| A diagram you composed exceeds 15 lines or 70 cols | Tighten it (drop a less load-bearing layer) or skip the section. Do not split one diagram across two fences. |
