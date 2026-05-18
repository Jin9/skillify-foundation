# Requirement ↔ code traceability (reverse direction)

Distilled from research report *Requirement-to-code pipeline*, applied in
reverse (code → spec → requirement xref). The pipeline literature defines six
edges (PRD→story→AC→spec→code→test→traceability); this skill walks the
spec↔code and AC↔code edges *backwards* to recover an audit-grade rule spec.

## Rule statements in a disciplined notation
Write each extracted rule in a constrained, testable form so it is
unambiguous and round-trippable to an acceptance criterion:
- **EARS**: `When <trigger>, the <system> shall <response>`; plus the
  Ubiquitous / Event-driven / State-driven / Optional / Unwanted templates for
  invariants, events, states, optional features, and forbidden behavior.
- **Given/When/Then** (Gherkin) when the rule is scenario-shaped.
These notations force a named subject, an explicit trigger, and an observable
response — exactly the fields a downstream reviewer or test needs. Free-form
prose rules are the documented failure mode (ambiguous trigger, missing
negative case, role confusion).

## Provenance is mandatory, hallucination is the top risk
~17% of LLM-written code references a non-existent identifier/API; the
symmetric risk here is asserting a rule the code does not implement. Every
rule MUST cite `path:line(-line)` and quote/point to the actual span. If you
cannot point at the code, the rule does not go in the spec — it goes in the
ledger as an uncertainty. Decision/branch logic becomes a **decision table**
(inputs → outcome) with each row cited, not a prose paragraph.

## Completeness: follow the callers
Silent-regression evidence: agentic edits test only plan-scope and miss
untouched callers. The extraction analogue: a rule is incomplete if its
guards/inputs are set by callers you did not inspect. For each extracted rule,
check the call sites that supply its inputs; note un-inspected callers as
coverage gaps in the ledger rather than assuming the rule is total.

## Traceability as the audit substrate (AI proposes, human disposes)
The production pattern across RE tools is identical: AI proposes a trace link
with a confidence score; a human (or a merge gate) confirms before it is
treated as authoritative. Mirror it: every rule→requirement link carries a
confidence and is marked `proposed` until confirmed. The end state is a
bidirectional map `requirement-id ↔ code(file:line) ↔ rule-id`, which is what
turns an ad-hoc extraction into a compliant, auditable spec. Where no
requirement source exists, the link is `none` — never invent a requirement id
(that is forward generation, out of scope).
