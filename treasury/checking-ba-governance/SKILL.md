---
name: checking-ba-governance
description: >
  Stage 3 of the BA pipeline: sweep a Story Set for governance and privacy gaps and emit a typed Governance Check contract (PII flags, compliance flags, and blockers each with a named resolve owner, plus a verdict of clear or blocked) wrapped in the shared envelope; on blocked it fires gate G3 and loops back to a named human. It flags, it never fixes, and it never resolves its own blocker. Use when the user is running the BA pipeline and asks to "run the governance check", "flag PII and compliance gaps", "is this clear to hand off", or "produce the Governance Check". Do NOT use for deep banking KYC, AML, or PCI compliance mapping and PRDs (use analyzing-banking-requirements); do NOT auto-resolve any blocker here.
compatibility: claude-code, codex, gemini-cli, opencode
---

# checking-ba-governance

## Purpose
Stage 3 (`governance-agent`) of the BA pipeline. Run a lightweight **surface-don't-repair** sweep over a Story Set and
emit a typed **Governance Check** — PII flags, compliance flags, and human-owned blockers, with a `clear`/`blocked`
verdict — wrapped in the shared envelope. On `blocked` it fires Gate G3 and the pipeline loops back to a named human.
This node **flags only**: it never edits the stories and never clears its own blocker.

## When to use
- Triggers: *"run the governance check"*, *"flag PII and compliance gaps"*, *"is this clear to hand off"*, *"produce the Governance Check"*.
- **Not this skill:** deep banking KYC / AML / PCI compliance mapping and PRDs → `analyzing-banking-requirements`.

## Model & cost
Runs on a **mid** model (e.g. Sonnet 4.6 / Gemini 3 Pro), medium effort — a checklist sweep. Escalate to a **frontier** model only when compliance wording is genuinely ambiguous.

## Input
- A **Story Set** (S2 output, `state: drafted`) — read `02-story-set/INDEX.md`, then only the `ST-NN-*.md` files you need to flag. Never load every story at once.

## Output
A **Governance Check** contract — the shared envelope plus the S3 fields. Flag every gap you find; resolve none.

```
# Governance Check — <task_id>
## envelope
task_id: REQ-<slug>            # same trace id as the Story Set
stage: S3
intent: surface governance + privacy gaps (flag only)
state: clear | blocked         # mirrors verdict
confidence: high | medium | low
provenance: { raw_ref: <source>, upstream_contract_ref: story_set@1.0 }
produced_by: governance-agent
owner: BA lead + named SME (Legal / DPO / Compliance)
schemaVersion: 1.0
created_at: <RFC-3339>
## contract
pii_flags:        [ { field, pii_class, where_seen, story_ref } ]               # story_ref = ST-NN (+ file) where the field is seen
compliance_flags: [ { topic, concern, regulator?, named_owner?, story_refs? } ] # story_refs = the ST-NN(s) it applies to
blockers:         [ { id, type, description, resolve_owner, story_refs? } ]      # a named human clears each; the agent never does
verdict:          clear | blocked
```

## Decision rules / anti-patterns (these are the guardrails — never break them)
1. **Flag, never fix.** Surface every gap as a flag or blocker; never edit the stories to "repair" it.
2. **Never auto-resolve a blocker** (DENY). Only a named human clears it; every blocker carries a `resolve_owner`.
3. **Redact any real PII** to `<PII:REDACTED:CLASS=...>`; record the field and class in `pii_flags`, never the value.
4. **Never conflate Compliance with Legal.** Compliance describes the rule; Legal interprets the wording — record them as different `named_owner`s.
5. Any open blocker ⟹ `verdict: blocked` ⟹ the pipeline cannot advance to feasibility or handoff.
6. A missing reviewer (Legal / DPO / Compliance absent on a flow touching PII or customer-facing copy) is itself a blocker, not a pass.
7. **Read selectively** — open `02-story-set/INDEX.md` first, then only the stories you need; reference each flag and blocker back to its `ST-NN` (and file). Never load every story at once.

## Checklist
- [ ] Every personal-data field is in `pii_flags` with its class and where seen
- [ ] Compliance concerns (retention, consent, citations, dual-approval, customer-facing copy) flagged with an owner
- [ ] Each blocker has a `resolve_owner` (a named human or role)
- [ ] `verdict` set and `state` mirrors it
- [ ] Each flag / blocker references the `ST-NN` (and file) it came from
- [ ] No story was edited; no blocker was self-resolved; no real PII echoed

## What pages a human
A `blocked` verdict pages the BA lead + the named SME to open Gate G3. Resolution happens outside this node; the
affected upstream stage then re-runs with the enriched input, and this node re-runs to confirm `clear`.

## Example (ShopPilot MVP)
On the business-only spec the sweep surfaces **seven gaps → `verdict: blocked`**: (1) PII collected with no per-field
inventory / lawful basis / retention → `resolve_owner: DPO`; (2) Legal/DPO absent on a PII + customer-facing flow →
Legal / DPO / Compliance; (3) retention period unstated → DPO + Finance; (4) regulatory citations unresolved →
Compliance + Legal; (5) dual-approval owner missing for admin cancel / refund → Finance / Ops; (6) compensating
action missing for paid-but-cancelled → Finance; (7) customer-facing copy risk (payment-failure / out-of-stock
wording) → Legal + UX. After named humans resolve each and the stories re-run, the check re-runs to `clear`.

## Human approval gate (G3)
**Stop on `blocked`.** A named SME (Legal / DPO / Compliance) plus the BA lead resolve each blocker before the
pipeline continues. This node never marks a blocker resolved and never relaxes the gate.
