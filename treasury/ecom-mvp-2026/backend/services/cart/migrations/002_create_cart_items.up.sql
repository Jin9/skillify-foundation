-- Migration: 002 — create cart_items table
-- UNIQUE(cart_id, product_id) enforces the same-SKU merge invariant (CART-006).
-- No cross-service FK on product_id — product details fetched live from catalog at read-time.

CREATE TABLE IF NOT EXISTS cart.cart_items (
    cart_item_id UUID        NOT NULL DEFAULT gen_random_uuid(),
    cart_id      UUID        NOT NULL,
    product_id   UUID        NOT NULL,
    qty          INT         NOT NULL CHECK (qty > 0),
    added_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT pk_cart_items         PRIMARY KEY (cart_item_id),
    CONSTRAINT fk_cart_items_cart_id FOREIGN KEY (cart_id) REFERENCES cart.carts (cart_id) ON DELETE CASCADE,
    CONSTRAINT uq_cart_items_sku     UNIQUE (cart_id, product_id)
);

CREATE INDEX IF NOT EXISTS idx_cart_items_cart_id ON cart.cart_items (cart_id);

COMMENT ON TABLE  cart.cart_items            IS 'Line items in a cart. Merged on same SKU via ON CONFLICT (CART-006).';
COMMENT ON COLUMN cart.cart_items.product_id IS 'Stored reference only; no cross-service FK. Details fetched from catalog at read-time.';
COMMENT ON COLUMN cart.cart_items.qty        IS 'Positive integer. qty=0 is handled at service layer before write (remove path).';
