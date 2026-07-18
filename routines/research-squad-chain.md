---
routine: research-squad-chain
summary: Chain the six squad-researcher skills end to end - plan, search, extract, synthesize, review, and report a research run, with each stage's JSON payload embedded in its node artifact.
version: 0.1.0
default_on_fail: stop
owner: chinnawat.w
tags: research, squad
requires: web search backend available to the search-sources skill
---

Runs the treasury squad-researcher chain as black boxes. Stage payloads are JSON-shaped per each skill's own I/O contract; every node artifact embeds that stage's JSON payload plus a short human-readable summary, so downstream nodes consume artifacts without re-deriving state.

## Inputs

- topic: research question or subject for the run (required)
- depth: quick, standard, or deep (optional, default standard)
- audience: who reads the final report (optional, default technical generalist)

## Nodes

### Node: plan

- purpose: Decompose the topic into a research plan of sub-questions and dimensions
- executor: plan-research
- tier: mid
- inputs: user.topic, user.depth
- outputs: 01-research-plan.md
- gate: after
- notes: Artifact embeds the research_plan JSON payload. Human approves scope before any searching starts.

### Node: search

- purpose: Search sources for every sub-question in the plan
- executor: search-sources
- tier: mid
- inputs: 01-research-plan.md
- outputs: 02-sources.md
- gate: none
- notes: Embeds the sources JSON payload. Requires the web search backend named in requires.

### Node: extract

- purpose: Extract findings and contradictions from the gathered sources against the plan
- executor: extract-findings
- tier: mid
- inputs: 01-research-plan.md, 02-sources.md
- outputs: 03-findings.md
- gate: none
- notes: Embeds the findings JSON payload.

### Node: synthesize

- purpose: Synthesize the findings into a cited draft report for the audience
- executor: synthesize-report
- tier: frontier
- inputs: 01-research-plan.md, 02-sources.md, 03-findings.md, user.audience
- outputs: 04-draft-report.md
- gate: none
- notes: Embeds draft_report and cited_sources payloads.

### Node: review

- purpose: Review and finalize the draft report for accuracy, citation integrity, and fit to audience
- executor: review-report
- tier: frontier
- inputs: 04-draft-report.md, 03-findings.md, user.topic, user.audience
- outputs: 05-final-report.md
- gate: after
- notes: Embeds final_report, cited_sources, and review_notes payloads. Human approves the final report here.

### Node: report-run

- purpose: Produce the canonical single-page summary of the whole run
- executor: reporting-research-run
- tier: small
- inputs: user.topic, user.depth, user.audience, 01-research-plan.md, 02-sources.md, 03-findings.md, 05-final-report.md
- outputs: 00-run_report.md
- gate: none
- on_fail: continue
- notes: The executor's own contract writes exactly 00-run_report.md into the run dir - distinct from the launcher's run-report.md. The validator's sequence-prefix warning on this node is expected and accepted.
