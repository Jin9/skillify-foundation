-- Migration 003: payment.outbox table
-- Transactional outbox for events.payment.completed | events.payment.failed.
-- Written in the SAME tx as the intent state transition.
-- Drained by the outbox publisher goroutine to Kafka topic ecom.payment.events.
-- Partition key = order_id (preserves per-order event ordering).

CREATE TABLE payment.outbox (
    id             UUID        PRIMARY KEY,
    aggregate_type TEXT        NOT NULL DEFAULT 'payment_intent',
    aggregate_id   UUID        NOT NULL,
    event_id       UUID        NOT NULL,
    event_type     TEXT        NOT NULL
                               CHECK (event_type IN ('payment.completed','payment.failed','payment.expired')),
    topic          TEXT        NOT NULL DEFAULT 'ecom.payment.events',
    partition_key  TEXT        NOT NULL,
    payload        JSONB       NOT NULL,
    status         TEXT        NOT NULL DEFAULT 'PENDING'
                               CHECK (status IN ('PENDING','PUBLISHED','FAILED')),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    published_at   TIMESTAMPTZ NULL,
    attempts       INT         NOT NULL DEFAULT 0,
    last_error     TEXT        NULL,

    CONSTRAINT uq_outbox_event_id UNIQUE (event_id)
);

-- Backs publisher polling for PENDING rows ordered by created_at.
CREATE INDEX idx_outbox_status_created_at
    ON payment.outbox (status, created_at)
    WHERE status = 'PENDING';
