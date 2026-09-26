package driven

import (
	"context"

	"github.com/printify-go/pkg/printify/domain"
	"github.com/printify-go/pkg/printify/ports/driving"
)

// WebhookGateway abstracts the Printify webhooks API.
type WebhookGateway interface {
	List(ctx context.Context, shop domain.ShopID, f driving.WebhookFilter) ([]domain.Webhook, error)
	Create(ctx context.Context, shop domain.ShopID, w domain.CreateWebhook) (*domain.Webhook, error)
	Modify(ctx context.Context, shop domain.ShopID, id domain.WebhookID, url string) (*domain.Webhook, error)
	Delete(ctx context.Context, shop domain.ShopID, id domain.WebhookID, host string) (domain.WebhookID, error)
	Simulate(ctx context.Context, shop domain.ShopID, id domain.WebhookID) error
}
