package organization

import (
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/organization/access"
)

type HandlerConfig struct {
	OrganizationStorage access.OrganizationStorage
}

type handler struct {
	organizationStorage access.OrganizationStorage
}

func NewHandler(cfg HandlerConfig) *handler {
	return &handler{
		organizationStorage: cfg.OrganizationStorage,
	}
}
