package checkout

// handler_commit.go — full orchestration for POST /api/v1/checkout/checkout/commit.
//
// ORCHESTRATION STEPS (per td.json §orchestration_pseudocode):
//   0.  Idempotency-Key lookup / INSERT INFLIGHT
//   1.  cart.read
//   2.  identity.profile.read → buyerEmail (REV-L2-001 wire)
//   3.  identity.address.list → ownership + snapshot
//   4.  catalog.product.detail (parallel fan-out) + inventory.stock.bulk-read
//   5.  Server-side pricing snapshot + orderId generation
//   6.  inventory.reservation.create   [COMPENSATION: none if fails here]
//   7.  order.create-from-checkout     [COMPENSATION: reservation.release]
//   8.  payment.intent.create          [COMPENSATION: order.cancel + reservation.release]
//   9.  cart.clear-on-checkout         [best-effort; do NOT compensate]
//  10.  Finalize idempotency_keys → COMPLETED, INSERT saga_log, COMMIT TX
//
// COMPENSATION MATRIX:
//   step 6 fails → no state to undo (reservation never created)
//   step 7 fails → inventory.reservation.release (3 retries [50,200,800]ms backoff)
//   step 8 fails → order.cancel-on-checkout-failure (3 retries) THEN reservation.release (3 retries)
//   step 9 fails → log warning ONLY; order is PENDING_PAYMENT; cart-clear is self-healing

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/sync/errgroup"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/checkout/app/checkout/access"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

// Commit handles POST /api/v1/checkout/checkout/commit.
func (s *Service) Commit(c *gin.Context) {
	claims, ok := requireCustomerClaims(c)
	if !ok {
		return
	}

	// --- Idempotency-Key header ---
	idemKey := c.GetHeader("Idempotency-Key")
	if !isValidIdempotencyKey(idemKey) {
		respondValidationError(c, "Idempotency-Key header required (max 128 chars, pattern ^[A-Za-z0-9._:-]{1,128}$)")
		return
	}

	// --- Request body validation (CHK-005, CHK-AMBIG-002) ---
	var rawBody map[string]json.RawMessage
	if err := c.ShouldBindJSON(&rawBody); err != nil {
		respondValidationError(c, "invalid request body")
		return
	}
	if err := rejectForbiddenFields(rawBody, commitForbiddenFields); err != nil {
		respondValidationError(c, err.Error())
		return
	}
	shippingAddressID, err := requireStringField(rawBody, "shippingAddressId", 64)
	if err != nil {
		respondValidationError(c, err.Error())
		return
	}

	bearerToken := extractBearer(c)
	traceID := traceIDFromContext(c)
	customerUserID := claims.Sub

	// Compute request hash for idempotency comparison (sha256 of canonical body excluding auth).
	requestHash := computeRequestHash(rawBody)

	// -------------------------------------------------------------------------
	// Open checkout DB transaction — wraps entire orchestration.
	// The idempotency_keys row lock holds INFLIGHT through step 10.
	// -------------------------------------------------------------------------
	tx, err := s.db.BeginTx(c.Request.Context(), pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		respondDBUnavailable(c, err, traceID)
		return
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(context.Background())
		}
	}()

	// Step 0: Atomic idempotency claim — INSERT ON CONFLICT DO NOTHING RETURNING.
	// TryClaimOrLookup closes the Lookup→InsertInflight race window (H02): the
	// single atomic INSERT prevents two concurrent first-time callers from both
	// seeing NOT FOUND and both proceeding to orchestration.
	winner, existingRow, claimErr := s.idempotencyStore.TryClaimOrLookup(c.Request.Context(), tx, idemKey, customerUserID, requestHash)
	if claimErr != nil {
		if errors.Is(claimErr, access.ErrNotFound) {
			// Row vanished between the conflict and the follow-up SELECT (extremely
			// rare: janitor deleted an ABANDONED row in that tiny window).  Re-claim
			// by looping: rollback and let the client retry (they will re-enter with
			// a fresh tx and TryClaimOrLookup will win the INSERT).
			_ = tx.Rollback(context.Background())
			committed = true
			c.Header("Retry-After", "1")
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusConflict,
				Code:       CodeIdempotencyKeyInflight,
				Message:    "transient conflict; retry after 1s",
				TraceID:    traceID,
			})
			return
		}
		respondDBUnavailable(c, claimErr, traceID)
		return
	}

	if !winner {
		// Conflict — a row already exists; inspect its status.
		row := existingRow
		switch row.Status {
		case access.StatusCompleted:
			if row.RequestHash == requestHash {
				// Cache hit — return verbatim
				_ = tx.Rollback(context.Background())
				committed = true
				httpStatus := http.StatusOK
				if row.HTTPStatus != nil {
					httpStatus = *row.HTTPStatus
				}
				c.Data(httpStatus, "application/json", row.ResponseEnvelope)
				return
			}
			// Same key, different payload → 409 IDEMPOTENCY_KEY_REUSED
			_ = tx.Rollback(context.Background())
			committed = true
			originalOrderID := extractOrderIDFromEnvelope(row.ResponseEnvelope)
			wrapper.Respond(c, wrapper.ResponseOption[map[string]string]{
				HTTPStatus: http.StatusConflict,
				Code:       CodeIdempotencyKeyReused,
				Message:    "Idempotency-Key was used with a different request body",
				Data:       &map[string]string{"originalOrderId": originalOrderID},
				TraceID:    traceID,
			})
			return
		case access.StatusInflight:
			// Another request is still in-flight for this key.
			_ = tx.Rollback(context.Background())
			committed = true
			c.Header("Retry-After", "1")
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusConflict,
				Code:       CodeIdempotencyKeyInflight,
				Message:    "concurrent request with same Idempotency-Key is in progress; retry after 1s",
				TraceID:    traceID,
			})
			return
		case access.StatusAbandoned:
			// Janitor marked this row ABANDONED (stale INFLIGHT > 60s).
			// DELETE the dead row and re-claim so this request can proceed as fresh.
			const deleteAbandoned = `
				DELETE FROM checkout.idempotency_keys
				WHERE key = $1 AND customer_user_id = $2 AND status = 'ABANDONED'`
			if _, delErr := tx.Exec(c.Request.Context(), deleteAbandoned, idemKey, customerUserID); delErr != nil {
				respondDBUnavailable(c, delErr, traceID)
				return
			}
			// Re-INSERT now that the slot is free; this INSERT must win because we
			// hold the row-lock from the FOR UPDATE above through the same tx.
			const reInsert = `
				INSERT INTO checkout.idempotency_keys
				    (key, customer_user_id, request_hash, status, response_envelope, http_status, created_at, updated_at, expires_at)
				VALUES ($1, $2, $3, 'INFLIGHT', NULL, NULL, NOW(), NOW(), NOW() + INTERVAL '24 hours')`
			if _, insErr := tx.Exec(c.Request.Context(), reInsert, idemKey, customerUserID, requestHash); insErr != nil {
				respondDBUnavailable(c, insErr, traceID)
				return
			}
			// Fall through — we are now the winner; orchestration continues below.
		}
	}

	// -------------------------------------------------------------------------
	// ORCHESTRATION — TX is open; row lock holds INFLIGHT.
	// On any decided error we COMMIT the tx (with COMPLETED envelope) so retries
	// see the cached outcome rather than re-executing side-effects.
	// -------------------------------------------------------------------------

	// Helper: commit the tx with a terminal envelope, then return the HTTP response.
	finalizeAndReturn := func(status access.IdempotencyStatus, envelope json.RawMessage, httpStatus int, sagaStatus access.SagaStatus, steps []access.SagaStep, compensation json.RawMessage) {
		if err := s.idempotencyStore.Finalize(c.Request.Context(), tx, idemKey, customerUserID, status, envelope, httpStatus); err != nil {
			slog.ErrorContext(c.Request.Context(), "checkout.commit.finalize_failed", slog.String("err", err.Error()))
		}
		_ = s.sagaLogStore.Write(c.Request.Context(), tx, "", customerUserID, traceID, sagaStatus, steps, compensation)
		if err := tx.Commit(c.Request.Context()); err != nil {
			slog.ErrorContext(c.Request.Context(), "checkout.commit.tx_commit_failed", slog.String("err", err.Error()))
		}
		committed = true
		c.Data(httpStatus, "application/json", envelope)
	}

	// ---- Step 1: cart.read ----
	slog.InfoContext(c.Request.Context(), "checkout.commit.step1.cart_read",
		slog.String("traceId", traceID), slog.String("customerUserId", customerUserID))
	cartItems, err := s.cart.ReadCart(c.Request.Context(), bearerToken)
	if err != nil {
		env := buildErrorEnvelope(CodeUpstreamTimeout, "cart service unavailable", traceID)
		finalizeAndReturn(access.StatusCompleted, env, http.StatusGatewayTimeout, access.SagaCompleted, nil, nil)
		return
	}
	if len(cartItems) == 0 {
		// CHK-001: empty cart
		env := buildErrorEnvelope(CodeValidationError, "cart is empty (CHK-001)", traceID)
		finalizeAndReturn(access.StatusCompleted, env, http.StatusBadRequest, access.SagaCompleted, nil, nil)
		return
	}

	// ---- Step 2: identity.profile.read — fetch buyerEmail for order snapshot ----
	// REV-L2-001 wire: order.create-from-checkout requires buyerEmail (binding:"required,email").
	// identity.profile.read is currently STUBBED (501) in the Identity service. This call
	// will block end-to-end checkout until Identity un-stubs the endpoint.
	// See KNOWN_ISSUES.md REV-L2-004 note on stub dependencies.
	slog.InfoContext(c.Request.Context(), "checkout.commit.step2.profile_read", slog.String("traceId", traceID))
	profile, err := s.identity.ReadProfile(c.Request.Context(), bearerToken)
	if err != nil {
		env := buildErrorEnvelope(CodeUpstreamTimeout, "identity service unavailable (profile.read)", traceID)
		finalizeAndReturn(access.StatusCompleted, env, http.StatusGatewayTimeout, access.SagaCompleted, nil, nil)
		return
	}

	// ---- Step 3: identity.address.list ----
	slog.InfoContext(c.Request.Context(), "checkout.commit.step3.address_list", slog.String("traceId", traceID))
	addresses, err := s.identity.ListAddresses(c.Request.Context(), bearerToken)
	if err != nil {
		env := buildErrorEnvelope(CodeUpstreamTimeout, "identity service unavailable", traceID)
		finalizeAndReturn(access.StatusCompleted, env, http.StatusGatewayTimeout, access.SagaCompleted, nil, nil)
		return
	}
	addrSnapshot, found := findAddress(addresses, shippingAddressID)
	if !found {
		env := buildErrorEnvelope(CodeAddressNotOwned, "shipping address not owned by caller (CHK-002)", traceID)
		finalizeAndReturn(access.StatusCompleted, env, http.StatusForbidden, access.SagaCompleted, nil, nil)
		return
	}
	if !addrSnapshot.IsComplete {
		env := buildErrorEnvelope(CodeAddressIncomplete, "shipping address is incomplete (CHK-002)", traceID)
		finalizeAndReturn(access.StatusCompleted, env, http.StatusBadRequest, access.SagaCompleted, nil, nil)
		return
	}

	// ---- Step 4a: catalog.product.detail (parallel fan-out) ----
	slog.InfoContext(c.Request.Context(), "checkout.commit.step4a.catalog_fanout", slog.String("traceId", traceID))
	details := make([]ProductDetail, len(cartItems))
	detailErrors := make([]error, len(cartItems))
	{
		g, gCtx := errgroup.WithContext(c.Request.Context())
		for i, item := range cartItems {
			i, item := i, item
			g.Go(func() error {
				d, err := s.catalog.GetProductDetail(gCtx, bearerToken, item.ProductID)
				details[i] = d
				detailErrors[i] = err
				return nil
			})
		}
		_ = g.Wait()
	}
	// Collect blockers — any INACTIVE/DELETED/NOT_FOUND blocks the commit
	type blockItem struct {
		ProductID string `json:"productId"`
		Code      string `json:"code"`
	}
	var catalogBlockers []blockItem
	for i, ci := range cartItems {
		d := details[i]
		if detailErrors[i] != nil || d.ProductID == "" {
			catalogBlockers = append(catalogBlockers, blockItem{ProductID: ci.ProductID, Code: "PRODUCT_INACTIVE"})
			continue
		}
		if d.Status == "INACTIVE" {
			catalogBlockers = append(catalogBlockers, blockItem{ProductID: ci.ProductID, Code: "PRODUCT_INACTIVE"})
		} else if d.Status == "DELETED" {
			catalogBlockers = append(catalogBlockers, blockItem{ProductID: ci.ProductID, Code: "PRODUCT_DELETED"})
		}
	}
	if len(catalogBlockers) > 0 {
		firstCode := catalogBlockers[0].Code
		env := buildDataErrorEnvelope(wrapper.Code(firstCode), "one or more cart items are inactive or deleted (CHK-003)", traceID, catalogBlockers)
		finalizeAndReturn(access.StatusCompleted, env, http.StatusConflict, access.SagaCompleted, nil, nil)
		return
	}

	// ---- Step 4b: inventory.stock.bulk-read (pre-reservation fail-fast) ----
	slog.InfoContext(c.Request.Context(), "checkout.commit.step4b.stock_bulk_read", slog.String("traceId", traceID))
	skus := make([]string, 0, len(cartItems))
	for _, item := range cartItems {
		skus = append(skus, item.SKU)
	}
	stocks, err := s.inventory.BulkReadStock(c.Request.Context(), bearerToken, skus)
	if err != nil {
		env := buildErrorEnvelope(CodeDBUnavailable, "inventory service unavailable", traceID)
		finalizeAndReturn(access.StatusCompleted, env, http.StatusServiceUnavailable, access.SagaCompleted, nil, nil)
		return
	}
	stockMap := stockBySkU(stocks)

	type stockShortage struct {
		ProductID    string `json:"productId"`
		RequestedQty int    `json:"requestedQty"`
		AvailableQty int    `json:"availableQty"`
	}
	var shortages []stockShortage
	for i, ci := range cartItems {
		avail := 0
		if s, ok := stockMap[ci.SKU]; ok {
			avail = s.AvailableQty
		}
		if ci.Qty > avail {
			shortages = append(shortages, stockShortage{
				ProductID:    details[i].ProductID,
				RequestedQty: ci.Qty,
				AvailableQty: avail,
			})
		}
	}
	if len(shortages) > 0 {
		env := buildDataErrorEnvelope(CodeInsufficientStock, "one or more items have insufficient stock (CHK-004)", traceID, shortages)
		finalizeAndReturn(access.StatusCompleted, env, http.StatusConflict, access.SagaCompleted, nil, nil)
		return
	}

	// ---- Step 5: generate orderId + server-side pricing snapshot ----
	orderID, err := uuid.NewV7()
	if err != nil {
		// fallback to v4 if v7 unavailable (should never happen in practice)
		orderID = uuid.New()
	}
	orderIDStr := orderID.String()

	lineItems := make([]LineItem, len(cartItems))
	var subtotalMinor int64
	for i, ci := range cartItems {
		d := details[i]
		lineItems[i] = LineItem{
			ProductID:          d.ProductID,
			SKU:                ci.SKU,
			CartItemID:         ci.CartItemID,
			Qty:                ci.Qty,
			PriceSnapshotMinor: d.PriceMinor,
			ProductName:        d.Name,
			ProductImageURL:    d.ImageURL,
		}
		subtotalMinor += d.PriceMinor * int64(ci.Qty)
	}
	shippingFeeMinor := ComputeShippingFee(subtotalMinor)       // single source — pricing.go
	couponDiscountMinor := int64(0)
	grandTotalMinor := ComputeGrandTotal(subtotalMinor, shippingFeeMinor, couponDiscountMinor)

	// ---- Step 6: inventory.reservation.create ----
	// COMPENSATION on failure: none (reservation was never created).
	slog.InfoContext(c.Request.Context(), "checkout.commit.step6.reservation_create",
		slog.String("orderId", orderIDStr), slog.String("traceId", traceID))
	reservationItems := make([]ReservationItem, len(lineItems))
	for i, li := range lineItems {
		reservationItems[i] = ReservationItem{SKU: li.SKU, Qty: li.Qty}
	}
	expiresAt := time.Now().Add(15 * time.Minute).UTC().Format(time.RFC3339)
	_, err = s.inventory.CreateReservation(c.Request.Context(), bearerToken, orderIDStr, reservationItems, expiresAt)
	if err != nil {
		if IsUpstreamError(err, "INSUFFICIENT_STOCK") {
			// Race condition: stock was consumed between step 4 and step 6 (CHK-004)
			env := buildErrorEnvelope(CodeInsufficientStock, "stock exhausted between check and reservation (CHK-004 race)", traceID)
			finalizeAndReturn(access.StatusCompleted, env, http.StatusConflict, access.SagaCompleted, nil, nil)
		} else {
			env := buildErrorEnvelope(CodeUpstreamTimeout, "inventory reservation failed", traceID)
			finalizeAndReturn(access.StatusCompleted, env, http.StatusGatewayTimeout, access.SagaCompleted, nil, nil)
		}
		return
	}

	// ---- Step 7: order.create-from-checkout (CHK-AMBIG-001) ----
	// COMPENSATION on failure: inventory.reservation.release (3 retries, 50/200/800ms backoff).
	slog.InfoContext(c.Request.Context(), "checkout.commit.step7.order_create",
		slog.String("orderId", orderIDStr), slog.String("traceId", traceID))
	orderLineItems := make([]OrderLineItem, len(lineItems))
	for i, li := range lineItems {
		orderLineItems[i] = OrderLineItem{
			ProductID:          li.ProductID,
			SKU:                li.SKU,
			CartItemID:         li.CartItemID,
			Qty:                li.Qty,
			PriceSnapshotMinor: li.PriceSnapshotMinor,
			ProductName:        li.ProductName,
			ProductImageURL:    li.ProductImageURL,
		}
	}
	orderRes, err := s.order.CreateFromCheckout(c.Request.Context(), bearerToken, orderCreateRequest{
		OrderID:          orderIDStr,
		CustomerUserID:   customerUserID,
		BuyerEmail:       profile.Email,
		Items:            orderLineItems,
		AddressSnapshot:  addrSnapshot,
		SubtotalMinor:    subtotalMinor,
		ShippingFeeMinor: shippingFeeMinor,
		CouponDiscount:   couponDiscountMinor,
		TotalMinor:       grandTotalMinor,
		IdempotencyKey:   orderIDStr,
	})
	if err != nil {
		// COMPENSATE: release reservation
		slog.ErrorContext(c.Request.Context(), "checkout.commit.step7.order_create_failed",
			slog.String("orderId", orderIDStr), slog.String("err", err.Error()))

		compErr := s.inventory.ReleaseReservation(context.Background(), bearerToken, orderIDStr, "ORDER_CANCELLED")
		var compJSON json.RawMessage
		sagaStatus := access.SagaCompensated
		if compErr != nil {
			slog.ErrorContext(c.Request.Context(), "checkout.commit.compensation.step7.release_failed",
				slog.String("orderId", orderIDStr), slog.String("compErr", compErr.Error()),
				slog.String("level", "CRITICAL"))
			compJSON, _ = json.Marshal(map[string]string{"error": compErr.Error(), "step": "step7_compensation"})
			sagaStatus = access.SagaDegraded
			// Reservation TTL (15min) is the safety net.
		}

		env := buildErrorEnvelope(CodeUpstreamTimeout, "order creation failed", traceID)
		finalizeAndReturn(access.StatusCompleted, env, http.StatusGatewayTimeout, sagaStatus, nil, compJSON)
		return
	}

	// ---- Step 8: payment.intent.create ----
	// COMPENSATION on failure: (A) order.cancel-on-checkout-failure, (B) reservation.release.
	slog.InfoContext(c.Request.Context(), "checkout.commit.step8.payment_intent_create",
		slog.String("orderId", orderIDStr), slog.String("traceId", traceID))
	intent, err := s.payment.CreateIntent(c.Request.Context(), bearerToken, orderIDStr, grandTotalMinor)
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "checkout.commit.step8.payment_intent_failed",
			slog.String("orderId", orderIDStr), slog.String("err", err.Error()))

		var compNotes []string

		// Compensation A: cancel the order
		cancelErr := s.order.CancelOnCheckoutFailure(context.Background(), bearerToken, orderIDStr, "PAYMENT_INTENT_FAILURE")
		if cancelErr != nil {
			slog.ErrorContext(c.Request.Context(), "checkout.commit.compensation.step8a.cancel_failed",
				slog.String("orderId", orderIDStr), slog.String("compErr", cancelErr.Error()),
				slog.String("level", "CRITICAL"))
			compNotes = append(compNotes, "order_cancel_failed: "+cancelErr.Error())
		}

		// Compensation B: release inventory reservation
		relErr := s.inventory.ReleaseReservation(context.Background(), bearerToken, orderIDStr, "ORDER_CANCELLED")
		if relErr != nil {
			slog.ErrorContext(c.Request.Context(), "checkout.commit.compensation.step8b.release_failed",
				slog.String("orderId", orderIDStr), slog.String("compErr", relErr.Error()),
				slog.String("level", "CRITICAL"))
			compNotes = append(compNotes, "reservation_release_failed: "+relErr.Error())
		}

		sagaStatus := access.SagaCompensated
		var compJSON json.RawMessage
		if len(compNotes) > 0 {
			sagaStatus = access.SagaDegraded
			compJSON, _ = json.Marshal(map[string]any{"notes": compNotes, "step": "step8_compensation"})
		}

		env := buildErrorEnvelope(CodeUpstreamTimeout, "payment intent creation failed", traceID)
		finalizeAndReturn(access.StatusCompleted, env, http.StatusGatewayTimeout, sagaStatus, nil, compJSON)
		return
	}

	// ---- Step 9: cart.clear-on-checkout (best-effort) ----
	// COMPENSATION: NONE. Order is PENDING_PAYMENT. Failure is non-fatal.
	slog.InfoContext(c.Request.Context(), "checkout.commit.step9.cart_clear",
		slog.String("orderId", orderIDStr), slog.String("traceId", traceID))
	cartItemIDs := make([]string, len(cartItems))
	for i, ci := range cartItems {
		cartItemIDs[i] = ci.CartItemID
	}
	if err := s.cart.ClearOnCheckout(c.Request.Context(), bearerToken, customerUserID, orderIDStr, cartItemIDs); err != nil {
		// Non-fatal — log warning and emit metric; proceed to success response.
		slog.WarnContext(c.Request.Context(), "checkout.commit.step9.cart_clear_failed",
			slog.String("orderId", orderIDStr),
			slog.String("err", err.Error()),
			slog.String("metric", "checkout_cart_clear_failed"))
	}

	// ---- Step 10: finalize idempotency + saga_log + COMMIT TX ----
	commitResponse := CommitResponse{
		OrderID:        orderIDStr,
		OrderNumber:    orderRes.OrderNumber,
		Status:         "PENDING_PAYMENT",
		Subtotal:       minorToTHB(subtotalMinor),
		ShippingFee:    minorToTHB(shippingFeeMinor),
		CouponDiscount: 0,
		Total:          minorToTHB(grandTotalMinor),
		PaymentIntent: PaymentIntentResult{
			PaymentIntentID: intent.PaymentIntentID,
			Status:          "REQUIRES_PAYMENT",
			Amount:          minorToTHB(grandTotalMinor),
		},
	}

	successEnv, _ := json.Marshal(wrapper.Response[CommitResponse]{
		Code:    CodeCreated,
		Message: "order created",
		Data:    &commitResponse,
	})

	finalizeAndReturn(access.StatusCompleted, json.RawMessage(successEnv), http.StatusCreated, access.SagaCompleted, nil, nil)
}

// ---------------------------------------------------------------------------
// Commit-local helpers
// ---------------------------------------------------------------------------

// computeRequestHash returns sha256(canonicalized request body) as a hex string.
// Canonicalization: marshal the map (Go map iteration is non-deterministic, so
// we re-marshal the parsed JSON which sorts keys stably in json.Marshal for
// simple flat objects). For full RFC correctness a dedicated canonical-JSON
// library should be used in production.
func computeRequestHash(body map[string]json.RawMessage) string {
	canonical, _ := json.Marshal(body)
	sum := sha256.Sum256(canonical)
	return fmt.Sprintf("%x", sum)
}

func isValidIdempotencyKey(key string) bool {
	if len(key) == 0 || len(key) > 128 {
		return false
	}
	for _, c := range key {
		if !isIdempotencyKeyChar(c) {
			return false
		}
	}
	return true
}

func isIdempotencyKeyChar(c rune) bool {
	return (c >= 'A' && c <= 'Z') ||
		(c >= 'a' && c <= 'z') ||
		(c >= '0' && c <= '9') ||
		c == '.' || c == '_' || c == '-' || c == ':'
}

// extractOrderIDFromEnvelope pulls orderId from a cached response envelope for
// IDEMPOTENCY_KEY_REUSED responses. Returns empty string if not parseable.
func extractOrderIDFromEnvelope(envelope json.RawMessage) string {
	var env struct {
		Data *struct {
			OrderID string `json:"orderId"`
		} `json:"data"`
	}
	if err := json.Unmarshal(envelope, &env); err != nil || env.Data == nil {
		return ""
	}
	return env.Data.OrderID
}

// buildErrorEnvelope returns a JSON-encoded error envelope for caching in idempotency_keys.
// traceID is accepted for symmetry but the standard wrapper.Response does not embed it;
// the middleware writes traceId to the structured log line instead.
func buildErrorEnvelope(code wrapper.Code, message wrapper.Message, _ string) json.RawMessage {
	env := wrapper.Response[any]{
		Code:    code,
		Message: message,
		Data:    nil,
	}
	b, _ := json.Marshal(env)
	return json.RawMessage(b)
}

// buildDataErrorEnvelope returns a JSON-encoded error envelope with a data payload.
func buildDataErrorEnvelope[T any](code wrapper.Code, message wrapper.Message, _ string, data T) json.RawMessage {
	env := wrapper.Response[T]{
		Code:    code,
		Message: message,
		Data:    &data,
	}
	b, _ := json.Marshal(env)
	return json.RawMessage(b)
}
