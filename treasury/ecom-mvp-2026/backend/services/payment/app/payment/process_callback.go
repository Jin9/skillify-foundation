package payment

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// processCallbackDeps are the storage dependencies used by processCallback.
// Injected by the service to keep the function pure and testable.
type processCallbackDeps struct {
	db               *pgxpool.Pool
	intentStore      IntentStorage
	dedupStore       CallbackDedupStorage
	outboxStore      OutboxStorage
	emitExpiredEvent bool
}

// ProcessCallback is the canonical 7-step handler for all payment outcomes.
// It is invoked by:
//   - handler_callback.go (HTTP webhook after HMAC verification)
//   - handler_simulate.go (in-process, skipping the HTTP hop)
//
// 7-step ordering (per td.json internal_handlers.process_callback):
//  1. Begin tx + SELECT payment_intents FOR UPDATE
//  2. Amount check (PAY-006): reject before deriving dedup key
//  3. Derive dedup_key = sha256(intentId || '|' || providerStatus)
//  4. Dedup lookup: if found, return cached envelope without re-firing outbox
//  5. Terminal-state guard: 409 if already terminal with a different providerStatus
//  6. Apply transition: UPDATE intent + INSERT outbox (SUCCEEDED/FAILED only)
//  7. INSERT dedup row with success envelope; COMMIT
//
// SECURITY invariants:
//   - Dedup key composed ONLY from server-validated fields (intentId + providerStatus).
//   - Amount compared against payment_intents.amount_minor (server truth), never trusted alone.
//   - Full callback bodies are NEVER logged; only dedup_key + intentId + status.
func ProcessCallback(ctx context.Context, deps processCallbackDeps, in ProcessCallbackInput) (ProcessCallbackOutput, error) {
	// ── Step 1: open tx, lock intent row ────────────────────────────────────
	tx, err := deps.db.Begin(ctx)
	if err != nil {
		return ProcessCallbackOutput{}, fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		// Best-effort rollback; ignored if already committed.
		_ = tx.Rollback(ctx)
	}()

	intent, err := deps.intentStore.GetByIDForUpdate(ctx, tx, in.IntentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errEnvelope(http.StatusNotFound, "PAYMENT_INTENT_NOT_FOUND", "Payment intent not found."), nil
		}
		return ProcessCallbackOutput{}, fmt.Errorf("lock intent: %w", err)
	}

	// ── Step 2: amount-mismatch check (PAY-006) ──────────────────────────────
	// Do NOT derive dedupKey yet — do NOT insert a dedup row for amount errors
	// per error_contract note: "no dedup-row write".
	if in.AmountFromCaller != intent.AmountMinor {
		slog.WarnContext(ctx, "amount mismatch",
			"intent_id", intent.IntentID,
			"expected", intent.AmountMinor,
			"got", in.AmountFromCaller,
		)
		_ = tx.Rollback(ctx)
		return errEnvelopeWithData(http.StatusConflict, "PAYMENT_AMOUNT_MISMATCH", "Amount does not match recorded intent amount.",
			map[string]any{
				"expected": intent.AmountMinor,
				"got":      in.AmountFromCaller,
			},
		), nil
	}

	// ── Step 3: derive dedup_key from server-validated fields only ───────────
	dedupKey := computeDedupKey(in.IntentID, in.ProviderStatus)

	// ── Step 4: replay detection (inside tx for serialisation) ──────────────
	existing, err := deps.dedupStore.LookupTx(ctx, tx, dedupKey)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return ProcessCallbackOutput{}, fmt.Errorf("dedup lookup: %w", err)
	}
	if err == nil {
		// REPLAY_HIT: return cached envelope; never re-fire outbox.
		slog.InfoContext(ctx, "dedup replay hit",
			"dedup_key", dedupKey,
			"intent_id", in.IntentID,
			"provider_status", in.ProviderStatus,
			"action", "REPLAY_HIT",
		)
		_ = tx.Commit(ctx) // no writes; commit is a no-op but keeps tx lifecycle clean
		return ProcessCallbackOutput{
			HTTPStatus: existing.HTTPStatus,
			Body:       existing.Envelope,
		}, nil
	}

	// ── Step 5: terminal-state guard ─────────────────────────────────────────
	// Same-status replay would have collided on dedup key at Step 4.
	// Reaching here with a terminal intent means a different providerStatus was attempted.
	if TerminalStatuses[intent.Status] {
		slog.InfoContext(ctx, "intent already terminal",
			"dedup_key", dedupKey,
			"intent_id", in.IntentID,
			"current_status", intent.Status,
			"attempted_status", in.ProviderStatus,
			"action", "TERMINAL",
		)
		envelope := buildErrorEnvelope(http.StatusConflict, "PAYMENT_INTENT_TERMINAL", "Payment intent is already in a terminal state.",
			map[string]any{
				"currentStatus":   string(intent.Status),
				"attemptedStatus": in.ProviderStatus,
			},
		)
		// Persist dedup row so retries return the same 409 without re-entering the guard.
		dedupRow := CallbackDedup{
			DedupKey:       dedupKey,
			IntentID:       in.IntentID,
			ProviderStatus: in.ProviderStatus,
			Envelope:       envelope,
			HTTPStatus:     http.StatusConflict,
			ExpiresAt:      time.Now().UTC().AddDate(0, 0, 30),
		}
		if insertErr := deps.dedupStore.Insert(ctx, tx, dedupRow); insertErr != nil {
			slog.ErrorContext(ctx, "failed to insert terminal dedup row", "error", insertErr, "dedup_key", dedupKey)
			// Non-fatal: still return the 409; dedup row will be missing but that's
			// a correctness degradation only (repeat calls re-enter guard).
		}
		if err := tx.Commit(ctx); err != nil {
			return ProcessCallbackOutput{}, fmt.Errorf("commit terminal dedup: %w", err)
		}
		return ProcessCallbackOutput{
			HTTPStatus: http.StatusConflict,
			Body:       envelope,
		}, nil
	}

	// ── Step 6: apply transition ─────────────────────────────────────────────
	// At this point intent.Status MUST be REQUIRES_PAYMENT.
	providerRef := in.MockPaymentRef
	now := time.Now().UTC()

	switch in.ProviderStatus {
	case "SUCCEEDED":
		paidAt := in.ProviderTimestamp
		if paidAt.IsZero() {
			paidAt = now
		}
		if err := deps.intentStore.MarkSucceeded(ctx, tx, intent.IntentID, intent.Version, providerRef, paidAt); err != nil {
			return ProcessCallbackOutput{}, fmt.Errorf("mark succeeded: %w", err)
		}
		// Emit outbox event: payment.completed
		if err := insertOutboxEvent(ctx, tx, deps.outboxStore, intent, "payment.completed", providerRef, now); err != nil {
			return ProcessCallbackOutput{}, fmt.Errorf("outbox insert (completed): %w", err)
		}

	case "FAILED":
		if err := deps.intentStore.MarkFailed(ctx, tx, intent.IntentID, intent.Version, providerRef); err != nil {
			return ProcessCallbackOutput{}, fmt.Errorf("mark failed: %w", err)
		}
		// Emit outbox event: payment.failed
		if err := insertOutboxEvent(ctx, tx, deps.outboxStore, intent, "payment.failed", providerRef, now); err != nil {
			return ProcessCallbackOutput{}, fmt.Errorf("outbox insert (failed): %w", err)
		}

	case "EXPIRED":
		if err := deps.intentStore.MarkExpired(ctx, tx, intent.IntentID, intent.Version); err != nil {
			return ProcessCallbackOutput{}, fmt.Errorf("mark expired: %w", err)
		}
		// Per expired_event_decision: DO NOT emit events.payment.expired in MVP.
		// Feature flag PAYMENT_EMIT_EXPIRED_EVENT (default false) gates future enablement.
		if deps.emitExpiredEvent {
			if err := insertOutboxEvent(ctx, tx, deps.outboxStore, intent, "payment.expired", providerRef, now); err != nil {
				return ProcessCallbackOutput{}, fmt.Errorf("outbox insert (expired): %w", err)
			}
		}

	default:
		return ProcessCallbackOutput{}, fmt.Errorf("unexpected providerStatus: %s", in.ProviderStatus)
	}

	// ── Step 7: persist success dedup row ────────────────────────────────────
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
	if err := deps.dedupStore.Insert(ctx, tx, dedupRow); err != nil {
		return ProcessCallbackOutput{}, fmt.Errorf("dedup insert (success): %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return ProcessCallbackOutput{}, fmt.Errorf("commit: %w", err)
	}

	slog.InfoContext(ctx, "callback applied",
		"dedup_key", dedupKey,
		"intent_id", in.IntentID,
		"order_id", intent.OrderID,
		"new_status", in.ProviderStatus,
		"action", "APPLIED",
	)

	return ProcessCallbackOutput{
		HTTPStatus: http.StatusOK,
		Body:       successBody,
	}, nil
}

// ─── helpers ─────────────────────────────────────────────────────────────────

// computeDedupKey returns sha256(intentID || '|' || providerStatus) as hex.
// ONLY server-validated fields are used; no user-controlled body content.
func computeDedupKey(intentID uuid.UUID, providerStatus string) string {
	h := sha256.Sum256([]byte(intentID.String() + "|" + providerStatus))
	return fmt.Sprintf("%x", h)
}

func insertOutboxEvent(ctx context.Context, tx pgx.Tx, store OutboxStorage, intent PaymentIntent, eventType, providerRef string, now time.Time) error {
	eventID := uuid.New()

	var rawPayload []byte
	var err error

	switch eventType {
	case "payment.completed":
		rawPayload, err = json.Marshal(PaymentCompletedPayload{
			EventID:         eventID.String(),
			OccurredAt:      now,
			OrderID:         intent.OrderID.String(),
			PaymentIntentID: intent.IntentID.String(),
			Amount:          intent.AmountMinor,
			MockPaymentRef:  providerRef,
		})
	case "payment.failed":
		rawPayload, err = json.Marshal(PaymentFailedPayload{
			EventID:         eventID.String(),
			OccurredAt:      now,
			OrderID:         intent.OrderID.String(),
			PaymentIntentID: intent.IntentID.String(),
			Reason:          "PROVIDER_DECLINED",
		})
	case "payment.expired":
		// Future: same shape; reserved for feature flag path.
		rawPayload, err = json.Marshal(map[string]any{
			"eventId":         eventID.String(),
			"occurredAt":      now,
			"orderId":         intent.OrderID.String(),
			"paymentIntentId": intent.IntentID.String(),
		})
	}
	if err != nil {
		return fmt.Errorf("marshal outbox payload: %w", err)
	}

	evt := OutboxEvent{
		ID:            uuid.New(),
		AggregateType: "payment_intent",
		AggregateID:   intent.IntentID,
		EventID:       eventID,
		EventType:     eventType,
		Topic:         "ecom.payment.events",
		PartitionKey:  intent.OrderID.String(), // per-order ordering
		Payload:       rawPayload,
	}
	return store.Insert(ctx, tx, evt)
}

// envelope helpers — produce the raw JSON body stored in dedup rows and returned to callers.

func buildSuccessEnvelope(httpStatus int, code string, data any) json.RawMessage {
	b, _ := json.Marshal(map[string]any{
		"code":    code,
		"message": "success",
		"data":    data,
	})
	return b
}

func buildErrorEnvelope(httpStatus int, code, message string, data any) json.RawMessage {
	b, _ := json.Marshal(map[string]any{
		"code":    code,
		"message": message,
		"data":    data,
	})
	return b
}

// errEnvelope wraps a non-dedup error as a ProcessCallbackOutput.
func errEnvelope(httpStatus int, code, message string) ProcessCallbackOutput {
	body := buildErrorEnvelope(httpStatus, code, message, nil)
	return ProcessCallbackOutput{HTTPStatus: httpStatus, Body: body}
}

func errEnvelopeWithData(httpStatus int, code, message string, data any) ProcessCallbackOutput {
	body := buildErrorEnvelope(httpStatus, code, message, data)
	return ProcessCallbackOutput{HTTPStatus: httpStatus, Body: body}
}
