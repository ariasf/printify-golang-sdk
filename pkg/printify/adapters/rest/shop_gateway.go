package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/printify-go/pkg/printify/domain"
)

// ShopGateway implements driven.ShopGateway over the Printify REST API.
type ShopGateway struct {
	client *Client
}

// NewShopGateway builds a ShopGateway using the shared transport client.
func NewShopGateway(client *Client) *ShopGateway {
	return &ShopGateway{client: client}
}

// shopDTO mirrors the wire format; transport tags stay out of the domain.
type shopDTO struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Channel string `json:"sales_channel"`
}

func (d shopDTO) toDomain() domain.Shop {
	return domain.Shop{ID: domain.ShopID(d.ID), Title: d.Title, Channel: d.Channel}
}

// ListShops fetches all shops for the account.
func (g *ShopGateway) ListShops(ctx context.Context) ([]domain.Shop, error) {
	var dtos []shopDTO
	if err := g.client.do(ctx, http.MethodGet, "/v1/shops.json", nil, nil, &dtos); err != nil {
		return nil, err
	}
	shops := make([]domain.Shop, 0, len(dtos))
	for _, d := range dtos {
		shops = append(shops, d.toDomain())
	}
	return shops, nil
}

// DisconnectShop removes the shop's sales-channel connection.
func (g *ShopGateway) DisconnectShop(ctx context.Context, id domain.ShopID) error {
	path := fmt.Sprintf("/v1/shops/%d/connection.json", id)
	return g.client.do(ctx, http.MethodDelete, path, nil, nil, nil)
}
