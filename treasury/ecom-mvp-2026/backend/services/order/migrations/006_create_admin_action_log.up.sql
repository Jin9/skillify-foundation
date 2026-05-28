-- Migration 006: admin_action_log table.
-- Local audit fallback per td.json §admin_action_log. Audit module is out-of-scope
-- (BA out_of_scope line 76) but admin endpoints write a local row in-tx for traceability.

CREATE TABLE IF NOT EXISTS "order".admin_action_log (
    id               UUID        NOT NULL DEFAULT gen_random_uuid(),
    actor_user_id    UUID        NOT NULL,
    endpoint         TEXT        NOT NULL,   -- e.g. 'order.update-status-admin'
    order_id         UUID        NULL,
    request_summary  JSONB       NOT NULL,   -- {toStatus, trackingNumber?, reason?} — no JWT/auth headers
    outcome_code     TEXT        NOT NULL,   -- from cross-cutting.error-codes registry
    occurred_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT admin_action_log_pk PRIMARY KEY (id)
);

-- Backs lookups by actor or order.
CREATE INDEX IF NOT EXISTS idx_admin_action_log_actor
    ON "order".admin_action_log (actor_user_id, occurred_at DESC);

CREATE INDEX IF NOT EXISTS idx_admin_action_log_order
    ON "order".admin_action_log (order_id)
    WHERE order_id IS NOT NULL;
