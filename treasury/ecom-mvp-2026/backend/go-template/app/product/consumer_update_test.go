package product_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/kafka"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/product"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/product/access"
	access_mocks "gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/product/access/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestOnProductUpdated(t *testing.T) {
	r := require.New(t)

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
					UpdateProduct(args.ctx, mock.Anything).
					Return(access.Product{}, nil)
			},
			args: args{
				ctx: context.Background(),
				msg: kafka.Message[json.RawMessage]{
					EventName:   "PRODUCT_UPDATED",
					AggregateID: uuid.New().String(),
					EventID:     uuid.New().String(),
					Timestamp:   time.Now(),
					Payload: func() json.RawMessage {
						b, _ := json.Marshal(product.UpdateProductMessage{
							ProductID:   "prod-123",
							Name:        "Updated Product",
							Description: "updated desc",
							Price:       20.00,
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
					EventName: "PRODUCT_UPDATED",
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
					EventName: "PRODUCT_UPDATED",
					Payload:   json.RawMessage(`{"name":"Updated"}`),
				},
			},
			wantErr: true,
		},
		{
			name: "fail, case storage error",
			prepare: func(m mockArgs, args args) {
				m.productStorage.
					EXPECT().
					UpdateProduct(args.ctx, mock.Anything).
					Return(access.Product{}, assert.AnError)
			},
			args: args{
				ctx: context.Background(),
				msg: kafka.Message[json.RawMessage]{
					EventName: "PRODUCT_UPDATED",
					Payload: func() json.RawMessage {
						b, _ := json.Marshal(product.UpdateProductMessage{
							ProductID: "prod-123",
							Name:      "Updated Product",
							Price:     20.00,
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

			err := h.OnProductUpdated(tt.args.ctx, tt.args.msg)
			if tt.wantErr {
				r.Error(err)
			} else {
				r.NoError(err)
			}
		})
	}
}
