-- Migration 004: Create outbox_events table.
-- Per cross-cutting.persistence.outbox_pattern.
-- Relayer goroutine polls WHERE published_at IS NULL ORDER BY id LIMIT 100,
-- publishes to Kafka (ecom.inventory.events), marks published_at.
-- At-least-once delivery.

CREATE TABLE IF NOT EXISTS inventory.outbox_events (
    id           UUID        NOT NULL,
    aggregate_id TEXT        NOT NULL,  -- = orderId for reservation.expired (partition key)
    event_type   TEXT        NOT NULL,
    payload_json JSONB       NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ,

    CONSTRAINT pk_outbox_events PRIMARY KEY (id)
);

-- idx_outbox_unpublished: partial index used by the relay goroutine poll.
CREATE INDEX IF NOT EXISTS idx_outbox_unpublished
    ON inventory.outbox_events (id)
    WHERE published_at IS NULL;
