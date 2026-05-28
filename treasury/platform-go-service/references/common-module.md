# Common Module — Quick Reference

The shared module at `gitdev.devops.krungthai.com/ktb-digital-lending/dgl-new/shareable/common` provides every infrastructure abstraction. Services **import** it; never vendor or copy.

---

## Packages

| Package | Purpose | Key types / functions |
|---|---|---|
| `logger` | Structured `slog`-based logging | `New()`, `AWSKeyReplacer`, `OpenSearchKeyReplacer`, `CensorReplacer`, `AddCensor`, `LoadCensorsFromMap` |
| `serror` | Source-location error wrapping | `Wrap(err)`, `New(msg)`, `.With(slog.Attr...)` |
| `wrapper` | HTTP response envelope | `Respond()`, `BindJSON[T]()`, `ResponseOption[T]` |
| `middleware` | Gin middleware chain | `SecurityHeaders`, `AccessControl`, `RefIDMiddleware`, `TraceContextTraceIDMiddleware`, `AutoLoggingMiddleware`, `Timeout`, `AccessLog` |
| `kafka` | Producer / Consumer / Event Router | `MustNewProducer`, `MustNewConsumerGroup`, `NewEventRouter`, `BindMessage`, `Message[T]`, `KafkaHandler`, `WithLogging`, `WithLoggingFromEnv` |
| `firestore` | Firestore client wrapper | `MustNewClient`, `Client.Inner()` |
| `database` | MySQL connection helper | `MustNewMySQLWithConfig` |
| `redis` | Redis client factory | `MustNew` |
| `httpclient` | HTTP client with middleware | `NewHTTPClient`, `DebugOption` |
| `crypt` | AES-GCM encryption | `MustNew`, `Cipher.Encrypt`, `Cipher.Decrypt` |
| `hash` | bcrypt hashing with pepper | `MustNewHashManager`, `HashManager` |
| `token` | JWT signing (ES256) | `MustNewJWTSigner`, `JWTSigner` |
| `codec` | Base64 encoding/decoding | `NewBase64Coder`, `DecodeBase64` |
| `config` | Env parsing + environment helpers | `ParseEnv[T]`, `IsLocalEnv`, `IsProdEnv`, `Env` |
| `health` | Liveness / readiness / metrics | `Liveness`, `Readiness`, `Metrics` |
| `app` | Shared constants + timestamp helpers | `CodeSuccess`, `CodeBadRequest`, `CodeInternalError`, `MessageSuccess`, `SetTimestamps` |

---

## Top external dependencies (third-party)

| Package | Purpose |
|---|---|
| `github.com/gin-gonic/gin` | HTTP framework |
| `github.com/IBM/sarama` | Kafka client |
| `cloud.google.com/go/firestore` | Firestore client |
| `github.com/redis/go-redis/v9` | Redis client |
| `github.com/caarlos0/env/v11` | Env config parsing |
| `github.com/google/uuid` | UUID generation |
| `github.com/stretchr/testify` | Test assertions and mocks |
| `github.com/vektra/mockery/v2` | Mock generation |
