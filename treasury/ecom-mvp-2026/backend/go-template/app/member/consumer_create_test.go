package member_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/kafka"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/member"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/member/access"
	access_mocks "gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/member/access/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestOnMemberCreated(t *testing.T) {
	r := require.New(t)

	type mockArgs struct {
		memberStorage *access_mocks.MemberStorageMock
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
				m.memberStorage.
					EXPECT().
					CreateMember(args.ctx, mock.Anything).
					Return(access.Member{}, nil)
			},
			args: args{
				ctx: context.Background(),
				msg: kafka.Message[json.RawMessage]{
					EventName:   "MEMBER_CREATED",
					AggregateID: uuid.New().String(),
					EventID:     uuid.New().String(),
					Timestamp:   time.Now(),
					Payload: func() json.RawMessage {
						b, _ := json.Marshal(member.CreateMemberMessage{
							Username:       "testuser",
							EncryptedEmail: "encrypted",
							HashedEmail:    "hashed",
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
					EventName: "MEMBER_CREATED",
					Payload:   json.RawMessage(`{invalid`),
				},
			},
			wantErr: true,
		},
		{
			name:    "fail, case validation error - missing required fields",
			prepare: func(m mockArgs, args args) {},
			args: args{
				ctx: context.Background(),
				msg: kafka.Message[json.RawMessage]{
					EventName: "MEMBER_CREATED",
					Payload:   json.RawMessage(`{}`),
				},
			},
			wantErr: true,
		},
		{
			name: "fail, case storage error",
			prepare: func(m mockArgs, args args) {
				m.memberStorage.
					EXPECT().
					CreateMember(args.ctx, mock.Anything).
					Return(access.Member{}, assert.AnError)
			},
			args: args{
				ctx: context.Background(),
				msg: kafka.Message[json.RawMessage]{
					EventName: "MEMBER_CREATED",
					Payload: func() json.RawMessage {
						b, _ := json.Marshal(member.CreateMemberMessage{
							Username:       "testuser",
							EncryptedEmail: "encrypted",
							HashedEmail:    "hashed",
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
				memberStorage: access_mocks.NewMemberStorageMock(t),
			}

			if tt.prepare != nil {
				tt.prepare(m, tt.args)
			}

			h := member.NewHandler(member.HandlerConfig{
				MemberStorage: m.memberStorage,
			})

			err := h.OnMemberCreated(tt.args.ctx, tt.args.msg)
			if tt.wantErr {
				r.Error(err)
			} else {
				r.NoError(err)
			}
		})
	}
}
