# Go Service Template

Blueprint repository for scaffolding A-Team platform microservices in Go. Clone this repo, replace the example domains under `app/`, and you have a production-ready service with HTTP routing, Kafka event handling, structured logging, and graceful shutdown — all wired up.

## Architectural Foundation

This template codifies four architectural pillars. Every file and convention traces back to one of them.

| Pillar | How It's Applied |
|--------|-----------------|
| **Domain-Driven Design** | One package per aggregate under `app/<domain>/`. The `access/` sub-package wraps infrastructure behind domain-specific interfaces: `storage_` for persistence, `cache_` for caching, `client_` for external APIs. Domain models co-locate with their access file. No cross-domain imports. |
| **CQRS** | Write paths (create, update, delete) and read paths (get, list) are separate handler files. Kafka events propagate state changes across bounded contexts. |
| **Fowler-Style Patterns** | Behavior-preserving refactoring via small safe steps. Logic moves to Domain Models (Tell, Don't Ask). Strict layer targets: Handlers for transport, Services for flow control, Domain for business rules, and Adapters (`storage_`, `cache_`, `client_`) for external protocols. |
| **Go Idioms** | Accept interfaces/return structs, compile-time interface checks, `sync.Once` singletons, `context.Context` propagation, `%w` error wrapping, table-driven tests, `log/slog` structured logging. |

## What Lives Where

**Core Application**

```text
main.go                    Bootstrap, graceful shutdown, runtime tuning
config/config.go           Environment-backed configuration (nested structs + sync.Once)
router/deps.go             Composition root — infrastructure client creation & cleanup
router/router.go           Gin engine, middleware chain, HTTP route registration
router/subscriber.go       Kafka event→handler routing, consumer lifecycle
app/<domain>/              One sub-package per business aggregate
app/<domain>/access/       Repository / Cache / Gateway layer (storage_, cache_, client_)
```

**Shared Infrastructure (external module)**

Infrastructure and utility packages are imported from the shared **common** module at `gitlab.com/b2c-e-commerce-platform/platform/backend/common`:

| Package | Purpose |
|---------|---------|
| `logger` | Structured `slog`-based logging with key replacers and field censoring |
| `middleware` | SecurityHeaders, CORS, RefID, TraceContext, AccessLog, Timeout |
| `wrapper` | Generic HTTP response envelope (`ResponseOption[T]`, `BindJSON[T]`) |
| `serror` | Source-location error wrapping with investigation context |
| `kafka` | Producer/Consumer, EventRouter, BindMessage, logging interceptor |
| `firestore` | Firestore client wrapper |
| `database` | MySQL connection helper |
| `redis` | Redis client factory |
| `httpclient` | HTTP client with middleware options |
| `crypt` | AES-GCM encryption |
| `hash` | bcrypt hashing with pepper |
| `token` | JWT signing (ES256) |
| `codec` | Base64 encoding/decoding |
| `config` | Env parsing, `ParseEnv[T]`, environment helpers |
| `health` | Liveness, readiness, Prometheus metrics |
| `app` | Shared constants (`CodeSuccess`, `SetTimestamps`) |

These packages do **not** live in this repository — services import them, never vendor or copy them.

---

## Quick Start

### Prerequisites

- Go 1.25+
- Access to the private module host `gitlab.com/b2c-e-commerce-platform/platform/backend`
- A valid `.env` file based on `.env.template`

Set `GOPRIVATE` for private modules:

```bash
go env -w GOPRIVATE=gitlab.com/b2c-e-commerce-platform/platform/backend
```

### Local Setup

1. Copy `.env.template` to `.env` and fill in the values.
2. Run `make setup` to install project tools and hooks.
3. Run `make test` to confirm the workspace is healthy.
4. Run `make run` to start the service in Docker.

---

## Generating a New Microservice

This repository is a **template** — it provides the scaffolding for new services. The example domains (`auth`, `member`, `organization`, `product`) demonstrate the conventions. When starting a new service:

1. **Clone** this repo (or use it as a GitLab template)
2. **Rename** the module in `go.mod`
3. **Delete** all `app/` domains
4. **Create** your own domains following the patterns in `app/product/` (see "How to Extend" below)
5. **Update** `config/config.go` to match your service's configuration
6. **Update** `router/router.go` and `router/subscriber.go` to wire your domains
7. **Update** `.env.template`, `VERSION`, `CHANGELOG.md` and CI variables

Keep everything else intact: `main.go`, `router/deps.go` structure, `Makefile`, `Dockerfile`, `.golangci.yaml`, `.mockery.yaml`.

---

## Using This Template with an AI Coding Assistant

This repository ships a self-contained skill at `.github/skills/platform-go-service/` that captures every convention the template enforces. When loaded, an LLM-based coding assistant (Claude Code, Cursor, GitHub Copilot Chat, Codex CLI, etc.) will scaffold new domains, handlers, consumers, and tests **in the exact shape this template expects** — file naming, `access/` layout, `serror.Wrap`, `wrapper.Respond`, `kafka.BindMessage`, table-driven tests with 100% coverage, and the Fowler refactoring playbook.

### What's in the skill

```text
.github/skills/platform-go-service/
├── SKILL.md                          # Lean entry point: triggers, conventions, workflow table
└── references/
    ├── architecture.md               # DDD / CQRS / Fowler / Go-idiom pillars + dep flow
    ├── scaffolding.md                # main.go, config/, router/{deps,router,subscriber}.go
    ├── domain-patterns.md            # access/, handler.go, handler_<action>.go, consumer_<action>.go
    ├── testing.md                    # mockArgs/args/want/prepare table-driven pattern
    ├── refactoring-fowler.md         # Smell detection + safe refactor actions
    ├── common-module.md              # Quick reference for the shared common/* packages
    └── dockerfile-and-ci.md          # Dockerfile, Makefile targets, GitLab pipeline
```

The `SKILL.md` is the entry point — assistants load it first, then progressively read references as the task requires.

### Apply it to your assistant

| Tool | One-time setup | Activation |
|---|---|---|
| **Claude Code** (`claude`) | Symlink or copy the skill into your skills dir: `ln -s "$(pwd)/.github/skills/platform-go-service" ~/.claude/skills/` (user-scope) **or** `mkdir -p .claude/skills && ln -s ../../.github/skills/platform-go-service .claude/skills/` (project-scope) | Auto-discovered; the description triggers on the phrases below. You can also force it: `Use the platform-go-service skill to ...` |
| **Cursor** | Add a rule in `.cursor/rules/platform-go-service.mdc` that references the SKILL.md content (or `@`-mention the file in chat) | Mention the rule, or `@.github/skills/platform-go-service/SKILL.md` in chat |
| **GitHub Copilot Chat** | Reference the file in context: `#file:.github/skills/platform-go-service/SKILL.md` | Add the `#file:` reference at the top of your prompt |
| **Codex CLI / ChatGPT** | Paste `SKILL.md` contents into the system prompt or attach the folder as project context | Mention "follow the platform-go-service skill" in the user prompt |
| **Generic / other LLMs** | Provide `SKILL.md` (and any referenced files it asks for) as system context | Prompt with one of the trigger phrases below |

### Trigger phrases (already in the skill description)

The skill description lists the phrases that should trigger it. Use any of these in your prompt and the assistant will follow the conventions:

- "create a new Go service / microservice"
- "add a new domain / aggregate / bounded context"
- "add an HTTP endpoint / handler"
- "add a Kafka consumer / event handler"
- "wire a new infrastructure client" (Firestore, MySQL, Redis, S3/GCS, external API)
- "generate mocks" / "write tests with 100% coverage"
- "refactor this handler/service" (invokes the Fowler workflow)

### Example prompts

```text
Add a new `invoice` domain aggregate with HTTP endpoints for create / detail /
list and a Kafka consumer for `INVOICE_PAID`. Persist in Firestore. Follow the
platform-go-service skill.
```

```text
Add a POST /api/v1/platform/product/archive endpoint to the existing product
domain. Wire it in router.go and write tests reaching 100% coverage.
```

```text
Refactor app/auth/handler_issue_token.go using the Fowler workflow in the
platform-go-service skill — the handler is mixing transport, orchestration, and
business rules.
```

### Verify the assistant respects the skill

Before merging AI-generated code, check it against the conventions checklist in `SKILL.md` §7. The most common drift signals are:

- Using `json.Unmarshal` instead of `kafka.BindMessage` in a consumer.
- Returning a raw error to `wrapper.Respond` instead of wrapping with `serror.Wrap(err).With(...)`.
- Putting models in a separate `model.go` rather than co-locating in the `access/` file that uses them.
- Hand-written mocks instead of running `mockery`.
- Missing test cases for `if err != nil` branches (coverage drops below 100%).

If the assistant ignores the skill, restate the trigger phrase explicitly (`"Use the platform-go-service skill to ..."`) or paste `SKILL.md` directly into the chat as context.

---

## Project Structure

```text
.
├── app/
│   ├── auth/
│   │   ├── handler.go                       # HandlerConfig + constructor
│   │   ├── handler_issue_token.go           # POST /auth/issue-token
│   │   ├── handler_resolve_identity.go      # POST /auth/resolve-identity
│   │   ├── service_authenticate_google.go   # Google OAuth verification logic
│   │   ├── service_resolve_existing_member.go
│   │   └── access/
│   │       ├── client_google.go             # client_ = external API gateway
│   │       ├── storage_member.go            # storage_ = persistence
│   │       ├── storage_organization.go
│   │       └── mocks/
│   ├── member/
│   │   ├── handler.go
│   │   ├── handler_get_info.go              # POST /member/me
│   │   ├── handler_register.go             # POST /member/register
│   │   ├── consumer_create.go              # MEMBER_CREATED event
│   │   └── access/
│   ├── organization/
│   │   ├── handler.go
│   │   ├── handler_get_current.go           # POST /organization/current
│   │   ├── handler_register.go             # POST /organization/register
│   │   ├── consumer_create.go              # ORGANIZATION_CREATED event
│   │   ├── consumer_add_member.go          # ORGANIZATION_MEMBER_ADDED event
│   │   └── access/
│   └── product/
│       ├── handler.go
│       ├── handler_create.go               # POST /product/create
│       ├── handler_delete.go               # POST /product/delete
│       ├── handler_get.go                  # POST /product/detail
│       ├── handler_list.go                 # POST /product/list
│       ├── handler_update.go               # POST /product/update
│       ├── consumer_create.go              # PRODUCT_CREATED event
│       ├── consumer_delete.go              # PRODUCT_DELETED event
│       ├── consumer_update.go              # PRODUCT_UPDATED event
│       └── access/
│           ├── storage_product.go           # storage_ = Firestore persistence
│           └── mocks/
├── config/
│   ├── config.go
│   └── config_test.go
├── router/
│   ├── deps.go
│   ├── router.go
│   └── subscriber.go
├── .github/
│   └── skills/                              # AI coding assistant skills
├── .scripts/                                # Dev setup, CI, deployment scripts
├── main.go
├── go.mod / go.sum
├── Makefile
├── Dockerfile
├── docker-compose.yml
├── gitlabci.yml
├── spec.md                                  # API request/response spec
├── VERSION / CHANGELOG.md
├── .env.template
├── .golangci.yaml
└── .mockery.yaml
```

Each `app/<domain>/access/` directory contains domain-specific infrastructure wrappers, categorized by prefix:

| Prefix | Fowler Pattern | Example |
|--------|---------------|----------|
| `storage_<dep>.go` | Repository | Firestore, MySQL, PostgreSQL, S3/GCS |
| `cache_<dep>.go` | Cache | Redis, Memcached |
| `client_<dep>.go` | Gateway | HTTP APIs, gRPC services |

Each file holds the interface, implementation, and domain model together — no separate `model.go` or `errors.go`. A `mocks/` subdirectory contains generated test doubles.

---

## Package Usage Examples

### `logger` — Structured Logging

Initialize once in `main.go`. Replacers are composable — each maps slog keys to platform-specific names.

```go
import "gitlab.com/b2c-e-commerce-platform/platform/backend/common/logger"

// Initialize with cloud-specific key replacers and sensitive-field masking.
_ = logger.New(
    logger.AWSKeyReplacer,        // msg→message, time→timestamp
    logger.OpenSearchKeyReplacer, // msg→message, time→@timestamp, level→log.level
    logger.CensorReplacer,        // masks password, access_token, api_key, secret
)
```

Register additional sensitive fields at startup:

```go
logger.AddCensor("card_number", "****-****-****-****")
logger.AddCensor("cvv", "***")

// Or bulk-load from config:
logger.LoadCensorsFromMap(map[string]string{
    "email": "***EMAIL***",
    "ssn":   "***SSN***",
})
```

---

### `serror` — Structured Error Wrapping

Wraps errors with source-location tracking and investigation context. The observability middleware automatically extracts and logs these fields.

```go
import "gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"

// Wrap an existing error (captures file:line:func)
if err := db.Insert(ctx, record); err != nil {
    return serror.Wrap(err)
}

// Create a new error
return serror.New("validation failed")

// Attach investigation context (flows to access log automatically)
return serror.Wrap(err).With(
    slog.String("organization_id", orgID),
    slog.String("product_name", req.Name),
)
```

---

### `wrapper` — Standard HTTP Response

Generic response envelope with metadata for the access log middleware.

```go
import "gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
import "gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"

// Request binding (returns (T, bool); auto-responds 400 on failure)
req, ok := wrapper.BindJSON[CreateProductRequest](c, slog.String("handler", "CreateProduct"))
if !ok {
    return
}

// Success response
wrapper.Respond(c, wrapper.ResponseOption[CreateProductResponse]{
    HTTPStatus: http.StatusCreated,
    Code:       app.CodeSuccess,
    Message:    app.MessageSuccess,
    Data:       &CreateProductResponse{ProductID: productID},
})

// Error response (Err flows to AccessLog middleware for structured logging)
wrapper.Respond(c, wrapper.ResponseOption[CreateProductResponse]{
    HTTPStatus: http.StatusInternalServerError,
    Code:       app.CodeInternalError,
    Message:    app.MessageInternalError,
    Err: serror.Wrap(err).With(
        slog.String("product_name", req.Name),
    ),
})
```

---

### `config` — Environment-Based Configuration

Loads configuration from environment variables using struct tags.

```go
import "gitlab.com/b2c-e-commerce-platform/platform/backend/common/config"
import myconfig "gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/config"

// Load the full application config (prefix-based)
cfg := myconfig.C(config.Env)

// Environment checks
if config.IsLocalEnv() {
    // local-only setup
}
if config.IsProdEnv() {
    // production-only setup
}
```

---

### `kafka` — Producer Logging Interceptor

Wraps a producer with structured logging. Three modes: Debug, Meta, Silent.

```go
import "gitlab.com/b2c-e-commerce-platform/platform/backend/common/kafka"

producer := kafka.MustNewProducer(cfg)

// Auto-select mode from ENV (LOCAL/DEV→Debug, UAT/PROD→Meta)
logged := kafka.WithLoggingFromEnv(producer, os.Getenv("ENV"))

// Or explicit mode
debug  := kafka.WithLogging(producer, kafka.LogInterceptorOption{Mode: kafka.LogModeDebug})  // full payload
meta   := kafka.WithLogging(producer, kafka.LogInterceptorOption{Mode: kafka.LogModeMeta})   // metadata only
silent := kafka.WithLogging(producer, kafka.LogInterceptorOption{Mode: kafka.LogModeSilent}) // no logging

// In PROD (Meta mode), attach investigation attrs — no payload is logged
logged.SendMessageWithOption("topic", msg, kafka.SendMessageOption{
    LogAttrs: []slog.Attr{
        slog.String("organization_id", orgID),
    },
})
```

---

### `kafka` — Consumer & Event Router

Event-driven consumer with routing by event name.

```go
import "gitlab.com/b2c-e-commerce-platform/platform/backend/common/kafka"

// Create a consumer group
group := kafka.MustNewConsumerGroup(kafka.ConsumerConfig{
    KafkaConf: kafka.NewConsumerConfigAtLeastOnce(),
    Brokers:   []string{"broker-1:9092"},
    GroupID:   "go-template",
})

// Register event handlers per aggregate
handlers := map[string]kafka.KafkaHandler{
    "PRODUCT_CREATED": productHandler.OnProductCreated,
    "PRODUCT_UPDATED": productHandler.OnProductUpdated,
    "MEMBER_CREATED":  memberHandler.OnMemberCreated,
}

// Create event router (acts as access-log for the async path)
processor := kafka.NewEventRouter(handlers)
handler := kafka.NewConsumerGroupHandler(ctx, processor)

// Consume
group.Consume(ctx, []string{"events"}, handler)
```

---

## Runtime Flow

1. `main.go` loads configuration and initializes structured logging (`logger.New`).
2. `router.StartSubscriber(...)` starts the Kafka consumer group (if broker settings are present). This starts **before** HTTP to ensure event processing begins immediately.
3. `router.New(...)` builds the Gin engine, creates infrastructure clients via `newDeps`, registers middleware and all API routes.
4. The HTTP server starts with graceful shutdown wired via `signal.NotifyContext`.
5. On SIGINT/SIGTERM: HTTP server drains → Kafka consumer stops → infrastructure clients close → process exits.

---

## Configuration

Configuration is loaded from environment variables in `config/`.

| Area | Variables |
|---|---|
| Server | `ENV`, `HOSTNAME`, `PORT` |
| CORS | `ACCESS_CONTROL_ALLOW_ORIGIN` |
| Headers | `REF_ID_HEADER_KEY` |
| JWT | `JWT_ISSUER`, `JWT_AUDIENCE`, `JWT_EXP_DURATION`, `SECRET_JWT_PRIVATE_KEY` |
| Firestore | `GCP_PROJECT_ID`, `GCP_CREDENTIALS_JSON`, `GCP_FIRESTORE_DATABASE_ID`, `GCP_FIRESTORE_CONNECT_TIMEOUT` |
| Google OAuth | `GOOGLE_OAUTH2_VERIFY_TOKEN`, `GOOGLE_OAUTH2_GET_USER_PROFILE`, `GOOGLE_OAUTH2_REVOKE_TOKEN` |
| Secrets | `SECRET_AESGCM_KEY`, `SECRET_HASH_PEPPER` |
| Kafka consumer | `KAFKA_BROKERS`, `KAFKA_GROUP_ID`, `KAFKA_TOPICS`, `KAFKA_OFFSETS_INITIAL`, `KAFKA_REBALANCE_STRATEGY` |
| Kafka producer | `KAFKA_PRODUCER_BROKERS` |
| MySQL | `MYSQL_HOST`, `MYSQL_PORT`, `SECRET_DB_USERNAME`, `SECRET_DB_PASSWORD`, `MYSQL_DATABASE` |
| Redis | `REDIS_ADDR`, `SECRET_REDIS_PASSWORD` |
| HTTP Client | `HTTP_CLIENT_ENABLE_LOG_DEBUG` |

Notes:

- `GCP_CREDENTIALS_JSON` and `SECRET_JWT_PRIVATE_KEY` are base64-decoded at startup.
- Kafka subscriber startup is skipped when `KAFKA_BROKERS`, `KAFKA_GROUP_ID`, or `KAFKA_TOPICS` are empty.
- Kafka producer is skipped if `KAFKA_PRODUCER_BROKERS` is empty.

---

## Commands

The `Makefile` is the main entry point for development tasks.

| Command | Description |
|---|---|
| `make setup` | Install tools and project hooks |
| `make upgrade` | Upgrade dev tools (golangci-lint, govulncheck, swagger, pkgsite) |
| `make mod` | Run `go fmt` and `go mod tidy` |
| `make lint` | Run `golangci-lint` with auto-fix |
| `make test` | Run full test suite with race detection |
| `make test-integration` | Run integration tests (requires external services) |
| `make coverage` | Generate and open HTML coverage report |
| `make vuln` | Run `govulncheck` |
| `make precommit` | Run mod + lint + test + vuln + vet before committing |
| `make ci` | Full CI pipeline (precommit + diff check) |
| `make docker` | Build the Docker image |
| `make run` | Start the container with `.env` |
| `make swagger` | Regenerate OpenAPI spec |
| `make openapi` | Serve Swagger UI locally on port 8910 |
| `make doc` | Serve Go package docs locally on port 6060 |
| `make bump-version` | Bump version (usage: `make bump-version version=v0.0.1`) |
| `make clean` | Remove build artifacts and caches |
| `make release-cloudrun` | Build, push, and deploy to Cloud Run |

---

## HTTP APIs

The HTTP routes are registered in `router/router.go`.

| Area | Routes |
|---|---|
| Health | `GET /liveness`, `GET /readiness`, `GET /metrics` |
| Auth | `POST /api/v1/platform/auth/resolve-identity`, `POST /api/v1/platform/auth/issue-token` |
| Organization | `POST /api/v1/platform/organization/current`, `POST /api/v1/platform/organization/register` |
| Member | `POST /api/v1/platform/member/me`, `POST /api/v1/platform/member/register` |
| Product | `POST /api/v1/platform/product/create`, `POST /api/v1/platform/product/detail`, `POST /api/v1/platform/product/list`, `POST /api/v1/platform/product/update`, `POST /api/v1/platform/product/delete` |

Detailed request and response shapes live in `spec.md`.

---

## Kafka Events

Kafka event routing is registered in `router/subscriber.go`.

| Event | Domain | Handler File |
|---|---|---|
| `PRODUCT_CREATED` | product | `app/product/consumer_create.go` |
| `PRODUCT_UPDATED` | product | `app/product/consumer_update.go` |
| `PRODUCT_DELETED` | product | `app/product/consumer_delete.go` |
| `ORGANIZATION_CREATED` | organization | `app/organization/consumer_create.go` |
| `ORGANIZATION_MEMBER_ADDED` | organization | `app/organization/consumer_add_member.go` |
| `MEMBER_CREATED` | member | `app/member/consumer_create.go` |

If you add a new aggregate, keep its event consumer inside the corresponding `app/<domain>/` package and register its events in `router/subscriber.go`.

---

## How to Extend the Codebase

### Add a New Domain Aggregate

1. Create `app/<domain>/`.
2. Create `app/<domain>/access/storage_<dep>.go` with interface + impl + domain model.
3. (Optional) Create `access/cache_<dep>.go` for caching, `access/client_<dep>.go` for external APIs.
4. Add a `handler.go` defining `HandlerConfig` + `handler` struct + `NewHandler`.
5. Add `handler_<action>.go` files for HTTP endpoints.
6. Add `consumer_<action>.go` files for Kafka events.
7. Wire HTTP routes in `router/router.go` via `register<Domain>Routes(r, d)`.
8. Wire Kafka events in `router/subscriber.go` via `register<Domain>Events(d)`.
9. Write tests aiming for 100% coverage.

### Add a New HTTP Endpoint

1. Add or update a handler in `app/<domain>/handler_<action>.go`.
2. Wire the route in `router/router.go` inside the matching `register<Domain>Routes` function.
3. Add or update tests beside the handler.

### Add a New Kafka Event

1. Add a `consumer_<action>.go` file in the domain package (e.g. `app/product/consumer_create.go`).
2. Register the event name in `router/subscriber.go`.
3. Add tests in `consumer_<action>_test.go`.

### Add a New Infrastructure Client

Dependencies fall into two categories that determine how they flow into domain code:

**Raw infrastructure SDKs** (Firestore, MySQL, Redis, S3/GCS, `*http.Client`) — must be wrapped:
1. Add the SDK client field to `deps` struct in `router/deps.go`.
2. Initialize it in `newDeps()` and add its `Close()` to the cleanup function.
3. Create `access/storage_<dep>.go`, `access/cache_<dep>.go`, or `access/client_<dep>.go` in the domain package.
4. In `register<Domain>Routes` / `register<Domain>Events`, construct the access-layer impl and pass it via `HandlerConfig`.

**Common module abstractions** (`hash.HashManager`, `crypt.Cipher`, `token.JWTSigner`, `kafka.Producer`) — pass directly:
1. Add the interface field to `deps` struct (if not already there).
2. Initialize it in `newDeps()` with the corresponding `MustNew*` constructor.
3. Pass it directly to `HandlerConfig` — no `access/` wrapper needed (they are already interfaces).

### Dependency Injection Pattern (`router/deps.go`)

All infrastructure clients are initialized once in `router/deps.go` and passed to handlers via the `deps` struct. This is the **composition root** — it keeps `main.go` thin and co-locates cleanup with initialization.

#### Step 1 — Register the client in `deps`

```go
// router/deps.go
type deps struct {
    cfg             config.Config

    // Raw SDK handles — wrapped in access/ before reaching domain code
    httpClient      *http.Client
    firestoreClient *commonfirestore.Client
    mysqlClient     *sql.DB
    redisClient     redis.UniversalClient

    // Common module abstractions — already interfaces, passed directly to HandlerConfig
    hash            hash.HashManager
    cipher          crypt.Cipher
    token           token.JWTSigner
    producer        kafka.Producer
}
```

#### Step 2 — Initialize and clean up in `newDeps`

```go
func newDeps(ctx context.Context, cfg config.Config) (deps, func()) {
    httpClient := httpclient.NewHTTPClient(...)
    fs := commonfirestore.MustNewClient(ctx, newFirestoreConfig(cfg))
    db := database.MustNewMySQLWithConfig(newMySQLConfig(cfg))
    rdb := commonredis.MustNew(cfg.Redis.Addr, cfg.Redis.Password)
    producer := newProducer(cfg)

    d := deps{
        cfg:             cfg,
        httpClient:      httpClient,
        firestoreClient: fs,
        mysqlClient:     db,
        redisClient:     rdb,
        hash:            newHashManager(cfg),
        cipher:          newCipher(cfg),
        token:           newTokenManager(cfg),
        producer:        producer,
    }

    cleanup := func() {
        _ = fs.Close()
        _ = db.Close()
        _ = rdb.Close()
        if producer != nil {
            _ = producer.Close()
        }
    }

    return d, cleanup
}
```

#### Step 3 — Pass to domain handlers via `router.go`

```go
// router/router.go
func registerProductRoutes(r *gin.Engine, d deps) {
    // Raw SDK → wrap in access/ layer (domain never sees the SDK directly)
    productStorage := productaccess.NewProductStorage(d.firestoreClient.Inner())

    h := product.NewHandler(product.HandlerConfig{
        // access/ interfaces — domain-specific wrappers
        ProductStorage: productStorage,
        // common module interfaces can be passed directly:
        // Hash: d.hash,
        // Cipher: d.cipher,
    })

    productGroup := r.Group("/api/v1/platform/product")
    {
        productGroup.POST("/create", h.CreateProduct)
    }
}
```

---

## Maintenance Notes

- Keep `main.go` thin: initialization and shutdown only.
- Keep domain logic in `app/` rather than in `router/` or helper packages.
- Prefer small files with responsibility-based names over large catch-all files.
- Use `spec.md` when changing request or response payloads.
- Update tests next to the code they validate.
- No cross-domain imports (`app/product` never imports `app/member`).
- Use compile-time interface checks: `var _ Interface = (*impl)(nil)`.
- File naming uses prefix convention: `handler_`, `consumer_`, `service_`, `storage_`, `cache_`, `client_`.
- Co-locate domain models and errors with the access file that uses them — no separate `model.go` or `errors.go` in `access/`.
- Raw SDK handles must be wrapped in `access/`; common module interfaces pass directly to `HandlerConfig`.

### Refactoring Guidelines
When addressing tech debt or making structural improvements, follow **Martin Fowler's refactoring principles**:
- **Layer Mapping**: Handlers strictly perform transport concerns. Service logic coordinates flow without business rules. Domain models contain all domain/business logic. Adapters manage pure I/O.
- **Workflow**: Identify smells (e.g., God handlers, primitive obsession, conditional explosions), make small isolated changes, verify with tests, and repeat. Do not rewrite wholesale.
- **Tell, Don't Ask**: Move data manipulation into the domain models representing the data; don't pull data out to manipulate it externally.

---

## Troubleshooting

- If startup fails early, check base64-encoded secrets (`SECRET_JWT_PRIVATE_KEY`, `GCP_CREDENTIALS_JSON`) first.
- If HTTP starts but Kafka does not, verify `KAFKA_BROKERS`, `KAFKA_GROUP_ID`, and `KAFKA_TOPICS` are present and non-empty.
- If Firestore access fails, check `GCP_PROJECT_ID`, `GCP_CREDENTIALS_JSON`, and `GCP_FIRESTORE_DATABASE_ID`.
- If MySQL connection fails, verify `MYSQL_HOST`, `MYSQL_PORT`, `SECRET_DB_USERNAME`, `SECRET_DB_PASSWORD`, and `MYSQL_DATABASE`.
- If Redis connection fails, verify `REDIS_ADDR` and `SECRET_REDIS_PASSWORD`.

---

## Repository Purpose

This repository serves as a **foundational template and reference implementation** for new A-Team microservices. The example domains (auth, member, organization, product) demonstrate the structural conventions — they are blueprints, not production code.

When creating a new service:

- **Keep**: `main.go` structure, `router/` wiring pattern, `config/` pattern, `Makefile`, `Dockerfile`, `.golangci.yaml`, `.mockery.yaml`, middleware chain, all tooling scripts.
- **Replace**: `app/` domains with your own business logic, `config/config.go` fields for your service, route and event registrations.
- **Import**: All infrastructure from the `common` module — never copy.
