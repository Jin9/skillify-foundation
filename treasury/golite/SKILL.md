---
name: golite
description: >
  Implement, improve, and refactor Go backend code as a senior engineer: simple,
  high-impact, readable code with lean one-line comments, idiomatic Go conventions,
  Martin Fowler refactoring idioms, a handler/service/access layering kernel, and
  table-driven stretchr/testify unit tests with mockery. Adapts to any existing
  project structure instead of imposing one, and recommends best practices as opt-in.
  Use when the user asks to implement a Go handler, service, consumer, or repository;
  write or refactor Go backend code; make Go code idiomatic, simpler, or more readable;
  clean up a Go function; write testify table-driven unit tests; or apply Go best
  practices. Do NOT use for non-Go code, frontend, or infrastructure; for migrating a
  messy microservice toward a DDD or CQRS target architecture use refactoring-go-services;
  for go-template scaffold-locked requirement work use implementing-go-template-requirements;
  for A-Team platform service scaffolding use platform-go-service.
---

# Golite

## Purpose
Act as a senior Go backend engineer: implement, improve, and refactor backend code that is simple, high-impact, and readable, conforming to the repository you are already in.

## The cognitive kernel
Boot this small, fixed core on every task; load everything else on demand.

1. **Three pillars** — simple, high-impact, readable. Clear beats clever. Ship the smallest change that solves the problem.
2. **Layering invariant** — separate transport, business, and integration (see Layering kernel).
3. **Simplicity gates** — never add abstraction the code does not yet need (see Simplicity gates).

Everything deeper — Go idioms, the Fowler catalog, the testify pattern — is *userland*: read the matching reference only when a step needs it. The kernel stays small and fixed; modules load when needed.

## Layering kernel
Three roles, regardless of the repo's folder names:

- **transport** — HTTP handler or queue consumer. Parse and validate input, call business logic, map the result or error to the response shape. Thin.
- **business** — request rules, orchestration, domain decisions. Lives in the handler method while small; extract to a `service` unit only when orchestration grows (Fowler Extract Function).
- **integration (access)** — `storage` (Repository), `client` (Gateway), `cache` (Cache). An interface plus a lean implementation that wraps infrastructure, maps infra errors to domain errors, and holds **no business logic**.

Dependency rule: business depends on access **interfaces** (defined at the access layer, injected through a constructor config, mocked in tests). Business rules never leak down into access; infrastructure types never leak up into transport.

Map these roles onto whatever structure exists — flat, `cmd/internal/pkg`, DDD, hexagonal, or custom names like `repository`/`usecase`/`store`. Conform to the repo's naming; do not impose `handler/service/access`. See `references/architecture-kernel.md`.

## When to use this skill
- Use when: implementing a Go handler, consumer, service, or repository.
- Use when: writing, simplifying, or refactoring Go backend code, or making it idiomatic and readable.
- Use when: writing table-driven `testify` unit tests for Go.
- Do NOT use: non-Go code, frontend, or infrastructure; DDD/CQRS-target migration (use `refactoring-go-services`); go-template scaffold work (use `implementing-go-template-requirements`); A-Team platform scaffolding (use `platform-go-service`).

## Workflow
1. **Detect — conform, don't impose.** Read `go.mod` (module path, Go version). Identify where the repo places transport / business / integration and what it names them. Read 1–2 nearest sibling files and one sibling `_test.go` for naming, error, logging, and test+mock style (testify? `.mockery.yaml`? gomock? hand fakes?). *Exit when* you can name the target package and the conventions to match. → `references/architecture-kernel.md`.
2. **Scope.** State the one responsibility in a sentence. Choose the smallest set of files to touch. If the task needs a structural rewrite, ask once before proceeding.
3. **Implement.** Place code by role (Layering kernel). Keep access lean behind interfaces; keep business out of access. Apply the Simplicity gates and Lean comments. Pull idioms from `references/go-conventions.md` only as needed.
4. **Refactor (cleanup tasks) — OODA loop.** Observe one smell, apply one Fowler action, run `go build` and `go test`, repeat. Preserve behavior. Stop after 5 iterations or when no smell remains, and report what is left. → `references/refactoring-playbook.md`.
5. **Test.** Write or extend table-driven `testify` tests in the repo's pattern: one case per behavior branch (success and each failure mode), `require.New(t)`, a `prepare` closure for mock setup, deps injected via the constructor config. When an interface changes, regenerate mocks with mockery — never hand-edit generated `mocks.go`. → `references/testing-testify.md`.
6. **Recommend (opt-in).** List best-practice improvements you did not make as a separate, numbered list. Do not apply them without approval.
7. **Verify and report.** Run whichever of `gofmt -l`, `go vet ./...`, `go build ./...`, `go test -race ./...` exist; report skipped tools plainly. `go build` catches hallucinated imports; `-race` catches concurrency bugs. Confirm every new branch has a test. Emit the Output contract.

## Simplicity gates
Refuse complexity the code has not earned:

- No abstraction before the second real call site (Rule of Three: extract on the second duplicate, not the first).
- No new interface until 2+ implementations exist — a genuine test fake counts as the second.
- Delete pass-through layers that only forward calls without adding logic.
- No new package for fewer than ~3 cohesive exported things.
- No new dependency when the standard library or an existing dependency suffices.
- Default to synchronous code; add goroutines only when concurrency is genuinely required, and own their lifecycle.
- No generics unless they remove real duplication without hurting readability.
- Maximum nesting depth 3; beyond that, return early or extract a function.
- No `init()` and no mutable package-level state.

## Lean comments
Comment only the non-obvious; one godoc line per exported identifier; never narrate. Tests keep the repo's convention (e.g. AAA markers). Full rules: `references/go-conventions.md`.

## Output contract
Produce edited or created `.go` files (including `_test.go`) that match the existing structure, plus a short summary:

- what changed and where, by role;
- Fowler idioms applied;
- tests added or updated and the branches they cover;
- opt-in recommendations not applied;
- verification results (`gofmt` / `vet` / `build` / `go test -race`, or which were skipped).

Production comments: one line, only where non-obvious.

## Constraints
- Conform to the repo before changing it; match its naming, error, and logging style.
- Keep business logic out of the access layer; business lives in transport/business and depends on access interfaces.
- Map the three roles onto the repo's own names; do not impose `handler/service/access` folders.
- Never swallow errors — wrap with `%w` and context.
- Treat exported names, JSON tags, and error sentinels as contracts; do not break them silently.
- Refactors must preserve behavior, including transaction and consistency boundaries.
- Give every new behavior branch a `testify` case; regenerate mocks with mockery, never hand-edit `mocks.go`.
- Keep recommendations separate from applied edits.
- Do not add a dependency, abstraction, package, or interface that the Simplicity gates forbid.

## References
| Need | Reference |
|------|-----------|
| handler/service/access roles, dependency rule, mapping onto any layout | `references/architecture-kernel.md` |
| idiomatic Go conventions, lean comments, readability | `references/go-conventions.md` |
| Fowler refactoring shortlist, smell to action map, OODA loop | `references/refactoring-playbook.md` |
| table-driven testify + mockery unit-test pattern | `references/testing-testify.md` |
