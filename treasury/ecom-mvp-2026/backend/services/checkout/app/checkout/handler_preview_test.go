package checkout

// handler_preview_test.go — table-driven coverage for the pre-client gates
// of Service.Preview plus the small pure helpers (stockBySkU, findAddress,
// pricing.Compute*).
//
// As with handler_commit_test.go, the full read-only orchestration (cart +
// identity + catalog fanout + inventory) is integration-test territory.
// This file pins the pre-call validation surface — the gates that map to
// CHK-005 / CHK-AMBIG-002 / auth/role enforcement — and the pricing-tier
// boundary (BA-PR-002) which is the single source per pricing.go.

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newPreviewRequest(t *testing.T, body string, role string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/checkout/checkout/preview", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	if role != "" {
		// Reuse the helper from handler_commit_test.go (same package).
		return newRequestWithClaims(t, body, role, "n/a-not-required-for-preview")
	}
	return req
}

// -----------------------------------------------------------------------------
// 1. Pre-call validation surface
// -----------------------------------------------------------------------------

func TestPreview_PreCallValidation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		role       string
		body       string
		wantStatus int
		wantCode   string
	}{
		{
			name:       "AUTH_MISSING when no claims in context",
			role:       "",
			body:       `{"shippingAddressId":"addr-1"}`,
			wantStatus: http.StatusUnauthorized,
			wantCode:   string(CodeAuthMissing),
		},
		{
			name:       "AUTH_FORBIDDEN when role is not CUSTOMER",
			role:       "ADMIN",
			body:       `{"shippingAddressId":"addr-1"}`,
			wantStatus: http.StatusForbidden,
			wantCode:   string(CodeAuthForbidden),
		},
		{
			name:       "VALIDATION_ERROR when body is malformed JSON",
			role:       "CUSTOMER",
			body:       `{not-json}`,
			wantStatus: http.StatusBadRequest,
			wantCode:   string(CodeValidationError),
		},
		{
			name:       "VALIDATION_ERROR when body carries forbidden 'subtotal' (CHK-005)",
			role:       "CUSTOMER",
			body:       `{"shippingAddressId":"addr-1","subtotal":1}`,
			wantStatus: http.StatusBadRequest,
			wantCode:   string(CodeValidationError),
		},
		{
			name:       "VALIDATION_ERROR when body carries forbidden 'couponCode' (CHK-AMBIG-002)",
			role:       "CUSTOMER",
			body:       `{"shippingAddressId":"addr-1","couponCode":"X"}`,
			wantStatus: http.StatusBadRequest,
			wantCode:   string(CodeValidationError),
		},
		{
			name:       "VALIDATION_ERROR when shippingAddressId missing",
			role:       "CUSTOMER",
			body:       `{}`,
			wantStatus: http.StatusBadRequest,
			wantCode:   string(CodeValidationError),
		},
		{
			name:       "VALIDATION_ERROR when shippingAddressId empty",
			role:       "CUSTOMER",
			body:       `{"shippingAddressId":""}`,
			wantStatus: http.StatusBadRequest,
			wantCode:   string(CodeValidationError),
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			w := runPreviewBare(t, newPreviewRequest(t, tc.body, tc.role))
			if w.Code != tc.wantStatus {
				t.Fatalf("status: want %d got %d body=%s", tc.wantStatus, w.Code, w.Body.String())
			}
			env := decodeEnv(t, w.Body.Bytes())
			if env.Code != tc.wantCode {
				t.Errorf("code: want %q got %q body=%s", tc.wantCode, env.Code, w.Body.String())
			}
		})
	}
}

// -----------------------------------------------------------------------------
// 2. Pricing tier (§11.2 / BA-PR-002) — single source = pricing.go
// -----------------------------------------------------------------------------

// TestPricingTier_BAPR002 pins the boundary-inclusive 1500-THB rule.
// STORY_CHECKOUT_SHIPPING_FEE.AC1 + AC2 in td.json §ba_acceptance_mapping.
func TestPricingTier_BAPR002(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name             string
		subtotalMinor    int64
		wantShippingFee  int64
		wantGrandTotal   int64
		acceptanceCriID  string
	}{
		{
			name:            "subtotal=0 → fee=60 (lower bound)",
			subtotalMinor:   0,
			wantShippingFee: 60_00,
			wantGrandTotal:  60_00,
			acceptanceCriID: "BA-PR-002",
		},
		{
			name:            "subtotal=149_900 (1499 THB) → fee=60, total=1559 THB",
			subtotalMinor:   1499_00,
			wantShippingFee: 60_00,
			wantGrandTotal:  1559_00,
			acceptanceCriID: "STORY_CHECKOUT_SHIPPING_FEE.AC1",
		},
		{
			name:            "subtotal=150_000 (1500 THB inclusive) → fee=0, total=1500 THB",
			subtotalMinor:   1500_00,
			wantShippingFee: 0,
			wantGrandTotal:  1500_00,
			acceptanceCriID: "STORY_CHECKOUT_SHIPPING_FEE.AC2 (BA-PR-002 inclusive)",
		},
		{
			name:            "subtotal=200_000 (2000 THB) → fee=0",
			subtotalMinor:   2000_00,
			wantShippingFee: 0,
			wantGrandTotal:  2000_00,
			acceptanceCriID: "BA-PR-002",
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			gotFee := ComputeShippingFee(c.subtotalMinor)
			if gotFee != c.wantShippingFee {
				t.Errorf("[%s] shippingFee: want %d got %d", c.acceptanceCriID, c.wantShippingFee, gotFee)
			}
			gotTotal := ComputeGrandTotal(c.subtotalMinor, gotFee, 0)
			if gotTotal != c.wantGrandTotal {
				t.Errorf("[%s] grandTotal: want %d got %d", c.acceptanceCriID, c.wantGrandTotal, gotTotal)
			}
		})
	}
}

// -----------------------------------------------------------------------------
// 3. Helper purity
// -----------------------------------------------------------------------------

func TestStockBySkU(t *testing.T) {
	t.Parallel()
	in := []StockLevel{
		{SKU: "A", AvailableQty: 1},
		{SKU: "B", AvailableQty: 2},
	}
	got := stockBySkU(in)
	if got["A"].AvailableQty != 1 || got["B"].AvailableQty != 2 {
		t.Errorf("stockBySkU dropped rows: %+v", got)
	}
}

func TestFindAddress(t *testing.T) {
	t.Parallel()
	addrs := []AddressSnapshot{
		{AddressID: "a1", IsComplete: true},
		{AddressID: "a2", IsComplete: false},
	}
	a, ok := findAddress(addrs, "a2")
	if !ok || a.AddressID != "a2" {
		t.Errorf("findAddress(a2) miss: %v ok=%v", a, ok)
	}
	_, missingOk := findAddress(addrs, "nope")
	if missingOk {
		t.Errorf("findAddress(nope) should miss")
	}
}

// TestPreview_EmptyCartShape is a structural lock on the empty-cart
// short-circuit envelope. We can't reach the cart client without spinning
// up a real cart httptest server, but we can lock the shape of the response
// the handler emits when cart-read returns no items: items=[], blockers=[],
// code=0000. Done by directly building the wire envelope and round-tripping.
func TestPreview_EmptyCartEnvelopeShape(t *testing.T) {
	t.Parallel()
	resp := PreviewResponse{
		Items:    []PreviewItem{},
		Blockers: []string{},
	}
	b, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !bytes.Contains(b, []byte(`"items":[]`)) {
		t.Errorf("items must serialise as [] not null: %s", string(b))
	}
	if !bytes.Contains(b, []byte(`"blockers":[]`)) {
		t.Errorf("blockers must serialise as [] not null: %s", string(b))
	}
	gin.SetMode(gin.TestMode) // touch gin to avoid an unused-import after refactors
}
