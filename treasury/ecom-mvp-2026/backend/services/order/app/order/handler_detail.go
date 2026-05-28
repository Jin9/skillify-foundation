package order

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

// Application-level codes for the order aggregate.
const (
	CodeOrderNotOwned    app.Code    = "ORD404"
	MessageOrderNotOwned app.Message = "Order not found"

	CodeOrderNotFound    app.Code    = "ORD4041"
	MessageOrderNotFound app.Message = "Order not found"
)

// HandleDetail implements POST /api/v1/order/order/detail.
//
// Auth: JWT required. Ownership rule (per td.json §order.detail.ownership_scope):
//   - CUSTOMER: order must belong to claims.sub; miss → 404 ORDER_NOT_OWNED (no existence leak)
//   - ADMIN:    any order; miss → 404 NOT_FOUND
func (s *Service) HandleDetail(c *gin.Context) {
	var req DetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusBadRequest,
			Code:       app.CodeBadRequest,
			Message:    app.Message(err.Error()),
		})
		return
	}

	claims, err := token.ClaimsFromContext(c.Request.Context())
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusUnauthorized,
			Code:       app.CodeUnauthorized,
			Message:    app.MessageUnauthorized,
		})
		return
	}

	// Extract role from claims.Extra (role is a domain-level claim stored in Extra).
	role := claimRole(claims)
	ctx := c.Request.Context()

	var o *Order
	if role == "ADMIN" {
		o, err = s.orders.GetByID(ctx, req.OrderID)
		if err != nil || o == nil {
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusNotFound,
				Code:       CodeOrderNotFound,
				Message:    MessageOrderNotFound,
			})
			return
		}
	} else {
		// CUSTOMER: scope to user — do NOT expose existence for other users' orders.
		// Read-only path: use GetByID then check ownership. No FOR UPDATE needed here
		// since we are not mutating state.
		o, err = s.orders.GetByID(ctx, req.OrderID)
		if err != nil || o == nil || o.UserID != claims.Sub {
			wrapper.Respond(c, wrapper.ResponseOption[any]{
				HTTPStatus: http.StatusNotFound,
				Code:       CodeOrderNotOwned,
				Message:    MessageOrderNotOwned,
			})
			return
		}
	}

	// Fetch line items (ORD-007 — all snapshot fields).
	items, err := s.items.ListByOrderID(ctx, o.ID)
	if err != nil {
		slog.ErrorContext(ctx, "detail: list items", "orderId", o.ID, "error", err)
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
		})
		return
	}

	// Fetch status history (ORD-009 — sorted ASC).
	history, err := s.statusHistory.ListByOrderIDASC(ctx, o.ID)
	if err != nil {
		slog.ErrorContext(ctx, "detail: list history", "orderId", o.ID, "error", err)
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
		})
		return
	}

	// Unmarshal address snapshot.
	var addr AddressSnapshot
	if len(o.AddressSnapshot) > 0 {
		_ = json.Unmarshal(o.AddressSnapshot, &addr)
	}

	itemResponses := make([]ItemResponse, 0, len(items))
	for _, it := range items {
		itemResponses = append(itemResponses, ItemResponse{
			ItemID:           it.ID,
			ProductID:        it.ProductID,
			SKU:              it.SKU,
			Qty:              it.Qty,
			NameSnapshot:     it.NameSnapshot,
			ImageURLSnapshot: it.ImageURLSnapshot,
			PriceSnapshot:    it.PriceSnapshot,
			LineSubtotal:     it.LineSubtotal,
		})
	}

	histResponses := make([]HistoryResponse, 0, len(history))
	for _, h := range history {
		// REV-L2-005: actorUserId must not be exposed to customer-context viewers
		// to prevent leaking admin user IDs. Only ADMIN callers see the field.
		var actorUserID *string
		if role == "ADMIN" {
			actorUserID = h.ActorUserID
		}
		histResponses = append(histResponses, HistoryResponse{
			FromStatus:  h.FromStatus,
			ToStatus:    h.ToStatus,
			ActorRole:   h.ActorRole,
			ActorUserID: actorUserID,
			Reason:      h.Reason,
			OccurredAt:  h.OccurredAt.Format(time.RFC3339),
		})
	}

	wrapper.Respond(c, wrapper.ResponseOption[DetailResponse]{
		HTTPStatus: http.StatusOK,
		Code:       app.CodeSuccess,
		Message:    app.MessageSuccess,
		Data: &DetailResponse{
			OrderID:         o.ID,
			OrderNumber:     o.OrderNumber,
			Status:          o.Status,
			Items:           itemResponses,
			AddressSnapshot: addr,
			Subtotal:        o.Subtotal,
			ShippingFee:     o.ShippingFee,
			CouponDiscount:  o.CouponDiscount,
			GrandTotal:      o.GrandTotal,
			Currency:        o.Currency,
			TrackingNumber:  o.TrackingNumber,
			StatusHistory:   histResponses,
			CreatedAt:       o.CreatedAt.Format(time.RFC3339),
			UpdatedAt:       o.UpdatedAt.Format(time.RFC3339),
		},
	})
}

// claimRole extracts the "role" key from claims.Extra. Defaults to "CUSTOMER".
func claimRole(claims *token.Claims) string {
	if claims.Extra == nil {
		return "CUSTOMER"
	}
	if r, ok := claims.Extra["role"].(string); ok {
		return r
	}
	return "CUSTOMER"
}
