-- Migration 004: order_items and order_status_history tables.

-- order_items: immutable price/name/image snapshots per ORD-007.
-- No UPDATE is allowed by the repo layer (Order service enforces this in code).
CREATE TABLE IF NOT EXISTS "order".order_items (
    id                UUID        NOT NULL DEFAULT gen_random_uuid(),
    order_id          UUID        NOT NULL,
    product_id        UUID        NOT NULL,
    sku               TEXT        NOT NULL,
    qty               INT         NOT NULL CHECK (qty > 0),
    name_snapshot     TEXT        NOT NULL,
    image_url_snapshot TEXT       NULL,
    price_snapshot    BIGINT      NOT NULL CHECK (price_snapshot >= 0),
    line_subtotal     BIGINT      NOT NULL CHECK (line_subtotal >= 0),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT order_items_pk PRIMARY KEY (id),
    CONSTRAINT order_items_order_fk FOREIGN KEY (order_id)
        REFERENCES "order".orders (id) ON DELETE RESTRICT
);

-- Backs order.detail JOIN on order_id.
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON "order".order_items (order_id);

-- order_status_history: append-only audit per ORD-009.
-- Written in the SAME transaction as the orders.status UPDATE.
CREATE TABLE IF NOT EXISTS "order".order_status_history (
    id           UUID        NOT NULL DEFAULT gen_random_uuid(),
    order_id     UUID        NOT NULL,
    from_status  TEXT        NULL,   -- NULL only for the initial PENDING_PAYMENT row
    to_status    TEXT        NOT NULL,
    actor_user_id UUID       NULL,   -- NULL when actor_role=SYSTEM
    actor_role   TEXT        NOT NULL,
    reason       TEXT        NULL,   -- required for CANCELLED (admin path) and SYSTEM events
    occurred_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT order_status_history_pk PRIMARY KEY (id),
    CONSTRAINT order_status_history_order_fk FOREIGN KEY (order_id)
        REFERENCES "order".orders (id),
    CONSTRAINT order_status_history_actor_role_check CHECK (actor_role IN ('CUSTOMER','ADMIN','SYSTEM'))
);

-- Backs order.detail.statusHistory ordered by occurred_at ASC.
CREATE INDEX IF NOT EXISTS idx_order_status_history_order_occurred
    ON "order".order_status_history (order_id, occurred_at ASC);
