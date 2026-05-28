package member

import (
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/crypt"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/hash"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/member/access"
)

type HandlerConfig struct {
	MemberStorage access.MemberStorage
	Cipher        crypt.Cipher
	Hash          hash.HashManager
}

type handler struct {
	memberStorage access.MemberStorage
	cipher        crypt.Cipher
	hash          hash.HashManager
}

func NewHandler(cfg HandlerConfig) *handler {
	return &handler{
		memberStorage: cfg.MemberStorage,
		cipher:        cfg.Cipher,
		hash:          cfg.Hash,
	}
}
