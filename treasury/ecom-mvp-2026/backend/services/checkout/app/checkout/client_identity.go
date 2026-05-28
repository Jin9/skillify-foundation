package checkout

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/httpclient"
)

// IdentityClient calls the identity service for address operations.
type IdentityClient struct {
	baseURL        string
	httpClient     *http.Client
	addressTimeout time.Duration
	profileTimeout time.Duration
}

func NewIdentityClient(baseURL string, addressTimeout time.Duration) *IdentityClient {
	return &IdentityClient{
		baseURL:        baseURL,
		httpClient:     httpclient.NewHTTPClient(),
		addressTimeout: addressTimeout,
		profileTimeout: 500 * time.Millisecond,
	}
}

type addressListRequest struct{}

type addressListData struct {
	Addresses []AddressSnapshot `json:"addresses"`
}

// profileReadRequest is the empty body sent to identity.profile.read.
type profileReadRequest struct{}

// profileReadData holds the customer profile fields returned by identity.profile.read.
// NOTE: identity.profile.read is currently a STUB (returns 501 NOT_IMPLEMENTED).
// This integration will become functional once Identity un-stubs the endpoint.
// See REV-L2-001 repair notes.
type profileReadData struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// ReadProfile calls POST /api/v1/identity/profile/read for the authenticated customer.
// Returns the profile including email. timeout: 500ms, retries: 1 on 5xx.
//
// DEPENDENCY NOTE (REV-L2-001): identity.profile.read is stubbed (501) in the current
// Identity implementation. End-to-end checkout will not succeed until Identity
// un-stubs this endpoint. The contract and wiring are in place.
func (c *IdentityClient) ReadProfile(ctx context.Context, bearerToken string) (profileReadData, error) {
	ctx, cancel := context.WithTimeout(ctx, c.profileTimeout)
	defer cancel()

	url := c.baseURL + "/api/v1/identity/profile/read"

	var lastErr error
	for attempt := 0; attempt <= 1; attempt++ {
		res, err := httpclient.PostWithOptions[profileReadRequest, svcEnvelope[profileReadData]](
			ctx, c.httpClient, url, profileReadRequest{},
			bearerAuthOption(bearerToken),
		)
		if err != nil {
			lastErr = err
			if isRetryableErr(err) {
				continue
			}
			return profileReadData{}, err
		}
		if res.Code >= 500 && attempt == 0 {
			lastErr = fmt.Errorf("identity.profile.read upstream %d", res.Code)
			continue
		}
		if res.Code >= 400 {
			return profileReadData{}, &UpstreamError{
				Service: "identity.profile.read", HTTPStatus: res.Code,
				Code: res.Response.Code, Message: res.Response.Message,
			}
		}
		if res.Response.Data == nil {
			return profileReadData{}, fmt.Errorf("identity.profile.read: empty data")
		}
		return *res.Response.Data, nil
	}
	return profileReadData{}, lastErr
}

// ListAddresses calls POST /api/v1/identity/address/list for the authenticated customer.
// timeout: 500ms, retries: 1 on 5xx.
func (c *IdentityClient) ListAddresses(ctx context.Context, bearerToken string) ([]AddressSnapshot, error) {
	ctx, cancel := context.WithTimeout(ctx, c.addressTimeout)
	defer cancel()

	url := c.baseURL + "/api/v1/identity/address/list"

	var lastErr error
	for attempt := 0; attempt <= 1; attempt++ {
		res, err := httpclient.PostWithOptions[addressListRequest, svcEnvelope[addressListData]](
			ctx, c.httpClient, url, addressListRequest{},
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
			lastErr = fmt.Errorf("identity.address.list upstream %d", res.Code)
			continue
		}
		if res.Code >= 400 {
			return nil, &UpstreamError{
				Service: "identity.address.list", HTTPStatus: res.Code,
				Code: res.Response.Code, Message: res.Response.Message,
			}
		}
		if res.Response.Data == nil {
			return nil, fmt.Errorf("identity.address.list: empty data")
		}
		return res.Response.Data.Addresses, nil
	}
	return nil, lastErr
}
