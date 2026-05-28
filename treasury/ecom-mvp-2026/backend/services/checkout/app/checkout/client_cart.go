package checkout

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/httpclient"
)

// CartClient calls the cart service.
type CartClient struct {
	baseURL      string
	httpClient   *http.Client
	readTimeout  time.Duration
	clearTimeout time.Duration
}

func NewCartClient(baseURL string, readTimeout, clearTimeout time.Duration) *CartClient {
	return &CartClient{
		baseURL:      baseURL,
		httpClient:   httpclient.NewHTTPClient(),
		readTimeout:  readTimeout,
		clearTimeout: clearTimeout,
	}
}

// cartReadRequest is sent to POST /api/v1/cart/cart/read.
// No body fields — cart service identifies caller from the Authorization header.
type cartReadRequest struct{}

// CartItem is one line from cart.read.
type CartItem struct {
	CartItemID string `json:"cartItemId"`
	ProductID  string `json:"productId"`
	SKU        string `json:"sku"`
	Qty        int    `json:"qty"`
}

// svcEnvelope is the standard downstream response shape: {code, message, data}.
// Used as the decode target for all downstream HTTP calls.
type svcEnvelope[T any] struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    *T     `json:"data,omitempty"`
}

type cartReadData struct {
	Items []CartItem `json:"items"`
}

// ReadCart calls POST /api/v1/cart/cart/read for the authenticated customer.
// timeout: 800ms, retries: 1 on 5xx/connection failure.
func (c *CartClient) ReadCart(ctx context.Context, bearerToken string) ([]CartItem, error) {
	ctx, cancel := context.WithTimeout(ctx, c.readTimeout)
	defer cancel()

	url := c.baseURL + "/api/v1/cart/cart/read"

	var lastErr error
	for attempt := 0; attempt <= 1; attempt++ {
		res, err := httpclient.PostWithOptions[cartReadRequest, svcEnvelope[cartReadData]](
			ctx, c.httpClient, url, cartReadRequest{},
			bearerAuthOption(bearerToken),
		)
		if err != nil {
			lastErr = err
			if isRetryableErr(err) {
				continue
			}
			return nil, err
		}
		if res.Code >= 500 && attempt == 0 {
			lastErr = fmt.Errorf("cart.read upstream %d", res.Code)
			continue
		}
		if res.Code >= 400 {
			return nil, &UpstreamError{
				Service: "cart.read", HTTPStatus: res.Code,
				Code: res.Response.Code, Message: res.Response.Message,
			}
		}
		if res.Response.Data != nil {
			return res.Response.Data.Items, nil
		}
		return []CartItem{}, nil
	}
	return nil, lastErr
}

type clearOnCheckoutRequest struct {
	CustomerUserID string   `json:"customerUserId"`
	CartItemIDs    []string `json:"cartItemIds"`
	OrderID        string   `json:"orderId"`
}

// ClearOnCheckout calls POST /api/v1/cart/cart/clear-on-checkout.
// Best-effort: caller MUST log warning on error and NOT compensate the order.
// timeout: 800ms, retries: 2 (50ms/250ms backoff per spec).
func (c *CartClient) ClearOnCheckout(ctx context.Context, bearerToken, customerUserID, orderID string, cartItemIDs []string) error {
	ctx, cancel := context.WithTimeout(ctx, c.clearTimeout)
	defer cancel()

	url := c.baseURL + "/api/v1/cart/cart/clear-on-checkout"
	req := clearOnCheckoutRequest{
		CustomerUserID: customerUserID,
		CartItemIDs:    cartItemIDs,
		OrderID:        orderID,
	}

	backoffs := []time.Duration{50 * time.Millisecond, 250 * time.Millisecond}
	var lastErr error
	for attempt := 0; attempt <= 2; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoffs[attempt-1]):
			}
		}
		res, err := httpclient.PostWithOptions[clearOnCheckoutRequest, svcEnvelope[struct{}]](
			ctx, c.httpClient, url, req,
			bearerAuthOption(bearerToken),
		)
		if err != nil {
			lastErr = err
			continue
		}
		if res.Code >= 500 {
			lastErr = fmt.Errorf("cart.clear-on-checkout upstream %d", res.Code)
			continue
		}
		return nil
	}
	return lastErr
}
