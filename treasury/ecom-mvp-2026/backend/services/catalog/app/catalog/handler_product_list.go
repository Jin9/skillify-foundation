package catalog

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

// ProductListRequest is the JSON body for POST /api/v1/catalog/product/list.
// CAT-001..005
type ProductListRequest struct {
	Page   int            `json:"page"`
	Limit  int            `json:"limit"`
	Sort   string         `json:"sort"`
	Filter *ProductFilter `json:"filter"`
}

// ProductFilter carries optional filter parameters.
type ProductFilter struct {
	CategoryID *string  `json:"categoryId"`
	MinPrice   *float64 `json:"minPrice"`
	MaxPrice   *float64 `json:"maxPrice"`
	InStock    *bool    `json:"inStock"`
	Search     *string  `json:"search"`
}

// ProductListResponseItem is a single product in the list response.
type ProductListResponseItem struct {
	ProductID    string           `json:"productId"`
	SKU          string           `json:"sku"`
	Name         string           `json:"name"`
	Price        string           `json:"price"`
	Category     ListItemCategory `json:"category"`
	ThumbnailURL string           `json:"thumbnailUrl"`
	Status       string           `json:"status"`
}

// ListItemCategory is the category sub-object in product list items.
type ListItemCategory struct {
	CategoryID string `json:"categoryId"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
}

// ProductListResponse is the data envelope for product.list.
type ProductListResponse struct {
	Items []ProductListResponseItem `json:"items"`
	Total int                       `json:"total"`
	Page  int                       `json:"page"`
	Limit int                       `json:"limit"`
}

// ListProducts handles POST /api/v1/catalog/product/list.
// Covers BA ACs: CAT-001 (active+visible filter), CAT-002 (pagination+sort),
// CAT-003 (categoryId filter), CAT-004 (price range filter), CAT-005 (inStock stub).
func (h *Handler) ListProducts(c *gin.Context) {
	ctx := c.Request.Context()

	req, ok := wrapper.BindJSON[ProductListRequest](c, slog.String("handler", "ListProducts"))
	if !ok {
		return
	}

	// --- Validation ---

	// page default 1; minimum 1
	page := req.Page
	if page == 0 {
		page = 1
	}
	if page < 1 {
		wrapper.Respond(c, wrapper.ResponseOption[ProductListResponse]{
			HTTPStatus: http.StatusBadRequest,
			Code:       app.CodeBadRequest,
			Message:    "VALIDATION_ERROR: page must be >= 1",
		})
		return
	}

	// limit default 20; range 1..100
	limit := req.Limit
	if limit == 0 {
		limit = 20
	}
	if limit < 1 || limit > 100 {
		wrapper.Respond(c, wrapper.ResponseOption[ProductListResponse]{
			HTTPStatus: http.StatusBadRequest,
			Code:       app.CodeBadRequest,
			Message:    "VALIDATION_ERROR: limit must be 1..100",
		})
		return
	}

	// sort default "created_desc"; must be in valid set
	sortStr := strings.TrimSpace(req.Sort)
	if sortStr == "" {
		sortStr = string(SortCreatedDesc)
	}
	sort := SortOption(sortStr)
	if _, ok := ValidSortOptions[sort]; !ok {
		wrapper.Respond(c, wrapper.ResponseOption[ProductListResponse]{
			HTTPStatus: http.StatusBadRequest,
			Code:       app.CodeBadRequest,
			Message:    "VALIDATION_ERROR: sort must be one of price_asc|price_desc|name_asc|name_desc|created_desc",
		})
		return
	}

	// Build domain filter
	filter := ProductListFilter{}

	if req.Filter != nil {
		f := req.Filter

		// categoryId: optional UUID
		if f.CategoryID != nil && *f.CategoryID != "" {
			catID, err := uuid.Parse(*f.CategoryID)
			if err != nil {
				wrapper.Respond(c, wrapper.ResponseOption[ProductListResponse]{
					HTTPStatus: http.StatusBadRequest,
					Code:       app.CodeBadRequest,
					Message:    "VALIDATION_ERROR: categoryId must be a valid UUID",
				})
				return
			}
			filter.CategoryID = &catID
		}

		// minPrice / maxPrice validation
		if f.MinPrice != nil && *f.MinPrice < 0 {
			wrapper.Respond(c, wrapper.ResponseOption[ProductListResponse]{
				HTTPStatus: http.StatusBadRequest,
				Code:       app.CodeBadRequest,
				Message:    "VALIDATION_ERROR: minPrice must be >= 0",
			})
			return
		}
		if f.MaxPrice != nil && *f.MaxPrice < 0 {
			wrapper.Respond(c, wrapper.ResponseOption[ProductListResponse]{
				HTTPStatus: http.StatusBadRequest,
				Code:       app.CodeBadRequest,
				Message:    "VALIDATION_ERROR: maxPrice must be >= 0",
			})
			return
		}
		if f.MinPrice != nil && f.MaxPrice != nil && *f.MaxPrice < *f.MinPrice {
			wrapper.Respond(c, wrapper.ResponseOption[ProductListResponse]{
				HTTPStatus: http.StatusBadRequest,
				Code:       app.CodeBadRequest,
				Message:    "VALIDATION_ERROR: maxPrice must be >= minPrice",
			})
			return
		}
		if f.MinPrice != nil {
			s := formatDecimal(*f.MinPrice)
			filter.MinPrice = &s
		}
		if f.MaxPrice != nil {
			s := formatDecimal(*f.MaxPrice)
			filter.MaxPrice = &s
		}

		// search: trim, ignore empty
		if f.Search != nil {
			trimmed := strings.TrimSpace(*f.Search)
			if len(trimmed) > 200 {
				trimmed = trimmed[:200]
			}
			if trimmed != "" {
				filter.Search = &trimmed
			}
		}

		// inStock: pass through; actual inventory call is stubbed (TODO below)
		filter.InStock = f.InStock
	}

	// TODO: wire inventory.stock.bulk-read once Inventory service is up.
	// When filter.InStock != nil && *filter.InStock == true:
	//   Phase 1 already done by service.ListProducts (SQL returns page of ACTIVE+visible rows).
	//   Phase 2: call client_inventory.BulkReadStock(ctx, skus) with 2s timeout.
	//   Filter out items where availableQty == 0 or sku not in response map.
	//   On inventory 503/timeout: return 504 UPSTREAM_TIMEOUT (do NOT silently return unfiltered list).

	result, err := h.svc.ListProducts(ctx, filter, sort, page, limit)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[ProductListResponse]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err:        serror.Wrap(err),
		})
		return
	}

	// Map to response DTOs
	items := make([]ProductListResponseItem, 0, len(result.Items))
	for _, p := range result.Items {
		items = append(items, ProductListResponseItem{
			ProductID: p.ProductID.String(),
			SKU:       p.SKU,
			Name:      p.Name,
			Price:     p.Price,
			Category: ListItemCategory{
				CategoryID: p.CategoryID.String(),
				Name:       p.CategoryName,
				Slug:       p.CategorySlug,
			},
			ThumbnailURL: p.ThumbnailURL,
			Status:       p.Status,
		})
	}

	wrapper.Respond(c, wrapper.ResponseOption[ProductListResponse]{
		HTTPStatus: http.StatusOK,
		Code:       app.CodeSuccess,
		Message:    app.MessageSuccess,
		Data: &ProductListResponse{
			Items: items,
			Total: result.Total,
			Page:  result.Page,
			Limit: result.Limit,
		},
	})
}

// formatDecimal converts a float64 price to a NUMERIC-safe string.
// e.g. 10.5 → "10.50", 10.0 → "10.00"
func formatDecimal(v float64) string {
	return fmt.Sprintf("%.2f", v)
}
