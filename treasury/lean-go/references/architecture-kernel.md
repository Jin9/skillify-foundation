# Architecture kernel: transport / business / integration

The kernel is three roles and one dependency rule. Folder names vary by repo; the roles do not.

## The three roles

### Transport (handler, consumer)
Entry point for one request or event. Thin.
- Bind and validate the incoming payload.
- Call business logic.
- Map the result or error to the response shape (HTTP envelope, ack/nack, returned `error`).
- No persistence calls inline, no domain rules — delegate both.

### Business (service, or the handler method while small)
Where decisions live.
- Request rules, orchestration across access calls, domain validation, idempotency.
- Starts as the body of the handler method. Extract to a `service` unit **only when** the method accumulates real orchestration (Fowler Extract Function) — not preemptively. A one-call passthrough does not need a service.
- Depends on access **interfaces**, never on concrete infrastructure.

### Integration (access)
Lean adapters to the outside world. **No business logic.**
- `storage` — Repository over a database.
- `client` — Gateway over an external HTTP/gRPC service.
- `cache` — Cache over Redis/in-memory.
- Each is an interface plus an unexported implementation that wraps the infra client, translates infra errors into domain errors (sentinels), and returns plain domain types.

## The dependency rule
```
transport ──> business ──> access interface
                                  ^
                                  └── access impl (wraps infra)
```
- Business depends on access **interfaces**, declared at the access layer (the data owner).
- Inject dependencies through a constructor config so they are swappable and mockable.
- Infrastructure types never travel up into transport; business rules never travel down into access.

## Access shape (Repository, grounded and portable)
```go
package access

var ErrProductNotFound = errors.New("product not found")

type ProductStorage interface {
	GetByID(ctx context.Context, id uuid.UUID) (Product, error)
	Create(ctx context.Context, p Product) (Product, error)
}

type productStorage struct{ db *infra.Client } // infra dependency only

var _ ProductStorage = (*productStorage)(nil) // compile-time conformance

func NewProductStorage(db *infra.Client) ProductStorage { return &productStorage{db: db} }

func (s *productStorage) GetByID(ctx context.Context, id uuid.UUID) (Product, error) {
	doc, err := s.db.Get(ctx, id.String())
	if err != nil {
		if errors.Is(err, infra.ErrNotFound) {
			return Product{}, ErrProductNotFound // translate infra error to domain error
		}
		return Product{}, fmt.Errorf("get product: %w", err)
	}
	var p Product
	if err := doc.DataTo(&p); err != nil {
		return Product{}, fmt.Errorf("decode product: %w", err)
	}
	return p, nil
}
```

The access constructor returns its **interface** (`NewProductStorage(...) ProductStorage`) — the DI seam owned by the data owner is the one accepted exception to "return concrete types."

## Transport + DI seam (grounded and portable)
```go
package product

type HandlerConfig struct {
	ProductStorage access.ProductStorage // interface, so tests inject a mock
}

type handler struct{ productStorage access.ProductStorage }

func NewHandler(cfg HandlerConfig) *handler { return &handler{productStorage: cfg.ProductStorage} }
```

Business logic sits in the handler method (or an extracted `service_<action>.go` unexported method on `*handler` once it grows): bind → decide → call access → map response.

## Map the roles onto any structure
Detect the repo's name for each role, then conform — do not impose `handler/service/access`.

| Layout | transport | business | integration |
|--------|-----------|----------|-------------|
| flat single package | `*Handler` funcs in the package | functions in the same package | `*Store`/`*Client` in the same package |
| `cmd/internal/pkg` | `internal/http` or `internal/api` | `internal/<domain>` | `internal/<domain>/repo` or `internal/store` |
| layered | `handler/` | `service/` or `usecase/` | `repository/` / `store/` / `gateway/` |
| DDD | `interfaces/`/`app/` | `domain/` + application service | `infrastructure/` / `adapters/` |
| hexagonal | inbound `adapters/` | core / `application` | outbound `adapters/` (ports impl) |

Rules that hold in every layout:
- Keep business out of the integration layer.
- Business depends on an interface owned by the integration layer.
- Add a layer only when there is logic to put in it; a pass-through layer is a smell (see `refactoring-playbook.md`).
