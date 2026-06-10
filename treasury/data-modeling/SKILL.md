---
name: data-modeling
description: >
  Design the persistence model safely (schema, source of truth, indexes, constraints) with an expand/contract migration and rollback plan that preserves backward compatibility. Use when the user asks "design the table/schema", "what indexes do we need", "plan this migration", or "where's the source of truth". Produces a Data Model plus Migration Plan artifact, a checklist, and a human approval gate before any migration runs. Do NOT use for query-cost or performance tuning of an existing system or chasing a specific slow path (use progressive-bug-hunter), or for domain boundaries, aggregates, and ownership (use domain-modeling).
---

# data-modeling

## Purpose
Design the persistence model safely — schema, source of truth, indexes, constraints — with a migration and rollback
plan that preserves backward compatibility. This is *persistence* (tables, indexes, migrations), not business
semantics.

## When to use
- Triggers: *"design the table/schema"*, *"what indexes do we need"*, *"plan this migration"*, *"where's the source of truth"*.
- **Not this skill:** query-cost/performance tuning of an existing system → a dedicated performance-cost review, or `progressive-bug-hunter` for a specific slow path.

## Input
- A domain model + access patterns + expected volume/growth (NFRs).

## Output
A **Data Model + Migration Plan** artifact + checklist. Skeleton:

```
# Data Model — <area>
Tables/collections:  <name> — columns (type, null?, default)
Source of truth:     <which store owns which fact>
Indexes:             <index> — serves <query pattern>
Constraints:         <unique / fk / check>
Migration (expand → migrate → contract):
  1. expand …  2. backfill …  3. switch …  4. contract …
Rollback:            <how to reverse each step>
```

## Decision rules
1. Derive **indexes from real access patterns + volume**, not guesses.
2. **One source of truth per fact**; keep domain rules out of the schema — domain/infra separation.
3. **Never migrate without a rollback plan** — deprecating/changing data without one is a forbidden move.
4. Prefer **reversible expand/contract** migrations; a data-migration *strategy* is a one-way door — decide-then-explain.
5. If only one person understands the schema, the **bus factor is too low** — document it.

## Checklist
- [ ] Tables/columns + types defined
- [ ] Source of truth designated per fact
- [ ] Indexes mapped to query patterns
- [ ] Constraints (unique/fk/check) set
- [ ] Migration steps (expand/contract) written
- [ ] Rollback for each step
- [ ] Volume/growth accounted for

## Anti-patterns (never do)
- Run a migration with no rollback.
- Delete or deprecate a table/service without a migration plan.
- Add indexes by guess instead of access pattern.
- Encode domain logic into the schema.
- Leave schema knowledge with a single person.

## Example
**Input:** wallet ledger. **Output (excerpt):** `ledger(entry_id, wallet_id, amount, type, idempotency_key UNIQUE,
created_at)`, append-only; source of truth for balance = sum of ledger; index `(wallet_id, created_at)`. Migration:
expand new column nullable → backfill → enforce. Rollback: drop column (no data loss).

## Human approval gate
**Stop.** A human approves the migration **and** its rollback before it runs — data migration is a one-way door.
The skill designs and plans; it does not execute migrations.
