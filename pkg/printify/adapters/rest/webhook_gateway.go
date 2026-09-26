package rest

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/printify-go/pkg/printify/domain"
	"github.com/printify-go/pkg/printify/ports/driving"
)

// WebhookGateway implements driven.WebhookGateway over the Printify REST API.
type WebhookGateway struct {
	client *Client
}

// NewWebhookGateway builds a WebhookGateway using the shared transport client.
func NewWebhookGateway(client *Client) *WebhookGateway {
	return &WebhookGateway{client: client}
}

type webhookDTO struct {
	ID     int64  `json:"id"`
	Topic  string `json:"topic"`
	URL    string `json:"url"`
	ShopID string `json:"shop_id"`
}

func webhookListQuery(f driving.WebhookFilter) url.Values {
	q := url.Values{}
	if f.Limit != nil {
		q.Set("limit", fmt.Sprint(*f.Limit))
	}
	if f.Page != nil {
		q.Set("page", fmt.Sprint(*f.Page))
	}
	return q
}

// List fetches a page of webhooks for a shop.
func (g *WebhookGateway) List(ctx context.Context, shop domain.ShopID, f driving.WebhookFilter) ([]domain.Webhook, error) {
	var dtos []webhookDTO
	path := fmt.Sprintf("/v1/shops/%d/webhooks.json", shop)
	if err := g.client.do(ctx, http.MethodGet, path, webhookListQuery(f), nil, &dtos); err != nil {
		return nil, err
	}
	out := make([]domain.Webhook, 0, len(dtos))
	for _, d := range dtos {
		out = append(out, domain.Webhook{ID: domain.WebhookID(fmt.Sprint(d.ID)), Topic: d.Topic, URL: d.URL})
	}
	return out, nil
}

// Create subscribes a URL to a topic.
func (g *WebhookGateway) Create(ctx context.Context, shop domain.ShopID, w domain.CreateWebhook) (*domain.Webhook, error) {
	body := map[string]any{"topic": w.Topic, "url": w.URL}
	var dto webhookDTO
	path := fmt.Sprintf("/v1/shops/%d/webhooks.json", shop)
	if err := g.client.do(ctx, http.MethodPost, path, nil, body, &dto); err != nil {
		return nil, err
	}
	wh := domain.Webhook{ID: domain.WebhookID(fmt.Sprint(dto.ID)), Topic: dto.Topic, URL: dto.URL, ShopID: shop}
	return &wh, nil
}

// Modify updates the delivery URL of a webhook.
func (g *WebhookGateway) Modify(ctx context.Context, shop domain.ShopID, id domain.WebhookID, url string) (*domain.Webhook, error) {
	body := map[string]any{"url": url}
	var dto webhookDTO
	path := fmt.Sprintf("/v1/shops/%d/webhooks/%s.json", shop, id)
	if err := g.client.do(ctx, http.MethodPut, path, nil, body, &dto); err != nil {
		return nil, err
	}
	wh := domain.Webhook{ID: domain.WebhookID(fmt.Sprint(dto.ID)), Topic: dto.Topic, URL: dto.URL, ShopID: shop}
	return &wh, nil
}

// Delete removes a webhook.
func (g *WebhookGateway) Delete(ctx context.Context, shop domain.ShopID, id domain.WebhookID, host string) (domain.WebhookID, error) {
	var dto webhookDTO
	path := fmt.Sprintf("/v1/shops/%d/webhooks/%s.json", shop, id)
	q := url.Values{"host": {host}}
	if err := g.client.do(ctx, http.MethodDelete, path, q, nil, &dto); err != nil {
		return "", err
	}
	return domain.WebhookID(fmt.Sprint(dto.ID)), nil
}

// Simulate triggers a test delivery of a webhook.
func (g *WebhookGateway) Simulate(ctx context.Context, shop domain.ShopID, id domain.WebhookID) error {
	path := fmt.Sprintf("/v1/shops/%d/webhooks/%s/simulate", shop, id)
	return g.client.do(ctx, http.MethodPost, path, nil, nil, nil)
}
