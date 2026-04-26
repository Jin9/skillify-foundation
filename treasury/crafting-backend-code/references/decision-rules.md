# Backend Decision Rules

- **Pattern first for templates**: extract reusable backend service patterns from examples without copying example package names, deployment shape, or business capabilities as mandatory rules.
- **Repo first for existing services**: match the existing architecture, naming, errors, logging, validation, dependency injection, tests, and package layout before importing a preferred pattern.
- **Contracts before code**: for APIs and messages, decide request/response shape, versioning, compatibility, errors, and authorization before implementation.
- **One owner per piece of state**: do not let multiple services or packages mutate the same table, aggregate, cache key, or event stream without a clear ownership rule.
- **Transactions are explicit**: name what is atomic, what is eventually consistent, and what compensation or retry handles partial failure.
- **Idempotency is required for retries**: commands, webhook handlers, queue consumers, and outbox dispatchers must tolerate duplicate delivery where duplicates are plausible.
- **Context propagates**: carry cancellation, deadlines, trace IDs, auth claims, and request metadata through service, repository, and client calls according to repo conventions.
- **Errors are operational signals**: preserve cause, classify client vs server vs dependency errors, avoid panics in request/worker paths, and avoid leaking secrets or internals in public responses.
- **Observability follows failure modes**: add logs, metrics, traces, and readiness checks where they help operators detect and diagnose real failures.
- **Security is not a cleanup task**: check authN/authZ, input validation, SSRF/deserialization risks, secret handling, SQL injection, PII logging, and crypto usage before shipping backend changes.
- **Generated and migrated artifacts are special**: do not hand-edit generated files or committed migration outputs unless the discovered workflow requires it; regenerate or add migrations using target tooling.
