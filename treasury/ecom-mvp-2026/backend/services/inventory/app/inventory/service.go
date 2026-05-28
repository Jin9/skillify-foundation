package inventory

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/inventory/config"
)

// ─── Storage interfaces ───────────────────────────────────────────────────────
//
// Defined here (in the domain package) to break the import cycle:
//   access/* imports inventory (for domain types).
//   service.go must NOT import access — store interfaces instead.
// Only main.go imports both packages and injects concrete access implementations.

// StockStorage is the storage interface for the stock_levels table.
type StockStorage interface {
	GetBySKU(ctx context.Context, sku string) (*StockLevel, error)
	GetBySKUForUpdate(ctx context.Context, tx pgx.Tx, sku string) (*StockLevel, error)
	GetManyBySKUsForUpdate(ctx context.Context, tx pgx.Tx, skus []string) (map[string]*StockLevel, error)
	DecrementAvailableIncrementReserved(ctx context.Context, tx pgx.Tx, sku string, qty int) error
	IncrementAvailableDecrementReserved(ctx context.Context, tx pgx.Tx, sku string, qty int) error
	DecrementReservedIncrementSold(ctx context.Context, tx pgx.Tx, sku string, qty int) error
	DecrementSoldIncrementAvailable(ctx context.Context, tx pgx.Tx, sku string, qty int) error
	AdjustQuantities(ctx context.Context, tx pgx.Tx, sku string, delta int) error
	Insert(ctx context.Context, tx pgx.Tx, sku string) error
}

// ReservationStorage is the storage interface for the reservations table.
type ReservationStorage interface {
	Insert(ctx context.Context, tx pgx.Tx, r *Reservation) error
	FindByOrderID(ctx context.Context, tx pgx.Tx, orderID uuid.UUID) ([]*Reservation, error)
	FindByOrderIDForUpdate(ctx context.Context, tx pgx.Tx, orderID uuid.UUID) ([]*Reservation, error)
	FindExpiredCandidates(ctx context.Context, tx pgx.Tx, batchSize int) ([]*Reservation, error)
	GetByIDForUpdate(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*Reservation, error)
	MarkExpired(ctx context.Context, tx pgx.Tx, id uuid.UUID) error
	MarkCommitted(ctx context.Context, tx pgx.Tx, id uuid.UUID) error
	MarkReleased(ctx context.Context, tx pgx.Tx, id uuid.UUID, reason ReleaseReason) error
}

// OutboxStorage is the storage interface for the outbox_events table.
type OutboxStorage interface {
	Insert(ctx context.Context, tx pgx.Tx, evt *OutboxEvent) error
	FetchUnpublished(ctx context.Context, tx pgx.Tx, limit int) ([]*OutboxEvent, error)
	MarkPublished(ctx context.Context, tx pgx.Tx, id interface{}) error
}

// ConsumedStorage is the storage interface for the consumed_events dedup table.
type ConsumedStorage interface {
	TryInsert(ctx context.Context, tx pgx.Tx, eventID uuid.UUID, consumerName string) (bool, error)
}

// ─── Service ──────────────────────────────────────────────────────────────────

// Service wires together all storage dependencies for the inventory domain.
// It is passed to handlers; each handler calls the appropriate storage methods
// directly within its own transaction boundary.
type Service struct {
	Pool            *pgxpool.Pool
	StockStorage    StockStorage
	ResvStorage     ReservationStorage
	OutboxStorage   OutboxStorage
	ConsumedStorage ConsumedStorage
	Cfg             config.Config
}

// ServiceConfig holds the dependencies for NewService.
// main.go constructs the concrete access implementations and passes them here
// so that service.go never needs to import the access package.
type ServiceConfig struct {
	Pool            *pgxpool.Pool
	StockStorage    StockStorage
	ResvStorage     ReservationStorage
	OutboxStorage   OutboxStorage
	ConsumedStorage ConsumedStorage
	Cfg             config.Config
}

// NewService constructs a Service from injected dependencies.
func NewService(cfg ServiceConfig) *Service {
	return &Service{
		Pool:            cfg.Pool,
		StockStorage:    cfg.StockStorage,
		ResvStorage:     cfg.ResvStorage,
		OutboxStorage:   cfg.OutboxStorage,
		ConsumedStorage: cfg.ConsumedStorage,
		Cfg:             cfg.Cfg,
	}
}
