-- Migration 002: payment_callback_dedup table
-- Strict webhook dedup keyed by sha256(intent_id || '|' || provider_status).
-- Replay returns the cached envelope without re-firing outbox.

CREATE TABLE payment.payment_callback_dedup (
    dedup_key        TEXT        PRIMARY KEY,
    intent_id        UUID        NOT NULL REFERENCES payment.payment_intents (intent_id),
    provider_status  TEXT        NOT NULL,
    envelope         JSONB       NOT NULL,
    http_status      INT         NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at       TIMESTAMPTZ NOT NULL DEFAULT now() + INTERVAL '30 days'
);

-- Backs intent-scoped queries (not required for critical path; useful for debugging).
CREATE INDEX idx_payment_callback_dedup_intent_id
    ON payment.payment_callback_dedup (intent_id);

-- Backs retention sweep of expired dedup rows.
CREATE INDEX idx_payment_callback_dedup_expires_at
    ON payment.payment_callback_dedup (expires_at);
