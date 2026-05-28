package order

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// StorageOrder defines the repository contract for the orders table.
// All methods that mutate orders.status use a pgx.Tx to enforce the
// SELECT ... FOR UPDATE → validate → write pattern.
type StorageOrder interface {
	Insert(ctx context.Context, tx pgx.Tx, o Order) error
	GetByIDForUpdate(ctx context.Context, tx pgx.Tx, id string) (*Order, error)
	GetByIDScopedForUpdate(ctx context.Context, tx pgx.Tx, id, userID string) (*Order, error)
	GetByID(ctx context.Context, id string) (*Order, error)
	UpdateStatus(ctx context.Context, tx pgx.Tx, id string, newStatus OrderStatus, version int, trackingNumber *string) error
	ListByUserID(ctx context.Context, userID string, status *OrderStatus, limit, offset int) ([]Order, int, error)
	NextOrderNumber(ctx context.Context, tx pgx.Tx, nowUTC string) (string, error)
}

// StorageOrderItem defines the repository contract for order_items.
// No Update method is exposed — immutability is enforced at the repo layer (ORD-007).
type StorageOrderItem interface {
	BulkInsert(ctx context.Context, tx pgx.Tx, items []OrderItem) error
	ListByOrderID(ctx context.Context, orderID string) ([]OrderItem, error)
}

// StorageStatusHistory defines the repository contract for order_status_history.
type StorageStatusHistory interface {
	Insert(ctx context.Context, tx pgx.Tx, row OrderStatusHistory) error
	ListByOrderIDASC(ctx context.Context, orderID string) ([]OrderStatusHistory, error)
}

// StorageOutbox defines the repository contract for outbox_events.
type StorageOutbox interface {
	Insert(ctx context.Context, tx pgx.Tx, event OutboxEvent) error
	MarkPublished(ctx context.Context, id string) error
	PollUnpublished(ctx context.Context, limit int) ([]OutboxEvent, error)
}

// StorageConsumedEvent defines the repository contract for consumed_events.
// InsertOrConflict returns true when the row was newly inserted (first-time processing).
type StorageConsumedEvent interface {
	InsertOrConflict(ctx context.Context, tx pgx.Tx, eventID, consumerName, eventType, orderID string) (inserted bool, err error)
}

// Service wires all storage dependencies together for the order aggregate.
type Service struct {
	pool          *pgxpool.Pool
	orders        StorageOrder
	items         StorageOrderItem
	statusHistory StorageStatusHistory
	outbox        StorageOutbox
	consumed      StorageConsumedEvent
}

// NewService constructs a Service with all dependencies.
func NewService(
	pool *pgxpool.Pool,
	orders StorageOrder,
	items StorageOrderItem,
	statusHistory StorageStatusHistory,
	outbox StorageOutbox,
	consumed StorageConsumedEvent,
) *Service {
	return &Service{
		pool:          pool,
		orders:        orders,
		items:         items,
		statusHistory: statusHistory,
		outbox:        outbox,
		consumed:      consumed,
	}
}

// beginTx is a thin helper that starts a serializable-read-committed pgx transaction.
func (s *Service) beginTx(ctx context.Context) (pgx.Tx, error) {
	return s.pool.Begin(ctx)
}
