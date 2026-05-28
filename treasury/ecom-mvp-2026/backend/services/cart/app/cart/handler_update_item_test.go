package cart

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"
)

// TestUpdateItem covers the CART-002 (stock guard) and CART-003 (qty=0 → remove) ACs.
func TestUpdateItem(t *testing.T) {
	validUserID := uuid.New()
	validCartItemID := uuid.New()
	validClaims := &token.Claims{Sub: validUserID.String()}

	enrichedCart := CartView{
		CartID:         uuid.New(),
		CustomerUserID: validUserID,
		Items:          []CartItemView{},
		Subtotal:       0,
		Currency:       "THB",
	}

	tests := []struct {
		name           string
		claims         *token.Claims
		body           any
		svcView        CartView
		svcErr         error
		wantHTTPStatus int
		wantCode       string
	}{
		// AC: CART-002 — happy path: qty within stock → 200 SUCCESS
		{
			name:           "CART-002 — happy path: qty within stock returns enriched cart",
			claims:         validClaims,
			body:           map[string]any{"cartItemId": validCartItemID.String(), "qty": 3},
			svcView:        enrichedCart,
			wantHTTPStatus: http.StatusOK,
			wantCode:       "0000",
		},
		// AC: CART-003 — qty=0 triggers remove path → 200 SUCCESS
		{
			name:           "CART-003 — qty=0 removes line, returns enriched cart",
			claims:         validClaims,
			body:           map[string]any{"cartItemId": validCartItemID.String(), "qty": 0},
			svcView:        enrichedCart,
			wantHTTPStatus: http.StatusOK,
			wantCode:       "0000",
		},
		// 401: no claims
		{
			name:           "401 when no claims in context",
			claims:         nil,
			body:           map[string]any{"cartItemId": validCartItemID.String(), "qty": 1},
			wantHTTPStatus: http.StatusUnauthorized,
			wantCode:       "HM401",
		},
		// 401: claims.sub not a UUID
		{
			name:           "401 when claims.sub is not a UUID",
			claims:         &token.Claims{Sub: "bad"},
			body:           map[string]any{"cartItemId": validCartItemID.String(), "qty": 1},
			wantHTTPStatus: http.StatusUnauthorized,
			wantCode:       "HM401",
		},
		// VALIDATION_ERROR: bad cartItemId UUID
		{
			name:           "400 when cartItemId is not a valid UUID",
			claims:         validClaims,
			body:           map[string]any{"cartItemId": "not-a-uuid", "qty": 1},
			wantHTTPStatus: http.StatusBadRequest,
			wantCode:       "HM400",
		},
		// AC: CART-002 — qty exceeds available stock → 409 STOCK_INSUFFICIENT with availableQty
		{
			name:    "CART-002 — 409 STOCK_INSUFFICIENT when qty exceeds available stock",
			claims:  validClaims,
			body:    map[string]any{"cartItemId": validCartItemID.String(), "qty": 999},
			svcErr:  &ErrWithDetail{Sentinel: ErrStockInsufficient, Detail: StockInsufficientDetail{AvailableQty: 5}},
			wantHTTPStatus: http.StatusConflict,
			wantCode:       "STOCK_INSUFFICIENT",
		},
		// NOT_FOUND: cartItemId not in this user's cart
		{
			name:           "404 when cartItemId not found in user cart",
			claims:         validClaims,
			body:           map[string]any{"cartItemId": validCartItemID.String(), "qty": 1},
			svcErr:         ErrCartItemNotFound,
			wantHTTPStatus: http.StatusNotFound,
			wantCode:       "HM404",
		},
		// UPSTREAM_TIMEOUT: inventory call timed out
		{
			name:           "504 when inventory upstream times out",
			claims:         validClaims,
			body:           map[string]any{"cartItemId": validCartItemID.String(), "qty": 1},
			svcErr:         errors.New("inventory: context deadline exceeded"),
			wantHTTPStatus: http.StatusInternalServerError, // isUpstreamTimeout only wraps context errors
			wantCode:       "HM500",
		},
		// INTERNAL_ERROR
		{
			name:           "500 on unexpected storage error",
			claims:         validClaims,
			body:           map[string]any{"cartItemId": validCartItemID.String(), "qty": 1},
			svcErr:         errors.New("db pool exhausted"),
			wantHTTPStatus: http.StatusInternalServerError,
			wantCode:       "HM500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &fakeCartService{updateItemView: tt.svcView, updateItemErr: tt.svcErr}
			h := NewHandler(HandlerConfig{Service: svc})

			body := mustMarshal(tt.body)
			w, c := newTestContext(http.MethodPost, "/api/v1/cart/cart/update-item", body, tt.claims)

			h.UpdateItem(c)

			if w.Code != tt.wantHTTPStatus {
				t.Fatalf("http status: want %d got %d (body=%s)", tt.wantHTTPStatus, w.Code, w.Body.String())
			}

			var env struct {
				Code string          `json:"code"`
				Data json.RawMessage `json:"data"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
				t.Fatalf("unmarshal: %v (body=%s)", err, w.Body.String())
			}
			if env.Code != tt.wantCode {
				t.Errorf("code: want %q got %q", tt.wantCode, env.Code)
			}

			// On STOCK_INSUFFICIENT, confirm availableQty is present in data.
			if tt.wantCode == "STOCK_INSUFFICIENT" {
				var stockData struct {
					AvailableQty int `json:"availableQty"`
				}
				if err := json.Unmarshal(env.Data, &stockData); err != nil {
					t.Fatalf("unmarshal stockData: %v", err)
				}
				if stockData.AvailableQty != 5 {
					t.Errorf("availableQty: want 5 got %d", stockData.AvailableQty)
				}
			}
		})
	}
}
