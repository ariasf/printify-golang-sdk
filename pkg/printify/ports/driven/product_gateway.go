package driven

import (
	"context"

	"github.com/printify-go/pkg/printify/domain"
	"github.com/printify-go/pkg/printify/ports/driving"
)

// ProductGateway abstracts the Printify products API.
type ProductGateway interface {
	List(ctx context.Context, shop domain.ShopID, f driving.ProductFilter) ([]domain.Product, error)
	Get(ctx context.Context, shop domain.ShopID, id domain.ProductID) (*domain.Product, error)
	Create(ctx context.Context, shop domain.ShopID, p domain.CreateProduct) (*domain.Product, error)
	Update(ctx context.Context, shop domain.ShopID, id domain.ProductID, p domain.UpdateProduct) (*domain.Product, error)
	Delete(ctx context.Context, shop domain.ShopID, id domain.ProductID) error
	Publish(ctx context.Context, shop domain.ShopID, id domain.ProductID, opts domain.PublishProduct) error
	Unpublish(ctx context.Context, shop domain.ShopID, id domain.ProductID) error
	SetPublishSucceeded(ctx context.Context, shop domain.ShopID, id domain.ProductID, ext domain.PublishSucceeded) error
	SetPublishFailed(ctx context.Context, shop domain.ShopID, id domain.ProductID, reason domain.PublishFailed) error
	ListGpsr(ctx context.Context, shop domain.ShopID, id domain.ProductID) ([]domain.GpsrInfo, error)
}
