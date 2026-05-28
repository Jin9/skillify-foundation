# Anti-patterns

Sweep this table before emitting `findings`. If any row applies to a finding
you're about to emit, fix it in place — do not emit and rely on downstream
stages to catch it. The same table is the basis for validation-gate item 10
in `SKILL.md`.

| # | Anti-pattern | Why dangerous | Correct alternative |
|---|---|---|---|
| 1 | **Inventing source ids** | Claims become unverifiable; downstream review can't trace them. Hallucinated citations are the canonical RAG failure. | Every `source_id` and entry in `supporting_source_ids` / `disputed_by` must be an id present in input `sources[]`. Unknown id = drop the finding. |
| 2 | **Claim without evidence** | The model writes a plausible claim and attaches a citation that doesn't actually support it ("post-rationalized" citation — ~57% of vanilla-RAG citations). | Extract the `evidence` excerpt first, then write the claim from it. If you can't quote/paraphrase the source on this point, drop the finding. |
| 3 | **Paraphrasing that strips hedging** | Source: "may reduce" → claim: "reduces". Changes the technical meaning and inflates apparent strength. | Preserve the source's hedging verbatim in the claim. If the evidence is hedged, the claim must be too. |
| 4 | **Collapsing contradictions** | Two sources disagree; model picks one silently. Loses signal that the question is contested — exactly what synthesis needs to know. | Emit both claims as separate findings, each with the opposing source in `disputed_by`. Let `synthesize-report` decide how to present the disagreement. |
| 5 | **Over-extraction (>7 per sub-question)** | More findings correlates with lower downstream factual accuracy ("information overload" in deep-research evals). | Cap at 7 per sub-question. Rank by centrality; drop the tail. If the cap is biting hard, the sub-question is probably too broad — flag it. |
| 6 | **Numeric confidence scores** | LLM-generated `0.0–1.0` confidences are poorly calibrated (ECE > 0.37 in evals); they collapse to overconfidence. | Use the qualitative `high \| medium \| low` rubric only. Map evidence properties → bucket, not vibes → number. |
| 7 | **Dumping into "unmapped"** | Unmapped becomes a junk drawer when the model can't be bothered to align findings to sub-questions; the outline is lost. | Use `unmapped` rarely and always with a `note` explaining why the finding is too important to drop. Many unmapped findings = the research_plan needs revision (a `review-report` issue, not this skill's). |
| 8 | **Silent skip of unreadable sources** | A source with a broken body silently disappears; the user can't tell whether it was off-topic or unprocessable. | Add the source to `unsupported_sources[]` with a reason. Downstream review depends on this signal. |
| 9 | **Mixing faithfulness with correctness** | The skill drops a claim because it seems false in the world, not because the source fails to support it. Suppresses signal that downstream `review-report` needs. | This skill enforces only faithfulness: the citation supports the claim as stated. Correctness in the world is `review-report`'s job. Surface contested claims, do not pre-filter. |
| 10 | **Reordering sub-questions** | Findings emitted in source-by-source order force `synthesize-report` to reorganize. Lost outline. | Iterate sub-question by sub-question, in plan order. `coverage[]` rows mirror `research_plan.sub_questions[]` order. |
