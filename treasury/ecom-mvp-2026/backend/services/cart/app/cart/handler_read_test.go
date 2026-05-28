package cart

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"
)

// TestReadCart covers CART-004 (server-side subtotal), CART-005 (availability flags),
// and CART-007 (updatedAt field) acceptance criteria, as well as fan-out error propagation.
func TestReadCart(t *testing.T) {
	validUserID := uuid.New()
	validClaims := &token.Claims{Sub: validUserID.String()}
	updatedAt := time.Date(2026, 5, 8, 10, 0, 0, 0, time.UTC)

	activeProductID := uuid.New()
	inactiveProductID := uuid.New()

	// Cart with one active and one inactive item; subtotal covers only active checkoutable line.
	// active: price=590, qty=2 → lineSubtotal=1180, checkoutable=true
	// inactive: price=300, qty=1 → lineSubtotal=300, checkoutable=false, excluded from subtotal
	enrichedCart := CartView{
		CartID:         uuid.New(),
		CustomerUserID: validUserID,
		Items: []CartItemView{
			{
				CartItemID:    uuid.New(),
				ProductID:     activeProductID,
				SKU:           "SHIRT-RED-M",
				Name:          "Red Cotton Shirt (M)",
				CurrentPrice:  590,
				Qty:           2,
				LineSubtotal:  1180,
				AvailableQty:  8,
				ProductStatus: ProductStatusActive,
				Checkoutable:  true,
			},
			{
				CartItemID:    uuid.New(),
				ProductID:     inactiveProductID,
				SKU:           "OLD-ITEM",
				Name:          "Old Item",
				CurrentPrice:  300,
				Qty:           1,
				LineSubtotal:  300,
				AvailableQty:  0,
				ProductStatus: ProductStatusInactive,
				Checkoutable:  false,
			},
		},
		Subtotal:  1180, // AC CART-004: only checkoutable lines
		Currency:  "THB",
		UpdatedAt: updatedAt, // AC CART-007
	}

	emptyCart := CartView{
		CustomerUserID: validUserID,
		Items:          []CartItemView{},
		Subtotal:       0,
		Currency:       "THB",
	}

	tests := []struct {
		name           string
		claims         *token.Claims
		svcView        CartView
		svcErr         error
		wantHTTPStatus int
		wantCode       string
		checkData      func(t *testing.T, data json.RawMessage)
	}{
		// AC: CART-004, CART-005, CART-007 — happy path with enriched items
		{
			name:           "CART-004/005/007 — returns enriched cart with availability flags and updatedAt",
			claims:         validClaims,
			svcView:        enrichedCart,
			wantHTTPStatus: http.StatusOK,
			wantCode:       "0000",
			checkData: func(t *testing.T, data json.RawMessage) {
				t.Helper()
				var d struct {
					Cart struct {
						Items     []CartItemView `json:"items"`
						Subtotal  float64        `json:"subtotal"`
						Currency  string         `json:"currency"`
						UpdatedAt time.Time      `json:"updatedAt"`
					} `json:"cart"`
				}
				if err := json.Unmarshal(data, &d); err != nil {
					t.Fatalf("unmarshal data: %v", err)
				}
				// CART-004: subtotal is only checkoutable lines
				if d.Cart.Subtotal != 1180 {
					t.Errorf("subtotal: want 1180 got %v", d.Cart.Subtotal)
				}
				// CART-007: updatedAt is present
				if d.Cart.UpdatedAt.IsZero() {
					t.Error("updatedAt should not be zero")
				}
				if len(d.Cart.Items) != 2 {
					t.Fatalf("items count: want 2 got %d", len(d.Cart.Items))
				}
				// CART-005: inactive product is not checkoutable
				inactiveItem := d.Cart.Items[1]
				if inactiveItem.Checkoutable {
					t.Error("inactive product item should have checkoutable=false")
				}
				// CART-004: active item is checkoutable
				activeItem := d.Cart.Items[0]
				if !activeItem.Checkoutable {
					t.Error("active product item should have checkoutable=true")
				}
			},
		},
		// Empty cart: no items, subtotal=0
		{
			name:           "empty cart returns items:[] subtotal:0",
			claims:         validClaims,
			svcView:        emptyCart,
			wantHTTPStatus: http.StatusOK,
			wantCode:       "0000",
			checkData: func(t *testing.T, data json.RawMessage) {
				t.Helper()
				var d struct {
					Cart struct {
						Items    []CartItemView `json:"items"`
						Subtotal float64        `json:"subtotal"`
					} `json:"cart"`
				}
				if err := json.Unmarshal(data, &d); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				if len(d.Cart.Items) != 0 {
					t.Errorf("items: want [] got %v", d.Cart.Items)
				}
				if d.Cart.Subtotal != 0 {
					t.Errorf("subtotal: want 0 got %v", d.Cart.Subtotal)
				}
			},
		},
		// 401: no claims
		{
			name:           "401 when no claims in context",
			claims:         nil,
			wantHTTPStatus: http.StatusUnauthorized,
			wantCode:       "HM401",
		},
		// 401: bad UUID in claims.sub
		{
			name:           "401 when claims.sub is not a valid UUID",
			claims:         &token.Claims{Sub: "bad-uuid"},
			wantHTTPStatus: http.StatusUnauthorized,
			wantCode:       "HM401",
		},
		// UPSTREAM_TIMEOUT from fan-out: any fan-out error surfaces as 500 INTERNAL_ERROR
		// (AMB-003 resolution: strict policy — whole read fails on any upstream error)
		{
			name:           "500 when fan-out returns error (catalog or inventory upstream failure)",
			claims:         validClaims,
			svcErr:         fmt.Errorf("cart.read fan-out: %w", context.DeadlineExceeded),
			wantHTTPStatus: http.StatusInternalServerError,
			wantCode:       "HM500",
		},
		// 500: unexpected DB error
		{
			name:           "500 on DB error fetching cart",
			claims:         validClaims,
			svcErr:         fmt.Errorf("pgx: connection refused"),
			wantHTTPStatus: http.StatusInternalServerError,
			wantCode:       "HM500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &fakeCartService{readCartView: tt.svcView, readCartErr: tt.svcErr}
			h := NewHandler(HandlerConfig{Service: svc})

			w, c := newTestContext(http.MethodPost, "/api/v1/cart/cart/read", nil, tt.claims)

			h.ReadCart(c)

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
			if tt.checkData != nil && tt.wantHTTPStatus == http.StatusOK {
				tt.checkData(t, env.Data)
			}
		})
	}
}
