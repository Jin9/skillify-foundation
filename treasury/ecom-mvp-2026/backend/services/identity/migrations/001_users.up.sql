BEGIN;

CREATE SCHEMA IF NOT EXISTS identity;

SET search_path TO identity;

CREATE TABLE IF NOT EXISTS users (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    email         TEXT        UNIQUE NOT NULL,
    password_hash TEXT        NOT NULL,
    name          TEXT        NOT NULL,
    phone         TEXT        NULL,
    role          TEXT        NOT NULL CHECK (role IN ('CUSTOMER','ADMIN')),
    status        TEXT        NOT NULL CHECK (status IN ('ACTIVE','SUSPENDED')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
