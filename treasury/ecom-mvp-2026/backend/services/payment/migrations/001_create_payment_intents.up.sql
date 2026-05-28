-- Migration 001: payment_intents table
-- Schema: payment
-- Lifecycle: status REQUIRES_PAYMENT → SUCCEEDED | FAILED | EXPIRED

CREATE SCHEMA IF NOT EXISTS payment;

CREATE TABLE payment.payment_intents (
    intent_id        UUID        PRIMARY KEY,
    order_id         UUID        NOT NULL,
    owner_user_id    UUID        NOT NULL,
    amount_minor     BIGINT      NOT NULL CHECK (amount_minor > 0),
    currency         CHAR(3)     NOT NULL DEFAULT 'THB',
    status           TEXT        NOT NULL DEFAULT 'REQUIRES_PAYMENT'
                                 CHECK (status IN ('REQUIRES_PAYMENT','SUCCEEDED','FAILED','EXPIRED')),
    mock_provider_ref TEXT       NULL,
    provider_status  TEXT        NULL,
    paid_at          TIMESTAMPTZ NULL,
    expires_at       TIMESTAMPTZ NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    version          INT         NOT NULL DEFAULT 0,

    CONSTRAINT uq_payment_intents_order_id UNIQUE (order_id)
);

-- Backs the expiry sweeper: REQUIRES_PAYMENT rows past expires_at.
CREATE INDEX idx_payment_intents_status_expires_at
    ON payment.payment_intents (status, expires_at)
    WHERE status = 'REQUIRES_PAYMENT';
