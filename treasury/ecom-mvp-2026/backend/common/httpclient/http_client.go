// Package httpclient provides functions for creating http.Client
// and calling other APIs
package httpclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"
)

const (
	maxIdleConns        = 100
	maxConnsPerHost     = 100
	maxIdleConnsPerHost = 100
	defaultTimeout      = 15 * time.Second
)

// OptionFunc defines the function to update *http.Request
type OptionFunc func(*http.Request, ...context.Context)

type transport struct {
	*http.Transport
	options []OptionFunc
}

func cloneDefaultTransport() *http.Transport {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = maxIdleConns
	t.MaxConnsPerHost = maxConnsPerHost
	t.MaxIdleConnsPerHost = maxIdleConnsPerHost
	return t
}

// NewHTTPClient returns *http.Client instance with the provided options
// and default timeout 5s.
func NewHTTPClient(options ...OptionFunc) *http.Client {
	t := cloneDefaultTransport()
	options = append(options, jsonOption)
	return &http.Client{
		Timeout: defaultTimeout,
		Transport: transport{
			Transport: t,
			options:   options,
		},
	}
}

// NewHTTPClientWithCA returns *http.Client instance with TLS configuration
// and the provided options and default timeout 5s.
// certPool is the root CA certificates to verify the server's certificate.
func NewHTTPClientWithCA(certPool *x509.CertPool, options ...OptionFunc) *http.Client {
	t := cloneDefaultTransport()
	t.TLSClientConfig = &tls.Config{
		RootCAs:    certPool,
		MinVersion: tls.VersionTLS12,
	}
	options = append(options, jsonOption)
	return &http.Client{
		Timeout: defaultTimeout,
		Transport: transport{
			Transport: t,
			options:   options,
		},
	}
}

// Get is a helper function to make a Get request
// and decode the response into the given type.
func Get[RES any](ctx context.Context, client *http.Client, url string) (response Response[RES], err error) {
	return do[any, RES](ctx, client, http.MethodGet, url, nil)
}

// GetWithOptions is a helper function to make a Get request
// with per-request options (e.g. dynamic headers).
func GetWithOptions[RES any](ctx context.Context, client *http.Client, url string, options ...OptionFunc) (response Response[RES], err error) {
	return doWithOptions[any, RES](ctx, client, http.MethodGet, url, nil, options...)
}

// Post is a helper function to make a Post request
// and decode the response into the given type.
func Post[REQ, RES any](ctx context.Context, client *http.Client, url string, payload REQ) (response Response[RES], err error) {
	return do[REQ, RES](ctx, client, http.MethodPost, url, payload)
}

// PostWithOptions is a helper function to make a Post request
// with per-request options (e.g. dynamic headers).
func PostWithOptions[REQ, RES any](ctx context.Context, client *http.Client, url string, payload REQ, options ...OptionFunc) (response Response[RES], err error) {
	return doWithOptions[REQ, RES](ctx, client, http.MethodPost, url, payload, options...)
}

// Put is a helper function to make a Put request
// and decode the response into the given type.
func Put[REQ, RES any](ctx context.Context, client *http.Client, url string, payload REQ) (response Response[RES], err error) {
	return do[REQ, RES](ctx, client, http.MethodPut, url, payload)
}

// PutWithOptions is a helper function to make a Put request
// with per-request options (e.g. dynamic headers).
func PutWithOptions[REQ, RES any](ctx context.Context, client *http.Client, url string, payload REQ, options ...OptionFunc) (response Response[RES], err error) {
	return doWithOptions[REQ, RES](ctx, client, http.MethodPut, url, payload, options...)
}

// Patch is a helper function to make a Patch request
// and decode the response into the given type.
func Patch[REQ, RES any](ctx context.Context, client *http.Client, url string, payload REQ) (response Response[RES], err error) {
	return do[REQ, RES](ctx, client, http.MethodPatch, url, payload)
}

// PatchWithOptions is a helper function to make a Patch request
// with per-request options (e.g. dynamic headers).
func PatchWithOptions[REQ, RES any](ctx context.Context, client *http.Client, url string, payload REQ, options ...OptionFunc) (response Response[RES], err error) {
	return doWithOptions[REQ, RES](ctx, client, http.MethodPatch, url, payload, options...)
}

// Delete is a helper function to make a Delete request
// and decode the response into the given type.
func Delete[RES any](ctx context.Context, client *http.Client, url string) (response Response[RES], err error) {
	return do[any, RES](ctx, client, http.MethodDelete, url, nil)
}

// DeleteWithOptions is a helper function to make a Delete request
// with per-request options (e.g. dynamic headers).
func DeleteWithOptions[RES any](ctx context.Context, client *http.Client, url string, options ...OptionFunc) (response Response[RES], err error) {
	return doWithOptions[any, RES](ctx, client, http.MethodDelete, url, nil, options...)
}

type Response[T any] struct {
	Code     int
	Response T
}

func do[REQ, RES any](ctx context.Context, client *http.Client, method, url string, payload REQ) (response Response[RES], err error) {
	req, err := newRequest(ctx, client, method, url, payload)
	if err != nil {
		return response, err
	}

	return doRequest[RES](client, req)
}

func doWithOptions[REQ, RES any](ctx context.Context, client *http.Client, method, url string, payload REQ, options ...OptionFunc) (response Response[RES], err error) {
	req, err := newRequest(ctx, client, method, url, payload)
	if err != nil {
		return response, err
	}

	applyOptions(req, ctx, options...)
	return doRequest[RES](client, req)
}

// DoRequest sends an HTTP request via client given
// and returns the Response with an HTTP response it gets
func DoRequest[RES any](client *http.Client, req *http.Request) (response Response[RES], err error) {
	return doRequest[RES](client, req)
}

func doRequest[RES any](client *http.Client, req *http.Request) (response Response[RES], err error) {
	debugEnabled := isDebugEnabled(req)
	if debugEnabled {
		if err = logRequest(req); err != nil {
			err = serror.WrapSkip(err, 3)
			return
		}
	}

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		err = serror.WrapSkip(err, 3)
		return
	}
	defer resp.Body.Close()

	var v RES
	if debugEnabled {
		var bodyBytes []byte
		bodyBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			err = serror.WrapSkip(err, 3)
			return
		}
		logResponse(req.Context(), resp, bodyBytes)
		if err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&v); err != nil {
			err = serror.WrapSkip(err, 3)
			return
		}
	} else {
		if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
			err = serror.WrapSkip(err, 3)
			return
		}
	}

	response = Response[RES]{
		Code:     resp.StatusCode,
		Response: v,
	}
	return
}

// NewRequest returns *http.Request
func NewRequest(ctx context.Context, client *http.Client, method, url string, payload any) (*http.Request, error) {
	return newRequest(ctx, client, method, url, payload)
}

func newRequest(ctx context.Context, client *http.Client, method, url string, payload any) (*http.Request, error) {
	var body io.Reader
	if payload != nil {
		if r, ok := payload.(io.Reader); ok {
			body = r
		} else {
			var buf bytes.Buffer
			if err := json.NewEncoder(&buf).Encode(payload); err != nil {
				return nil, serror.WrapSkip(err, 3)
			}
			body = &buf
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, serror.WrapSkip(err, 3)
	}

	if t, ok := client.Transport.(transport); ok {
		for _, option := range t.options {
			option(req, ctx)
		}
	}

	return req, nil
}

func applyOptions(req *http.Request, ctx context.Context, options ...OptionFunc) {
	for _, option := range options {
		option(req, ctx)
	}
}

func isDebugEnabled(req *http.Request) bool {
	if req == nil {
		return false
	}
	value := req.Context().Value(debugContextKey{})
	if value == nil {
		return false
	}
	enabled, ok := value.(bool)
	return ok && enabled
}

func logRequest(req *http.Request) error {
	if req == nil {
		return nil
	}

	var bodyBytes []byte
	if req.GetBody != nil {
		body, err := req.GetBody()
		if err != nil {
			return err
		}
		defer body.Close()
		bodyBytes, err = io.ReadAll(body)
		if err != nil {
			return err
		}
	} else if req.Body != nil {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return err
		}
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	slog.DebugContext(
		req.Context(),
		"httpclient.request",
		slog.String("method", req.Method),
		slog.String("url", req.URL.String()),
		slog.String("body", string(bodyBytes)),
	)

	return nil
}

func logResponse(ctx context.Context, resp *http.Response, bodyBytes []byte) {
	if resp == nil {
		return
	}

	slog.DebugContext(
		ctx,
		"httpclient.response",
		slog.Int("status", resp.StatusCode),
		slog.String("body", string(bodyBytes)),
	)
}
