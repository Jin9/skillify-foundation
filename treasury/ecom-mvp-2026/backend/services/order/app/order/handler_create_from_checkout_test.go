package order

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

// --- Inline stub implementations of storage interfaces ---
// Minimal, test-only stubs. No mockery dependency.

type stubOrders struct {
	nextOrderNumber string
	nextOrderErr    error
	insertErr       error
}

func (s *stubOrders) Insert(_ context.Context, _ pgx.Tx, _ Order) error {
	return s.insertErr
}
func (s *stubOrders) GetByIDForUpdate(_ context.Context, _ pgx.Tx, _ string) (*Order, error) {
	return nil, nil
}
func (s *stubOrders) GetByIDScopedForUpdate(_ context.Context, _ pgx.Tx, _, _ string) (*Order, error) {
	return nil, nil
}
func (s *stubOrders) GetByID(_ context.Context, _ string) (*Order, error) { return nil, nil }
func (s *stubOrders) UpdateStatus(_ context.Context, _ pgx.Tx, _ string, _ OrderStatus, _ int, _ *string) error {
	return nil
}
func (s *stubOrders) ListByUserID(_ context.Context, _ string, _ *OrderStatus, _, _ int) ([]Order, int, error) {
	return nil, 0, nil
}
func (s *stubOrders) NextOrderNumber(_ context.Context, _ pgx.Tx, _ string) (string, error) {
	return s.nextOrderNumber, s.nextOrderErr
}

type stubItems struct{ bulkInsertErr error }

func (s *stubItems) BulkInsert(_ context.Context, _ pgx.Tx, _ []OrderItem) error {
	return s.bulkInsertErr
}
func (s *stubItems) ListByOrderID(_ context.Context, _ string) ([]OrderItem, error) {
	return nil, nil
}

type stubStatusHistory struct{ insertErr error }

func (s *stubStatusHistory) Insert(_ context.Context, _ pgx.Tx, _ OrderStatusHistory) error {
	return s.insertErr
}
func (s *stubStatusHistory) ListByOrderIDASC(_ context.Context, _ string) ([]OrderStatusHistory, error) {
	return nil, nil
}

type stubOutbox struct{}

func (s *stubOutbox) Insert(_ context.Context, _ pgx.Tx, _ OutboxEvent) error { return nil }
func (s *stubOutbox) MarkPublished(_ context.Context, _ string) error          { return nil }
func (s *stubOutbox) PollUnpublished(_ context.Context, _ int) ([]OutboxEvent, error) {
	return nil, nil
}

type stubConsumed struct{}

func (s *stubConsumed) InsertOrConflict(_ context.Context, _ pgx.Tx, _, _, _, _ string) (bool, error) {
	return true, nil
}

// --- Test helpers ---

// newTestService wires a Service with a nil pool (no DB). Only use for binding-
// error tests where the handler never reaches s.beginTx().
func newTestService(orders StorageOrder, items StorageOrderItem, hist StorageStatusHistory) *Service {
	return &Service{
		pool:          nil,
		orders:        orders,
		items:         items,
		statusHistory: hist,
		outbox:        &stubOutbox{},
		consumed:      &stubConsumed{},
	}
}

func setupRouter(svc *Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/create-from-checkout", svc.HandleCreateFromCheckout)
	return r
}

func doPost(r *gin.Engine, path string, body any) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func decodeResponse(t *testing.T, w *httptest.ResponseRecorder) wrapper.Response[json.RawMessage] {
	t.Helper()
	var resp wrapper.Response[json.RawMessage]
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response body: %v\nbody: %s", err, w.Body.String())
	}
	return resp
}

func validCreatePayload() map[string]any {
	return map[string]any{
		"idempotencyKey": "idem-001",
		"customerUserId": "550e8400-e29b-41d4-a716-446655440000",
		"buyerEmail":     "buyer@example.com",
		"addressSnapshot": map[string]any{
			"receiverName": "Test User",
			"phone":        "0800000000",
			"line":         "123 Main St",
			"province":     "Bangkok",
			"district":     "Watthana",
			"postalCode":   "10110",
		},
		"items": []map[string]any{
			{
				"productId":    "660e8400-e29b-41d4-a716-446655440001",
				"sku":          "SKU-001",
				"qty":          2,
				"priceSnapshot": 5000,
				"nameSnapshot": "Widget",
				"lineSubtotal": 10000,
			},
		},
		"subtotal":    10000,
		"shippingFee": 60,
		"grandTotal":  10060,
	}
}

// --- Bad-request tests (nil pool safe: handler returns before beginTx) ---

// TestHandleCreateFromCheckout_BadRequest verifies that missing or invalid
// required fields return 400 before any storage is touched.
func TestHandleCreateFromCheckout_BadRequest(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		payload map[string]any
	}{
		{
			name:    "empty body",
			payload: map[string]any{},
		},
		{
			name: "missing idempotencyKey",
			payload: func() map[string]any {
				p := validCreatePayload()
				delete(p, "idempotencyKey")
				return p
			}(),
		},
		{
			name: "invalid customerUserId (not uuid)",
			payload: func() map[string]any {
				p := validCreatePayload()
				p["customerUserId"] = "not-a-uuid"
				return p
			}(),
		},
		{
			name: "missing items (empty slice)",
			payload: func() map[string]any {
				p := validCreatePayload()
				p["items"] = []any{}
				return p
			}(),
		},
		{
			name: "invalid buyerEmail",
			payload: func() map[string]any {
				p := validCreatePayload()
				p["buyerEmail"] = "not-an-email"
				return p
			}(),
		},
	}

	svc := newTestService(&stubOrders{nextOrderNumber: "ORD-20260508-000001"}, &stubItems{}, &stubStatusHistory{})
	r := setupRouter(svc)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			w := doPost(r, "/create-from-checkout", tc.payload)

			if w.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d\nbody: %s", w.Code, w.Body.String())
			}
			resp := decodeResponse(t, w)
			if string(resp.Code) != string(wrapper.CodeBadRequest) {
				t.Errorf("expected code %s, got %s", wrapper.CodeBadRequest, resp.Code)
			}
		})
	}
}

// TestHandleCreateFromCheckout_InvalidJSON verifies that a malformed JSON body
// returns 400 before hitting any storage.
func TestHandleCreateFromCheckout_InvalidJSON(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	svc := newTestService(&stubOrders{nextOrderNumber: "ORD-20260508-000001"}, &stubItems{}, &stubStatusHistory{})
	r := setupRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/create-from-checkout", bytes.NewReader([]byte(`{bad json`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d\nbody: %s", w.Code, w.Body.String())
	}
}

// TestHandleCreateFromCheckout_ItemValidation verifies item-level binding rules
// (qty >= 1 per CreateItemRequest).
func TestHandleCreateFromCheckout_ItemValidation(t *testing.T) {
	t.Parallel()

	payload := validCreatePayload()
	items := payload["items"].([]map[string]any)
	items[0]["qty"] = 0 // invalid: min=1
	payload["items"] = items

	svc := newTestService(&stubOrders{nextOrderNumber: "ORD-20260508-000001"}, &stubItems{}, &stubStatusHistory{})
	r := setupRouter(svc)

	w := doPost(r, "/create-from-checkout", payload)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for qty=0, got %d\nbody: %s", w.Code, w.Body.String())
	}
}

// TestHandleCreateFromCheckout_ItemMissingSKU verifies that a missing SKU
// in an item triggers a 400 binding error.
func TestHandleCreateFromCheckout_ItemMissingSKU(t *testing.T) {
	t.Parallel()

	payload := validCreatePayload()
	payload["items"] = []map[string]any{
		{
			"productId":    "660e8400-e29b-41d4-a716-446655440001",
			"qty":          1,
			"nameSnapshot": "Widget",
			// sku absent
		},
	}

	svc := newTestService(&stubOrders{nextOrderNumber: "ORD-20260508-000001"}, &stubItems{}, &stubStatusHistory{})
	r := setupRouter(svc)

	w := doPost(r, "/create-from-checkout", payload)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing sku, got %d\nbody: %s", w.Code, w.Body.String())
	}
}

// --- Domain-object / DTO tests (no DB, no gin) ---

// TestHandleCreateFromCheckout_RequestSerialisationRoundTrip verifies that
// a complete CreateFromCheckoutRequest round-trips through JSON without data loss.
// This covers the happy-path binding logic at the struct level, decoupled from
// the gin engine and pool.
func TestHandleCreateFromCheckout_RequestSerialisationRoundTrip(t *testing.T) {
	t.Parallel()

	original := CreateFromCheckoutRequest{
		IdempotencyKey: "idem-roundtrip",
		CustomerUserID: "550e8400-e29b-41d4-a716-446655440000",
		BuyerEmail:     "buyer@example.com",
		AddressSnapshot: AddressSnapshot{
			ReceiverName: "Alice",
			Phone:        "0812345678",
			Line:         "1 Example Rd",
			Province:     "Chiang Mai",
			District:     "Mueang",
			PostalCode:   "50000",
		},
		Items: []CreateItemRequest{
			{
				ProductID:     "660e8400-e29b-41d4-a716-446655440001",
				SKU:           "SKU-RT-001",
				Qty:           3,
				PriceSnapshot: 100,
				NameSnapshot:  "Gadget",
				LineSubtotal:  300,
			},
		},
		Subtotal:    300,
		ShippingFee: 60,
		GrandTotal:  360,
	}

	b, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded CreateFromCheckoutRequest
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if decoded.IdempotencyKey != original.IdempotencyKey {
		t.Errorf("IdempotencyKey: want %s got %s", original.IdempotencyKey, decoded.IdempotencyKey)
	}
	if decoded.CustomerUserID != original.CustomerUserID {
		t.Errorf("CustomerUserID: want %s got %s", original.CustomerUserID, decoded.CustomerUserID)
	}
	if decoded.BuyerEmail != original.BuyerEmail {
		t.Errorf("BuyerEmail: want %s got %s", original.BuyerEmail, decoded.BuyerEmail)
	}
	if decoded.AddressSnapshot.ReceiverName != original.AddressSnapshot.ReceiverName {
		t.Errorf("AddressSnapshot.ReceiverName: want %s got %s",
			original.AddressSnapshot.ReceiverName, decoded.AddressSnapshot.ReceiverName)
	}
	if len(decoded.Items) != 1 {
		t.Fatalf("Items count: want 1 got %d", len(decoded.Items))
	}
	if decoded.Items[0].Qty != 3 {
		t.Errorf("Items[0].Qty: want 3 got %d", decoded.Items[0].Qty)
	}
	if decoded.GrandTotal != 360 {
		t.Errorf("GrandTotal: want 360 got %d", decoded.GrandTotal)
	}
}

// TestHandleCreateFromCheckout_ResponseShape verifies the CreateFromCheckoutResponse
// JSON shape matches the contract expected by Checkout service.
func TestHandleCreateFromCheckout_ResponseShape(t *testing.T) {
	t.Parallel()

	resp := CreateFromCheckoutResponse{
		OrderID:     "some-uuid",
		OrderNumber: "ORD-20260508-000001",
		Status:      string(StatusPendingPayment),
		CreatedAt:   "2026-05-08T00:00:00Z",
	}

	b, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	for _, key := range []string{"orderId", "orderNumber", "status", "createdAt"} {
		if _, ok := m[key]; !ok {
			t.Errorf("response missing key %q", key)
		}
	}
	if m["status"] != "PENDING_PAYMENT" {
		t.Errorf("status: want PENDING_PAYMENT got %v", m["status"])
	}
}

// TestHandleCreateFromCheckout_InternalAuthMiddleware_Rejects verifies that a
// request without the X-Internal-Secret header is rejected with 401 when the
// middleware is applied (unit-level, no pool needed).
func TestHandleCreateFromCheckout_InternalAuthMiddleware_Rejects(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	svc := newTestService(&stubOrders{nextOrderNumber: "ORD-20260508-000001"}, &stubItems{}, &stubStatusHistory{})

	r := gin.New()
	r.POST("/create-from-checkout",
		InternalAuthMiddleware("correct-secret"),
		svc.HandleCreateFromCheckout,
	)

	// No X-Internal-Secret header → 401.
	w := doPost(r, "/create-from-checkout", validCreatePayload())
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 without secret, got %d\nbody: %s", w.Code, w.Body.String())
	}
}

// TestHandleCreateFromCheckout_InternalAuthMiddleware_WrongSecret verifies that
// a wrong shared secret returns 401.
func TestHandleCreateFromCheckout_InternalAuthMiddleware_WrongSecret(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	svc := newTestService(&stubOrders{nextOrderNumber: "ORD-20260508-000001"}, &stubItems{}, &stubStatusHistory{})

	r := gin.New()
	r.POST("/create-from-checkout",
		InternalAuthMiddleware("correct-secret"),
		svc.HandleCreateFromCheckout,
	)

	b, _ := json.Marshal(validCreatePayload())
	req := httptest.NewRequest(http.MethodPost, "/create-from-checkout", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(HeaderInternalSecret, "wrong-secret")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 with wrong secret, got %d\nbody: %s", w.Code, w.Body.String())
	}
}
