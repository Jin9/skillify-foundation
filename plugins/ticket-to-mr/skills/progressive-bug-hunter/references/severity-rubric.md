# Severity rubric (P0–P4)

How to summarize a diagnosed bug on a five-band severity scale. Severity rates
*how bad the bug is if the hypothesis holds* — it is a triage signal for the
fixer, not a measure of how sure you are. P0 is worst (direction matches
`treasury/incident-response`). Assign **exactly one** band to the bug.

## The bands

| Band | Name | Trigger — from observable code/trace/test evidence |
|------|------|----------------------------------------------------|
| P0 | Critical | outage / data loss or corruption / security breach on a **reachable core path**; fix immediately, blocks release |
| P1 | High | core feature broken, **no workaround**, **wide** blast-radius |
| P2 | Medium | impaired but a **workaround exists**; **bounded** blast-radius |
| P3 | Low | minor / edge-case defect, **narrow** reach |
| P4 | Trivial | cosmetic / negligible |

## How to assign

Read the band off concrete evidence already gathered in the hunt — do not guess:

- **Failure kind** — unhandled crash / exception type, panic, hang, silent wrong
  result. Data loss or corruption and memory-/type-safety breaches push upward.
- **Security-reachability** — does an attacker-influenced or untrusted input reach
  the faulty span on a real path? A reachable security defect floors at P0/P1; an
  unreachable one does not.
- **Blast-radius / reachability** — reuse the Tier-2 traversal from
  `symbol-graph-and-ast.md` (*Regression localization = impact / blast-radius*).
  Wide reachability (many callers / core path) → higher; an isolated leaf → lower.
  Distinguish "exists in the dependency tree" from "actually executed on this path".
- **Workaround** — a viable workaround caps the band at P2; none, on a core path,
  argues P0/P1.

When two signals disagree, take the higher band and state why in the report.

## Severity ≠ confidence

These are independent axes; report both and never collapse one into the other:

- **Severity** — impact *if the root-cause hypothesis is true*.
- **Confidence** — certainty that the hypothesis *is* true (the existing
  high/med/low signal on suspects and the hypothesis).

A low-confidence diagnosis of a P0-class failure is still reported as P0 (with low
confidence) — under-rating severity because the diagnosis is uncertain hides risk.

## Technical vs business impact

This skill sees only code, traces, and tests, so it rates **technical** severity.
When the final band genuinely depends on production/business context it cannot
observe (traffic on the path, regulatory exposure, customer reach), give the
evidence-supported band as **provisional** and name the missing context — do not
invent a business impact to justify a number.

## Worked example

A `NullPointerException` in a payment-settlement writer, reached on every
checkout, no guard, no workaround → **P0 (critical)**: data-correctness defect on
a reachable core path, wide blast-radius. Confidence may still be *medium* if the
exact trigger value is unconfirmed — severity stays P0, confidence stays medium.
