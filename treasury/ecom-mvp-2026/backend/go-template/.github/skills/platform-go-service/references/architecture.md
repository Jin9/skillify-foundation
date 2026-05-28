# Architectural Foundations

This template enforces four architectural pillars. Every convention traces back to one of them.

## 1. Domain-Driven Design (DDD)

- **Aggregate per package**: Each bounded context gets its own package under `app/<domain>/`.
- **Access layer**: Storage, cache, and external client interfaces live in `app/<domain>/access/` — the domain's **Repository / Gateway** layer in DDD terms. Each file bundles the interface, unexported implementation, and related domain models together.
- **Ubiquitous language**: Types, event names, and handler methods use the domain's vocabulary (`Product`, `PRODUCT_CREATED`, `OnProductCreated`).
- **Bounded context isolation**: Domain packages never import each other. If two domains need to collaborate, they do so via Kafka events or through the infrastructure layer.
- **Domain model**: Entities reside in `access/` alongside their repository interface. The model is the single source of truth for persistence shape (`firestore:`, `json:` struct tags).

## 2. CQRS (Command Query Responsibility Segregation)

- **Write path (Commands)**: HTTP `POST` endpoints and Kafka consumers that mutate state.
- **Read path (Queries)**: HTTP `POST` endpoints that retrieve state (`/detail`, `/list`, `/current`, `/me`).
- **Event-driven side effects**: Write operations on one aggregate emit domain events via Kafka. Other aggregates subscribe to those events.
- **Separate handler files**: `handler_<action>.go` (HTTP) and `consumer_<action>.go` (Kafka) naturally separate command and query codepaths.

## 3. Fowler-Style Patterns

| Pattern | Where Applied |
|---|---|
| **Repository** | `access/storage_<dep>.go` — interface + concrete impl wrapping the data store (Firestore, MySQL, S3/GCS) |
| **Cache** | `access/cache_<dep>.go` — interface + impl wrapping a cache layer (Redis, Memcached) |
| **Service Layer** | `handler` struct methods — thin orchestration in `handler_<action>.go`; private helpers extracted to `service_<action>.go` when purpose differs |
| **Domain Model** | Entity structs in `access/` (`Product`, `Member`) |
| **Gateway** | `access/client_<dep>.go` — wraps external HTTP APIs behind an interface |
| **Dependency Injection** | `router/deps.go` — composition root, manual DI (no framework) |
| **Event-Driven** | `router/subscriber.go` — event router dispatches to handler methods |
| **Value Object** | Status types as `type ProductStatusType string` with defined constants |
| **Continuous Refactoring** | Behavior-preserving, small-step refactoring. See `refactoring-fowler.md`. |

## 4. Go Idioms

- **Accept interfaces, return structs**: `NewHandler(cfg HandlerConfig) *handler` returns a concrete type; `NewProductStorage(fs *gcpfirestore.Client) ProductStorage` returns the interface.
- **Compile-time interface checks**: `var _ ProductStorage = (*productStorage)(nil)`.
- **Unexported implementation types**: `handler` and `productStorage` are unexported; only the constructor and interface are public.
- **`sync.Once` singletons**: Config loading is lazy and safe for concurrent use.
- **Error wrapping with `%w`**: Storage methods use `fmt.Errorf("failed to ...: %w", err)`.
- **Structured logging via `log/slog`**: No `log.Printf`, no third-party loggers.
- **Table-driven tests**: With local `mockArgs` / `args` / `want` types and a `prepare` function.
- **Context propagation**: Every I/O method accepts `context.Context` as first argument.
- **Embed directives**: `//go:embed VERSION` for build metadata.
- **Graceful shutdown**: `signal.NotifyContext` + deferred teardown chain.

---

## Architecture Overview

```text
<root>/
├── main.go                  # Bootstrap: config → deps → router + subscriber → graceful shutdown
├── go.mod                   # Module path  (replace directive points common → common.git vX.Y.Z)
├── VERSION                  # Semver file (read via //go:embed)
├── Makefile                 # Dev & CI targets
├── Dockerfile               # Multi-platform alpine build
├── docker-compose.yml       # Local infra
├── gitlabci.yml             # CI pipeline
├── .env.template            # Environment variable reference
├── .golangci.yaml           # Linter config (golangci-lint v2)
├── .mockery.yaml            # Mockery v2 config
├── config/
│   ├── config.go            # Nested Config struct + sync.Once loader C()
│   └── config_test.go
├── router/
│   ├── deps.go              # deps struct — all shared infrastructure clients (Composition Root)
│   ├── router.go            # Gin engine: middleware chain, health probes, route groups
│   └── subscriber.go        # Kafka event→handler routing
└── app/
    └── <domain>/
        ├── handler.go                   # HandlerConfig + handler struct + NewHandler
        ├── handler_<action>.go          # HTTP handler (Command or Query)
        ├── handler_<action>_test.go
        ├── consumer_<action>.go         # Kafka consumer (Event handler)
        ├── consumer_<action>_test.go
        ├── service_<action>.go          # Private service helpers grouped by purpose
        └── access/
            ├── storage_<dependency>.go
            ├── cache_<dependency>.go
            ├── client_<dependency>.go
            └── mocks/mocks.go
```

## Dependency Flow (DDD Onion)

```text
  main.go
    └── router/ (Composition Root — infrastructure wiring)
          ├── deps.go       ← creates infra clients (Firestore, MySQL, Redis, Kafka, HTTP)
          ├── router.go     ← injects deps into domain handlers
          └── subscriber.go ← injects deps into event handlers
               └── app/<domain>/     (Application/Domain layer)
                     ├── handler.go          ← depends on access interfaces only
                     ├── handler_<action>.go  ← orchestrates domain logic
                     ├── service_<action>.go  ← private helpers
                     ├── consumer_<action>.go ← handles domain events
                     └── access/             (Infrastructure/Repository layer)
                           ├── storage_<dep>.go ← implements interfaces against real infra
                           ├── cache_<dep>.go
                           └── client_<dep>.go
```

Dependencies point **inward**: domain handlers depend on `access/` interfaces; `access/` implementations depend on infrastructure SDKs. The `router/` package is the **composition root** that assembles everything.

## Two categories of dependencies in `deps`

1. **Raw infrastructure SDK handles** (Firestore, MySQL, Redis, S3/GCS, `*http.Client`) — **never** passed directly to domain handlers. Wrapped in `access/` behind a domain-specific interface (`storage_`, `cache_`, or `client_` files). The `register<Domain>Routes` / `register<Domain>Events` functions create the access-layer implementations and inject them via `HandlerConfig`.

2. **Common-module abstractions** (`hash.HashManager`, `crypt.Cipher`, `token.JWTSigner`, `kafka.Producer`) — already interfaces from the common module; passed directly to `HandlerConfig`. They don't need an `access/` wrapper because they are already infrastructure-agnostic and mockable.
