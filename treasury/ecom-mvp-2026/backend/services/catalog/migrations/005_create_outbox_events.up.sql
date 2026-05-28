-- Migration 005: Create outbox_events table
-- Outbox pattern per cross-cutting.persistence.
-- id = UUID v7 (monotonic); relay scans ORDER BY id ASC.
-- published_at IS NULL = pending; relay marks non-null after Kafka ACK.
-- Partial index on (published_at) WHERE published_at IS NULL keeps the relay scan fast.

CREATE TABLE IF NOT EXISTS catalog.outbox_events (
    id           UUID        PRIMARY KEY,  -- UUID v7 generated in Go
    aggregate_id UUID        NOT NULL,     -- = product_id
    event_type   TEXT        NOT NULL,     -- 'product.created' | 'product.updated'
    payload_json JSONB       NOT NULL,     -- full event payload including eventId, occurredAt, productId, sku
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ NULL          -- NULL = pending; relay marks non-null after Kafka ACK
);

-- Relay scan index: partial, keeps it tiny since only pending rows are scanned
CREATE INDEX IF NOT EXISTS idx_outbox_events_pending
    ON catalog.outbox_events (published_at)
    WHERE published_at IS NULL;
