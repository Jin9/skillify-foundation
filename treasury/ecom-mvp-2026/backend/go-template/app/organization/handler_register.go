package organization

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/organization/access"
)

type RegisterOrganizationRequest struct {
	Email    string    `json:"email" binding:"required"`
	Name     string    `json:"name" binding:"required"`
	MemberID uuid.UUID `json:"memberId" binding:"required"`
}

type RegisterOrganizationResponse struct {
	OrganizationID uuid.UUID                         `json:"organizationId"`
	MemberID       uuid.UUID                         `json:"memberId"`
	Role           access.OrganizationMemberRoleType `json:"role"`
}

func (h *handler) RegisterOrganization(c *gin.Context) {
	ctx := c.Request.Context()

	req, ok := wrapper.BindJSON[RegisterOrganizationRequest](c, slog.String("handler", "RegisterOrganization"))
	if !ok {
		return
	}

	// Create organization using provided email as name
	org := access.Organization{
		Name:   req.Email, // Use email as organization name
		Status: access.OrganizationStatusActive,
	}

	createdOrg, err := h.organizationStorage.CreateOrganization(ctx, org)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[RegisterOrganizationResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err:        serror.Wrap(err).With(slog.String("member_id", req.MemberID.String())),
		})
		return
	}

	// Create organization-member relationship
	now := time.Now()
	orgMember := access.OrganizationMember{
		OrganizationID: createdOrg.OrganizationID,
		MemberID:       req.MemberID.String(),
		Role:           access.OrganizationMemberRoleAdmin,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := h.organizationStorage.CreateOrganizationMember(ctx, orgMember); err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[RegisterOrganizationResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err: serror.Wrap(err).With(
				slog.String("member_id", req.MemberID.String()),
				slog.String("organization_id", createdOrg.OrganizationID),
			),
		})
		return
	}

	// Return organization ID, member ID, and role
	orgID, err := createdOrg.GetOrganizationID()
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[RegisterOrganizationResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err:        serror.Wrap(err).With(slog.String("organization_id", createdOrg.OrganizationID)),
		})
		return
	}

	wrapper.Respond(c, wrapper.ResponseOption[RegisterOrganizationResponse]{
		HTTPStatus: http.StatusOK,
		Code:       app.CodeSuccess,
		Message:    app.MessageSuccess,
		Data: &RegisterOrganizationResponse{
			OrganizationID: orgID,
			MemberID:       req.MemberID,
			Role:           access.OrganizationMemberRoleAdmin,
		},
	})
}
