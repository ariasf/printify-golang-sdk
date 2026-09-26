package rest

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/printify-go/pkg/printify/domain"
)

// CatalogGateway implements driven.CatalogGateway over the Printify REST API.
type CatalogGateway struct {
	client *Client
}

// NewCatalogGateway builds a CatalogGateway using the shared transport client.
func NewCatalogGateway(client *Client) *CatalogGateway {
	return &CatalogGateway{client: client}
}

type blueprintDTO struct {
	ID               int64    `json:"id"`
	Title            string   `json:"title"`
	Description      *string  `json:"description"`
	Brand            string   `json:"brand"`
	Model            string   `json:"model"`
	Images           []string `json:"images"`
	Tags             []string `json:"tags"`
	Features         []string `json:"features"`
	CareInstructions []string `json:"care_instructions"`
}

func (d blueprintDTO) toDomain() domain.Blueprint {
	return domain.Blueprint{
		ID:               domain.BlueprintID(d.ID),
		Title:            d.Title,
		Description:      d.Description,
		Brand:            d.Brand,
		Model:            d.Model,
		Images:           d.Images,
		Tags:             d.Tags,
		Features:         d.Features,
		CareInstructions: d.CareInstructions,
	}
}

type locationDTO struct {
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	City     string `json:"city"`
	Country  string `json:"country"`
	Region   string `json:"region"`
	Zip      string `json:"zip"`
}

type printProviderDTO struct {
	ID                int64          `json:"id"`
	Title             string         `json:"title"`
	Location          *locationDTO   `json:"location"`
	Blueprints        []blueprintDTO `json:"blueprints"`
	DecorationMethods []string       `json:"decoration_methods"`
}

func (d printProviderDTO) toDomain() domain.PrintProvider {
	pp := domain.PrintProvider{
		ID:                domain.PrintProviderID(d.ID),
		Title:             d.Title,
		DecorationMethods: d.DecorationMethods,
	}
	if d.Location != nil {
		pp.Location = &domain.Location{
			Address1: d.Location.Address1,
			Address2: d.Location.Address2,
			City:     d.Location.City,
			Country:  d.Location.Country,
			Region:   d.Location.Region,
			Zip:      d.Location.Zip,
		}
	}
	for _, b := range d.Blueprints {
		pp.Blueprints = append(pp.Blueprints, b.toDomain())
	}
	return pp
}

type printProviderRefDTO struct {
	ID                int64    `json:"id"`
	Title             string   `json:"title"`
	DecorationMethods []string `json:"decoration_methods"`
}

type variantPlaceholderDTO struct {
	Position         string `json:"position"`
	DecorationMethod string `json:"decoration_method"`
	Height           int    `json:"height"`
	Width            int    `json:"width"`
}

type variantOptionsDTO struct {
	Color string `json:"color"`
	Size  string `json:"size"`
}

type variantDTO struct {
	ID                int64                   `json:"id"`
	Title             string                  `json:"title"`
	Options           *variantOptionsDTO      `json:"options"`
	Placeholders      []variantPlaceholderDTO `json:"placeholders"`
	DecorationMethods []string                `json:"decoration_methods"`
}

func (d variantDTO) toDomain() domain.Variant {
	v := domain.Variant{
		ID:                domain.VariantID(d.ID),
		Title:             d.Title,
		Placeholders:      make([]domain.VariantPlaceholder, 0, len(d.Placeholders)),
		DecorationMethods: d.DecorationMethods,
	}
	if d.Options != nil {
		v.Options = &domain.VariantOptions{Color: d.Options.Color, Size: d.Options.Size}
	}
	for _, p := range d.Placeholders {
		v.Placeholders = append(v.Placeholders, domain.VariantPlaceholder{
			Position:         p.Position,
			DecorationMethod: p.DecorationMethod,
			Height:           p.Height,
			Width:            p.Width,
		})
	}
	return v
}

type variantsDTO struct {
	ID       int64        `json:"id"`
	Title    string       `json:"title"`
	Variants []variantDTO `json:"variants"`
}

type handlingTimeDTO struct {
	Value int    `json:"value"`
	Unit  string `json:"unit"`
}

type shippingCostDTO struct {
	Cost     int    `json:"cost"`
	Currency string `json:"currency"`
}

type shippingProfileDTO struct {
	VariantIDs      []int64         `json:"variant_ids"`
	FirstItem       shippingCostDTO `json:"first_item"`
	AdditionalItems shippingCostDTO `json:"additional_items"`
	Countries       []string        `json:"countries"`
}

type shippingDTO struct {
	HandlingTime handlingTimeDTO      `json:"handling_time"`
	Profiles     []shippingProfileDTO `json:"profiles"`
}

type sizeRangeDTO struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type sizeTypeDTO struct {
	Name   string         `json:"name"`
	Unit   string         `json:"unit"`
	Values []sizeRangeDTO `json:"values"`
}

type sizeGuideDTO struct {
	Sizes []string      `json:"sizes"`
	Types []sizeTypeDTO `json:"types"`
}

// ListBlueprints fetches all catalog blueprints.
func (g *CatalogGateway) ListBlueprints(ctx context.Context) ([]domain.Blueprint, error) {
	var dtos []blueprintDTO
	if err := g.client.do(ctx, http.MethodGet, "/v1/catalog/blueprints.json", nil, nil, &dtos); err != nil {
		return nil, err
	}
	out := make([]domain.Blueprint, 0, len(dtos))
	for _, d := range dtos {
		out = append(out, d.toDomain())
	}
	return out, nil
}

// GetBlueprint fetches one blueprint.
func (g *CatalogGateway) GetBlueprint(ctx context.Context, id domain.BlueprintID) (*domain.Blueprint, error) {
	var dto blueprintDTO
	path := fmt.Sprintf("/v1/catalog/blueprints/%d.json", id)
	if err := g.client.do(ctx, http.MethodGet, path, nil, nil, &dto); err != nil {
		return nil, err
	}
	bp := dto.toDomain()
	return &bp, nil
}

// GetSizeGuide fetches the size guide for a blueprint.
func (g *CatalogGateway) GetSizeGuide(ctx context.Context, id domain.BlueprintID) (*domain.SizeGuide, error) {
	var dto sizeGuideDTO
	path := fmt.Sprintf("/v1/catalog/blueprints/%d/size_guide.json", id)
	if err := g.client.do(ctx, http.MethodGet, path, nil, nil, &dto); err != nil {
		return nil, err
	}
	return &domain.SizeGuide{Sizes: dto.Sizes, Types: toSizeTypes(dto.Types)}, nil
}

// ListPrintProviders fetches all print providers.
func (g *CatalogGateway) ListPrintProviders(ctx context.Context) ([]domain.PrintProvider, error) {
	var dtos []printProviderDTO
	if err := g.client.do(ctx, http.MethodGet, "/v1/catalog/print_providers.json", nil, nil, &dtos); err != nil {
		return nil, err
	}
	out := make([]domain.PrintProvider, 0, len(dtos))
	for _, d := range dtos {
		out = append(out, d.toDomain())
	}
	return out, nil
}

// GetPrintProvider fetches one print provider.
func (g *CatalogGateway) GetPrintProvider(ctx context.Context, id domain.PrintProviderID) (*domain.PrintProvider, error) {
	var dto printProviderDTO
	path := fmt.Sprintf("/v1/catalog/print_providers/%d.json", id)
	if err := g.client.do(ctx, http.MethodGet, path, nil, nil, &dto); err != nil {
		return nil, err
	}
	pp := dto.toDomain()
	return &pp, nil
}

// ListProvidersForBlueprint fetches the providers fulfilling a blueprint.
func (g *CatalogGateway) ListProvidersForBlueprint(ctx context.Context, blueprint domain.BlueprintID) ([]domain.PrintProviderRef, error) {
	var dtos []printProviderRefDTO
	path := fmt.Sprintf("/v1/catalog/blueprints/%d/print_providers.json", blueprint)
	if err := g.client.do(ctx, http.MethodGet, path, nil, nil, &dtos); err != nil {
		return nil, err
	}
	out := make([]domain.PrintProviderRef, 0, len(dtos))
	for _, d := range dtos {
		out = append(out, domain.PrintProviderRef{
			ID:                domain.PrintProviderID(d.ID),
			Title:             d.Title,
			DecorationMethods: d.DecorationMethods,
		})
	}
	return out, nil
}

// GetVariants fetches the variants of a blueprint for a provider.
func (g *CatalogGateway) GetVariants(ctx context.Context, blueprint domain.BlueprintID, provider domain.PrintProviderID, includeOutOfStock bool) (*domain.Variants, error) {
	path := fmt.Sprintf("/v1/catalog/blueprints/%d/print_providers/%d/variants.json", blueprint, provider)
	query := url.Values{}
	if includeOutOfStock {
		query.Set("show-out-of-stock", "1")
	}
	var dto variantsDTO
	if err := g.client.do(ctx, http.MethodGet, path, query, nil, &dto); err != nil {
		return nil, err
	}
	v := &domain.Variants{
		BlueprintID: domain.BlueprintID(dto.ID),
		Title:       dto.Title,
		Variants:    make([]domain.Variant, 0, len(dto.Variants)),
	}
	for _, d := range dto.Variants {
		v.Variants = append(v.Variants, d.toDomain())
	}
	return v, nil
}

// GetShipping fetches shipping information for a blueprint.
func (g *CatalogGateway) GetShipping(ctx context.Context, blueprint domain.BlueprintID) (*domain.BlueprintShipping, error) {
	var dto shippingDTO
	path := fmt.Sprintf("/v1/catalog/blueprints/%d/print_providers/shipping.json", blueprint)
	if err := g.client.do(ctx, http.MethodGet, path, nil, nil, &dto); err != nil {
		return nil, err
	}
	s := &domain.BlueprintShipping{
		HandlingTime: domain.HandlingTime{Value: dto.HandlingTime.Value, Unit: dto.HandlingTime.Unit},
		Profiles:     make([]domain.ShippingProfile, 0, len(dto.Profiles)),
	}
	for _, p := range dto.Profiles {
		s.Profiles = append(s.Profiles, domain.ShippingProfile{
			VariantIDs:      toVariantIDs(p.VariantIDs),
			FirstItem:       domain.ShippingCost{Cost: p.FirstItem.Cost, Currency: p.FirstItem.Currency},
			AdditionalItems: domain.ShippingCost{Cost: p.AdditionalItems.Cost, Currency: p.AdditionalItems.Currency},
			Countries:       p.Countries,
		})
	}
	return s, nil
}

// GetShippingMethod fetches raw v2 shipping data for one method.
func (g *CatalogGateway) GetShippingMethod(ctx context.Context, blueprint domain.BlueprintID, method domain.ShippingMethod) (map[string]any, error) {
	var out map[string]any
	path := fmt.Sprintf("/v1/catalog/blueprints/%d/print_providers/shipping/%s.json", blueprint, method)
	if err := g.client.do(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// GetShippingV2 fetches the raw v2 shipping list.
func (g *CatalogGateway) GetShippingV2(ctx context.Context, blueprint domain.BlueprintID, provider domain.PrintProviderID) (map[string]any, error) {
	var out map[string]any
	path := fmt.Sprintf("/v2/catalog/blueprints/%d/print_providers/%d/shipping.json", blueprint, provider)
	if err := g.client.do(ctx, http.MethodGet, path, nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func toVariantIDs(ids []int64) []domain.VariantID {
	out := make([]domain.VariantID, 0, len(ids))
	for _, id := range ids {
		out = append(out, domain.VariantID(id))
	}
	return out
}

func toSizeTypes(dtos []sizeTypeDTO) []domain.SizeType {
	out := make([]domain.SizeType, 0, len(dtos))
	for _, d := range dtos {
		st := domain.SizeType{Name: d.Name, Unit: d.Unit}
		for _, v := range d.Values {
			st.Values = append(st.Values, domain.SizeRange{From: v.From, To: v.To})
		}
		out = append(out, st)
	}
	return out
}
