package services

import (
	"context"
	"fmt"

	"github.com/printify-go/pkg/printify/domain"
	"github.com/printify-go/pkg/printify/ports/driven"
	"github.com/printify-go/pkg/printify/ports/driving"
)

// ProductService implements driving.ProductService.
type ProductService struct {
	products driven.ProductGateway
}

// NewProductService builds a ProductService backed by the given gateway.
func NewProductService(products driven.ProductGateway) *ProductService {
	return &ProductService{products: products}
}

func validShop(id domain.ShopID) error {
	if !id.Valid() {
		return fmt.Errorf("shop id %d: %w", id, domain.ErrInvalidInput)
	}
	return nil
}

func validProduct(id domain.ProductID) error {
	if !id.Valid() {
		return fmt.Errorf("product id %q: %w", id, domain.ErrInvalidInput)
	}
	return nil
}

// List returns a page of products for a shop.
func (s *ProductService) List(ctx context.Context, shop domain.ShopID, f driving.ProductFilter) ([]domain.Product, error) {
	if err := validShop(shop); err != nil {
		return nil, err
	}
	products, err := s.products.List(ctx, shop, f)
	if err != nil {
		return nil, fmt.Errorf("listing products in shop %d: %w", shop, err)
	}
	return products, nil
}

// Get returns one product.
func (s *ProductService) Get(ctx context.Context, shop domain.ShopID, id domain.ProductID) (*domain.Product, error) {
	if err := validShop(shop); err != nil {
		return nil, err
	}
	if err := validProduct(id); err != nil {
		return nil, err
	}
	p, err := s.products.Get(ctx, shop, id)
	if err != nil {
		return nil, fmt.Errorf("getting product %s in shop %d: %w", id, shop, err)
	}
	return p, nil
}

// Create creates a new product in a shop.
func (s *ProductService) Create(ctx context.Context, shop domain.ShopID, p domain.CreateProduct) (*domain.Product, error) {
	if err := validShop(shop); err != nil {
		return nil, err
	}
	if err := validateCreateProduct(p); err != nil {
		return nil, err
	}
	created, err := s.products.Create(ctx, shop, p)
	if err != nil {
		return nil, fmt.Errorf("creating product in shop %d: %w", shop, err)
	}
	return created, nil
}

// Update updates an existing product.
func (s *ProductService) Update(ctx context.Context, shop domain.ShopID, id domain.ProductID, p domain.UpdateProduct) (*domain.Product, error) {
	if err := validShop(shop); err != nil {
		return nil, err
	}
	if err := validProduct(id); err != nil {
		return nil, err
	}
	updated, err := s.products.Update(ctx, shop, id, p)
	if err != nil {
		return nil, fmt.Errorf("updating product %s in shop %d: %w", id, shop, err)
	}
	return updated, nil
}

// Delete removes a product from a shop.
func (s *ProductService) Delete(ctx context.Context, shop domain.ShopID, id domain.ProductID) error {
	if err := validShop(shop); err != nil {
		return err
	}
	if err := validProduct(id); err != nil {
		return err
	}
	if err := s.products.Delete(ctx, shop, id); err != nil {
		return fmt.Errorf("deleting product %s in shop %d: %w", id, shop, err)
	}
	return nil
}

// Publish pushes a product to its sales channel.
func (s *ProductService) Publish(ctx context.Context, shop domain.ShopID, id domain.ProductID, opts domain.PublishProduct) error {
	if err := validShop(shop); err != nil {
		return err
	}
	if err := validProduct(id); err != nil {
		return err
	}
	if err := s.products.Publish(ctx, shop, id, opts); err != nil {
		return fmt.Errorf("publishing product %s in shop %d: %w", id, shop, err)
	}
	return nil
}

// Unpublish notifies that a product has been unpublished.
func (s *ProductService) Unpublish(ctx context.Context, shop domain.ShopID, id domain.ProductID) error {
	if err := validShop(shop); err != nil {
		return err
	}
	if err := validProduct(id); err != nil {
		return err
	}
	if err := s.products.Unpublish(ctx, shop, id); err != nil {
		return fmt.Errorf("unpublishing product %s in shop %d: %w", id, shop, err)
	}
	return nil
}

// SetPublishSucceeded records a successful channel publish.
func (s *ProductService) SetPublishSucceeded(ctx context.Context, shop domain.ShopID, id domain.ProductID, ext domain.PublishSucceeded) error {
	if err := validShop(shop); err != nil {
		return err
	}
	if err := validProduct(id); err != nil {
		return err
	}
	if err := s.products.SetPublishSucceeded(ctx, shop, id, ext); err != nil {
		return fmt.Errorf("setting publish succeeded for product %s in shop %d: %w", id, shop, err)
	}
	return nil
}

// SetPublishFailed records a failed channel publish.
func (s *ProductService) SetPublishFailed(ctx context.Context, shop domain.ShopID, id domain.ProductID, reason string) error {
	if err := validShop(shop); err != nil {
		return err
	}
	if err := validProduct(id); err != nil {
		return err
	}
	if reason == "" {
		return fmt.Errorf("publish failure reason: %w", domain.ErrInvalidInput)
	}
	if err := s.products.SetPublishFailed(ctx, shop, id, domain.PublishFailed{Reason: reason}); err != nil {
		return fmt.Errorf("setting publish failed for product %s in shop %d: %w", id, shop, err)
	}
	return nil
}

// ListGpsr returns GPSR compliance information for a product.
func (s *ProductService) ListGpsr(ctx context.Context, shop domain.ShopID, id domain.ProductID) ([]domain.GpsrInfo, error) {
	if err := validShop(shop); err != nil {
		return nil, err
	}
	if err := validProduct(id); err != nil {
		return nil, err
	}
	info, err := s.products.ListGpsr(ctx, shop, id)
	if err != nil {
		return nil, fmt.Errorf("listing gpsr for product %s in shop %d: %w", id, shop, err)
	}
	return info, nil
}

func validateCreateProduct(p domain.CreateProduct) error {
	if p.Title == "" {
		return fmt.Errorf("title: %w", domain.ErrInvalidInput)
	}
	if !p.BlueprintID.Valid() {
		return fmt.Errorf("blueprint id %d: %w", p.BlueprintID, domain.ErrInvalidInput)
	}
	if !p.PrintProviderID.Valid() {
		return fmt.Errorf("print provider id %d: %w", p.PrintProviderID, domain.ErrInvalidInput)
	}
	if len(p.Variants) == 0 {
		return fmt.Errorf("variants: at least one required, %w", domain.ErrInvalidInput)
	}
	if len(p.PrintAreas) == 0 {
		return fmt.Errorf("print areas: at least one required, %w", domain.ErrInvalidInput)
	}
	return nil
}
