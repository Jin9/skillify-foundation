package organization_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/kafka"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/organization"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/organization/access"
	access_mocks "gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/organization/access/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestOnOrganizationMemberAdded(t *testing.T) {
	r := require.New(t)

	type mockArgs struct {
		organizationStorage *access_mocks.OrganizationStorageMock
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
				m.organizationStorage.
					EXPECT().
					CreateOrganizationMember(args.ctx, mock.Anything).
					Return(nil)
			},
			args: args{
				ctx: context.Background(),
				msg: kafka.Message[json.RawMessage]{
					EventName:   "ORGANIZATION_MEMBER_ADDED",
					AggregateID: uuid.New().String(),
					EventID:     uuid.New().String(),
					Timestamp:   time.Now(),
					Payload: func() json.RawMessage {
						b, _ := json.Marshal(organization.AddOrganizationMemberMessage{
							OrganizationID: "org-123",
							MemberID:       "member-456",
							Role:           access.OrganizationMemberRoleAdmin,
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
					EventName: "ORGANIZATION_MEMBER_ADDED",
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
					EventName: "ORGANIZATION_MEMBER_ADDED",
					Payload:   json.RawMessage(`{}`),
				},
			},
			wantErr: true,
		},
		{
			name: "fail, case storage error",
			prepare: func(m mockArgs, args args) {
				m.organizationStorage.
					EXPECT().
					CreateOrganizationMember(args.ctx, mock.Anything).
					Return(assert.AnError)
			},
			args: args{
				ctx: context.Background(),
				msg: kafka.Message[json.RawMessage]{
					EventName: "ORGANIZATION_MEMBER_ADDED",
					Payload: func() json.RawMessage {
						b, _ := json.Marshal(organization.AddOrganizationMemberMessage{
							OrganizationID: "org-123",
							MemberID:       "member-456",
							Role:           access.OrganizationMemberRoleUser,
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
				organizationStorage: access_mocks.NewOrganizationStorageMock(t),
			}

			if tt.prepare != nil {
				tt.prepare(m, tt.args)
			}

			h := organization.NewHandler(organization.HandlerConfig{
				OrganizationStorage: m.organizationStorage,
			})

			err := h.OnOrganizationMemberAdded(tt.args.ctx, tt.args.msg)
			if tt.wantErr {
				r.Error(err)
			} else {
				r.NoError(err)
			}
		})
	}
}
