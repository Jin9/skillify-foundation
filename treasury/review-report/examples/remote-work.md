# Worked example: remote-work

Example invocation of the `review-report` skill demonstrating each bounded edit type: verified-claim recording, citation swap, sentence revise + citation add (for an unsupported uncited claim), and resolve-contradiction (for a conclusion that doesn't follow from the corrected claims). Shows the `review_notes` shape with all required sections.

## Inputs (abbreviated)

```
topic:    "Effect of remote work on engineering productivity, 2020-2025"
audience: "engineering manager"

cited_sources: [
  {id: "S1", title: "GitHub Productivity Report 2022", url: "..."},
  {id: "S2", title: "Stanford WFH study 2024",         url: "..."},
  {id: "S3", title: "Atlassian State of Teams 2023",   url: "..."}
]

findings: [
  {claim: "Devs working fully remote shipped 13% more PRs in 2022 vs 2019", evidence: "...", source_id: "S1", confidence: "high"},
  {claim: "Hybrid teams reported 20% higher retention than fully in-office",  evidence: "...", source_id: "S2", confidence: "high"},
  {claim: "Code-review latency rose 40% during 2020 lockdowns",               evidence: "...", source_id: "S4", confidence: "medium"},
  {claim: "Atlassian survey: 67% of teams adopted async-first practices",     evidence: "...", source_id: "S3", confidence: "high"}
]

draft_report contains:
  ...full-remote developers shipped 13% more PRs in 2022 [S1].
  Hybrid teams saw a 35% retention boost [S2].
  Code-review latency stayed roughly flat through the pandemic.
  Most teams now run async-first [S2].
  Therefore, going fully remote is the highest-leverage productivity bet.
```

## Procedure trace

Step 2 finds:

- "13% more PRs … [S1]" → `verified` (finding matches exactly).
- "35% retention boost [S2]" → `unsupported_by_cited_source` (finding S2 says 20%, not 35%).
- "Most teams … async-first [S2]" → `citation_broken` (S2 doesn't cover this; S3 does).

Step 3 finds:

- "Code-review latency stayed roughly flat" → uncited factual assertion. Finding from S4 contradicts it (40% rise). → either swap claim and cite S4, or it's an internal contradiction.

Step 4 finds:

- The conclusion "going fully remote is the highest-leverage productivity bet" doesn't follow from the (corrected) claims — hybrid teams have the retention advantage. `internal_contradiction`, resolvable: weaken the conclusion.
- Audience fit: "engineering manager" — fine, no jargon issues.

Step 5 decisions:

- Swap S2→nothing for retention: there is a finding (S2, 20%). Revise the sentence to match: "Hybrid teams saw a 20% retention boost [S2]" (revise_sentence + adjusted citation).
- Swap citation S2→S3 for the async-first claim (swap_citation).
- Revise the uncited code-review-latency sentence to match S4 (revise_sentence + add_citation to S4).
- Soften the conclusion: revise from "going fully remote is the highest-leverage…" to "the data does not support a single-mode optimum; hybrid teams show the strongest retention while fully-remote shows the strongest throughput" (revise_sentence + resolve_contradiction).

## Emitted output

Emitted `final_report`: the draft with the four revisions applied.

Emitted `cited_sources`: original [S1, S2, S3] + S4 added.

Emitted `review_notes` (abbreviated):

````markdown
## Review summary
- claims_verified: 1
- citations_swapped: 1
- citations_added: 1
- citations_dropped: 0
- sentences_revised: 3
- flags_raised: 0
- unresolvable_gaps: 0

## Applied changes
- [revise_sentence] §Retention: "35% retention boost" → "20% retention boost" — finding S2 reports 20%, not 35%.
- [swap_citation]   §Async work: async-first claim — S2→S3 (Atlassian survey is the actual source).
- [revise_sentence + add_citation] §Latency: replaced "stayed roughly flat" with "rose ~40% during 2020 lockdowns" [S4].
- [resolve_contradiction] §Conclusion: weakened "fully remote is the highest-leverage bet" to a two-sided summary; the original conclusion was not entailed by the (corrected) claims.

## Flags (not changed)
(none)

## Unresolvable gaps
(none)

## Top risks
- The original draft's conclusion contradicted its own evidence; the corrected conclusion is now consistent but materially different in recommendation. The operator should review whether downstream readers expecting the original directional claim need an explicit "revised vs. previous draft" note.
````
