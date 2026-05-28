-- Migration 005: Create consumed_events dedup table.
-- Per cross-cutting.persistence pattern for idempotent Kafka consumers.
-- INSERT inside the same tx as the event side-effect;
-- ON CONFLICT (event_id, consumer_name) → already processed, treat as no-op.

CREATE TABLE IF NOT EXISTS inventory.consumed_events (
    event_id      UUID        NOT NULL,
    consumer_name TEXT        NOT NULL,
    processed_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT pk_consumed_events PRIMARY KEY (event_id, consumer_name)
);
