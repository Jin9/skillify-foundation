-- Migration 001: Create inventory schema and stock_levels table.
-- INV-001 / INV-008: no-negative invariant enforced by CHECK constraints (belt-and-braces;
-- the service layer also verifies before UPDATE under FOR UPDATE row lock).

CREATE SCHEMA IF NOT EXISTS inventory;

CREATE TABLE IF NOT EXISTS inventory.stock_levels (
    sku           TEXT        NOT NULL,
    available_qty INTEGER     NOT NULL DEFAULT 0,
    reserved_qty  INTEGER     NOT NULL DEFAULT 0,
    sold_qty      INTEGER     NOT NULL DEFAULT 0,
    version       INTEGER     NOT NULL DEFAULT 0,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT pk_stock_levels PRIMARY KEY (sku),

    -- No-negative invariant: DB-level defense (belt-and-braces).
    -- Service layer also checks under FOR UPDATE before UPDATE.
    CONSTRAINT chk_stock_available_non_negative CHECK (available_qty >= 0),
    CONSTRAINT chk_stock_reserved_non_negative  CHECK (reserved_qty  >= 0),
    CONSTRAINT chk_stock_sold_non_negative      CHECK (sold_qty      >= 0)
);
