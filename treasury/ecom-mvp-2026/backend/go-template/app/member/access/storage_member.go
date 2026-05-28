package access

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"

	gcpfirestore "cloud.google.com/go/firestore"
	"github.com/google/uuid"
)

// ErrMemberNotFound is returned when no member matches the query.
var ErrMemberNotFound = errors.New("member not found")

type MemberStorage interface {
	GetMemberByID(ctx context.Context, id uuid.UUID) (Member, error)
	CreateMember(ctx context.Context, member Member) (Member, error)
}

type memberStorage struct {
	fs *gcpfirestore.Client
}

var _ MemberStorage = (*memberStorage)(nil)

const memberCollection = "member"

func NewMemberStorage(fs *gcpfirestore.Client) MemberStorage {
	return &memberStorage{
		fs: fs,
	}
}

func (s *memberStorage) GetMemberByID(ctx context.Context, id uuid.UUID) (Member, error) {

	doc, err := s.fs.Collection(memberCollection).
		Doc(id.String()).
		Get(ctx)
	if err != nil {
		return Member{}, fmt.Errorf("failed to query member: %w", err)
	}

	var member Member
	if err := doc.DataTo(&member); err != nil {
		return Member{}, fmt.Errorf("failed to parse member data: %w", err)
	}

	return member, nil
}

func (s *memberStorage) CreateMember(ctx context.Context, member Member) (Member, error) {
	if member.MemberID == "" {
		member.MemberID = uuid.New().String()
	}

	app.SetTimestamps(&member.CreatedAt, &member.UpdatedAt)

	_, err := s.fs.Collection(memberCollection).Doc(member.MemberID).Set(ctx, member)
	if err != nil {
		return Member{}, fmt.Errorf("failed to create member: %w", err)
	}
	return member, nil
}

type MemberStatusType string

var (
	MemberStatusActive    MemberStatusType = "ACTIVE"
	MemberStatusSuspended MemberStatusType = "SUSPENDED"
)

type Member struct {
	MemberID       string           `firestore:"member_id" json:"memberId"`
	Username       string           `firestore:"username" json:"username"`
	EncryptedEmail string           `firestore:"encrypted_email" json:"encryptedEmail"`
	HashedEmail    string           `firestore:"hashed_email" json:"hashedEmail"`
	Status         MemberStatusType `firestore:"status" json:"status"`
	CreatedAt      time.Time        `firestore:"created_at" json:"createdAt"`
	UpdatedAt      time.Time        `firestore:"updated_at" json:"updatedAt"`
	DeletedAt      *time.Time       `firestore:"deleted_at,omitempty" json:"deletedAt,omitempty"`
}

// GetID parses the MemberID string as a UUID.
// Returns uuid.Nil and an error if the stored ID is malformed.
func (m *Member) GetID() (uuid.UUID, error) {
	return uuid.Parse(m.MemberID)
}
