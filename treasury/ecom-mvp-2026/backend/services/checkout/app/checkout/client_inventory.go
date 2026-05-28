package checkout

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/httpclient"
)

// InventoryClient calls the inventory service.
type InventoryClient struct {
	baseURL            string
	httpClient         *http.Client
	internalSecret     string
	bulkReadTimeout    time.Duration
	reservationTimeout time.Duration
	releaseTimeout     time.Duration
}

func NewInventoryClient(baseURL, internalSecret string, bulkReadTimeout, reservationTimeout, releaseTimeout time.Duration) *InventoryClient {
	return &InventoryClient{
		baseURL:            baseURL,
		httpClient:         httpclient.NewHTTPClient(),
		internalSecret:     internalSecret,
		bulkReadTimeout:    bulkReadTimeout,
		reservationTimeout: reservationTimeout,
		releaseTimeout:     releaseTimeout,
	}
}

// StockLevel is the availableQty for one SKU from inventory.stock.bulk-read.
type StockLevel struct {
	SKU          string `json:"sku"`
	AvailableQty int    `json:"availableQty"`
}

type bulkReadRequest struct {
	SKUs []string `json:"skus"`
}

type bulkReadData struct {
	Stocks []StockLevel `json:"stocks"`
}

// BulkReadStock calls POST /api/v1/inventory/stock/bulk-read.
// timeout: 800ms, retries: 1 on 5xx.
func (c *InventoryClient) BulkReadStock(ctx context.Context, bearerToken string, skus []string) ([]StockLevel, error) {
	ctx, cancel := context.WithTimeout(ctx, c.bulkReadTimeout)
	defer cancel()

	url := c.baseURL + "/api/v1/inventory/stock/bulk-read"
	req := bulkReadRequest{SKUs: skus}

	var lastErr error
	for attempt := 0; attempt <= 1; attempt++ {
		res, err := httpclient.PostWithOptions[bulkReadRequest, svcEnvelope[bulkReadData]](
			ctx, c.httpClient, url, req,
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
			lastErr = fmt.Errorf("inventory.stock.bulk-read upstream %d", res.Code)
			continue
		}
		if res.Code >= 400 {
			return nil, &UpstreamError{
				Service: "inventory.stock.bulk-read", HTTPStatus: res.Code,
				Code: res.Response.Code, Message: res.Response.Message,
			}
		}
		if res.Response.Data == nil {
			return nil, fmt.Errorf("inventory.stock.bulk-read: empty data")
		}
		return res.Response.Data.Stocks, nil
	}
	return nil, lastErr
}

// ReservationItem is one line in the reservation request.
type ReservationItem struct {
	SKU string `json:"sku"`
	Qty int    `json:"qty"`
}

type reservationCreateRequest struct {
	OrderID        string            `json:"orderId"`
	Items          []ReservationItem `json:"items"`
	ExpiresAt      string            `json:"expiresAt"` // RFC3339
	IdempotencyKey string            `json:"idempotencyKey"`
}

type reservationCreateData struct {
	ReservationID string `json:"reservationId"`
}

// CreateReservation calls POST /api/v1/inventory/reservation/create.
// IDEMPOTENT on orderId. timeout: 1500ms, retries: 1 on 5xx.
func (c *InventoryClient) CreateReservation(ctx context.Context, bearerToken, orderID string, items []ReservationItem, expiresAt string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.reservationTimeout)
	defer cancel()

	url := c.baseURL + "/api/v1/inventory/reservation/create"
	req := reservationCreateRequest{
		OrderID:        orderID,
		Items:          items,
		ExpiresAt:      expiresAt,
		IdempotencyKey: orderID, // orderId doubles as idempotency key per spec
	}

	var lastErr error
	for attempt := 0; attempt <= 1; attempt++ {
		res, err := httpclient.PostWithOptions[reservationCreateRequest, svcEnvelope[reservationCreateData]](
			ctx, c.httpClient, url, req,
			bearerAuthOption(bearerToken),
			internalSecretOption(c.internalSecret),
		)
		if err != nil {
			lastErr = err
			if isRetryableErr(err) {
				continue
			}
			return "", err
		}
		if res.Code >= 500 && attempt == 0 {
			lastErr = fmt.Errorf("inventory.reservation.create upstream %d", res.Code)
			continue
		}
		if res.Code == 409 && res.Response.Code == "INSUFFICIENT_STOCK" {
			return "", &UpstreamError{
				Service: "inventory.reservation.create", HTTPStatus: 409,
				Code: "INSUFFICIENT_STOCK", Message: res.Response.Message,
			}
		}
		if res.Code >= 400 {
			return "", &UpstreamError{
				Service: "inventory.reservation.create", HTTPStatus: res.Code,
				Code: res.Response.Code, Message: res.Response.Message,
			}
		}
		if res.Response.Data == nil {
			return "", fmt.Errorf("inventory.reservation.create: empty data")
		}
		return res.Response.Data.ReservationID, nil
	}
	return "", lastErr
}

type reservationReleaseRequest struct {
	OrderID string `json:"orderId"`
	Reason  string `json:"reason"`
}

// ReleaseReservation calls POST /api/v1/inventory/reservation/release with exponential backoff.
// Compensation: timeout 1s per attempt, 3 retries with 50ms/200ms/800ms backoff.
// IDEMPOTENT on orderId.
func (c *InventoryClient) ReleaseReservation(ctx context.Context, bearerToken, orderID, reason string) error {
	url := c.baseURL + "/api/v1/inventory/reservation/release"
	req := reservationReleaseRequest{OrderID: orderID, Reason: reason}

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

		callCtx, cancel := context.WithTimeout(ctx, c.releaseTimeout)
		res, err := httpclient.PostWithOptions[reservationReleaseRequest, svcEnvelope[struct{}]](
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
			lastErr = fmt.Errorf("inventory.reservation.release upstream %d", res.Code)
			continue
		}
		if res.Code == 404 {
			// Already released or never existed — idempotent; treat as success.
			return nil
		}
		if res.Code >= 400 {
			return &UpstreamError{
				Service: "inventory.reservation.release", HTTPStatus: res.Code,
				Code: res.Response.Code, Message: res.Response.Message,
			}
		}
		return nil
	}
	return lastErr
}
