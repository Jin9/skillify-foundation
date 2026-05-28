# Conventions

Mandatory rules for changes to `common`. Read before adding a new package or modifying an existing one.

## Architectural Pillars

The library is designed around three complementary pillars.

### Domain-Driven Design (DDD)

- The `common` module is a *Shared Kernel* (Evans). It contains zero domain logic; every package is a reusable infrastructure concern.
- Domain-specific codes live in the consuming service's `app/codes.go`, never here. The `app/` package in common is a thin re-export layer that bridges `wrapper` response types to the consuming service's domain namespace.
- Aggregate-per-package convention: each consuming service organises domain logic under `app/<aggregate>/`, with `access/` sub-packages defining port interfaces for external dependencies (Ports-and-Adapters / Hexagonal).

### CQRS (Command Query Responsibility Segregation)

- Command path (HTTP `POST/PUT/DELETE` and Kafka consumer handlers) is separated from query path (`GET` handlers) at the handler level.
- The `kafka.Message[T]` envelope carries domain events through the async command channel. `NewEventRouter` routes events to per-aggregate handlers.
- `wrapper.Respond` provides a uniform response envelope (`code`, `message`, `data`) for both commands and queries.

### Fowler-Style Patterns

| Pattern | Implementation |
|---|---|
| Gateway (PoEAA) | `httpclient`, `database`, `redis`, `firestore`, `cloud-storage`, `sftp` |
| Clock Wrapper | `clock.Clock` interface — mockable system time |
| Special Case / Null Object | `safe.Deref`, `safe.DerefOr`, `safe.Ptr` |
| Decorator / Interceptor | `kafka.WithLogging(producer, …)` wraps Producer |
| Event Router | `kafka.NewEventRouter(handlers)` — content-based router |
| Registry (Replacer Functions) | `logger.New(replacers…)` — pluggable log key mapping |
| Separated Interface | Every external dep exports an interface + concrete + mock |
| Value Object | `wrapper.Code`, `wrapper.Message`, `token.Claims` |
| Money / Envelope | `kafka.Message[T]`, `wrapper.Response[T]` |

## Go Idioms

- Accept interfaces, return structs — every constructor returns a concrete type wrapped in a named interface (`Producer`, `Clock`, `Cipher`, etc.).
- Functional Options via `OptionFunc` (`httpclient`); variadic `ReplacerFunc` (`logger`).
- `Must*` / non-`Must*` pairs — `NewProducer` returns `(Producer, error)`, `MustNewProducer` panics. `Must*` is deprecated in library code; prefer returning errors.
- Compile-time interface checks — `var _ Interface = (*impl)(nil)`.
- Table-driven tests with `t.Run` and `testify/assert` + `testify/require`.
- `context.Context` as first arg in all I/O-facing functions.
- Errors are values — wrap with `serror.Wrap(err)` for source-location tracking and structured investigation context.
- Generics — `wrapper.Response[T]`, `wrapper.BindJSON[T]`, `kafka.Message[T]`, `httpclient.Get[RES]`, `config.ParseEnv[T]`, `safe.Deref[T]`.

## Interface + Mock Layout

Every external dependency follows this layout:

```
<package>/
├── <package>.go       # Exported interface + concrete implementation
├── <package>_test.go  # Tests (white-box, same package)
└── mocks/
    └── mocks.go       # mockery-generated mock
```

Repo root holds `.mockery.yaml`:

```yaml
structname: "{{.InterfaceName}}Mock"
pkgname: "{{ .SrcPackageName }}_mocks"
template: testify
```

Compile-time interface check in every package:

```go
var _ Producer = (*producer)(nil)
```

## Adding a New Package

1. Create `<package>/<package>.go` with an exported interface.
2. Add concrete implementation behind the interface.
3. Add `var _ Interface = (*impl)(nil)` compile-time check.
4. Create `<package>/mocks/` via mockery.
5. Write `<package>/<package>_test.go` — table-driven, 100% branch coverage.
6. Use `testify/assert` + `testify/require`. Gin tests use `httptest.NewRecorder()` + `gin.CreateTestContext(w)`.
7. Run `make precommit` before pushing.

## Library Rules

- No `main.go` — this is a library.
- No domain business logic — keep packages infrastructure-only.
- Use `slog` for logging (never `log` or `fmt.Println`).
- Return errors, don't panic (except `Must*` startup helpers, which are deprecated).
- Configuration via `env` struct tags, parsed by `config.ParseEnv`.

## Conventions Checklist

When adding or modifying packages in `common`:

- [ ] Package has a clear, single responsibility
- [ ] Public API defined via interface (where applicable)
- [ ] `var _ Interface = (*impl)(nil)` compile-time check
- [ ] Mocks auto-generated with mockery in `mocks/` sub-directory
- [ ] Tests achieve 100% branch coverage (table-driven, `t.Run`)
- [ ] No `main.go` — library only
- [ ] No domain-specific business logic
- [ ] HTTP responses via `wrapper.Respond`
- [ ] Logging via `slog` (never `log` or `fmt.Println`)
- [ ] Config uses `env` struct tags, parsed via `config.ParseEnv`
- [ ] Errors returned, not panicked (`Must*` deprecated for new code)
- [ ] `context.Context` as first parameter for I/O functions
- [ ] `make precommit` passes before committing
