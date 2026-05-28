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

type GetCurrentRequest struct {
	MemberID uuid.UUID `json:"memberId" binding:"required"`
}

type GetCurrentResponse struct {
	OrganizationID uuid.UUID                         `json:"organizationId"`
	Name           string                            `json:"name"`
	Status         access.OrganizationStatusType     `json:"status"`
	Role           access.OrganizationMemberRoleType `json:"role"`
	JoinedAt       *time.Time                        `json:"joinedAt,omitempty"`
	CreatedAt      time.Time                         `json:"createdAt"`
	UpdatedAt      time.Time                         `json:"updatedAt"`
}

func (h *handler) GetCurrent(c *gin.Context) {
	ctx := c.Request.Context()

	req, ok := wrapper.BindJSON[GetCurrentRequest](c, slog.String("handler", "GetCurrentOrganization"))
	if !ok {
		return
	}

	orgMember, err := h.organizationStorage.GetMemberOrganization(ctx, req.MemberID)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[GetCurrentResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err:        serror.Wrap(err).With(slog.String("member_id", req.MemberID.String())),
		})
		return
	}

	orgMemberOrgID, err := orgMember.GetOrganizationID()
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[GetCurrentResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err:        serror.Wrap(err).With(slog.String("member_id", req.MemberID.String())),
		})
		return
	}

	org, err := h.organizationStorage.GetOrganizationByID(ctx, orgMemberOrgID)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[GetCurrentResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
		})
		return
	}

	orgID, err := org.GetOrganizationID()
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[GetCurrentResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err:        serror.Wrap(err).With(slog.String("member_id", req.MemberID.String())),
		})
		return
	}

	wrapper.Respond(c, wrapper.ResponseOption[GetCurrentResponse]{
		HTTPStatus: http.StatusOK,
		Code:       app.CodeSuccess,
		Message:    app.MessageSuccess,
		Data: &GetCurrentResponse{
			OrganizationID: orgID,
			Name:           org.Name,
			Status:         org.Status,
			Role:           orgMember.Role,
			JoinedAt:       orgMember.JoinedAt,
			CreatedAt:      org.CreatedAt,
			UpdatedAt:      org.UpdatedAt,
		},
	})
}
