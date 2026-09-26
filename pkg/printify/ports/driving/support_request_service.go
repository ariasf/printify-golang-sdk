package driving

import (
	"context"

	"github.com/printify-go/pkg/printify/domain"
)

// SupportRequestService exposes order support requests to SDK consumers.
type SupportRequestService interface {
	List(ctx context.Context, shop domain.ShopID, order domain.OrderID) ([]domain.SupportRequest, error)
	Get(ctx context.Context, shop domain.ShopID, order domain.OrderID, id domain.SupportRequestID) (*domain.SupportRequest, error)
	RequestReprint(ctx context.Context, shop domain.ShopID, order domain.OrderID, r domain.ReprintRequest) (*domain.SupportRequest, error)
	RequestRefund(ctx context.Context, shop domain.ShopID, order domain.OrderID, r domain.RefundRequest) (*domain.SupportRequest, error)
}
