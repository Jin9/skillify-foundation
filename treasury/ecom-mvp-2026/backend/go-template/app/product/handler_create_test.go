package product_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/product"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/product/access"
	access_mocks "gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/product/access/mocks"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateProduct(t *testing.T) {
	r := require.New(t)

	productID := uuid.New()

	type mockArgs struct {
		productStorage *access_mocks.ProductStorageMock
	}

	type args struct {
		ctx context.Context
		req product.CreateProductRequest
	}

	type want struct {
		err     bool
		code    wrapper.Code
		Message wrapper.Message
	}

	tests := []struct {
		name    string
		prepare func(m mockArgs, args args)
		args    args
		want    want
	}{
		{
			name: "success, case valid request",
			prepare: func(m mockArgs, args args) {
				m.productStorage.
					EXPECT().
					CreateProduct(args.ctx, mock.MatchedBy(func(p access.Product) bool {
						return p.Name == args.req.Name &&
							p.Description == args.req.Description &&
							p.Price == args.req.Price &&
							p.OrganizationID == args.req.OrganizationID &&
							p.Status == access.ProductStatusActive
					})).
					Return(access.Product{
						ProductID: productID.String(),
					}, nil)
			},
			args: args{
				req: product.CreateProductRequest{
					Name:           "Test Product",
					Description:    "A test product",
					Price:          99.99,
					OrganizationID: uuid.New().String(),
				},
			},
			want: want{
				err:     false,
				code:    app.CodeSuccess,
				Message: app.MessageSuccess,
			},
		},
		{
			name: "fail, case invalid request body - missing name",
			prepare: func(m mockArgs, args args) {
				// no calls
			},
			args: args{
				req: product.CreateProductRequest{
					Price:          99.99,
					OrganizationID: uuid.New().String(),
				},
			},
			want: want{
				err:     true,
				code:    app.CodeBadRequest,
				Message: app.MessageBadRequest,
			},
		},
		{
			name: "fail, case CreateProduct storage error",
			prepare: func(m mockArgs, args args) {
				m.productStorage.
					EXPECT().
					CreateProduct(args.ctx, mock.Anything).
					Return(access.Product{}, assert.AnError)
			},
			args: args{
				req: product.CreateProductRequest{
					Name:           "Test Product",
					Description:    "A test product",
					Price:          99.99,
					OrganizationID: uuid.New().String(),
				},
			},
			want: want{
				err:     true,
				code:    app.CodeInternalError,
				Message: app.MessageInternalError,
			},
		},
		{
			name: "fail, case GetID error - invalid product ID returned",
			prepare: func(m mockArgs, args args) {
				m.productStorage.
					EXPECT().
					CreateProduct(args.ctx, mock.Anything).
					Return(access.Product{
						ProductID: "not-a-valid-uuid",
					}, nil)
			},
			args: args{
				req: product.CreateProductRequest{
					Name:           "Test Product",
					Description:    "A test product",
					Price:          99.99,
					OrganizationID: uuid.New().String(),
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

			req := httptest.NewRequest(http.MethodPost, "http://0.0.0.0/api/v1/product", &payload)
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
			h.CreateProduct(ctx)

			// Assert
			var resp wrapper.ResponseOption[product.CreateProductResponse]
			json.NewDecoder(w.Body).Decode(&resp)

			if tt.want.err {
				r.NotEqual(http.StatusOK, w.Code)
				r.Equal(tt.want.code, resp.Code)
				r.Equal(tt.want.Message, resp.Message)
			} else {
				r.Equal(http.StatusCreated, w.Code)
				r.Equal(tt.want.code, resp.Code)
				r.Equal(tt.want.Message, resp.Message)
				r.NotNil(resp.Data)
				r.Equal(productID, resp.Data.ProductID)
			}
		})
	}
}
