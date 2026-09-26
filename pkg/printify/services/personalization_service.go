package services

import (
	"context"
	"fmt"

	"github.com/printify-go/pkg/printify/domain"
	"github.com/printify-go/pkg/printify/ports/driven"
)

// PersonalizationService implements driving.PersonalizationService.
type PersonalizationService struct {
	personalization driven.PersonalizationGateway
}

// NewPersonalizationService builds a PersonalizationService backed by the given gateway.
func NewPersonalizationService(gw driven.PersonalizationGateway) *PersonalizationService {
	return &PersonalizationService{personalization: gw}
}

func validTask(id domain.TaskID) error {
	if !id.Valid() {
		return fmt.Errorf("task id %q: %w", id, domain.ErrInvalidInput)
	}
	return nil
}

// ListOptions returns the personalization fields of a product.
func (s *PersonalizationService) ListOptions(ctx context.Context, shop domain.ShopID, product domain.ProductID) ([]domain.PersonalizationOption, error) {
	if err := validShop(shop); err != nil {
		return nil, err
	}
	if err := validProduct(product); err != nil {
		return nil, err
	}
	opts, err := s.personalization.ListOptions(ctx, shop, product)
	if err != nil {
		return nil, fmt.Errorf("listing personalization options for product %s in shop %d: %w", product, shop, err)
	}
	return opts, nil
}

// CreateConfig configures personalization for a product variant.
func (s *PersonalizationService) CreateConfig(ctx context.Context, shop domain.ShopID, product domain.ProductID, cfg domain.CreatePersonalizationConfig) (*domain.PersonalizationConfig, error) {
	if err := validShop(shop); err != nil {
		return nil, err
	}
	if err := validProduct(product); err != nil {
		return nil, err
	}
	if err := validatePersonalizationConfig(cfg); err != nil {
		return nil, err
	}
	config, err := s.personalization.CreateConfig(ctx, shop, product, cfg)
	if err != nil {
		return nil, fmt.Errorf("creating personalization config for product %s in shop %d: %w", product, shop, err)
	}
	return config, nil
}

// CreatePreviewTask queues a personalization preview render.
func (s *PersonalizationService) CreatePreviewTask(ctx context.Context, shop domain.ShopID, product domain.ProductID, task domain.CreatePreviewTask) (*domain.PreviewTask, error) {
	if err := validShop(shop); err != nil {
		return nil, err
	}
	if err := validProduct(product); err != nil {
		return nil, err
	}
	if len(task.VariantIDs) == 0 {
		return nil, fmt.Errorf("variant ids: at least one required, %w", domain.ErrInvalidInput)
	}
	created, err := s.personalization.CreatePreviewTask(ctx, shop, product, task)
	if err != nil {
		return nil, fmt.Errorf("creating preview task for product %s in shop %d: %w", product, shop, err)
	}
	return created, nil
}

// GetPreviewTask returns the status of a preview task.
func (s *PersonalizationService) GetPreviewTask(ctx context.Context, shop domain.ShopID, product domain.ProductID, id domain.TaskID) (*domain.PreviewTask, error) {
	if err := validShop(shop); err != nil {
		return nil, err
	}
	if err := validProduct(product); err != nil {
		return nil, err
	}
	if err := validTask(id); err != nil {
		return nil, err
	}
	task, err := s.personalization.GetPreviewTask(ctx, shop, product, id)
	if err != nil {
		return nil, fmt.Errorf("getting preview task %s for product %s in shop %d: %w", id, product, shop, err)
	}
	return task, nil
}

func validatePersonalizationConfig(cfg domain.CreatePersonalizationConfig) error {
	if !cfg.VariantID.Valid() {
		return fmt.Errorf("variant id %d: %w", cfg.VariantID, domain.ErrInvalidInput)
	}
	if len(cfg.Items) == 0 {
		return fmt.Errorf("items: at least one required, %w", domain.ErrInvalidInput)
	}
	for i, item := range cfg.Items {
		if item.FieldID == "" {
			return fmt.Errorf("item %d: field id required, %w", i, domain.ErrInvalidInput)
		}
	}
	return nil
}
