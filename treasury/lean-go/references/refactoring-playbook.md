# Refactoring playbook (Fowler, applied to Go)

Behavior-preserving changes only. One smell, one action, verify, repeat.

## The OODA refactor loop
1. **Observe** — find one code smell (table below).
2. **Orient** — pick the one Fowler action that removes it.
3. **Act** — apply only that action.
4. **Decide** — run `go build` and `go test` (`-race` if available). Green: keep going. Red: revert that step.
Stop after 5 iterations or when no smell remains. Report any smell left untouched and why.

Never mix a refactor with a behavior change in the same step. Do not change transaction or consistency boundaries while refactoring.

## Smell to action map
| Smell | Action |
|-------|--------|
| Long Function | Extract Function |
| Deep Nesting | Replace Nested Conditional with Guard Clauses |
| Long Parameter List (~4+) | Introduce Parameter Object / Preserve Whole Object |
| Duplicated Code (2nd time) | Extract Function and reuse (Rule of Three) |
| Feature Envy (method uses another type's data) | Move Function |
| Primitive Obsession | Encapsulate into a small type with methods |
| Unclear name | Rename |
| Speculative Generality / pass-through layer | Inline Function/Variable; delete the layer |
| Large `switch` on a type that drives behavior | Replace Conditional with Polymorphism (interface) |

## The shortlist (these earn their keep)
- **Extract Function** — name a block; the name documents intent and shrinks the caller.
- **Inline Function/Variable** — remove indirection that adds no value; the cure for premature abstraction.
- **Rename** — make the name state the intent; cheapest high-impact change.
- **Introduce Parameter Object / Preserve Whole Object** — collapse a noisy argument list into one struct (only at ~4+ related args).
- **Replace Nested Conditional with Guard Clauses** — early returns flatten the happy path.
- **Move Function/Field** — put behavior next to the data it uses.
- **Encapsulate Variable** — wrap a bare primitive/field in a type or accessor when invariants matter.
- **Replace Conditional with Polymorphism** — *rarely*; only when behavior genuinely varies by case and a Go interface removes a growing `switch`.

## Drop as noise
Do not reach for: an exhaustive Fowler catalog, inheritance-shaped patterns (no inheritance in Go), Extract Class for its own sake, clean-architecture mandates, default layering, enterprise DI containers, premature generics, or extreme DRY. A little duplication beats the wrong abstraction.

## Before / after
```go
// before: deep nesting, long function, mixed levels (router already enforces POST)
func handle(w http.ResponseWriter, r *http.Request) {
	var in CreateReq
	if err := json.NewDecoder(r.Body).Decode(&in); err == nil {
		if in.Name != "" {
			id, err := store.Create(r.Context(), in)
			if err == nil {
				_ = json.NewEncoder(w).Encode(id)
			} else {
				http.Error(w, "internal", 500)
			}
		} else {
			http.Error(w, "name required", 400)
		}
	} else {
		http.Error(w, "bad body", 400)
	}
}
```
```go
// after: guard clauses (flat happy path) + Extract Function (validation named)
func handle(w http.ResponseWriter, r *http.Request) {
	in, err := decodeCreate(r)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest) // log err internally; don't leak it
		return
	}
	id, err := store.Create(r.Context(), in)
	if err != nil {
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(id) // header already written
}

func decodeCreate(r *http.Request) (CreateReq, error) {
	var in CreateReq
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		return CreateReq{}, fmt.Errorf("bad body: %w", err)
	}
	if in.Name == "" {
		return CreateReq{}, errors.New("name required")
	}
	return in, nil
}
```
