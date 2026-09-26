// Package driven defines outbound ports: interfaces the services need
// from the outside world. Interfaces only — no logic.
package driven

import (
	"context"

	"github.com/printify-go/pkg/printify/domain"
)

// ShopGateway abstracts the Printify shops API.
type ShopGateway interface {
	ListShops(ctx context.Context) ([]domain.Shop, error)
	DisconnectShop(ctx context.Context, id domain.ShopID) error
}
