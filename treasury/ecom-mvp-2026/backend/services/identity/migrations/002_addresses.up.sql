BEGIN;

SET search_path TO identity;

CREATE TABLE IF NOT EXISTS addresses (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID        NOT NULL REFERENCES users(id),
    receiver_name TEXT        NOT NULL,
    phone         TEXT        NOT NULL,
    address_line1 TEXT        NOT NULL,
    address_line2 TEXT        NULL,
    province      TEXT        NOT NULL,
    district      TEXT        NOT NULL,
    sub_district  TEXT        NULL,
    postal_code   TEXT        NOT NULL,
    is_default    BOOLEAN     NOT NULL DEFAULT false,
    deleted_at    TIMESTAMPTZ NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_addresses_user_id_active
    ON addresses (user_id)
    WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS uidx_addresses_user_default
    ON addresses (user_id)
    WHERE is_default = true AND deleted_at IS NULL;

COMMIT;
