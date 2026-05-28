# Common — Shared Go Infrastructure Library

Reusable infrastructure and utility packages for A-Team platform Go microservices.

```bash
go get gitlab.com/b2c-e-commerce-platform/platform/backend/common
```

> **Library contract** — this module contains zero domain business logic.
> It is the *Shared Kernel* (Evans, DDD) that every microservice imports.
> Domain-specific codes, models, and handlers live in each service's own repository.

---

## Architectural Foundations

| Pillar | How it manifests here |
|---|---|
| **DDD — Shared Kernel** | Generic infrastructure only; consuming services own their Aggregates and Bounded Contexts |
| **CQRS** | `wrapper.Respond` (command/query response envelope), `kafka.Message[T]` (async command channel), `kafka.NewEventRouter` (event dispatch) |
| **Fowler — Gateway** | `httpclient`, `database`, `redis`, `firestore`, `cloud-storage`, `sftp` |
| **Fowler — Clock Wrapper** | `clock.Clock` interface for deterministic testing |
| **Fowler — Special Case** | `safe.Deref`, `safe.DerefOr`, `safe.Ptr` |
| **Decorator / Interceptor** | `kafka.WithLogging(producer, …)` — logging decorator for Kafka Producer |
| **Separated Interface** | Every package exports an interface + concrete + `mocks/` |
| **Go idiom** | Accept interfaces, return structs · `Must*` / non-`Must*` pairs · `var _ I = (*s)(nil)` · `context.Context` first · Generics (`[T any]`) |

---

## Packages

**HTTP & Web**

```text
wrapper/               Standardised JSON response envelope (Respond[T], BindJSON[T],
                       Response[T], ResponseOption[T], Code, Message)
health/                Liveness, readiness, and runtime-metrics endpoints
middleware/            AccessLog, SecurityHeaders, AccessControl, RefIDMiddleware,
                       TraceContextTraceIDMiddleware, JWT, AutoLoggingMiddleware,
                       Timeout, response-writer interceptor
```

**Messaging**

```text
kafka/                 Sarama-based Kafka producer & consumer group, event-routing
                       processor (NewEventRouter), Message[T] envelope, BindMessage[T],
                       logging interceptor (WithLogging / WithLoggingFromEnv),
                       SSL/TLS + SASL/SCRAM security, backoff retry
```

**Data Stores**

```text
database/              PostgreSQL (pgx pool) & MySQL connection helpers with
                       sane pool defaults (lifetime jitter, idle timeout, health check)
firestore/             GCP Firestore client wrapper with emulator & named-database support
redis/                 Redis universal, cluster, and sentinel-failover clients
```

**Cloud & File**

```text
cloud-storage/         GCS and S3 (AWS SDK v2) compatible storage clients
sftp/                  SFTP file transfer (key or password auth, context-aware I/O)
zip/                   ZIP compression and decompression (Zip Slip protected)
```

**Security & Cryptography**

```text
token/                 JWT manager — signing (RS256, ES256), verification (RS256,
                       ES256, HS512), parsing, claims validation, and context
                       propagation (WithClaims / ClaimsFromContext).
                       Claims.Extra supports service-specific flat claims (RFC 7519).
crypt/                 AES-256-GCM encryption and RSA encrypt/decrypt/sign/verify
hash/                  SHA-256 hashing with salt and pepper
```

**Configuration**

```text
config/                Environment-variable parsing (ParseEnv[T]) with prefix support.
                       Environment helpers: IsLocalEnv, IsDevEnv, IsUATEnv, IsProdEnv.
validator/             Struct validation via go-playground/validator (Validate returns error)
```

**Utilities**

```text
logger/                Structured slog logging with pluggable ReplacerFuncs
                       (AWSKeyReplacer, GCPKeyReplacer, OpenSearchKeyReplacer,
                       CensorReplacer). Runtime censor management (AddCensor,
                       LoadCensorsFromMap). ENV=local → text handler; else → JSON.
httpclient/            HTTP client with connection pooling, ref-id forwarding,
                       per-request debug logging, TLS/CA support. Generic helpers:
                       Get[RES], Post[REQ,RES], Put, Patch, Delete, WithOptions variants.
generator/             UUID v4/v7, secure random string, W3C traceparent
                       (generate, parse, create child span)
codec/                 JSON and Base64 encoding/decoding behind interfaces
clock/                 Mockable time interface (Now, NowIn, Parse, ParseIn) plus
                       convenience functions (NowUTC, NowBangkok, ToBangkok)
safe/                  Null-safe pointer helpers: Deref[T], DerefOr[T], Ptr[T], IsNil[T]
serror/                Structured error wrapping with file/line/func context and
                       slog.Attr investigation fields. Flows to AccessLog middleware.
```

**Application**

```text
app/                   Domain response codes re-exported from wrapper (Code, Message,
                       Response[T] type aliases). Timestamp helper (SetTimestamps).
```

---

## Usage Examples

### `wrapper` — Standard HTTP Response

```go
import "gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"

// Success
wrapper.Respond(c, wrapper.ResponseOption[*MyData]{
    HTTPStatus: http.StatusOK,
    Code:       wrapper.CodeSuccess,
    Message:    wrapper.MessageSuccess,
    Data:       &result,
})

// Error (Err flows to AccessLog middleware → structured log)
wrapper.Respond(c, wrapper.ResponseOption[any]{
    HTTPStatus: http.StatusInternalServerError,
    Code:       wrapper.CodeInternalError,
    Message:    wrapper.MessageInternalError,
    Err:        serror.Wrap(err),
})
```

### `wrapper` — Request Binding

```go
req, ok := wrapper.BindJSON[CreateRequest](c,
    slog.String("handler", "create"),
)
if !ok {
    return // 400 already written
}
```

### `logger` — Structured Logging

```go
import "gitlab.com/b2c-e-commerce-platform/platform/backend/common/logger"

_ = logger.New(
    logger.AWSKeyReplacer,        // msg→message, time→timestamp
    logger.OpenSearchKeyReplacer, // msg→message, time→@timestamp, level→log.level
    logger.CensorReplacer,        // masks password, access_token, api_key, secret
)

logger.AddCensor("ssn", "***SSN***")
```

### `serror` — Structured Error Wrapping

```go
import "gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"

if err := db.Insert(ctx, record); err != nil {
    return serror.Wrap(err)
}

// Attach investigation context
return serror.Wrap(err).With(
    slog.String("user_id", userID),
    slog.String("action", actionName),
)
```

### `config` — Environment Parsing

```go
import "gitlab.com/b2c-e-commerce-platform/platform/backend/common/config"

type Config struct {
    APIKey string `env:"API_KEY,required"`
    Port   int    `env:"PORT" envDefault:"8080"`
}

cfg, err := config.ParseEnv[Config](env.Options{Prefix: "MYAPP_"})
if config.IsLocalEnv() { /* … */ }
```

### `kafka` — Producer & Logging Interceptor

```go
import "gitlab.com/b2c-e-commerce-platform/platform/backend/common/kafka"

producer, _ := kafka.NewProducer(kafka.ProducerConfig{
    KafkaConf: kafka.NewSyncProducerGuarantee(),
    Brokers:   []string{"broker-1:9092"},
})

// Decorator: auto-select log mode from ENV
logged := kafka.WithLoggingFromEnv(producer, config.Env)

msg := kafka.NewMessage("USER_CREATED", userID, payload)
_ = logged.SendMessageWithOption(ctx, "topic", msg, kafka.SendMessageOption{
    LogAttrs: []slog.Attr{slog.String("user_id", userID)},
})
```

### `kafka` — Consumer & Event Router

```go
group, _ := kafka.NewConsumerGroup(kafka.ConsumerConfig{
    KafkaConf: kafka.NewConsumerConfigAtLeastOnce(),
    Brokers:   []string{"broker-1:9092"},
    GroupID:   "my-service",
})

handlers := map[string]kafka.KafkaHandler{
    "USER_CREATED": onUserCreated,
    "USER_UPDATED": onUserUpdated,
}

processor := kafka.NewEventRouter(handlers)
group.Consume(ctx, []string{"events-topic"},
    kafka.NewConsumerGroupHandler(ctx, processor))
```

### `middleware` — Gin Stack

```go
import "gitlab.com/b2c-e-commerce-platform/platform/backend/common/middleware"

r := gin.New()
r.Use(middleware.AccessLog())
r.Use(middleware.SecurityHeaders())
r.Use(middleware.RefIDMiddleware("X-Request-ID"))
r.Use(middleware.TraceContextTraceIDMiddleware("traceparent"))
r.Use(middleware.JWT(jwtParser, jwtVerifier))
```

### `token` — JWT Lifecycle

```go
import "gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"

signer, _ := token.NewJWTSigner(token.JWTSignerConfig{
    PrivateKey: rsaPEM, Alg: "RS256",
    Issuer: "svc", Audience: "partner", Expire: 15 * time.Minute,
})
jwt, _ := signer.SignRS256(token.Claims{Sub: userID, Extra: map[string]any{"role": "admin"}})

// Verify & parse
verifier, _ := token.NewJWTVerifier(token.JWTVerifierConfig{PublicKey: pubPEM, Alg: "RS256"})
parser, _ := token.NewJWTParser(token.JWTParserConfig{Issuer: "svc", Audience: "partner"})
_ = verifier.VerifySignatureRS256(jwt)
claims, _ := parser.ParseToken(jwt)

// Context propagation
ctx := token.WithClaims(ctx, claims)
c, _ := token.ClaimsFromContext(ctx)
```

### `clock` — Mockable Time (Fowler: Clock Wrapper)

```go
import "gitlab.com/b2c-e-commerce-platform/platform/backend/common/clock"

clk := clock.New()
now := clk.Now()                        // UTC
bkk, _ := clk.NowIn(clock.Bangkok)     // Asia/Bangkok
t, _ := clk.Parse(clock.DefaultLayout, "2025-01-15T10:30:00")
```

---

## Project Structure

```text
.
├── app/                 Domain response codes + timestamp helpers
├── clock/               Mockable time interface
├── cloud-storage/       GCS and S3 storage clients
├── codec/               JSON and Base64 codecs
├── config/              Environment variable parsing
├── crypt/               AES-GCM and RSA encryption
├── database/            PostgreSQL (pgx) and MySQL connection helpers
├── firestore/           GCP Firestore client wrapper
├── generator/           UUID, secure random, traceparent generators
├── hash/                SHA-256 hashing with salt/pepper
├── health/              Health check endpoints
├── httpclient/          HTTP client with pooling and ref-id forwarding
├── kafka/               Producer, consumer, event router, interceptor, security
├── logger/              Structured slog logging with replacers and censoring
├── middleware/          Gin middleware (auth, observability, security, tracing)
├── redis/               Redis universal, cluster, failover clients
├── safe/                Null-safe pointer helpers
├── serror/              Structured error wrapping with source location
├── sftp/                SFTP file transfer client
├── token/               JWT signing, verification, parsing, context claims
├── validator/           Struct validation (go-playground/validator)
├── wrapper/             Standard JSON response envelope
├── zip/                 ZIP compression (Zip Slip protected)
├── .golangci.yaml       Linter configuration
├── .mockery.yaml        Mock generation configuration
├── CHANGELOG.md         Release history
├── Makefile             Development targets
├── README.md            This file
└── VERSION              Current version tag
```

## Prerequisites

- Go 1.25+
- Access to the private module host `gitlab.com/b2c-e-commerce-platform/platform/backend`

```bash
go env -w GOPRIVATE=gitlab.com/b2c-e-commerce-platform/platform/backend
```

## Installation

```go
import (
    "gitlab.com/b2c-e-commerce-platform/platform/backend/common/logger"
    "gitlab.com/b2c-e-commerce-platform/platform/backend/common/middleware"
    "gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)
```

```bash
go mod tidy
```

## Local Development

1. Run `make setup` to install dev tools and git hooks.
2. Run `make test` to confirm everything passes.
3. Run `make doc` to browse package documentation locally on `:6060`.

## Commands

| Command | Description |
|---|---|
| `make setup` | Install dev tools and git hooks |
| `make mod` | Run `go fmt` and `go mod tidy` |
| `make lint` | Run `golangci-lint` with auto-fix |
| `make test` | Run tests with race detection and coverage |
| `make coverage` | Generate HTML coverage report |
| `make vuln` | Run `govulncheck` |
| `make precommit` | Full validation: lint + test + vuln + vet + mod verify |
| `make ci` | CI pipeline: precommit + uncommitted-changes check |
| `make doc` | Serve package docs locally via pkgsite |
| `make upgrade` | Upgrade dev tools (golangci-lint, govulncheck, etc.) |
| `make bump-version version=vX.Y.Z` | Tag a new release version |

## How to Add a New Package

1. Create a new directory at the repository root (e.g. `newpkg/`).
2. Define an exported interface and at least one concrete implementation.
3. Add `var _ Interface = (*impl)(nil)` compile-time interface check.
4. Add a `mocks/` subdirectory with mockery-generated mocks.
5. Write tests alongside the code — table-driven, aim for 100% branch coverage.
6. Run `make precommit` to validate before pushing.

**Rules for this library:**
- No `main.go` — this is a library, not an application.
- No domain business logic — keep packages infrastructure-only.
- Use `slog` for logging (never `log` or `fmt.Println`).
- Return errors, don't panic (`Must*` constructors are deprecated for new code).
- Configuration via `env` struct tags, parsed by `config.ParseEnv`.
- `context.Context` as first parameter for all I/O-facing functions.
- Wrap errors with `serror.Wrap(err)` for source-location tracking.

## Versioning

This project follows [Semantic Versioning](https://semver.org/). All notable changes are recorded in `CHANGELOG.md`.

```bash
make bump-version version=v1.2.3
```

The current version is tracked in `VERSION`.

## Maintenance Notes

- Keep each package focused on a single responsibility.
- Provide a `mocks/` subdirectory for every exported interface.
- Prefer small files with clear, responsibility-based names.
- Update tests next to the code they validate.
