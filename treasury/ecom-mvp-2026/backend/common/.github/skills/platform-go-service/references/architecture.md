# Service Architecture

Layout, DDD vocabulary, and CQRS request flow for an A-Team Go microservice that consumes the `common` library.

## Directory Layout

```
<root>/
├── main.go                  # Bootstrap: config → logger → deps → router → graceful shutdown
├── go.mod                   # gitdev.devops.krungthai.com/a-team/poc-structure/<service>
├── Makefile                 # Dev targets
├── Dockerfile               # Multi-stage alpine build
├── docker-compose.yml       # Local infra (Postgres, Redis, Kafka, emulators)
├── gitlabci.yml             # CI pipeline
├── config/
│   ├── config.go            # Config struct (env tags)
│   ├── env.go               # Env-name constants
│   └── parser.go            # Load + validate via config.ParseEnv[T]
├── router/
│   ├── deps.go              # Deps struct — all shared dependencies
│   ├── router.go            # Gin engine, middleware stack, route registration
│   └── subscriber.go        # Kafka topic→handler wiring
├── app/
│   ├── codes.go             # Domain response codes (extends wrapper via type alias)
│   ├── timestamps.go        # SetTimestamps helper
│   └── <aggregate>/         # One sub-package per DDD aggregate
│       ├── handler.go               # HandlerConfig struct + NewHandler constructor
│       ├── <action>_handler.go      # HTTP handler (command or query)
│       ├── <action>_handler_test.go # HTTP handler test (100% coverage)
│       ├── <event>_consumer.go      # Kafka consumer handler
│       ├── <event>_consumer_test.go # Kafka consumer test
│       └── access/
│           ├── <port>.go            # Interface per external dependency (Port)
│           └── model.go             # Domain/persistence models
```

## DDD Vocabulary

| DDD Term | Implementation |
|---|---|
| Bounded Context | One Git repository per microservice |
| Shared Kernel | The `common` module — infrastructure only, no domain logic |
| Aggregate | One sub-package under `app/<aggregate>/` |
| Port (Hexagonal) | Interface files in `access/` (e.g. `Storage`, `Producer`) |
| Adapter | Concrete implementations wired in `router/deps.go` |
| Domain Event | `kafka.Message[T]` envelope with `event_name`, `aggregate_id`, `event_id` |
| Application Service | `handler` struct with injected port interfaces |
| Value Object | `app.Code`, `app.Message`, `access.Model` |

## CQRS Flow

```
HTTP POST/PUT/DELETE → handler (command) → access.Storage → DB
                                         → kafka.Producer  → topic
                                         ← wrapper.Respond(code, message, data)

Kafka topic          → NewEventRouter    → KafkaHandler (command on consumer side)
                                         → access.Storage → DB

HTTP GET             → handler (query)   → access.Storage → DB
                                         ← wrapper.Respond(code, message, data)
```
