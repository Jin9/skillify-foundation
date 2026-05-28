# API Spec (Request/Response)

This document lists all HTTP paths from router/router.go with request/response payloads derived from handler structs.

## Response envelope (all JSON responses)

Responses are wrapped by the common response envelope. Successful responses include:

```json
{
  "code": "SUCCESS",
  "message": "SUCCESS",
  "data": {},
  "traceId": "string (optional)"
}
```

Error responses use the same envelope with `code`/`message` indicating the error and `data` omitted or null. HTTP status codes are set per handler.

---

## Health

### GET /liveness

- Request: none
- Response: health payload (version/commit)

### GET /readiness

- Request: none
- Response: readiness payload

### GET /metrics

- Request: none
- Response: Prometheus metrics text

---

## Auth

### POST /api/v1/platform/auth/resolve-identity

**Authentication**: Public (no auth required)

Request:

```json
{
  "googleAccessToken": "string"
}
```

Response data:

```json
{
  "memberId": "uuid",
  "username": "string",
  "email": "string",
  "hashedEmail": "string",
  "status": "string",
  "organizationID": "uuid",
  "role": "string",
  "profileImage": "string"
}
```

### POST /api/v1/platform/auth/issue-token

**Authentication**: Public (no auth required)

Request:

```json
{
  "memberId": "uuid",
  "organizationId": "uuid",
  "memberRole": "string"
}
```

Response data:

```json
{
  "accessToken": "string"
}
```

---

## Organization

### POST /api/v1/platform/organization/current

**Authentication**: Public (no auth required)

Request:

```json
{
  "memberId": "uuid"
}
```

Response data:

```json
{
  "organizationId": "uuid",
  "name": "string",
  "status": "string",
  "role": "string",
  "joinedAt": "time",
  "createdAt": "time",
  "updatedAt": "time"
}
```

### POST /api/v1/platform/organization/register

**Authentication**: Public (no auth required)

Request:

```json
{
  "email": "string",
  "name": "string",
  "memberId": "uuid"
}
```

Response data:

```json
{
  "organizationId": "uuid",
  "memberId": "uuid",
  "role": "string"
}
```

---

## Member

### POST /api/v1/platform/member/me

**Authentication**: Public (no auth required)

Request:

```json
{
  "memberId": "uuid"
}
```

Response data:

```json
{
  "memberId": "uuid",
  "username": "string",
  "email": "string",
  "hashedEmail": "string",
  "status": "string",
  "createdAt": "time",
  "updatedAt": "time"
}
```

### POST /api/v1/platform/member/register

**Authentication**: Public (no auth required)

Request:

```json
{
  "email": "string",
  "name": "string"
}
```

Response data:

```json
{
  "memberId": "uuid"
}
```
