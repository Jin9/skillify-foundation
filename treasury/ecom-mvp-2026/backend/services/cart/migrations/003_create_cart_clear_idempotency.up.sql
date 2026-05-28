-- Migration: 003 — idempotency anchor for cart.clear-on-checkout (CHK-010)
-- Keyed on order_id; ON CONFLICT DO NOTHING prevents double-processing on replay.

CREATE TABLE IF NOT EXISTS cart.cart_clear_idempotency (
    order_id   UUID        NOT NULL,
    removed    INT         NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT pk_cart_clear_idempotency PRIMARY KEY (order_id)
);

COMMENT ON TABLE  cart.cart_clear_idempotency          IS 'Idempotency table for cart.clear-on-checkout. Keyed by orderId (CHK-010).';
COMMENT ON COLUMN cart.cart_clear_idempotency.removed  IS 'Number of cart_items rows deleted in the original call. Returned verbatim on replay.';
