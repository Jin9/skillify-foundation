---
name: refactoring-go-services
description: >
  Incrementally refactor messy Go microservices toward clean DDD/CQRS
  architecture while preserving behavior — one code smell and one Fowler-style
  refactoring action per iteration, verified by tests or build after each step.
  Use when the user asks "clean up this Go service", "refactor this handler",
  "move logic into the domain layer", or "remove Go service code smells". Do NOT
  use for new feature work, bug fixes, frontend code, or greenfield service
  scaffolding; not for deciding whether a refactor is worth doing (use
  refactor-decision); and not for broad go-template service conventions or
  scaffolding (use platform-go-service).
---

# Refactoring Skill: Incrementally Clean Go Microservices

Step-by-step workflow for incrementally refactoring messy Go code into a clean,
Domain-Driven Design (DDD) / CQRS architecture, inspired by Martin Fowler's
refactoring principles. Work one smell and one action at a time; never batch.

## 1. Refactoring philosophy

- **No behavior change**: refactoring must preserve the external behavior of the system. Never mix refactoring with bug fixes or feature additions.
- **Small, safe steps**: execute small, localized changes; commits are atomic.
- **Test & verify**: always run tests or compile the code after each incremental step.
- **Tell, don't ask**: encapsulate data and move behavior to the objects that hold the data.
- **Continuous improvement**: leave the codebase cleaner and more loosely coupled in each iteration.

## 2. Execution workflow (agent playbook)

Run this strict loop. DO NOT batch unrelated operations.

1. **Identify ONE smell.** Analyze the selected file and name exactly one
   distinct problem (e.g. "this orchestrator has business logic for calculating
   fees"). Detection catalog: `references/code-smells.md`.
2. **Choose ONE refactor action.** Pick the matching technique (e.g. "Extract
   Function" + "Move Function to Domain"). Action catalog with smell → action
   mappings: `references/refactoring-actions.md`.
3. **Apply the minimal change.** Make exclusively the edit required for that
   single action, and keep each layer within its boundary —
   `references/layer-target-mapping.md` (Handler / Orchestrator / Domain / Adapter
   duties and restrictions).
4. **Verify behavior.** State how the isolated behavior works now; emphasize
   running unit/integration tests or `go build` / `go vet` without breaking the
   build.
5. **Repeat** from step 1 until the file is cleanly aligned with the target
   architecture and the smell is extinguished.

A complete before/after walk-through (Primitive Obsession + Orchestrator Leak) is
in `examples/go-refactor-example.md`.

## 3. Guardrails

- **DO NOT rewrite the entire module.** Refactoring is iterative; large-scale rewrites induce broken application states.
- **DO NOT introduce unnecessary abstractions.** Prefer an idiomatic Go interface over a complex framework. Stop when the code is clean and cohesive.
- **DO NOT mix behavior change with refactoring.** Complete the structural cleanup first; do not slide feature additions in parallel.
- **STOP when the code is "good enough".** Over-engineering violates lean design; evaluate boundaries sensibly.

## 4. Output contract

Per refactoring iteration, report exactly:

1. **Smell + action** — the one detected smell and the single technique chosen
   (e.g. "Orchestrator leak → Move Function to Domain + Extract Function").
2. **Minimal diff** — only the lines changed for that one action, as a unified
   diff or a fenced before/after snippet. No unrelated edits.
3. **Behavior preserved — evidence** — the verification actually run: tests
   executed and their result, or `go build` / `go vet` status, plus one line on
   why behavior is unchanged. If you could not run them, say so explicitly and
   state what the reviewer must run.

When no smell remains, emit a **stop line**:
`Refactor complete — <file/package> aligned to <target layer>; N iterations; behavior preserved (verified by <tests/build>).`
Do not continue past "good enough".

## References

| Need | File |
|------|------|
| Smell detection catalog (step 1) | `references/code-smells.md` |
| Smell → refactoring-action mappings (step 2) | `references/refactoring-actions.md` |
| DDD layer duties & restrictions (step 3) | `references/layer-target-mapping.md` |
| Worked before/after example | `examples/go-refactor-example.md` |
