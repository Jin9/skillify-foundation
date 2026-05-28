package access

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/services/order/app/order"
)

// StatusHistoryStorage implements order.StorageStatusHistory against PostgreSQL.
// Append-only — no Update or Delete methods are exposed (ORD-009).
type StatusHistoryStorage struct {
	pool *pgxpool.Pool
}

func NewStatusHistoryStorage(pool *pgxpool.Pool) *StatusHistoryStorage {
	return &StatusHistoryStorage{pool: pool}
}

// Insert writes a single status-history row inside the given transaction.
// Must be called in the SAME tx as the orders UPDATE (ORD-009 atomicity requirement).
func (s *StatusHistoryStorage) Insert(ctx context.Context, tx pgx.Tx, row order.OrderStatusHistory) error {
	const q = `
		INSERT INTO "order".order_status_history
			(id, order_id, from_status, to_status, actor_user_id, actor_role, reason, occurred_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	var fromStatus *string
	if row.FromStatus != nil {
		s := string(*row.FromStatus)
		fromStatus = &s
	}

	if _, err := tx.Exec(ctx, q,
		row.ID,
		row.OrderID,
		fromStatus,
		string(row.ToStatus),
		row.ActorUserID,
		string(row.ActorRole),
		row.Reason,
		row.OccurredAt,
	); err != nil {
		return fmt.Errorf("insert status history: %w", err)
	}
	return nil
}

// ListByOrderIDASC returns all history rows for an order sorted by occurred_at ASC.
// Used by order.detail to populate statusHistory[].
func (s *StatusHistoryStorage) ListByOrderIDASC(ctx context.Context, orderID string) ([]order.OrderStatusHistory, error) {
	const q = `
		SELECT id, order_id, from_status, to_status, actor_user_id, actor_role, reason, occurred_at
		FROM "order".order_status_history
		WHERE order_id = $1
		ORDER BY occurred_at ASC`

	rows, err := s.pool.Query(ctx, q, orderID)
	if err != nil {
		return nil, fmt.Errorf("list status history: %w", err)
	}
	defer rows.Close()

	result := make([]order.OrderStatusHistory, 0)
	for rows.Next() {
		var h order.OrderStatusHistory
		var fromStatus *string
		var toStatus, actorRole string
		if err := rows.Scan(
			&h.ID, &h.OrderID,
			&fromStatus, &toStatus,
			&h.ActorUserID, &actorRole,
			&h.Reason, &h.OccurredAt,
		); err != nil {
			return nil, fmt.Errorf("scan status history: %w", err)
		}
		if fromStatus != nil {
			s := order.OrderStatus(*fromStatus)
			h.FromStatus = &s
		}
		h.ToStatus = order.OrderStatus(toStatus)
		h.ActorRole = order.ActorRole(actorRole)
		result = append(result, h)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("status history rows: %w", err)
	}
	return result, nil
}
