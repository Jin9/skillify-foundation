# Dev Template

**Parent:** [`../templates.md`](../templates.md)
**Owner role:** [Dev](../roles.md) · **`template_version`:** 0.1.0

Dev does NOT emit a JSON spec; Dev emits **code on disk**. This template defines the structure-mirror checklist Dev must satisfy + the rejection_block Dev returns when the TD spec is ambiguous.

**Code lands at:** `<workflow_root>/backend/services/<svc>/` (or `<workflow_root>/frontend/` for frontend-web).

**Hard rule (from [`../roles.md`](../roles.md)):** if the TD spec is ambiguous or incomplete, Dev rejects the task with a structured `rejection_block` and routes back to Tech-Designer. **No improvising.** Hallucinated gap-filling is the documented failure mode for parallel agents.

---

## Backend service mandatory structure (mirror `backend/go-template`)

> **Pinned baseline:** `backend/go-template@08eccda` (branch `develop`, full SHA `08eccdaef909e9673b29d375f9dc0b09ac19910d`). Re-pinning on every iteration plan; bumping requires a deliberate `dev.md` edit + a v0.X+1 entry in the README change log. The structure-mirror checklist below is checked against this SHA.

Every backend service must produce these files exactly (file paths relative to `backend/services/<svc>/`):

```
backend/services/<svc>/
├── main.go                         # Gin engine boot + graceful shutdown
├── go.mod                          # module github.com/<org>/<svc>; replace common -> ../../common
├── go.sum                          # COMMITTED (not gitignored); make precommit verifies
├── Dockerfile                      # multi-stage; CGO_ENABLED=0; distroless or alpine runtime; healthcheck
├── docker-compose.yml              # service + Postgres + (Kafka if KAFKA_ENABLED) for local dev
├── Makefile                        # targets: mod, lint, test, vuln, precommit, ci, docker, run
├── .golangci.yaml                  # mirror backend/go-template/.golangci.yaml (17 linters)
├── .mockery.yaml                   # mirror backend/go-template/.mockery.yaml
├── README.md                       # what's implemented + what's stubbed (per-handler 501 list); env vars; run instructions
├── config/
│   └── config.go                   # env-backed; ParseEnv pattern from common/config
├── router/
│   ├── deps.go                     # composition root
│   └── router.go                   # gin Engine wiring; route groups by auth tier (public, customer JWT, admin JWT, internal-secret)
├── app/<domain>/
│   ├── domain.go                   # types
│   ├── service.go                  # interface + impl
│   ├── handler_<action>.go         # one file per endpoint
│   ├── handler_<action>_test.go    # SIBLING test file MANDATORY for FULL handlers
│   ├── consumer.go                 # if event-consuming
│   ├── sweeper.go                  # if has scheduled background work
│   ├── middleware_<...>.go         # custom middleware (e.g., internal_auth)
│   └── access/
│       ├── storage_<entity>.go     # repository pattern, pgx
│       ├── storage_<entity>_test.go
│       ├── cache_<entity>.go       # if caching
│       └── client_<external>.go    # HTTP / gRPC clients to other services
└── migrations/
    ├── 001_<name>.up.sql
    ├── 001_<name>.down.sql
    └── seed/                       # version-aligned seed scripts per env
        ├── local.sql
        └── sit.sql
```

### Hard rules

- **`go.sum` MUST be committed.** Generate with `go mod tidy && go mod verify`. Repos missing `go.sum` cannot build in production (caught by FEEDBACK.md ☐ "go.sum committed").
- **`_test.go` siblings MANDATORY** for every FULL handler and every storage method. Table-driven; mockery-generated mocks (run `make mod` to regenerate after interface changes).
- **`make lint` MUST pass clean** before commit. `.golangci.yaml` enforces the 17-linter policy from `backend/go-template`.
- **Stubbed handlers** still register their route and return `wrapper.Error(c, 501, "NOT_IMPLEMENTED_MVP", "Endpoint stubbed for MVP — see README")`.
- **Use `common/*` libraries; do NOT vendor or copy them.** Specifically: `wrapper`, `serror`, `logger` (slog only), `database` (pgx), `kafka`, `token`, `validator`, `hash`, `generator`, `health`, `middleware`, `config`.
- **No `fmt.Println`.** Logging via `slog` only.
- **No `Must*` constructors.** Return errors.
- **Wrap errors with `serror.Wrap(err)`** for source-location tracking.
- **Migrations use schema-per-service:** `CREATE SCHEMA IF NOT EXISTS <svc>;` first migration; `SET search_path TO <svc>;` in app config.

### Structure-mirror checklist (Dev self-checks before reporting back)

- [ ] `go.mod` declares the right module name; `replace` points to `../../common`.
- [ ] `go.sum` exists and is committed.
- [ ] `go build ./...` passes.
- [ ] `go vet ./...` passes.
- [ ] `make lint` passes (golangci-lint clean against the 17 rules).
- [ ] Every FULL handler has a sibling `_test.go` with at least 1 happy-path and 1 negative-path test.
- [ ] Every STUB handler is registered and returns `501 NOT_IMPLEMENTED_MVP`.
- [ ] Migrations are sequentially numbered with both `.up.sql` and `.down.sql`; `down` actually reverses `up`.
- [ ] Dockerfile is multi-stage with `CGO_ENABLED=0`; healthcheck via `/liveness`.
- [ ] README enumerates: what's FULL, what's STUB, env vars, run instructions, dependency on other services.

---

## Frontend (Next.js) mandatory structure

```
frontend/
├── package.json                    # next 14+, react 18, typescript 5+, tailwindcss 3+
├── package-lock.json (or pnpm-lock.yaml)  # COMMITTED for reproducible builds
├── next.config.js                  # security headers (CSP baseline, X-Frame-Options, X-Content-Type-Options, Referrer-Policy)
├── tsconfig.json                   # strict mode
├── tailwind.config.ts
├── postcss.config.js
├── .eslintrc.json                  # next/core-web-vitals
├── .env.example                    # backend service URLs; never commit .env
├── README.md
├── app/
│   ├── layout.tsx
│   ├── globals.css
│   ├── page.tsx                    # Home (Server Component default)
│   ├── <route>/page.tsx            # one folder per page
│   └── api/
│       ├── auth/session/route.ts   # GET — session probe (does NOT trigger refresh)
│       └── proxy/<service>/<path>/route.ts  # backend proxies; use lib/auth.withAuth
├── lib/
│   ├── auth.ts                     # 6-step withAuth flow per cross-cutting.auth contract
│   ├── http.ts                     # forwarder helpers
│   └── idempotency.ts              # crypto.randomUUID() + in-flight Map
└── components/
    └── *.tsx
```

### Hard rules (frontend)

- **HttpOnly cookies only** for auth tokens; never `document.cookie` or `localStorage`/`sessionStorage`.
- **`crypto.randomUUID()`** for Idempotency-Key generation, never `Math.random()` or timestamp+counter.
- **Server Components by default**; `'use client'` only when interactivity demands it.
- **No `dangerouslySetInnerHTML`.**
- **CSP baseline in `next.config.js`** — no `unsafe-eval`; `script-src 'unsafe-inline'` only with a documented escape plan.
- **Backend URLs are server-only** (`BACKEND_*_BASE_URL` in `.env`, used by route handlers; never exposed to the browser).

---

## Rejection block (when Dev rejects ambiguous spec)

When TD spec is ambiguous or incomplete, Dev returns this structure (JSON) and does NOT write code:

```json
{
  "rejected": true,
  "component_name": "<svc>",
  "rejection_blocks": [
    {
      "ambiguity_locus": "td.json#/endpoints[name=login]/edge_behavior/transactionality",
      "missing_field": "transactionality",
      "spec_says": "n/a",
      "dev_question": "Login does NOT mutate state, so 'n/a' is correct — but should the storage_user.GetByEmail call be in a tx for read-isolation under concurrent password rotations?",
      "suggested_clarification": "Add 'read-only; no tx required; pgx pool acquire is sufficient' OR 'wrap GetByEmail + bcrypt.Compare in a SERIALIZABLE tx if password-rotation is in scope.'"
    }
  ]
}
```

The orchestrator routes a non-empty `rejection_blocks` array back to Tech-Designer with tag `spec_incomplete` (cap 1).

---

## Negative examples

### Negative #1 — Dev gap-filled instead of rejecting

TD spec says: "user can log in." Dev produces a login handler that also issues a special "remember-me" 30-day cookie because Dev assumed it.

QA-L1 + Reviewer-L1 catch:

1. The "remember-me" cookie is not in TD spec. Tag: `code_mismatch` (routes Dev, cap 2). Worse: Reviewer-L1 may also tag `code_security_issue` because a 30-day cookie wasn't reviewed for token-rotation implications.
2. Dev should have rejected with a `rejection_block` asking TD to clarify session lifetime instead of inventing one.

### Negative #2 — Missing `_test.go` siblings on FULL handlers

Dev ships `handler_register.go`, `handler_login.go`, `handler_refresh.go` (all FULL), no test files.

QA-L1 catches:

1. Hard rule violation. Tag: `code_mismatch` (routes Dev, cap 2). README enumerates the handlers as "FULL" but they have no automated coverage.
2. Reviewer-L1 may additionally fire `code_quality_issue` (medium) for absent test scaffolding when the structure-mirror checklist explicitly demands it.
