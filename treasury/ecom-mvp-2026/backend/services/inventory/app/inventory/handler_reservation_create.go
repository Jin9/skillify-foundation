package inventory

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/middleware"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

const (
	idemEndpointReservationCreate = "inventory.reservation.create"
)

// ReservationCreateHandler is the FULL handler for POST /api/v1/inventory/reservation/create.
//
// Lock-order pin (TD spec concurrency_model.lock_ordering):
//
//	Within the tx: lock stock_levels rows in lexicographic SKU order FIRST,
//	then INSERT reservation rows. This prevents deadlocks under concurrent
//	multi-SKU checkouts.
//
// Idempotency: keyed on (orderId, endpoint). Same orderId + same items → replay.
// Same orderId + different items → 409 IDEMPOTENCY_KEY_REUSED.
func (svc *Service) ReservationCreateHandler(c *gin.Context) {
	traceID := middleware.RefID(c)
	ctx := c.Request.Context()

	req, ok := wrapper.BindJSON[ReservationCreateRequest](c)
	if !ok {
		return
	}

	// Validate idempotencyKey == orderId per contract.
	if req.IdempotencyKey != req.OrderID {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusBadRequest,
			Code:       CodeValidationError,
			Message:    "idempotencyKey must equal orderId",
			TraceID:    traceID,
		})
		return
	}

	if len(req.Items) == 0 {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusBadRequest,
			Code:       CodeValidationError,
			Message:    "items must be non-empty",
			TraceID:    traceID,
		})
		return
	}

	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusBadRequest,
			Code:       CodeValidationError,
			Message:    "orderId must be a valid UUID",
			TraceID:    traceID,
		})
		return
	}

	// Compute canonical request hash: sort items by sku first.
	sortedItems := make([]ReservationItem, len(req.Items))
	copy(sortedItems, req.Items)
	sort.Slice(sortedItems, func(i, j int) bool {
		return sortedItems[i].SKU < sortedItems[j].SKU
	})
	requestHash := canonicalHash(req.OrderID, sortedItems)

	// Determine expiresAt.
	expiresAt := req.ExpiresAt
	if expiresAt.IsZero() {
		expiresAt = time.Now().UTC().Add(
			time.Duration(svc.Cfg.Sweeper.ReservationTTLMinutes) * time.Minute,
		)
	}

	// Deduplicate unique SKUs for the lock scan; maintain lex order.
	skuSet := make(map[string]struct{}, len(sortedItems))
	for _, it := range sortedItems {
		skuSet[it.SKU] = struct{}{}
	}
	sortedSKUs := make([]string, 0, len(skuSet))
	for sku := range skuSet {
		sortedSKUs = append(sortedSKUs, sku)
	}
	sort.Strings(sortedSKUs) // lexicographic order — LOCK-ORDER PIN

	// Begin transaction.
	tx, err := svc.Pool.Begin(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "reservation.create begin tx", slog.String("error", err.Error()))
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusServiceUnavailable,
			Code:       CodeDatabaseUnavailable,
			Message:    MessageDatabaseUnavailable,
			TraceID:    traceID,
		})
		return
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// ── Idempotency check (FOR UPDATE on idempotency_keys serializes concurrent same-key calls) ──
	existing, err := idemFetch(ctx, tx, req.OrderID, idemEndpointReservationCreate)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.ErrorContext(ctx, "reservation.create idem fetch", slog.String("error", err.Error()))
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusServiceUnavailable,
			Code:       CodeDatabaseUnavailable,
			Message:    MessageDatabaseUnavailable,
			TraceID:    traceID,
		})
		return
	}

	if existing != nil {
		if existing.RequestHash == requestHash {
			// Idempotent replay: return cached envelope verbatim.
			_ = tx.Rollback(ctx)
			c.Data(existing.HTTPStatus, "application/json", existing.ResponseEnvelope)
			return
		}
		// Same key, different payload → conflict.
		_ = tx.Rollback(ctx)
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusConflict,
			Code:       CodeIdempotencyKeyReused,
			Message:    MessageIdempotencyKeyReused,
			TraceID:    traceID,
		})
		return
	}

	// ── LOCK-ORDER PIN: lock stock_levels in lexicographic SKU order ──────────
	// "ALWAYS lock stock_levels(sku) BEFORE reservations(id) inside any
	//  state-mutating tx. Within a single multi-SKU reservation.create call,
	//  lock SKUs in lexicographic order (ORDER BY sku)." — td.json concurrency_model.lock_ordering
	stockRows, err := svc.StockStorage.GetManyBySKUsForUpdate(ctx, tx, sortedSKUs)
	if err != nil {
		slog.ErrorContext(ctx, "reservation.create lock stock", slog.String("error", err.Error()))
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusServiceUnavailable,
			Code:       CodeDatabaseUnavailable,
			Message:    MessageDatabaseUnavailable,
			TraceID:    traceID,
		})
		return
	}

	// ── Per-SKU availability check ────────────────────────────────────────────
	var insufficientList []InsufficientDetail
	for _, item := range sortedItems {
		sl, ok := stockRows[item.SKU]
		if !ok {
			// SKU not found — treat as INSUFFICIENT_STOCK (no stock at all).
			insufficientList = append(insufficientList, InsufficientDetail{
				SKU:          item.SKU,
				RequestedQty: item.Qty,
				AvailableQty: 0,
			})
			continue
		}
		if sl.AvailableQty < item.Qty {
			insufficientList = append(insufficientList, InsufficientDetail{
				SKU:          item.SKU,
				RequestedQty: item.Qty,
				AvailableQty: sl.AvailableQty,
			})
		}
	}

	if len(insufficientList) > 0 {
		// Rollback entire tx and return 409 with per-SKU detail.
		_ = tx.Rollback(ctx)
		type insufficientBody struct {
			InsufficientItems []InsufficientDetail `json:"insufficientItems"`
		}
		body := insufficientBody{InsufficientItems: insufficientList}
		wrapper.Respond(c, wrapper.ResponseOption[insufficientBody]{
			HTTPStatus: http.StatusConflict,
			Code:       CodeInsufficientStock,
			Message:    MessageInsufficientStock,
			Data:       &body,
			TraceID:    traceID,
		})
		return
	}

	// ── Mutate: UPDATE stock + INSERT reservations (still inside same tx) ─────
	// Items are already in lex order (sortedItems). The stock rows are already locked above.
	var reservationIDs []uuid.UUID
	for _, item := range sortedItems {
		// No-negative invariant enforced:
		//   (a) tx-level: verified available_qty >= qty in the check above (on the locked row).
		//   (b) DB-level: CHECK (available_qty >= 0) constraint fires if any race slips past.
		if err := svc.StockStorage.DecrementAvailableIncrementReserved(ctx, tx, item.SKU, item.Qty); err != nil {
			slog.ErrorContext(ctx, "reservation.create update stock",
				slog.String("sku", item.SKU),
				slog.String("error", err.Error()),
			)
			// DB CHECK constraint firing maps to STOCK_NEGATIVE_INVARIANT (defensive; unreachable under normal flow).
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusConflict,
				Code:       CodeStockNegativeInvariant,
				Message:    MessageStockNegativeInvariant,
				TraceID:    traceID,
				Err:        serror.Wrap(err),
			})
			return
		}

		rid := uuid.New()
		resv := &Reservation{
			ID:        rid,
			OrderID:   orderID,
			SKU:       item.SKU,
			Qty:       item.Qty,
			Status:    StatusReserved,
			ExpiresAt: expiresAt,
		}
		if err := svc.ResvStorage.Insert(ctx, tx, resv); err != nil {
			slog.ErrorContext(ctx, "reservation.create insert reservation",
				slog.String("sku", item.SKU),
				slog.String("error", err.Error()),
			)
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusServiceUnavailable,
				Code:       CodeDatabaseUnavailable,
				Message:    MessageDatabaseUnavailable,
				TraceID:    traceID,
				Err:        serror.Wrap(err),
			})
			return
		}
		reservationIDs = append(reservationIDs, rid)
	}

	// ── Build response envelope and persist idempotency key atomically ────────
	respItems := make([]ReservationCreateResponseItem, len(sortedItems))
	for i, item := range sortedItems {
		respItems[i] = ReservationCreateResponseItem{
			SKU:      item.SKU,
			Qty:      item.Qty,
			Reserved: true,
		}
	}
	responseData := ReservationCreateResponse{
		ReservationID: reservationIDs[0].String(),
		Items:         respItems,
		ExpiresAt:     expiresAt,
	}

	envelopeBody, err := json.Marshal(map[string]any{
		"code":    string(CodeSuccess),
		"message": "reserved",
		"data":    responseData,
		"traceId": traceID,
	})
	if err != nil {
		slog.ErrorContext(ctx, "reservation.create marshal envelope", slog.String("error", err.Error()))
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       CodeInternalError,
			Message:    MessageInternalError,
			TraceID:    traceID,
		})
		return
	}

	if err := idemInsert(ctx, tx, IdempotencyKey{
		Key:              req.OrderID,
		Endpoint:         idemEndpointReservationCreate,
		RequestHash:      requestHash,
		ResponseEnvelope: envelopeBody,
		HTTPStatus:       http.StatusOK,
		ExpiresAt:        time.Now().UTC().Add(24 * time.Hour),
	}); err != nil {
		slog.ErrorContext(ctx, "reservation.create persist idem key", slog.String("error", err.Error()))
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusServiceUnavailable,
			Code:       CodeDatabaseUnavailable,
			Message:    MessageDatabaseUnavailable,
			TraceID:    traceID,
		})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		slog.ErrorContext(ctx, "reservation.create commit", slog.String("error", err.Error()))
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusServiceUnavailable,
			Code:       CodeDatabaseUnavailable,
			Message:    MessageDatabaseUnavailable,
			TraceID:    traceID,
		})
		return
	}

	wrapper.Respond(c, wrapper.ResponseOption[ReservationCreateResponse]{
		HTTPStatus: http.StatusOK,
		Code:       CodeSuccess,
		Message:    "reserved",
		Data:       &responseData,
		TraceID:    traceID,
	})
}

// ReservationReleaseHandler is the FULL handler for POST /api/v1/inventory/reservation/release.
//
// Lock-order pin (TD spec concurrency_model.lock_ordering):
//
//	ALWAYS lock stock_levels(sku) BEFORE reservations(id) in any state-mutating tx.
//	Implementation (Option A two-phase read):
//	  Phase 1 — Read reservations without FOR UPDATE as a snapshot (FindByOrderID).
//	            This gives us the ordered SKU set without holding any row-lock.
//	  Phase 2 — Lock stock_levels in lexicographic SKU order (GetBySKUForUpdate).
//	            Then immediately re-lock the specific reservation row by PK
//	            (GetByIDForUpdate) before mutating it.
//	This matches the sweeper's lock order and eliminates the cross-path deadlock
//	between concurrent sweep + payment.completed / release on the same orderId.
func (svc *Service) ReservationReleaseHandler(c *gin.Context) {
	traceID := middleware.RefID(c)
	ctx := c.Request.Context()

	req, ok := wrapper.BindJSON[ReservationReleaseRequest](c)
	if !ok {
		return
	}

	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusBadRequest,
			Code:       CodeValidationError,
			Message:    "orderId must be a valid UUID",
			TraceID:    traceID,
		})
		return
	}

	releaseReason := ReleaseReason(req.Reason)

	tx, err := svc.Pool.Begin(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "reservation.release begin tx", slog.String("error", err.Error()))
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusServiceUnavailable,
			Code:       CodeDatabaseUnavailable,
			Message:    MessageDatabaseUnavailable,
			TraceID:    traceID,
		})
		return
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// ── Phase 1: snapshot read (no lock) ────────────────────────────────────────
	// LOCK-ORDER PIN (Option A): read reservations WITHOUT FOR UPDATE first so we
	// know the SKU set. No row-locks acquired here. This prevents the inversion
	// where reservations were locked before stock_levels.
	rows, err := svc.ResvStorage.FindByOrderID(ctx, tx, orderID)
	if err != nil {
		slog.ErrorContext(ctx, "reservation.release snapshot reservations", slog.String("error", err.Error()))
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusServiceUnavailable,
			Code:       CodeDatabaseUnavailable,
			Message:    MessageDatabaseUnavailable,
			TraceID:    traceID,
		})
		return
	}
	if len(rows) == 0 {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusNotFound,
			Code:       CodeReservationNotFound,
			Message:    MessageReservationNotFound,
			TraceID:    traceID,
		})
		return
	}

	// ── Phase 2: lock stock_levels first, then re-lock reservation by PK ────────
	// rows are already ORDER BY sku ASC from the snapshot query. Iterate in that
	// order so stock_levels locks are acquired in lexicographic SKU order, which
	// matches the sweeper and reservation.create paths.
	for _, snap := range rows {
		switch snap.Status {
		case StatusReleased, StatusExpired:
			// Idempotent: already released/expired — no locks needed.
			continue

		case StatusReserved:
			// LOCK-ORDER PIN: acquire stock_levels(sku) lock FIRST.
			if _, err := svc.StockStorage.GetBySKUForUpdate(ctx, tx, snap.SKU); err != nil {
				slog.ErrorContext(ctx, "reservation.release lock stock", slog.String("sku", snap.SKU), slog.String("error", err.Error()))
				wrapper.Respond(c, wrapper.ResponseOption[any]{
					HTTPStatus: http.StatusServiceUnavailable,
					Code:       CodeDatabaseUnavailable,
					Message:    MessageDatabaseUnavailable,
					TraceID:    traceID,
				})
				return
			}
			// Now lock the reservation row by PK (stock_levels already held above).
			r, err := svc.ResvStorage.GetByIDForUpdate(ctx, tx, snap.ID)
			if err != nil {
				slog.ErrorContext(ctx, "reservation.release re-lock reservation", slog.String("id", snap.ID.String()), slog.String("error", err.Error()))
				wrapper.Respond(c, wrapper.ResponseOption[any]{
					HTTPStatus: http.StatusServiceUnavailable,
					Code:       CodeDatabaseUnavailable,
					Message:    MessageDatabaseUnavailable,
					TraceID:    traceID,
				})
				return
			}
			if r.Status != StatusReserved {
				// Raced to RELEASED/EXPIRED between snapshot and re-lock — skip.
				continue
			}
			if err := svc.StockStorage.IncrementAvailableDecrementReserved(ctx, tx, r.SKU, r.Qty); err != nil {
				slog.ErrorContext(ctx, "reservation.release update stock", slog.String("sku", r.SKU), slog.String("error", err.Error()))
				wrapper.Respond(c, wrapper.ResponseOption[any]{
					HTTPStatus: http.StatusServiceUnavailable,
					Code:       CodeDatabaseUnavailable,
					Message:    MessageDatabaseUnavailable,
					TraceID:    traceID,
				})
				return
			}
			if err := svc.ResvStorage.MarkReleased(ctx, tx, r.ID, releaseReason); err != nil {
				slog.ErrorContext(ctx, "reservation.release mark released", slog.String("id", r.ID.String()), slog.String("error", err.Error()))
				wrapper.Respond(c, wrapper.ResponseOption[any]{
					HTTPStatus: http.StatusServiceUnavailable,
					Code:       CodeDatabaseUnavailable,
					Message:    MessageDatabaseUnavailable,
					TraceID:    traceID,
				})
				return
			}

		case StatusCommitted:
			// PR-006 soft-cancel of a PAID order: sold → available.
			// LOCK-ORDER PIN: acquire stock_levels(sku) lock FIRST.
			if _, err := svc.StockStorage.GetBySKUForUpdate(ctx, tx, snap.SKU); err != nil {
				slog.ErrorContext(ctx, "reservation.release lock stock (COMMITTED)", slog.String("sku", snap.SKU), slog.String("error", err.Error()))
				wrapper.Respond(c, wrapper.ResponseOption[any]{
					HTTPStatus: http.StatusServiceUnavailable,
					Code:       CodeDatabaseUnavailable,
					Message:    MessageDatabaseUnavailable,
					TraceID:    traceID,
				})
				return
			}
			// Now lock the reservation row by PK (stock_levels already held above).
			r, err := svc.ResvStorage.GetByIDForUpdate(ctx, tx, snap.ID)
			if err != nil {
				slog.ErrorContext(ctx, "reservation.release re-lock reservation (COMMITTED)", slog.String("id", snap.ID.String()), slog.String("error", err.Error()))
				wrapper.Respond(c, wrapper.ResponseOption[any]{
					HTTPStatus: http.StatusServiceUnavailable,
					Code:       CodeDatabaseUnavailable,
					Message:    MessageDatabaseUnavailable,
					TraceID:    traceID,
				})
				return
			}
			if r.Status != StatusCommitted {
				// Status changed between snapshot and re-lock — skip.
				continue
			}
			if err := svc.StockStorage.DecrementSoldIncrementAvailable(ctx, tx, r.SKU, r.Qty); err != nil {
				slog.ErrorContext(ctx, "reservation.release sold→available", slog.String("sku", r.SKU), slog.String("error", err.Error()))
				wrapper.Respond(c, wrapper.ResponseOption[any]{
					HTTPStatus: http.StatusServiceUnavailable,
					Code:       CodeDatabaseUnavailable,
					Message:    MessageDatabaseUnavailable,
					TraceID:    traceID,
				})
				return
			}
			if err := svc.ResvStorage.MarkReleased(ctx, tx, r.ID, ReasonPaidSoftCancel); err != nil {
				slog.ErrorContext(ctx, "reservation.release mark PAID_SOFT_CANCEL", slog.String("id", r.ID.String()), slog.String("error", err.Error()))
				wrapper.Respond(c, wrapper.ResponseOption[any]{
					HTTPStatus: http.StatusServiceUnavailable,
					Code:       CodeDatabaseUnavailable,
					Message:    MessageDatabaseUnavailable,
					TraceID:    traceID,
				})
				return
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		slog.ErrorContext(ctx, "reservation.release commit", slog.String("error", err.Error()))
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusServiceUnavailable,
			Code:       CodeDatabaseUnavailable,
			Message:    MessageDatabaseUnavailable,
			TraceID:    traceID,
		})
		return
	}

	respData := ReservationReleaseResponse{
		ReservationID: rows[0].ID.String(),
		Status:        "RELEASED",
	}
	wrapper.Respond(c, wrapper.ResponseOption[ReservationReleaseResponse]{
		HTTPStatus: http.StatusOK,
		Code:       CodeSuccess,
		Message:    MessageSuccess,
		Data:       &respData,
		TraceID:    traceID,
	})
}

// ─── Idempotency key DB helpers ──────────────────────────────────────────────

// idemFetch fetches and row-locks an idempotency_keys row.
// Returns (nil, pgx.ErrNoRows) if no matching row exists.
func idemFetch(ctx context.Context, tx pgx.Tx, key, endpoint string) (*IdempotencyKey, error) {
	const q = `
		SELECT key, endpoint, request_hash, response_envelope, http_status, created_at, expires_at
		FROM inventory.idempotency_keys
		WHERE key = $1 AND endpoint = $2
		FOR UPDATE`

	var ik IdempotencyKey
	err := tx.QueryRow(ctx, q, key, endpoint).
		Scan(&ik.Key, &ik.Endpoint, &ik.RequestHash, &ik.ResponseEnvelope, &ik.HTTPStatus, &ik.CreatedAt, &ik.ExpiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("idemFetch scan: %w", err)
	}
	return &ik, nil
}

// idemInsert persists an idempotency_keys row inside an open transaction.
func idemInsert(ctx context.Context, tx pgx.Tx, ik IdempotencyKey) error {
	const q = `
		INSERT INTO inventory.idempotency_keys
		    (key, endpoint, request_hash, response_envelope, http_status, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := tx.Exec(ctx, q,
		ik.Key, ik.Endpoint, ik.RequestHash, ik.ResponseEnvelope, ik.HTTPStatus, ik.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("idemInsert exec: %w", err)
	}
	return nil
}

// canonicalHash computes sha256 of the canonical JSON {orderId, items[sorted by sku]}.
func canonicalHash(orderID string, sortedItems []ReservationItem) string {
	payload := struct {
		OrderID string            `json:"orderId"`
		Items   []ReservationItem `json:"items"`
	}{
		OrderID: orderID,
		Items:   sortedItems,
	}
	b, _ := json.Marshal(payload)
	return fmt.Sprintf("%x", sha256.Sum256(b))
}
