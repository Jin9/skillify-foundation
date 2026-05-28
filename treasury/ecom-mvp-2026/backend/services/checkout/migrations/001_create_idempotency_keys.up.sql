-- Migration 001: Create checkout schema and idempotency_keys table.
-- Schema: checkout (isolated per cross-cutting.persistence; no cross-service joins).

CREATE SCHEMA IF NOT EXISTS checkout;

CREATE TABLE checkout.idempotency_keys (
    key              VARCHAR(128)     NOT NULL,
    customer_user_id UUID             NOT NULL,
    request_hash     CHAR(64)         NOT NULL,   -- sha256 hex of canonicalized body
    status           VARCHAR(16)      NOT NULL CHECK (status IN ('INFLIGHT', 'COMPLETED', 'ABANDONED')),
    response_envelope JSONB           NULL,       -- NULL until status = COMPLETED | ABANDONED
    http_status      SMALLINT         NULL,       -- NULL until status = COMPLETED | ABANDONED
    created_at       TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    expires_at       TIMESTAMPTZ      NOT NULL,   -- NOW() + 24h per cross-cutting.idempotency

    PRIMARY KEY (key, customer_user_id)
);

-- Index for TTL cleanup janitor (delete COMPLETED/ABANDONED after 24h)
CREATE INDEX idx_idempotency_expires_at
    ON checkout.idempotency_keys (expires_at)
    WHERE status IN ('COMPLETED', 'ABANDONED');

-- Index for stale-INFLIGHT janitor (sweep INFLIGHT > 60s old)
CREATE INDEX idx_idempotency_inflight
    ON checkout.idempotency_keys (created_at)
    WHERE status = 'INFLIGHT';
