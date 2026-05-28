-- Migration: 001 — create cart schema and carts table
-- Schema: cart
-- One cart per customer (UNIQUE user_id). cart_id UUID v7 PK.

CREATE SCHEMA IF NOT EXISTS cart;

CREATE TABLE IF NOT EXISTS cart.carts (
    cart_id    UUID        NOT NULL DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT pk_carts PRIMARY KEY (cart_id),
    CONSTRAINT uq_carts_user_id UNIQUE (user_id)
);

CREATE INDEX IF NOT EXISTS idx_carts_user_id ON cart.carts (user_id);

COMMENT ON TABLE  cart.carts          IS 'One cart per customer. Keyed by user_id from JWT claims.sub.';
COMMENT ON COLUMN cart.carts.user_id  IS 'UUID of the owning customer (from identity service, JWT claims.sub).';
