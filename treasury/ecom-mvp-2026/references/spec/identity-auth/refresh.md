# POST /api/v1/identity/auth/refresh

## Summary

Rotates the refresh token in a single atomic Postgres transaction (ID-DEC-1). The old refresh-token jti is revoked and linked via `replaced_by_jti` to the new jti, then a new ES256 JWT pair (access + refresh) is signed and returned. Any replay of a revoked jti returns `AUTH_REVOKED` immediately (no soft-idempotency window at MVP per ID-DEC-2). Implements AUTH-003.

## Story refs

- `STORY_AUTH_REGISTER` (AC-07, AC-08)

## Contract ref

[`identity.refresh`](../../architecture/contracts.json) — `contract_version: 0.1.0`

## Auth

- Tier: **none** (public endpoint — presents refresh token in request body, not Authorization header).

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["refreshToken"],
  "properties": {
    "refreshToken": { "type": "string", "minLength": 1 }
  }
}
```

## Response (success)

HTTP 200:

```json
{
  "code": "SUCCESS",
  "message": "Token refreshed",
  "data": {
    "accessToken":           "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9...",
    "accessTokenExpiresAt":  "2026-05-08T10:15:00Z",
    "refreshToken":          "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshTokenExpiresAt": "2026-05-22T10:00:00Z"
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `VALIDATION_ERROR` | 400 | `refreshToken` missing or empty |
| `AUTH_INVALID`     | 401 | JWT signature invalid, expired, iss/aud mismatch, or `tokenType` claim is not `REFRESH` |
| `AUTH_REVOKED`     | 401 | jti not found in `refresh_tokens` OR `revoked_at IS NOT NULL` (covers replay and logout-then-refresh) |

## Business logic steps

1. Bind and validate request (`common/validator`); return `VALIDATION_ERROR` if `refreshToken` is missing or empty.
2. Parse and verify JWT via `common/token.ParseES256`; validate `exp`, `iss`, `aud`, and `tokenType == "REFRESH"` claims. Return `AUTH_INVALID` on any verification failure.
3. Extract `old_jti` (JWT `jti` claim) and `user_id` (JWT `sub` claim).
4. Generate `new_jti = uuid.NewV4()`.
5. Open a `pgx` transaction.
6. `UPDATE refresh_tokens SET revoked_at = NOW(), replaced_by_jti = $new_jti WHERE jti = $old_jti AND revoked_at IS NULL`. Check `RowsAffected`.
7. If `RowsAffected == 0`: `ROLLBACK`; return `AUTH_REVOKED` (covers both never-existed and already-revoked tokens).
8. `INSERT INTO refresh_tokens(jti = $new_jti, user_id, issued_at = NOW(), expires_at = NOW() + 14 days)`.
9. `COMMIT` the transaction.
10. After successful commit, sign new access token (ES256, 15 min) and new refresh token (ES256, 14 days) via `common/token.SignES256`.
11. Log `old_jti → new_jti` rotation at INFO level (never log token bytes).
12. Return `common/wrapper.Success(ctx, data)` with the new token pair.

## Side effects

- Single Postgres transaction mutates two rows in `identity.refresh_tokens`: revokes old jti (sets `revoked_at`, `replaced_by_jti`), inserts new jti row.
- Audit chain maintained via `replaced_by_jti` column.
- No outbox event.
- Token signing is performed after `COMMIT` to prevent returning tokens when the DB write rolled back.

## Idempotency

Non-idempotent (ID-DEC-2). The 5-second soft-idempotency window is deferred to v2. Replaying a revoked jti returns `AUTH_REVOKED` immediately. Frontend route handler may call `/login` to recover.

## Performance

- p95 < 300ms (no bcrypt; Postgres tx dominates at ~5–20ms).
- Expected QPS at MVP: 5–20.

## Sequence diagram

```mermaid
sequenceDiagram
    participant FE as Frontend
    participant ID as Identity
    participant DB as Postgres

    FE->>+ID: POST /api/v1/identity/auth/refresh {refreshToken}
    ID->>ID: ParseES256 — verify sig, exp, iss, aud, tokenType=REFRESH
    alt invalid JWT
        ID-->>FE: 401 AUTH_INVALID
    else valid JWT
        ID->>ID: extract old_jti, user_id; generate new_jti
        ID->>+DB: BEGIN Tx
        ID->>DB: UPDATE refresh_tokens SET revoked_at=NOW(), replaced_by_jti=$new_jti<br/>WHERE jti=$old_jti AND revoked_at IS NULL
        DB-->>ID: RowsAffected
        alt RowsAffected == 0 (already revoked or unknown jti)
            ID->>DB: ROLLBACK
            DB-->>-ID: OK
            ID-->>FE: 401 AUTH_REVOKED
        else RowsAffected == 1
            ID->>DB: INSERT refresh_tokens(jti=$new_jti, user_id, issued_at, expires_at)
            DB-->>ID: OK
            ID->>DB: COMMIT
            DB-->>-ID: OK
            ID->>ID: SignES256 → new accessToken + refreshToken
            ID->>ID: LOG INFO old_jti → new_jti
            ID-->>-FE: 200 {accessToken, refreshToken, ...expiresAt}
        end
    end
```

## Test cases

- `refresh_rotates_jti_atomically_and_returns_new_pair`
- `refresh_with_revoked_jti_returns_AUTH_REVOKED`
- `refresh_with_expired_jwt_returns_AUTH_INVALID`
- `refresh_with_invalid_signature_returns_AUTH_INVALID`
- `refresh_with_wrong_token_type_returns_AUTH_INVALID`
- `refresh_with_missing_refreshToken_returns_VALIDATION_ERROR`
- `refresh_rotation_is_atomic_under_concurrency`
- `refresh_audit_chain_replaced_by_jti_is_set`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Sonnet 4.6, dry-run #2) | Created |
