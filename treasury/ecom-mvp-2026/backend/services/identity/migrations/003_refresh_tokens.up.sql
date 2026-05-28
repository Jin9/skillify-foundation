BEGIN;

SET search_path TO identity;

CREATE TABLE IF NOT EXISTS refresh_tokens (
    jti             UUID        PRIMARY KEY,
    user_id         UUID        NOT NULL REFERENCES users(id),
    issued_at       TIMESTAMPTZ NOT NULL,
    expires_at      TIMESTAMPTZ NOT NULL,
    revoked_at      TIMESTAMPTZ NULL,
    replaced_by_jti UUID        NULL
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_expires
    ON refresh_tokens (user_id, expires_at);

COMMIT;
