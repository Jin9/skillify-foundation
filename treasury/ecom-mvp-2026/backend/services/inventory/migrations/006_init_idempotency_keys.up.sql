-- Migration 006: Create idempotency_keys table.
-- Server-side idempotency store for inventory.reservation.create.
-- Key = orderId, endpoint = 'inventory.reservation.create'.
-- TTL = 24h per cross-cutting.idempotency.server_side_store.ttl.

CREATE TABLE IF NOT EXISTS inventory.idempotency_keys (
    key               TEXT        NOT NULL,
    endpoint          TEXT        NOT NULL,
    request_hash      TEXT        NOT NULL,   -- sha256 of canonical {orderId, items[sorted by sku]}
    response_envelope JSONB       NOT NULL,
    http_status       INTEGER     NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at        TIMESTAMPTZ NOT NULL,   -- created_at + 24h

    CONSTRAINT pk_idempotency_keys PRIMARY KEY (key, endpoint)
);

-- idx_idem_expires: enables TTL cleanup job (future).
CREATE INDEX IF NOT EXISTS idx_idem_expires
    ON inventory.idempotency_keys (expires_at);
