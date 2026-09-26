package driven

import (
	"context"

	"github.com/printify-go/pkg/printify/domain"
	"github.com/printify-go/pkg/printify/ports/driving"
)

// OrderGateway abstracts the Printify orders API.
type OrderGateway interface {
	List(ctx context.Context, shop domain.ShopID, f driving.OrderFilter) ([]domain.Order, error)
	Get(ctx context.Context, shop domain.ShopID, id domain.OrderID) (*domain.Order, error)
	Submit(ctx context.Context, shop domain.ShopID, o domain.SubmitOrder) (*domain.OrderIDResult, error)
	SubmitExpress(ctx context.Context, shop domain.ShopID, o domain.ExpressOrder) ([]map[string]any, error)
	Cancel(ctx context.Context, shop domain.ShopID, id domain.OrderID) (*domain.Order, error)
	SendToProduction(ctx context.Context, shop domain.ShopID, id domain.OrderID) (*domain.OrderIDResult, error)
	CalculateShipping(ctx context.Context, shop domain.ShopID, o domain.SubmitOrder) (*domain.ShippingCosts, error)
	ChangeAddress(ctx context.Context, shop domain.ShopID, id domain.OrderID, a domain.AddressChange) (*domain.AddressChangeResult, error)
}
