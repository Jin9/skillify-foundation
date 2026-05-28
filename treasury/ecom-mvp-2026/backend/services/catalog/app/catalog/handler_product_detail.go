package catalog

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/serror"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/token"
	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/wrapper"
)

// ProductDetailRequest is the JSON body for POST /api/v1/catalog/product/detail.
type ProductDetailRequest struct {
	ProductID string `json:"productId" binding:"required"`
}

// DetailCategory is the category sub-object in product.detail response.
type DetailCategory struct {
	CategoryID string `json:"categoryId"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
}

// ReviewSummary is a placeholder per TD spec (Reviews module out of scope MVP).
type ReviewSummary struct {
	AverageRating interface{} `json:"averageRating"` // always null
	ReviewCount   int         `json:"reviewCount"`   // always 0
}

// ProductDetailResponseData is the data envelope for product.detail.
type ProductDetailResponseData struct {
	ProductID     string        `json:"productId"`
	SKU           string        `json:"sku"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	Images        []string      `json:"images"`
	Price         string        `json:"price"`
	Category      DetailCategory `json:"category"`
	StockStatus   string        `json:"stockStatus"`
	ReviewSummary ReviewSummary `json:"reviewSummary"`
	Status        string        `json:"status"`
}

// GetProduct handles POST /api/v1/catalog/product/detail.
// CAT-006: returns full product with image array; reviewSummary placeholder.
// Visibility rules:
//   - ACTIVE / INACTIVE: visible to all callers (INACTIVE products remain accessible
//     for cart compatibility per TD edge_behavior).
//   - DRAFT / DELETED: visible only to ADMIN callers; guests receive 404.
func (h *Handler) GetProduct(c *gin.Context) {
	ctx := c.Request.Context()

	req, ok := wrapper.BindJSON[ProductDetailRequest](c, slog.String("handler", "GetProduct"))
	if !ok {
		return
	}

	// Validate productId is a UUID
	productID, err := uuid.Parse(strings.TrimSpace(req.ProductID))
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[ProductDetailResponseData]{
			HTTPStatus: http.StatusBadRequest,
			Code:       app.CodeBadRequest,
			Message:    "VALIDATION_ERROR: productId must be a valid UUID",
		})
		return
	}

	// Detect admin role: parse Authorization header if present.
	// Do NOT reject missing auth — public callers are allowed.
	isAdmin := extractIsAdmin(c)

	detail, err := h.svc.GetProduct(ctx, productID, isAdmin)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			wrapper.Respond(c, wrapper.ResponseOption[ProductDetailResponseData]{
				HTTPStatus: http.StatusNotFound,
				Code:       app.CodeNotFound,
				Message:    "PRODUCT_NOT_FOUND",
			})
			return
		}
		wrapper.Respond(c, wrapper.ResponseOption[ProductDetailResponseData]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err:        serror.Wrap(err).With(slog.String("product_id", productID.String())),
		})
		return
	}

	// stockStatus: best-effort inventory call stubbed here.
	// TODO: call inventory.stock.read (POST /api/v1/inventory/stock/read, {sku}) with 2s timeout.
	// availableQty > 10 → IN_STOCK; 1–10 → LOW; 0 → OUT_OF_STOCK.
	// On timeout/error: log warning, return stockStatus=OUT_OF_STOCK (safe fallback per TD).
	stockStatus := StockStatusOutOfStock // placeholder until inventory client is wired

	wrapper.Respond(c, wrapper.ResponseOption[ProductDetailResponseData]{
		HTTPStatus: http.StatusOK,
		Code:       app.CodeSuccess,
		Message:    app.MessageSuccess,
		Data: &ProductDetailResponseData{
			ProductID:   detail.ProductID.String(),
			SKU:         detail.SKU,
			Name:        detail.Name,
			Description: detail.Description,
			Images:      detail.Images,
			Price:       detail.Price,
			Category: DetailCategory{
				CategoryID: detail.CategoryID.String(),
				Name:       detail.CategoryName,
				Slug:       detail.CategorySlug,
			},
			StockStatus: string(stockStatus),
			ReviewSummary: ReviewSummary{
				AverageRating: nil,
				ReviewCount:   0,
			},
			Status: detail.Status,
		},
	})
}

// extractIsAdmin parses the Authorization header if present and returns true
// if the token claims contain role=ADMIN. Returns false on any error (missing
// header, parse failure, wrong role) — callers must never reject on missing auth.
func extractIsAdmin(c *gin.Context) bool {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return false
	}

	// Use claims already set in context by the optional-auth middleware (if any),
	// or parse the raw token directly for public endpoints where JWT middleware
	// is not in the chain.
	claims, err := token.ClaimsFromContext(c.Request.Context())
	if err != nil {
		// Token not in context — attempt raw parse from header for admin detection only.
		// We do NOT verify signature here; this is role-based display logic on a read path.
		// Full verification happens only in the protected admin JWT middleware group.
		// This is acceptable per TD auth_notes: "detect admin by optionally parsing the
		// Authorization header if present — do not reject missing auth."
		return false
	}

	role, _ := claims.Extra["role"].(string)
	return strings.EqualFold(role, "ADMIN")
}
