# Recipes

Concrete code patterns for the most common usages of `common`. Read this when you need a working snippet for HTTP responses, request binding, logging, error wrapping, middleware wiring, Kafka, config, clock, or JWT.

## 1. Response Envelope

All HTTP responses go through `wrapper.Respond`. It writes JSON and attaches a `ResponseMeta` to the Gin context (consumed by `AccessLog` middleware).

```go
// Success
wrapper.Respond(c, wrapper.ResponseOption[*MyData]{
    HTTPStatus: http.StatusOK,
    Code:       wrapper.CodeSuccess,
    Message:    wrapper.MessageSuccess,
    Data:       &result,
})

// Error (Err flows to AccessLog → structured log with source location)
wrapper.Respond(c, wrapper.ResponseOption[any]{
    HTTPStatus: http.StatusInternalServerError,
    Code:       wrapper.CodeInternalError,
    Message:    wrapper.MessageInternalError,
    Err:        serror.Wrap(err).With(slog.String("user_id", id)),
})
```

Wire shape:

```json
{ "code": "0000", "message": "Success", "data": { "...": "..." } }
```

Standard codes live in `wrapper/codes.go`. Domain codes live in the consuming service's `app/codes.go` (re-exported via type aliases from `wrapper`).

## 2. Request Binding

```go
// Returns (*T, true) on success, writes 400 + returns (nil, false) on failure.
req, ok := wrapper.BindJSON[CreateProductRequest](c,
    slog.String("handler", "create_product"),
)
if !ok {
    return // 400 already written
}
```

## 3. Logger with Replacers and Censoring

```go
_ = logger.New(
    logger.AWSKeyReplacer,        // msg→message, time→timestamp
    logger.GCPKeyReplacer,        // msg→message, time→timestamp, level→severity
    logger.OpenSearchKeyReplacer, // msg→message, time→@timestamp, level→log.level
    logger.CensorReplacer,        // masks password, access_token, api_key, secret
)

// Add runtime censors:
logger.AddCensor("ssn", "***SSN***")
logger.LoadCensorsFromMap(map[string]string{"credit_card": "***CC***"})
```

- `ENV=local` → text handler; all others → JSON handler.
- `LOG_LEVEL` env var controls level (DEBUG, INFO, WARN, ERROR). Default: ERROR.

## 4. Structured Error Wrapping

```go
// Wrap with source location
if err := db.Insert(ctx, record); err != nil {
    return serror.Wrap(err)
}

// Attach investigation context (auto-flows to AccessLog middleware)
return serror.Wrap(err).With(
    slog.String("user_id", userID),
    slog.String("action", "create_product"),
)

// Skip caller frames (when wrapping inside a helper)
return serror.WrapSkip(err, 1)
```

The `AccessLog` middleware decodes SError messages and emits structured logs with `error_source{file, line, func}` and `error_context{...attrs}` groups.

## 5. Middleware Stack (Gin)

Typical ordering:

```go
r := gin.New()
r.Use(middleware.AccessLog())                                  // access log + error logging
r.Use(middleware.SecurityHeaders())                            // OWASP security headers
r.Use(middleware.RefIDMiddleware("X-Request-ID"))              // extract/generate ref-id
r.Use(middleware.TraceContextTraceIDMiddleware("traceparent")) // W3C Trace Context
r.Use(middleware.AutoLoggingMiddleware(wrapper.CodeSuccess))   // non-OK response auto-log
r.Use(middleware.Timeout(30 * time.Second))                    // request timeout
r.Use(middleware.AccessControl("*", []string{"Authorization", "Content-Type"}))

// Protected routes
protected := r.Group("/v1")
protected.Use(middleware.JWT(jwtParser, jwtVerifier))
```

`JWT` middleware verifies ES256 signature, parses claims, validates iss/aud/exp/iat, and injects claims into `context.Context` via `token.WithClaims`. Retrieve with `token.ClaimsFromContext(ctx)`.

## 6. Kafka Producer + Interceptor (Decorator Pattern)

```go
producer, _ := kafka.NewProducer(kafka.ProducerConfig{
    KafkaConf: kafka.NewSyncProducerGuarantee(),
    Brokers:   []string{"broker-1:9092"},
})

// Decorator: auto-selects log mode from ENV (LOCAL/DEV→Debug, UAT/PROD→Meta)
logged := kafka.WithLoggingFromEnv(producer, config.Env)

msg := kafka.NewMessage("USER_CREATED", userID, payload)
err := logged.SendMessageWithOption(ctx, "topic", msg, kafka.SendMessageOption{
    WithRetry:   true,
    MaxRetries:  3,
    FailedTopic: "topic.failed",
    Key:         userID,
    LogAttrs:    []slog.Attr{slog.String("user_id", userID)},
})
```

Log modes:

- `LogModeDebug` — full payload + envelope + custom attrs (LOCAL/DEV)
- `LogModeMeta` — envelope + custom attrs only, no payload (UAT/PROD)
- `LogModeSilent` — no produce logs (high-frequency handlers)

## 7. Kafka Consumer + Event Router

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

Each `KafkaHandler` receives `kafka.Message[json.RawMessage]`. Use `kafka.BindMessage[T]` to unmarshal and validate:

```go
func onUserCreated(ctx context.Context, msg kafka.Message[json.RawMessage]) error {
    var payload UserCreatedEvent
    if err := kafka.BindMessage(msg.Payload, &payload); err != nil {
        return serror.Wrap(err)
    }
    // ... domain logic ...
    return nil
}
```

## 8. Config Parsing

```go
type Config struct {
    APIKey string `env:"API_KEY,required"`
    Port   int    `env:"PORT" envDefault:"8080"`
}

cfg, err := config.ParseEnv[Config](env.Options{Prefix: "MYAPP_"})
if config.IsLocalEnv() { /* ... */ }
```

## 9. Clock Wrapper (Fowler)

```go
// Production: real system clock
clk := clock.New()
now := clk.Now()         // UTC
bkk, _ := clk.NowIn(clock.Bangkok)

// Testing: inject mock
type MockClock struct { clock.Clock }
func (m *MockClock) Now() time.Time { return fixedTime }
```

## 10. JWT Token Lifecycle

```go
// Signing (service-to-service)
signer, _ := token.NewJWTSigner(token.JWTSignerConfig{
    PrivateKey: rsaPEM,
    Alg:        string(token.RS256),
    Issuer:     "my-service",
    Audience:   "partner-service",
    Expire:     15 * time.Minute,
})
jwt, _ := signer.SignRS256(token.Claims{Sub: userID, Extra: map[string]any{"role": "admin"}})

// Verification + parsing
verifier, _ := token.NewJWTVerifier(token.JWTVerifierConfig{PublicKey: pubPEM, Alg: "RS256"})
parser, _ := token.NewJWTParser(token.JWTParserConfig{Issuer: "my-service", Audience: "partner-service"})

_ = verifier.VerifySignatureRS256(jwt)
claims, _ := parser.ParseToken(jwt)
_ = parser.ValidateClaims(claims)

// Context propagation (via middleware)
ctx := token.WithClaims(ctx, claims)
c, _ := token.ClaimsFromContext(ctx)
role := c.Extra["role"]  // "admin"
```
