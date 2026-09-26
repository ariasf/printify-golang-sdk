package services

import (
	"context"
	"fmt"

	"github.com/printify-go/pkg/printify/domain"
	"github.com/printify-go/pkg/printify/ports/driven"
	"github.com/printify-go/pkg/printify/ports/driving"
)

// WebhookService implements driving.WebhookService.
type WebhookService struct {
	webhooks driven.WebhookGateway
}

// NewWebhookService builds a WebhookService backed by the given gateway.
func NewWebhookService(webhooks driven.WebhookGateway) *WebhookService {
	return &WebhookService{webhooks: webhooks}
}

func validWebhook(id domain.WebhookID) error {
	if !id.Valid() {
		return fmt.Errorf("webhook id %q: %w", id, domain.ErrInvalidInput)
	}
	return nil
}

// List returns a page of webhooks for a shop.
func (s *WebhookService) List(ctx context.Context, shop domain.ShopID, f driving.WebhookFilter) ([]domain.Webhook, error) {
	if err := validShop(shop); err != nil {
		return nil, err
	}
	webhooks, err := s.webhooks.List(ctx, shop, f)
	if err != nil {
		return nil, fmt.Errorf("listing webhooks in shop %d: %w", shop, err)
	}
	return webhooks, nil
}

// Create subscribes a URL to a topic.
func (s *WebhookService) Create(ctx context.Context, shop domain.ShopID, topic domain.WebhookTopic, url string) (*domain.Webhook, error) {
	if err := validShop(shop); err != nil {
		return nil, err
	}
	if err := validateWebhookCreate(topic, url); err != nil {
		return nil, err
	}
	w, err := s.webhooks.Create(ctx, shop, domain.CreateWebhook{Topic: string(topic), URL: url})
	if err != nil {
		return nil, fmt.Errorf("creating webhook in shop %d: %w", shop, err)
	}
	return w, nil
}

// Modify updates the delivery URL of a webhook.
func (s *WebhookService) Modify(ctx context.Context, shop domain.ShopID, id domain.WebhookID, url string) (*domain.Webhook, error) {
	if err := validShop(shop); err != nil {
		return nil, err
	}
	if err := validWebhook(id); err != nil {
		return nil, err
	}
	if url == "" {
		return nil, fmt.Errorf("url: %w", domain.ErrInvalidInput)
	}
	w, err := s.webhooks.Modify(ctx, shop, id, url)
	if err != nil {
		return nil, fmt.Errorf("modifying webhook %s in shop %d: %w", id, shop, err)
	}
	return w, nil
}

// Delete removes a webhook. Host is a safeguard that must match the webhook's host.
func (s *WebhookService) Delete(ctx context.Context, shop domain.ShopID, id domain.WebhookID, host string) (domain.WebhookID, error) {
	if err := validShop(shop); err != nil {
		return "", err
	}
	if err := validWebhook(id); err != nil {
		return "", err
	}
	if host == "" {
		return "", fmt.Errorf("host: %w", domain.ErrInvalidInput)
	}
	deleted, err := s.webhooks.Delete(ctx, shop, id, host)
	if err != nil {
		return "", fmt.Errorf("deleting webhook %s in shop %d: %w", id, shop, err)
	}
	return deleted, nil
}

// Simulate triggers a test delivery of a webhook.
func (s *WebhookService) Simulate(ctx context.Context, shop domain.ShopID, id domain.WebhookID) error {
	if err := validShop(shop); err != nil {
		return err
	}
	if err := validWebhook(id); err != nil {
		return err
	}
	if err := s.webhooks.Simulate(ctx, shop, id); err != nil {
		return fmt.Errorf("simulating webhook %s in shop %d: %w", id, shop, err)
	}
	return nil
}

func validateWebhookCreate(topic domain.WebhookTopic, url string) error {
	if topic == "" {
		return fmt.Errorf("topic: %w", domain.ErrInvalidInput)
	}
	if url == "" {
		return fmt.Errorf("url: %w", domain.ErrInvalidInput)
	}
	return nil
}
