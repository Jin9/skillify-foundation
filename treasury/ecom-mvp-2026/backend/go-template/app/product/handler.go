package product

import (
	"gitlab.com/b2c-e-commerce-platform/platform/backend/go-template/app/product/access"
)

type HandlerConfig struct {
	ProductStorage access.ProductStorage
}

type handler struct {
	productStorage access.ProductStorage
}

func NewHandler(cfg HandlerConfig) *handler {
	return &handler{
		productStorage: cfg.ProductStorage,
	}
}
