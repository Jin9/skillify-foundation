-- Migration 002: Create products table
-- price stored as NUMERIC(12,2) — NOT INT (TD pattern review requirement)
-- Partial indexes exclude DELETED/DRAFT/INACTIVE rows from storefront queries (PERF-001)

CREATE TABLE IF NOT EXISTS catalog.products (
    id          UUID            PRIMARY KEY,
    sku         TEXT            NOT NULL,
    name        TEXT            NOT NULL,
    description TEXT            NOT NULL DEFAULT '',
    price       NUMERIC(12, 2)  NOT NULL CHECK (price > 0),
    category_id UUID            NOT NULL REFERENCES catalog.categories(id),
    status      TEXT            NOT NULL CHECK (status IN ('DRAFT', 'ACTIVE', 'INACTIVE', 'DELETED')),
    visible     BOOLEAN         NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

-- Case-insensitive SKU uniqueness
CREATE UNIQUE INDEX IF NOT EXISTS products_sku_lower_key
    ON catalog.products (LOWER(sku));

-- Listing index: status + visible + category_id — partial (ACTIVE+visible only)
CREATE INDEX IF NOT EXISTS idx_products_listing
    ON catalog.products (status, visible, category_id)
    WHERE status = 'ACTIVE' AND visible = TRUE;

-- Price-range filter index — partial (ACTIVE+visible only)
CREATE INDEX IF NOT EXISTS idx_products_price
    ON catalog.products (status, visible, price)
    WHERE status = 'ACTIVE' AND visible = TRUE;

-- Sort-by-created index — partial (ACTIVE+visible only)
CREATE INDEX IF NOT EXISTS idx_products_created
    ON catalog.products (status, visible, created_at)
    WHERE status = 'ACTIVE' AND visible = TRUE;
