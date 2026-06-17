# Loss bounding

How to preserve the same context as the source while keeping the contract lean. The goal is **bounded** loss — not zero loss, but loss that is salient-aware, explicit, and recoverable.

## The principle

A contract is the *smallest sufficient* set of high-signal content the next stage needs to act — plus a record of what was set aside and a pointer back to the rest. Extract for the consumer's decision, not for completeness. Completeness lives in the source; the contract carries the decision-relevant subset and a way back.

## Salience tiers

Tag every unit of source content into exactly one tier:

- **must-preserve** — the next stage cannot act correctly without it. Goes into `contract`.
- **droppable** — relevant but reconstructable from the source if ever needed. Goes into `_meta.dropped` (summarized), reachable via `source_ref`.
- **out-of-scope** — not relevant to the consumer. Goes into `_meta.dropped` only if a reader might *expect* it to be there (set expectations); otherwise omit.

When `focus` is supplied, treat it as the salience oracle: content that serves the focus is must-preserve; content that does not is droppable or out-of-scope.

## Preserve-vs-drop rules

PRESERVE:
- Decisions, constraints, obligations, and the conditions that trigger them.
- Identifiers, keys, names, amounts, dates, units, and enumerations — the literal values a downstream stage will key on.
- Edge cases and exceptions (the "except when…" clauses) — these are the first thing a lossy summary drops and the most expensive to lose.
- Anything `focus` names.

DROP (record in `_meta.dropped`):
- Restatements, throat-clearing, and narrative connective tissue.
- Examples that only illustrate a rule already captured.
- Long verbatim spans — replace with a short description + `source_ref` anchor.

## Recording loss (the inline ledger)

`_meta.dropped[]` is the ledger. Each entry is `{ item, reason }`:
- `item` — a short description of what was left out (not the full content).
- `reason` — why: `redundant`, `out-of-scope`, `reconstructable-from-source`, `not-in-source` (conform mode, required field the source lacks), `size-bounded` (pushed behind source_ref).

Set `_meta.coverage` honestly: `high` when all must-preserve content is in the contract; `partial` when some must-preserve content was dropped (each such drop reason-tagged); `low` when the source resisted faithful extraction (say why in `notes`).

## Size discipline (context economy)

The contract is consumed inside a finite context window, sometimes by a small/mid-tier model. Keep it cheap to read:
- Prefer **flat** objects to deep nesting; deep nesting is expensive for downstream stages to traverse and reason over.
- Bound field sizes. If a value would be a large blob, store a short form plus a `source_ref` anchor to the full text — do not inline it.
- Honor `tier_hint`: `small` → flattest, smallest contract; `frontier` → more structure is acceptable.
- One source → one contract. Do not concatenate multiple sources into one envelope.

## Recoverability test

Before emitting, ask: *using only `contract`, `_meta.dropped`, and `_meta.source_ref`, could a reader reconstruct the source's decision-relevant context?* If a dropped item is neither described in `dropped` nor reachable via `source_ref`, it was lost silently — fix it before emitting.
