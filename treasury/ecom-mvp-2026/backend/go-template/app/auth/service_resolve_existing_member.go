package auth

import (
	"context"
	"fmt"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/auth/access"
)

// resolveExistingMember builds the response for an already-registered member.
func (h *handler) resolveExistingMember(ctx context.Context, memberInfo access.Member, profileImage string) (*ResolveIdentityResponse, error) {
	memberID, err := memberInfo.GetID()
	if err != nil {
		return nil, fmt.Errorf("parse member ID: %w", err)
	}

	organizationMember, err := h.organizationStorage.GetMemberOrganization(ctx, memberID)
	if err != nil {
		return nil, fmt.Errorf("get member organization: %w", err)
	}

	email, err := h.cipher.Decrypt(memberInfo.EncryptedEmail)
	if err != nil {
		return nil, fmt.Errorf("decrypt email: %w", err)
	}

	orgID, err := organizationMember.GetOrganizationID()
	if err != nil {
		return nil, fmt.Errorf("parse organization ID: %w", err)
	}

	return &ResolveIdentityResponse{
		IsMember:       true,
		MemberID:       memberID,
		Username:       memberInfo.Username,
		Email:          email,
		HashedEmail:    memberInfo.HashedEmail,
		Status:         memberInfo.Status,
		OrganizationID: orgID,
		Role:           organizationMember.Role,
		ProfileImage:   profileImage,
	}, nil
}
