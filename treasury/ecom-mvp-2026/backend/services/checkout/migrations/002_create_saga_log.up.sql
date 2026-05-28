-- Migration 002: Create saga_log table for operational forensics.
-- Write-only; NOT on the request path. Retained 90 days.

CREATE TABLE checkout.saga_log (
    saga_id              UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id             UUID         NOT NULL,
    customer_user_id     UUID         NOT NULL,
    status               VARCHAR(16)  NOT NULL CHECK (status IN ('COMPLETED', 'COMPENSATED', 'DEGRADED')),
    steps                JSONB        NOT NULL,   -- [{step, name, startedAt, endedAt, outcome}]
    compensation_outcome JSONB        NULL,       -- NULL unless status IN ('COMPENSATED','DEGRADED')
    trace_id             VARCHAR(64)  NOT NULL,
    created_at           TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Index for lookup by orderId during postmortems
CREATE INDEX idx_saga_order_id
    ON checkout.saga_log (order_id);

-- Partial index for ops alerting on DEGRADED entries
CREATE INDEX idx_saga_status_created
    ON checkout.saga_log (status, created_at)
    WHERE status = 'DEGRADED';
