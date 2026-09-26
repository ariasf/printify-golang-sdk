package driving

import (
	"context"

	"github.com/printify-go/pkg/printify/domain"
)

// CatalogService exposes catalog browsing to SDK consumers.
type CatalogService interface {
	ListBlueprints(ctx context.Context) ([]domain.Blueprint, error)
	GetBlueprint(ctx context.Context, id domain.BlueprintID) (*domain.Blueprint, error)
	GetSizeGuide(ctx context.Context, id domain.BlueprintID) (*domain.SizeGuide, error)
	ListPrintProviders(ctx context.Context) ([]domain.PrintProvider, error)
	GetPrintProvider(ctx context.Context, id domain.PrintProviderID) (*domain.PrintProvider, error)
	ListProvidersForBlueprint(ctx context.Context, blueprint domain.BlueprintID) ([]domain.PrintProviderRef, error)
	GetVariants(ctx context.Context, blueprint domain.BlueprintID, provider domain.PrintProviderID, includeOutOfStock bool) (*domain.Variants, error)
	GetShipping(ctx context.Context, blueprint domain.BlueprintID) (*domain.BlueprintShipping, error)
	GetShippingMethod(ctx context.Context, blueprint domain.BlueprintID, method domain.ShippingMethod) (map[string]any, error)
	GetShippingV2(ctx context.Context, blueprint domain.BlueprintID, provider domain.PrintProviderID) (map[string]any, error)
}
