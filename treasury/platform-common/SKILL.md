---
name: platform-common
description: >-
  Shared Go infrastructure library for A-Team Krungthai DGL microservices —
  Gin middleware, Kafka producer/consumer, JWT, structured slog logging,
  response envelopes, database/Redis/Firestore/GCS/S3 connectors, AES/RSA
  crypto, structured errors, mockable clock, config parser. Use when the user
  asks to "import common", "use the common library", "wrapper.Respond",
  "BindJSON", "kafka.NewProducer", "kafka.NewEventRouter", "logger.New",
  "JWT middleware", "AccessLog middleware", "serror.Wrap", "add a package to
  common", "extend wrapper", or how a common library package API works.
  Do NOT use for scaffolding a new microservice or adding a domain aggregate
  (use platform-go-service skill instead). Do NOT use for repo-policy,
  branching, or CI configuration (those live in AGENTS.md / CLAUDE.md).
argument-hint: Describe which common package you want to use, extend, or integrate
---

# Platform Common — Shared Go Library Skill

Use this skill to use, extend, or integrate the `common` shared library that every A-Team Go microservice depends on.

- **Module:** `gitdev.devops.krungthai.com/ktb-digital-lending/dgl-new/shareable/common`
- **Go:** 1.25+ | **Framework:** Gin | **Logger:** `log/slog` (JSON)
- **Boundary:** infrastructure only. No domain logic. Domain code lives in the consuming service.

## Workflow

Pick the branch that matches the user's request and follow its steps. If multiple branches apply, run them in the order listed.

### A. Use an existing package in a consuming service

1. Identify the capability the user wants (HTTP response, request binding, Kafka publish, JWT verify, logger, etc.).
2. Look up the package and key APIs in `references/packages.md`.
3. Open `references/recipes.md` and copy the matching recipe section verbatim, then adapt names and types.
4. Wire dependencies through interfaces — never pass concrete structs across package boundaries.
5. Wrap I/O errors with `serror.Wrap(err)` and attach context with `.With(slog.Attr…)`.
6. Verify the change builds with `go build ./...`.

### B. Add a new package to `common`

1. Confirm the package belongs here: it must be infrastructure, reusable across services, and contain zero domain logic. If not, redirect the user to the consuming service's `app/<aggregate>/`.
2. Read `references/conventions.md` end to end before writing code.
3. Create `<package>/<package>.go` with an exported interface + concrete implementation + `var _ Interface = (*impl)(nil)` compile-time check.
4. Create `<package>/mocks/mocks.go` via mockery (uses repo-root `.mockery.yaml`).
5. Write `<package>/<package>_test.go` — table-driven (`t.Run`) with 100% branch coverage; use `testify/assert` and `testify/require`.
6. Run `make precommit` (see `references/commands.md`).

### C. Extend an existing package

1. Re-read the public interface in the package's primary `.go` file.
2. Prefer additive changes (new methods on the interface, new options) over breaking changes.
3. If a breaking change is unavoidable, surface it explicitly to the user and pause for confirmation before editing.
4. Update mocks: re-run mockery after changing the interface.
5. Add or extend table-driven tests for every new branch.
6. Run `make precommit`.

### D. Explain a `common/<package>` API to the user

1. Find the package in `references/packages.md`.
2. Show the recipe for it from `references/recipes.md` (copy the relevant section).
3. If the recipe is missing, read the package source directly and explain in terms of the established patterns from `references/conventions.md`.

## Output Contract

- Code edits stay inside the consuming service or inside `common/<package>/` — never cross both in a single change without flagging it.
- Every new public function/method has a matching test case in the same package's `_test.go`.
- Every external dependency is exposed through an interface with a generated mock under `mocks/`.
- Commit only after `make precommit` passes.

## Hard Constraints

- DO NOT add `main.go` to `common`. It is a library.
- DO NOT add domain business logic to `common`.
- DO NOT use `log` or `fmt.Println` for logging — use `slog`.
- DO NOT introduce new `Must*` constructors; return `(T, error)` instead.
- DO NOT panic in library code outside startup paths.
- DO NOT bypass `wrapper.Respond` for HTTP responses (it powers `AccessLog`).
- DO NOT hardcode secrets or environment values; use `config.ParseEnv` with `env` struct tags.
- DO NOT skip `make precommit` before committing.

## References

| Need | File |
|---|---|
| Package directory and key APIs | `references/packages.md` |
| Working code snippets (response, bind, kafka, JWT, …) | `references/recipes.md` |
| Architectural rules, idioms, mock layout, new-package checklist | `references/conventions.md` |
| Make targets | `references/commands.md` |
