-- Migration 003: orders table.
-- Aggregate root. One row per customer order.
-- Order is the SOLE WRITER of the status column (sole_writer_invariant per td.json).
-- BIGINT for all monetary fields (whole-THB integers per cross-cutting.pricing currency rule).

CREATE TABLE IF NOT EXISTS "order".orders (
    id                   UUID        NOT NULL DEFAULT gen_random_uuid(),
    order_number         TEXT        NOT NULL,
    user_id              UUID        NOT NULL,
    status               TEXT        NOT NULL,
    subtotal             BIGINT      NOT NULL CHECK (subtotal >= 0),
    shipping_fee         BIGINT      NOT NULL CHECK (shipping_fee >= 0),
    coupon_discount      BIGINT      NOT NULL DEFAULT 0 CHECK (coupon_discount = 0), -- reserved; coupon engine out of scope
    grand_total          BIGINT      NOT NULL CHECK (grand_total >= 0),
    currency             TEXT        NOT NULL DEFAULT 'THB' CHECK (currency = 'THB'),
    address_snapshot     JSONB       NOT NULL,
    buyer_email_snapshot TEXT        NOT NULL,
    tracking_number      TEXT        NULL,
    idempotency_key      TEXT        NULL,
    version              INT         NOT NULL DEFAULT 0,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT orders_pk PRIMARY KEY (id),
    CONSTRAINT orders_status_check CHECK (status IN (
        'PENDING_PAYMENT','PAID','PAYMENT_FAILED','PAYMENT_EXPIRED',
        'PACKING','SHIPPED','DELIVERED','CANCELLED'
    ))
);

-- Human-readable order numbers must be unique.
CREATE UNIQUE INDEX IF NOT EXISTS idx_orders_order_number ON "order".orders (order_number);

-- Backs order.list-mine (sort createdAt DESC, scoped by user_id).
CREATE INDEX IF NOT EXISTS idx_orders_user_id_created_at_desc
    ON "order".orders (user_id, created_at DESC);

-- Backs order.list-admin status filter.
CREATE INDEX IF NOT EXISTS idx_orders_status_created_at_desc
    ON "order".orders (status, created_at DESC);

-- Backs order.list-admin search-by-email (buyer_email_snapshot).
CREATE INDEX IF NOT EXISTS idx_orders_buyer_email
    ON "order".orders (buyer_email_snapshot);

-- Backs replay/debug correlation via idempotency_key.
CREATE INDEX IF NOT EXISTS idx_orders_idempotency_key
    ON "order".orders (idempotency_key);
