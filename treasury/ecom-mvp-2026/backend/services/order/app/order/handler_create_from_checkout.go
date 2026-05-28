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

// HandleCreateFromCheckout implements POST /api/v1/order/order/create-from-checkout.
//
// This endpoint is INTERNAL — guarded by InternalAuthMiddleware (X-Internal-Secret header).
// It is the ONLY entry point into the order state machine; it sets status=PENDING_PAYMENT.
//
// Transaction:
//  1. BEGIN
//  2. SELECT nextval('order_number_seq') → format ORD-YYYYMMDD-NNNNNN
//  3. INSERT orders (status=PENDING_PAYMENT, version=0)
//  4. INSERT order_items (snapshot fields frozen — ORD-007)
//  5. INSERT order_status_history (from=NULL, to=PENDING_PAYMENT, actor=SYSTEM)
//  6. COMMIT
//
// No outbox event is emitted on creation (PENDING_PAYMENT is not broadcast — td.json line 234).
func (s *Service) HandleCreateFromCheckout(c *gin.Context) {
	var req CreateFromCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusBadRequest,
			Code:       app.CodeBadRequest,
			Message:    app.Message(err.Error()),
		})
		return
	}

	// REV-L2-101: use the IdempotencyKey (set by Checkout to the Checkout-generated
	// orderId) as the orders.id so that reservation/payment-intent foreign keys align.
	// The key must be a valid UUID; reject with 400 if it is not.
	parsedID, err := uuid.Parse(req.IdempotencyKey)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusBadRequest,
			Code:       app.CodeBadRequest,
			Message:    app.Message("idempotencyKey must be a valid UUID"),
		})
		return
	}
	orderID := parsedID.String()

	ctx := c.Request.Context()

	// Idempotent replay: if the row already exists return 200 with the existing data.
	existing, err := s.orders.GetByID(ctx, orderID)
	if err != nil {
		slog.ErrorContext(ctx, "create-from-checkout: check existing order", "error", err)
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
		})
		return
	}
	if existing != nil {
		wrapper.Respond(c, wrapper.ResponseOption[CreateFromCheckoutResponse]{
			HTTPStatus: http.StatusOK,
			Code:       app.CodeSuccess,
			Message:    app.MessageSuccess,
			Data: &CreateFromCheckoutResponse{
				OrderID:     existing.ID,
				OrderNumber: existing.OrderNumber,
				Status:      string(existing.Status),
				CreatedAt:   existing.CreatedAt.Format(time.RFC3339),
			},
		})
		return
	}

	now := time.Now().UTC()

	tx, err := s.beginTx(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "create-from-checkout: begin tx", "error", err)
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
		})
		return
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Step 2: allocate order_number via global sequence (ORD-008, td.json order_number_strategy).
	orderNumber, err := s.orders.NextOrderNumber(ctx, tx, now.Format("20060102"))
	if err != nil {
		slog.ErrorContext(ctx, "create-from-checkout: next order number", "error", err)
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
		})
		return
	}

	// Marshal address snapshot for JSONB storage.
	addrJSON, err := json.Marshal(req.AddressSnapshot)
	if err != nil {
		slog.ErrorContext(ctx, "create-from-checkout: marshal address", "error", err)
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
		})
		return
	}

	idempKey := req.IdempotencyKey
	o := Order{
		ID:                 orderID,
		OrderNumber:        orderNumber,
		UserID:             req.CustomerUserID,
		Status:             StatusPendingPayment,
		Subtotal:           req.Subtotal,
		ShippingFee:        req.ShippingFee,
		CouponDiscount:     req.CouponDiscount,
		GrandTotal:         req.GrandTotal,
		Currency:           "THB",
		AddressSnapshot:    addrJSON,
		BuyerEmailSnapshot: req.BuyerEmail,
		IdempotencyKey:     &idempKey,
		Version:            0,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	// Step 3: INSERT orders.
	if err := s.orders.Insert(ctx, tx, o); err != nil {
		slog.ErrorContext(ctx, "create-from-checkout: insert order", "error", err)
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
		})
		return
	}

	// Step 4: INSERT order_items (immutable snapshots — ORD-007).
	items := make([]OrderItem, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, OrderItem{
			ID:               uuid.New().String(),
			OrderID:          orderID,
			ProductID:        it.ProductID,
			SKU:              it.SKU,
			Qty:              it.Qty,
			NameSnapshot:     it.NameSnapshot,
			ImageURLSnapshot: it.ImageURLSnapshot,
			PriceSnapshot:    it.PriceSnapshot,
			LineSubtotal:     it.LineSubtotal,
			CreatedAt:        now,
		})
	}
	if err := s.items.BulkInsert(ctx, tx, items); err != nil {
		slog.ErrorContext(ctx, "create-from-checkout: bulk insert items", "error", err)
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
		})
		return
	}

	// Step 5: INSERT initial status_history row (from=NULL, ORD-009).
	reason := "checkout.commit"
	histRow := OrderStatusHistory{
		ID:          uuid.New().String(),
		OrderID:     orderID,
		FromStatus:  nil, // NULL for initial PENDING_PAYMENT row per td.json
		ToStatus:    StatusPendingPayment,
		ActorUserID: nil,
		ActorRole:   RoleSystem,
		Reason:      &reason,
		OccurredAt:  now,
	}
	if err := s.statusHistory.Insert(ctx, tx, histRow); err != nil {
		slog.ErrorContext(ctx, "create-from-checkout: insert status history", "error", err)
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
		})
		return
	}

	// Step 6: COMMIT.
	if err := tx.Commit(ctx); err != nil {
		slog.ErrorContext(ctx, "create-from-checkout: commit", "error", err)
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
		})
		return
	}

	slog.InfoContext(ctx, "order created",
		"orderId", orderID,
		"orderNumber", orderNumber,
		"userId", req.CustomerUserID,
	)

	wrapper.Respond(c, wrapper.ResponseOption[CreateFromCheckoutResponse]{
		HTTPStatus: http.StatusCreated,
		Code:       app.CodeSuccess,
		Message:    app.MessageSuccess,
		Data: &CreateFromCheckoutResponse{
			OrderID:     orderID,
			OrderNumber: orderNumber,
			Status:      string(StatusPendingPayment),
			CreatedAt:   now.Format(time.RFC3339),
		},
	})
}
