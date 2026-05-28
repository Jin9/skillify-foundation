package order

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

const (
	defaultListLimit = 20
	maxListLimit     = 100
	minPage          = 1
)

// HandleListMine implements POST /api/v1/order/order/list-mine.
//
// Auth: JWT required. Scoped to claims.sub (user_id). Per td.json §order.list-mine:
//   - sort: createdAt DESC
//   - default page=1, limit=20; max limit=100
//   - optional status filter
func (s *Service) HandleListMine(c *gin.Context) {
	var req ListMineRequest
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

	// Normalise pagination — server enforces bounds, never trusts client.
	page := req.Page
	if page < minPage {
		page = minPage
	}
	limit := req.Limit
	if limit <= 0 {
		limit = defaultListLimit
	}
	if limit > maxListLimit {
		limit = maxListLimit
	}
	offset := (page - 1) * limit

	orders, total, err := s.orders.ListByUserID(
		c.Request.Context(),
		claims.Sub,
		req.Status,
		limit,
		offset,
	)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
		})
		return
	}

	summaries := make([]OrderSummary, 0, len(orders))
	for _, o := range orders {
		summaries = append(summaries, OrderSummary{
			OrderID:     o.ID,
			OrderNumber: o.OrderNumber,
			Status:      o.Status,
			GrandTotal:  o.GrandTotal,
			Currency:    o.Currency,
			CreatedAt:   o.CreatedAt.Format(time.RFC3339),
		})
	}

	wrapper.Respond(c, wrapper.ResponseOption[ListMineResponse]{
		HTTPStatus: http.StatusOK,
		Code:       app.CodeSuccess,
		Message:    app.MessageSuccess,
		Data: &ListMineResponse{
			Orders: summaries,
			Total:  total,
			Page:   page,
		},
	})
}
