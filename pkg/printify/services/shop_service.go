// Package services contains all business logic. Services implement the
// driving ports and depend only on driven ports.
package services

import (
	"context"
	"fmt"

	"github.com/printify-go/pkg/printify/domain"
	"github.com/printify-go/pkg/printify/ports/driven"
)

// ShopService implements driving.ShopService.
type ShopService struct {
	shops driven.ShopGateway
}

// NewShopService builds a ShopService backed by the given gateway.
func NewShopService(shops driven.ShopGateway) *ShopService {
	return &ShopService{shops: shops}
}

// ListShops returns all shops connected to the account.
func (s *ShopService) ListShops(ctx context.Context) ([]domain.Shop, error) {
	shops, err := s.shops.ListShops(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing shops: %w", err)
	}
	return shops, nil
}

// DisconnectShop removes the sales-channel connection for a shop.
func (s *ShopService) DisconnectShop(ctx context.Context, id domain.ShopID) error {
	if !id.Valid() {
		return fmt.Errorf("shop id %d: %w", id, domain.ErrInvalidInput)
	}
	if err := s.shops.DisconnectShop(ctx, id); err != nil {
		return fmt.Errorf("disconnecting shop %d: %w", id, err)
	}
	return nil
}
