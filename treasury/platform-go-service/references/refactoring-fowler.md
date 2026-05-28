# Refactoring Workflow — Fowler Principles

This repository enforces continuous, behavior-preserving refactoring based on Martin Fowler's principles. Use this iterative playbook when modifying existing code, eliminating tech debt, or addressing messy legacy implementations.

---

## 1. Refactoring Philosophy

- **No behavior change**: Refactoring strictly improves internal structure and design without altering observable behavior.
- **Small safe steps**: Do not attempt sweeping rewrites. Apply minimal, atomic changes.
- **Test/verify constantly**: Run tests, check linters, or verify runtime logs after every single step.

---

## 2. Layer Target Mapping (crucial)

Enforce strict separation of concerns and proper ownership of logic:

| Layer | Owns | Does NOT own |
|---|---|---|
| **Handler** | Decoding HTTP/Kafka requests, basic format validation, encoding responses | Flow control, business rules |
| **Service / Orchestrator** | Sequencing steps, fetching from DB, calling external APIs | Business rules, transport |
| **Domain Model** | All business logic — state mutation, calculations, invariants (Tell, Don't Ask) | Persistence detail, transport |
| **Adapter** (`access/`) | External protocol handling — translate domain requests to infrastructure | Business logic, flow control |

---

## 3. Detection Phase — code smells to scan for

- **God Handler / Service**: Functions mixing HTTP processing, DB orchestration, and deep business logic.
- **Primitive Obsession**: Loose primitives (`string`, `int`) passed around instead of Domain or Value Objects (e.g. `ProductStatusType`).
- **Feature Envy**: A function intensely interrogating the inner data of another object instead of asking that object to perform the operation.
- **Conditional Explosion**: Deeply nested `if/else` or large `switch` blocks handling specialised logic.
- **Mixed Abstraction Levels**: High-level flow (e.g. `UpdateProduct()`) mixed with byte-manipulation in the same function.
- **Leakage**: Repository details (SQL connection handles, Firestore constructs) showing up in the domain orchestration layer.

---

## 4. Refactoring Actions

Apply these tactics safely and incrementally:

| Action | What it does |
|---|---|
| **Move Function** | Push business rules into Domain Model structs ("Tell, Don't Ask"). |
| **Extract Function & Rename** | Break down long methods into sub-functions with clear, intention-revealing names. |
| **Introduce Value Objects** | Replace `string`/`int` parameters with strongly typed domain definitions. |
| **Replace Conditionals with Strategy** | Map complex, expanding conditional logic into a strategy map / polymorphism. |
| **Separate Formatting** | Decouple data calculation from output formatting via intermediate DTOs / View Models. |
| **Encapsulate State** | Hide internal struct variables and manage access to reduce shared mutable data. |

---

## 5. Execution Workflow (Agent Playbook)

1. **Identify** one instance of a code smell or structural issue.
2. **Choose** a single isolated refactoring action.
3. **Apply** the minimal required change to fix the isolated issue.
4. **Verify** behavior (via unit tests or syntax checks).
5. **Repeat** the cycle until the module is clean.

---

## 6. Guardrails

- **DO NOT** rewrite entire modules or domains from scratch.
- **DO NOT** introduce preemptive or unnecessary abstractions. Keep things concrete until flexibility is actually required.
- **DO NOT** mix feature changes (changing behavior) with refactoring. Do one, commit, then do the other.
- **STOP** when the code cleanly separates domain logic from infrastructure and is "good enough".
