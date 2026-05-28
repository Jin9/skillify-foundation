-- Migration 001: Create catalog schema and categories table
-- Schema: catalog
-- PERF-001: partial indexes on active/level for storefront queries

CREATE SCHEMA IF NOT EXISTS catalog;

CREATE TABLE IF NOT EXISTS catalog.categories (
    id                UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name              TEXT         NOT NULL,
    slug              TEXT         NOT NULL,
    parent_category_id UUID        NULL REFERENCES catalog.categories(id),
    level             INT          NOT NULL,
    active            BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Globally unique slug
CREATE UNIQUE INDEX IF NOT EXISTS categories_slug_key
    ON catalog.categories (slug);

-- Name unique within the same parent (non-root)
CREATE UNIQUE INDEX IF NOT EXISTS categories_parent_category_id_name_key
    ON catalog.categories (parent_category_id, name)
    WHERE parent_category_id IS NOT NULL;

-- Name unique at root level (parent IS NULL)
CREATE UNIQUE INDEX IF NOT EXISTS categories_name_root_key
    ON catalog.categories (name)
    WHERE parent_category_id IS NULL;

-- Index for storefront category listing
CREATE INDEX IF NOT EXISTS idx_categories_active_level
    ON catalog.categories (active, level);
