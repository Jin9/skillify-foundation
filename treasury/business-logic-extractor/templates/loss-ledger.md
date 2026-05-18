# Information-loss / coverage ledger: <scope>

> Mandatory companion to `business-logic-spec.md`. Gate on THIS, not on
> faithfulness alone. An omission-heavy extraction fails here even if every
> stated rule is correct.

## Decision-site coverage (the denominator)
Inventory of decision sites found in scope (step 1) vs. rules produced:

| module / file | decision sites found | rules emitted | status |
|---|---|---|---|
| `path/a.ext` | 12 | 12 | covered |
| `path/b.ext` | 9  | 5  | partial — see omissions |
| `path/c.ext` | 4  | 0  | uncovered — not yet extracted |

- coverage = covered_sites / total_sites = **<n/m>**
- alignment = no contradicted/invented rule? **<yes/no>**
- **gate = min(coverage, alignment)** → <pass | fail> (threshold: <state it>)

## Preserved verbatim vs. abstracted
- Verbatim (code quoted): R1, R3, R7 …
- Abstracted (paraphrased over selected spans): R2, R5 … (each still cited)
- Nothing was restated without a `file:line` anchor: <yes/no>

## Omissions & uncertainties (enumerated — do not hide)
- `path/b.ext:* ` — <which decision sites have no rule yet and why>
- R4 — requirement link unconfirmed (proposed, confidence low)
- R6 — `not-observed` in traces; cannot confirm it is live
- callers of `funcX` not inspected → R2 input domain may be incomplete
- dead-code candidates: <list, flagged not dropped>

## Confidence band
Automated checks are a noisy filter (consistency metrics ~60–75% ceiling);
this ledger is not a sufficient sign-off on its own.

## Human spot-check checklist (route these before sign-off)
- [ ] every `partial`/`uncovered` module reviewed
- [ ] each `proposed` requirement link confirmed or rejected
- [ ] each `not-observed`/`dead` rule adjudicated (rare path vs. dead code)
- [ ] edge-case rows sampled against source
- [ ] no rule asserts behavior absent from the cited span
