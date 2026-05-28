package catalog

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// fakeDetailSvc is a test double for the catalogService interface used by GetProduct.
type fakeDetailSvc struct {
	detail    ProductDetail
	detailErr error
}

func (f *fakeDetailSvc) ListProducts(_ context.Context, _ ProductListFilter, _ SortOption, _, _ int) (ListProductsResult, error) {
	panic("not used in detail tests")
}

func (f *fakeDetailSvc) GetProduct(_ context.Context, _ uuid.UUID, _ bool) (ProductDetail, error) {
	return f.detail, f.detailErr
}

func (f *fakeDetailSvc) CreateProduct(_ context.Context, _ CreateProductInput) (uuid.UUID, error) {
	panic("not used in detail tests")
}

func setupGinDetailHandler(svc catalogService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := &Handler{svc: svc}
	r.POST("/product/detail", h.GetProduct)
	return r
}

func TestGetProduct(t *testing.T) {
	prodID := uuid.New()
	catID := uuid.New()

	okDetail := ProductDetail{
		ProductID:    prodID,
		SKU:          "SKU-001",
		Name:         "Widget",
		Description:  "A fine widget",
		Images:       []string{"https://cdn.example.com/w.jpg"},
		Price:        "9.99",
		CategoryID:   catID,
		CategoryName: "Gadgets",
		CategorySlug: "gadgets",
		Status:       string(ProductStatusActive),
	}

	tests := []struct {
		name       string
		body       any
		svc        catalogService
		wantStatus int
		wantCode   string
	}{
		{
			// CAT-006: happy path — public caller gets ACTIVE product
			name:       "happy_path_active_product",
			body:       map[string]any{"productId": prodID.String()},
			svc:        &fakeDetailSvc{detail: okDetail},
			wantStatus: http.StatusOK,
			wantCode:   "0000",
		},
		{
			// CAT-006 negative: product not found → 404
			name:       "product_not_found",
			body:       map[string]any{"productId": prodID.String()},
			svc:        &fakeDetailSvc{detailErr: ErrProductNotFound},
			wantStatus: http.StatusNotFound,
			wantCode:   "HM404",
		},
		{
			// negative: productId missing → ShouldBindJSON required tag fails → 400
			name:       "missing_product_id",
			body:       map[string]any{},
			svc:        &fakeDetailSvc{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "HM400",
		},
		{
			// negative: productId not a valid UUID
			name:       "invalid_product_id_uuid",
			body:       map[string]any{"productId": "not-a-uuid"},
			svc:        &fakeDetailSvc{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "HM400",
		},
		{
			// negative: service error other than not-found → 500
			name:       "service_internal_error",
			body:       map[string]any{"productId": prodID.String()},
			svc:        &fakeDetailSvc{detailErr: ErrDuplicateSKU},
			wantStatus: http.StatusInternalServerError,
			wantCode:   "HM500",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := setupGinDetailHandler(tc.svc)

			body, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/product/detail", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			require.Equal(t, tc.wantStatus, w.Code, "body: %s", w.Body.String())

			var resp map[string]any
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
			require.Equal(t, tc.wantCode, resp["code"], "body: %s", w.Body.String())
		})
	}
}
