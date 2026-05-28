package cart

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/httpclient"
)

type inventoryBulkReadRequest struct {
	ProductIDs []uuid.UUID `json:"productIds"`
}

// StockEntry holds the available quantity for a single product.
type StockEntry struct {
	ProductID    uuid.UUID `json:"productId"`
	AvailableQty int       `json:"availableQty"`
}

type inventoryBulkReadData struct {
	Items []StockEntry `json:"items"`
}

type inventoryResponseEnvelope struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Data    *inventoryBulkReadData `json:"data,omitempty"`
}

// InventoryClient calls the inventory service.
type InventoryClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewInventoryClient(httpClient *http.Client, baseURL string) *InventoryClient {
	return &InventoryClient{httpClient: httpClient, baseURL: baseURL}
}

// BulkReadStock returns a map[productID]availableQty for the given product IDs.
// Products absent from the response are treated as availableQty=0 by callers.
func (c *InventoryClient) BulkReadStock(ctx context.Context, productIDs []uuid.UUID) (map[uuid.UUID]int, error) {
	url := fmt.Sprintf("%s/api/v1/inventory/stock/bulk-read", c.baseURL)
	req := inventoryBulkReadRequest{ProductIDs: productIDs}

	resp, err := httpclient.Post[inventoryBulkReadRequest, inventoryResponseEnvelope](ctx, c.httpClient, url, req)
	if err != nil {
		return nil, err
	}

	if resp.Code != http.StatusOK || resp.Response.Data == nil {
		return nil, fmt.Errorf("inventory: unexpected response code=%d envelope_code=%s", resp.Code, resp.Response.Code)
	}

	result := make(map[uuid.UUID]int, len(resp.Response.Data.Items))
	for _, entry := range resp.Response.Data.Items {
		result[entry.ProductID] = entry.AvailableQty
	}
	return result, nil
}

// SingleStock is a convenience wrapper used by cart.update-item to stock-check one SKU.
func (c *InventoryClient) SingleStock(ctx context.Context, productID uuid.UUID) (int, error) {
	m, err := c.BulkReadStock(ctx, []uuid.UUID{productID})
	if err != nil {
		return 0, err
	}
	return m[productID], nil
}
