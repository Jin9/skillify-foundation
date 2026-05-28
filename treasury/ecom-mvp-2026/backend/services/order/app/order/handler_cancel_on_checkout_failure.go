package order

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

// cancelOnCheckoutFailureRequest is the inbound payload.
// orderID is required; the reason is fixed to "checkout-failure" by spec — not caller-supplied.
type cancelOnCheckoutFailureRequest struct {
	OrderID string `json:"orderId" binding:"required,uuid"`
}

// HandleCancelOnCheckoutFailure implements POST /api/v1/order/order/cancel-on-checkout-failure.
// INTERNAL endpoint — guarded by InternalAuthMiddleware (X-Internal-Secret header).
//
// CHK-009 compensation contract: checkout.commit calls this when payment.intent.create
// fails AFTER the order row was created. Transitions PENDING_PAYMENT → CANCELLED and
// emits events.order.cancelled via the outbox.
//
// Transaction:
//  1. BEGIN
//  2. SELECT orders WHERE id=$1 FOR UPDATE
//  3. If status=CANCELLED → 200 no-op (idempotent replay)
//  4. If status≠PENDING_PAYMENT → 409 INVALID_ORDER_STATE
//  5. ValidateTransition(PENDING_PAYMENT, CANCELLED, SYSTEM)
//  6. UPDATE orders SET status=CANCELLED, version=version+1
//  7. INSERT order_status_history (actor_role=SYSTEM, reason="checkout-failure")
//  8. INSERT outbox_events (events.order.cancelled)
//  9. COMMIT
func (s *Service) HandleCancelOnCheckoutFailure(c *gin.Context) {
	var req cancelOnCheckoutFailureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusBadRequest,
			Code:       app.CodeBadRequest,
			Message:    app.Message(err.Error()),
		})
		return
	}

	ctx := c.Request.Context()

	tx, err := s.beginTx(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "cancel-on-checkout-failure: begin tx", "error", err, "orderId", req.OrderID)
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
		})
		return
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Step 2: SELECT ... FOR UPDATE (sole-writer invariant).
	o, err := s.orders.GetByIDForUpdate(ctx, tx, req.OrderID)
	if err != nil {
		slog.ErrorContext(ctx, "cancel-on-checkout-failure: select for update", "error", err, "orderId", req.OrderID)
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
		})
		return
	}
	if o == nil {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusNotFound,
			Code:       app.CodeNotFound,
			Message:    "order not found",
		})
		return
	}

	// Step 3: idempotent no-op if already CANCELLED.
	if o.Status == StatusCancelled {
		slog.InfoContext(ctx, "cancel-on-checkout-failure: already cancelled, no-op",
			"orderId", req.OrderID)
		if err := tx.Commit(ctx); err != nil {
			slog.ErrorContext(ctx, "cancel-on-checkout-failure: commit no-op", "error", err)
		}
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusOK,
			Code:       app.CodeSuccess,
			Message:    app.MessageSuccess,
		})
		return
	}

	// Step 4: guard — only PENDING_PAYMENT may be cancelled via this path.
	if o.Status != StatusPendingPayment {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusConflict,
			Code:       "ORD409",
			Message:    app.Message("order is not in PENDING_PAYMENT state"),
		})
		return
	}

	// Step 5: state-machine gate (belt-and-suspenders; always passes for PENDING_PAYMENT→CANCELLED/SYSTEM).
	if err := ValidateTransition(o.Status, StatusCancelled, RoleSystem); err != nil {
		slog.ErrorContext(ctx, "cancel-on-checkout-failure: invalid transition", "error", err, "orderId", req.OrderID)
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusConflict,
			Code:       "ORD409",
			Message:    app.Message(err.Error()),
		})
		return
	}

	now := time.Now().UTC()

	// Step 6: UPDATE orders.status → CANCELLED (ADR-003 sole-writer).
	if err := s.orders.UpdateStatus(ctx, tx, o.ID, StatusCancelled, o.Version, nil); err != nil {
		slog.ErrorContext(ctx, "cancel-on-checkout-failure: update status", "error", err, "orderId", req.OrderID)
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
		})
		return
	}

	// Step 7: INSERT order_status_history (actor=SYSTEM, reason="checkout-failure").
	reason := "checkout-failure"
	histRow := OrderStatusHistory{
		ID:          uuid.New().String(),
		OrderID:     o.ID,
		FromStatus:  ptrStatus(StatusPendingPayment),
		ToStatus:    StatusCancelled,
		ActorUserID: nil,
		ActorRole:   RoleSystem,
		Reason:      &reason,
		OccurredAt:  now,
	}
	if err := s.statusHistory.Insert(ctx, tx, histRow); err != nil {
		slog.ErrorContext(ctx, "cancel-on-checkout-failure: insert history", "error", err, "orderId", req.OrderID)
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
		})
		return
	}

	// Step 8: INSERT outbox_events (events.order.cancelled).
	fromStatus := string(StatusPendingPayment)
	cancelActor := string(RoleSystem)
	outboxPayload := OrderCancelledPayload{
		EventID:     uuid.New().String(),
		OccurredAt:  now.Format(time.RFC3339),
		OrderID:     o.ID,
		CancelActor: cancelActor,
		FromStatus:  fromStatus,
		Reason:      &reason,
	}
	payloadJSON, err := json.Marshal(outboxPayload)
	if err != nil {
		slog.ErrorContext(ctx, "cancel-on-checkout-failure: marshal outbox payload", "error", err, "orderId", req.OrderID)
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
		})
		return
	}
	outboxRow := OutboxEvent{
		ID:          uuid.New().String(),
		AggregateID: o.ID,
		EventType:   "order.cancelled",
		PayloadJSON: payloadJSON,
		CreatedAt:   now,
	}
	if err := s.outbox.Insert(ctx, tx, outboxRow); err != nil {
		slog.ErrorContext(ctx, "cancel-on-checkout-failure: insert outbox", "error", err, "orderId", req.OrderID)
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
		})
		return
	}

	// Step 9: COMMIT.
	if err := tx.Commit(ctx); err != nil {
		slog.ErrorContext(ctx, "cancel-on-checkout-failure: commit", "error", err, "orderId", req.OrderID)
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
		})
		return
	}

	slog.InfoContext(ctx, "cancel-on-checkout-failure: order cancelled",
		"orderId", o.ID,
		"orderNumber", o.OrderNumber,
	)

	wrapper.Respond(c, wrapper.ResponseOption[any]{
		HTTPStatus: http.StatusOK,
		Code:       app.CodeSuccess,
		Message:    app.MessageSuccess,
	})
}
