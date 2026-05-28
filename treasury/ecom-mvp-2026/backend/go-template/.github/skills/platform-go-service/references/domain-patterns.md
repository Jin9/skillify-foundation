# Domain Package Patterns — `app/<domain>/`

Each domain aggregate lives in its own sub-package under `app/`. This file holds the canonical code shape for `access/`, `handler.go`, `handler_<action>.go`, and `consumer_<action>.go`.

---

## 1. `access/` — Repository Pattern

The `access/` layer is the **infrastructure boundary** — the thinnest possible adapter between the domain and external systems. Every file in `access/` follows the same co-location pattern:

| File prefix | Purpose | Wraps | Example |
|---|---|---|---|
| `storage_<dep>.go` | Persistence repository (CRUD) | Firestore, MySQL, PostgreSQL, S3/GCS | `storage_product.go` |
| `cache_<dep>.go` | Cache repository (get/set/delete with TTL) | Redis, Memcached | `cache_product.go` |
| `client_<dep>.go` | External API gateway | HTTP APIs, gRPC services | `client_google.go` |

Each file bundles together:

- **Interface** — the domain contract (e.g. `ProductStorage`, `ProductCache`, `GoogleClient`).
- **Unexported impl struct** — holds the SDK client handle.
- **Compile-time check** — `var _ Interface = (*impl)(nil)`.
- **Constructor** — `New<Interface>(sdkClient) Interface` (returns the interface).
- **Domain model / types** — entity structs, value objects, sentinel errors, constants.

> **Nothing else** goes in `access/`. No separate `model.go`, `errors.go`, or `constants.go` — models, errors, and constants are co-located in the file that uses them.

### Storage example — Firestore

```go
// access/storage_product.go
package access

type ProductStorage interface {
    GetProductByID(ctx context.Context, id uuid.UUID) (Product, error)
    ListProducts(ctx context.Context, organizationID uuid.UUID) ([]Product, error)
    CreateProduct(ctx context.Context, product Product) (Product, error)
    UpdateProduct(ctx context.Context, product Product) (Product, error)
    DeleteProduct(ctx context.Context, id uuid.UUID) error
}

// Domain model (Entity in DDD terms)
type ProductStatusType string

var (
    ProductStatusActive   ProductStatusType = "ACTIVE"
    ProductStatusInactive ProductStatusType = "INACTIVE"
)

type Product struct {
    ProductID      string            `firestore:"product_id" json:"productId"`
    Name           string            `firestore:"name" json:"name"`
    Description    string            `firestore:"description" json:"description"`
    Price          float64           `firestore:"price" json:"price"`
    OrganizationID string            `firestore:"organization_id" json:"organizationId"`
    Status         ProductStatusType `firestore:"status" json:"status"`
    CreatedAt      time.Time         `firestore:"created_at" json:"createdAt"`
    UpdatedAt      time.Time         `firestore:"updated_at" json:"updatedAt"`
}

// Model getters return (T, error) for parsed fields
func (p Product) GetID() (uuid.UUID, error) {
    return uuid.Parse(p.ProductID)
}

// productStorage implements ProductStorage (unexported — consumers use the interface)
type productStorage struct {
    fs *gcpfirestore.Client
}

var _ ProductStorage = (*productStorage)(nil) // compile-time check

func NewProductStorage(fs *gcpfirestore.Client) ProductStorage {
    return &productStorage{fs: fs}
}

func (s *productStorage) CreateProduct(ctx context.Context, product Product) (Product, error) {
    if product.ProductID == "" {
        product.ProductID = uuid.New().String()
    }
    app.SetTimestamps(&product.CreatedAt, &product.UpdatedAt)

    _, err := s.fs.Collection(productCollection).Doc(product.ProductID).Set(ctx, product)
    if err != nil {
        return Product{}, fmt.Errorf("failed to create product: %w", err)
    }
    return product, nil
}
```

### Gateway example — external HTTP API

```go
// access/client_google.go
type GoogleClient interface {
    VerifyToken(ctx context.Context, accessToken string) (*TokenInfo, error)
    GetUserProfile(ctx context.Context, accessToken string) (*UserProfile, error)
    RevokeToken(ctx context.Context, accessToken string) error
}
```

### Cache example — Redis

```go
// access/cache_product.go
package access

type ProductCache interface {
    GetProduct(ctx context.Context, id string) (Product, error)
    SetProduct(ctx context.Context, product Product, ttl time.Duration) error
    DeleteProduct(ctx context.Context, id string) error
}

type productCache struct {
    rdb redis.UniversalClient
}

var _ ProductCache = (*productCache)(nil)

func NewProductCache(rdb redis.UniversalClient) ProductCache {
    return &productCache{rdb: rdb}
}
```

---

## 2. `handler.go` — Config + Struct + Constructor

```go
package <domain>

import "<module>/app/<domain>/access"

type HandlerConfig struct {
    // access/ interfaces — domain-specific wrappers around raw infrastructure SDKs
    ProductStorage access.ProductStorage
    ProductCache   access.ProductCache   // only if domain uses caching
    // common module interfaces — already abstract, passed directly from deps
    // (hash.HashManager, crypt.Cipher, token.JWTSigner, kafka.Producer, etc.)
}

type handler struct {
    productStorage access.ProductStorage
    productCache   access.ProductCache
}

func NewHandler(cfg HandlerConfig) *handler {
    return &handler{
        productStorage: cfg.ProductStorage,
        productCache:   cfg.ProductCache,
    }
}
```

---

## 3. HTTP Handler — `handler_<action>.go`

```go
func (h *handler) CreateProduct(c *gin.Context) {
    ctx := c.Request.Context()

    req, ok := wrapper.BindJSON[CreateProductRequest](c, slog.String("handler", "CreateProduct"))
    if !ok {
        return
    }

    product := access.Product{
        Name:           req.Name,
        Description:    req.Description,
        Price:          req.Price,
        OrganizationID: req.OrganizationID,
        Status:         access.ProductStatusActive,
    }

    created, err := h.productStorage.CreateProduct(ctx, product)
    if err != nil {
        wrapper.Respond(c, wrapper.ResponseOption[CreateProductResponse]{
            HTTPStatus: http.StatusInternalServerError,
            Code:       app.CodeInternalError,
            Message:    app.MessageInternalError,
            Err: serror.Wrap(err).With(
                slog.String("product_name", req.Name),
                slog.String("organization_id", req.OrganizationID),
            ),
        })
        return
    }

    productID, err := created.GetID()
    if err != nil {
        wrapper.Respond(c, wrapper.ResponseOption[CreateProductResponse]{
            HTTPStatus: http.StatusInternalServerError,
            Code:       app.CodeInternalError,
            Message:    app.MessageInternalError,
            Err:        serror.Wrap(err),
        })
        return
    }

    wrapper.Respond(c, wrapper.ResponseOption[CreateProductResponse]{
        HTTPStatus: http.StatusCreated,
        Code:       app.CodeSuccess,
        Message:    app.MessageSuccess,
        Data:       &CreateProductResponse{ProductID: productID},
    })
}
```

Key conventions:

- Extract `ctx := c.Request.Context()` at the top.
- Use `wrapper.BindJSON[T]` for request binding (returns `(T, bool)` — auto-responds 400 on failure).
- Use `wrapper.ResponseOption[T]` generic — the `T` matches the response struct.
- Use `app.CodeSuccess` / `app.CodeBadRequest` / `app.CodeInternalError` (from common module).
- Wrap errors with `serror.Wrap(err).With(slog attrs...)` for structured observability.
- Each `if err != nil` branch returns immediately with the appropriate HTTP status.

---

## 4. Kafka Consumer — `consumer_<action>.go`

```go
func (h *handler) OnProductCreated(ctx context.Context, msg kafka.Message[json.RawMessage]) error {
    var payload CreateProductMessage
    if err := kafka.BindMessage(msg.Payload, &payload); err != nil {
        return serror.Wrap(err).With(slog.String("event_id", msg.EventID))
    }

    product := access.Product{
        Name:           payload.Name,
        Description:    payload.Description,
        Price:          payload.Price,
        OrganizationID: payload.OrganizationID,
        Status:         access.ProductStatusActive,
    }

    _, err := h.productStorage.CreateProduct(ctx, product)
    if err != nil {
        return serror.Wrap(err).With(
            slog.String("event_id", msg.EventID),
            slog.String("product_name", payload.Name),
            slog.String("organization_id", payload.OrganizationID),
        )
    }
    return nil
}
```

Key conventions:

- Signature: `func(ctx context.Context, msg kafka.Message[json.RawMessage]) error` — includes `ctx`.
- Use `kafka.BindMessage(msg.Payload, &payload)` — **NOT** `json.Unmarshal`. `BindMessage` deserializes **and** validates (binding tags are enforced).
- The `kafka.Message[T]` fields include `EventID`, `AggregateID`, `EventName`, `Timestamp`, `Payload`.
- Return `serror.Wrap` errors with investigation context.

---

## 5. Response envelope (`wrapper`)

```go
wrapper.Respond(c, wrapper.ResponseOption[MyResponse]{
    HTTPStatus: http.StatusOK,
    Code:       app.CodeSuccess,
    Message:    app.MessageSuccess,
    Data:       &myResp,
})
```

Produces:

```json
{ "code": "0000", "message": "success", "data": { ... } }
```

Response codes and messages live in the **common** module as typed constants:

| Constant | Value |
|---|---|
| `app.CodeSuccess` | `"0000"` |
| `app.MessageSuccess` | `"success"` |
| `app.CodeBadRequest` | varies |
| `app.CodeInternalError` | varies |

For the `Err` field, use `serror.Wrap(err)` for structured error logging (the error is logged but not exposed in the JSON response).

---

## 6. Service helpers — `service_<action>.go`

When a handler accumulates private methods that share a single responsibility (e.g. external API orchestration, data assembly), extract them into `service_<action>.go` files grouped by purpose. The handler file keeps only the public HTTP/Kafka entry point and its request/response types.

The service file holds private methods on `*handler`. They share `*handler`'s injected dependencies; they are **not** a separate struct.

---

## 7. Kafka message envelope

```go
type Message[T any] struct {
    EventID     string    `json:"eventId"`
    EventName   string    `json:"eventName"`
    AggregateID string    `json:"aggregateId"`
    Timestamp   time.Time `json:"timestamp"`
    Payload     T         `json:"payload"`
}
```

Consumers receive `kafka.Message[json.RawMessage]` and decode `Payload` via `kafka.BindMessage`.

### Producer logging interceptor

```go
producer := kafka.MustNewProducer(cfg)

// Auto-select mode from ENV (LOCAL/DEV→Debug, UAT/PROD→Meta)
logged := kafka.WithLoggingFromEnv(producer, os.Getenv("ENV"))

// Or explicit mode
debug  := kafka.WithLogging(producer, kafka.LogInterceptorOption{Mode: kafka.LogModeDebug})
meta   := kafka.WithLogging(producer, kafka.LogInterceptorOption{Mode: kafka.LogModeMeta})
silent := kafka.WithLogging(producer, kafka.LogInterceptorOption{Mode: kafka.LogModeSilent})

// In PROD (Meta mode), attach investigation attrs — no payload is logged
logged.SendMessageWithOption("topic", msg, kafka.SendMessageOption{
    LogAttrs: []slog.Attr{
        slog.String("organization_id", orgID),
    },
})
```
