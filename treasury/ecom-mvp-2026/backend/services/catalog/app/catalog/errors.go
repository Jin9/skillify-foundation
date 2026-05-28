package catalog

import "errors"

// Sentinel errors used by service layer and mapped to HTTP codes in handlers.
var (
	ErrProductNotFound       = errors.New("catalog: product not found")
	ErrDuplicateSKU          = errors.New("catalog: duplicate sku")
	ErrCategoryNotFound      = errors.New("catalog: category not found or inactive")
	ErrValidation            = errors.New("catalog: validation error")
	ErrDuplicateCategoryName = errors.New("catalog: duplicate category name")
	ErrDuplicateSlug         = errors.New("catalog: duplicate slug")
)
