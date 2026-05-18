# Loss-bounding extraction method

Distilled from research report *Mitigating lossy compression in summarization*.
Extraction is lossy compression by construction — a blind feed-forward step
that cannot know which currently-peripheral rule a future reader needs, and a
mis-extracted or omitted rule is unrecoverable from the spec. The job is to
*bound* that loss and *account* for it, not to pretend it is zero.

## Principles (apply in workflow order)

1. **Stay on the faithful/extractive side of the frontier.** Faithfulness and
   compression are not independently dialable; "more faithful" summaries are
   often just more extractive. For a business-rule spec that is the *correct*
   regime — quote the load-bearing code, do not paraphrase a rule into
   plausible prose. An unsupported-but-true addition is still a provenance
   break for a reader who cannot verify it.
2. **Salience-first, before any abstraction.** Compute the rule/entity set
   (conditions, thresholds, named domain entities, decision paths) and pass it
   as an explicit retention constraint into any abstractive pass. Providing
   salient entities as a control signal demonstrably reduces entity loss; an
   unconstrained abstractive pass drops mid-importance rules.
3. **Order matters (positional U-shaped loss).** Models under-use mid-context
   content. When serializing many decision sites for extraction, do not bury
   the highest-value rules in the middle of a long pass; chunk by module so
   each rule sits near a boundary.
4. **Bound recursion depth.** Cumulative degradation — not single-step error —
   governs viability; rare/edge rules are the first casualties of repeated
   compaction. Never re-summarize an already-extracted spec; re-extract from
   source. Treat the spec as durable externalized state, not a transient
   summary that gets re-compacted.
5. **Extract-then-abstract pipeline.** Select the exact code spans first
   (extractive selector), then write the rule statement over the reduced
   input. This caps token pressure on large modules and keeps the abstractive
   step anchored to selected evidence.
6. **Gate on coverage/omission, not faithfulness alone.** Faithfulness answers
   "is what is here correct?"; it does not answer "what is missing?". Take the
   minimum of a coverage score (decision sites with a rule) and an alignment
   score (no contradicted/invented rule). An omission-heavy spec must not pass
   on faithfulness.
7. **Instruments are weak — keep a human in the loop.** Automated consistency
   metrics top out ~60–75% balanced accuracy and degrade on long outputs;
   treat any automated check as a noisy filter and route the ledger's flagged
   gaps to a human spot-check, not as a sufficient sign-off.

## The dangerous signature
"Knowledge collapse": factual accuracy degrades while surface fluency
persists — a confidently-wrong spec that passes a casual read. The defenses
above (provenance binding, coverage gate, bounded recursion, human spot-check)
exist specifically to catch it. Calibration priors only — never paste a
research statistic into the spec as if it described this codebase.
