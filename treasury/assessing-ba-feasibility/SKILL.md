---
name: assessing-ba-feasibility
description: >
  Stage 4 of the BA pipeline: consume the Story Set, a clear Governance Check, and the raw requirement, and emit a typed Feasibility Note contract (at least two options with pros and cons, dependencies, risks with mitigations, time-boxed spikes, and a verdict of buildable, not-now, or phased) wrapped in the shared envelope, for the Tech Lead verdict gate (G4). It sets build scope before engineering commits. Use when the user is running the BA pipeline and asks to "assess feasibility", "produce the Feasibility Note", "give the TL options and risks", or "is this buildable for the pipeline". Do NOT use for standalone feasibility outside the pipeline (use technical-feasibility); do NOT run it while governance is blocked.
compatibility: claude-code, codex, gemini-cli, opencode
---

# assessing-ba-feasibility

## Purpose
Stage 4 (`feasibility-agent`) of the BA pipeline. Consume the **Story Set**, a **clear** Governance Check, and the
**raw requirement**, and emit a typed **Feasibility Note** — options with trade-offs, dependencies, risks with
mitigations, and time-boxed spikes, ending in a `buildable`/`not-now`/`phased` verdict — wrapped in the shared
envelope, ready for Gate G4. It sets build scope before engineering commits.

## When to use
- Triggers: *"assess feasibility"*, *"produce the Feasibility Note"*, *"give the TL options and risks"*, *"is this buildable for the pipeline"*.
- **Not this skill:** standalone feasibility outside the pipeline → `technical-feasibility`.

## Model & cost
Runs on a **frontier** model (e.g. Opus 4.7 / top Gemini), high reasoning effort. This is the one node that genuinely wants the heavy model — it weighs real technical trade-offs.

## Input
- The **Story Set** (S2) — read `02-story-set/INDEX.md`, then only the `ST-NN-*.md` you need to weigh — the **Governance Check** (S3, must be `verdict: clear`), and the **raw requirement**.

## Output
A **Feasibility Note** contract — the shared envelope plus the S4 fields.

```
# Feasibility Note — <task_id>
## envelope
task_id: REQ-<slug>            # same trace id throughout
stage: S4
intent: assess buildability + set scope
state: buildable | not-now | phased    # mirrors verdict
confidence: high | medium | low
provenance: { raw_ref: <raw requirement>, upstream_contract_ref: [ story_set@1.0, governance_check@1.0 (clear) ] }
produced_by: feasibility-agent
owner: Tech Lead
schemaVersion: 1.0
created_at: <RFC-3339>
## contract
options:      [ { name, pros: [ ... ], cons: [ ... ] } ]   # at least 2
dependencies: [ ... ]                                      # systems, vendors, data feeds, prerequisite work
risks:        [ { risk, mitigation } ]                     # every risk carries a mitigation
spikes:       [ { question, time_box } ]                   # time-boxed unknowns
verdict:      buildable | not-now | phased
```

## Decision rules
1. **Refuse to run unless the Governance Check is `clear`.** An open blocker means stop — do not assess feasibility past it.
2. Read the **raw requirement** as well as the contracts, so the verdict reflects the original ask, not only the distillation.
3. Give **≥2 options**, each with pros and cons — never a single forced answer.
4. Every **risk carries a mitigation**; an unmitigated risk is itself a finding.
5. Turn unknowns into **time-boxed spikes** (`{ question, time_box }`), not open-ended research.
6. Choose the verdict honestly: `buildable` (MVP scope holds), `phased` (split now / later), or `not-now` (a dependency or risk blocks).
7. **Redact any real PII** to `<PII:REDACTED:CLASS=...>`.
8. **Read the Story Set selectively** — the INDEX first, then only the stories your options and risks touch; reference findings back to their `ST-NN`. Never load every story at once.

## Checklist
- [ ] Governance Check input is `clear` (else stopped)
- [ ] Story Set read via INDEX + selected stories (not the whole set)
- [ ] Raw requirement was in the input
- [ ] ≥2 options, each with pros and cons
- [ ] Dependencies listed (systems, vendors, data feeds, prerequisite work)
- [ ] Every risk has a mitigation
- [ ] Unknowns are time-boxed spikes
- [ ] `verdict` set; envelope complete; no real PII

## Anti-patterns (never do)
- Assess feasibility while governance is `blocked`.
- Offer a single option as if there were no alternative.
- List a risk with no mitigation, or an unknown with no time-box.
- Silently widen scope beyond the stories and the raw requirement.

## Example (ShopPilot MVP)
**options:** (a) atomic DB-row stock decrement at confirm — simple, fits MVP; (b) reservation table + TTL sweep — more
moving parts, needed for a hold window. **dependencies:** mock payment provider, mock courier, region data-residency,
e-Tax broker for VAT. **risks:** last-item race → atomic decrement + integration test; webhook double-delivery →
idempotency key per event id; price drift → snapshot price/name into the order at creation. **spikes:** stock
reservation TTL & sweep (2 days); idempotent webhook handler (1 day). **verdict:** `buildable`.

## Human approval gate (G4)
**Stop.** The Tech Lead gives the feasibility verdict (sync, named). The verdict sets build scope and cost, so the TL
owns it — the agent proposes options and risks; it does not decide buildability.
