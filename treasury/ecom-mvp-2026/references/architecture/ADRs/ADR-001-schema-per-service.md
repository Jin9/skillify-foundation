# ADR-001 — Schema-per-service for Postgres

- **Status:** Accepted (locked decision; carried over from dry-run #1)
- **Date:** 2026-05-08
- **Deciders:** Tech-Lead (Claude Opus 4.7 1M, dry-run #2)
- **Context:** B2C E-Commerce Platform MVP — 7 backend services on PostgreSQL via `common/database` (pgx).

## Context

The orchestrator locked PostgreSQL as the only persistence engine. The seven services own different aggregates (users, products, stock_levels, carts, idempotency_keys, orders, payment_intents) and the BA spec's failure_tolerance section requires that no service may read another's tables for failover.

Two viable layouts existed:
1. One Postgres database per service (the cleanest physical isolation).
2. One Postgres instance with one schema per service (cheap to operate locally; same access-pattern guarantees if no cross-schema joins are written).

## Decision

Schema-per-service is the canonical layout. In local Docker Compose this is one Postgres instance with one schema per service so the developer doesn't need 7 PG containers; in SIT/UAT/PRD the operator MAY split into separate Postgres instances per service without any code change because no cross-schema joins exist in any service's source.

NO cross-service tables. NO shared tables (no shared `users`, no shared `audit_log`). Each service owns its tables and its migrations. Cross-service data flows through HTTP contracts and Kafka events only.

Schemas:
- identity_db: users, addresses, refresh_tokens
- catalog_db: products, categories, product_images, product_status_history, outbox_events
- inventory_db: stock_levels, reservations, stock_adjustments, outbox_events, consumed_events, idempotency_keys
- cart_db: carts, cart_items
- checkout_db: idempotency_keys, saga_log
- order_db: orders, order_items, order_status_history, order_addresses, outbox_events, consumed_events, admin_action_log
- payment_db: payment_intents, payment_records, payment_callback_dedup, outbox_events

## Consequences

- **Positive:** services can be deployed/scaled/migrated independently; no cross-service lock contention; per-service backup/restore granularity.
- **Positive:** local cost is one PG container, not seven.
- **Negative:** denormalized snapshots (order_items.priceSnapshot, order_items.buyerEmailSnapshot) must be written at commit-time by the owning service because no JOIN is permitted later.
- **Negative:** ad-hoc cross-domain reporting requires a downstream warehouse or per-service queries — not in MVP scope.
