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
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"
)

// fakeCreateSvc is a test double for CreateProduct.
type fakeCreateSvc struct {
	productID uuid.UUID
	createErr error
}

func (f *fakeCreateSvc) ListProducts(_ context.Context, _ ProductListFilter, _ SortOption, _, _ int) (ListProductsResult, error) {
	panic("not used in create tests")
}

func (f *fakeCreateSvc) GetProduct(_ context.Context, _ uuid.UUID, _ bool) (ProductDetail, error) {
	panic("not used in create tests")
}

func (f *fakeCreateSvc) CreateProduct(_ context.Context, _ CreateProductInput) (uuid.UUID, error) {
	return f.productID, f.createErr
}

// adminClaims builds a *token.Claims with role=ADMIN and a valid UUID sub.
func adminClaims(sub string) *token.Claims {
	return &token.Claims{
		Sub: sub,
		Extra: map[string]any{
			"role": "ADMIN",
		},
	}
}

// setupGinCreateHandler registers POST /product/create and injects JWT claims
// into the gin context via a before-handler middleware, simulating the JWT
// middleware that would normally run in production.
func setupGinCreateHandler(svc catalogService, claims *token.Claims) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := &Handler{svc: svc}

	r.POST("/product/create", func(c *gin.Context) {
		if claims != nil {
			c.Request = c.Request.WithContext(token.WithClaims(c.Request.Context(), claims))
		}
		h.CreateProduct(c)
	})
	return r
}

func TestCreateProduct(t *testing.T) {
	actorID := uuid.New()
	catID := uuid.New()
	newProdID := uuid.New()

	validBody := map[string]any{
		"sku":        "SKU-001",
		"name":       "Widget",
		"price":      9.99,
		"categoryId": catID.String(),
		"status":     "DRAFT",
	}

	tests := []struct {
		name       string
		body       any
		claims     *token.Claims
		svc        catalogService
		wantStatus int
		wantCode   string
	}{
		{
			// CAT-007: happy path — ADMIN creates a DRAFT product
			name:       "happy_path_create_draft",
			body:       validBody,
			claims:     adminClaims(actorID.String()),
			svc:        &fakeCreateSvc{productID: newProdID},
			wantStatus: http.StatusCreated,
			wantCode:   "CREATED",
		},
		{
			// CAT-007: create with ACTIVE status
			name: "happy_path_create_active",
			body: map[string]any{
				"sku": "SKU-002", "name": "Gadget",
				"price": 19.99, "categoryId": catID.String(), "status": "ACTIVE",
			},
			claims:     adminClaims(actorID.String()),
			svc:        &fakeCreateSvc{productID: newProdID},
			wantStatus: http.StatusCreated,
			wantCode:   "CREATED",
		},
		{
			// CAT-007: with images array
			name: "with_images",
			body: map[string]any{
				"sku": "SKU-003", "name": "Item",
				"price": 5.0, "categoryId": catID.String(),
				"images": []string{"https://cdn.example.com/img1.jpg"},
			},
			claims:     adminClaims(actorID.String()),
			svc:        &fakeCreateSvc{productID: newProdID},
			wantStatus: http.StatusCreated,
			wantCode:   "CREATED",
		},
		{
			// auth: no claims in context → 401
			name:       "no_jwt_claims",
			body:       validBody,
			claims:     nil,
			svc:        &fakeCreateSvc{},
			wantStatus: http.StatusUnauthorized,
			wantCode:   "HM401",
		},
		{
			// auth: non-ADMIN role → 403
			name: "non_admin_role",
			body: validBody,
			claims: &token.Claims{
				Sub:   actorID.String(),
				Extra: map[string]any{"role": "CUSTOMER"},
			},
			svc:        &fakeCreateSvc{},
			wantStatus: http.StatusForbidden,
			wantCode:   "HM403",
		},
		{
			// auth: claims.Sub not a valid UUID → 401
			name:       "invalid_sub_uuid",
			body:       validBody,
			claims:     adminClaims("not-a-uuid"),
			svc:        &fakeCreateSvc{},
			wantStatus: http.StatusUnauthorized,
			wantCode:   "HM401",
		},
		{
			// validation: sku empty → 400
			name: "empty_sku",
			body: map[string]any{
				"sku": "", "name": "Widget",
				"price": 9.99, "categoryId": catID.String(),
			},
			claims:     adminClaims(actorID.String()),
			svc:        &fakeCreateSvc{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "HM400",
		},
		{
			// validation: name empty → 400
			name: "empty_name",
			body: map[string]any{
				"sku": "SKU-X", "name": "",
				"price": 9.99, "categoryId": catID.String(),
			},
			claims:     adminClaims(actorID.String()),
			svc:        &fakeCreateSvc{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "HM400",
		},
		{
			// validation: price zero → 400
			name: "zero_price",
			body: map[string]any{
				"sku": "SKU-X", "name": "Widget",
				"price": 0.0, "categoryId": catID.String(),
			},
			claims:     adminClaims(actorID.String()),
			svc:        &fakeCreateSvc{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "HM400",
		},
		{
			// validation: invalid status → 400
			name: "invalid_status",
			body: map[string]any{
				"sku": "SKU-X", "name": "Widget",
				"price": 9.99, "categoryId": catID.String(), "status": "DELETED",
			},
			claims:     adminClaims(actorID.String()),
			svc:        &fakeCreateSvc{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "HM400",
		},
		{
			// validation: categoryId not a UUID → 400
			name: "invalid_category_uuid",
			body: map[string]any{
				"sku": "SKU-X", "name": "Widget",
				"price": 9.99, "categoryId": "bad-uuid",
			},
			claims:     adminClaims(actorID.String()),
			svc:        &fakeCreateSvc{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "HM400",
		},
		{
			// validation: images[1] is empty string → 400
			name: "empty_image_url",
			body: map[string]any{
				"sku": "SKU-X", "name": "Widget",
				"price": 9.99, "categoryId": catID.String(),
				"images": []string{"https://valid.example.com/img.jpg", ""},
			},
			claims:     adminClaims(actorID.String()),
			svc:        &fakeCreateSvc{},
			wantStatus: http.StatusBadRequest,
			wantCode:   "HM400",
		},
		{
			// service: duplicate SKU → 409
			name:       "duplicate_sku",
			body:       validBody,
			claims:     adminClaims(actorID.String()),
			svc:        &fakeCreateSvc{createErr: ErrDuplicateSKU},
			wantStatus: http.StatusConflict,
			wantCode:   "DUPLICATE_SKU",
		},
		{
			// service: category not found → 404
			name:       "category_not_found",
			body:       validBody,
			claims:     adminClaims(actorID.String()),
			svc:        &fakeCreateSvc{createErr: ErrCategoryNotFound},
			wantStatus: http.StatusNotFound,
			wantCode:   "HM404",
		},
		{
			// service: internal error → 500
			name:       "service_internal_error",
			body:       validBody,
			claims:     adminClaims(actorID.String()),
			svc:        &fakeCreateSvc{createErr: ErrValidation},
			wantStatus: http.StatusInternalServerError,
			wantCode:   "HM500",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := setupGinCreateHandler(tc.svc, tc.claims)

			body, err := json.Marshal(tc.body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/product/create", bytes.NewReader(body))
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
