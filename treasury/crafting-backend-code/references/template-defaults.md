# Microservice Template Defaults

Use these defaults only when they match the target template or the task is greenfield. Treat any existing Go services as examples of reusable backend patterns, not as a required deployment model or folder layout.

## Service Structure

- Prefer one deployable service per bounded capability, explicit bootstrap/config/router or transport wiring, health/readiness endpoints, dependency injection, typed config validation, graceful shutdown, and standard observability.
- Prefer thin handlers, use-case/service logic outside transport handlers, repository/adapters for persistence and external dependencies, and small interfaces at the consumer boundary.

## Go Services

- Prefer clear package boundaries, `context.Context` propagation, explicit errors, table-driven tests, and standard library primitives where sufficient.

## CQRS-Style Services

- Keep command handlers task-based and invariant-enforcing; keep query handlers side-effect free and read-model oriented.

## Event-Driven Flows

- Prefer additive event versioning, idempotent consumers, durable outbox/inbox semantics when crossing process boundaries, and observable retry/backoff behavior.

## PostgreSQL-Backed Services

- Prefer migrations, parameterized queries, explicit transaction scopes, indexes tied to query plans, and tests for concurrency-sensitive paths.

## HTTP APIs

- Prefer stable request/response contracts, typed validation, consistent error envelopes, OpenAPI/docs updates when already used, and backwards-compatible additions.

## Security-Sensitive Code

- Prefer vetted libraries, repo-approved crypto/auth helpers, centralized secret handling, and no custom cryptography.
