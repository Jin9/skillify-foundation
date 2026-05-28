package catalog

import (
	"errors"
	"fmt"
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

// CreateProductRequest is the JSON body for POST /api/v1/catalog/product/create.
type CreateProductRequest struct {
	SKU         string   `json:"sku"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
	Price       float64  `json:"price"`
	CategoryID  string   `json:"categoryId"`
	Status      string   `json:"status"`
}

// CreateProductResponseData is the data returned on successful product creation.
type CreateProductResponseData struct {
	ProductID string `json:"productId"`
}

// CreateProduct handles POST /api/v1/catalog/product/create.
// CAT-007: validates request, INSERTs product + outbox row in same tx.
// Auth: ADMIN role required (enforced by JWT middleware group in router).
func (h *Handler) CreateProduct(c *gin.Context) {
	ctx := c.Request.Context()

	req, ok := wrapper.BindJSON[CreateProductRequest](c, slog.String("handler", "CreateProduct"))
	if !ok {
		return
	}

	// --- Claim extraction (JWT middleware already ran; actor must be ADMIN) ---
	claims, err := token.ClaimsFromContext(ctx)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[CreateProductResponseData]{
			HTTPStatus: http.StatusUnauthorized,
			Code:       app.CodeUnauthorized,
			Message:    app.MessageUnauthorized,
		})
		return
	}

	role, _ := claims.Extra["role"].(string)
	if !strings.EqualFold(role, "ADMIN") {
		wrapper.Respond(c, wrapper.ResponseOption[CreateProductResponseData]{
			HTTPStatus: http.StatusForbidden,
			Code:       app.CodeForbidden,
			Message:    "AUTH_FORBIDDEN",
		})
		return
	}

	actorUserID, err := uuid.Parse(claims.Sub)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[CreateProductResponseData]{
			HTTPStatus: http.StatusUnauthorized,
			Code:       app.CodeUnauthorized,
			Message:    app.MessageUnauthorized,
		})
		return
	}

	// --- Field-level validation ---

	sku := strings.TrimSpace(req.SKU)
	if sku == "" || len(sku) > 100 {
		wrapper.Respond(c, wrapper.ResponseOption[CreateProductResponseData]{
			HTTPStatus: http.StatusBadRequest,
			Code:       app.CodeBadRequest,
			Message:    "VALIDATION_ERROR: sku is required and must be <= 100 chars",
		})
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" || len(name) > 255 {
		wrapper.Respond(c, wrapper.ResponseOption[CreateProductResponseData]{
			HTTPStatus: http.StatusBadRequest,
			Code:       app.CodeBadRequest,
			Message:    "VALIDATION_ERROR: name is required and must be <= 255 chars",
		})
		return
	}

	if req.Price <= 0 {
		wrapper.Respond(c, wrapper.ResponseOption[CreateProductResponseData]{
			HTTPStatus: http.StatusBadRequest,
			Code:       app.CodeBadRequest,
			Message:    "VALIDATION_ERROR: price must be > 0",
		})
		return
	}
	// Format price as NUMERIC string with max 2 decimal places
	priceStr := fmt.Sprintf("%.2f", req.Price)

	categoryIDStr := strings.TrimSpace(req.CategoryID)
	if categoryIDStr == "" {
		wrapper.Respond(c, wrapper.ResponseOption[CreateProductResponseData]{
			HTTPStatus: http.StatusBadRequest,
			Code:       app.CodeBadRequest,
			Message:    "VALIDATION_ERROR: categoryId is required",
		})
		return
	}
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		wrapper.Respond(c, wrapper.ResponseOption[CreateProductResponseData]{
			HTTPStatus: http.StatusBadRequest,
			Code:       app.CodeBadRequest,
			Message:    "VALIDATION_ERROR: categoryId must be a valid UUID",
		})
		return
	}

	// Status: default DRAFT; only DRAFT|ACTIVE allowed on create
	statusStr := strings.ToUpper(strings.TrimSpace(req.Status))
	if statusStr == "" {
		statusStr = string(ProductStatusDraft)
	}
	status := ProductStatus(statusStr)
	if status != ProductStatusDraft && status != ProductStatusActive {
		wrapper.Respond(c, wrapper.ResponseOption[CreateProductResponseData]{
			HTTPStatus: http.StatusBadRequest,
			Code:       app.CodeBadRequest,
			Message:    "VALIDATION_ERROR: status must be DRAFT or ACTIVE",
		})
		return
	}

	// Images: optional; max 20 entries; each non-empty
	if len(req.Images) > 20 {
		wrapper.Respond(c, wrapper.ResponseOption[CreateProductResponseData]{
			HTTPStatus: http.StatusBadRequest,
			Code:       app.CodeBadRequest,
			Message:    "VALIDATION_ERROR: images must not exceed 20 entries",
		})
		return
	}
	for i, img := range req.Images {
		if strings.TrimSpace(img) == "" {
			wrapper.Respond(c, wrapper.ResponseOption[CreateProductResponseData]{
				HTTPStatus: http.StatusBadRequest,
				Code:       app.CodeBadRequest,
				Message:    wrapper.Message(fmt.Sprintf("VALIDATION_ERROR: images[%d] must be a non-empty URL", i)),
			})
			return
		}
	}

	// --- Service call ---
	input := CreateProductInput{
		SKU:         sku,
		Name:        name,
		Description: req.Description,
		Images:      req.Images,
		Price:       priceStr,
		CategoryID:  categoryID,
		Status:      status,
		ActorUserID: actorUserID,
	}

	productID, err := h.svc.CreateProduct(ctx, input)
	if err != nil {
		if errors.Is(err, ErrDuplicateSKU) {
			wrapper.Respond(c, wrapper.ResponseOption[CreateProductResponseData]{
				HTTPStatus: http.StatusConflict,
				Code:       app.Code("DUPLICATE_SKU"),
				Message:    "DUPLICATE_SKU: sku already exists",
			})
			return
		}
		if errors.Is(err, ErrCategoryNotFound) {
			wrapper.Respond(c, wrapper.ResponseOption[CreateProductResponseData]{
				HTTPStatus: http.StatusNotFound,
				Code:       app.CodeNotFound,
				Message:    "NOT_FOUND: categoryId does not exist or is inactive",
			})
			return
		}
		wrapper.Respond(c, wrapper.ResponseOption[CreateProductResponseData]{
			HTTPStatus: http.StatusInternalServerError,
			Code:       app.CodeInternalError,
			Message:    app.MessageInternalError,
			Err: serror.Wrap(err).With(
				slog.String("sku", sku),
				slog.String("category_id", categoryID.String()),
			),
		})
		return
	}

	wrapper.Respond(c, wrapper.ResponseOption[CreateProductResponseData]{
		HTTPStatus: http.StatusCreated,
		Code:       app.Code("CREATED"),
		Message:    "Product Created",
		Data:       &CreateProductResponseData{ProductID: productID.String()},
	})
}
