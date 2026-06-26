# Go conventions (high-impact only)

Curated for backend work. Conform to the repo first; where it is silent, use these.

## Control flow and readability
- **Guard clauses.** Handle errors and edge cases first with early returns; keep the happy path left-aligned and unindented.
- **Shallow nesting.** Max depth 3. Beyond that, return early or extract a function.
- **Small functions.** One job each. If you must scroll to read it, extract.
```go
// good: guard clauses, happy path last and flat
func (s *svc) charge(ctx context.Context, id uuid.UUID, amount Money) error {
	if amount.IsZero() {
		return ErrZeroAmount
	}
	acct, err := s.accounts.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("charge: load account: %w", err)
	}
	return s.accounts.Debit(ctx, acct, amount)
}
```

## Errors
- Handle every error; never discard with `_` unless you can justify it in one line.
- Wrap with context and `%w` so callers can `errors.Is`/`errors.As`: `fmt.Errorf("create product: %w", err)`.
- Define sentinel errors for conditions callers branch on: `var ErrNotFound = errors.New("not found")`.
- Do not `panic` for expected failures; reserve panic for programmer errors.
- Translate infrastructure errors into domain errors at the access boundary, not in business code.

## Context
- `context.Context` is the first parameter of any request-scoped function: `func Do(ctx context.Context, ...)`.
- Pass it down; never store it in a struct.
- Honor cancellation and deadlines on I/O and loops.

## Interfaces and types
- **Accept interfaces, return structs.** Functions take the narrow interface they use and return concrete types.
- **Define interfaces at the consumer** (or at the access layer that owns the data), kept small — often one or two methods.
- Do not create an interface until there are 2+ implementations or a real test fake (Simplicity gates).
- Make the zero value useful where you can; prefer it to an `Init()` step.

## Naming
- Short names in small scopes (`r`, `i`, `ctx`); descriptive names for package-level and exported identifiers.
- No stutter: in package `user`, name it `user.Store`, not `user.UserStore`.
- Exported identifiers get one short godoc line starting with the identifier name.

## Concurrency
- Default to synchronous code. Add a goroutine only when concurrency is genuinely needed.
- Own every goroutine's lifecycle: a clear stop signal (context), and a way to wait — no leaks.
- Guard shared state with a mutex, or pass ownership over a channel; do not mix both on the same data.
- Use `errgroup` for fan-out that must collect errors.

## Packages and dependencies
- Package name is short, lowercase, no underscores, and matches its directory.
- No `init()` for setup and no mutable package-level globals; pass dependencies explicitly.
- Reach for the standard library before adding a dependency; `a little copying is better than a little dependency`.

## Lean comments (production code)
- Comment only what the code cannot say: a magic number, a regex, a counter-intuitive "why" or business rule.
- Delete comments that restate code; rename instead of explaining.
- One godoc line per exported identifier; no decorative dividers, banners, or commented-out code.
- (Tests are exempt: follow the repo's test convention, e.g. AAA markers — see `testing-testify.md`.)

## Verify
Run whichever exist; report skipped ones:
- `gofmt -l .` (or `goimports`) — formatting.
- `go vet ./...` — suspicious constructs.
- `go build ./...` — catches hallucinated imports and type errors.
- `go test -race ./...` — behavior and data races.
