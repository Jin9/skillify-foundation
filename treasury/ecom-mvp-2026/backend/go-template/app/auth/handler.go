package auth

import (
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/crypt"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/hash"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/auth/access"
)

type HandlerConfig struct {
	GoogleClient        access.GoogleClient
	MemberStorage       access.MemberStorage
	OrganizationStorage access.OrganizationStorage
	Hash                hash.HashManager
	Cipher              crypt.Cipher
	Token               token.JWTSigner
}

type handler struct {
	googleClient        access.GoogleClient
	memberStorage       access.MemberStorage
	organizationStorage access.OrganizationStorage
	hash                hash.HashManager
	cipher              crypt.Cipher
	token               token.JWTSigner
}

func NewHandler(cfg HandlerConfig) *handler {
	return &handler{
		googleClient:        cfg.GoogleClient,
		memberStorage:       cfg.MemberStorage,
		organizationStorage: cfg.OrganizationStorage,
		hash:                cfg.Hash,
		cipher:              cfg.Cipher,
		token:               cfg.Token,
	}
}
