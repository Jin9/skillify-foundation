package product_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/kafka"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/product"
	access_mocks "gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/product/access/mocks"
)

func TestOnProductDeleted(t *testing.T) {
	r := require.New(t)

	productID := uuid.New()

	type mockArgs struct {
		productStorage *access_mocks.ProductStorageMock
	}

	type args struct {
		ctx context.Context
		msg kafka.Message[json.RawMessage]
	}

	tests := []struct {
		name    string
		prepare func(m mockArgs, args args)
		args    args
		wantErr bool
	}{
		{
			name: "success, case valid event",
			prepare: func(m mockArgs, args args) {
				m.productStorage.
					EXPECT().
					DeleteProduct(args.ctx, productID).
					Return(nil)
			},
			args: args{
				ctx: context.Background(),
				msg: kafka.Message[json.RawMessage]{
					EventName:   "PRODUCT_DELETED",
					AggregateID: productID.String(),
					EventID:     uuid.New().String(),
					Timestamp:   time.Now(),
					Payload: func() json.RawMessage {
						b, _ := json.Marshal(product.DeleteProductMessage{
							ProductID: productID.String(),
						})
						return b
					}(),
				},
			},
			wantErr: false,
		},
		{
			name:    "fail, case invalid JSON payload",
			prepare: func(m mockArgs, args args) {},
			args: args{
				ctx: context.Background(),
				msg: kafka.Message[json.RawMessage]{
					EventName: "PRODUCT_DELETED",
					Payload:   json.RawMessage(`{invalid`),
				},
			},
			wantErr: true,
		},
		{
			name:    "fail, case validation error - missing required productId",
			prepare: func(m mockArgs, args args) {},
			args: args{
				ctx: context.Background(),
				msg: kafka.Message[json.RawMessage]{
					EventName: "PRODUCT_DELETED",
					Payload:   json.RawMessage(`{}`),
				},
			},
			wantErr: true,
		},
		{
			name:    "fail, case GetID error - invalid product ID",
			prepare: func(m mockArgs, args args) {},
			args: args{
				ctx: context.Background(),
				msg: kafka.Message[json.RawMessage]{
					EventName: "PRODUCT_DELETED",
					Payload: func() json.RawMessage {
						b, _ := json.Marshal(product.DeleteProductMessage{
							ProductID: "not-a-valid-uuid",
						})
						return b
					}(),
				},
			},
			wantErr: true,
		},
		{
			name: "fail, case DeleteProduct storage error",
			prepare: func(m mockArgs, args args) {
				m.productStorage.
					EXPECT().
					DeleteProduct(args.ctx, productID).
					Return(assert.AnError)
			},
			args: args{
				ctx: context.Background(),
				msg: kafka.Message[json.RawMessage]{
					EventName: "PRODUCT_DELETED",
					Payload: func() json.RawMessage {
						b, _ := json.Marshal(product.DeleteProductMessage{
							ProductID: productID.String(),
						})
						return b
					}(),
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mockArgs{
				productStorage: access_mocks.NewProductStorageMock(t),
			}

			if tt.prepare != nil {
				tt.prepare(m, tt.args)
			}

			h := product.NewHandler(product.HandlerConfig{
				ProductStorage: m.productStorage,
			})

			err := h.OnProductDeleted(tt.args.ctx, tt.args.msg)
			if tt.wantErr {
				r.Error(err)
			} else {
				r.NoError(err)
			}
		})
	}
}
