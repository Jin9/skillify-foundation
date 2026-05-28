package cart

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"
)

// fakeCartService is a minimal in-package stub satisfying CartService.
// Only the method under test needs a real implementation; all others return zero values.
type fakeCartService struct {
	addItemView CartView
	addItemErr  error

	updateItemView CartView
	updateItemErr  error

	readCartView CartView
	readCartErr  error

	removeItemView CartView
	removeItemErr  error

	clearOnCheckoutRemoved int
	clearOnCheckoutErr     error
}

func (f *fakeCartService) AddItem(_ context.Context, _, _ uuid.UUID, _ int) (CartView, error) {
	return f.addItemView, f.addItemErr
}
func (f *fakeCartService) UpdateItem(_ context.Context, _, _ uuid.UUID, _ int) (CartView, error) {
	return f.updateItemView, f.updateItemErr
}
func (f *fakeCartService) RemoveItem(_ context.Context, _, _ uuid.UUID) (CartView, error) {
	return f.removeItemView, f.removeItemErr
}
func (f *fakeCartService) ReadCart(_ context.Context, _ uuid.UUID) (CartView, error) {
	return f.readCartView, f.readCartErr
}
func (f *fakeCartService) ClearOnCheckout(_ context.Context, _ uuid.UUID, _ []uuid.UUID, _ uuid.UUID) (int, error) {
	return f.clearOnCheckoutRemoved, f.clearOnCheckoutErr
}

// --- helpers -----------------------------------------------------------------

func newTestContext(method, path string, body []byte, claims *token.Claims) (*httptest.ResponseRecorder, *gin.Context) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if claims != nil {
		req = req.WithContext(token.WithClaims(req.Context(), claims))
	}
	c.Request = req
	return w, c
}

func mustMarshal(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

// --- add-item tests ----------------------------------------------------------

func TestAddItem(t *testing.T) {
	validUserID := uuid.New()
	validProductID := uuid.New()
	validClaims := &token.Claims{Sub: validUserID.String()}

	enrichedCart := CartView{
		CartID:         uuid.New(),
		CustomerUserID: validUserID,
		Items: []CartItemView{
			{
				CartItemID:    uuid.New(),
				ProductID:     validProductID,
				SKU:           "SHIRT-RED-M",
				Name:          "Red Cotton Shirt (M)",
				CurrentPrice:  590,
				Qty:           2,
				LineSubtotal:  1180,
				AvailableQty:  8,
				ProductStatus: ProductStatusActive,
				Checkoutable:  true,
			},
		},
		Subtotal: 1180,
		Currency: "THB",
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
		// AC: CART-001 — authenticated customer adds active product → 200 SUCCESS
		{
			name:           "CART-001 — happy path: new SKU creates line, returns enriched cart",
			claims:         validClaims,
			body:           map[string]any{"productId": validProductID.String(), "qty": 2},
			svcView:        enrichedCart,
			wantHTTPStatus: http.StatusOK,
			wantCode:       "0000",
		},
		// AC: CART-001 — unauthenticated → 401
		{
			name:           "401 when no claims in context",
			claims:         nil,
			body:           map[string]any{"productId": validProductID.String(), "qty": 1},
			wantHTTPStatus: http.StatusUnauthorized,
			wantCode:       "HM401",
		},
		// AC: CART-001 — malformed claims.Sub → 401
		{
			name:           "401 when claims.sub is not a UUID",
			claims:         &token.Claims{Sub: "not-a-uuid"},
			body:           map[string]any{"productId": validProductID.String(), "qty": 1},
			wantHTTPStatus: http.StatusUnauthorized,
			wantCode:       "HM401",
		},
		// VALIDATION_ERROR: invalid productId UUID
		{
			name:           "400 when productId is not a valid UUID",
			claims:         validClaims,
			body:           map[string]any{"productId": "bad-uuid", "qty": 1},
			wantHTTPStatus: http.StatusBadRequest,
			wantCode:       "HM400",
		},
		// VALIDATION_ERROR: qty=0 (min=1 per spec)
		{
			name:           "400 when qty is zero (below minimum)",
			claims:         validClaims,
			body:           map[string]any{"productId": validProductID.String(), "qty": 0},
			wantHTTPStatus: http.StatusBadRequest,
			wantCode:       "HM400",
		},
		// PRODUCT_INACTIVE (CART-001 rejection)
		{
			name:           "422 when product is INACTIVE",
			claims:         validClaims,
			body:           map[string]any{"productId": validProductID.String(), "qty": 1},
			svcErr:         ErrProductInactive,
			wantHTTPStatus: http.StatusUnprocessableEntity,
			wantCode:       "PRODUCT_INACTIVE",
		},
		// PRODUCT_DELETED (CART-001 rejection)
		{
			name:           "422 when product is DELETED",
			claims:         validClaims,
			body:           map[string]any{"productId": validProductID.String(), "qty": 1},
			svcErr:         ErrProductDeleted,
			wantHTTPStatus: http.StatusUnprocessableEntity,
			wantCode:       "PRODUCT_DELETED",
		},
		// NOT_FOUND: catalog returned 404
		{
			name:           "404 when product not found in catalog",
			claims:         validClaims,
			body:           map[string]any{"productId": validProductID.String(), "qty": 1},
			svcErr:         ErrProductNotFound,
			wantHTTPStatus: http.StatusNotFound,
			wantCode:       "HM404",
		},
		// INTERNAL_ERROR: unexpected storage failure
		{
			name:           "500 on unexpected storage error",
			claims:         validClaims,
			body:           map[string]any{"productId": validProductID.String(), "qty": 1},
			svcErr:         errors.New("connection refused"),
			wantHTTPStatus: http.StatusInternalServerError,
			wantCode:       "HM500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &fakeCartService{addItemView: tt.svcView, addItemErr: tt.svcErr}
			h := NewHandler(HandlerConfig{Service: svc})

			body := mustMarshal(tt.body)
			w, c := newTestContext(http.MethodPost, "/api/v1/cart/cart/add-item", body, tt.claims)

			h.AddItem(c)

			if w.Code != tt.wantHTTPStatus {
				t.Fatalf("http status: want %d got %d (body=%s)", tt.wantHTTPStatus, w.Code, w.Body.String())
			}

			var env struct {
				Code string `json:"code"`
				Data *struct {
					Cart struct {
						Items    []CartItemView `json:"items"`
						Subtotal float64        `json:"subtotal"`
					} `json:"cart"`
				} `json:"data"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
				t.Fatalf("unmarshal: %v (body=%s)", err, w.Body.String())
			}
			if env.Code != tt.wantCode {
				t.Errorf("code: want %q got %q", tt.wantCode, env.Code)
			}
			if tt.wantHTTPStatus == http.StatusOK {
				if env.Data == nil {
					t.Fatal("expected data on success")
				}
				if len(env.Data.Cart.Items) != 1 {
					t.Errorf("items count: want 1 got %d", len(env.Data.Cart.Items))
				}
				if env.Data.Cart.Subtotal != enrichedCart.Subtotal {
					t.Errorf("subtotal: want %v got %v", enrichedCart.Subtotal, env.Data.Cart.Subtotal)
				}
			}
		})
	}
}
