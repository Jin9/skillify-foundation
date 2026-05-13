# Go Code Style

## Contents
- Project structure
- Principles
- Error handling
- Refactoring (Fowler-style) workflow
- Code smells → Go-specific refactoring actions
- Hard rules

## Project structure

NEVER use: `common/`, `utils/`, `helpers/`.

Prefer a domain-oriented layout:

```
app/
  codes.go                        # Domain-wide response codes
  timestamps.go                   # Domain-wide timestamp helpers
  <domain>/                       # One sub-package per domain aggregate
    handler.go                    # HandlerConfig + handler struct + constructors
    handler_<action>.go           # HTTP handler (Swagger-annotated)
    handler_<action>_test.go      # HTTP handler test
    consumer_<action>.go          # Kafka consumer handler
    consumer_<action>_test.go     # Kafka consumer test
    service_<action>.go           # Domain service logic (optional)
    access/
      <dependency>.go             # Interface per external dependency
      types.go                    # Domain types (value objects)
```

- One package per aggregate. Action-suffixed files keep packages navigable at scale.
- `access/` isolates dependency interfaces from handler/service code — consumer-side interface definition.
- Use `types.go` for domain structs, NOT `model.go` (Go convention; "model" is ORM/MVC terminology).

## Principles

- Interface-driven design: define interfaces at the consumer side, not the producer.
- Explicit dependencies: inject everything. No package-level globals.
- Testable components: inject time, IO, and external clients.
- Prefer table-driven tests.
- `context.Context` as the first parameter, always.

## Error handling

- Errors are values. Wrap with context via `fmt.Errorf("doing X: %w", err)`.
- Use **sentinel errors** (`var ErrNotFound = errors.New(...)`) for expected outcomes that callers must handle.
- Use **typed errors** (a struct implementing `error`) when callers need to branch on domain classification (validation vs conflict vs authorization).
- NEVER swallow errors. Handle, wrap, or return.
- Log errors at the top of the call chain (handler/entrypoint), not deep in domain code.

## Refactoring (Fowler-style) workflow

Follow Martin Fowler's discipline: behavior-preserving, incremental, test-backed transformation. Never refactor and change behavior in the same commit.

1. **Green tests first** — if no tests cover the target code, write characterization tests before touching anything.
2. **One refactoring move per commit** — Extract Function, Move Function, Inline Variable, etc. Name the move in the commit message.
3. **Run tests after every move** — if red, revert immediately. Do not debug a broken refactor.
4. **No feature changes during refactoring** — separate "refactor" commits from "feat/fix" commits in history.

## Code smells → Go-specific refactoring actions

| Smell | Fowler refactoring | Go application |
|---|---|---|
| Long function (>40 lines) | Extract Function | Extract into a named method on the receiver or a standalone func; keep in the same package unless reuse is needed |
| Feature Envy (func uses another struct's fields more than its own) | Move Function | Move the method to the struct it actually operates on |
| Primitive Obsession (string IDs, raw int status) | Replace Primitive with Type | Introduce a named type (`type LoanID string`, `type Status int`) with methods |
| Data Clump (same 3+ params always travel together) | Introduce Parameter Object | Define a struct; pass it as a single param |
| Shotgun Surgery (one change → edits in 5+ files) | Move Function, Inline Class | Consolidate scattered logic into the owning domain package |
| Divergent Change (one file changes for unrelated reasons) | Extract Function, Split Package | Split into action-suffixed files (`service_create.go`, `service_cancel.go`) |
| Middle Man (func that only delegates) | Inline Function | Remove the passthrough; call the real impl directly |
| Dead Code | Remove Dead Code | Delete unused exports; `go vet` + linter catch unexported dead code |
| God Struct (struct with 10+ fields, 15+ methods) | Extract Struct, Move Method | Split into focused structs with a clear single responsibility |
| Large Package (>15 files, >2000 LOC) | Split Package | Break into sub-packages by aggregate or concern |

## Hard rules

- NEVER rename public API symbols in a refactoring commit — that is a breaking change, not a refactoring.
- NEVER refactor without test coverage. If tests do not exist, write them first (in a separate commit).
- NEVER mix "move code" with "change logic" — the diff should show structural change only, no new behavior.
- Prefer Extract over Rewrite — extracting preserves behavior by default; rewriting risks introducing bugs.
- When unsure a change is behavior-preserving, add an assertion/test that proves it before and after.
- Refactoring is driven by the next change that needs to be made. Do not refactor speculatively.
