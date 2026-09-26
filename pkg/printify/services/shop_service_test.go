package services

import (
	"context"
	"errors"
	"testing"

	"github.com/printify-go/pkg/printify/domain"
)

type stubShopGateway struct {
	listFn       func(ctx context.Context) ([]domain.Shop, error)
	disconnectFn func(ctx context.Context, id domain.ShopID) error
}

func (s stubShopGateway) ListShops(ctx context.Context) ([]domain.Shop, error) {
	return s.listFn(ctx)
}

func (s stubShopGateway) DisconnectShop(ctx context.Context, id domain.ShopID) error {
	return s.disconnectFn(ctx, id)
}

func TestShopService_ListShops(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		gateway stubShopGateway
		want    int
		wantErr error
	}{
		{
			name: "returns shops from gateway",
			gateway: stubShopGateway{listFn: func(context.Context) ([]domain.Shop, error) {
				return []domain.Shop{{ID: 1, Title: "My Store", Channel: "shopify"}}, nil
			}},
			want: 1,
		},
		{
			name: "propagates gateway error",
			gateway: stubShopGateway{listFn: func(context.Context) ([]domain.Shop, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewShopService(tt.gateway)
			got, err := svc.ListShops(context.Background())
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if len(got) != tt.want {
				t.Fatalf("len(shops) = %d, want %d", len(got), tt.want)
			}
		})
	}
}

func TestShopService_DisconnectShop(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		id      domain.ShopID
		gateway stubShopGateway
		wantErr error
	}{
		{
			name:    "rejects non-positive id without calling gateway",
			id:      0,
			gateway: stubShopGateway{disconnectFn: func(context.Context, domain.ShopID) error { t.Fatal("gateway must not be called"); return nil }},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "disconnects valid shop",
			id:   42,
			gateway: stubShopGateway{disconnectFn: func(_ context.Context, id domain.ShopID) error {
				if id != 42 {
					t.Fatalf("gateway got id %d, want 42", id)
				}
				return nil
			}},
		},
		{
			name:    "propagates gateway error",
			id:      42,
			gateway: stubShopGateway{disconnectFn: func(context.Context, domain.ShopID) error { return gatewayErr }},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewShopService(tt.gateway)
			err := svc.DisconnectShop(context.Background(), tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
