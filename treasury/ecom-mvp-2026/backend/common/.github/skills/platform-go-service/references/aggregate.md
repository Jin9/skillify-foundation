# Domain Aggregate Recipes

Patterns for `app/<aggregate>/`. Read this when adding a new aggregate, a new handler, or a new consumer.

## 1. Handler Config + Constructor

```go
// handler.go
package product

type HandlerConfig struct {
    Storage  access.Storage
    Producer kafka.Producer
    Clock    clock.Clock
}

type handler struct {
    storage  access.Storage
    producer kafka.Producer
    clock    clock.Clock
}

func NewHandler(cfg HandlerConfig) *handler {
    return &handler{
        storage:  cfg.Storage,
        producer: cfg.Producer,
        clock:    cfg.Clock,
    }
}
```

## 2. HTTP Command Handler

```go
// create_handler.go
func (h *handler) CreateProduct(c *gin.Context) {
    req, ok := wrapper.BindJSON[CreateProductRequest](c)
    if !ok {
        return
    }

    result, err := h.storage.Insert(c.Request.Context(), req.ToModel())
    if err != nil {
        wrapper.Respond(c, wrapper.ResponseOption[any]{
            HTTPStatus: http.StatusInternalServerError,
            Code:       app.CodeInternalError,
            Message:    app.MessageInternalError,
            Err:        serror.Wrap(err),
        })
        return
    }

    msg := kafka.NewMessage("PRODUCT_CREATED", result.ID, result)
    _ = h.producer.SendMessage(c.Request.Context(), "product-events", msg)

    wrapper.Respond(c, wrapper.ResponseOption[*ProductResponse]{
        HTTPStatus: http.StatusCreated,
        Code:       app.CodeProductCreated,
        Message:    app.MessageProductCreated,
        Data:       toResponse(result),
    })
}
```

## 3. Kafka Consumer Handler

```go
// consumer_create.go
func (h *handler) OnProductCreated(ctx context.Context, msg kafka.Message[json.RawMessage]) error {
    var payload ProductCreatedEvent
    if err := kafka.BindMessage(msg.Payload, &payload); err != nil {
        return serror.Wrap(err)
    }

    if err := h.storage.Upsert(ctx, payload.ToModel()); err != nil {
        return serror.Wrap(err).With(
            slog.String("product_id", payload.ID),
        )
    }
    return nil
}
```

## 4. Access Layer (Ports)

`access/storage.go` — Port interface:

```go
package access

type Storage interface {
    Insert(ctx context.Context, model *Model) (*Model, error)
    FindByID(ctx context.Context, id string) (*Model, error)
    Upsert(ctx context.Context, model *Model) error
    Delete(ctx context.Context, id string) error
}
```

`access/model.go` — persistence/domain model:

```go
type Model struct {
    ID        string    `db:"id" firestore:"id"`
    Name      string    `db:"name" firestore:"name"`
    CreatedAt time.Time `db:"created_at" firestore:"created_at"`
    UpdatedAt time.Time `db:"updated_at" firestore:"updated_at"`
}
```

## 5. Response Codes — `app/codes.go`

```go
package app

import "gitdev.devops.krungthai.com/ktb-digital-lending/dgl-new/shareable/common/wrapper"

type Code = wrapper.Code
type Message = wrapper.Message
type Response[T any] = wrapper.Response[T]

// Re-export standard codes
const (
    CodeSuccess    = wrapper.CodeSuccess
    MessageSuccess = wrapper.MessageSuccess
    CodeBadRequest = wrapper.CodeBadRequest
    // ... etc.
)

// Domain-specific codes — prefix per aggregate
const (
    CodeProductCreated    Code    = "PD201"
    MessageProductCreated Message = "Product Created"
)
```
