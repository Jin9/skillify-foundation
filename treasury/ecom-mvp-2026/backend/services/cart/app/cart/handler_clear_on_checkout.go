package cart

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

// ClearOnCheckout handles POST /api/v1/cart/cart/clear-on-checkout.
// STUB — returns 501 Not Implemented.
// Full implementation:
//   - Validate customerUserId == claims.sub (no cross-customer leakage — AMB-002).
//   - Validate cartItemIds non-empty, orderId non-empty.
//   - Idempotency: check cart_clear_idempotency by orderId; return cached result on hit.
//   - Fetch cart_id for user; DELETE FROM cart_items WHERE cart_id=$1 AND id=ANY($2::uuid[]).
//   - Touch carts.updated_at.
//   - Persist idempotency record (INSERT ON CONFLICT DO NOTHING).
//   - Return {removed: N} (CHK-010).
func (h *handler) ClearOnCheckout(c *gin.Context) {
	wrapper.Respond(c, wrapper.ResponseOption[any]{
		HTTPStatus: http.StatusNotImplemented,
		Code:       "NOT_IMPLEMENTED",
		Message:    "cart.clear-on-checkout is not yet implemented",
	})
}
