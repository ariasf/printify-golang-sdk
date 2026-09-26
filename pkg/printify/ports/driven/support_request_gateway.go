package driven

import (
	"context"

	"github.com/printify-go/pkg/printify/domain"
)

// SupportRequestGateway abstracts the Printify order support requests API.
type SupportRequestGateway interface {
	List(ctx context.Context, shop domain.ShopID, order domain.OrderID) ([]domain.SupportRequest, error)
	Get(ctx context.Context, shop domain.ShopID, order domain.OrderID, id domain.SupportRequestID) (*domain.SupportRequest, error)
	Reprint(ctx context.Context, shop domain.ShopID, order domain.OrderID, r domain.ReprintRequest) (*domain.SupportRequest, error)
	Refund(ctx context.Context, shop domain.ShopID, order domain.OrderID, r domain.RefundRequest) (*domain.SupportRequest, error)
}
