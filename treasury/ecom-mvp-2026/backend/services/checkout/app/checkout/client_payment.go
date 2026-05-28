package checkout

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/httpclient"
)

// PaymentClient calls the payment service.
type PaymentClient struct {
	baseURL        string
	httpClient     *http.Client
	internalSecret string
	intentTimeout  time.Duration
}

func NewPaymentClient(baseURL, internalSecret string, intentTimeout time.Duration) *PaymentClient {
	return &PaymentClient{
		baseURL:        baseURL,
		httpClient:     httpclient.NewHTTPClient(),
		internalSecret: internalSecret,
		intentTimeout:  intentTimeout,
	}
}

type paymentIntentCreateRequest struct {
	OrderID string `json:"orderId"`
	Amount  int64  `json:"amount"` // THB minor units
}

// PaymentIntent is the result from payment.intent.create.
type PaymentIntent struct {
	PaymentIntentID string `json:"paymentIntentId"`
	Status          string `json:"status"` // REQUIRES_PAYMENT
	Amount          int64  `json:"amount"`
}

type paymentIntentData struct {
	PaymentIntent PaymentIntent `json:"paymentIntent"`
}

// CreateIntent calls POST /api/v1/payment/intent/create.
// IDEMPOTENT on orderId. timeout: 1500ms, retries: 1 on 5xx.
func (c *PaymentClient) CreateIntent(ctx context.Context, bearerToken, orderID string, totalMinor int64) (PaymentIntent, error) {
	ctx, cancel := context.WithTimeout(ctx, c.intentTimeout)
	defer cancel()

	url := c.baseURL + "/api/v1/payment/intent/create"
	req := paymentIntentCreateRequest{OrderID: orderID, Amount: totalMinor}

	var lastErr error
	for attempt := 0; attempt <= 1; attempt++ {
		res, err := httpclient.PostWithOptions[paymentIntentCreateRequest, svcEnvelope[paymentIntentData]](
			ctx, c.httpClient, url, req,
			bearerAuthOption(bearerToken),
			internalSecretOption(c.internalSecret),
		)
		if err != nil {
			lastErr = err
			if isRetryableErr(err) {
				continue
			}
			return PaymentIntent{}, err
		}
		if res.Code >= 500 && attempt == 0 {
			lastErr = fmt.Errorf("payment.intent.create upstream %d", res.Code)
			continue
		}
		if res.Code >= 400 {
			return PaymentIntent{}, &UpstreamError{
				Service: "payment.intent.create", HTTPStatus: res.Code,
				Code: res.Response.Code, Message: res.Response.Message,
			}
		}
		if res.Response.Data == nil {
			return PaymentIntent{}, fmt.Errorf("payment.intent.create: empty data")
		}
		return res.Response.Data.PaymentIntent, nil
	}
	return PaymentIntent{}, lastErr
}
