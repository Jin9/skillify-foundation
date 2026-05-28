---
name: platform-go-service-common
description: >-
  Scaffold and extend A-Team Krungthai DGL Go microservices that consume the
  shared common library. Covers DDD aggregates with Ports-and-Adapters,
  CQRS HTTP+Kafka handlers, dependency wiring, table-driven tests with 100%
  coverage, Dockerfile, and GitLab CI. Use when the user asks to "create a
  new service", "add a new aggregate", "scaffold an aggregate package",
  "wire deps.go", "add a Kafka consumer handler", "add an HTTP handler",
  "set up router.go", "graceful shutdown main.go", "100% test coverage", or
  "Dockerfile for service". Do NOT use for adding a new package to the
  shared common library or explaining a common library package API (use
  platform-common skill instead). Do NOT use for repo-policy or branching
  rules (those live in AGENTS.md / CLAUDE.md).
argument-hint: Describe the new service, aggregate, or feature to scaffold
---

# Platform Go Service — Scaffold & Convention Guide

Use this skill to create new Go microservices, or extend existing ones, following A-Team platform conventions. All services import `gitdev.devops.krungthai.com/ktb-digital-lending/dgl-new/shareable/common`.

- **Boundary:** service-side only. Reusable infrastructure goes in `common` (use `platform-common` for that).

## Workflow

Pick the branch that matches the user's request and follow its steps.

### A. Create a new microservice

1. Confirm scope: one Bounded Context per repo. If the user is adding behaviour to an existing service, use branch B or C instead.
2. Read `references/architecture.md` to internalise the directory layout, DDD vocabulary, and CQRS flow.
3. Initialise the Go module and pull `common` (see `references/scaffold.md` step 1).
4. Create `main.go`, `config/`, `router/deps.go`, `router/router.go`, `router/subscriber.go` from `references/scaffold.md` steps 2–6. Adapt names; keep middleware ordering verbatim.
5. Add at least one aggregate via branch B before committing.
6. Add `Dockerfile` and `gitlabci.yml` from `references/deployment.md`.
7. Run `make precommit`.

### B. Add a new domain aggregate (`app/<aggregate>/`)

1. Confirm the aggregate is a true DDD aggregate root, not a helper or DTO bag. If unclear, ask the user before creating the package.
2. Create `app/<aggregate>/access/` with one Port interface per external dependency (Storage, Producer, …) and a `model.go` for persistence/domain models. See `references/aggregate.md` §4.
3. Create `app/<aggregate>/handler.go` with `HandlerConfig` + `NewHandler` (see `references/aggregate.md` §1).
4. Add command, query, and consumer handlers as needed (`references/aggregate.md` §§2–3).
5. Register routes in `router/router.go`; register Kafka events in `router/subscriber.go`.
6. Add domain response codes in `app/codes.go` with a per-aggregate prefix (`references/aggregate.md` §5).
7. Generate mocks for every Port via mockery into `access/mocks/`.
8. Add tests per branch D until coverage hits 100% for the new package.

### C. Add a single handler to an existing aggregate

1. Re-read the aggregate's existing handlers and `access/` interfaces; reuse Ports if possible.
2. For HTTP, follow `references/aggregate.md` §2: `wrapper.BindJSON` → port call → `wrapper.Respond`. Wrap errors with `serror.Wrap(err)`.
3. For Kafka consumers, follow `references/aggregate.md` §3: `kafka.BindMessage` → port call → return wrapped error. Register in `router/subscriber.go`.
4. Add a domain code/message pair in `app/codes.go` if the response is new.
5. Add table-driven tests covering every error branch (`references/tests.md`).

### D. Add or extend tests

1. Follow the table-driven shape in `references/tests.md` §§1–2.
2. Use mockery `.EXPECT()` mocks; set expectations only for the call shapes the branch exercises.
3. Cover every `if err != nil` branch with a dedicated case.
4. Run `go test ./... -race -cover` and confirm 100% statement coverage on changed `app/` packages.

## Output Contract

- New service skeleton matches the layout in `references/architecture.md` exactly.
- Every aggregate has `handler.go`, `*_handler.go`, `*_consumer.go` (where applicable), `access/<port>.go`, `access/model.go`, mocks under `access/mocks/`, and 100%-covered `_test.go` files.
- Routes live only in `router/router.go`; Kafka event→handler wiring lives only in `router/subscriber.go`.
- All HTTP I/O passes through `wrapper.BindJSON` and `wrapper.Respond`.
- All errors crossing layer boundaries are wrapped with `serror.Wrap(err)`.
- Commit only after `make precommit` passes.

## Hard Constraints

- DO NOT put domain logic in `common`. Domain logic lives in `app/<aggregate>/`.
- DO NOT skip Ports. Every external dependency goes through an interface in `access/`.
- DO NOT bypass `wrapper.Respond` or `wrapper.BindJSON` for HTTP I/O.
- DO NOT skip mock generation for new Port interfaces.
- DO NOT panic in handlers. Return wrapped errors via `serror.Wrap`.
- DO NOT register routes outside `router/router.go` or Kafka handlers outside `router/subscriber.go`.
- DO NOT commit if any `app/` package coverage is below 100%.

## References

| Need | File |
|---|---|
| Directory layout, DDD vocabulary, CQRS flow | `references/architecture.md` |
| `main.go`, `config/`, `router/deps.go`, `router/router.go`, `router/subscriber.go` | `references/scaffold.md` |
| Handler / consumer / access / codes recipes | `references/aggregate.md` |
| Table-driven HTTP and Kafka tests, coverage rules | `references/tests.md` |
| Dockerfile, GitLab CI, runtime deps | `references/deployment.md` |
