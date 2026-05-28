-- Migration 005: outbox_events and consumed_events tables.

-- outbox_events: transactional outbox for events.order.cancelled (the only event emitted).
-- Rows are inserted in the SAME tx as the status mutation; a publisher goroutine
-- polls and ships to Kafka, then sets published_at.
CREATE TABLE IF NOT EXISTS "order".outbox_events (
    id           UUID        NOT NULL DEFAULT gen_random_uuid(),
    aggregate_id UUID        NOT NULL,   -- = orders.id; Kafka partition key
    event_type   TEXT        NOT NULL,   -- MVP: only 'order.cancelled'
    payload_json JSONB       NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    published_at TIMESTAMPTZ NULL,       -- NULL until publisher ships to Kafka

    CONSTRAINT outbox_events_pk PRIMARY KEY (id)
);

-- Partial index for the publisher polling query (unpublished only).
CREATE INDEX IF NOT EXISTS idx_outbox_unpublished
    ON "order".outbox_events (id)
    WHERE published_at IS NULL;

-- consumed_events: per-service idempotency for inbound Kafka events.
-- PK on event_id: conflict on insert = duplicate delivery → ack silently.
CREATE TABLE IF NOT EXISTS "order".consumed_events (
    event_id      UUID        NOT NULL,
    consumer_name TEXT        NOT NULL,
    event_type    TEXT        NOT NULL,
    order_id      UUID        NOT NULL,
    processed_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT consumed_events_pk PRIMARY KEY (event_id)
);

-- Backs debug queries by order_id.
CREATE INDEX IF NOT EXISTS idx_consumed_events_order_id
    ON "order".consumed_events (order_id);
