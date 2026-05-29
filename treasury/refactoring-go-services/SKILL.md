---
name: refactoring-go-services
description: >
  Incrementally refactor messy Go microservices toward clean DDD/CQRS architecture while preserving behavior. Use when the user asks "clean up this Go service", "refactor this handler", "move logic into the domain layer", or "remove Go service code smells". Do NOT use for new feature work, bug fixes, frontend code, or greenfield service scaffolding.
---

# Refactoring Skill: Incrementally Clean Go Microservices

This skill defines the step-by-step workflow for an AI agent to incrementally refactor messy Go code into a clean, Domain-Driven Design (DDD), and CQRS-based architecture, heavily inspired by Martin Fowler’s refactoring principles.

## 1. Refactoring Philosophy

- **No Behavior Change**: Refactoring must permanently preserve the external behavior of the system. Never mix refactoring with bug fixes or feature additions.
- **Small, Safe Steps**: Execute small, localized changes. Commits should be atomic.
- **Test & Verify**: Always run tests or compile code after each incremental step.
- **Tell, Don't Ask**: Encapsulate data and move behavior to the objects that hold the data.
- **Continuous Improvement**: Leave the codebase cleaner and loosely coupled in each iteration.

---

## 2. Detection Phase (Smell Detection)

Before touching the codebase, analyze the files for these common code smells:

- **Long Function**: A function that spans more than ~50 lines and handles multiple distinct operations.
- **God Handler / Orchestrator**: A single handler or service managing DB connections, formatting external API requests, and running core business logic simultaneously.
- **Primitive Obsession**: Passing loosely typed primitives (e.g., `string`, `int`, `float64`) instead of rich Domain Objects (e.g., `MerchantID`, `Money`, `Email`).
- **Conditional Explosion**: Deeply nested `if/else` or gigantic `switch` blocks that control application behavior.
- **Feature Envy**: A function residing in module A, but extracting and manipulating variables exclusively from module B.
- **Mixed Abstraction Levels**: E.g., combining high-level logic (validating an account) with low-level protocol details (parsing a JWT token) in the same function.
- **Infrastructure Leakage**: Emitting SQL-specific errors, caching logic, or Kafka `msg` objects directly into the domain layer.

---

## 3. Refactoring Actions (Executable Operations)

Address each detected smell with **ONE** precise refactoring action per iteration:

- **Extract Function / Extract Method**
  - **Smell**: Long function / Mixed abstractions.
  - **Action**: Isolate a distinct chunk of code into a private, intention-revealing method. Ensure required parameters are explicitly passed.

- **Rename for Intention**
  - **Smell**: Vague or overly technical names (e.g., `Process()`, `CalcData()`).
  - **Action**: Rename to strictly reflect the business intent (e.g., `CalculateMonthlyInterest()`).

- **Move Function to Domain Layer (Tell, Don't Ask)**
  - **Smell**: Orchestrator forcibly mutating an object's fields directly.
  - **Action**: Shift the exact mutation logic to a method directly onto the target `Domain Model`.

- **Introduce Value Object**
  - **Smell**: Signatures requiring bare metrics (e.g., `amount float64, currency string`).
  - **Action**: Encapsulate these parameters into a domain struct (e.g., `Money{}`) to guard against invalid states.

- **Replace Conditional with Strategy / Polymorphism**
  - **Smell**: A giant `switch` dictating variations in logic.
  - **Action**: Define a common interface and delegate the logic to bounded implementations. Use a mapping/factory for the lookup execution.

- **Separate Calculation from Presentation (Introduce Intermediate DTO)**
  - **Smell**: Writing JSON marshalling tags or HTTP codes directly on a rich DB-coupled Domain model.
  - **Action**: Create pure `ResponseDTO` / `RequestDTO` structs in the handler layer and a simple map function (`ToDTO()`).

- **Encapsulate State**
  - **Smell**: Exported core entity fields randomly operated on across multiple handler files.
  - **Action**: Make the struct variables private. Provide controlled mutator functions indicating domain events (e.g., `account.ActivateProfile()`).

---

## 4. Layer Target Mapping

Refactored code must strictly align with the following domain boundaries:

- **[Handler Layer] / Transport**
  - **Duties**: Validates raw input interfaces/JSON, builds transport responses (HTTP).
  - **Restriction**: NO core business logic, NO repository calls.

- **[Orchestrator Layer] / Application**
  - **Duties**: Standard flow control. Retrieves objects from the Repository layer, delegates work to Domain methods, and persists results.
  - **Restriction**: NO domain decision-making (calculating discounts, validating state transitions).

- **[Domain Layer] / Core**
  - **Duties**: Pure Go structs, strict types, specific domain errors, and Value Objects. All dynamic business logic lives here locally via pointer methods.
  - **Restriction**: ZERO external framework dependencies (No HTTP routers, no SQL drivers, no messaging queues).

- **[Adapter Layer] / Infrastructure**
  - **Duties**: Implementation of outward/inward protocols (Postgres Repository, third-party REST client, gRPC handler).
  - **Restriction**: Exclusively maps transport/vendor data formats to pure Domain interfaces.

---

## 5. Execution Workflow (Agent Playbook)

During refactoring tasks, execute the following strict loop. DO NOT batch unrelated operations.

1. **Identify ONE Smell**: Analyze the selected file and identify exactly ONE distinct problem (e.g., "This orchestrator has business logic for calculating fees").
2. **Choose ONE Refactor Action**: Pick the matching refactoring technique (e.g., "Extract Function" + "Move Function to Domain").
3. **Apply Minimal Change**: Make exclusively the code edit required for that single action. Check boundary layers.
4. **Verify Behavior**: State how the isolated behavior works now. Emphasize running unit/integration tests without breaking the application build.
5. **Repeat**: Loop back to Step 1 until the file is cleanly aligned with the target architecture and the smell is extinguished.

---

## 6. Guardrails

- **DO NOT rewrite the entire module.** Refactoring is iterative. Large scale rewrites induce broken application states.
- **DO NOT introduce unnecessary abstractions.** Avoid complex frameworks if an idiomatic Go interface functions better. Stop when the code is clean and cohesive.
- **DO NOT mix behavior change with refactoring.** If the user asks for a structural cleanup, DO NOT slide feature additions in parallel. Complete the refactor first.
- **STOP when code is "good enough"**. Over-engineering violates lean design. Evaluate boundaries sensibly.

---

## 7. Output Contract

Per refactoring iteration, report exactly:

1. **Smell + action** — the one detected smell and the single technique chosen for it (e.g. "Orchestrator leak → Move Function to Domain + Extract Function").
2. **Minimal diff** — only the lines changed for that one action, as a unified diff or a fenced before/after snippet. No unrelated edits.
3. **Behavior preserved — evidence** — the verification actually run: tests executed and their result, or `go build` / `go vet` status, plus one line on why behavior is unchanged. If you could not run them, say so explicitly and state what the reviewer must run.

When no smell remains, emit a **stop line**: `Refactor complete — <file/package> aligned to <target layer>; N iterations; behavior preserved (verified by <tests/build>).` Do not continue past "good enough".

---

### Appendix: Go-Style Pseudo Refactor

**Smell: Primitive Obsession & Orchestrator Leak**
```go
// BEFORE
func (s *Service) ProcessAccount(balance float64, currency string, accountType string) error {
    if accountType == "premium" {
        balance = balance - 10.0 // Business logic leaked into the service layer
    }
    // ... Save balance ...
    return nil
}
```

**Refactored: Rich Objects & Tell, Don't Ask**
```go
// AFTER

// Domain Layer
type Money struct {
    Amount   float64
    Currency string
}

func (m Money) Subtract(o Money) Money {
    return Money{Amount: m.Amount - o.Amount, Currency: m.Currency}
}

type Account struct {
    Type    AccountType
    Balance Money // was: balance float64 — primitive promoted to a value object
}

func (a *Account) ApplyMonthlyFee() {
    if a.Type == AccountPremium {
        // Core business logic operates cleanly on the domain itself
        a.Balance = a.Balance.Subtract(Money{Amount: 10, Currency: "USD"})
    }
}

// Orchestrator Layer
func (s *Service) ProcessAccount(act *Account) error {
    act.ApplyMonthlyFee() // Tell, Don't Ask
    return s.repo.Save(act)
}
```
