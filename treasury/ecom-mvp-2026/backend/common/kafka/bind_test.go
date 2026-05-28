package kafka_test

import (
	"encoding/json"
	"testing"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/kafka"

	"github.com/stretchr/testify/require"
)

type testPayload struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required,gt=0"`
}

func TestBindMessage(t *testing.T) {
	r := require.New(t)

	tests := []struct {
		name    string
		data    json.RawMessage
		wantErr bool
	}{
		{
			name:    "success, case valid payload",
			data:    json.RawMessage(`{"name":"product","price":10.5}`),
			wantErr: false,
		},
		{
			name:    "fail, case invalid JSON",
			data:    json.RawMessage(`{invalid`),
			wantErr: true,
		},
		{
			name:    "fail, case missing required field",
			data:    json.RawMessage(`{"price":10.5}`),
			wantErr: true,
		},
		{
			name:    "fail, case validation gt=0 violated",
			data:    json.RawMessage(`{"name":"product","price":-1}`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dest testPayload
			err := kafka.BindMessage(tt.data, &dest)
			if tt.wantErr {
				r.Error(err)
			} else {
				r.NoError(err)
				r.Equal("product", dest.Name)
				r.Equal(10.5, dest.Price)
			}
		})
	}
}
