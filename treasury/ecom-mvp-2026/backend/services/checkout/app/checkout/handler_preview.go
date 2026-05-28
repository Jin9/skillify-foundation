package checkout

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

// Preview handles POST /api/v1/checkout/checkout/preview.
//
// Orchestration (read-only, no side-effects):
//   Step 1: cart.read
//   Step 2: identity.address.list → ownership + snapshot
//   Step 3: catalog.product.detail (parallel fan-out)
//   Step 4: inventory.stock.bulk-read
//   Pricing: server-side via pricing.go (single source)
func (s *Service) Preview(c *gin.Context) {
	claims, ok := requireCustomerClaims(c)
	if !ok {
		return
	}

	// --- Request validation (CHK-005, CHK-AMBIG-002) ---
	var rawBody map[string]json.RawMessage
	if err := c.ShouldBindJSON(&rawBody); err != nil {
		respondValidationError(c, "invalid request body")
		return
	}
	if err := rejectForbiddenFields(rawBody, previewForbiddenFields); err != nil {
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

	// Step 1: cart.read
	slog.InfoContext(c.Request.Context(), "checkout.preview.step1.cart_read",
		slog.String("customerUserId", claims.Sub), slog.String("traceId", traceID))
	cartItems, err := s.cart.ReadCart(c.Request.Context(), bearerToken)
	if err != nil {
		respondUpstreamError(c, err, traceID)
		return
	}
	if len(cartItems) == 0 {
		wrapper.Respond(c, wrapper.ResponseOption[PreviewResponse]{
			HTTPStatus: http.StatusOK,
			Code:       CodeSuccess,
			Message:    "cart is empty",
			Data: &PreviewResponse{
				Items:    []PreviewItem{},
				Blockers: []string{},
			},
			TraceID: traceID,
		})
		return
	}

	// Step 2: identity.address.list → ownership + completeness check
	slog.InfoContext(c.Request.Context(), "checkout.preview.step2.address_list", slog.String("traceId", traceID))
	addresses, err := s.identity.ListAddresses(c.Request.Context(), bearerToken)
	if err != nil {
		respondUpstreamError(c, err, traceID)
		return
	}
	addrSnapshot, found := findAddress(addresses, shippingAddressID)
	if !found {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusForbidden,
			Code:       CodeAddressNotOwned,
			Message:    "shipping address not owned by caller",
			TraceID:    traceID,
		})
		return
	}
	if !addrSnapshot.IsComplete {
		wrapper.Respond(c, wrapper.ResponseOption[any]{
			HTTPStatus: http.StatusBadRequest,
			Code:       CodeAddressIncomplete,
			Message:    "shipping address is incomplete",
			TraceID:    traceID,
		})
		return
	}

	// Step 3: catalog.product.detail — parallel fan-out (cap 50 per spec)
	slog.InfoContext(c.Request.Context(), "checkout.preview.step3.catalog_fanout", slog.String("traceId", traceID))
	details := make([]ProductDetail, len(cartItems))
	detailErrors := make([]error, len(cartItems))
	{
		g, gCtx := errgroup.WithContext(c.Request.Context())
		for i, item := range cartItems {
			i, item := i, item // capture
			g.Go(func() error {
				d, err := s.catalog.GetProductDetail(gCtx, bearerToken, item.ProductID)
				details[i] = d
				detailErrors[i] = err
				return nil // collect individually; never abort group
			})
		}
		_ = g.Wait()
	}

	// Step 4: inventory.stock.bulk-read
	slog.InfoContext(c.Request.Context(), "checkout.preview.step4.stock_bulk_read", slog.String("traceId", traceID))
	skus := make([]string, 0, len(cartItems))
	for _, item := range cartItems {
		skus = append(skus, item.SKU)
	}
	stocks, err := s.inventory.BulkReadStock(c.Request.Context(), bearerToken, skus)
	if err != nil {
		respondUpstreamError(c, err, traceID)
		return
	}
	stockMap := stockBySkU(stocks)

	// Pricing (server-side only, using pricing.go)
	var subtotalMinor int64
	var blockers []string
	items := make([]PreviewItem, 0, len(cartItems))

	for i, ci := range cartItems {
		d := details[i]
		dErr := detailErrors[i]

		checkoutable := true

		if dErr != nil || d.Status != "ACTIVE" {
			code := "PRODUCT_INACTIVE"
			if d.Status == "DELETED" {
				code = "PRODUCT_DELETED"
			}
			blockers = append(blockers, code+":"+ci.ProductID)
			checkoutable = false
		}

		availQty := 0
		if s, ok := stockMap[ci.SKU]; ok {
			availQty = s.AvailableQty
		}
		if checkoutable && ci.Qty > availQty {
			blockers = append(blockers, "INSUFFICIENT_STOCK:"+ci.SKU)
			checkoutable = false
		}

		var lineSubtotalMinor int64
		if checkoutable {
			lineSubtotalMinor = d.PriceMinor * int64(ci.Qty)
			subtotalMinor += lineSubtotalMinor
		}

		items = append(items, PreviewItem{
			ProductID:    ci.ProductID,
			SKU:          ci.SKU,
			Qty:          ci.Qty,
			CurrentPrice: minorToTHB(d.PriceMinor),
			LineSubtotal: minorToTHB(lineSubtotalMinor),
			AvailableQty: availQty,
			Checkoutable: checkoutable,
		})
	}

	shippingFeeMinor := ComputeShippingFee(subtotalMinor)
	grandTotalMinor := ComputeGrandTotal(subtotalMinor, shippingFeeMinor, 0)

	if blockers == nil {
		blockers = []string{}
	}

	wrapper.Respond(c, wrapper.ResponseOption[PreviewResponse]{
		HTTPStatus: http.StatusOK,
		Code:       CodeSuccess,
		Message:    "preview computed",
		Data: &PreviewResponse{
			Items:           items,
			Subtotal:        minorToTHB(subtotalMinor),
			ShippingFee:     minorToTHB(shippingFeeMinor),
			CouponDiscount:  0,
			Total:           minorToTHB(grandTotalMinor),
			Blockers:        blockers,
			AddressSnapshot: addrSnapshot,
		},
		TraceID: traceID,
	})
}

// ---------------------------------------------------------------------------
// Helpers (preview-local)
// ---------------------------------------------------------------------------

func stockBySkU(stocks []StockLevel) map[string]StockLevel {
	m := make(map[string]StockLevel, len(stocks))
	for _, s := range stocks {
		m[s.SKU] = s
	}
	return m
}

func findAddress(addresses []AddressSnapshot, id string) (AddressSnapshot, bool) {
	for _, a := range addresses {
		if a.AddressID == id {
			return a, true
		}
	}
	return AddressSnapshot{}, false
}
