package checkout

// client_order.go — internal caller for order.create-from-checkout and
// order.cancel-on-checkout-failure.
//
// CHK-AMBIG-001: these endpoints are not yet first-class contracts in contracts.json.
// Routes agreed per td.json CHK-AMBIG-001 chosen_path_for_this_TD:
//   POST /api/v1/order/order/create-from-checkout
//   POST /api/v1/order/order/cancel-on-checkout-failure
// Internal auth: X-Internal-Secret header (INTERNAL_SHARED_SECRET env var).
// TL action: add these routes to contracts.json and to cross-cutting.route-convention.

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/httpclient"
)

// OrderClient calls the order service internal endpoints.
type OrderClient struct {
	baseURL        string
	httpClient     *http.Client
	internalSecret string
	createTimeout  time.Duration
	cancelTimeout  time.Duration
}

func NewOrderClient(baseURL, internalSecret string, createTimeout, cancelTimeout time.Duration) *OrderClient {
	return &OrderClient{
		baseURL:        baseURL,
		httpClient:     httpclient.NewHTTPClient(),
		internalSecret: internalSecret,
		createTimeout:  createTimeout,
		cancelTimeout:  cancelTimeout,
	}
}

// OrderLineItem is the price-snapshot line sent to order.create-from-checkout.
type OrderLineItem struct {
	ProductID          string `json:"productId"`
	SKU                string `json:"sku"`
	CartItemID         string `json:"cartItemId"`
	Qty                int    `json:"qty"`
	PriceSnapshotMinor int64  `json:"priceSnapshotMinor"`
	ProductName        string `json:"productName"`
	ProductImageURL    string `json:"productImageUrl"`
}

type orderCreateRequest struct {
	OrderID          string          `json:"orderId"`
	CustomerUserID   string          `json:"customerUserId"`
	BuyerEmail       string          `json:"buyerEmail"`
	Items            []OrderLineItem `json:"items"`
	AddressSnapshot  AddressSnapshot `json:"addressSnapshot"`
	SubtotalMinor    int64           `json:"subtotalMinor"`
	ShippingFeeMinor int64           `json:"shippingFeeMinor"`
	CouponDiscount   int64           `json:"couponDiscount"`
	TotalMinor       int64           `json:"totalMinor"`
	IdempotencyKey   string          `json:"idempotencyKey"`
}

type orderCreateData struct {
	OrderID     string `json:"orderId"`
	OrderNumber string `json:"orderNumber"`
	Status      string `json:"status"`
}

// CreateFromCheckout calls POST /api/v1/order/order/create-from-checkout.
// IDEMPOTENT on orderId. timeout: 1000ms, retries: 1.
// Auth: Bearer + X-Internal-Secret.
func (c *OrderClient) CreateFromCheckout(ctx context.Context, bearerToken string, req orderCreateRequest) (orderCreateData, error) {
	ctx, cancel := context.WithTimeout(ctx, c.createTimeout)
	defer cancel()

	url := c.baseURL + "/api/v1/order/order/create-from-checkout"

	var lastErr error
	for attempt := 0; attempt <= 1; attempt++ {
		res, err := httpclient.PostWithOptions[orderCreateRequest, svcEnvelope[orderCreateData]](
			ctx, c.httpClient, url, req,
			bearerAuthOption(bearerToken),
			internalSecretOption(c.internalSecret),
		)
		if err != nil {
			lastErr = err
			if isRetryableErr(err) {
				continue
			}
			return orderCreateData{}, err
		}
		if res.Code >= 500 && attempt == 0 {
			lastErr = fmt.Errorf("order.create-from-checkout upstream %d", res.Code)
			continue
		}
		if res.Code >= 400 {
			return orderCreateData{}, &UpstreamError{
				Service: "order.create-from-checkout", HTTPStatus: res.Code,
				Code: res.Response.Code, Message: res.Response.Message,
			}
		}
		if res.Response.Data == nil {
			return orderCreateData{}, fmt.Errorf("order.create-from-checkout: empty data")
		}
		return *res.Response.Data, nil
	}
	return orderCreateData{}, lastErr
}

type orderCancelRequest struct {
	OrderID string `json:"orderId"`
	Reason  string `json:"reason"`
}

// CancelOnCheckoutFailure calls POST /api/v1/order/order/cancel-on-checkout-failure.
// Compensation: timeout 1s per attempt, 3 retries with 50ms/200ms/800ms backoff.
// IDEMPOTENT on orderId.
func (c *OrderClient) CancelOnCheckoutFailure(ctx context.Context, bearerToken, orderID, reason string) error {
	url := c.baseURL + "/api/v1/order/order/cancel-on-checkout-failure"
	req := orderCancelRequest{OrderID: orderID, Reason: reason}

	backoffs := []time.Duration{50 * time.Millisecond, 200 * time.Millisecond, 800 * time.Millisecond}
	var lastErr error
	for attempt := 0; attempt <= 3; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoffs[attempt-1]):
			}
		}

		callCtx, cancel := context.WithTimeout(ctx, c.cancelTimeout)
		res, err := httpclient.PostWithOptions[orderCancelRequest, svcEnvelope[struct{}]](
			callCtx, c.httpClient, url, req,
			bearerAuthOption(bearerToken),
			internalSecretOption(c.internalSecret),
		)
		cancel()

		if err != nil {
			lastErr = err
			continue
		}
		if res.Code >= 500 {
			lastErr = fmt.Errorf("order.cancel-on-checkout-failure upstream %d", res.Code)
			continue
		}
		if res.Code == 404 || res.Code == 409 {
			// Order not found or already in terminal state — treat as success.
			return nil
		}
		if res.Code >= 400 {
			return &UpstreamError{
				Service: "order.cancel-on-checkout-failure", HTTPStatus: res.Code,
				Code: res.Response.Code, Message: res.Response.Message,
			}
		}
		return nil
	}
	return lastErr
}
