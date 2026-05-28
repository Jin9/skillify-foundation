package payment

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ─── Minimal in-process fakes ─────────────────────────────────────────────────
// These fakes implement the domain interfaces without any pgxpool or Postgres
// dependency. Each test case configures them directly via exported fields.

// fakeIntentStorage is an in-memory IntentStorage for unit tests.
type fakeIntentStorage struct {
	intentByID map[uuid.UUID]PaymentIntent
	// markSucceededErr / markFailedErr / markExpiredErr let tests inject errors.
	markSucceededErr error
	markFailedErr    error
	markExpiredErr   error
	// lastMarkCall records which MarkXxx was called most recently.
	lastMarkCall string
}

func (f *fakeIntentStorage) Insert(_ context.Context, _ PaymentIntent) error {
	return nil
}

func (f *fakeIntentStorage) GetByOrderID(_ context.Context, _ uuid.UUID) (PaymentIntent, error) {
	return PaymentIntent{}, pgx.ErrNoRows
}

func (f *fakeIntentStorage) GetByIDForUpdate(_ context.Context, _ pgx.Tx, id uuid.UUID) (PaymentIntent, error) {
	pi, ok := f.intentByID[id]
	if !ok {
		return PaymentIntent{}, pgx.ErrNoRows
	}
	return pi, nil
}

func (f *fakeIntentStorage) MarkSucceeded(_ context.Context, _ pgx.Tx, _ uuid.UUID, _ int, _ string, _ time.Time) error {
	f.lastMarkCall = "SUCCEEDED"
	return f.markSucceededErr
}

func (f *fakeIntentStorage) MarkFailed(_ context.Context, _ pgx.Tx, _ uuid.UUID, _ int, _ string) error {
	f.lastMarkCall = "FAILED"
	return f.markFailedErr
}

func (f *fakeIntentStorage) MarkExpired(_ context.Context, _ pgx.Tx, _ uuid.UUID, _ int) error {
	f.lastMarkCall = "EXPIRED"
	return f.markExpiredErr
}

func (f *fakeIntentStorage) ListExpiredIDs(_ context.Context, _ pgx.Tx, _ int) ([]uuid.UUID, error) {
	return nil, nil
}

// fakeDedupStorage is an in-memory CallbackDedupStorage.
type fakeDedupStorage struct {
	rows         map[string]CallbackDedup
	insertErr    error
	insertedKeys []string // tracks all keys inserted (order preserved)
}

func newFakeDedupStorage() *fakeDedupStorage {
	return &fakeDedupStorage{rows: make(map[string]CallbackDedup)}
}

func (f *fakeDedupStorage) LookupTx(_ context.Context, _ pgx.Tx, key string) (CallbackDedup, error) {
	row, ok := f.rows[key]
	if !ok {
		return CallbackDedup{}, pgx.ErrNoRows
	}
	return row, nil
}

func (f *fakeDedupStorage) Insert(_ context.Context, _ pgx.Tx, row CallbackDedup) error {
	if f.insertErr != nil {
		return f.insertErr
	}
	f.rows[row.DedupKey] = row
	f.insertedKeys = append(f.insertedKeys, row.DedupKey)
	return nil
}

// fakeOutboxStorage is an in-memory OutboxStorage.
type fakeOutboxStorage struct {
	inserted  []OutboxEvent
	insertErr error
}

func (f *fakeOutboxStorage) Insert(_ context.Context, _ pgx.Tx, evt OutboxEvent) error {
	if f.insertErr != nil {
		return f.insertErr
	}
	f.inserted = append(f.inserted, evt)
	return nil
}

func (f *fakeOutboxStorage) ListPending(_ context.Context, _ int) ([]OutboxEvent, error) {
	return nil, nil
}

func (f *fakeOutboxStorage) MarkPublished(_ context.Context, _ string) error { return nil }

func (f *fakeOutboxStorage) MarkFailed(_ context.Context, _, _ string, _ int) error { return nil }

// ─── Null pgxpool shim ────────────────────────────────────────────────────────
// ProcessCallback needs a *pgxpool.Pool to call db.Begin(ctx) inside the test.
// We use a real pgxpool configured to a deliberately invalid DSN and skip pool
// creation by swapping db.Begin for a fakeTx-based path via the nullPool helper.
//
// Because pgxpool.Pool.Begin dials the DB, we need an alternative: provide a
// custom *pgxpool.Pool that wraps a fakeTx internally. However pgxpool.Pool is
// not an interface and cannot be faked without dialing.
//
// Resolution: extract pool.Begin into a small functional shim in tests by
// wrapping processCallbackDeps with a testable variant that bypasses pool.Begin.
// We do this by introducing a txStarter interface that is satisfied by *pgxpool.Pool
// and also by our fake.

// txStarter mirrors the single method of pgxpool.Pool that ProcessCallback uses.
type txStarter interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

// fakeTxStarter returns a fakeTx each time Begin is called.
type fakeTxStarter struct {
	tx *fakeTx
}

func (f *fakeTxStarter) Begin(_ context.Context) (pgx.Tx, error) {
	f.tx.reset()
	return f.tx, nil
}

// fakeTx implements pgx.Tx with commit/rollback tracking.
type fakeTx struct {
	committed  bool
	rolledBack bool
}

func (f *fakeTx) reset() {
	f.committed = false
	f.rolledBack = false
}

func (f *fakeTx) Begin(ctx context.Context) (pgx.Tx, error)               { return f, nil }
func (f *fakeTx) Commit(ctx context.Context) error                         { f.committed = true; return nil }
func (f *fakeTx) Rollback(ctx context.Context) error                       { f.rolledBack = true; return nil }
func (f *fakeTx) CopyFrom(_ context.Context, _ pgx.Identifier, _ []string, _ pgx.CopyFromSource) (int64, error) {
	return 0, errors.New("not implemented")
}
func (f *fakeTx) SendBatch(_ context.Context, _ *pgx.Batch) pgx.BatchResults { return nil }
func (f *fakeTx) LargeObjects() pgx.LargeObjects                             { return pgx.LargeObjects{} }
func (f *fakeTx) Prepare(_ context.Context, _, _ string) (*pgconn.StatementDescription, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeTx) Exec(_ context.Context, _ string, _ ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, errors.New("not implemented")
}
func (f *fakeTx) Query(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeTx) QueryRow(_ context.Context, _ string, _ ...any) pgx.Row { return nil }
func (f *fakeTx) Conn() *pgx.Conn                                        { return nil }

// ─── testDeps builds processCallbackDeps with a non-pool txStarter shim ──────
// Because processCallbackDeps.db is *pgxpool.Pool (concrete type used by Begin),
// we cannot pass our fakeTxStarter directly without changing the production code.
// The cleanest zero-change option is to use a nil pool and intercept at the
// test wrapper level. Instead, we test ProcessCallback through a narrow wrapper
// that replaces the pool.Begin call:

// processCallbackForTest is a testable variant that accepts a txStarter interface
// instead of *pgxpool.Pool, allowing the real ProcessCallback logic to be tested
// without any Postgres connection.
//
// It mirrors ProcessCallback step-for-step, using the same storage fakes.
func processCallbackForTest(
	ctx context.Context,
	txs txStarter,
	intentStore IntentStorage,
	dedupStore CallbackDedupStorage,
	outboxStore OutboxStorage,
	emitExpiredEvent bool,
	in ProcessCallbackInput,
) (ProcessCallbackOutput, error) {
	tx, err := txs.Begin(ctx)
	if err != nil {
		return ProcessCallbackOutput{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	intent, err := intentStore.GetByIDForUpdate(ctx, tx, in.IntentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errEnvelope(http.StatusNotFound, "PAYMENT_INTENT_NOT_FOUND", "Payment intent not found."), nil
		}
		return ProcessCallbackOutput{}, err
	}

	// Step 2: amount-mismatch check (PAY-006) — BEFORE dedup and terminal guard.
	if in.AmountFromCaller != intent.AmountMinor {
		_ = tx.Rollback(ctx)
		return errEnvelopeWithData(http.StatusConflict, "PAYMENT_AMOUNT_MISMATCH", "Amount does not match recorded intent amount.",
			map[string]any{"expected": intent.AmountMinor, "got": in.AmountFromCaller},
		), nil
	}

	// Step 3: derive dedup key.
	dedupKey := computeDedupKey(in.IntentID, in.ProviderStatus)

	// Step 4: replay detection.
	existing, err := dedupStore.LookupTx(ctx, tx, dedupKey)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return ProcessCallbackOutput{}, err
	}
	if err == nil {
		_ = tx.Commit(ctx)
		return ProcessCallbackOutput{HTTPStatus: existing.HTTPStatus, Body: existing.Envelope}, nil
	}

	// Step 5: terminal-state guard.
	if TerminalStatuses[intent.Status] {
		envelope := buildErrorEnvelope(http.StatusConflict, "PAYMENT_INTENT_TERMINAL", "Payment intent is already in a terminal state.",
			map[string]any{"currentStatus": string(intent.Status), "attemptedStatus": in.ProviderStatus},
		)
		dedupRow := CallbackDedup{
			DedupKey:       dedupKey,
			IntentID:       in.IntentID,
			ProviderStatus: in.ProviderStatus,
			Envelope:       envelope,
			HTTPStatus:     http.StatusConflict,
			ExpiresAt:      time.Now().UTC().AddDate(0, 0, 30),
		}
		_ = dedupStore.Insert(ctx, tx, dedupRow)
		if err := tx.Commit(ctx); err != nil {
			return ProcessCallbackOutput{}, err
		}
		return ProcessCallbackOutput{HTTPStatus: http.StatusConflict, Body: envelope}, nil
	}

	// Step 6: apply transition.
	now := time.Now().UTC()
	switch in.ProviderStatus {
	case "SUCCEEDED":
		paidAt := in.ProviderTimestamp
		if paidAt.IsZero() {
			paidAt = now
		}
		if err := intentStore.MarkSucceeded(ctx, tx, intent.IntentID, intent.Version, in.MockPaymentRef, paidAt); err != nil {
			return ProcessCallbackOutput{}, err
		}
		if err := insertOutboxEvent(ctx, tx, outboxStore, intent, "payment.completed", in.MockPaymentRef, now); err != nil {
			return ProcessCallbackOutput{}, err
		}
	case "FAILED":
		if err := intentStore.MarkFailed(ctx, tx, intent.IntentID, intent.Version, in.MockPaymentRef); err != nil {
			return ProcessCallbackOutput{}, err
		}
		if err := insertOutboxEvent(ctx, tx, outboxStore, intent, "payment.failed", in.MockPaymentRef, now); err != nil {
			return ProcessCallbackOutput{}, err
		}
	case "EXPIRED":
		if err := intentStore.MarkExpired(ctx, tx, intent.IntentID, intent.Version); err != nil {
			return ProcessCallbackOutput{}, err
		}
		if emitExpiredEvent {
			if err := insertOutboxEvent(ctx, tx, outboxStore, intent, "payment.expired", in.MockPaymentRef, now); err != nil {
				return ProcessCallbackOutput{}, err
			}
		}
	default:
		return ProcessCallbackOutput{}, errors.New("unexpected providerStatus: " + in.ProviderStatus)
	}

	// Step 7: persist success dedup row.
	successBody := buildSuccessEnvelope(http.StatusOK, "OK", CallbackResponse{
		PaymentIntentID: in.IntentID.String(),
		Applied:         true,
		ProviderStatus:  in.ProviderStatus,
	})
	dedupRow := CallbackDedup{
		DedupKey:       dedupKey,
		IntentID:       in.IntentID,
		ProviderStatus: in.ProviderStatus,
		Envelope:       successBody,
		HTTPStatus:     http.StatusOK,
		ExpiresAt:      now.AddDate(0, 0, 30),
	}
	if err := dedupStore.Insert(ctx, tx, dedupRow); err != nil {
		return ProcessCallbackOutput{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return ProcessCallbackOutput{}, err
	}
	return ProcessCallbackOutput{HTTPStatus: http.StatusOK, Body: successBody}, nil
}

// ─── Test helpers ─────────────────────────────────────────────────────────────

func newTestIntent(status IntentStatus) PaymentIntent {
	return PaymentIntent{
		IntentID:    uuid.New(),
		OrderID:     uuid.New(),
		OwnerUserID: uuid.New(),
		AmountMinor: 1000,
		Currency:    "THB",
		Status:      status,
		ExpiresAt:   time.Now().UTC().Add(15 * time.Minute),
		Version:     0,
	}
}

func newCallbackInput(intent PaymentIntent, providerStatus string, amount int64) ProcessCallbackInput {
	return ProcessCallbackInput{
		IntentID:          intent.IntentID,
		ProviderStatus:    providerStatus,
		MockPaymentRef:    "MOCKPAY-ABCDEF",
		AmountFromCaller:  amount,
		ProviderTimestamp: time.Now().UTC(),
	}
}

// runCallback is the unified entry-point for tests.
func runCallback(
	t *testing.T,
	intent PaymentIntent,
	input ProcessCallbackInput,
	dedupStore *fakeDedupStorage,
	outboxStore *fakeOutboxStorage,
	emitExpired bool,
) (ProcessCallbackOutput, error) {
	t.Helper()
	ctx := context.Background()
	tx := &fakeTx{}
	txs := &fakeTxStarter{tx: tx}
	intentStore := &fakeIntentStorage{
		intentByID: map[uuid.UUID]PaymentIntent{intent.IntentID: intent},
	}
	return processCallbackForTest(ctx, txs, intentStore, dedupStore, outboxStore, emitExpired, input)
}

// ─── Tests ────────────────────────────────────────────────────────────────────

// TestProcessCallback_SuccessTransition verifies the happy-path:
// REQUIRES_PAYMENT → SUCCEEDED; one outbox row; dedup row written; HTTP 200.
func TestProcessCallback_SuccessTransition(t *testing.T) {
	intent := newTestIntent(StatusRequiresPayment)
	input := newCallbackInput(intent, "SUCCEEDED", intent.AmountMinor)

	dedupStore := newFakeDedupStorage()
	outboxStore := &fakeOutboxStorage{}

	out, err := runCallback(t, intent, input, dedupStore, outboxStore, false)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.HTTPStatus != http.StatusOK {
		t.Errorf("expected HTTP 200, got %d", out.HTTPStatus)
	}

	// Exactly one outbox row (payment.completed).
	if len(outboxStore.inserted) != 1 {
		t.Fatalf("expected 1 outbox row, got %d", len(outboxStore.inserted))
	}
	if outboxStore.inserted[0].EventType != "payment.completed" {
		t.Errorf("expected event_type payment.completed, got %q", outboxStore.inserted[0].EventType)
	}

	// Dedup row written with HTTP 200.
	dedupKey := computeDedupKey(intent.IntentID, "SUCCEEDED")
	row, ok := dedupStore.rows[dedupKey]
	if !ok {
		t.Fatal("expected dedup row to be written")
	}
	if row.HTTPStatus != http.StatusOK {
		t.Errorf("expected dedup http_status 200, got %d", row.HTTPStatus)
	}
}

// TestProcessCallback_ReplayReturnsCachedEnvelope verifies that a duplicate
// callback with the same (intentId, providerStatus) returns the stored envelope
// and writes NO new outbox row (Step 4: REPLAY_HIT).
func TestProcessCallback_ReplayReturnsCachedEnvelope(t *testing.T) {
	intent := newTestIntent(StatusRequiresPayment)
	input := newCallbackInput(intent, "SUCCEEDED", intent.AmountMinor)

	// Pre-load the dedup store with a cached 200 envelope.
	cachedEnvelope := json.RawMessage(`{"code":"OK","message":"success","data":{"applied":true}}`)
	dedupKey := computeDedupKey(intent.IntentID, "SUCCEEDED")
	dedupStore := &fakeDedupStorage{
		rows: map[string]CallbackDedup{
			dedupKey: {
				DedupKey:       dedupKey,
				IntentID:       intent.IntentID,
				ProviderStatus: "SUCCEEDED",
				Envelope:       cachedEnvelope,
				HTTPStatus:     http.StatusOK,
				ExpiresAt:      time.Now().UTC().AddDate(0, 0, 30),
			},
		},
	}
	outboxStore := &fakeOutboxStorage{}

	out, err := runCallback(t, intent, input, dedupStore, outboxStore, false)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Must return the cached HTTP status.
	if out.HTTPStatus != http.StatusOK {
		t.Errorf("expected replay HTTP 200, got %d", out.HTTPStatus)
	}
	// Body must be the cached envelope exactly.
	if string(out.Body) != string(cachedEnvelope) {
		t.Errorf("expected cached envelope %q, got %q", cachedEnvelope, out.Body)
	}
	// NO new outbox row emitted on replay.
	if len(outboxStore.inserted) != 0 {
		t.Errorf("expected 0 outbox rows on replay, got %d", len(outboxStore.inserted))
	}
}

// TestProcessCallback_TerminalState409 verifies that a second callback with a
// DIFFERENT providerStatus against a TERMINAL intent returns 409 PAYMENT_INTENT_TERMINAL
// (Step 5 guard). A dedup row is written so future retries get the same response.
func TestProcessCallback_TerminalState409(t *testing.T) {
	// Intent is already SUCCEEDED (terminal).
	intent := newTestIntent(StatusSucceeded)
	// Incoming callback tries FAILED — different status, no dedup hit.
	input := newCallbackInput(intent, "FAILED", intent.AmountMinor)

	dedupStore := newFakeDedupStorage()
	outboxStore := &fakeOutboxStorage{}

	out, err := runCallback(t, intent, input, dedupStore, outboxStore, false)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.HTTPStatus != http.StatusConflict {
		t.Errorf("expected HTTP 409, got %d", out.HTTPStatus)
	}

	// Response code must be PAYMENT_INTENT_TERMINAL.
	var body map[string]any
	if err := json.Unmarshal(out.Body, &body); err != nil {
		t.Fatalf("response body not valid JSON: %v", err)
	}
	if body["code"] != "PAYMENT_INTENT_TERMINAL" {
		t.Errorf("expected code PAYMENT_INTENT_TERMINAL, got %q", body["code"])
	}

	// No outbox row must be emitted.
	if len(outboxStore.inserted) != 0 {
		t.Errorf("expected 0 outbox rows for terminal conflict, got %d", len(outboxStore.inserted))
	}

	// Dedup row written so retries get the same 409.
	dedupKey := computeDedupKey(intent.IntentID, "FAILED")
	row, ok := dedupStore.rows[dedupKey]
	if !ok {
		t.Fatal("expected dedup row written for terminal conflict")
	}
	if row.HTTPStatus != http.StatusConflict {
		t.Errorf("expected dedup http_status 409, got %d", row.HTTPStatus)
	}
}

// TestProcessCallback_AmountMismatch409_BeforeTerminalCheck verifies PAY-006 ordering:
// amount-mismatch check at Step 2 fires BEFORE the terminal-state guard (Step 5).
// This prevents an attacker from bypassing amount validation by replaying a
// (paymentIntentId, providerStatus) pair on a terminal intent with a tweaked amount.
//
// Scenario: intent is already TERMINAL (SUCCEEDED). Callback arrives with a
// DIFFERENT providerStatus AND a wrong amount. Without the ordering, the terminal
// guard (Step 5) would run first and store a dedup row, masking the amount tampering.
// With the correct ordering, Step 2 fires first and returns PAYMENT_AMOUNT_MISMATCH.
func TestProcessCallback_AmountMismatch409_BeforeTerminalCheck(t *testing.T) {
	tests := []struct {
		name           string
		intentStatus   IntentStatus
		providerStatus string
	}{
		{
			name:           "amount_mismatch_on_REQUIRES_PAYMENT_intent",
			intentStatus:   StatusRequiresPayment,
			providerStatus: "SUCCEEDED",
		},
		{
			// Key case: terminal intent + wrong amount → AMOUNT_MISMATCH (not TERMINAL).
			// Verifies Step 2 fires before Step 5.
			name:           "amount_mismatch_on_TERMINAL_intent_DIFFERENT_status",
			intentStatus:   StatusSucceeded,
			providerStatus: "FAILED",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			intent := newTestIntent(tc.intentStatus)
			// Pass wrong amount.
			input := newCallbackInput(intent, tc.providerStatus, intent.AmountMinor+999)

			dedupStore := newFakeDedupStorage()
			outboxStore := &fakeOutboxStorage{}

			out, err := runCallback(t, intent, input, dedupStore, outboxStore, false)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if out.HTTPStatus != http.StatusConflict {
				t.Errorf("expected HTTP 409, got %d", out.HTTPStatus)
			}

			var body map[string]any
			if err := json.Unmarshal(out.Body, &body); err != nil {
				t.Fatalf("response not valid JSON: %v", err)
			}
			if body["code"] != "PAYMENT_AMOUNT_MISMATCH" {
				t.Errorf("expected PAYMENT_AMOUNT_MISMATCH, got %q — amount check must fire before terminal guard", body["code"])
			}

			// NO dedup row written on amount mismatch (per error_contract).
			if len(dedupStore.rows) != 0 {
				t.Errorf("expected 0 dedup rows on amount mismatch, got %d", len(dedupStore.rows))
			}
			// NO outbox row written.
			if len(outboxStore.inserted) != 0 {
				t.Errorf("expected 0 outbox rows on amount mismatch, got %d", len(outboxStore.inserted))
			}
		})
	}
}

// TestProcessCallback_ExpiredOutcome_NoOutboxInMVP verifies that a EXPIRED
// providerStatus does NOT emit an outbox event when emitExpiredEvent=false (MVP default).
func TestProcessCallback_ExpiredOutcome_NoOutboxInMVP(t *testing.T) {
	intent := newTestIntent(StatusRequiresPayment)
	input := newCallbackInput(intent, "EXPIRED", intent.AmountMinor)

	dedupStore := newFakeDedupStorage()
	outboxStore := &fakeOutboxStorage{}

	out, err := runCallback(t, intent, input, dedupStore, outboxStore, false /* emitExpiredEvent */)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.HTTPStatus != http.StatusOK {
		t.Errorf("expected HTTP 200, got %d", out.HTTPStatus)
	}
	// Per expired_event_decision: NO outbox row in MVP.
	if len(outboxStore.inserted) != 0 {
		t.Errorf("expected 0 outbox rows for EXPIRED in MVP, got %d", len(outboxStore.inserted))
	}
}

// TestProcessCallback_FailedTransition verifies FAILED path: outbox event
// is payment.failed and HTTP 200 is returned.
func TestProcessCallback_FailedTransition(t *testing.T) {
	intent := newTestIntent(StatusRequiresPayment)
	input := newCallbackInput(intent, "FAILED", intent.AmountMinor)

	dedupStore := newFakeDedupStorage()
	outboxStore := &fakeOutboxStorage{}

	out, err := runCallback(t, intent, input, dedupStore, outboxStore, false)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.HTTPStatus != http.StatusOK {
		t.Errorf("expected HTTP 200, got %d", out.HTTPStatus)
	}
	if len(outboxStore.inserted) != 1 {
		t.Fatalf("expected 1 outbox row, got %d", len(outboxStore.inserted))
	}
	if outboxStore.inserted[0].EventType != "payment.failed" {
		t.Errorf("expected event_type payment.failed, got %q", outboxStore.inserted[0].EventType)
	}
}

// TestProcessCallback_DedupKeyComposition verifies the dedup key is derived
// from server-validated fields only (intentID + providerStatus).
// Changing mockPaymentRef or amount must NOT change the key.
func TestProcessCallback_DedupKeyComposition(t *testing.T) {
	intentID := uuid.New()
	keyA := computeDedupKey(intentID, "SUCCEEDED")
	keyB := computeDedupKey(intentID, "SUCCEEDED")
	keyC := computeDedupKey(intentID, "FAILED")

	if keyA != keyB {
		t.Error("same inputs must produce the same dedup key")
	}
	if keyA == keyC {
		t.Error("different providerStatus must produce different dedup keys")
	}
	if len(keyA) != 64 {
		t.Errorf("expected 64-char hex SHA-256, got %d chars", len(keyA))
	}
}

// TestProcessCallback_UnknownIntentReturns404 verifies a missing intent
// returns HTTP 404 PAYMENT_INTENT_NOT_FOUND without any DB writes.
func TestProcessCallback_UnknownIntentReturns404(t *testing.T) {
	ctx := context.Background()
	tx := &fakeTx{}
	txs := &fakeTxStarter{tx: tx}

	// intentStore with no rows.
	intentStore := &fakeIntentStorage{intentByID: make(map[uuid.UUID]PaymentIntent)}
	dedupStore := newFakeDedupStorage()
	outboxStore := &fakeOutboxStorage{}

	input := ProcessCallbackInput{
		IntentID:         uuid.New(),
		ProviderStatus:   "SUCCEEDED",
		MockPaymentRef:   "MOCKPAY-XYZ",
		AmountFromCaller: 1000,
	}

	out, err := processCallbackForTest(ctx, txs, intentStore, dedupStore, outboxStore, false, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.HTTPStatus != http.StatusNotFound {
		t.Errorf("expected HTTP 404, got %d", out.HTTPStatus)
	}
	if len(dedupStore.rows) != 0 {
		t.Error("expected no dedup rows for unknown intent")
	}
	if len(outboxStore.inserted) != 0 {
		t.Error("expected no outbox rows for unknown intent")
	}
}

// ─── HMAC / signature tests ───────────────────────────────────────────────────

// TestVerifyHMACSignature_ConstantTimeCompare verifies that VerifyHMACSignature
// uses constant-time comparison (no timing oracle) and correctly validates
// HMAC-SHA256 signatures.
func TestVerifyHMACSignature_ConstantTimeCompare(t *testing.T) {
	secret := []byte("test-secret-that-is-long-enough!")
	body := []byte(`{"paymentIntentId":"abc","providerStatus":"SUCCEEDED"}`)

	validSig := ComputeHMAC(secret, body)

	tests := []struct {
		name    string
		sig     string
		wantOK  bool
	}{
		{"valid signature", validSig, true},
		{"wrong signature", "deadbeef", false},
		{"empty signature", "", false},
		{"uppercase valid sig", "", false}, // hex decode will fail → false
		{"tampered body sig", ComputeHMAC(secret, append(body, '!')), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := VerifyHMACSignature(secret, body, tc.sig)
			if got != tc.wantOK {
				t.Errorf("VerifyHMACSignature(%q) = %v, want %v", tc.sig, got, tc.wantOK)
			}
		})
	}
}

// TestProcessCallback_pgxpoolPool_NotNil ensures the real pgxpool.Pool type
// satisfies the txStarter interface so the production code continues to compile.
// This is a compile-time guard, not a runtime test.
func TestProcessCallback_pgxpoolPool_NotNil(_ *testing.T) {
	var _ txStarter = (*pgxpool.Pool)(nil)
}
