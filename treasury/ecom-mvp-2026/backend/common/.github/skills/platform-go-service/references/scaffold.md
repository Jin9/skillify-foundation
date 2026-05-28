# Service Scaffold

Boilerplate for a brand-new service: module init, `main.go`, `config/`, `router/deps.go`, `router/router.go`, `router/subscriber.go`. Read this when creating a new service or wiring a new top-level dependency.

## 1. Initialise the Module

```bash
mkdir <service-name> && cd <service-name>
go mod init gitdev.devops.krungthai.com/a-team/poc-structure/<service-name>
go get gitdev.devops.krungthai.com/ktb-digital-lending/dgl-new/shareable/common
```

## 2. `main.go` — Bootstrap

```go
package main

import (
    "context"
    "fmt"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "gitdev.devops.krungthai.com/ktb-digital-lending/dgl-new/shareable/common/logger"

    "<module>/config"
    "<module>/router"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        panic(fmt.Sprintf("config: %v", err))
    }

    _ = logger.New(logger.CensorReplacer, logger.GCPKeyReplacer)

    deps := router.BuildDeps(cfg)
    defer deps.Close()

    r := router.New(cfg, deps)

    srv := &http.Server{
        Addr:    fmt.Sprintf(":%s", cfg.Port),
        Handler: r,
    }
    go func() {
        slog.Info("server started", slog.String("addr", srv.Addr))
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            slog.Error("server failed", slog.Any("error", err))
            os.Exit(1)
        }
    }()

    signalCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer cancel()
    router.StartConsumers(signalCtx, cfg, deps)

    <-signalCtx.Done()
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer shutdownCancel()
    _ = srv.Shutdown(shutdownCtx)
}
```

## 3. `config/`

`config.go`:

```go
type Config struct {
    Port string `env:"PORT" envDefault:"8080"`
    Env  string `env:"ENV"  envDefault:"LOCAL"`

    DBHost     string `env:"DATABASE_HOST"`
    DBPort     string `env:"DATABASE_PORT" envDefault:"5432"`
    DBName     string `env:"DATABASE_NAME"`
    DBUser     string `env:"DATABASE_USER"`
    DBPassword string `env:"DATABASE_PASSWORD"`

    KafkaBrokers string `env:"KAFKA_BROKERS"`
    KafkaGroupID string `env:"KAFKA_GROUP_ID"`
}
```

`parser.go`:

```go
func Load() (*Config, error) {
    return config.ParseEnv[Config](env.Options{Prefix: "APP_"})
}
```

## 4. `router/deps.go` — Dependency Container

```go
type Deps struct {
    DB            *pgxpool.Pool
    Producer      kafka.Producer
    TokenParser   token.JWTParser
    TokenVerifier token.JWTVerifier
    Clock         clock.Clock
}

func BuildDeps(cfg *config.Config) *Deps {
    db, err := database.ConnectPostgresDB(database.PostgresConfig{ /* ... */ })
    // ... wire all adapters ...
    return &Deps{DB: db /* , Producer: logged, ... */}
    _ = err
}

func (d *Deps) Close() {
    d.DB.Close()
    _ = d.Producer.Close()
}
```

## 5. `router/router.go` — Gin Engine + Middleware

```go
func New(cfg *config.Config, deps *Deps) *gin.Engine {
    r := gin.New()
    r.Use(
        middleware.AccessLog(),
        middleware.SecurityHeaders(),
        middleware.RefIDMiddleware("X-Request-ID"),
        middleware.TraceContextTraceIDMiddleware("traceparent"),
        middleware.Timeout(30 * time.Second),
    )

    r.GET("/liveness", health.Liveness(version, commit))
    r.GET("/readiness", health.Readiness())
    r.GET("/metrics", health.Metrics())

    v1 := r.Group("/v1")
    v1.Use(middleware.JWT(deps.TokenParser, deps.TokenVerifier))
    {
        productHandler := product.NewHandler(product.HandlerConfig{ /* ... */ })
        v1.POST("/products", productHandler.CreateProduct)
        v1.GET("/products/:id", productHandler.GetProduct)
    }
    return r
}
```

## 6. `router/subscriber.go` — Kafka Event→Handler Map

```go
func EventHandlers(deps *Deps) map[string]kafka.KafkaHandler {
    productHandler := product.NewHandler(product.HandlerConfig{ /* ... */ })
    return map[string]kafka.KafkaHandler{
        "PRODUCT_CREATED": productHandler.OnProductCreated,
    }
}

func StartConsumers(ctx context.Context, cfg *config.Config, deps *Deps) {
    group, _ := kafka.NewConsumerGroup(kafka.ConsumerConfig{ /* ... */ })
    processor := kafka.NewEventRouter(EventHandlers(deps))
    go group.Consume(ctx, cfg.KafkaTopics, kafka.NewConsumerGroupHandler(ctx, processor))
}
```
