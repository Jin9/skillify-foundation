# Tech-Designer (per-component) Template

**Parent:** [`../templates.md`](../templates.md)
**Owner role:** [Tech-Designer](../roles.md) · **`template_version`:** 0.1.0

The detailed per-component spec Dev consumes. One TD per component (8 components in the dry-run #1 reference). Promoted from inline v0.1 used in dry-run #1.

**File location:** `<workflow_root>/design/components/<component_name>/td.json`

**Companion artifacts owned by the same role:**
- `<workflow_root>/design/components/<component_name>/erd.md` — ERD per microservice (see [`./erd.md`](./erd.md))
- `<workflow_root>/design/spec/<service>-<aggregate>/<API_NAME>.md` — per-API spec MDs (see [`./api-spec.md`](./api-spec.md))

---

## Fields

- `template_version` — string (semver), required
- `component_name` — string, required, mirrored from `tl-components.json`
- `required` — boolean, required, mirrored from TL
- `complexity` — enum, required, mirrored from TL
- `scaffold` — object, required: `{from_template, service_path, go_module, common_lib_import}`
- `persistence` — object, required: `{schema_owner, tables[], migrations[]}` — table list is high-level (full ERD lives in `erd.md`)
- `endpoints` — array, required (≥1): one per HTTP endpoint owned by this component (see endpoint object below)
- `event_consumers` — array, optional: each element `{event_contract_ref, dispatch_function, idempotency_strategy, state_driven: bool}`
- `event_producers` — array, optional: each element `{event_contract_ref, emit_path, partition_key}`
- `internal_modules` — object, required: `{domain_pkg, key_files[]}` — file paths Dev uses verbatim
- `test_strategy` — object, required: `{unit, integration, happy_paths[], negative_paths[]}` — see field detail below
- `ba_acceptance_mapping` — array, required (≥1): one entry per AC this component is responsible for

## Endpoint object

- `name` — string, required (`register`, `login`, etc.)
- `contract_ref` — string, required, must resolve in `contracts.json`
- `method_path` — string, required, format `POST /api/v1/<domain>/<aggregate>/<action>`
- `auth_required` — enum: `none` | `customer_jwt` | `admin_jwt` | `internal_secret`, required
- `request_schema` — object (JSON Schema), required
- `response_schema` — object (JSON Schema, envelope-wrapped), required
- `validation_rules` — array of strings, required: which field uses `common/validator` tag and what
- `error_contract` — array, required: `[{code, http_status, trigger}]`
- `edge_behavior` — object, required: `{idempotency, side_effects[], transactionality}`
- `implementation_hints` — array of strings: which `common/*` modules Dev uses, NOT actual code

## test_strategy detail

- `unit` — string description of unit-test approach (mockery + table-driven; mirror go-template)
- `integration` — string description (docker-compose + curl smoke)
- `happy_paths` — list of test names QA will look for
- `negative_paths` — list of test names QA will look for

---

## Schema (abbreviated; full schema generated alongside the run)

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "tech-designer/v0.1.0",
  "type": "object",
  "additionalProperties": false,
  "required": [
    "template_version", "component_name", "required", "complexity",
    "scaffold", "persistence", "endpoints",
    "internal_modules", "test_strategy", "ba_acceptance_mapping"
  ]
}
```

(Full schema enforces endpoint object shape, etc. Authored as a JSON Schema file alongside the template when the squad runs in a CI gate.)

---

## Example (known-good — abbreviated; full file in dry-run #1's `design/components/identity/td.json`)

```json
{
  "template_version": "0.1.0",
  "component_name": "identity",
  "required": true,
  "complexity": "standard",
  "scaffold": {
    "from_template": "B2C E-Commerce Platform/backend/go-template",
    "service_path": "B2C E-Commerce Platform/backend/services/identity",
    "go_module": "github.com/example/shoppilot/identity",
    "common_lib_import": "B2C E-Commerce Platform/backend/common"
  },
  "persistence": {
    "schema_owner": "identity",
    "tables": [
      {"name": "users", "primary_key": "id (UUID)", "indexes": ["UNIQUE(email)"]},
      {"name": "addresses", "primary_key": "id (UUID)", "indexes": ["INDEX(user_id) WHERE deleted_at IS NULL"]},
      {"name": "refresh_tokens", "primary_key": "jti (UUID)", "indexes": ["INDEX(user_id, expires_at)"]}
    ],
    "migrations": ["001_users.up.sql", "002_addresses.up.sql", "003_refresh_tokens.up.sql"]
  },
  "endpoints": [
    {
      "name": "register",
      "contract_ref": "identity.register",
      "method_path": "POST /api/v1/identity/auth/register",
      "auth_required": "none",
      "request_schema": {"type": "object", "required": ["email","password","name"], "properties": {"email": {"type":"string","format":"email"}}},
      "response_schema": {"type": "object", "required": ["code","message","data","traceId"]},
      "validation_rules": ["email format via common/validator", "password length 8-128", "name 1-80 chars"],
      "error_contract": [
        {"code": "VALIDATION_FAILED", "http_status": 400, "trigger": "any field violates schema"},
        {"code": "EMAIL_ALREADY_REGISTERED", "http_status": 409, "trigger": "email already exists"}
      ],
      "edge_behavior": {
        "idempotency": "none — duplicate email returns 409",
        "side_effects": ["INSERT users row with bcrypt hash, default role=CUSTOMER, status=ACTIVE"],
        "transactionality": "single-row insert"
      },
      "implementation_hints": ["use bcrypt cost 12 directly via golang.org/x/crypto/bcrypt", "use common/database pgx pool", "wrap response with common/wrapper.Success"]
    }
  ],
  "internal_modules": {
    "domain_pkg": "app/identity",
    "key_files": ["app/identity/handler_register.go", "app/identity/access/storage_user.go"]
  },
  "test_strategy": {
    "unit": "table-driven tests per handler with pgxmock or testcontainers Postgres",
    "integration": "docker-compose up; curl smoke through register → login → refresh",
    "happy_paths": ["register_with_unique_email_succeeds"],
    "negative_paths": ["register_with_duplicate_email_returns_409"]
  },
  "ba_acceptance_mapping": [
    {"story_id": "STORY_AUTH_REGISTER", "ac_index": 0, "endpoint": "register", "test_name": "register_with_unique_email_succeeds"}
  ]
}
```

---

## Negative examples

### Negative #1 — Endpoint with no error contract

```json
{
  "name": "login",
  "contract_ref": "identity.login",
  "method_path": "POST /api/v1/identity/auth/login",
  "auth_required": "none",
  "request_schema": {"type": "object"},
  "response_schema": {"type": "object"},
  "validation_rules": [],
  "error_contract": [],
  "edge_behavior": {"idempotency": "n/a", "side_effects": [], "transactionality": "n/a"},
  "implementation_hints": []
}
```

QA-L1 will reject (`spec_incomplete`, routes to TD, cap 1):

1. `error_contract: []` — every public endpoint has at least 1 error path. The login endpoint has at minimum AUTH_INVALID and VALIDATION_FAILED. Tag: `spec_incomplete`.
2. `validation_rules: []` — login takes credentials; field-level validation rules MUST be enumerated.
3. `request_schema` is an empty object — what fields, what types? Tag: `spec_incomplete`.

### Negative #2 — TD spec referencing a non-existent contract

```json
{
  "name": "create-from-checkout",
  "contract_ref": "order.create-from-checkout",
  "method_path": "POST /api/v1/order/internal/create-from-checkout",
  "auth_required": "internal_secret"
}
```

QA-L1 + Plan-Reviewer should catch:

1. `contract_ref: "order.create-from-checkout"` resolves only if Tech-Lead authored that contract. In dry-run #1, this gap was the cause of REV-L2-001 (buyerEmail missing). Tag: `spec_wrong` (routes to TD, cap 1) — TD must either (a) wait for TL to author the contract, or (b) restructure the orchestration to use a contract that exists.
2. Use the orphan-dep check from `tech-lead-components.md`'s validator to catch this at TL emission time, not at TD emission time.
