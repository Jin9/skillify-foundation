# SQL & Migrations

House defaults — when the host repo's stated conventions (AGENTS.md or
equivalent) conflict with a rule here, the host repo wins; note the
divergence in the hand-off summary.

## Contents
- Migration rules
- Query conventions

## Migration rules

- One change per migration file. NEVER bundle unrelated schema changes.
- Forward-only: no `DOWN` migrations in production. Write a new `UP` to reverse.
- Every table MUST have `created_at` and `updated_at` (server-default timestamps).
- Soft-delete (`deleted_at`) only when audit or compliance requires it. Prefer hard delete otherwise.

## Query conventions

- NEVER `SELECT *`. Always list explicit columns.
- Justify every index in the migration comment (which queries it serves).
- Use `EXPLAIN ANALYZE` to validate index usage before merging.
- Prefer CTEs over nested subqueries for readability. Benchmark if performance-critical.
