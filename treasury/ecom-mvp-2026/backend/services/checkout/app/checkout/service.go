package checkout

import (
	"gitlab.com/b2c-e-commerce-platform/platform/backend/checkout/app/checkout/access"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service bundles all dependencies for the checkout handlers.
// Both Preview and Commit are methods on this struct.
type Service struct {
	db               *pgxpool.Pool
	idempotencyStore *access.IdempotencyStorage
	sagaLogStore     *access.SagaLogStorage

	cart     *CartClient
	catalog  *CatalogClient
	inventory *InventoryClient
	payment  *PaymentClient
	order    *OrderClient
	identity *IdentityClient
}

// ServiceConfig groups all constructor arguments for clarity.
type ServiceConfig struct {
	DB               *pgxpool.Pool
	IdempotencyStore *access.IdempotencyStorage
	SagaLogStore     *access.SagaLogStorage
	Cart             *CartClient
	Catalog          *CatalogClient
	Inventory        *InventoryClient
	Payment          *PaymentClient
	Order            *OrderClient
	Identity         *IdentityClient
}

func NewService(cfg ServiceConfig) *Service {
	return &Service{
		db:               cfg.DB,
		idempotencyStore: cfg.IdempotencyStore,
		sagaLogStore:     cfg.SagaLogStore,
		cart:             cfg.Cart,
		catalog:          cfg.Catalog,
		inventory:        cfg.Inventory,
		payment:          cfg.Payment,
		order:            cfg.Order,
		identity:         cfg.Identity,
	}
}
