package cart

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

// RemoveItem handles POST /api/v1/cart/cart/remove-item.
// STUB — returns 501 Not Implemented.
// Full implementation: validate cartItemId UUID, delete cart_items row (idempotent),
// touch carts.updated_at, return enriched cart (subtotal recalculated — CART-003/CART-004).
func (h *handler) RemoveItem(c *gin.Context) {
	wrapper.Respond(c, wrapper.ResponseOption[any]{
		HTTPStatus: http.StatusNotImplemented,
		Code:       "NOT_IMPLEMENTED",
		Message:    "cart.remove-item is not yet implemented",
	})
}
