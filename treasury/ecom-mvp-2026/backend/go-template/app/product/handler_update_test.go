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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateProduct(t *testing.T) {
	r := require.New(t)

	productID := uuid.New()
	now := time.Now()

	type mockArgs struct {
		productStorage *access_mocks.ProductStorageMock
	}

	type args struct {
		ctx context.Context
		req product.UpdateProductRequest
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
			name: "success, case valid request with partial update",
			prepare: func(m mockArgs, args args) {
				m.productStorage.
					EXPECT().
					GetProductByID(args.ctx, args.req.ProductID).
					Return(access.Product{
						ProductID:      productID.String(),
						Name:           "Old Name",
						Description:    "Old Description",
						Price:          50.00,
						OrganizationID: uuid.New().String(),
						Status:         access.ProductStatusActive,
						CreatedAt:      now,
						UpdatedAt:      now,
					}, nil)

				m.productStorage.
					EXPECT().
					UpdateProduct(args.ctx, mock.MatchedBy(func(p access.Product) bool {
						return p.Name == "New Name" &&
							p.Description == "Old Description" &&
							p.Price == 50.00
					})).
					Return(access.Product{
						ProductID: productID.String(),
					}, nil)
			},
			args: args{
				req: product.UpdateProductRequest{
					ProductID: productID,
					Name:      "New Name",
				},
			},
			want: want{
				err:     false,
				code:    app.CodeSuccess,
				Message: app.MessageSuccess,
			},
		},
		{
			name: "success, case valid request with all fields updated",
			prepare: func(m mockArgs, args args) {
				m.productStorage.
					EXPECT().
					GetProductByID(args.ctx, args.req.ProductID).
					Return(access.Product{
						ProductID:      productID.String(),
						Name:           "Old Name",
						Description:    "Old Description",
						Price:          50.00,
						OrganizationID: uuid.New().String(),
						Status:         access.ProductStatusActive,
						CreatedAt:      now,
						UpdatedAt:      now,
					}, nil)

				m.productStorage.
					EXPECT().
					UpdateProduct(args.ctx, mock.MatchedBy(func(p access.Product) bool {
						return p.Name == "New Name" &&
							p.Description == "New Description" &&
							p.Price == 199.99
					})).
					Return(access.Product{
						ProductID: productID.String(),
					}, nil)
			},
			args: args{
				req: product.UpdateProductRequest{
					ProductID:   productID,
					Name:        "New Name",
					Description: "New Description",
					Price:       199.99,
				},
			},
			want: want{
				err:     false,
				code:    app.CodeSuccess,
				Message: app.MessageSuccess,
			},
		},
		{
			name: "fail, case invalid request body - missing product ID",
			prepare: func(m mockArgs, args args) {
				// no calls
			},
			args: args{
				req: product.UpdateProductRequest{},
			},
			want: want{
				err:     true,
				code:    app.CodeBadRequest,
				Message: app.MessageBadRequest,
			},
		},
		{
			name: "fail, case GetProductByID error",
			prepare: func(m mockArgs, args args) {
				m.productStorage.
					EXPECT().
					GetProductByID(args.ctx, args.req.ProductID).
					Return(access.Product{}, assert.AnError)
			},
			args: args{
				req: product.UpdateProductRequest{
					ProductID: productID,
					Name:      "New Name",
				},
			},
			want: want{
				err:     true,
				code:    app.CodeInternalError,
				Message: app.MessageInternalError,
			},
		},
		{
			name: "fail, case UpdateProduct storage error",
			prepare: func(m mockArgs, args args) {
				m.productStorage.
					EXPECT().
					GetProductByID(args.ctx, args.req.ProductID).
					Return(access.Product{
						ProductID: productID.String(),
						Name:      "Old Name",
					}, nil)

				m.productStorage.
					EXPECT().
					UpdateProduct(args.ctx, mock.Anything).
					Return(access.Product{}, assert.AnError)
			},
			args: args{
				req: product.UpdateProductRequest{
					ProductID: productID,
					Name:      "New Name",
				},
			},
			want: want{
				err:     true,
				code:    app.CodeInternalError,
				Message: app.MessageInternalError,
			},
		},
		{
			name: "fail, case GetID error - invalid product ID returned from update",
			prepare: func(m mockArgs, args args) {
				m.productStorage.
					EXPECT().
					GetProductByID(args.ctx, args.req.ProductID).
					Return(access.Product{
						ProductID: productID.String(),
						Name:      "Old Name",
					}, nil)

				m.productStorage.
					EXPECT().
					UpdateProduct(args.ctx, mock.Anything).
					Return(access.Product{
						ProductID: "not-a-valid-uuid",
					}, nil)
			},
			args: args{
				req: product.UpdateProductRequest{
					ProductID: productID,
					Name:      "New Name",
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

			req := httptest.NewRequest(http.MethodPut, "http://0.0.0.0/api/v1/product", &payload)
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
			h.UpdateProduct(ctx)

			// Assert
			var resp wrapper.ResponseOption[product.UpdateProductResponse]
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
				r.Equal(productID, resp.Data.ProductID)
			}
		})
	}
}
