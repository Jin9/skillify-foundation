# Package Reference

Authoritative table of every public package in `common`. Read when you need to find which package owns a capability or look up a key API surface.

## HTTP & Web

| Package | Key API | Purpose |
|---|---|---|
| `wrapper` | `Respond[T]`, `BindJSON[T]`, `Response[T]`, `ResponseOption[T]`, `Code`, `Message` | Standard JSON response envelope and request binding |
| `health` | `Liveness(version, commit)`, `Readiness()`, `Metrics()` | Health check / readiness / runtime metrics endpoints |
| `middleware` | `AccessLog()`, `SecurityHeaders()`, `AccessControl(origin, headers)`, `RefIDMiddleware(header)`, `TraceContextTraceIDMiddleware(header)`, `JWT(parser, verifier)`, `AutoLoggingMiddleware(successCode)`, `Timeout(d)` | Gin middleware stack |

## Messaging

| Package | Key API | Purpose |
|---|---|---|
| `kafka` | `NewProducer`, `NewConsumerGroup`, `NewEventRouter`, `NewMessage[T]`, `BindMessage[T]`, `WithLogging`, `WithLoggingFromEnv`, `SendMessageOption`, `KafkaHandler`, `Processor` | Kafka producer, consumer group, event routing, message envelope, logging interceptor, SSL/SASL security |

## Data Stores

| Package | Key API | Purpose |
|---|---|---|
| `database` | `ConnectPostgresDB(cfg)`, `ConnectPostgresDBWithContext(ctx, cfg)`, `ParseConfig(cfg)`, `MustNewPostgresDB`, `PostgresConfig` | PostgreSQL (pgx pool) connection helpers |
| `database` | `NewMySqlDB(cfg)`, `MustNewMySqlDB`, `MySqlConfig` | MySQL connection helpers |
| `firestore` | `New(ctx, cfg)`, `NewClient(ctx, cfg)`, `MustNewClient`, `Config{ProjectID, EmulatorHost, DatabaseID, …}` | GCP Firestore client wrapper with emulator support |
| `redis` | `Connect`, `ConnectWithConfig`, `ConnectClusterWithConfig`, `ConnectFailOverWithConfig` | Redis universal, cluster, and failover clients |

## Cloud & File

| Package | Key API | Purpose |
|---|---|---|
| `cloud-storage` | `New(ctx, cfg)`, `NewClient(ctx, cfg)`, `Config{ProjectID, EmulatorHost, …}` | GCS client wrapper with emulator support |
| `cloud-storage` | `NewS3Client(cfg)`, `S3Config` | S3-compatible storage (AWS SDK v2) |
| `sftp` | `NewClient(cfg)`, `Client{Upload, Download, ListDirectory, Delete, Close}` | SFTP file transfer (key or password auth) |
| `zip` | `NewZipper()`, `Zipper{CompressStream, UnzipSecure}` | ZIP compression/decompression with Zip Slip protection |

## Security & Cryptography

| Package | Key API | Purpose |
|---|---|---|
| `token` | `NewJWTSigner(cfg)`, `NewJWTVerifier(cfg)`, `NewJWTParser(cfg)`, `Claims{Sub, Iss, Aud, Exp, Extra}`, `WithClaims(ctx, claims)`, `ClaimsFromContext(ctx)` | JWT manager (RS256, ES256, HS512) with context-propagated claims |
| `crypt` | `New(cfg)`, `Cipher{GenerateIV, Encrypt, Decrypt}` | AES-256-GCM encryption |
| `crypt` | `NewRSA(cfg)`, `RSA{Encrypt, Decrypt, Sign, Verify}` | RSA encryption and signing |
| `hash` | `NewHashManager(cfg)`, `HashManager{GenerateSalt, HashSha256Encode, HashSha256EncodeSalt, HashSha256EncodePepper}` | SHA-256 hashing with salt and pepper |

## Configuration

| Package | Key API | Purpose |
|---|---|---|
| `config` | `ParseEnv[T](opts)`, `IsLocalEnv()`, `IsDevEnv()`, `IsUATEnv()`, `IsProdEnv()`, `Env` | Environment variable parsing with prefix support |
| `validator` | `Validate(i any) error`, `MustValid(i any)` (deprecated) | Struct validation via `go-playground/validator` |

## Utilities

| Package | Key API | Purpose |
|---|---|---|
| `logger` | `New(replacers…)`, `AWSKeyReplacer`, `GCPKeyReplacer`, `OpenSearchKeyReplacer`, `CensorReplacer`, `AddCensor(key, mask)`, `LoadCensorsFromMap(m)` | Structured `slog` logging with cloud-specific key mapping and sensitive-data masking |
| `httpclient` | `NewHTTPClient(opts…)`, `NewHTTPClientWithCA(pool, opts…)`, `Get[RES]`, `Post[REQ,RES]`, `Put[REQ,RES]`, `Patch[REQ,RES]`, `Delete[RES]`, `ForwardRefIDOption` | HTTP client with connection pooling, ref-id forwarding, debug logging |
| `generator` | `NewAppUUID()`, `GenerateSecureString(length, charset)`, `NewTraceParent()`, `Parse(traceparent)`, `TraceParent.CreateChild()` | UUID v4/v7, secure random strings, W3C Trace Context |
| `codec` | `NewJSONCoder()`, `NewBase64Coder()` | JSON and Base64 encoding/decoding behind interfaces |
| `clock` | `New()`, `NowUTC()`, `NowBangkok()`, `ParseUTC(value)`, `ParseBangkok(value)`, `ToBangkok(t)` | Mockable time interface (Fowler: Clock Wrapper) |
| `safe` | `Deref[T](p)`, `DerefOr[T](p, fallback)`, `Ptr[T](v)`, `IsNil[T](p)` | Null-safe pointer helpers |
| `serror` | `New(s)`, `Wrap(err)`, `WrapSkip(err, skip)`, `SError.With(attrs…)`, `SError.Attrs()`, `DecodeMessage(s)` | Structured errors with file/line/func context and slog attrs |
| `app` | `Code`, `Message`, `Response[T]`, `SetTimestamps(createdAt, updatedAt)` | Domain response code re-exports from `wrapper`; timestamp helpers |
