package driving

import (
	"context"

	"github.com/printify-go/pkg/printify/domain"
)

// WebhookFilter narrows the webhook list.
type WebhookFilter struct {
	Limit *int
	Page  *int
}

// WebhookService exposes webhook management to SDK consumers.
type WebhookService interface {
	List(ctx context.Context, shop domain.ShopID, f WebhookFilter) ([]domain.Webhook, error)
	Create(ctx context.Context, shop domain.ShopID, topic domain.WebhookTopic, url string) (*domain.Webhook, error)
	Modify(ctx context.Context, shop domain.ShopID, id domain.WebhookID, url string) (*domain.Webhook, error)
	Delete(ctx context.Context, shop domain.ShopID, id domain.WebhookID, host string) (domain.WebhookID, error)
	Simulate(ctx context.Context, shop domain.ShopID, id domain.WebhookID) error
}
