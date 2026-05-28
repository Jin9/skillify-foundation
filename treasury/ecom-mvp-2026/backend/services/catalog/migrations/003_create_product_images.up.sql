-- Migration 003: Create product_images table
-- Images stored as URL strings. Image set replaced atomically (DELETE+INSERT) during product updates.

CREATE TABLE IF NOT EXISTS catalog.product_images (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID        NOT NULL REFERENCES catalog.products(id) ON DELETE CASCADE,
    url        TEXT        NOT NULL,
    sort_order INT         NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_product_images_product_sort
    ON catalog.product_images (product_id, sort_order);
