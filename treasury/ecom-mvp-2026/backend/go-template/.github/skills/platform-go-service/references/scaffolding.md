# Scaffolding — `main.go`, `config/`, `router/`

Templates for the bootstrap and composition-root layer. Pair with `domain-patterns.md` when adding actual business code.

---

## 1. Initialise the module

```bash
mkdir <service-name> && cd <service-name>
go mod init gitdev.devops.krungthai.com/a-team/<team>/<service-name>
```

Add the `common` module via a `replace` directive in `go.mod`:

```text
replace gitdev.devops.krungthai.com/ktb-digital-lending/dgl-new/shareable/common => gitdev.devops.krungthai.com/ktb-digital-lending/dgl-new/shareable/common.git v0.0.14
```

---

## 2. `main.go` — Bootstrap

Key conventions: `signal.NotifyContext` (not manual `signal.Notify`), embedded `VERSION`, explicit timeout constants, `router.StartSubscriber` **before** HTTP server, deferred teardown order.

```go
package main

import (
    "context"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "runtime"
    "runtime/debug"
    "syscall"
    "time"

    commonconfig "gitdev.devops.krungthai.com/ktb-digital-lending/dgl-new/shareable/common/config"
    "gitdev.devops.krungthai.com/ktb-digital-lending/dgl-new/shareable/common/logger"

    "<module>/config"
    "<module>/router"

    _ "embed"
    _ "time/tzdata"
)

const (
    gracefulShutdownDuration = 10 * time.Second
    serverReadHeaderTimeout  = 5 * time.Second
    serverReadTimeout        = 5 * time.Second
    serverWriteTimeout       = 10 * time.Second
    handlerTimeout           = serverWriteTimeout - (time.Millisecond * 100)
)

var commit string // go build -ldflags "-X main.commit=abc123"

//go:embed VERSION
var version string

func init() {
    if os.Getenv("GOMAXPROCS") != "" {
        runtime.GOMAXPROCS(0)
    } else {
        runtime.GOMAXPROCS(1)
    }
    if os.Getenv("GOMEMLIMIT") != "" {
        debug.SetMemoryLimit(-1)
    }
}

func main() {
    cfg := config.C(commonconfig.Env)
    _ = logger.New(
        logger.AWSKeyReplacer,
        logger.OpenSearchKeyReplacer,
        logger.CensorReplacer,
    )

    ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer cancel()

    // 1. Start Kafka subscriber (before HTTP so events are consumed immediately)
    consumerDone, stopConsumer := router.StartSubscriber(ctx, cfg)

    defer func() {
        cancel()
        stopConsumer()
        router.WaitForSubscriber(consumerDone, gracefulShutdownDuration)
    }()

    // Monitor unexpected subscriber exit → trigger global shutdown
    go func() {
        <-consumerDone
        cancel()
    }()

    // 2. Build router (creates deps internally, returns cleanup)
    r, cleanupRouter := router.New(cfg, version, commit, handlerTimeout)
    defer cleanupRouter()

    // 3. Start HTTP server
    srv := newServer(cfg, r)
    go gracefulShutdown(ctx, srv, gracefulShutdownDuration)

    slog.Info("run", "port", cfg.Server.Port)

    defer func() {
        if r := recover(); r != nil {
            slog.Error("HTTP server panicked", slog.Any("panic", r))
        }
    }()

    if err := srv.ListenAndServe(); err != http.ErrServerClosed {
        slog.Error("HTTP server ListenAndServe", "error", err)
        return
    }

    slog.Info("bye")
}
```

---

## 3. `config/` — single `config.go`

Nested structs + `sync.Once` singleton loader `C()`. No separate `env.go` or `parser.go`.

```go
package config

import (
    "fmt"
    "log"
    "sync"
    "time"

    "gitdev.devops.krungthai.com/ktb-digital-lending/dgl-new/shareable/common/codec"
    commonconfig "gitdev.devops.krungthai.com/ktb-digital-lending/dgl-new/shareable/common/config"
    env "github.com/caarlos0/env/v11"
)

type Config struct {
    Server        Server
    AccessControl AccessControl
    Firestore     Firestore
    Header        Header
    JWT           JWT
    Consumer      Consumer
    Producer      Producer
    HttpClient    HttpClient
    MySQL         MySQL
    Redis         Redis
    // Add domain-specific groups here
}

type Server struct {
    Hostname string `env:"HOSTNAME"`
    Port     string `env:"PORT,notEmpty"`
}

type JWT struct {
    Issuer      string        `env:"JWT_ISSUER,notEmpty"`
    Audience    string        `env:"JWT_AUDIENCE,notEmpty"`
    ExpDuration time.Duration `env:"JWT_EXP_DURATION,notEmpty"`
    PrivateKey  string        `env:"SECRET_JWT_PRIVATE_KEY,notEmpty"`
}

type Consumer struct {
    Brokers           string `env:"KAFKA_BROKERS"`
    GroupID           string `env:"KAFKA_GROUP_ID"`
    Topics            string `env:"KAFKA_TOPICS"`
    OffsetsInitial    string `env:"KAFKA_OFFSETS_INITIAL" envDefault:"latest"`
    RebalanceStrategy string `env:"KAFKA_REBALANCE_STRATEGY" envDefault:"range"`
}

type Producer struct {
    Brokers string `env:"KAFKA_PRODUCER_BROKERS"`
    Env     string `env:"ENV" envDefault:"LOCAL"`
}

type MySQL struct {
    Host     string `env:"MYSQL_HOST"`
    Port     string `env:"MYSQL_PORT"`
    User     string `env:"SECRET_DB_USERNAME"`
    Password string `env:"SECRET_DB_PASSWORD"`
    DBName   string `env:"MYSQL_DATABASE"`
}

type Redis struct {
    Addr     string `env:"REDIS_ADDR"`
    Password string `env:"SECRET_REDIS_PASSWORD"`
}

var once sync.Once
var config Config

func prefix(e string) string {
    if e == "" {
        return ""
    }
    return fmt.Sprintf("%s_", e)
}

// C returns the singleton config. envPrefix is typically commonconfig.Env.
func C(envPrefix string) Config {
    once.Do(func() {
        opts := env.Options{Prefix: prefix(envPrefix)}
        var err error
        config, err = commonconfig.ParseEnv[Config](opts)
        if err != nil {
            log.Fatal(err)
        }
        // Post-processing: decode base64-encoded secrets
        base64Coder := codec.NewBase64Coder()
        rawJWTPrivateKey, err := base64Coder.DecodeBase64(config.JWT.PrivateKey)
        if err != nil {
            log.Fatal(err)
        }
        config.JWT.PrivateKey = rawJWTPrivateKey
    })
    return config
}
```

---

## 4. `router/deps.go` — Composition Root

Unexported `deps` struct holding all infrastructure clients. `newDeps` returns the struct plus a `cleanup` function.

```go
package router

type deps struct {
    cfg             config.Config
    httpClient      *http.Client
    firestoreClient *commonfirestore.Client
    mysqlClient     *sql.DB
    redisClient     redis.UniversalClient
    producer        kafka.Producer
    hash            hash.HashManager
    cipher          crypt.Cipher
    token           token.JWTSigner
}

func newDeps(ctx context.Context, cfg config.Config) (deps, func()) {
    httpClient := httpclient.NewHTTPClient(
        middleware.ForwardRefIDOption,
        httpclient.DebugOption(cfg.HttpClient.EnableLogDebug),
    )
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
        _ = producer.Close()
    }
    return d, cleanup
}
```

> **Two categories of deps:**
>
> 1. **Raw SDK handles** (Firestore, MySQL, Redis, S3/GCS, `*http.Client`) — must be wrapped in `access/` before reaching domain handlers.
> 2. **Common-module abstractions** (`hash.HashManager`, `crypt.Cipher`, `token.JWTSigner`, `kafka.Producer`) — already interfaces; passed directly to `HandlerConfig`.

---

## 5. `router/router.go` — Gin engine + middleware + route registration

Middleware order matters. Health probes register **before** the auth chain.

```go
func New(cfg config.Config, version, commit string, timeoutDuration time.Duration) (*gin.Engine, func()) {
    r := gin.New()
    r.Use(gin.Recovery())

    if commonconfig.IsLocalEnv() {
        r.Use(gin.Logger())
    }

    ctx := context.Background()

    // Health probes (before auth middleware)
    r.GET("/liveness", health.Liveness(version, commit))
    r.GET("/metrics", health.Metrics())
    r.GET("/readiness", health.Readiness())

    // Middleware chain — order matters
    r.Use(
        middleware.SecurityHeaders(),
        middleware.AccessControl(cfg.AccessControl.AllowOrigin, allowedHeaders(cfg.Header.RefIDHeaderKey)),
        middleware.TraceContextTraceIDMiddleware(""),
        middleware.RefIDMiddleware(cfg.Header.RefIDHeaderKey),
        middleware.AutoLoggingMiddleware(app.CodeSuccess),
        middleware.Timeout(timeoutDuration),
        middleware.AccessLog(),
    )

    d, cleanupDeps := newDeps(ctx, cfg)

    // Register routes per domain (each domain is a bounded context)
    registerAuthRoutes(r, d)
    registerOrganizationRoutes(r, d)
    registerMemberRoutes(r, d)
    registerProductRoutes(r, d)

    return r, cleanupDeps
}
```

Each domain gets a `register<Domain>Routes` function:

```go
func registerProductRoutes(r *gin.Engine, d deps) {
    productStorage := productaccess.NewProductStorage(d.firestoreClient.Inner())
    h := product.NewHandler(product.HandlerConfig{ProductStorage: productStorage})

    g := r.Group("/api/v1/platform/product")
    {
        g.POST("/create", h.CreateProduct)
        g.POST("/detail", h.GetProduct)
        g.POST("/list", h.ListProducts)
        g.POST("/update", h.UpdateProduct)
        g.POST("/delete", h.DeleteProduct)
    }
}
```

---

## 6. `router/subscriber.go` — Kafka event router

`mergeRoutes` panics on duplicate event names, catching wiring mistakes at startup.

```go
// StartSubscriber starts Kafka consumer group; returns done channel + stop function.
func StartSubscriber(ctx context.Context, cfg config.Config) (<-chan struct{}, func()) {
    done := make(chan struct{})
    if !subscriberConfigured(cfg) {
        close(done)
        return done, func() {}
    }
    d, cleanupDeps := newDeps(ctx, cfg)
    eventHandlers := registerEventRoutes(d)
    consumerCtx, cancel := context.WithCancel(ctx)

    go func() {
        defer close(done)
        defer cleanupDeps()
        // ... run consumer loop
    }()

    return done, cancel
}

// WaitForSubscriber blocks until done closes or timeout.
func WaitForSubscriber(done <-chan struct{}, timeout time.Duration) { /* ... */ }

// registerEventRoutes builds the full event→handler map.
func registerEventRoutes(d deps) map[string]kafka.KafkaHandler {
    routes := make(map[string]kafka.KafkaHandler)
    mergeRoutes(routes, registerProductEvents(d))
    mergeRoutes(routes, registerOrganizationEvents(d))
    mergeRoutes(routes, registerMemberEvents(d))
    return routes
}

// mergeRoutes copies src into dst, panicking on duplicate event names.
func mergeRoutes(dst, src map[string]kafka.KafkaHandler) {
    for event, handler := range src {
        if _, exists := dst[event]; exists {
            panic(fmt.Sprintf("duplicate event route: %s", event))
        }
        dst[event] = handler
    }
}

func registerProductEvents(d deps) map[string]kafka.KafkaHandler {
    productStorage := productaccess.NewProductStorage(d.firestoreClient.Inner())
    h := product.NewHandler(product.HandlerConfig{ProductStorage: productStorage})
    return map[string]kafka.KafkaHandler{
        "PRODUCT_CREATED": h.OnProductCreated,
        "PRODUCT_UPDATED": h.OnProductUpdated,
        "PRODUCT_DELETED": h.OnProductDeleted,
    }
}
```

> **`kafka.KafkaHandler`** type is `func(ctx context.Context, msg kafka.Message[json.RawMessage]) error`.

---

## 7. Adding a new infrastructure client

**Raw infrastructure SDK** (Firestore, MySQL, Redis, S3/GCS, `*http.Client`):

1. Add the SDK client field to `deps` struct in `router/deps.go`.
2. Initialize it in `newDeps()` and add its `Close()` to the cleanup function.
3. Create `access/storage_<dep>.go`, `access/cache_<dep>.go`, or `access/client_<dep>.go` in the consuming domain.
4. In `register<Domain>Routes` / `register<Domain>Events`, construct the access-layer impl and pass it via `HandlerConfig`.

**Common-module abstraction** (`hash.HashManager`, `crypt.Cipher`, `token.JWTSigner`, `kafka.Producer`):

1. Add the interface field to `deps` struct (if not already there).
2. Initialize it in `newDeps()` with the corresponding `MustNew*` constructor.
3. Pass it directly to `HandlerConfig` (no `access/` wrapper needed).
