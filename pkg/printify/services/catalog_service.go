package services

import (
	"context"
	"fmt"

	"github.com/printify-go/pkg/printify/domain"
	"github.com/printify-go/pkg/printify/ports/driven"
)

// CatalogService implements driving.CatalogService.
type CatalogService struct {
	catalog driven.CatalogGateway
}

// NewCatalogService builds a CatalogService backed by the given gateway.
func NewCatalogService(catalog driven.CatalogGateway) *CatalogService {
	return &CatalogService{catalog: catalog}
}

// ListBlueprints returns all blueprints in the catalog.
func (s *CatalogService) ListBlueprints(ctx context.Context) ([]domain.Blueprint, error) {
	bps, err := s.catalog.ListBlueprints(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing blueprints: %w", err)
	}
	return bps, nil
}

// GetBlueprint returns one blueprint.
func (s *CatalogService) GetBlueprint(ctx context.Context, id domain.BlueprintID) (*domain.Blueprint, error) {
	if !id.Valid() {
		return nil, fmt.Errorf("blueprint id %d: %w", id, domain.ErrInvalidInput)
	}
	bp, err := s.catalog.GetBlueprint(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting blueprint %d: %w", id, err)
	}
	return bp, nil
}

// GetSizeGuide returns the size guide for a blueprint.
func (s *CatalogService) GetSizeGuide(ctx context.Context, id domain.BlueprintID) (*domain.SizeGuide, error) {
	if !id.Valid() {
		return nil, fmt.Errorf("blueprint id %d: %w", id, domain.ErrInvalidInput)
	}
	guide, err := s.catalog.GetSizeGuide(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting size guide for blueprint %d: %w", id, err)
	}
	return guide, nil
}

// ListPrintProviders returns all print providers.
func (s *CatalogService) ListPrintProviders(ctx context.Context) ([]domain.PrintProvider, error) {
	providers, err := s.catalog.ListPrintProviders(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing print providers: %w", err)
	}
	return providers, nil
}

// GetPrintProvider returns one print provider.
func (s *CatalogService) GetPrintProvider(ctx context.Context, id domain.PrintProviderID) (*domain.PrintProvider, error) {
	if !id.Valid() {
		return nil, fmt.Errorf("print provider id %d: %w", id, domain.ErrInvalidInput)
	}
	provider, err := s.catalog.GetPrintProvider(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting print provider %d: %w", id, err)
	}
	return provider, nil
}

// ListProvidersForBlueprint returns the providers fulfilling a blueprint.
func (s *CatalogService) ListProvidersForBlueprint(ctx context.Context, blueprint domain.BlueprintID) ([]domain.PrintProviderRef, error) {
	if !blueprint.Valid() {
		return nil, fmt.Errorf("blueprint id %d: %w", blueprint, domain.ErrInvalidInput)
	}
	refs, err := s.catalog.ListProvidersForBlueprint(ctx, blueprint)
	if err != nil {
		return nil, fmt.Errorf("listing providers for blueprint %d: %w", blueprint, err)
	}
	return refs, nil
}

// GetVariants returns the variants of a blueprint for a provider.
func (s *CatalogService) GetVariants(ctx context.Context, blueprint domain.BlueprintID, provider domain.PrintProviderID, includeOutOfStock bool) (*domain.Variants, error) {
	if !blueprint.Valid() {
		return nil, fmt.Errorf("blueprint id %d: %w", blueprint, domain.ErrInvalidInput)
	}
	if !provider.Valid() {
		return nil, fmt.Errorf("print provider id %d: %w", provider, domain.ErrInvalidInput)
	}
	v, err := s.catalog.GetVariants(ctx, blueprint, provider, includeOutOfStock)
	if err != nil {
		return nil, fmt.Errorf("getting variants for blueprint %d provider %d: %w", blueprint, provider, err)
	}
	return v, nil
}

// GetShipping returns shipping information for a blueprint.
func (s *CatalogService) GetShipping(ctx context.Context, blueprint domain.BlueprintID) (*domain.BlueprintShipping, error) {
	if !blueprint.Valid() {
		return nil, fmt.Errorf("blueprint id %d: %w", blueprint, domain.ErrInvalidInput)
	}
	shipping, err := s.catalog.GetShipping(ctx, blueprint)
	if err != nil {
		return nil, fmt.Errorf("getting shipping for blueprint %d: %w", blueprint, err)
	}
	return shipping, nil
}

// GetShippingMethod returns raw v2 shipping data for one shipping method.
func (s *CatalogService) GetShippingMethod(ctx context.Context, blueprint domain.BlueprintID, method domain.ShippingMethod) (map[string]any, error) {
	if !blueprint.Valid() {
		return nil, fmt.Errorf("blueprint id %d: %w", blueprint, domain.ErrInvalidInput)
	}
	data, err := s.catalog.GetShippingMethod(ctx, blueprint, method)
	if err != nil {
		return nil, fmt.Errorf("getting shipping method %q for blueprint %d: %w", method, blueprint, err)
	}
	return data, nil
}

// GetShippingV2 returns the raw v2 shipping list for a blueprint.
func (s *CatalogService) GetShippingV2(ctx context.Context, blueprint domain.BlueprintID, provider domain.PrintProviderID) (map[string]any, error) {
	if !blueprint.Valid() {
		return nil, fmt.Errorf("blueprint id %d: %w", blueprint, domain.ErrInvalidInput)
	}
	if !provider.Valid() {
		return nil, fmt.Errorf("print provider id %d: %w", provider, domain.ErrInvalidInput)
	}
	data, err := s.catalog.GetShippingV2(ctx, blueprint, provider)
	if err != nil {
		return nil, fmt.Errorf("getting v2 shipping for blueprint %d provider %d: %w", blueprint, provider, err)
	}
	return data, nil
}
