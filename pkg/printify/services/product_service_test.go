package services

import (
	"context"
	"errors"
	"testing"

	"github.com/printify-go/pkg/printify/domain"
	"github.com/printify-go/pkg/printify/ports/driving"
)

type stubProductGateway struct {
	listFn           func(ctx context.Context, shop domain.ShopID, f driving.ProductFilter) ([]domain.Product, error)
	getFn            func(ctx context.Context, shop domain.ShopID, id domain.ProductID) (*domain.Product, error)
	createFn         func(ctx context.Context, shop domain.ShopID, p domain.CreateProduct) (*domain.Product, error)
	updateFn         func(ctx context.Context, shop domain.ShopID, id domain.ProductID, p domain.UpdateProduct) (*domain.Product, error)
	deleteFn         func(ctx context.Context, shop domain.ShopID, id domain.ProductID) error
	publishFn        func(ctx context.Context, shop domain.ShopID, id domain.ProductID, opts domain.PublishProduct) error
	unpublishFn      func(ctx context.Context, shop domain.ShopID, id domain.ProductID) error
	setPublishSuccFn func(ctx context.Context, shop domain.ShopID, id domain.ProductID, ext domain.PublishSucceeded) error
	setPublishFailFn func(ctx context.Context, shop domain.ShopID, id domain.ProductID, failed domain.PublishFailed) error
	listGpsrFn       func(ctx context.Context, shop domain.ShopID, id domain.ProductID) ([]domain.GpsrInfo, error)
}

func (s stubProductGateway) List(ctx context.Context, shop domain.ShopID, f driving.ProductFilter) ([]domain.Product, error) {
	return s.listFn(ctx, shop, f)
}

func (s stubProductGateway) Get(ctx context.Context, shop domain.ShopID, id domain.ProductID) (*domain.Product, error) {
	return s.getFn(ctx, shop, id)
}

func (s stubProductGateway) Create(ctx context.Context, shop domain.ShopID, p domain.CreateProduct) (*domain.Product, error) {
	return s.createFn(ctx, shop, p)
}

func (s stubProductGateway) Update(ctx context.Context, shop domain.ShopID, id domain.ProductID, p domain.UpdateProduct) (*domain.Product, error) {
	return s.updateFn(ctx, shop, id, p)
}

func (s stubProductGateway) Delete(ctx context.Context, shop domain.ShopID, id domain.ProductID) error {
	return s.deleteFn(ctx, shop, id)
}

func (s stubProductGateway) Publish(ctx context.Context, shop domain.ShopID, id domain.ProductID, opts domain.PublishProduct) error {
	return s.publishFn(ctx, shop, id, opts)
}

func (s stubProductGateway) Unpublish(ctx context.Context, shop domain.ShopID, id domain.ProductID) error {
	return s.unpublishFn(ctx, shop, id)
}

func (s stubProductGateway) SetPublishSucceeded(ctx context.Context, shop domain.ShopID, id domain.ProductID, ext domain.PublishSucceeded) error {
	return s.setPublishSuccFn(ctx, shop, id, ext)
}

func (s stubProductGateway) SetPublishFailed(ctx context.Context, shop domain.ShopID, id domain.ProductID, failed domain.PublishFailed) error {
	return s.setPublishFailFn(ctx, shop, id, failed)
}

func (s stubProductGateway) ListGpsr(ctx context.Context, shop domain.ShopID, id domain.ProductID) ([]domain.GpsrInfo, error) {
	return s.listGpsrFn(ctx, shop, id)
}

func validCreateProduct() domain.CreateProduct {
	return domain.CreateProduct{
		Title:           "T-Shirt",
		BlueprintID:     101,
		PrintProviderID: 9,
		Variants:        []domain.ProductVariant{{ID: 1}},
		PrintAreas:      []domain.PrintArea{{VariantIDs: []domain.VariantID{1}}},
	}
}

func TestProductService_List(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		shop    domain.ShopID
		gateway stubProductGateway
		want    int
		wantErr error
	}{
		{
			name: "rejects non-positive shop without calling gateway",
			shop: 0,
			gateway: stubProductGateway{listFn: func(context.Context, domain.ShopID, driving.ProductFilter) ([]domain.Product, error) {
				t.Fatal("gateway must not be called")
				return nil, nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "returns products and passes filter through",
			shop: 1,
			gateway: stubProductGateway{listFn: func(_ context.Context, shop domain.ShopID, f driving.ProductFilter) ([]domain.Product, error) {
				if shop != 1 {
					t.Fatalf("gateway got shop %d, want 1", shop)
				}
				if f.Limit == nil || *f.Limit != 25 {
					t.Fatal("gateway got unexpected filter")
				}
				return []domain.Product{{ID: "abc"}}, nil
			}},
			want: 1,
		},
		{
			name: "propagates gateway error",
			shop: 1,
			gateway: stubProductGateway{listFn: func(context.Context, domain.ShopID, driving.ProductFilter) ([]domain.Product, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewProductService(tt.gateway)
			limit := 25
			got, err := svc.List(context.Background(), tt.shop, driving.ProductFilter{Limit: &limit})
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if len(got) != tt.want {
				t.Fatalf("len(products) = %d, want %d", len(got), tt.want)
			}
		})
	}
}

func TestProductService_Get(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		shop    domain.ShopID
		id      domain.ProductID
		gateway stubProductGateway
		wantErr error
	}{
		{
			name: "rejects non-positive shop",
			shop: 0,
			id:   "abc",
			gateway: stubProductGateway{getFn: func(context.Context, domain.ShopID, domain.ProductID) (*domain.Product, error) {
				t.Fatal("gateway must not be called")
				return nil, nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects empty product id",
			shop: 1,
			id:   "",
			gateway: stubProductGateway{getFn: func(context.Context, domain.ShopID, domain.ProductID) (*domain.Product, error) {
				t.Fatal("gateway must not be called")
				return nil, nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "returns product for valid ids",
			shop: 1,
			id:   "abc",
			gateway: stubProductGateway{getFn: func(_ context.Context, shop domain.ShopID, id domain.ProductID) (*domain.Product, error) {
				if shop != 1 || id != "abc" {
					t.Fatalf("gateway got shop=%d id=%q, want 1 abc", shop, id)
				}
				return &domain.Product{ID: "abc"}, nil
			}},
		},
		{
			name: "propagates gateway error",
			shop: 1,
			id:   "abc",
			gateway: stubProductGateway{getFn: func(context.Context, domain.ShopID, domain.ProductID) (*domain.Product, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewProductService(tt.gateway)
			got, err := svc.Get(context.Background(), tt.shop, tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got == nil {
				t.Fatal("product = nil, want non-nil")
			}
		})
	}
}

func TestProductService_Create(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		shop    domain.ShopID
		payload domain.CreateProduct
		gateway stubProductGateway
		wantErr error
	}{
		{
			name:    "rejects non-positive shop",
			shop:    0,
			payload: validCreateProduct(),
			gateway: stubProductGateway{createFn: func(context.Context, domain.ShopID, domain.CreateProduct) (*domain.Product, error) {
				t.Fatal("gateway must not be called")
				return nil, nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects empty title",
			shop: 1,
			payload: func() domain.CreateProduct {
				p := validCreateProduct()
				p.Title = ""
				return p
			}(),
			gateway: stubProductGateway{createFn: func(context.Context, domain.ShopID, domain.CreateProduct) (*domain.Product, error) {
				t.Fatal("gateway must not be called")
				return nil, nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects invalid blueprint id",
			shop: 1,
			payload: func() domain.CreateProduct {
				p := validCreateProduct()
				p.BlueprintID = 0
				return p
			}(),
			gateway: stubProductGateway{createFn: func(context.Context, domain.ShopID, domain.CreateProduct) (*domain.Product, error) {
				t.Fatal("gateway must not be called")
				return nil, nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects invalid print provider id",
			shop: 1,
			payload: func() domain.CreateProduct {
				p := validCreateProduct()
				p.PrintProviderID = 0
				return p
			}(),
			gateway: stubProductGateway{createFn: func(context.Context, domain.ShopID, domain.CreateProduct) (*domain.Product, error) {
				t.Fatal("gateway must not be called")
				return nil, nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects missing variants",
			shop: 1,
			payload: func() domain.CreateProduct {
				p := validCreateProduct()
				p.Variants = nil
				return p
			}(),
			gateway: stubProductGateway{createFn: func(context.Context, domain.ShopID, domain.CreateProduct) (*domain.Product, error) {
				t.Fatal("gateway must not be called")
				return nil, nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects missing print areas",
			shop: 1,
			payload: func() domain.CreateProduct {
				p := validCreateProduct()
				p.PrintAreas = nil
				return p
			}(),
			gateway: stubProductGateway{createFn: func(context.Context, domain.ShopID, domain.CreateProduct) (*domain.Product, error) {
				t.Fatal("gateway must not be called")
				return nil, nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "creates valid product",
			shop: 1,
			payload: func() domain.CreateProduct {
				p := validCreateProduct()
				return p
			}(),
			gateway: stubProductGateway{createFn: func(_ context.Context, shop domain.ShopID, p domain.CreateProduct) (*domain.Product, error) {
				if shop != 1 || p.Title != "T-Shirt" {
					t.Fatalf("gateway got unexpected shop=%d title=%q", shop, p.Title)
				}
				return &domain.Product{ID: "new"}, nil
			}},
		},
		{
			name:    "propagates gateway error",
			shop:    1,
			payload: validCreateProduct(),
			gateway: stubProductGateway{createFn: func(context.Context, domain.ShopID, domain.CreateProduct) (*domain.Product, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewProductService(tt.gateway)
			got, err := svc.Create(context.Background(), tt.shop, tt.payload)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got == nil {
				t.Fatal("product = nil, want non-nil")
			}
		})
	}
}

func TestProductService_Update(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		shop    domain.ShopID
		id      domain.ProductID
		gateway stubProductGateway
		wantErr error
	}{
		{
			name: "rejects non-positive shop",
			shop: 0,
			id:   "abc",
			gateway: stubProductGateway{updateFn: func(context.Context, domain.ShopID, domain.ProductID, domain.UpdateProduct) (*domain.Product, error) {
				t.Fatal("gateway must not be called")
				return nil, nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects empty product id",
			shop: 1,
			id:   "",
			gateway: stubProductGateway{updateFn: func(context.Context, domain.ShopID, domain.ProductID, domain.UpdateProduct) (*domain.Product, error) {
				t.Fatal("gateway must not be called")
				return nil, nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "updates valid product",
			shop: 1,
			id:   "abc",
			gateway: stubProductGateway{updateFn: func(_ context.Context, shop domain.ShopID, id domain.ProductID, p domain.UpdateProduct) (*domain.Product, error) {
				if shop != 1 || id != "abc" {
					t.Fatalf("gateway got shop=%d id=%q, want 1 abc", shop, id)
				}
				return &domain.Product{ID: "abc", Title: p.Title}, nil
			}},
		},
		{
			name: "propagates gateway error",
			shop: 1,
			id:   "abc",
			gateway: stubProductGateway{updateFn: func(context.Context, domain.ShopID, domain.ProductID, domain.UpdateProduct) (*domain.Product, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewProductService(tt.gateway)
			got, err := svc.Update(context.Background(), tt.shop, tt.id, domain.UpdateProduct{Title: "New Title"})
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got == nil {
				t.Fatal("product = nil, want non-nil")
			}
		})
	}
}

func TestProductService_Delete(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		shop    domain.ShopID
		id      domain.ProductID
		gateway stubProductGateway
		wantErr error
	}{
		{
			name: "rejects non-positive shop",
			shop: 0,
			id:   "abc",
			gateway: stubProductGateway{deleteFn: func(context.Context, domain.ShopID, domain.ProductID) error {
				t.Fatal("gateway must not be called")
				return nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects empty product id",
			shop: 1,
			id:   "",
			gateway: stubProductGateway{deleteFn: func(context.Context, domain.ShopID, domain.ProductID) error {
				t.Fatal("gateway must not be called")
				return nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "deletes valid product",
			shop: 1,
			id:   "abc",
			gateway: stubProductGateway{deleteFn: func(_ context.Context, shop domain.ShopID, id domain.ProductID) error {
				if shop != 1 || id != "abc" {
					t.Fatalf("gateway got shop=%d id=%q, want 1 abc", shop, id)
				}
				return nil
			}},
		},
		{
			name:    "propagates gateway error",
			shop:    1,
			id:      "abc",
			gateway: stubProductGateway{deleteFn: func(context.Context, domain.ShopID, domain.ProductID) error { return gatewayErr }},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewProductService(tt.gateway)
			err := svc.Delete(context.Background(), tt.shop, tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestProductService_Publish(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		shop    domain.ShopID
		id      domain.ProductID
		gateway stubProductGateway
		wantErr error
	}{
		{
			name: "rejects non-positive shop",
			shop: 0,
			id:   "abc",
			gateway: stubProductGateway{publishFn: func(context.Context, domain.ShopID, domain.ProductID, domain.PublishProduct) error {
				t.Fatal("gateway must not be called")
				return nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects empty product id",
			shop: 1,
			id:   "",
			gateway: stubProductGateway{publishFn: func(context.Context, domain.ShopID, domain.ProductID, domain.PublishProduct) error {
				t.Fatal("gateway must not be called")
				return nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "publishes valid product",
			shop: 1,
			id:   "abc",
			gateway: stubProductGateway{publishFn: func(_ context.Context, shop domain.ShopID, id domain.ProductID, opts domain.PublishProduct) error {
				if shop != 1 || id != "abc" {
					t.Fatalf("gateway got shop=%d id=%q, want 1 abc", shop, id)
				}
				return nil
			}},
		},
		{
			name: "propagates gateway error",
			shop: 1,
			id:   "abc",
			gateway: stubProductGateway{publishFn: func(context.Context, domain.ShopID, domain.ProductID, domain.PublishProduct) error {
				return gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewProductService(tt.gateway)
			err := svc.Publish(context.Background(), tt.shop, tt.id, domain.PublishProduct{})
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestProductService_Unpublish(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		shop    domain.ShopID
		id      domain.ProductID
		gateway stubProductGateway
		wantErr error
	}{
		{
			name: "rejects non-positive shop",
			shop: 0,
			id:   "abc",
			gateway: stubProductGateway{unpublishFn: func(context.Context, domain.ShopID, domain.ProductID) error {
				t.Fatal("gateway must not be called")
				return nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects empty product id",
			shop: 1,
			id:   "",
			gateway: stubProductGateway{unpublishFn: func(context.Context, domain.ShopID, domain.ProductID) error {
				t.Fatal("gateway must not be called")
				return nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "unpublishes valid product",
			shop: 1,
			id:   "abc",
			gateway: stubProductGateway{unpublishFn: func(_ context.Context, shop domain.ShopID, id domain.ProductID) error {
				if shop != 1 || id != "abc" {
					t.Fatalf("gateway got shop=%d id=%q, want 1 abc", shop, id)
				}
				return nil
			}},
		},
		{
			name:    "propagates gateway error",
			shop:    1,
			id:      "abc",
			gateway: stubProductGateway{unpublishFn: func(context.Context, domain.ShopID, domain.ProductID) error { return gatewayErr }},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewProductService(tt.gateway)
			err := svc.Unpublish(context.Background(), tt.shop, tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestProductService_SetPublishSucceeded(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		shop    domain.ShopID
		id      domain.ProductID
		gateway stubProductGateway
		wantErr error
	}{
		{
			name: "rejects non-positive shop",
			shop: 0,
			id:   "abc",
			gateway: stubProductGateway{setPublishSuccFn: func(context.Context, domain.ShopID, domain.ProductID, domain.PublishSucceeded) error {
				t.Fatal("gateway must not be called")
				return nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects empty product id",
			shop: 1,
			id:   "",
			gateway: stubProductGateway{setPublishSuccFn: func(context.Context, domain.ShopID, domain.ProductID, domain.PublishSucceeded) error {
				t.Fatal("gateway must not be called")
				return nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "records successful publish",
			shop: 1,
			id:   "abc",
			gateway: stubProductGateway{setPublishSuccFn: func(_ context.Context, shop domain.ShopID, id domain.ProductID, ext domain.PublishSucceeded) error {
				if shop != 1 || id != "abc" || ext.ExternalID != "ch-1" {
					t.Fatalf("gateway got shop=%d id=%q ext=%+v", shop, id, ext)
				}
				return nil
			}},
		},
		{
			name: "propagates gateway error",
			shop: 1,
			id:   "abc",
			gateway: stubProductGateway{setPublishSuccFn: func(context.Context, domain.ShopID, domain.ProductID, domain.PublishSucceeded) error {
				return gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewProductService(tt.gateway)
			err := svc.SetPublishSucceeded(context.Background(), tt.shop, tt.id, domain.PublishSucceeded{ExternalID: "ch-1", ExternalHandle: "handle"})
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestProductService_SetPublishFailed(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		shop    domain.ShopID
		id      domain.ProductID
		reason  string
		gateway stubProductGateway
		wantErr error
	}{
		{
			name:   "rejects non-positive shop",
			shop:   0,
			id:     "abc",
			reason: "failed",
			gateway: stubProductGateway{setPublishFailFn: func(context.Context, domain.ShopID, domain.ProductID, domain.PublishFailed) error {
				t.Fatal("gateway must not be called")
				return nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:   "rejects empty product id",
			shop:   1,
			id:     "",
			reason: "failed",
			gateway: stubProductGateway{setPublishFailFn: func(context.Context, domain.ShopID, domain.ProductID, domain.PublishFailed) error {
				t.Fatal("gateway must not be called")
				return nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:   "rejects empty reason without calling gateway",
			shop:   1,
			id:     "abc",
			reason: "",
			gateway: stubProductGateway{setPublishFailFn: func(context.Context, domain.ShopID, domain.ProductID, domain.PublishFailed) error {
				t.Fatal("gateway must not be called")
				return nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:   "records failed publish with reason",
			shop:   1,
			id:     "abc",
			reason: "channel rejected",
			gateway: stubProductGateway{setPublishFailFn: func(_ context.Context, shop domain.ShopID, id domain.ProductID, failed domain.PublishFailed) error {
				if shop != 1 || id != "abc" || failed.Reason != "channel rejected" {
					t.Fatalf("gateway got shop=%d id=%q reason=%q", shop, id, failed.Reason)
				}
				return nil
			}},
		},
		{
			name:   "propagates gateway error",
			shop:   1,
			id:     "abc",
			reason: "failed",
			gateway: stubProductGateway{setPublishFailFn: func(context.Context, domain.ShopID, domain.ProductID, domain.PublishFailed) error {
				return gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewProductService(tt.gateway)
			err := svc.SetPublishFailed(context.Background(), tt.shop, tt.id, tt.reason)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestProductService_ListGpsr(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		shop    domain.ShopID
		id      domain.ProductID
		gateway stubProductGateway
		want    int
		wantErr error
	}{
		{
			name: "rejects non-positive shop",
			shop: 0,
			id:   "abc",
			gateway: stubProductGateway{listGpsrFn: func(context.Context, domain.ShopID, domain.ProductID) ([]domain.GpsrInfo, error) {
				t.Fatal("gateway must not be called")
				return nil, nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects empty product id",
			shop: 1,
			id:   "",
			gateway: stubProductGateway{listGpsrFn: func(context.Context, domain.ShopID, domain.ProductID) ([]domain.GpsrInfo, error) {
				t.Fatal("gateway must not be called")
				return nil, nil
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "returns gpsr info",
			shop: 1,
			id:   "abc",
			gateway: stubProductGateway{listGpsrFn: func(_ context.Context, shop domain.ShopID, id domain.ProductID) ([]domain.GpsrInfo, error) {
				if shop != 1 || id != "abc" {
					t.Fatalf("gateway got shop=%d id=%q, want 1 abc", shop, id)
				}
				return []domain.GpsrInfo{{Title: "gpsr"}}, nil
			}},
			want: 1,
		},
		{
			name: "propagates gateway error",
			shop: 1,
			id:   "abc",
			gateway: stubProductGateway{listGpsrFn: func(context.Context, domain.ShopID, domain.ProductID) ([]domain.GpsrInfo, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewProductService(tt.gateway)
			got, err := svc.ListGpsr(context.Background(), tt.shop, tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if len(got) != tt.want {
				t.Fatalf("len(gpsr) = %d, want %d", len(got), tt.want)
			}
		})
	}
}
