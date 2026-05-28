package cart

import (
	"context"

	"github.com/google/uuid"
)

// CartService is the service interface consumed by HTTP handlers.
// Defined here to break any handler → service concrete-type coupling and to
// let tests inject a lightweight stub without spinning up real storage or HTTP clients.
type CartService interface {
	AddItem(ctx context.Context, userID, productID uuid.UUID, qty int) (CartView, error)
	UpdateItem(ctx context.Context, userID, cartItemID uuid.UUID, qty int) (CartView, error)
	RemoveItem(ctx context.Context, userID, cartItemID uuid.UUID) (CartView, error)
	ReadCart(ctx context.Context, userID uuid.UUID) (CartView, error)
	ClearOnCheckout(ctx context.Context, userID uuid.UUID, cartItemIDs []uuid.UUID, orderID uuid.UUID) (int, error)
}

// Compile-time check: *Service must satisfy CartService.
var _ CartService = (*Service)(nil)

// HandlerConfig holds all dependencies for the cart handler group.
type HandlerConfig struct {
	Service CartService
}

type handler struct {
	service CartService
}

// NewHandler constructs the cart handler.
func NewHandler(cfg HandlerConfig) *handler {
	return &handler{
		service: cfg.Service,
	}
}
