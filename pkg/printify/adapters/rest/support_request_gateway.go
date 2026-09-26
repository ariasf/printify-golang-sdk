package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/printify-go/pkg/printify/domain"
)

// SupportRequestGateway implements driven.SupportRequestGateway over the Printify REST API.
type SupportRequestGateway struct {
	client *Client
}

// NewSupportRequestGateway builds a SupportRequestGateway using the shared transport client.
func NewSupportRequestGateway(client *Client) *SupportRequestGateway {
	return &SupportRequestGateway{client: client}
}

type supportRequestDTO struct {
	ID            string   `json:"id"`
	OrderID       string   `json:"order_id"`
	Type          string   `json:"type"`
	Status        string   `json:"status"`
	StatusDetails string   `json:"status_details"`
	Outcome       string   `json:"outcome"`
	LineItemIDs   []string `json:"line_item_ids"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
}

func (d supportRequestDTO) toDomain() domain.SupportRequest {
	return domain.SupportRequest{
		ID: domain.SupportRequestID(d.ID), OrderID: domain.OrderID(d.OrderID),
		Type: d.Type, Status: d.Status, StatusDetails: d.StatusDetails, Outcome: d.Outcome,
		LineItemIDs: d.LineItemIDs, CreatedAt: parseTime(d.CreatedAt), UpdatedAt: parseTime(d.UpdatedAt),
	}
}

func toSupportRequestDTOList(dtos []supportRequestDTO) []domain.SupportRequest {
	out := make([]domain.SupportRequest, 0, len(dtos))
	for _, d := range dtos {
		out = append(out, d.toDomain())
	}
	return out
}

func toSupportRequestLineItemDTO(items []domain.SupportRequestLineItem) []map[string]any {
	out := make([]map[string]any, 0, len(items))
	for _, li := range items {
		out = append(out, map[string]any{"line_item_id": li.LineItemID, "quantity": li.Quantity})
	}
	return out
}

// List fetches the support requests for an order.
func (g *SupportRequestGateway) List(ctx context.Context, shop domain.ShopID, order domain.OrderID) ([]domain.SupportRequest, error) {
	var dtos []supportRequestDTO
	path := fmt.Sprintf("/v1/shops/%d/orders/%s/support-requests", shop, order)
	if err := g.client.do(ctx, http.MethodGet, path, nil, nil, &dtos); err != nil {
		return nil, err
	}
	return toSupportRequestDTOList(dtos), nil
}

// Get fetches one support request.
func (g *SupportRequestGateway) Get(ctx context.Context, shop domain.ShopID, order domain.OrderID, id domain.SupportRequestID) (*domain.SupportRequest, error) {
	var dto supportRequestDTO
	path := fmt.Sprintf("/v1/shops/%d/orders/%s/support-requests/%s", shop, order, id)
	if err := g.client.do(ctx, http.MethodGet, path, nil, nil, &dto); err != nil {
		return nil, err
	}
	sr := dto.toDomain()
	return &sr, nil
}

// Reprint creates a reprint support request.
func (g *SupportRequestGateway) Reprint(ctx context.Context, shop domain.ShopID, order domain.OrderID, r domain.ReprintRequest) (*domain.SupportRequest, error) {
	body := map[string]any{
		"reason":      r.Reason,
		"description": r.Description,
		"line_items":  toSupportRequestLineItemDTO(r.LineItems),
	}
	if len(r.ImageURLs) > 0 {
		body["image_urls"] = r.ImageURLs
	}
	var dto supportRequestDTO
	path := fmt.Sprintf("/v1/shops/%d/orders/%s/support-requests/reprint", shop, order)
	if err := g.client.do(ctx, http.MethodPost, path, nil, body, &dto); err != nil {
		return nil, err
	}
	sr := dto.toDomain()
	return &sr, nil
}

// Refund creates a refund support request.
func (g *SupportRequestGateway) Refund(ctx context.Context, shop domain.ShopID, order domain.OrderID, r domain.RefundRequest) (*domain.SupportRequest, error) {
	body := map[string]any{
		"reason":      r.Reason,
		"description": r.Description,
		"line_items":  toSupportRequestLineItemDTO(r.LineItems),
	}
	if len(r.ImageURLs) > 0 {
		body["image_urls"] = r.ImageURLs
	}
	var dto supportRequestDTO
	path := fmt.Sprintf("/v1/shops/%d/orders/%s/support-requests/refund", shop, order)
	if err := g.client.do(ctx, http.MethodPost, path, nil, body, &dto); err != nil {
		return nil, err
	}
	sr := dto.toDomain()
	return &sr, nil
}
