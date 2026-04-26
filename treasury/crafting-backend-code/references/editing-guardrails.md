# Editing Guardrails

## Read Before Write

- Read before writing. Do not rewrite files from memory or apply broad mechanical changes without inspecting current code.

## Do Not Introduce Without Request

- Do not introduce framework migrations, new ORMs, new messaging clients, new auth stacks, new config systems, or new CI gates unless requested.

## Protected Artifacts — Require Explicit Approval

- Do not change authentication, authorization, token/session storage, public API semantics, message schemas, database schema, generated clients, or persisted data formats without explicit approval or a migration plan accepted by the user.

## Scope Discipline

- Do not delete files, mass-format unrelated files, or clean up code outside the requested scope.
- Keep dependency changes separate from feature fixes when possible. If a dependency is required, explain why the existing stack cannot solve the problem.

## Safe Change Patterns

- Prefer feature flags, compatibility wrappers, additive fields, shadow reads/writes, or dual-read rollout patterns for risky backend changes.
- Preserve public APIs and package interfaces unless the task is explicitly a breaking refactor.
