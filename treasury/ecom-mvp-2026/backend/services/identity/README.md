# Identity Service

Handles user registration, authentication, and address management for the ShopPilot B2C platform.

## What is Implemented

| Endpoint | Handler | Status | Test file |
|---|---|---|---|
| `POST /api/v1/identity/auth/register` | `handler_register.go` | FULL | — |
| `POST /api/v1/identity/auth/login` | `handler_login.go` | FULL | `handler_login_test.go` (7 cases) |
| `POST /api/v1/identity/auth/refresh` | `handler_refresh.go` | FULL | `handler_refresh_test.go` (6 cases) |
| `POST /api/v1/identity/auth/logout` | `handler_logout.go` | STUB — 501 | — |
| `GET  /api/v1/identity/profile/me` | `handler_profile_read.go` | FULL | `handler_profile_read_test.go` (5 cases) |
| `PATCH /api/v1/identity/profile/me` | `handler_profile_update.go` | STUB — 501 | — |
| `POST /api/v1/identity/address/create` | `handler_address_create.go` | STUB — 501 | — |
| `GET  /api/v1/identity/address/list` | `handler_address_list.go` | STUB — 501 | — |
| `PATCH /api/v1/identity/address/update` | `handler_address_update.go` | STUB — 501 | — |
| `DELETE /api/v1/identity/address/delete` | `handler_address_delete.go` | STUB — 501 | — |
| `POST /api/v1/identity/address/set-default` | `handler_address_set_default.go` | FULL | `handler_address_set_default_test.go` — pending |

All 11 routes are registered. Remaining stubs return:
```json
{"code": "NOT_IMPLEMENTED_MVP", "message": "Endpoint stubbed for MVP — see README"}
```

## Un-stubbed in Phase 2 (dry-run #2)

- **`address.set-default`** — fully implemented. Atomic two-step Postgres transaction in `access/storage_address.go#SetDefault`: (1) clear `is_default` on all other live addresses for the user, (2) set `is_default = true` on the target. Returns `NOT_FOUND` (404) when address is missing, soft-deleted, or belongs to another user (`RowsAffected == 0` → ROLLBACK).
- **`handler_login_test.go`** — added (7 table-driven cases covering 200, 400 ×3, 401, 403, 500).
- **`handler_refresh_test.go`** — added (6 table-driven cases covering 200, 400 ×2, 401 ×2, 500).

## NOT_IMPLEMENTED_MVP List (v2 backlog)

- `logout` — revoke refresh token by jti scoped to authenticated user
- `profile.update` — PATCH name/phone with dynamic SET clause
- `address.create` — INSERT with optional default-clearing transaction
- `address.list` — SELECT WHERE deleted_at IS NULL
- `address.update` — dynamic PATCH with RowsAffected check
- `address.delete` — soft delete; idempotent (ID-DEC-3)

## Design Decisions

- **ID-DEC-1** Refresh token rotation in a single Postgres transaction (UPDATE old → INSERT new → COMMIT). Replay of a revoked JTI returns `AUTH_REVOKED` immediately.
- **ID-DEC-2** No 5-second soft idempotency window in MVP. Clients get `AUTH_REVOKED` on replay; v2 enhancement.
- **ID-DEC-4** Password hashing via `golang.org/x/crypto/bcrypt` at cost 12 (`bcrypt.DefaultCost` is 10; 12 is chosen for stronger work factor). Replaces the previous SHA-256+pepper approach (CWE-916/759, REV-IDENT-001) which lacked an adaptive work factor. bcrypt outputs ~60 chars; the `password_hash` column is TEXT, which is adequate.
- **AUTH-007** Login returns identical error code+message for unknown email and wrong password to prevent user enumeration.

## Run Locally

### 1. Generate ES256 keypair (dev only)

```bash
openssl ecparam -genkey -name prime256v1 -noout -out priv.pem
openssl ec -in priv.pem -pubout -out pub.pem
```

Export for docker-compose:

```bash
export JWT_PRIVATE_KEY_PEM="$(cat priv.pem)"
export JWT_PUBLIC_KEY_PEM="$(cat pub.pem)"
export HASH_PEPPER="my16charpepper!!"  # must be exactly 16 chars
```

### 2. Apply migrations

Migrations are plain SQL under `migrations/`. Apply in order before starting the service:

```bash
psql "postgresql://identity:identity_secret@localhost:5432/identity" \
  -f migrations/001_users.up.sql \
  -f migrations/002_addresses.up.sql \
  -f migrations/003_refresh_tokens.up.sql
```

### 3. Start with docker-compose

```bash
docker-compose up --build
```

### 4. Smoke test

**Register:**
```bash
curl -s -X POST http://localhost:8080/api/v1/identity/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123","name":"Alice"}' | jq .
```

**Login:**
```bash
curl -s -X POST http://localhost:8080/api/v1/identity/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}' | jq .
```

**Refresh (use refreshToken from login response):**
```bash
curl -s -X POST http://localhost:8080/api/v1/identity/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refreshToken":"<token>"}' | jq .
```

**Liveness:**
```bash
curl http://localhost:8080/liveness
```

## Environment Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `PORT` | Yes | — | HTTP listen port |
| `DB_HOST` | Yes | `localhost` | Postgres host |
| `DB_PORT` | No | `5432` | Postgres port |
| `DB_USER` | Yes | — | Postgres user |
| `DB_PASSWORD` | No | — | Postgres password |
| `DB_NAME` | Yes | — | Postgres database name |
| `JWT_PRIVATE_KEY_PEM` | Yes | — | ES256 private key PEM |
| `JWT_PUBLIC_KEY_PEM` | Yes | — | ES256 public key PEM |
| `JWT_ISSUER` | Yes | — | JWT iss claim |
| `JWT_AUDIENCE` | Yes | — | JWT aud claim |
| `JWT_ACCESS_TOKEN_TTL` | No | `15m` | Access token lifetime |
| `JWT_REFRESH_TOKEN_TTL` | No | `720h` | Refresh token lifetime (30 days) |
| `HASH_PEPPER` | Yes | — | 16-char password pepper |

## Module

```
gitlab.com/b2c-e-commerce-platform/platform/backend/services/identity
```

Common library is imported via local `replace` directive pointing to `../../common`.
