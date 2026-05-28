package payment

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ─── Storage interfaces ───────────────────────────────────────────────────────
// Interfaces live here (in the domain package) to break the import cycle:
//   access → payment (for domain types) + payment → access (for interfaces)
// The access package's concrete types implement these interfaces without the cycle.

// IntentStorage defines persistence operations for payment_intents.
type IntentStorage interface {
	// Insert inserts a new PaymentIntent. Returns ErrDuplicateOrderID on conflict.
	Insert(ctx context.Context, pi PaymentIntent) error

	// GetByOrderID looks up an existing intent by order_id.
	// Returns pgx.ErrNoRows if not found.
	GetByOrderID(ctx context.Context, orderID uuid.UUID) (PaymentIntent, error)

	// GetByIDForUpdate fetches a row with SELECT ... FOR UPDATE (must be called
	// inside a transaction obtained from the pool).
	GetByIDForUpdate(ctx context.Context, tx pgx.Tx, intentID uuid.UUID) (PaymentIntent, error)

	// MarkSucceeded transitions status to SUCCEEDED within an existing tx.
	MarkSucceeded(ctx context.Context, tx pgx.Tx, intentID uuid.UUID, version int, providerRef string, paidAt time.Time) error

	// MarkFailed transitions status to FAILED within an existing tx.
	MarkFailed(ctx context.Context, tx pgx.Tx, intentID uuid.UUID, version int, providerRef string) error

	// MarkExpired transitions status to EXPIRED within an existing tx.
	MarkExpired(ctx context.Context, tx pgx.Tx, intentID uuid.UUID, version int) error

	// ListExpiredIDs returns up to batchSize intent_ids whose
	// status=REQUIRES_PAYMENT and expires_at < now, using FOR UPDATE SKIP LOCKED.
	// Called by the expiry sweeper.
	ListExpiredIDs(ctx context.Context, tx pgx.Tx, batchSize int) ([]uuid.UUID, error)
}

// CallbackDedupStorage defines persistence for payment_callback_dedup.
type CallbackDedupStorage interface {
	// LookupTx returns the dedup row for the given key within tx, or pgx.ErrNoRows if absent.
	LookupTx(ctx context.Context, tx pgx.Tx, dedupKey string) (CallbackDedup, error)

	// Insert writes a new dedup row inside tx.
	// ON CONFLICT DO NOTHING is the DB-level defence-in-depth against concurrent racers.
	Insert(ctx context.Context, tx pgx.Tx, row CallbackDedup) error
}

// OutboxStorage defines persistence for the payment.outbox table.
type OutboxStorage interface {
	// Insert writes an outbox row within the provided transaction.
	Insert(ctx context.Context, tx pgx.Tx, evt OutboxEvent) error

	// ListPending returns up to limit PENDING outbox rows ordered by created_at.
	ListPending(ctx context.Context, limit int) ([]OutboxEvent, error)

	// MarkPublished marks a row as PUBLISHED and records published_at.
	MarkPublished(ctx context.Context, id string) error

	// MarkFailed increments attempts and records last_error; marks status=FAILED
	// when attempts >= maxAttempts.
	MarkFailed(ctx context.Context, id, lastError string, maxAttempts int) error
}

// ─── Service ─────────────────────────────────────────────────────────────────

// Service defines all operations exposed by handlers.
type Service interface {
	// CreateIntent is idempotent on order_id.
	// Returns the existing intent if one already exists for the order.
	// Returns ErrConflict if the existing intent is in a terminal state.
	CreateIntent(ctx context.Context, req IntentCreateRequest) (IntentCreateResponse, error)

	// Simulate maps a customer outcome → providerStatus and invokes ProcessCallback
	// in-process. Auth/ownership checks are performed before calling this.
	Simulate(ctx context.Context, intentID uuid.UUID, outcome string) (ProcessCallbackOutput, error)

	// Callback is invoked by the HTTP handler after HMAC verification.
	Callback(ctx context.Context, in ProcessCallbackInput) (ProcessCallbackOutput, error)

	// GetIntentByID returns the intent for ownership checking in simulate.
	GetIntentByID(ctx context.Context, intentID uuid.UUID) (PaymentIntent, error)
}

// ServiceConfig holds constructor dependencies.
type ServiceConfig struct {
	DB               *pgxpool.Pool
	IntentStore      IntentStorage
	DedupStore       CallbackDedupStorage
	OutboxStore      OutboxStorage
	IntentTTL        time.Duration
	EmitExpiredEvent bool
}

type svc struct {
	db               *pgxpool.Pool
	intentStore      IntentStorage
	dedupStore       CallbackDedupStorage
	outboxStore      OutboxStorage
	intentTTL        time.Duration
	emitExpiredEvent bool
}

var _ Service = (*svc)(nil)

// NewService constructs the payment Service.
func NewService(cfg ServiceConfig) Service {
	return &svc{
		db:               cfg.DB,
		intentStore:      cfg.IntentStore,
		dedupStore:       cfg.DedupStore,
		outboxStore:      cfg.OutboxStore,
		intentTTL:        cfg.IntentTTL,
		emitExpiredEvent: cfg.EmitExpiredEvent,
	}
}

func (s *svc) CreateIntent(ctx context.Context, req IntentCreateRequest) (IntentCreateResponse, error) {
	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		return IntentCreateResponse{}, err
	}
	ownerID, err := uuid.Parse(req.OwnerUserID)
	if err != nil {
		return IntentCreateResponse{}, err
	}

	currency := req.Currency
	if currency == "" {
		currency = "THB"
	}

	expiresAt := time.Now().UTC().Add(s.intentTTL)

	// Idempotency: look up existing intent for this order.
	existing, err := s.intentStore.GetByOrderID(ctx, orderID)
	if err == nil {
		// Row exists.
		if TerminalStatuses[existing.Status] {
			return IntentCreateResponse{}, ErrConflict
		}
		// Status == REQUIRES_PAYMENT: return existing.
		return IntentCreateResponse{
			PaymentIntentID: existing.IntentID.String(),
			Status:          string(existing.Status),
			Amount:          existing.AmountMinor,
			ExpiresAt:       existing.ExpiresAt,
		}, nil
	}

	// New intent.
	pi := PaymentIntent{
		IntentID:    uuid.New(),
		OrderID:     orderID,
		OwnerUserID: ownerID,
		AmountMinor: req.Amount,
		Currency:    currency,
		Status:      StatusRequiresPayment,
		ExpiresAt:   expiresAt,
	}

	if err := s.intentStore.Insert(ctx, pi); err != nil {
		if err == ErrDuplicateOrderID {
			// Race: another request beat us; retry the lookup.
			existing, err2 := s.intentStore.GetByOrderID(ctx, orderID)
			if err2 != nil {
				return IntentCreateResponse{}, err2
			}
			if TerminalStatuses[existing.Status] {
				return IntentCreateResponse{}, ErrConflict
			}
			return IntentCreateResponse{
				PaymentIntentID: existing.IntentID.String(),
				Status:          string(existing.Status),
				Amount:          existing.AmountMinor,
				ExpiresAt:       existing.ExpiresAt,
			}, nil
		}
		return IntentCreateResponse{}, err
	}

	return IntentCreateResponse{
		PaymentIntentID: pi.IntentID.String(),
		Status:          string(pi.Status),
		Amount:          pi.AmountMinor,
		ExpiresAt:       pi.ExpiresAt,
	}, nil
}

func (s *svc) GetIntentByID(ctx context.Context, intentID uuid.UUID) (PaymentIntent, error) {
	// Reuse GetByOrderID-equivalent via a direct DB query path.
	// We need GetByID (non-locking) for the simulate ownership check.
	// Storage layer exposes GetByIDForUpdate which requires a tx; for a read-only
	// lookup we query directly. This is a thin wrapper for now.
	//
	// NOTE: GetByIDForUpdate is the tx-aware version; here we run outside a tx
	// so we open our own read-only query. Add GetByID to IntentStorage if needed;
	// for scaffold we replicate the scan inline via the pool.
	const q = `
		SELECT intent_id, order_id, owner_user_id, amount_minor, currency,
		       status, mock_provider_ref, provider_status, paid_at,
		       expires_at, created_at, updated_at, version
		FROM payment.payment_intents
		WHERE intent_id = $1`

	row := s.db.QueryRow(ctx, q, intentID)
	var pi PaymentIntent
	err := row.Scan(
		&pi.IntentID, &pi.OrderID, &pi.OwnerUserID, &pi.AmountMinor, &pi.Currency,
		&pi.Status, &pi.MockProviderRef, &pi.ProviderStatus, &pi.PaidAt,
		&pi.ExpiresAt, &pi.CreatedAt, &pi.UpdatedAt, &pi.Version,
	)
	if err != nil {
		return PaymentIntent{}, err
	}
	return pi, nil
}

func (s *svc) Simulate(ctx context.Context, intentID uuid.UUID, outcome string) (ProcessCallbackOutput, error) {
	// Map customer-friendly outcome to providerStatus.
	var providerStatus string
	switch outcome {
	case "success":
		providerStatus = "SUCCEEDED"
	case "failed":
		providerStatus = "FAILED"
	case "timeout":
		// Timeout is sweeper-driven; simulate does not replicate the sweeper path.
		// Return a clear error; the TD says refuse or no-op with a log.
		return ProcessCallbackOutput{}, ErrTimeoutNotSupported
	default:
		return ProcessCallbackOutput{}, ErrTimeoutNotSupported
	}

	// Look up intent for the amount (server-side; never trust caller-supplied amount).
	intent, err := s.GetIntentByID(ctx, intentID)
	if err != nil {
		return ProcessCallbackOutput{}, ErrIntentNotFound
	}

	// Generate deterministic mock ref: "MOCKPAY-" + first 8 chars of sha256(intentId+outcome).
	hashHex := computeDedupKey(intentID, outcome)
	mockRef := "MOCKPAY-" + hashHex[:8]

	in := ProcessCallbackInput{
		IntentID:          intentID,
		ProviderStatus:    providerStatus,
		MockPaymentRef:    mockRef,
		AmountFromCaller:  intent.AmountMinor, // server-derived; not from request body
		ProviderTimestamp: time.Now().UTC(),
	}

	deps := processCallbackDeps{
		db:               s.db,
		intentStore:      s.intentStore,
		dedupStore:       s.dedupStore,
		outboxStore:      s.outboxStore,
		emitExpiredEvent: s.emitExpiredEvent,
	}
	return ProcessCallback(ctx, deps, in)
}

func (s *svc) Callback(ctx context.Context, in ProcessCallbackInput) (ProcessCallbackOutput, error) {
	deps := processCallbackDeps{
		db:               s.db,
		intentStore:      s.intentStore,
		dedupStore:       s.dedupStore,
		outboxStore:      s.outboxStore,
		emitExpiredEvent: s.emitExpiredEvent,
	}
	return ProcessCallback(ctx, deps, in)
}
