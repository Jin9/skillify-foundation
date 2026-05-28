# Story Template

**Parent:** [`../templates.md`](../templates.md)
**Owner role:** [BA](../roles.md) · **`template_version`:** 0.1.0

A Story is a single user-facing behavior with verifiable acceptance criteria. It belongs to exactly one EPIC. The BA emits one Story MD per behavior; downstream roles (Tech-Lead, Tech-Designer, Dev, QA) reference stories by `story_id`.

**File location:** `<workflow_root>/requirement/<EPIC_SLUG>/<STORY_SLUG>.md`
**ID convention:** `STORY_<DOMAIN>_<ACTION>` (uppercase, underscore-separated). Examples: `STORY_AUTH_LOGIN`, `STORY_AUTH_REFRESH`, `STORY_CATALOG_LIST`, `STORY_CHECKOUT_COMMIT`. Multi-word actions: `STORY_INVENTORY_RESERVATION_CREATE`. Domain matches the parent EPIC.

---

## Fields

- `story_id` — string, required, format `^STORY_[A-Z][A-Z0-9_]*$`
- `epic_id` — string, required, must match an existing `EPIC_<SLUG>`
- `title` — single-line, required
- `as_a` — string, required, the user role (`Customer`, `Admin`, `Guest`, `System`)
- `i_want` — string, required, the capability
- `so_that` — string, required, the benefit
- `acceptance_criteria` — list of `{given, when, then}` objects, required (≥1, each independently checkable)
- `priority` — enum (`Must` | `Should` | `Could` | `Wont`), required
- `size` — enum (`XS` | `S` | `M` | `L` | `XL`), required (relative effort)
- `edge_cases` — list of `{case, expected}` objects, required (empty array is a smell)
- `definition_of_ready` — list, required (defaults to global DoR; see below)
- `definition_of_done` — list, required (defaults to global DoD; see below)
- `requirement_refs` — list of strings, required (refs to requirement IDs like AUTH-001, CHK-005, etc.)
- `change_log` — table, required

**Default Definition of Ready (used unless story overrides):**
1. Story has a clear `as_a / i_want / so_that` triple.
2. Every AC has a falsifiable `then` clause (no "is good", "works correctly").
3. Edge cases enumerated (no empty `edge_cases`).
4. Requirement refs cite at least one ID from the source doc.
5. Story has been reviewed by Plan-Reviewer or has been stable for ≥1 cycle.

**Default Definition of Done (used unless story overrides):**
1. Code merged in the relevant service repo.
2. Unit tests for every AC (mapped via QA `ac_id`).
3. Integration test exercising the AC in `tests/` or via Playwright.
4. QA-L1 verdict `pass` for the owning component.
5. Reviewer-L1 has no unresolved `high` or `medium` findings on the story's code paths.
6. Documentation: API spec MD authored if the story adds a new endpoint; ERD updated if schema changed.

---

## File shape

```markdown
---
story_id: STORY_AUTH_LOGIN
epic_id: EPIC_AUTH
title: Customer logs in with email and password
as_a: Customer
i_want: to log in with my email and password
so_that: I can access my profile, cart, and order history
acceptance_criteria:
  - given: A customer with email "alice@example.com" and password "correct-pw" exists
    when: The customer POSTs /api/v1/identity/auth/login with those credentials
    then: |
      Response is 200 with code=SUCCESS and data contains accessToken (15min ttl) +
      refreshToken (14d ttl) + role=CUSTOMER. The accessToken is a valid ES256 JWT
      with sub=userId, role, exp claims.
  - given: A customer exists with email "alice@example.com" but the password is wrong
    when: The customer POSTs /api/v1/identity/auth/login with email and a wrong password
    then: |
      Response is 401 with code=AUTH_INVALID and message="Invalid email or password.".
      No tokens are returned. Latency MUST be within 50ms of the unknown-email path
      (constant-time-ish to prevent enumeration).
  - given: No user exists with email "ghost@example.com"
    when: The customer POSTs /api/v1/identity/auth/login with that email
    then: |
      Response is 401 with the EXACT same code=AUTH_INVALID and message=
      "Invalid email or password." as the wrong-password case. (AUTH-007.)
priority: Must
size: M
edge_cases:
  - case: Empty email or empty password
    expected: 400 VALIDATION_FAILED with field-level error
  - case: Email with leading/trailing whitespace ("  alice@example.com  ")
    expected: Email is normalized (trimmed + lowercased) before lookup; login succeeds if otherwise valid
  - case: 100 sequential failed login attempts in 60s from the same IP
    expected: Out of scope for MVP — flagged as v2 (rate-limiting) per ADR
requirement_refs: [AUTH-002, AUTH-005, AUTH-006, AUTH-007]
definition_of_ready: [default]
definition_of_done:
  - default
  - Constant-time-ish login latency verified by performance smoke test
change_log:
  - date: 2026-05-07
    author: BA (Claude Opus 4.7 xHigh)
    change: Created from requirement §8.1 AUTH-002 + AUTH-007
---

# STORY_AUTH_LOGIN — Customer logs in with email and password

## User narrative

As a **Customer**, I want to **log in with my email and password** so that **I can access my profile, cart, and order history**.

## Why this story exists in EPIC_AUTH

Login is the gateway to every authenticated path in the platform. Without a stable, secure login the rest of the customer journey (cart, checkout, orders) doesn't function.

## Acceptance criteria (Given/When/Then)

…rendered from frontmatter for readability…

## Edge cases

| Case | Expected |
|---|---|
| Empty email or empty password | 400 VALIDATION_FAILED with field-level error |
| Email with leading/trailing whitespace | Normalized (trim + lowercase); login succeeds if otherwise valid |
| 100 sequential failed login attempts in 60s | Out of scope for MVP — v2 rate-limiting (ADR-005) |

## Implementation notes (advisory, BA does NOT prescribe)

- Tech-Designer should specify ES256 key handling (env-PEM vs file, rotation strategy).
- Reviewer-L1 will check constant-time compare on the password and AUTH-007 message equivalence.

## Test mapping (QA-L1 fills this column)

| AC # | Test name |
|---|---|
| 1 | `login_with_correct_credentials_returns_token_pair` |
| 2 | `login_with_wrong_password_returns_AUTH_INVALID` |
| 3 | `login_with_unknown_email_returns_AUTH_INVALID_same_message` |

## Change log

| Date | Author | Change |
|---|---|---|
| 2026-05-07 | BA (Claude Opus 4.7 xHigh) | Created from requirement §8.1 AUTH-002 + AUTH-007 |
```

---

## Schema

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "story-output/v0.1.0",
  "type": "object",
  "additionalProperties": false,
  "required": [
    "template_version", "story_id", "epic_id", "title",
    "as_a", "i_want", "so_that",
    "acceptance_criteria", "priority", "size", "edge_cases",
    "definition_of_ready", "definition_of_done",
    "requirement_refs", "change_log"
  ],
  "properties": {
    "template_version": { "type": "string", "pattern": "^\\d+\\.\\d+\\.\\d+$" },
    "story_id":         { "type": "string", "pattern": "^STORY_[A-Z][A-Z0-9_]*$" },
    "epic_id":          { "type": "string", "pattern": "^EPIC_[A-Z][A-Z0-9_]*$" },
    "title":            { "type": "string", "minLength": 1, "maxLength": 120 },
    "as_a":             { "type": "string", "minLength": 1 },
    "i_want":           { "type": "string", "minLength": 1 },
    "so_that":          { "type": "string", "minLength": 1 },
    "acceptance_criteria": {
      "type": "array", "minItems": 1,
      "items": {
        "type": "object", "additionalProperties": false,
        "required": ["given", "when", "then"],
        "properties": {
          "given": { "type": "string", "minLength": 1 },
          "when":  { "type": "string", "minLength": 1 },
          "then":  { "type": "string", "minLength": 1 }
        }
      }
    },
    "priority": { "enum": ["Must", "Should", "Could", "Wont"] },
    "size":     { "enum": ["XS", "S", "M", "L", "XL"] },
    "edge_cases": {
      "type": "array",
      "items": {
        "type": "object", "additionalProperties": false,
        "required": ["case", "expected"],
        "properties": {
          "case":     { "type": "string", "minLength": 1 },
          "expected": { "type": "string", "minLength": 1 }
        }
      }
    },
    "definition_of_ready": { "type": "array", "items": { "type": "string", "minLength": 1 } },
    "definition_of_done":  { "type": "array", "items": { "type": "string", "minLength": 1 } },
    "requirement_refs":    { "type": "array", "minItems": 1, "items": { "type": "string", "minLength": 1 } },
    "change_log": {
      "type": "array", "minItems": 1,
      "items": {
        "type": "object", "additionalProperties": false,
        "required": ["date", "author", "change"],
        "properties": {
          "date":   { "type": "string", "pattern": "^\\d{4}-\\d{2}-\\d{2}$" },
          "author": { "type": "string" },
          "change": { "type": "string" }
        }
      }
    }
  }
}
```

`definition_of_ready` and `definition_of_done` accept the literal string `"default"` to inherit the global default in this template; or list overrides explicitly.

---

## Negative examples

### Negative #1 — Unfalsifiable `then`, empty edges, no requirement refs

```yaml
story_id: STORY_AUTH_LOGIN
epic_id: EPIC_AUTH
title: Customer logs in
as_a: Customer
i_want: to log in
so_that: I can use the site
acceptance_criteria:
  - given: Customer exists
    when: They log in
    then: It works
priority: Must
size: M
edge_cases: []
definition_of_ready: [default]
definition_of_done: [default]
requirement_refs: []
change_log:
  - {date: 2026-05-07, author: BA, change: Created}
```

What Plan-Reviewer should catch:

1. `then: It works` is unfalsifiable. Tag: `requirement_ambiguity` (high). Same failure mode as `ba.md`'s negative example #1.
2. `edge_cases: []` for an auth story — empty-input, malformed email, wrong-password, unknown-email all missing. Tag: `missing_edge_case` (high).
3. `requirement_refs: []` — the story doesn't trace to any requirement ID; Plan-Reviewer can't verify it's grounded. Tag: `requirements_gap` (high).
4. `as_a / i_want / so_that` is content-thin ("to use the site"). Tag: `requirement_ambiguity` (medium).

### Negative #2 — Story with implementation prescription leaked into AC

```yaml
story_id: STORY_AUTH_LOGIN
epic_id: EPIC_AUTH
title: Customer logs in with email and password
as_a: Customer
i_want: to log in with my email and password
so_that: I can access my account
acceptance_criteria:
  - given: A customer exists in the users table
    when: The frontend sends a POST to the Next.js route handler /api/proxy/identity/auth/login
          which forwards to identity service handler_login.go which calls bcrypt.CompareHashAndPassword
          on the password and pgx Query for the row
    then: A row is INSERTed into the refresh_tokens table with jti = uuid.NewV7() and the access_token
          is signed with the private key from the JWT_PRIVATE_KEY_PEM env var
priority: Must
size: M
edge_cases:
  - case: pgx connection is dropped mid-query
    expected: Retry with exponential backoff up to 3 times
requirement_refs: [AUTH-002]
definition_of_ready: [default]
definition_of_done: [default]
change_log:
  - {date: 2026-05-07, author: BA, change: Created}
```

What Plan-Reviewer should catch:

1. The AC `when` clause prescribes implementation details (Next.js route handler, handler_login.go, bcrypt, pgx). BA owns *behavior*, not implementation. The same behavior could be tested without naming the language or library. Tag: `unstated_assumption` (high) — BA crossed into Tech-Designer's territory.
2. The AC `then` clause names internal storage (refresh_tokens table, jti uuid v7) and env vars. These are implementation choices Tech-Designer should make. Tag: `unstated_assumption` (high).
3. `edge_cases` covers a low-level infrastructure case (pgx dropped); should cover business-visible edges (empty creds, wrong password, etc). Tag: `missing_edge_case` (medium).

Expected routing in both: Plan-Reviewer → BA, plan-stage cap 1, then HardFail.

---

## Story↔AC trace contract

The orchestrator validates before fan-out:

- Every AC in the workflow's `ba.json` matches at least one `STORY_*` (via story's `acceptance_criteria` content + `requirement_refs`).
- Every story has at least one AC.
- Every story's `epic_id` resolves to an existing EPIC MD on disk.
- Story IDs are unique within the workflow.

Validation failures route to BA with tag `requirements_gap` (high), plan-stage cap 1.
