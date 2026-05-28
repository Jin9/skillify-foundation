package httpclient

import (
	"bytes"
	"context"
	"crypto/x509"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type APIResponse struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

func withTestLogger(t *testing.T, level slog.Level) *bytes.Buffer {
	t.Helper()

	buf := &bytes.Buffer{}
	defaultLogger := slog.Default()
	handler := slog.NewTextHandler(buf, &slog.HandlerOptions{Level: level})
	slog.SetDefault(slog.New(handler))

	t.Cleanup(func() {
		slog.SetDefault(defaultLogger)
	})

	return buf
}

func TestHTTPClient(t *testing.T) {
	t.Run("should be able to GET request succesfuly", func(t *testing.T) {
		client := NewHTTPClient()
		serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code": 200, "message": "OK", "description": "Success"}`))
		}))
		defer serv.Close()
		ctx := context.Background()

		resp, err := Get[APIResponse](ctx, client, serv.URL)

		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 200, resp.Code)
		assert.Equal(t, "OK", resp.Response.Message)
		assert.Equal(t, "Success", resp.Response.Description)
	})

	t.Run("should be able to GET with dynamic header option", func(t *testing.T) {
		client := NewHTTPClient()
		traceID := "trace-123"
		serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got := r.Header.Get("X-Trace-Id")
			assert.Equal(t, traceID, got)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code": 200, "message": "OK", "description": "Success"}`))
		}))
		defer serv.Close()
		ctx := context.WithValue(context.Background(), "trace-id", traceID)

		opt := func(r *http.Request, ctx ...context.Context) {
			if len(ctx) == 0 || ctx[0] == nil {
				return
			}
			if v := ctx[0].Value("trace-id"); v != nil {
				r.Header.Set("X-Trace-Id", v.(string))
			}
		}

		resp, err := GetWithOptions[APIResponse](ctx, client, serv.URL, opt)

		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 200, resp.Code)
		assert.Equal(t, "OK", resp.Response.Message)
		assert.Equal(t, "Success", resp.Response.Description)
	})

	t.Run("should be able to POST request succesfuly", func(t *testing.T) {
		type RESQ struct {
			Name string `json:"name"`
		}
		client := NewHTTPClient()
		payload := RESQ{Name: "somebody name"}
		serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code": 200, "message": "OK", "description": "Success"}`))
		}))
		defer serv.Close()
		ctx := context.Background()

		resp, err := Post[RESQ, APIResponse](ctx, client, serv.URL, payload)

		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 200, resp.Code)
		assert.Equal(t, "OK", resp.Response.Message)
		assert.Equal(t, "Success", resp.Response.Description)
	})

	t.Run("should be able to POST with per-request header option", func(t *testing.T) {
		type RESQ struct {
			Name string `json:"name"`
		}
		client := NewHTTPClient()
		payload := RESQ{Name: "somebody name"}
		serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got := r.Header.Get("X-Request-Id")
			assert.Equal(t, "req-123", got)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code": 200, "message": "OK", "description": "Success"}`))
		}))
		defer serv.Close()
		ctx := context.Background()

		resp, err := PostWithOptions[RESQ, APIResponse](ctx, client, serv.URL, payload, HeaderOption("X-Request-Id", "req-123"))

		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 200, resp.Code)
		assert.Equal(t, "OK", resp.Response.Message)
		assert.Equal(t, "Success", resp.Response.Description)
	})

	t.Run("should be able to PUT request succesfuly", func(t *testing.T) {
		type REQ struct {
			Name string `json:"name"`
		}
		client := NewHTTPClient()
		payload := REQ{Name: "updated"}
		serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPut, r.Method)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code": 200, "message": "OK", "description": "Success"}`))
		}))
		defer serv.Close()
		ctx := context.Background()

		resp, err := Put[REQ, APIResponse](ctx, client, serv.URL, payload)

		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 200, resp.Code)
		assert.Equal(t, "OK", resp.Response.Message)
		assert.Equal(t, "Success", resp.Response.Description)
	})

	t.Run("should be able to PATCH request succesfuly", func(t *testing.T) {
		type REQ struct {
			Name string `json:"name"`
		}
		client := NewHTTPClient()
		payload := REQ{Name: "patched"}
		serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPatch, r.Method)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code": 200, "message": "OK", "description": "Success"}`))
		}))
		defer serv.Close()
		ctx := context.Background()

		resp, err := Patch[REQ, APIResponse](ctx, client, serv.URL, payload)

		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 200, resp.Code)
		assert.Equal(t, "OK", resp.Response.Message)
		assert.Equal(t, "Success", resp.Response.Description)
	})

	t.Run("should be able to DELETE request succesfuly", func(t *testing.T) {
		client := NewHTTPClient()
		serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodDelete, r.Method)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code": 200, "message": "OK", "description": "Success"}`))
		}))
		defer serv.Close()
		ctx := context.Background()

		resp, err := Delete[APIResponse](ctx, client, serv.URL)

		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 200, resp.Code)
		assert.Equal(t, "OK", resp.Response.Message)
		assert.Equal(t, "Success", resp.Response.Description)
	})

	t.Run("should be able to PUT with per-request header option", func(t *testing.T) {
		type REQ struct {
			Name string `json:"name"`
		}
		client := NewHTTPClient()
		payload := REQ{Name: "updated"}
		serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPut, r.Method)
			assert.Equal(t, "put-req", r.Header.Get("X-Request-Id"))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code": 200, "message": "OK", "description": "Success"}`))
		}))
		defer serv.Close()
		ctx := context.Background()

		resp, err := PutWithOptions[REQ, APIResponse](ctx, client, serv.URL, payload, HeaderOption("X-Request-Id", "put-req"))

		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 200, resp.Code)
		assert.Equal(t, "OK", resp.Response.Message)
		assert.Equal(t, "Success", resp.Response.Description)
	})

	t.Run("should be able to PATCH with per-request header option", func(t *testing.T) {
		type REQ struct {
			Name string `json:"name"`
		}
		client := NewHTTPClient()
		payload := REQ{Name: "patched"}
		serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPatch, r.Method)
			assert.Equal(t, "patch-req", r.Header.Get("X-Request-Id"))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code": 200, "message": "OK", "description": "Success"}`))
		}))
		defer serv.Close()
		ctx := context.Background()

		resp, err := PatchWithOptions[REQ, APIResponse](ctx, client, serv.URL, payload, HeaderOption("X-Request-Id", "patch-req"))

		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 200, resp.Code)
		assert.Equal(t, "OK", resp.Response.Message)
		assert.Equal(t, "Success", resp.Response.Description)
	})

	t.Run("should be able to DELETE with per-request header option", func(t *testing.T) {
		client := NewHTTPClient()
		serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodDelete, r.Method)
			assert.Equal(t, "del-req", r.Header.Get("X-Request-Id"))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code": 200, "message": "OK", "description": "Success"}`))
		}))
		defer serv.Close()
		ctx := context.Background()

		resp, err := DeleteWithOptions[APIResponse](ctx, client, serv.URL, HeaderOption("X-Request-Id", "del-req"))

		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 200, resp.Code)
		assert.Equal(t, "OK", resp.Response.Message)
		assert.Equal(t, "Success", resp.Response.Description)
	})
}

func TestNewRequest(t *testing.T) {
	t.Run("error when convert to json request", func(t *testing.T) {
		ctx := context.Background()
		client := NewHTTPClient()
		req, err := NewRequest(ctx, client, http.MethodPost, "http://0.0.0.0/api", make(chan int))

		if req != nil {
			t.Errorf("unsupported type expect expect nil request but actual %v\n", req)
		}
		if err == nil {
			t.Error("unsupported type expect not nil error")
		}
	})

	t.Run("error when create http request", func(t *testing.T) {
		ctx := context.Background()
		client := NewHTTPClient()
		req, err := NewRequest(ctx, client, "\\", "http://0.0.0.0/api", bytes.NewBufferString(`{}`))

		if req != nil {
			t.Errorf("invalid method expect expect nil request but actual %v\n", req)
		}
		if err == nil {
			t.Error("invalid method expect not nil error")
		}
	})

	t.Run("should set default Content-Type and apply client options", func(t *testing.T) {
		client := NewHTTPClient(HeaderOption("X-Client", "client"))
		ctx := context.Background()
		req, err := NewRequest(ctx, client, http.MethodGet, "http://localhost:8080", nil)

		assert.Nil(t, err)
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
		assert.Equal(t, "client", req.Header.Get("X-Client"))
	})

	t.Run("should accept io.Reader payload", func(t *testing.T) {
		client := NewHTTPClient()
		ctx := context.Background()
		body := bytes.NewBufferString("raw")
		req, err := NewRequest(ctx, client, http.MethodPost, "http://localhost:8080", body)

		assert.Nil(t, err)
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
	})

	t.Run("should create client with CA", func(t *testing.T) {
		pool := x509.NewCertPool()
		client := NewHTTPClientWithCA(pool)

		tt, ok := client.Transport.(transport)
		assert.True(t, ok)
		assert.NotNil(t, tt.TLSClientConfig)
		assert.Equal(t, pool, tt.TLSClientConfig.RootCAs)
	})
}

func TestDoRequest(t *testing.T) {
	t.Run("client do error", func(t *testing.T) {
		type RESQ struct {
			Name string `json:"name"`
		}
		client := NewHTTPClient()
		req := httptest.NewRequest(http.MethodPost, "http://0.0.0.0/api", bytes.NewBufferString("{}"))

		_, err := func() (Response[RESQ], error) {
			return DoRequest[RESQ](client, req)
		}()

		if err == nil {
			t.Error("invalid url expect not nil error")
		}
	})
	t.Run("error when decode json response", func(t *testing.T) {
		type RESQ struct {
			Name string `json:"name"`
		}
		client := NewHTTPClient()

		serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Add("Content-Type", "text/json")
			w.Write(nil)
		}))
		defer serv.Close()

		ctx := context.Background()
		payload := RESQ{Name: "somebody name"}
		req, _ := NewRequest(ctx, client, http.MethodPost, serv.URL, payload)

		_, err := func() (Response[RESQ], error) {
			return DoRequest[RESQ](client, req)
		}()

		if err == nil {
			t.Error("invalid response expect not nil error")
		}
	})

	t.Run("debug logging enabled", func(t *testing.T) {
		type RESQ struct {
			Code int `json:"code"`
		}
		type REQ struct {
			Name string `json:"name"`
		}

		buf := withTestLogger(t, slog.LevelDebug)
		client := NewHTTPClient(DebugOption(true))
		serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code": 200}`))
		}))
		defer serv.Close()

		_, err := Post[REQ, RESQ](context.Background(), client, serv.URL, REQ{Name: "demo"})

		assert.Nil(t, err)
		assert.Contains(t, buf.String(), "httpclient.request")
		assert.Contains(t, buf.String(), "httpclient.response")
	})
}

func TestDo(t *testing.T) {
	t.Run("POST method with request body and response", func(t *testing.T) {
		type RESP struct {
			Code        int    `json:"code"`
			Message     string `json:"message"`
			Description string `json:"description"`
		}

		type RESQ struct {
			Name string `json:"name"`
		}

		ctx := context.Background()
		client := NewHTTPClient()
		method := "POST"
		payload := RESQ{Name: "test"}
		serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code": 200, "message": "OK", "description": "Success"}`))
		}))
		defer serv.Close()

		resp, err := do[RESQ, RESP](ctx, client, method, serv.URL, payload)

		assert.Nil(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 200, resp.Code)
		assert.Equal(t, "OK", resp.Response.Message)
		assert.Equal(t, "Success", resp.Response.Description)
	})

	t.Run("GET method with no request body", func(t *testing.T) {
		type RESP struct {
			Code        int    `json:"code"`
			Message     string `json:"message"`
			Description string `json:"description"`
		}

		ctx := context.Background()
		client := NewHTTPClient()
		method := "GET"
		var payload bytes.Buffer
		serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code": 200, "message": "OK", "description": "Success"}`))
		}))
		defer serv.Close()

		reps, err := do[bytes.Buffer, RESP](ctx, client, method, serv.URL, payload)

		assert.Nil(t, err)
		assert.NotNil(t, reps)
		assert.Equal(t, 200, reps.Code)
		assert.Equal(t, "OK", reps.Response.Message)
		assert.Equal(t, "Success", reps.Response.Description)
	})

	t.Run("error newRequest", func(t *testing.T) {
		type RESP struct {
			Code        int    `json:"code"`
			Message     string `json:"message"`
			Description string `json:"description"`
		}

		ctx := context.Background()
		client := NewHTTPClient()
		method := "\\"
		var payload bytes.Buffer
		serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"code": 200, "message": "OK", "description": "Success"}`))
		}))
		defer serv.Close()

		_, err := do[bytes.Buffer, RESP](ctx, client, method, serv.URL, payload)

		if err == nil {
			t.Error("can not new http request expect not nil error")
		}
	})

	t.Run("error in doWithOptions from newRequest", func(t *testing.T) {
		client := NewHTTPClient()
		_, err := doWithOptions[any, any](context.Background(), client, "\\", "http://0.0.0.0/api", nil)
		if err == nil {
			t.Error("expected error from newRequest inside doWithOptions")
		}
	})

	t.Run("isDebugEnabled nil request", func(t *testing.T) {
		assert.False(t, isDebugEnabled(nil))
	})

	t.Run("logRequest nil request", func(t *testing.T) {
		assert.Nil(t, logRequest(nil))
	})

	t.Run("logResponse nil response", func(t *testing.T) {
		logResponse(context.Background(), nil, nil) // Should not panic
	})

	t.Run("logRequest body string error fallback", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "http://test", errReader{})
		err := logRequest(req)
		assert.Error(t, err)
	})

	t.Run("logRequest GetBody error", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "http://test", nil)
		req.GetBody = func() (io.ReadCloser, error) {
			return nil, context.Canceled
		}
		err := logRequest(req)
		assert.Error(t, err)
	})

	t.Run("logRequest GetBody read error", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "http://test", nil)
		req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(errReader{}), nil
		}
		err := logRequest(req)
		assert.Error(t, err)
	})

	t.Run("debug enabled decode failure", func(t *testing.T) {
		client := NewHTTPClient(DebugOption(true))
		serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{invalid-json}`))
		}))
		defer serv.Close()

		_, err := Get[APIResponse](context.Background(), client, serv.URL)
		assert.Error(t, err)
	})

	t.Run("debug enabled body read failure", func(t *testing.T) {
		client := &http.Client{
			Transport: mockErrorTransport{},
		}
		req, _ := http.NewRequestWithContext(
			context.WithValue(context.Background(), debugContextKey{}, true),
			http.MethodGet, "http://test", nil)

		_, err := DoRequest[APIResponse](client, req)
		assert.Error(t, err)
	})
}

type errReader struct{}

func (errReader) Read(p []byte) (n int, err error) {
	return 0, context.Canceled
}

type mockErrorTransport struct{}

func (mockErrorTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(errReader{}),
	}, nil
}
