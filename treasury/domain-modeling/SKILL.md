---
name: domain-modeling
description: >
  Defines domain boundaries and DDD building blocks so ownership is clear, the model matches the business, and domain logic stays separate from infrastructure. Use when the user asks "define the domain model", "where are the bounded contexts/boundaries", "what are the aggregates and events", or "who owns this part of the business". Produces a Domain Model artifact with ubiquitous language, bounded contexts and owners, aggregates and invariants, commands versus events, and an explicit domain/infra separation note. Do NOT use for reverse-engineering the model from existing code (that belongs to business-logic-extractor), or for persistence schema, index, and migration design (use data-modeling).
---

# domain-modeling

## Purpose
Define the domain boundaries and DDD building blocks so ownership is clear, the model matches the business, and domain
logic stays separate from infrastructure.

## When to use
- Triggers: *"define the domain model"*, *"where are the bounded contexts/boundaries"*, *"what are the aggregates and events"*, *"who owns this part of the business"*.
- **Not this skill:** reverse-engineering the model from existing code → `business-logic-extractor`.

## Input
- Requirement + business context (+ existing ubiquitous language, if any).

## Output
A **Domain Model** artifact + checklist. Skeleton:

```
# Domain Model — <area>
Ubiquitous language:  <term = meaning>
Bounded contexts:  <context> — responsibility — owner
Per aggregate:
  Aggregate:  <name> (root)
  Invariants:  <rules that must always hold>
  Entities / Value objects:
  Commands (intent):   Events (fact):
Domain ↔ infra separation note:
```

## Decision rules
1. **Always separate domain logic from infrastructure** — this is an iron rule.
2. Scope by **business domain**; define the responsibility and authority of each boundary.
3. **Conway's Law** → align services to domains; apply DDD + CQRS where they fit.
4. Get **aggregate boundaries** right early: wrong domain design can spread tech debt across the whole system and force a full re-scaffold. The path that works: DDD → aggregate boundary → standard template → self-sufficient team.
5. Distinguish **commands** (intent, may be rejected) from **events** (facts, already happened).

## Checklist
- [ ] Ubiquitous language captured
- [ ] Bounded contexts named with responsibility + owner
- [ ] Aggregates with their invariants
- [ ] Entities vs value objects distinguished
- [ ] Commands and events listed per context
- [ ] Domain/infra separation stated explicitly

## Anti-patterns (never do)
- Leak infrastructure concerns into the domain model.
- Leave aggregate boundaries fuzzy — a leading cause of severe domain tech debt.
- Model one giant blob with no boundaries.
- Skip the ubiquitous language (juniors then ask the same thing 3× → docs missing).

## Example
**Input:** wallet top-up feature. **Output (excerpt):** *Contexts:* Wallet (owns balance), Ledger (owns entries).
*Aggregate:* Wallet — invariant `balance ≥ 0`. *Commands:* `TopUpWallet`. *Events:* `WalletToppedUp`. *Separation:*
balance rules live in the domain; the payment-provider client is infra.

## Human approval gate
**Stop.** The team validates the boundaries and ownership before the model hardens — these are expensive to change
later, so treat them like a one-way door.
