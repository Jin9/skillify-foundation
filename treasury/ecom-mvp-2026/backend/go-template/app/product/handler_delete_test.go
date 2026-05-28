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
	access_mocks "gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/product/access/mocks"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteProduct(t *testing.T) {
	r := require.New(t)

	productID := uuid.New()

	type mockArgs struct {
		productStorage *access_mocks.ProductStorageMock
	}

	type args struct {
		ctx context.Context
		req product.DeleteProductRequest
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
					DeleteProduct(args.ctx, args.req.ProductID).
					Return(nil)
			},
			args: args{
				req: product.DeleteProductRequest{
					ProductID: productID,
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
				req: product.DeleteProductRequest{},
			},
			want: want{
				err:     true,
				code:    app.CodeBadRequest,
				Message: app.MessageBadRequest,
			},
		},
		{
			name: "fail, case DeleteProduct storage error",
			prepare: func(m mockArgs, args args) {
				m.productStorage.
					EXPECT().
					DeleteProduct(args.ctx, args.req.ProductID).
					Return(assert.AnError)
			},
			args: args{
				req: product.DeleteProductRequest{
					ProductID: productID,
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

			req := httptest.NewRequest(http.MethodDelete, "http://0.0.0.0/api/v1/product", &payload)
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
			h.DeleteProduct(ctx)

			// Assert
			var resp wrapper.ResponseOption[product.DeleteProductResponse]
			json.NewDecoder(w.Body).Decode(&resp)

			if tt.want.err {
				r.NotEqual(http.StatusOK, w.Code)
				r.Equal(tt.want.code, resp.Code)
				r.Equal(tt.want.Message, resp.Message)
			} else {
				r.Equal(http.StatusOK, w.Code)
				r.Equal(tt.want.code, resp.Code)
				r.Equal(tt.want.Message, resp.Message)
			}
		})
	}
}
