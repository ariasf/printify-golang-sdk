package services

import (
	"context"
	"errors"
	"testing"

	"github.com/printify-go/pkg/printify/domain"
)

type stubCatalogGateway struct {
	listBlueprintsFn     func(ctx context.Context) ([]domain.Blueprint, error)
	getBlueprintFn       func(ctx context.Context, id domain.BlueprintID) (*domain.Blueprint, error)
	getSizeGuideFn       func(ctx context.Context, id domain.BlueprintID) (*domain.SizeGuide, error)
	listPrintProvidersFn func(ctx context.Context) ([]domain.PrintProvider, error)
	getPrintProviderFn   func(ctx context.Context, id domain.PrintProviderID) (*domain.PrintProvider, error)
	listProvidersForBpFn func(ctx context.Context, blueprint domain.BlueprintID) ([]domain.PrintProviderRef, error)
	getVariantsFn        func(ctx context.Context, blueprint domain.BlueprintID, provider domain.PrintProviderID, includeOutOfStock bool) (*domain.Variants, error)
	getShippingFn        func(ctx context.Context, blueprint domain.BlueprintID) (*domain.BlueprintShipping, error)
	getShippingMethodFn  func(ctx context.Context, blueprint domain.BlueprintID, method domain.ShippingMethod) (map[string]any, error)
	getShippingV2Fn      func(ctx context.Context, blueprint domain.BlueprintID, provider domain.PrintProviderID) (map[string]any, error)
}

func (s stubCatalogGateway) ListBlueprints(ctx context.Context) ([]domain.Blueprint, error) {
	return s.listBlueprintsFn(ctx)
}

func (s stubCatalogGateway) GetBlueprint(ctx context.Context, id domain.BlueprintID) (*domain.Blueprint, error) {
	return s.getBlueprintFn(ctx, id)
}

func (s stubCatalogGateway) GetSizeGuide(ctx context.Context, id domain.BlueprintID) (*domain.SizeGuide, error) {
	return s.getSizeGuideFn(ctx, id)
}

func (s stubCatalogGateway) ListPrintProviders(ctx context.Context) ([]domain.PrintProvider, error) {
	return s.listPrintProvidersFn(ctx)
}

func (s stubCatalogGateway) GetPrintProvider(ctx context.Context, id domain.PrintProviderID) (*domain.PrintProvider, error) {
	return s.getPrintProviderFn(ctx, id)
}

func (s stubCatalogGateway) ListProvidersForBlueprint(ctx context.Context, blueprint domain.BlueprintID) ([]domain.PrintProviderRef, error) {
	return s.listProvidersForBpFn(ctx, blueprint)
}

func (s stubCatalogGateway) GetVariants(ctx context.Context, blueprint domain.BlueprintID, provider domain.PrintProviderID, includeOutOfStock bool) (*domain.Variants, error) {
	return s.getVariantsFn(ctx, blueprint, provider, includeOutOfStock)
}

func (s stubCatalogGateway) GetShipping(ctx context.Context, blueprint domain.BlueprintID) (*domain.BlueprintShipping, error) {
	return s.getShippingFn(ctx, blueprint)
}

func (s stubCatalogGateway) GetShippingMethod(ctx context.Context, blueprint domain.BlueprintID, method domain.ShippingMethod) (map[string]any, error) {
	return s.getShippingMethodFn(ctx, blueprint, method)
}

func (s stubCatalogGateway) GetShippingV2(ctx context.Context, blueprint domain.BlueprintID, provider domain.PrintProviderID) (map[string]any, error) {
	return s.getShippingV2Fn(ctx, blueprint, provider)
}

func TestCatalogService_ListBlueprints(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		gateway stubCatalogGateway
		want    int
		wantErr error
	}{
		{
			name: "returns blueprints from gateway",
			gateway: stubCatalogGateway{listBlueprintsFn: func(context.Context) ([]domain.Blueprint, error) {
				return []domain.Blueprint{{ID: 101, Title: "T-Shirt"}}, nil
			}},
			want: 1,
		},
		{
			name: "propagates gateway error",
			gateway: stubCatalogGateway{listBlueprintsFn: func(context.Context) ([]domain.Blueprint, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewCatalogService(tt.gateway)
			got, err := svc.ListBlueprints(context.Background())
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if len(got) != tt.want {
				t.Fatalf("len(blueprints) = %d, want %d", len(got), tt.want)
			}
		})
	}
}

func TestCatalogService_GetBlueprint(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		id      domain.BlueprintID
		gateway stubCatalogGateway
		wantErr error
	}{
		{
			name: "rejects non-positive id without calling gateway",
			id:   0,
			gateway: stubCatalogGateway{getBlueprintFn: func(context.Context, domain.BlueprintID) (*domain.Blueprint, error) {
				t.Fatal("gateway must not be called")
				return nil, nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "returns blueprint for valid id",
			id:   101,
			gateway: stubCatalogGateway{getBlueprintFn: func(_ context.Context, id domain.BlueprintID) (*domain.Blueprint, error) {
				if id != 101 {
					t.Fatalf("gateway got id %d, want 101", id)
				}
				return &domain.Blueprint{ID: 101}, nil
			}},
		},
		{
			name: "propagates gateway error",
			id:   101,
			gateway: stubCatalogGateway{getBlueprintFn: func(context.Context, domain.BlueprintID) (*domain.Blueprint, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewCatalogService(tt.gateway)
			got, err := svc.GetBlueprint(context.Background(), tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got == nil {
				t.Fatal("blueprint = nil, want non-nil")
			}
		})
	}
}

func TestCatalogService_GetSizeGuide(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		id      domain.BlueprintID
		gateway stubCatalogGateway
		wantErr error
	}{
		{
			name: "rejects non-positive id without calling gateway",
			id:   0,
			gateway: stubCatalogGateway{getSizeGuideFn: func(context.Context, domain.BlueprintID) (*domain.SizeGuide, error) {
				t.Fatal("gateway must not be called")
				return nil, nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "returns size guide for valid id",
			id:   101,
			gateway: stubCatalogGateway{getSizeGuideFn: func(_ context.Context, id domain.BlueprintID) (*domain.SizeGuide, error) {
				if id != 101 {
					t.Fatalf("gateway got id %d, want 101", id)
				}
				return &domain.SizeGuide{}, nil
			}},
		},
		{
			name: "propagates gateway error",
			id:   101,
			gateway: stubCatalogGateway{getSizeGuideFn: func(context.Context, domain.BlueprintID) (*domain.SizeGuide, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewCatalogService(tt.gateway)
			got, err := svc.GetSizeGuide(context.Background(), tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got == nil {
				t.Fatal("size guide = nil, want non-nil")
			}
		})
	}
}

func TestCatalogService_ListPrintProviders(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		gateway stubCatalogGateway
		want    int
		wantErr error
	}{
		{
			name: "returns print providers from gateway",
			gateway: stubCatalogGateway{listPrintProvidersFn: func(context.Context) ([]domain.PrintProvider, error) {
				return []domain.PrintProvider{{ID: 9, Title: "Provider"}}, nil
			}},
			want: 1,
		},
		{
			name: "propagates gateway error",
			gateway: stubCatalogGateway{listPrintProvidersFn: func(context.Context) ([]domain.PrintProvider, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewCatalogService(tt.gateway)
			got, err := svc.ListPrintProviders(context.Background())
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if len(got) != tt.want {
				t.Fatalf("len(providers) = %d, want %d", len(got), tt.want)
			}
		})
	}
}

func TestCatalogService_GetPrintProvider(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		id      domain.PrintProviderID
		gateway stubCatalogGateway
		wantErr error
	}{
		{
			name: "rejects non-positive id without calling gateway",
			id:   0,
			gateway: stubCatalogGateway{getPrintProviderFn: func(context.Context, domain.PrintProviderID) (*domain.PrintProvider, error) {
				t.Fatal("gateway must not be called")
				return nil, nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "returns provider for valid id",
			id:   9,
			gateway: stubCatalogGateway{getPrintProviderFn: func(_ context.Context, id domain.PrintProviderID) (*domain.PrintProvider, error) {
				if id != 9 {
					t.Fatalf("gateway got id %d, want 9", id)
				}
				return &domain.PrintProvider{ID: 9}, nil
			}},
		},
		{
			name: "propagates gateway error",
			id:   9,
			gateway: stubCatalogGateway{getPrintProviderFn: func(context.Context, domain.PrintProviderID) (*domain.PrintProvider, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewCatalogService(tt.gateway)
			got, err := svc.GetPrintProvider(context.Background(), tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got == nil {
				t.Fatal("provider = nil, want non-nil")
			}
		})
	}
}

func TestCatalogService_ListProvidersForBlueprint(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		id      domain.BlueprintID
		gateway stubCatalogGateway
		want    int
		wantErr error
	}{
		{
			name: "rejects non-positive id without calling gateway",
			id:   0,
			gateway: stubCatalogGateway{listProvidersForBpFn: func(context.Context, domain.BlueprintID) ([]domain.PrintProviderRef, error) {
				t.Fatal("gateway must not be called")
				return nil, nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "returns provider refs for valid id",
			id:   101,
			gateway: stubCatalogGateway{listProvidersForBpFn: func(_ context.Context, id domain.BlueprintID) ([]domain.PrintProviderRef, error) {
				if id != 101 {
					t.Fatalf("gateway got id %d, want 101", id)
				}
				return []domain.PrintProviderRef{{ID: 9}}, nil
			}},
			want: 1,
		},
		{
			name: "propagates gateway error",
			id:   101,
			gateway: stubCatalogGateway{listProvidersForBpFn: func(context.Context, domain.BlueprintID) ([]domain.PrintProviderRef, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewCatalogService(tt.gateway)
			got, err := svc.ListProvidersForBlueprint(context.Background(), tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if len(got) != tt.want {
				t.Fatalf("len(refs) = %d, want %d", len(got), tt.want)
			}
		})
	}
}

func TestCatalogService_GetVariants(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name      string
		blueprint domain.BlueprintID
		provider  domain.PrintProviderID
		gateway   stubCatalogGateway
		wantErr   error
	}{
		{
			name:      "rejects invalid blueprint without calling gateway",
			blueprint: 0,
			provider:  9,
			gateway: stubCatalogGateway{getVariantsFn: func(context.Context, domain.BlueprintID, domain.PrintProviderID, bool) (*domain.Variants, error) {
				t.Fatal("gateway must not be called")
				return nil, nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:      "rejects invalid provider without calling gateway",
			blueprint: 101,
			provider:  0,
			gateway: stubCatalogGateway{getVariantsFn: func(context.Context, domain.BlueprintID, domain.PrintProviderID, bool) (*domain.Variants, error) {
				t.Fatal("gateway must not be called")
				return nil, nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:      "passes both ids and out-of-stock flag through",
			blueprint: 101,
			provider:  9,
			gateway: stubCatalogGateway{getVariantsFn: func(_ context.Context, bp domain.BlueprintID, pp domain.PrintProviderID, outOfStock bool) (*domain.Variants, error) {
				if bp != 101 || pp != 9 || !outOfStock {
					t.Fatalf("gateway got blueprint=%d provider=%d outOfStock=%v, want 101 9 true", bp, pp, outOfStock)
				}
				return &domain.Variants{}, nil
			}},
		},
		{
			name:      "propagates gateway error",
			blueprint: 101,
			provider:  9,
			gateway: stubCatalogGateway{getVariantsFn: func(context.Context, domain.BlueprintID, domain.PrintProviderID, bool) (*domain.Variants, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewCatalogService(tt.gateway)
			got, err := svc.GetVariants(context.Background(), tt.blueprint, tt.provider, true)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got == nil {
				t.Fatal("variants = nil, want non-nil")
			}
		})
	}
}

func TestCatalogService_GetShipping(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		id      domain.BlueprintID
		gateway stubCatalogGateway
		wantErr error
	}{
		{
			name: "rejects non-positive id without calling gateway",
			id:   0,
			gateway: stubCatalogGateway{getShippingFn: func(context.Context, domain.BlueprintID) (*domain.BlueprintShipping, error) {
				t.Fatal("gateway must not be called")
				return nil, nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "returns shipping for valid id",
			id:   101,
			gateway: stubCatalogGateway{getShippingFn: func(_ context.Context, id domain.BlueprintID) (*domain.BlueprintShipping, error) {
				if id != 101 {
					t.Fatalf("gateway got id %d, want 101", id)
				}
				return &domain.BlueprintShipping{}, nil
			}},
		},
		{
			name: "propagates gateway error",
			id:   101,
			gateway: stubCatalogGateway{getShippingFn: func(context.Context, domain.BlueprintID) (*domain.BlueprintShipping, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewCatalogService(tt.gateway)
			got, err := svc.GetShipping(context.Background(), tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got == nil {
				t.Fatal("shipping = nil, want non-nil")
			}
		})
	}
}

func TestCatalogService_GetShippingMethod(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		id      domain.BlueprintID
		gateway stubCatalogGateway
		wantErr error
	}{
		{
			name: "rejects non-positive id without calling gateway",
			id:   0,
			gateway: stubCatalogGateway{getShippingMethodFn: func(context.Context, domain.BlueprintID, domain.ShippingMethod) (map[string]any, error) {
				t.Fatal("gateway must not be called")
				return nil, nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "returns raw data for valid id and method",
			id:   101,
			gateway: stubCatalogGateway{getShippingMethodFn: func(_ context.Context, id domain.BlueprintID, method domain.ShippingMethod) (map[string]any, error) {
				if id != 101 || method != domain.ShippingStandard {
					t.Fatalf("gateway got id %d method %q, want 101 %q", id, method, domain.ShippingStandard)
				}
				return map[string]any{"cost": 5}, nil
			}},
		},
		{
			name: "propagates gateway error",
			id:   101,
			gateway: stubCatalogGateway{getShippingMethodFn: func(context.Context, domain.BlueprintID, domain.ShippingMethod) (map[string]any, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewCatalogService(tt.gateway)
			got, err := svc.GetShippingMethod(context.Background(), tt.id, domain.ShippingStandard)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && len(got) == 0 {
				t.Fatal("data = empty, want non-empty")
			}
		})
	}
}

func TestCatalogService_GetShippingV2(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name      string
		blueprint domain.BlueprintID
		provider  domain.PrintProviderID
		gateway   stubCatalogGateway
		wantErr   error
	}{
		{
			name:      "rejects invalid blueprint without calling gateway",
			blueprint: 0,
			provider:  9,
			gateway: stubCatalogGateway{getShippingV2Fn: func(context.Context, domain.BlueprintID, domain.PrintProviderID) (map[string]any, error) {
				t.Fatal("gateway must not be called")
				return nil, nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:      "rejects invalid provider without calling gateway",
			blueprint: 101,
			provider:  0,
			gateway: stubCatalogGateway{getShippingV2Fn: func(context.Context, domain.BlueprintID, domain.PrintProviderID) (map[string]any, error) {
				t.Fatal("gateway must not be called")
				return nil, nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:      "returns raw v2 data",
			blueprint: 101,
			provider:  9,
			gateway: stubCatalogGateway{getShippingV2Fn: func(_ context.Context, bp domain.BlueprintID, pp domain.PrintProviderID) (map[string]any, error) {
				if bp != 101 || pp != 9 {
					t.Fatalf("gateway got blueprint=%d provider=%d, want 101 9", bp, pp)
				}
				return map[string]any{"shipping": []any{}}, nil
			}},
		},
		{
			name:      "propagates gateway error",
			blueprint: 101,
			provider:  9,
			gateway: stubCatalogGateway{getShippingV2Fn: func(context.Context, domain.BlueprintID, domain.PrintProviderID) (map[string]any, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewCatalogService(tt.gateway)
			got, err := svc.GetShippingV2(context.Background(), tt.blueprint, tt.provider)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && len(got) == 0 {
				t.Fatal("data = empty, want non-empty")
			}
		})
	}
}
