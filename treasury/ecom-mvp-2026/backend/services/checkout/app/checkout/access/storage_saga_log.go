package access

// storage_saga_log.go — write-only saga audit log.
// NOT on the request path. Written once at commit success or at any
// compensation-failure point for postmortem / forensics.
// See td.json §persistence.tables.saga_log.

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// SagaStatus values for the saga_log.status column.
type SagaStatus string

const (
	SagaCompleted   SagaStatus = "COMPLETED"
	SagaCompensated SagaStatus = "COMPENSATED"
	SagaDegraded    SagaStatus = "DEGRADED"
)

// SagaStep records one orchestration step's outcome.
type SagaStep struct {
	Step      int       `json:"step"`
	Name      string    `json:"name"`
	StartedAt time.Time `json:"startedAt"`
	EndedAt   time.Time `json:"endedAt"`
	Outcome   string    `json:"outcome"`
}

// SagaLogStorage is a simple write-only store for the saga_log table.
type SagaLogStorage struct{}

func NewSagaLogStorage() *SagaLogStorage { return &SagaLogStorage{} }

// Write inserts one saga_log row inside the caller's transaction.
// compensationOutcome is nil for COMPLETED status.
func (s *SagaLogStorage) Write(
	ctx context.Context,
	tx pgx.Tx,
	orderID, customerUserID, traceID string,
	status SagaStatus,
	steps []SagaStep,
	compensationOutcome json.RawMessage,
) error {
	stepsJSON, err := json.Marshal(steps)
	if err != nil {
		return fmt.Errorf("saga_log marshal steps: %w", err)
	}

	const q = `
		INSERT INTO checkout.saga_log
		    (order_id, customer_user_id, status, steps, compensation_outcome, trace_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())`

	_, err = tx.Exec(ctx, q, orderID, customerUserID, string(status), stepsJSON, []byte(compensationOutcome), traceID)
	if err != nil {
		return fmt.Errorf("saga_log insert: %w", err)
	}
	return nil
}
