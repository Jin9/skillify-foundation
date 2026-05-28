package access

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gitlab.com/b2c-e-commerce-platform/platform/backend/common/app"

	gcpfirestore "cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

var ErrProductNotFound = errors.New("product not found")

type ProductStorage interface {
	GetProductByID(ctx context.Context, id uuid.UUID) (Product, error)
	ListProducts(ctx context.Context, organizationID uuid.UUID) ([]Product, error)
	CreateProduct(ctx context.Context, product Product) (Product, error)
	UpdateProduct(ctx context.Context, product Product) (Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
}

type productStorage struct {
	fs *gcpfirestore.Client
}

var _ ProductStorage = (*productStorage)(nil)

const productCollection = "product"

func NewProductStorage(fs *gcpfirestore.Client) ProductStorage {
	return &productStorage{fs: fs}
}

func (s *productStorage) GetProductByID(ctx context.Context, id uuid.UUID) (Product, error) {
	doc, err := s.fs.Collection(productCollection).
		Doc(id.String()).
		Get(ctx)
	if err != nil {
		return Product{}, fmt.Errorf("failed to query product: %w", err)
	}

	var product Product
	if err := doc.DataTo(&product); err != nil {
		return Product{}, fmt.Errorf("failed to parse product data: %w", err)
	}
	return product, nil
}

func (s *productStorage) ListProducts(ctx context.Context, organizationID uuid.UUID) ([]Product, error) {
	iter := s.fs.Collection(productCollection).
		Where("organization_id", "==", organizationID.String()).
		Documents(ctx)
	defer iter.Stop()

	var products []Product
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate products: %w", err)
		}
		var p Product
		if err := doc.DataTo(&p); err != nil {
			return nil, fmt.Errorf("failed to parse product data: %w", err)
		}
		products = append(products, p)
	}
	return products, nil
}

func (s *productStorage) CreateProduct(ctx context.Context, product Product) (Product, error) {
	if product.ProductID == "" {
		product.ProductID = uuid.New().String()
	}

	app.SetTimestamps(&product.CreatedAt, &product.UpdatedAt)

	_, err := s.fs.Collection(productCollection).Doc(product.ProductID).Set(ctx, product)
	if err != nil {
		return Product{}, fmt.Errorf("failed to create product: %w", err)
	}
	return product, nil
}

func (s *productStorage) UpdateProduct(ctx context.Context, product Product) (Product, error) {
	product.UpdatedAt = time.Now().Local()

	_, err := s.fs.Collection(productCollection).Doc(product.ProductID).Set(ctx, product)
	if err != nil {
		return Product{}, fmt.Errorf("failed to update product: %w", err)
	}
	return product, nil
}

func (s *productStorage) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	_, err := s.fs.Collection(productCollection).Doc(id.String()).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}

type ProductStatusType string

var (
	ProductStatusActive   ProductStatusType = "ACTIVE"
	ProductStatusInactive ProductStatusType = "INACTIVE"
)

type Product struct {
	ProductID      string            `firestore:"product_id" json:"productId"`
	Name           string            `firestore:"name" json:"name"`
	Description    string            `firestore:"description" json:"description"`
	Price          float64           `firestore:"price" json:"price"`
	OrganizationID string            `firestore:"organization_id" json:"organizationId"`
	Status         ProductStatusType `firestore:"status" json:"status"`
	CreatedAt      time.Time         `firestore:"created_at" json:"createdAt"`
	UpdatedAt      time.Time         `firestore:"updated_at" json:"updatedAt"`
}

func (p Product) GetID() (uuid.UUID, error) {
	return uuid.Parse(p.ProductID)
}
