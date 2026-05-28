package cart

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/httpclient"
)

// ProductDetail is the response shape from catalog.product.detail.
type ProductDetail struct {
	ProductID    uuid.UUID     `json:"productId"`
	SKU          string        `json:"sku"`
	Name         string        `json:"name"`
	Image        string        `json:"image"`
	CurrentPrice float64       `json:"currentPrice"`
	Status       ProductStatus `json:"status"`
}

type catalogProductDetailRequest struct {
	ProductID uuid.UUID `json:"productId"`
}

type catalogProductDetailData struct {
	Product ProductDetail `json:"product"`
}

type catalogResponseEnvelope struct {
	Code    string                   `json:"code"`
	Message string                   `json:"message"`
	Data    *catalogProductDetailData `json:"data,omitempty"`
}

// CatalogClient calls the catalog service.
type CatalogClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewCatalogClient(httpClient *http.Client, baseURL string) *CatalogClient {
	return &CatalogClient{httpClient: httpClient, baseURL: baseURL}
}

// GetProductDetail fetches live product detail from the catalog service.
// Returns (zero, context.DeadlineExceeded) on timeout — callers should
// surface this as UPSTREAM_TIMEOUT per the TD contract.
func (c *CatalogClient) GetProductDetail(ctx context.Context, productID uuid.UUID) (ProductDetail, error) {
	url := fmt.Sprintf("%s/api/v1/catalog/product/detail", c.baseURL)
	req := catalogProductDetailRequest{ProductID: productID}

	resp, err := httpclient.Post[catalogProductDetailRequest, catalogResponseEnvelope](ctx, c.httpClient, url, req)
	if err != nil {
		return ProductDetail{}, err
	}

	if resp.Code != http.StatusOK || resp.Response.Data == nil {
		return ProductDetail{}, fmt.Errorf("catalog: unexpected response code=%d envelope_code=%s", resp.Code, resp.Response.Code)
	}

	return resp.Response.Data.Product, nil
}
