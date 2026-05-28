package product_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/product"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/product/access"
	access_mocks "gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/product/access/mocks"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListProducts(t *testing.T) {
	r := require.New(t)

	orgID := uuid.New()
	productID1 := uuid.New()
	productID2 := uuid.New()
	now := time.Now()

	type mockArgs struct {
		productStorage *access_mocks.ProductStorageMock
	}

	type args struct {
		ctx context.Context
		req product.ListProductsRequest
	}

	type want struct {
		err          bool
		code         wrapper.Code
		Message      wrapper.Message
		productCount int
	}

	tests := []struct {
		name    string
		prepare func(m mockArgs, args args)
		args    args
		want    want
	}{
		{
			name: "success, case valid request with products",
			prepare: func(m mockArgs, args args) {
				m.productStorage.
					EXPECT().
					ListProducts(args.ctx, args.req.OrganizationID).
					Return([]access.Product{
						{
							ProductID:      productID1.String(),
							Name:           "Product 1",
							Description:    "Description 1",
							Price:          10.00,
							OrganizationID: orgID.String(),
							Status:         access.ProductStatusActive,
							CreatedAt:      now,
							UpdatedAt:      now,
						},
						{
							ProductID:      productID2.String(),
							Name:           "Product 2",
							Description:    "Description 2",
							Price:          20.00,
							OrganizationID: orgID.String(),
							Status:         access.ProductStatusActive,
							CreatedAt:      now,
							UpdatedAt:      now,
						},
					}, nil)
			},
			args: args{
				req: product.ListProductsRequest{
					OrganizationID: orgID,
				},
			},
			want: want{
				err:          false,
				code:         app.CodeSuccess,
				Message:      app.MessageSuccess,
				productCount: 2,
			},
		},
		{
			name: "success, case products with invalid ID are skipped",
			prepare: func(m mockArgs, args args) {
				m.productStorage.
					EXPECT().
					ListProducts(args.ctx, args.req.OrganizationID).
					Return([]access.Product{
						{
							ProductID:      productID1.String(),
							Name:           "Valid Product",
							OrganizationID: orgID.String(),
							Status:         access.ProductStatusActive,
							CreatedAt:      now,
							UpdatedAt:      now,
						},
						{
							ProductID:      "not-a-valid-uuid",
							Name:           "Invalid Product",
							OrganizationID: orgID.String(),
							Status:         access.ProductStatusActive,
							CreatedAt:      now,
							UpdatedAt:      now,
						},
					}, nil)
			},
			args: args{
				req: product.ListProductsRequest{
					OrganizationID: orgID,
				},
			},
			want: want{
				err:          false,
				code:         app.CodeSuccess,
				Message:      app.MessageSuccess,
				productCount: 1,
			},
		},
		{
			name: "fail, case invalid request body - missing organization ID",
			prepare: func(m mockArgs, args args) {
				// no calls
			},
			args: args{
				req: product.ListProductsRequest{},
			},
			want: want{
				err:     true,
				code:    app.CodeBadRequest,
				Message: app.MessageBadRequest,
			},
		},
		{
			name: "fail, case ListProducts storage error",
			prepare: func(m mockArgs, args args) {
				m.productStorage.
					EXPECT().
					ListProducts(args.ctx, args.req.OrganizationID).
					Return(nil, assert.AnError)
			},
			args: args{
				req: product.ListProductsRequest{
					OrganizationID: orgID,
				},
			},
			want: want{
				err:     true,
				code:    app.CodeInternalError,
				Message: app.MessageInternalError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			// Arrange
			var payload bytes.Buffer
			json.NewEncoder(&payload).Encode(tt.args.req)

			req := httptest.NewRequest(http.MethodPost, "http://0.0.0.0/api/v1/product/list", &payload)
			req.Header.Set("Content-Type", "application/json")

			ctx.Request = req

			m := mockArgs{
				productStorage: access_mocks.NewProductStorageMock(t),
			}

			if tt.prepare != nil {
				tt.args.ctx = ctx.Request.Context()
				tt.prepare(m, tt.args)
			}

			h := product.NewHandler(product.HandlerConfig{
				ProductStorage: m.productStorage,
			})

			// Act
			h.ListProducts(ctx)

			// Assert
			var resp wrapper.ResponseOption[product.ListProductsResponse]
			json.NewDecoder(w.Body).Decode(&resp)

			if tt.want.err {
				r.NotEqual(http.StatusOK, w.Code)
				r.Equal(tt.want.code, resp.Code)
				r.Equal(tt.want.Message, resp.Message)
			} else {
				r.Equal(http.StatusOK, w.Code)
				r.Equal(tt.want.code, resp.Code)
				r.Equal(tt.want.Message, resp.Message)
				r.NotNil(resp.Data)
				r.Len(resp.Data.Products, tt.want.productCount)
			}
		})
	}
}
