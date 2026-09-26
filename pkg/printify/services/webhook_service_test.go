package services

import (
	"context"
	"errors"
	"testing"

	"github.com/printify-go/pkg/printify/domain"
	"github.com/printify-go/pkg/printify/ports/driving"
)

type stubWebhookGateway struct {
	listFn     func(ctx context.Context, shop domain.ShopID, f driving.WebhookFilter) ([]domain.Webhook, error)
	createFn   func(ctx context.Context, shop domain.ShopID, w domain.CreateWebhook) (*domain.Webhook, error)
	modifyFn   func(ctx context.Context, shop domain.ShopID, id domain.WebhookID, url string) (*domain.Webhook, error)
	deleteFn   func(ctx context.Context, shop domain.ShopID, id domain.WebhookID, host string) (domain.WebhookID, error)
	simulateFn func(ctx context.Context, shop domain.ShopID, id domain.WebhookID) error
}

func (s stubWebhookGateway) List(ctx context.Context, shop domain.ShopID, f driving.WebhookFilter) ([]domain.Webhook, error) {
	return s.listFn(ctx, shop, f)
}

func (s stubWebhookGateway) Create(ctx context.Context, shop domain.ShopID, w domain.CreateWebhook) (*domain.Webhook, error) {
	return s.createFn(ctx, shop, w)
}

func (s stubWebhookGateway) Modify(ctx context.Context, shop domain.ShopID, id domain.WebhookID, url string) (*domain.Webhook, error) {
	return s.modifyFn(ctx, shop, id, url)
}

func (s stubWebhookGateway) Delete(ctx context.Context, shop domain.ShopID, id domain.WebhookID, host string) (domain.WebhookID, error) {
	return s.deleteFn(ctx, shop, id, host)
}

func (s stubWebhookGateway) Simulate(ctx context.Context, shop domain.ShopID, id domain.WebhookID) error {
	return s.simulateFn(ctx, shop, id)
}

func TestWebhookService_List(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		shop    domain.ShopID
		gateway stubWebhookGateway
		want    int
		wantErr error
	}{
		{
			name: "rejects non-positive shop",
			shop: 0,
			gateway: stubWebhookGateway{listFn: func(context.Context, domain.ShopID, driving.WebhookFilter) ([]domain.Webhook, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "returns webhooks for valid shop",
			shop: 1,
			gateway: stubWebhookGateway{listFn: func(_ context.Context, shop domain.ShopID, _ driving.WebhookFilter) ([]domain.Webhook, error) {
				if shop != 1 {
					t.Fatalf("gateway got shop %d, want 1", shop)
				}
				return []domain.Webhook{{ID: "wh-1", Topic: "order.created"}}, nil
			}},
			want: 1,
		},
		{
			name: "propagates gateway error",
			shop: 1,
			gateway: stubWebhookGateway{listFn: func(context.Context, domain.ShopID, driving.WebhookFilter) ([]domain.Webhook, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewWebhookService(tt.gateway)
			got, err := svc.List(context.Background(), tt.shop, driving.WebhookFilter{})
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if len(got) != tt.want {
				t.Fatalf("len(webhooks) = %d, want %d", len(got), tt.want)
			}
		})
	}
}

func TestWebhookService_Create(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		shop    domain.ShopID
		topic   domain.WebhookTopic
		url     string
		gateway stubWebhookGateway
		wantErr error
	}{
		{
			name:  "rejects non-positive shop",
			shop:  0,
			topic: "order.created",
			url:   "https://example.com/hook",
			gateway: stubWebhookGateway{createFn: func(context.Context, domain.ShopID, domain.CreateWebhook) (*domain.Webhook, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:  "rejects empty topic",
			shop:  1,
			topic: "",
			url:   "https://example.com/hook",
			gateway: stubWebhookGateway{createFn: func(context.Context, domain.ShopID, domain.CreateWebhook) (*domain.Webhook, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:  "rejects empty url",
			shop:  1,
			topic: "order.created",
			url:   "",
			gateway: stubWebhookGateway{createFn: func(context.Context, domain.ShopID, domain.CreateWebhook) (*domain.Webhook, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:  "creates webhook for valid payload",
			shop:  1,
			topic: "order.created",
			url:   "https://example.com/hook",
			gateway: stubWebhookGateway{createFn: func(_ context.Context, shop domain.ShopID, w domain.CreateWebhook) (*domain.Webhook, error) {
				if shop != 1 || w.Topic != "order.created" || w.URL != "https://example.com/hook" {
					t.Fatalf("gateway got shop=%d webhook=%+v", shop, w)
				}
				return &domain.Webhook{ID: "wh-9", Topic: "order.created"}, nil
			}},
		},
		{
			name:  "propagates gateway error",
			shop:  1,
			topic: "order.created",
			url:   "https://example.com/hook",
			gateway: stubWebhookGateway{createFn: func(context.Context, domain.ShopID, domain.CreateWebhook) (*domain.Webhook, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewWebhookService(tt.gateway)
			got, err := svc.Create(context.Background(), tt.shop, tt.topic, tt.url)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got == nil {
				t.Fatal("webhook = nil, want non-nil")
			}
		})
	}
}

func TestWebhookService_Modify(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		shop    domain.ShopID
		id      domain.WebhookID
		url     string
		gateway stubWebhookGateway
		wantErr error
	}{
		{
			name: "rejects non-positive shop",
			shop: 0,
			id:   "wh-1",
			url:  "https://example.com/new",
			gateway: stubWebhookGateway{modifyFn: func(context.Context, domain.ShopID, domain.WebhookID, string) (*domain.Webhook, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects empty webhook id",
			shop: 1,
			id:   "",
			url:  "https://example.com/new",
			gateway: stubWebhookGateway{modifyFn: func(context.Context, domain.ShopID, domain.WebhookID, string) (*domain.Webhook, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects empty url",
			shop: 1,
			id:   "wh-1",
			url:  "",
			gateway: stubWebhookGateway{modifyFn: func(context.Context, domain.ShopID, domain.WebhookID, string) (*domain.Webhook, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "modifies webhook for valid payload",
			shop: 1,
			id:   "wh-1",
			url:  "https://example.com/new",
			gateway: stubWebhookGateway{modifyFn: func(_ context.Context, shop domain.ShopID, id domain.WebhookID, url string) (*domain.Webhook, error) {
				if shop != 1 || id != "wh-1" || url != "https://example.com/new" {
					t.Fatalf("gateway got shop=%d id=%q url=%q", shop, id, url)
				}
				return &domain.Webhook{ID: "wh-1", URL: url}, nil
			}},
		},
		{
			name: "propagates gateway error",
			shop: 1,
			id:   "wh-1",
			url:  "https://example.com/new",
			gateway: stubWebhookGateway{modifyFn: func(context.Context, domain.ShopID, domain.WebhookID, string) (*domain.Webhook, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewWebhookService(tt.gateway)
			got, err := svc.Modify(context.Background(), tt.shop, tt.id, tt.url)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got == nil {
				t.Fatal("webhook = nil, want non-nil")
			}
		})
	}
}

func TestWebhookService_Delete(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		shop    domain.ShopID
		id      domain.WebhookID
		host    string
		gateway stubWebhookGateway
		wantID  domain.WebhookID
		wantErr error
	}{
		{
			name: "rejects non-positive shop",
			shop: 0,
			id:   "wh-1",
			host: "example.com",
			gateway: stubWebhookGateway{deleteFn: func(context.Context, domain.ShopID, domain.WebhookID, string) (domain.WebhookID, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects empty webhook id",
			shop: 1,
			id:   "",
			host: "example.com",
			gateway: stubWebhookGateway{deleteFn: func(context.Context, domain.ShopID, domain.WebhookID, string) (domain.WebhookID, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects empty host",
			shop: 1,
			id:   "wh-1",
			host: "",
			gateway: stubWebhookGateway{deleteFn: func(context.Context, domain.ShopID, domain.WebhookID, string) (domain.WebhookID, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "deletes webhook for valid payload",
			shop: 1,
			id:   "wh-1",
			host: "example.com",
			gateway: stubWebhookGateway{deleteFn: func(_ context.Context, shop domain.ShopID, id domain.WebhookID, host string) (domain.WebhookID, error) {
				if shop != 1 || id != "wh-1" || host != "example.com" {
					t.Fatalf("gateway got shop=%d id=%q host=%q", shop, id, host)
				}
				return "wh-1", nil
			}},
			wantID: "wh-1",
		},
		{
			name: "propagates gateway error",
			shop: 1,
			id:   "wh-1",
			host: "example.com",
			gateway: stubWebhookGateway{deleteFn: func(context.Context, domain.ShopID, domain.WebhookID, string) (domain.WebhookID, error) {
				return "", gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewWebhookService(tt.gateway)
			gotID, err := svc.Delete(context.Background(), tt.shop, tt.id, tt.host)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && gotID != tt.wantID {
				t.Fatalf("id = %q, want %q", gotID, tt.wantID)
			}
		})
	}
}

func TestWebhookService_Simulate(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		shop    domain.ShopID
		id      domain.WebhookID
		gateway stubWebhookGateway
		wantErr error
	}{
		{
			name: "rejects non-positive shop",
			shop: 0,
			id:   "wh-1",
			gateway: stubWebhookGateway{simulateFn: func(context.Context, domain.ShopID, domain.WebhookID) error {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects empty webhook id",
			shop: 1,
			id:   "",
			gateway: stubWebhookGateway{simulateFn: func(context.Context, domain.ShopID, domain.WebhookID) error {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "simulates webhook for valid ids",
			shop: 1,
			id:   "wh-1",
			gateway: stubWebhookGateway{simulateFn: func(_ context.Context, shop domain.ShopID, id domain.WebhookID) error {
				if shop != 1 || id != "wh-1" {
					t.Fatalf("gateway got shop=%d id=%q, want 1 wh-1", shop, id)
				}
				return nil
			}},
		},
		{
			name: "propagates gateway error",
			shop: 1,
			id:   "wh-1",
			gateway: stubWebhookGateway{simulateFn: func(context.Context, domain.ShopID, domain.WebhookID) error {
				return gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewWebhookService(tt.gateway)
			err := svc.Simulate(context.Background(), tt.shop, tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
