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

// fakeListSvc is a test double for the catalogService interface used by ListProducts.
type fakeListSvc struct {
	listResult ListProductsResult
	listErr    error
}

func (f *fakeListSvc) ListProducts(_ context.Context, _ ProductListFilter, _ SortOption, _, _ int) (ListProductsResult, error) {
	return f.listResult, f.listErr
}

func (f *fakeListSvc) GetProduct(_ context.Context, _ uuid.UUID, _ bool) (ProductDetail, error) {
	panic("not used in list tests")
}

func (f *fakeListSvc) CreateProduct(_ context.Context, _ CreateProductInput) (uuid.UUID, error) {
	panic("not used in list tests")
}

func setupGinListHandler(svc catalogService) (*gin.Engine, *Handler) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := &Handler{svc: svc}
	r.POST("/product/list", h.ListProducts)
	return r, h
}

func TestListProducts(t *testing.T) {
	catID := uuid.New()
	prodID := uuid.New()

	okResult := ListProductsResult{
		Items: []ProductListItem{
			{
				ProductID:    prodID,
				SKU:          "SKU-001",
				Name:         "Widget",
				Price:        "9.99",
				CategoryID:   catID,
				CategoryName: "Gadgets",
				CategorySlug: "gadgets",
				ThumbnailURL: "https://cdn.example.com/widget.jpg",
				Status:       string(ProductStatusActive),
			},
		},
		Total: 1,
		Page:  1,
		Limit: 20,
	}

	tests := []struct {
		name       string
		body       any
		svc        catalogService
		wantStatus int
		wantCode   string
	}{
		{
			// CAT-001 + CAT-002: happy path, defaults apply
			name:       "happy_path_defaults",
			body:       map[string]any{},
			svc:        &fakeListSvc{listResult: okResult},
			wantStatus: http.StatusOK,
			wantCode:   "0000",
		},
		{
			// CAT-002: explicit valid pagination and sort
			name:       "explicit_pagination_and_sort",
			body:       map[string]any{"page": 2, "limit": 10, "sort": "price_asc"},
			svc:        &fakeListSvc{listResult: okResult},
			wantStatus: http.StatusOK,
			wantCode:   "0000",
		},
		{
			// CAT-003: categoryId filter (valid UUID)
			name: "category_filter",
			body: map[string]any{
				"filter": map[string]any{"categoryId": catID.String()},
			},
			svc:        &fakeListSvc{listResult: okResult},
			wantStatus: http.StatusOK,
			wantCode:   "0000",
		},
		{
			// CAT-004: price range filter
			name: "price_range_filter",
			body: map[string]any{
				"filter": map[string]any{"minPrice": 5.0, "maxPrice": 50.0},
			},
			svc:        &fakeListSvc{listResult: okResult},
			wantStatus: http.StatusOK,
			wantCode:   "0000",
		},
		{
			// negative: page=0 is coerced to 1 by handler; page < 0 is rejected
			name:       "invalid_page_negative",
			body:       map[string]any{"page": -1},
			svc:        &fakeListSvc{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "HM400",
		},
		{
			// negative: limit out of range
			name:       "invalid_limit_over_100",
			body:       map[string]any{"limit": 101},
			svc:        &fakeListSvc{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "HM400",
		},
		{
			// negative: unknown sort value
			name:       "invalid_sort",
			body:       map[string]any{"sort": "random_order"},
			svc:        &fakeListSvc{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "HM400",
		},
		{
			// negative: invalid categoryId UUID
			name: "invalid_category_uuid",
			body: map[string]any{
				"filter": map[string]any{"categoryId": "not-a-uuid"},
			},
			svc:        &fakeListSvc{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "HM400",
		},
		{
			// negative: maxPrice < minPrice
			name: "max_price_less_than_min_price",
			body: map[string]any{
				"filter": map[string]any{"minPrice": 100.0, "maxPrice": 10.0},
			},
			svc:        &fakeListSvc{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "HM400",
		},
		{
			// negative: service error → 500
			name:       "service_error",
			body:       map[string]any{},
			svc:        &fakeListSvc{listErr: ErrProductNotFound},
			wantStatus: http.StatusInternalServerError,
			wantCode:   "HM500",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := setupGinListHandler(tc.svc)

			body, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/product/list", bytes.NewReader(body))
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
