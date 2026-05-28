-- Migration 002: Create reservations table.
-- INV-003: one row per SKU per order_id; linked via order_id.
-- INV-004: expires_at drives the sweeper (idx_reservations_sweeper partial index).
-- LOCK-ORDER PIN: stock_levels(sku) MUST be locked BEFORE reservations(id) in all txs.

CREATE TABLE IF NOT EXISTS inventory.reservations (
    id             UUID        NOT NULL,
    order_id       UUID        NOT NULL,
    sku            TEXT        NOT NULL,
    qty            INTEGER     NOT NULL,
    status         TEXT        NOT NULL,
    expires_at     TIMESTAMPTZ NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    released_at    TIMESTAMPTZ,
    release_reason TEXT,

    CONSTRAINT pk_reservations PRIMARY KEY (id),

    CONSTRAINT chk_reservation_qty_positive  CHECK (qty > 0),
    CONSTRAINT chk_reservation_status        CHECK (status IN ('RESERVED','COMMITTED','RELEASED','EXPIRED')),
    CONSTRAINT chk_reservation_release_reason CHECK (
        release_reason IS NULL OR
        release_reason IN (
            'PAYMENT_FAILED','PAYMENT_EXPIRED','ORDER_CANCELLED',
            'ADMIN_FORCE','EXPIRED','PAID_SOFT_CANCEL'
        )
    )
);

-- idx_reservations_order_id: used by event consumers to look up all rows for an order.
CREATE INDEX IF NOT EXISTS idx_reservations_order_id
    ON inventory.reservations (order_id);

-- idx_reservations_sweeper: partial index; drives the sweeper WHERE status='RESERVED' scan cheaply.
CREATE INDEX IF NOT EXISTS idx_reservations_sweeper
    ON inventory.reservations (status, expires_at)
    WHERE status = 'RESERVED';
