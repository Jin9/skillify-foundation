# Extraction failure modes & anti-extrapolation

Cross-cutting, distilled from all three source reports. Sweep this before
emitting the spec; every applicable item must be absent or recorded in the
loss ledger.

## Failure modes
- **Lossy paraphrase as a rule.** Restating logic in plausible prose without a
  code span is a provenance break — a downstream agent cannot verify it.
  Fix: quote/point to `file:line`; if you can't, it's a ledger uncertainty.
- **Invented rule / identifier (hallucination).** Intrinsic (contradicts the
  code) or extrinsic (unsupported addition); ~17% identifier-hallucination
  base rate in LLM code work is the symmetric risk here. Fix: every rule
  traces to real code; no inferred APIs.
- **Edge-case / rare-branch drop.** Tail content is the first casualty of
  compaction and the highest-value part of a rules spec. Fix: each branch is
  its own rule, flowchart node, and pseudo-code line; never merge away a guard
  "for brevity".
- **Hallucinated flowchart edge / representation drift.** A flowchart node,
  arrow, or pseudo-code line with no cited branch — or the flowchart and
  pseudo-code disagreeing about the branch set. Same defect class as an
  invented rule. Fix: derive both from the cited decision sites, keep them in
  parity, and cite `file:line` on every node/line (see `ascii-logic-method.md`).
- **Recursive compaction drift.** Re-summarizing the spec compounds error.
  Fix: re-extract from source; never compact the spec onto itself.
- **Confidently-wrong (knowledge collapse).** Fluent spec, degraded accuracy,
  passes a casual read. Fix: provenance binding + coverage gate + human
  spot-check of flagged gaps.
- **Coverage blind spot.** Faithfulness metrics don't detect omission. Fix:
  ledger reports coverage against the step-1 decision-site inventory; gate on
  the minimum of coverage and alignment.
- **Dead-vs-live confusion.** Coded branch never exercised reported as an
  active business rule. Fix: cross-reference traces; mark
  exercised/not-observed/dead (see `trace-cross-referencing.md`).
- **Caller-scope incompleteness.** A rule whose inputs are set by
  un-inspected callers reported as total. Fix: follow callers; record
  un-inspected ones as gaps.

## Anti-extrapolation
The source research is calibration for *method*, not content. Never paste a
benchmark statistic, a vendor claim, or a contested research finding into the
spec as if it were a fact about the target codebase. The spec states only what
*this* code does, with *this* code's provenance. Where the evidence base
itself is contested (e.g., whether all unsupported content is a defect), this
skill takes the conservative provenance-first stance: unverifiable ⇒ ledger,
not spec.

## Scope guard
Extraction only. NOT: generating code, generating/“improving” requirements
(forward direction), writing tests, drawing architecture / system-topology /
sequence diagrams (a decision-logic flowchart of the extracted rules IS in
scope — see `ascii-logic-method.md`), or producing a general lossy
summary/docstrings. Those are different skills; the
deliverable here is a descriptive, traced rules spec + a loss ledger.
