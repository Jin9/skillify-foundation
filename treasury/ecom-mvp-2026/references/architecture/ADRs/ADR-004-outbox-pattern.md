# ADR-004 — Outbox pattern for reliable event emission

- **Status:** Accepted (LOCKED; carried over from dry-run #1)
- **Date:** 2026-05-08
- **Deciders:** Tech-Lead

## Context

Four services produce domain events that drive cross-service state transitions:
- catalog -> `events.product.created` (seeds inventory.stock_levels rows)
- payment -> `events.payment.completed | failed | expired` (drives Order + Inventory)
- order   -> `events.order.cancelled` (drives Inventory + Payment)
- inventory.sweeper -> `events.reservation.expired` (drives Order)

Naive `producer.Send(...)` directly from inside the business transaction has two failure modes: (a) Kafka unavailable -> business write succeeds but no event ever lands -> downstream is silently inconsistent; (b) Kafka succeeds but the business tx rolls back -> phantom event with no underlying state change.

## Decision

Every producer service writes the event into an `outbox_events` table inside the SAME transaction that mutates the business row, then a separate goroutine relays the row to Kafka.

Schema (per producer service):
```sql
CREATE TABLE outbox_events (
  id           UUID PRIMARY KEY,           -- UUID v7
  aggregate_id TEXT NOT NULL,              -- orderId or sku
  event_type   TEXT NOT NULL,              -- e.g. 'payment.completed'
  payload_json JSONB NOT NULL,             -- canonicalized event body
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  published_at TIMESTAMPTZ NULL
);
CREATE INDEX outbox_events_pending ON outbox_events (id) WHERE published_at IS NULL;
```

Relay loop (one goroutine per service, started at boot):
```
every 5s:
  rows = SELECT id, aggregate_id, event_type, payload_json
         FROM outbox_events
         WHERE published_at IS NULL
         ORDER BY id   -- UUID v7 = roughly time-ordered, preserves per-aggregate order at the producer
         LIMIT 100
  for each row:
     producer.Send(topic_for(event_type), key=aggregate_id, value=payload_json)
  UPDATE outbox_events SET published_at = NOW() WHERE id IN (rows.ids)
```

- **At-least-once delivery:** the relay is single-threaded per service; if a Send succeeds but UPDATE crashes, the row will be re-published on the next tick. Consumers MUST dedup on `eventId` (the `id` column) via per-service `consumed_events(eventId PK, consumerName, processedAt)` with `INSERT ... ON CONFLICT DO NOTHING`.
- **Per-aggregate order preserved at producer:** `ORDER BY id` (UUID v7) is monotonic per the v7 spec; the single-threaded relay sends in id order; Kafka `key=aggregate_id` lands them on one partition; consumer reads in partition order. NO global ordering across aggregates.
- **Cross-topic order NOT preserved:** see ADR-009 (state-driven consumers handle this).

The four producer services are: catalog, inventory, order, payment. Cart, checkout, identity, frontend-web do NOT have outboxes.

REV-L2 dry-run #1 (REV-L2-011) recommended promoting this to a `backend/common/outbox` package so all four implementations don't drift. That refactor is queued as a non-blocking followup; in MVP the four implementations remain independent but each must satisfy:
- single-threaded relay (no parallel publishers per service)
- ORDER BY id ASC
- `published_at` set in a separate UPDATE after `Send` returns success
- panic-recovery on the relay goroutine so one bad row doesn't kill the loop
- per-tick batch ceiling (recommended 100)

## Consequences

- **Positive:** business writes never have a phantom event or a missing event.
- **Positive:** Kafka outage is absorbed for as long as `outbox_events` has disk; the relay drains the backlog when Kafka returns.
- **Positive:** a metric `outbox_pending_count` gives ops a clear backlog signal (alert > 1000 for 10min).
- **Negative:** every consumer must be idempotent on `eventId` AND state-driven (ADR-009).
- **Negative:** four near-duplicate implementations until the common package lands (REV-L2-011 followup).
