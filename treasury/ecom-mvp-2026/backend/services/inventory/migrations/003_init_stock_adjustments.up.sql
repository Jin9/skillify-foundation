-- Migration 003: Create stock_adjustments audit table.
-- INV-007: written in the same tx as inventory.stock.adjust mutations.

CREATE TABLE IF NOT EXISTS inventory.stock_adjustments (
    id            UUID        NOT NULL,
    sku           TEXT        NOT NULL,
    delta         INTEGER     NOT NULL,
    reason        TEXT        NOT NULL,
    actor_user_id UUID        NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT pk_stock_adjustments PRIMARY KEY (id),
    CONSTRAINT chk_adjustment_delta_nonzero CHECK (delta <> 0)
);

-- idx_stock_adjustments_sku: supports audit queries by SKU.
CREATE INDEX IF NOT EXISTS idx_stock_adjustments_sku
    ON inventory.stock_adjustments (sku, created_at DESC);
