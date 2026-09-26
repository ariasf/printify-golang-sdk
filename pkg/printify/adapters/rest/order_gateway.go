package rest

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/printify-go/pkg/printify/domain"
	"github.com/printify-go/pkg/printify/ports/driving"
)

// OrderGateway implements driven.OrderGateway over the Printify REST API.
type OrderGateway struct {
	client *Client
}

// NewOrderGateway builds an OrderGateway using the shared transport client.
func NewOrderGateway(client *Client) *OrderGateway {
	return &OrderGateway{client: client}
}

// ---- DTOs ----

type addressDTO struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Region    string `json:"region"`
	Address1  string `json:"address1"`
	Address2  string `json:"address2"`
	City      string `json:"city"`
	Zip       string `json:"zip"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Country   string `json:"country"`
	Company   string `json:"company"`
}

type lineItemMetadataDTO struct {
	Title        string `json:"title"`
	Price        int    `json:"price"`
	VariantLabel string `json:"variant_label"`
	SKU          string `json:"sku"`
	Country      string `json:"country"`
	ExternalID   string `json:"external_id"`
}

type lineItemDTO struct {
	ProductID          string               `json:"product_id"`
	Quantity           int                  `json:"quantity"`
	VariantID          int64                `json:"variant_id"`
	PrintProviderID    int64                `json:"print_provider_id"`
	Cost               int                  `json:"cost"`
	ShippingCost       int                  `json:"shipping_cost"`
	Status             string               `json:"status"`
	Metadata           *lineItemMetadataDTO `json:"metadata"`
	SentToProductionAt string               `json:"sent_to_production_at"`
	FulfilledAt        string               `json:"fulfilled_at"`
}

type orderMetadataDTO struct {
	OrderType              string   `json:"order_type"`
	ShopOrderID            int64    `json:"shop_order_id"`
	ShopOrderLabel         string   `json:"shop_order_label"`
	ShopFulfilledAt        string   `json:"shop_fulfilled_at"`
	IsReprint              bool     `json:"is_reprint"`
	ReprintedOrderIDs      []string `json:"reprinted_order_ids"`
	ChildReprintedOrderIDs []string `json:"child_reprinted_order_ids"`
}

type shipmentDTO struct {
	Carrier     string `json:"carrier"`
	Number      string `json:"number"`
	URL         string `json:"url"`
	DeliveredAt string `json:"delivered_at"`
}

type orderDTO struct {
	ID                 string            `json:"id"`
	AppOrderID         string            `json:"app_order_id"`
	AddressTo          addressDTO        `json:"address_to"`
	LineItems          []lineItemDTO     `json:"line_items"`
	Metadata           *orderMetadataDTO `json:"metadata"`
	TotalPrice         int               `json:"total_price"`
	TotalShipping      int               `json:"total_shipping"`
	TotalTax           int               `json:"total_tax"`
	Status             string            `json:"status"`
	ShippingMethod     int               `json:"shipping_method"`
	IsPrintifyExpress  bool              `json:"is_printify_express"`
	IsEconomyShipping  bool              `json:"is_economy_shipping"`
	Shipments          []shipmentDTO     `json:"shipments"`
	CreatedAt          string            `json:"created_at"`
	SentToProductionAt string            `json:"sent_to_production_at"`
	FulfilledAt        string            `json:"fulfilled_at"`
}

type orderCreatedDTO struct {
	ID string `json:"id"`
}

type shippingCostsDTO struct {
	Standard        *int `json:"standard"`
	Express         *int `json:"express"`
	Priority        *int `json:"priority"`
	PrintifyExpress *int `json:"printify_express"`
	Economy         *int `json:"economy"`
}

type addressChangeResultDTO struct {
	Resolution string      `json:"resolution"`
	OrderID    string      `json:"order_id"`
	AddressTo  *addressDTO `json:"address_to"`
	ID         string      `json:"id"`
	Status     string      `json:"status"`
	Outcome    string      `json:"outcome"`
}

// ---- mapping ----

func parseTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	if t, err := time.Parse("2006-01-02 15:04:05-07:00", s); err == nil {
		return t
	}
	return time.Time{}
}

func (d addressDTO) toDomain() domain.Address {
	return domain.Address{
		FirstName: d.FirstName, LastName: d.LastName, Region: d.Region,
		Address1: d.Address1, Address2: d.Address2, City: d.City, Zip: d.Zip,
		Email: d.Email, Phone: d.Phone, Country: d.Country, Company: d.Company,
	}
}

func (d orderDTO) toDomain() domain.Order {
	o := domain.Order{
		ID: domain.OrderID(d.ID), AppOrderID: d.AppOrderID,
		AddressTo:  d.AddressTo.toDomain(),
		TotalPrice: d.TotalPrice, TotalShipping: d.TotalShipping, TotalTax: d.TotalTax,
		Status: d.Status, ShippingMethod: d.ShippingMethod,
		IsPrintifyExpress: d.IsPrintifyExpress, IsEconomyShipping: d.IsEconomyShipping,
		CreatedAt: parseTime(d.CreatedAt),
	}
	for _, li := range d.LineItems {
		item := domain.LineItem{
			ProductID: domain.ProductID(li.ProductID), Quantity: li.Quantity,
			VariantID: domain.VariantID(li.VariantID), PrintProviderID: domain.PrintProviderID(li.PrintProviderID),
			Cost: li.Cost, ShippingCost: li.ShippingCost, Status: li.Status,
		}
		if li.Metadata != nil {
			item.Metadata = &domain.LineItemMetadata{
				Title: li.Metadata.Title, Price: li.Metadata.Price,
				VariantLabel: li.Metadata.VariantLabel, SKU: li.Metadata.SKU,
				Country: li.Metadata.Country, ExternalID: li.Metadata.ExternalID,
			}
		}
		if t := parseTime(li.SentToProductionAt); !t.IsZero() {
			item.SentToProductionAt = &t
		}
		if t := parseTime(li.FulfilledAt); !t.IsZero() {
			item.FulfilledAt = &t
		}
		o.LineItems = append(o.LineItems, item)
	}
	if d.Metadata != nil {
		md := &domain.OrderMetadata{
			OrderType: d.Metadata.OrderType, ShopOrderID: int(d.Metadata.ShopOrderID),
			ShopOrderLabel:         d.Metadata.ShopOrderLabel,
			IsReprint:              d.Metadata.IsReprint,
			ReprintedOrderIDs:      d.Metadata.ReprintedOrderIDs,
			ChildReprintedOrderIDs: d.Metadata.ChildReprintedOrderIDs,
		}
		if t := parseTime(d.Metadata.ShopFulfilledAt); !t.IsZero() {
			md.ShopFulfilledAt = &t
		}
		o.Metadata = md
	}
	for _, s := range d.Shipments {
		sh := domain.Shipment{Carrier: s.Carrier, Number: s.Number, URL: s.URL}
		if t := parseTime(s.DeliveredAt); !t.IsZero() {
			sh.DeliveredAt = &t
		}
		o.Shipments = append(o.Shipments, sh)
	}
	if t := parseTime(d.SentToProductionAt); !t.IsZero() {
		o.SentToProductionAt = &t
	}
	if t := parseTime(d.FulfilledAt); !t.IsZero() {
		o.FulfilledAt = &t
	}
	return o
}

func toAddressDTO(a domain.Address) addressDTO {
	return addressDTO{
		FirstName: a.FirstName, LastName: a.LastName, Region: a.Region,
		Address1: a.Address1, Address2: a.Address2, City: a.City, Zip: a.Zip,
		Email: a.Email, Phone: a.Phone, Country: a.Country, Company: a.Company,
	}
}

func toLineItemDTOList(items []domain.LineItemBlueprint) []map[string]any {
	out := make([]map[string]any, 0, len(items))
	for _, li := range items {
		m := map[string]any{
			"print_provider_id": int64(li.PrintProviderID),
			"blueprint_id":      int64(li.BlueprintID),
			"variant_id":        int64(li.VariantID),
			"quantity":          li.Quantity,
		}
		if li.PrintAreas != nil {
			m["print_areas"] = li.PrintAreas
		}
		if li.ExternalID != "" {
			m["external_id"] = li.ExternalID
		}
		if li.Personalisation != nil {
			m["personalisation"] = map[string]any{
				"personalisation_strategy":     li.Personalisation.Strategy,
				"personalisation_instructions": li.Personalisation.Instructions,
				"buyer_request":                li.Personalisation.BuyerRequest,
			}
		}
		out = append(out, m)
	}
	return out
}

// ---- gateway ----

func orderListQuery(f driving.OrderFilter) url.Values {
	q := url.Values{}
	if f.Limit != nil {
		q.Set("limit", fmt.Sprint(*f.Limit))
	}
	if f.Page != nil {
		q.Set("page", fmt.Sprint(*f.Page))
	}
	if f.Status != "" {
		q.Set("status", f.Status)
	}
	if f.SKU != "" {
		q.Set("sku", f.SKU)
	}
	return q
}

// List fetches a page of orders for a shop.
func (g *OrderGateway) List(ctx context.Context, shop domain.ShopID, f driving.OrderFilter) ([]domain.Order, error) {
	var dtos []orderDTO
	path := fmt.Sprintf("/v1/shops/%d/orders.json", shop)
	if err := g.client.do(ctx, http.MethodGet, path, orderListQuery(f), nil, &dtos); err != nil {
		return nil, err
	}
	out := make([]domain.Order, 0, len(dtos))
	for _, d := range dtos {
		out = append(out, d.toDomain())
	}
	return out, nil
}

// Get fetches one order.
func (g *OrderGateway) Get(ctx context.Context, shop domain.ShopID, id domain.OrderID) (*domain.Order, error) {
	var dto orderDTO
	path := fmt.Sprintf("/v1/shops/%d/orders/%s.json", shop, id)
	if err := g.client.do(ctx, http.MethodGet, path, nil, nil, &dto); err != nil {
		return nil, err
	}
	o := dto.toDomain()
	return &o, nil
}

// Submit creates an order.
func (g *OrderGateway) Submit(ctx context.Context, shop domain.ShopID, o domain.SubmitOrder) (*domain.OrderIDResult, error) {
	body := map[string]any{
		"line_items": toLineItemDTOList(o.LineItems),
		"address_to": toAddressDTO(o.AddressTo),
	}
	setNullable(body, "external_id", o.ExternalID)
	setNullable(body, "label", o.Label)
	if o.ShippingMethod != nil {
		body["shipping_method"] = *o.ShippingMethod
	}
	if o.IsPrintifyExpress != nil {
		body["is_printify_express"] = *o.IsPrintifyExpress
	}
	if o.IsEconomyShipping != nil {
		body["is_economy_shipping"] = *o.IsEconomyShipping
	}
	if o.SendShippingNotification != nil {
		body["send_shipping_notification"] = *o.SendShippingNotification
	}
	var dto orderCreatedDTO
	path := fmt.Sprintf("/v1/shops/%d/orders.json", shop)
	if err := g.client.do(ctx, http.MethodPost, path, nil, body, &dto); err != nil {
		return nil, err
	}
	return &domain.OrderIDResult{ID: domain.OrderID(dto.ID)}, nil
}

// SubmitExpress creates a Printify Express order.
func (g *OrderGateway) SubmitExpress(ctx context.Context, shop domain.ShopID, o domain.ExpressOrder) ([]map[string]any, error) {
	body := map[string]any{
		"line_items": toLineItemDTOList(o.LineItems),
		"address_to": toAddressDTO(o.AddressTo),
	}
	setNullable(body, "external_id", o.ExternalID)
	setNullable(body, "label", o.Label)
	if o.ShippingMethod != nil {
		body["shipping_method"] = *o.ShippingMethod
	}
	if o.SendShippingNotification != nil {
		body["send_shipping_notification"] = *o.SendShippingNotification
	}
	var resp struct {
		Data []map[string]any `json:"data"`
	}
	path := fmt.Sprintf("/v1/shops/%d/orders/express.json", shop)
	if err := g.client.do(ctx, http.MethodPost, path, nil, body, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Cancel cancels an order.
func (g *OrderGateway) Cancel(ctx context.Context, shop domain.ShopID, id domain.OrderID) (*domain.Order, error) {
	var dto orderDTO
	path := fmt.Sprintf("/v1/shops/%d/orders/%s/cancel.json", shop, id)
	if err := g.client.do(ctx, http.MethodPost, path, nil, nil, &dto); err != nil {
		return nil, err
	}
	o := dto.toDomain()
	return &o, nil
}

// SendToProduction sends an order to production.
func (g *OrderGateway) SendToProduction(ctx context.Context, shop domain.ShopID, id domain.OrderID) (*domain.OrderIDResult, error) {
	var dto orderCreatedDTO
	path := fmt.Sprintf("/v1/shops/%d/orders/%s/send_to_production.json", shop, id)
	if err := g.client.do(ctx, http.MethodPost, path, nil, nil, &dto); err != nil {
		return nil, err
	}
	return &domain.OrderIDResult{ID: domain.OrderID(dto.ID)}, nil
}

// CalculateShipping returns shipping costs.
func (g *OrderGateway) CalculateShipping(ctx context.Context, shop domain.ShopID, o domain.SubmitOrder) (*domain.ShippingCosts, error) {
	body := map[string]any{
		"line_items": toLineItemDTOList(o.LineItems),
		"address_to": toAddressDTO(o.AddressTo),
	}
	var dto shippingCostsDTO
	path := fmt.Sprintf("/v1/shops/%d/orders/shipping.json", shop)
	if err := g.client.do(ctx, http.MethodPost, path, nil, body, &dto); err != nil {
		return nil, err
	}
	return &domain.ShippingCosts{
		Standard: dto.Standard, Express: dto.Express, Priority: dto.Priority,
		PrintifyExpress: dto.PrintifyExpress, Economy: dto.Economy,
	}, nil
}

// ChangeAddress attempts to change an order's shipping address.
func (g *OrderGateway) ChangeAddress(ctx context.Context, shop domain.ShopID, id domain.OrderID, a domain.AddressChange) (*domain.AddressChangeResult, error) {
	body := map[string]any{
		"first_name": a.FirstName, "last_name": a.LastName, "address1": a.Address1,
		"city": a.City, "country": a.Country, "zip": a.Zip,
	}
	if a.Address2 != "" {
		body["address2"] = a.Address2
	}
	if a.Region != "" {
		body["region"] = a.Region
	}
	if a.Email != "" {
		body["email"] = a.Email
	}
	if a.Phone != "" {
		body["phone"] = a.Phone
	}
	var dto addressChangeResultDTO
	path := fmt.Sprintf("/v1/shops/%d/orders/%s/support-requests/address-change", shop, id)
	if err := g.client.do(ctx, http.MethodPost, path, nil, body, &dto); err != nil {
		return nil, err
	}
	result := &domain.AddressChangeResult{Resolution: dto.Resolution, OrderID: domain.OrderID(dto.OrderID)}
	if dto.AddressTo != nil {
		addr := dto.AddressTo.toDomain()
		result.AddressTo = &addr
	}
	if dto.Resolution == "support_request_created" {
		result.SupportRequest = &domain.SupportRequest{
			ID: domain.SupportRequestID(dto.ID), Type: "order_address_editing",
			Status: dto.Status, Outcome: dto.Outcome,
		}
	}
	return result, nil
}

func setNullable(m map[string]any, key, val string) {
	if val != "" {
		m[key] = val
	}
}
