# POST /api/v1/identity/auth/register

## Summary

Guest-facing self-registration. Creates a new CUSTOMER account with email, password, and name. Validates uniqueness of email, hashes the password with bcrypt cost 12 + pepper, and returns the created user's public profile. No session tokens are issued — the caller must subsequently call login. Implements AUTH-001 and AUTH-006.

## Story refs

- `STORY_AUTH_REGISTER`

## Contract ref

[`identity.register`](../../architecture/contracts.json) — `contract_version: 0.1.0`

## Auth

- Tier: **none** (public endpoint; no Authorization header required).

## Request

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "additionalProperties": false,
  "required": ["email", "password", "name"],
  "properties": {
    "email":    { "type": "string", "format": "email", "maxLength": 254 },
    "password": { "type": "string", "minLength": 8, "maxLength": 128 },
    "name":     { "type": "string", "minLength": 1, "maxLength": 120 },
    "phone":    { "type": "string", "pattern": "^[0-9+\\- ]{7,20}$" }
  }
}
```

## Response (success)

HTTP 200:

```json
{
  "code": "SUCCESS",
  "message": "Registration successful",
  "data": {
    "userId":    "01935b9c-1234-7000-abcd-000000000001",
    "email":     "alice@example.com",
    "name":      "Alice",
    "role":      "CUSTOMER",
    "createdAt": "2026-05-08T10:00:00Z"
  },
  "traceId": "0af7651916cd43dd8448eb211c80319c"
}
```

## Response (errors)

| Code | HTTP | Trigger |
|---|--:|---|
| `VALIDATION_ERROR` | 400 | Any required field missing, email format invalid, password < 8 chars, name empty |
| `DUPLICATE_EMAIL` | 409 | Email already registered (pgx error 23505 on UNIQUE(email)) |

## Business logic steps

1. Bind and validate request body (`common/validator`); return `VALIDATION_ERROR` on any schema violation.
2. Normalise email: `strings.TrimSpace(strings.ToLower(req.Email))`.
3. Hash password: `common/hash.Hash(password)` — bcrypt cost 12 with pepper from env `BCRYPT_PEPPER`.
4. Generate `userId = uuid.NewV7()`.
5. `INSERT INTO users(id, email, password_hash, name, phone, role='CUSTOMER', status='ACTIVE', created_at=NOW(), updated_at=NOW())`.
6. On pgx error code `23505` (unique violation on `email`): return `DUPLICATE_EMAIL` (HTTP 409).
7. Wrap and return success envelope with `common/wrapper.Success(ctx, data)` — **do not include `password_hash` in data**.

## Side effects

- INSERT one row into `identity.users`.
- No outbox event (MVP).
- No session cookies set (caller must call login separately).

## Idempotency

None — operation is non-idempotent. Repeat with same email returns 409 `DUPLICATE_EMAIL`. No Idempotency-Key header required or accepted.

## Performance

- p95 < 500ms (bcrypt cost 12 dominates at ~250–350ms wall-clock).
- Expected QPS at MVP: 1–5 (registration is a one-time event per user).

## Test cases

- `register_with_unique_email_succeeds`
- `register_with_duplicate_email_returns_409_DUPLICATE_EMAIL`
- `register_with_invalid_email_format_returns_VALIDATION_ERROR`
- `register_with_short_password_returns_VALIDATION_ERROR`
- `register_with_empty_name_returns_VALIDATION_ERROR`
- `register_with_email_whitespace_normalised_succeeds`
- `register_response_does_not_contain_password_hash`
- `register_persisted_password_is_bcrypt_not_plaintext`

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-08 | Tech-Designer (Sonnet 4.6, dry-run #2) | Created |
