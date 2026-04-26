# Testing Strategy

## Contents
- Testing layers
- Principles
- Go test conventions
- Load / performance

## Testing layers

| Layer | Scope | Tooling | Owns |
|---|---|---|---|
| Unit | Single function/method | Native test framework | Business logic, pure transforms |
| Integration | Service + real dependencies | Testcontainers, docker-compose | Repository, adapter correctness |
| Contract | Cross-domain API surface | Pact, schema validation | API compatibility between domains |
| E2E | Full flow, staging env | API tests, seeded data | Critical business paths only |

## Principles

- Unit tests: fast, no IO, table-driven (Go). No mocking frameworks — use interfaces.
- Integration tests: real DB (Testcontainers), real queue. No in-memory fakes for critical paths.
- Contract tests: MUST exist for every cross-domain API. Producer and consumer both verify.
- E2E: minimal set. Cover happy path + top 3 failure scenarios only.
- NEVER mock what is not owned. Wrap external dependencies behind interfaces; test the wrapper with integration tests.

## Go test conventions

- Use `testify`: `assert` for soft checks, `require` for hard stops, `mock` for dependency stubs.
- Build tags for separation: `//go:build integration` for integration tests.
- CI runs unit tests by default; integration tests run via explicit `-tags integration`.
- Deterministic tests: NEVER call `time.Now()` directly — inject a clock. NEVER rely on map iteration order.
- Table-driven tests: every case MUST have a descriptive name.

## Load / performance

- Define SLOs before writing load tests (p99 latency, throughput target).
- Use realistic data volumes — not toy datasets.
- Run against staging with production-like topology.
