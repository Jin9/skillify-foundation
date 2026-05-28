-- Migration 004: Create product_status_history table
-- Audit-light pattern per CAT-008.
-- Full AUD module is out-of-scope MVP.
-- from_status IS NULL indicates initial product creation (no previous status).

CREATE TABLE IF NOT EXISTS catalog.product_status_history (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id    UUID        NOT NULL REFERENCES catalog.products(id),
    from_status   TEXT        NULL,
    to_status     TEXT        NOT NULL,
    changed_field TEXT        NULL, -- NULL = status-only change
    old_value     JSONB       NULL,
    new_value     JSONB       NULL,
    actor_user_id UUID        NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_product_status_history_product_created
    ON catalog.product_status_history (product_id, created_at DESC);
