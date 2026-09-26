// Package driving defines inbound ports: the interfaces the SDK exposes
// to consumers. Interfaces only — no logic.
package driving

import (
	"context"

	"github.com/printify-go/pkg/printify/domain"
)

// ShopService exposes shop operations to SDK consumers.
type ShopService interface {
	ListShops(ctx context.Context) ([]domain.Shop, error)
	DisconnectShop(ctx context.Context, id domain.ShopID) error
}
