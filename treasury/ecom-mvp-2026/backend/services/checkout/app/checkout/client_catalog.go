package checkout

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/httpclient"
)

// CatalogClient calls the catalog service.
type CatalogClient struct {
	baseURL       string
	httpClient    *http.Client
	detailTimeout time.Duration
}

func NewCatalogClient(baseURL string, detailTimeout time.Duration) *CatalogClient {
	return &CatalogClient{
		baseURL:       baseURL,
		httpClient:    httpclient.NewHTTPClient(),
		detailTimeout: detailTimeout,
	}
}

type productDetailRequest struct {
	ProductID string `json:"productId"`
}

// ProductDetail is the relevant subset of catalog.product.detail response.
type ProductDetail struct {
	ProductID  string `json:"productId"`
	SKU        string `json:"sku"`
	Status     string `json:"status"` // ACTIVE | INACTIVE | DELETED
	PriceMinor int64  `json:"priceMinor"`
	Name       string `json:"name"`
	ImageURL   string `json:"imageUrl"`
}

type productDetailData struct {
	Product ProductDetail `json:"product"`
}

// GetProductDetail calls POST /api/v1/catalog/product/detail for one productId.
// timeout: 500ms, retries: 1 on 5xx/connection failure.
func (c *CatalogClient) GetProductDetail(ctx context.Context, bearerToken, productID string) (ProductDetail, error) {
	ctx, cancel := context.WithTimeout(ctx, c.detailTimeout)
	defer cancel()

	url := c.baseURL + "/api/v1/catalog/product/detail"
	req := productDetailRequest{ProductID: productID}

	var lastErr error
	for attempt := 0; attempt <= 1; attempt++ {
		res, err := httpclient.PostWithOptions[productDetailRequest, svcEnvelope[productDetailData]](
			ctx, c.httpClient, url, req,
			bearerAuthOption(bearerToken),
		)
		if err != nil {
			lastErr = err
			if isRetryableErr(err) {
				continue
			}
			return ProductDetail{}, err
		}
		if res.Code >= 500 && attempt == 0 {
			lastErr = fmt.Errorf("catalog.product.detail upstream %d", res.Code)
			continue
		}
		if res.Code == 404 {
			return ProductDetail{}, &UpstreamError{Service: "catalog.product.detail", HTTPStatus: 404, Code: "NOT_FOUND", Message: "product not found"}
		}
		if res.Code >= 400 {
			return ProductDetail{}, &UpstreamError{
				Service: "catalog.product.detail", HTTPStatus: res.Code,
				Code: res.Response.Code, Message: res.Response.Message,
			}
		}
		if res.Response.Data == nil {
			return ProductDetail{}, fmt.Errorf("catalog.product.detail: empty data for productId=%s", productID)
		}
		return res.Response.Data.Product, nil
	}
	return ProductDetail{}, lastErr
}
