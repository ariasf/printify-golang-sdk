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

// ProductGateway implements driven.ProductGateway over the Printify REST API.
type ProductGateway struct {
	client *Client
}

// NewProductGateway builds a ProductGateway using the shared transport client.
func NewProductGateway(client *Client) *ProductGateway {
	return &ProductGateway{client: client}
}

// ---- DTOs ----

type optionValueDTO struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type productOptionDTO struct {
	Name   string           `json:"name"`
	Type   string           `json:"type"`
	Values []optionValueDTO `json:"values"`
}

type productVariantDTO struct {
	ID                        int64  `json:"id"`
	SKU                       string `json:"sku"`
	Cost                      int    `json:"cost"`
	Price                     int    `json:"price"`
	Title                     string `json:"title"`
	Grams                     int    `json:"grams"`
	IsEnabled                 bool   `json:"is_enabled"`
	IsDefault                 bool   `json:"is_default"`
	IsAvailable               bool   `json:"is_available"`
	IsPrintifyExpressEligible bool   `json:"is_printify_express_eligible"`
	Options                   []int  `json:"options"`
}

type productImageDTO struct {
	Src       string  `json:"src"`
	VariantID []int64 `json:"variant_ids"`
	Position  string  `json:"position"`
	IsDefault bool    `json:"is_default"`
}

type patternDTO struct {
	SpacingX int `json:"spacing_x"`
	SpacingY int `json:"spacing_y"`
	Scale    int `json:"scale"`
	Offset   int `json:"offset"`
}

type placeholderImageDTO struct {
	ID         string      `json:"id"`
	Src        *string     `json:"src"`
	Name       *string     `json:"name"`
	Type       string      `json:"type"`
	Height     int         `json:"height"`
	Width      int         `json:"width"`
	X          *float64    `json:"x"`
	Y          *float64    `json:"y"`
	Scale      *float64    `json:"scale"`
	Angle      int         `json:"angle"`
	FontFamily *string     `json:"font_family"`
	FontSize   *int        `json:"font_size"`
	FontWeight *int        `json:"font_weight"`
	FontColor  *string     `json:"font_color"`
	FontStyle  *string     `json:"font_style"`
	InputText  *string     `json:"input_text"`
	TextAlign  *string     `json:"text_align"`
	Pattern    *patternDTO `json:"pattern"`
}

type placeholderDTO struct {
	Position string                `json:"position"`
	Images   []placeholderImageDTO `json:"images"`
}

type printAreaDTO struct {
	VariantIDs   []int64          `json:"variant_ids"`
	Placeholders []placeholderDTO `json:"placeholders"`
	Background   string           `json:"background"`
}

type viewFileDTO struct {
	Src       string  `json:"src"`
	VariantID []int64 `json:"variant_ids"`
}

type viewDTO struct {
	ID       int           `json:"id"`
	Label    string        `json:"label"`
	Position string        `json:"position"`
	Files    []viewFileDTO `json:"files"`
}

type productDTO struct {
	ID                        string              `json:"id"`
	Title                     string              `json:"title"`
	Description               *string             `json:"description"`
	SafetyInformation         *string             `json:"safety_information"`
	Tags                      []string            `json:"tags"`
	Options                   []productOptionDTO  `json:"options"`
	Variants                  []productVariantDTO `json:"variants"`
	Images                    []productImageDTO   `json:"images"`
	CreatedAt                 string              `json:"created_at"`
	UpdatedAt                 string              `json:"updated_at"`
	Visible                   bool                `json:"visible"`
	IsLocked                  bool                `json:"is_locked"`
	IsPrintifyExpressEligible bool                `json:"is_printify_express_eligible"`
	IsPrintifyExpressEnabled  bool                `json:"is_printify_express_enabled"`
	IsEconomyShippingEligible bool                `json:"is_economy_shipping_eligible"`
	IsEconomyShippingEnabled  bool                `json:"is_economy_shipping_enabled"`
	BlueprintID               int64               `json:"blueprint_id"`
	UserID                    int64               `json:"user_id"`
	ShopID                    int64               `json:"shop_id"`
	PrintProviderID           int64               `json:"print_provider_id"`
	PrintAreas                []printAreaDTO      `json:"print_areas"`
	Views                     []viewDTO           `json:"views"`
	SalesChannelProps         []any               `json:"sales_channel_properties"`
}

type gpsrDTO struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// ---- mapping ----

func (d productDTO) toDomain() domain.Product {
	p := domain.Product{
		ID:                domain.ProductID(d.ID),
		Title:             d.Title,
		Description:       d.Description,
		SafetyInformation: d.SafetyInformation,
		Tags:              d.Tags,
		Visible:           d.Visible,
		Locked:            d.IsLocked,
		ExpressEligible:   d.IsPrintifyExpressEligible,
		ExpressEnabled:    d.IsPrintifyExpressEnabled,
		EconomyEligible:   d.IsEconomyShippingEligible,
		EconomyEnabled:    d.IsEconomyShippingEnabled,
		BlueprintID:       domain.BlueprintID(d.BlueprintID),
		UserID:            int(d.UserID),
		ShopID:            domain.ShopID(d.ShopID),
		PrintProviderID:   domain.PrintProviderID(d.PrintProviderID),
		SalesChannelProps: d.SalesChannelProps,
	}
	if t, err := time.Parse("2006-01-02 15:04:05-07:00", d.CreatedAt); err == nil {
		p.CreatedAt = t
	}
	if t, err := time.Parse("2006-01-02 15:04:05-07:00", d.UpdatedAt); err == nil {
		p.UpdatedAt = t
	}
	for _, o := range d.Options {
		p.Options = append(p.Options, domain.ProductOption{Name: o.Name, Type: o.Type, Values: toOptionValues(o.Values)})
	}
	for _, v := range d.Variants {
		p.Variants = append(p.Variants, toProductVariant(v))
	}
	for _, i := range d.Images {
		p.Images = append(p.Images, domain.ProductImage{
			Src: i.Src, VariantID: toVariantIDs(i.VariantID), Position: i.Position, IsDefault: i.IsDefault,
		})
	}
	for _, pa := range d.PrintAreas {
		p.PrintAreas = append(p.PrintAreas, toPrintArea(pa))
	}
	for _, v := range d.Views {
		p.Views = append(p.Views, toView(v))
	}
	return p
}

func toOptionValues(dtos []optionValueDTO) []domain.OptionValue {
	out := make([]domain.OptionValue, 0, len(dtos))
	for _, d := range dtos {
		out = append(out, domain.OptionValue{ID: d.ID, Title: d.Title})
	}
	return out
}

func toProductVariant(d productVariantDTO) domain.ProductVariant {
	return domain.ProductVariant{
		ID: domain.VariantID(d.ID), SKU: d.SKU, Cost: d.Cost, Price: d.Price,
		Title: d.Title, Grams: d.Grams, Enabled: d.IsEnabled, IsDefault: d.IsDefault,
		Available: d.IsAvailable, ExpressEligible: d.IsPrintifyExpressEligible, Options: d.Options,
	}
}

func toPrintArea(d printAreaDTO) domain.PrintArea {
	pa := domain.PrintArea{VariantIDs: toVariantIDs(d.VariantIDs), Background: d.Background}
	for _, ph := range d.Placeholders {
		p := domain.Placeholder{Position: ph.Position}
		for _, img := range ph.Images {
			p.Images = append(p.Images, toPlaceholderImage(img))
		}
		pa.Placeholders = append(pa.Placeholders, p)
	}
	return pa
}

func toPlaceholderImage(d placeholderImageDTO) domain.PlaceholderImage {
	pi := domain.PlaceholderImage{
		ID: d.ID, Src: d.Src, Name: d.Name, Type: d.Type,
		Height: d.Height, Width: d.Width, X: d.X, Y: d.Y, Scale: d.Scale, Angle: d.Angle,
		FontFamily: d.FontFamily, FontSize: d.FontSize, FontWeight: d.FontWeight,
		FontColor: d.FontColor, FontStyle: d.FontStyle, InputText: d.InputText, TextAlign: d.TextAlign,
	}
	if d.Pattern != nil {
		pi.Pattern = &domain.Pattern{SpacingX: d.Pattern.SpacingX, SpacingY: d.Pattern.SpacingY, Scale: d.Pattern.Scale, Offset: d.Pattern.Offset}
	}
	return pi
}

func toView(d viewDTO) domain.View {
	v := domain.View{ID: d.ID, Label: d.Label, Position: d.Position}
	for _, f := range d.Files {
		v.Files = append(v.Files, domain.ViewFile{Src: f.Src, VariantID: toVariantIDs(f.VariantID)})
	}
	return v
}

func toVariantDTO(p domain.ProductVariant) productVariantDTO {
	return productVariantDTO{
		ID: int64(p.ID), SKU: p.SKU, Cost: p.Cost, Price: p.Price,
		Title: p.Title, Grams: p.Grams, IsEnabled: p.Enabled, IsDefault: p.IsDefault,
		IsAvailable: p.Available, IsPrintifyExpressEligible: p.ExpressEligible, Options: p.Options,
	}
}

func toPrintAreaDTO(pa domain.PrintArea) printAreaDTO {
	d := printAreaDTO{VariantIDs: toInt64VariantIDs(pa.VariantIDs), Background: pa.Background}
	for _, ph := range pa.Placeholders {
		p := placeholderDTO{Position: ph.Position}
		for _, img := range ph.Images {
			p.Images = append(p.Images, toPlaceholderImageDTO(img))
		}
		d.Placeholders = append(d.Placeholders, p)
	}
	return d
}

func toInt64VariantIDs(ids []domain.VariantID) []int64 {
	out := make([]int64, 0, len(ids))
	for _, id := range ids {
		out = append(out, int64(id))
	}
	return out
}

func toPlaceholderImageDTO(pi domain.PlaceholderImage) placeholderImageDTO {
	d := placeholderImageDTO{
		ID: pi.ID, Src: pi.Src, Name: pi.Name, Type: pi.Type,
		Height: pi.Height, Width: pi.Width, X: pi.X, Y: pi.Y, Scale: pi.Scale, Angle: pi.Angle,
		FontFamily: pi.FontFamily, FontSize: pi.FontSize, FontWeight: pi.FontWeight,
		FontColor: pi.FontColor, FontStyle: pi.FontStyle, InputText: pi.InputText, TextAlign: pi.TextAlign,
	}
	if pi.Pattern != nil {
		d.Pattern = &patternDTO{SpacingX: pi.Pattern.SpacingX, SpacingY: pi.Pattern.SpacingY, Scale: pi.Pattern.Scale, Offset: pi.Pattern.Offset}
	}
	return d
}

// ---- gateway ----

func listQuery(f driving.ProductFilter) url.Values {
	q := url.Values{}
	if f.Limit != nil {
		q.Set("limit", fmt.Sprint(*f.Limit))
	}
	if f.Page != nil {
		q.Set("page", fmt.Sprint(*f.Page))
	}
	return q
}

// List fetches a page of products for a shop.
func (g *ProductGateway) List(ctx context.Context, shop domain.ShopID, f driving.ProductFilter) ([]domain.Product, error) {
	var dtos []productDTO
	path := fmt.Sprintf("/v1/shops/%d/products.json", shop)
	if err := g.client.do(ctx, http.MethodGet, path, listQuery(f), nil, &dtos); err != nil {
		return nil, err
	}
	out := make([]domain.Product, 0, len(dtos))
	for _, d := range dtos {
		out = append(out, d.toDomain())
	}
	return out, nil
}

// Get fetches one product.
func (g *ProductGateway) Get(ctx context.Context, shop domain.ShopID, id domain.ProductID) (*domain.Product, error) {
	var dto productDTO
	path := fmt.Sprintf("/v1/shops/%d/products/%s.json", shop, id)
	if err := g.client.do(ctx, http.MethodGet, path, nil, nil, &dto); err != nil {
		return nil, err
	}
	p := dto.toDomain()
	return &p, nil
}

// Create creates a product.
func (g *ProductGateway) Create(ctx context.Context, shop domain.ShopID, p domain.CreateProduct) (*domain.Product, error) {
	body := map[string]any{
		"title":             p.Title,
		"blueprint_id":      int64(p.BlueprintID),
		"print_provider_id": int64(p.PrintProviderID),
		"variants":          toVariantDTOList(p.Variants),
		"print_areas":       toPrintAreaDTOList(p.PrintAreas),
	}
	if p.Description != nil {
		body["description"] = *p.Description
	}
	if p.SafetyInformation != nil {
		body["safety_information"] = *p.SafetyInformation
	}
	var dto productDTO
	path := fmt.Sprintf("/v1/shops/%d/products.json", shop)
	if err := g.client.do(ctx, http.MethodPost, path, nil, body, &dto); err != nil {
		return nil, err
	}
	pd := dto.toDomain()
	return &pd, nil
}

// Update updates a product.
func (g *ProductGateway) Update(ctx context.Context, shop domain.ShopID, id domain.ProductID, p domain.UpdateProduct) (*domain.Product, error) {
	body := map[string]any{
		"title":             p.Title,
		"blueprint_id":      int64(p.BlueprintID),
		"print_provider_id": int64(p.PrintProviderID),
		"variants":          toVariantDTOList(p.Variants),
		"print_areas":       toPrintAreaDTOList(p.PrintAreas),
	}
	if p.Description != nil {
		body["description"] = *p.Description
	}
	if p.SafetyInformation != nil {
		body["safety_information"] = *p.SafetyInformation
	}
	var dto productDTO
	path := fmt.Sprintf("/v1/shops/%d/products/%s.json", shop, id)
	if err := g.client.do(ctx, http.MethodPut, path, nil, body, &dto); err != nil {
		return nil, err
	}
	pd := dto.toDomain()
	return &pd, nil
}

// Delete removes a product.
func (g *ProductGateway) Delete(ctx context.Context, shop domain.ShopID, id domain.ProductID) error {
	path := fmt.Sprintf("/v1/shops/%d/products/%s.json", shop, id)
	return g.client.do(ctx, http.MethodDelete, path, nil, nil, nil)
}

// Publish pushes a product to its channel.
func (g *ProductGateway) Publish(ctx context.Context, shop domain.ShopID, id domain.ProductID, opts domain.PublishProduct) error {
	body := map[string]any{
		"title": opts.Title, "description": opts.Description, "images": opts.Images,
		"variants": opts.Variants, "tags": opts.Tags, "keyFeatures": opts.KeyFeatures,
		"shipping_template": opts.ShippingTemplate,
	}
	path := fmt.Sprintf("/v1/shops/%d/products/%s/publish.json", shop, id)
	return g.client.do(ctx, http.MethodPost, path, nil, body, nil)
}

// Unpublish notifies the API that a product was unpublished.
func (g *ProductGateway) Unpublish(ctx context.Context, shop domain.ShopID, id domain.ProductID) error {
	path := fmt.Sprintf("/v1/shops/%d/products/%s/unpublish.json", shop, id)
	return g.client.do(ctx, http.MethodPost, path, nil, nil, nil)
}

// SetPublishSucceeded records a successful publish.
func (g *ProductGateway) SetPublishSucceeded(ctx context.Context, shop domain.ShopID, id domain.ProductID, ext domain.PublishSucceeded) error {
	body := map[string]any{"external": map[string]any{"id": ext.ExternalID, "handle": ext.ExternalHandle}}
	path := fmt.Sprintf("/v1/shops/%d/products/%s/publishing_succeeded.json", shop, id)
	return g.client.do(ctx, http.MethodPost, path, nil, body, nil)
}

// SetPublishFailed records a failed publish.
func (g *ProductGateway) SetPublishFailed(ctx context.Context, shop domain.ShopID, id domain.ProductID, reason domain.PublishFailed) error {
	body := map[string]any{"reason": reason.Reason}
	path := fmt.Sprintf("/v1/shops/%d/products/%s/publishing_failed.json", shop, id)
	return g.client.do(ctx, http.MethodPost, path, nil, body, nil)
}

// ListGpsr returns GPSR info for a product.
func (g *ProductGateway) ListGpsr(ctx context.Context, shop domain.ShopID, id domain.ProductID) ([]domain.GpsrInfo, error) {
	var dtos []gpsrDTO
	path := fmt.Sprintf("/v1/shops/%d/products/%s/gpsr.json", shop, id)
	if err := g.client.do(ctx, http.MethodGet, path, nil, nil, &dtos); err != nil {
		return nil, err
	}
	out := make([]domain.GpsrInfo, 0, len(dtos))
	for _, d := range dtos {
		out = append(out, domain.GpsrInfo{Title: d.Title, Text: d.Text})
	}
	return out, nil
}

func toVariantDTOList(vs []domain.ProductVariant) []productVariantDTO {
	out := make([]productVariantDTO, 0, len(vs))
	for _, v := range vs {
		out = append(out, toVariantDTO(v))
	}
	return out
}

func toPrintAreaDTOList(pas []domain.PrintArea) []printAreaDTO {
	out := make([]printAreaDTO, 0, len(pas))
	for _, pa := range pas {
		out = append(out, toPrintAreaDTO(pa))
	}
	return out
}
