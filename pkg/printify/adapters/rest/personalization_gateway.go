package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/printify-go/pkg/printify/domain"
)

// PersonalizationGateway implements driven.PersonalizationGateway over the Printify REST API.
type PersonalizationGateway struct {
	client *Client
}

// NewPersonalizationGateway builds a PersonalizationGateway using the shared transport client.
func NewPersonalizationGateway(client *Client) *PersonalizationGateway {
	return &PersonalizationGateway{client: client}
}

type personalizationOptionDTO struct {
	FieldID        string `json:"field_id"`
	Label          string `json:"label"`
	Type           string `json:"type"`
	CharacterLimit *int   `json:"character_limit"`
}

type personalizationConfigItemDTO struct {
	FieldID string         `json:"field_id"`
	Type    string         `json:"type"`
	Input   map[string]any `json:"input"`
}

type personalizationConfigDTO struct {
	Strategy     string `json:"personalisation_strategy"`
	Instructions string `json:"personalisation_instructions"`
}

type previewMockupDTO struct {
	VariantID int64  `json:"variant_id"`
	MockupID  string `json:"mockup_id"`
	Src       string `json:"src"`
}

type previewTaskDTO struct {
	TaskID       string             `json:"task_id"`
	Status       string             `json:"status"`
	PreviewCount int                `json:"preview_count"`
	Mockups      []previewMockupDTO `json:"mockups"`
	Error        *string            `json:"error"`
}

func (d previewTaskDTO) toDomain() domain.PreviewTask {
	t := domain.PreviewTask{
		TaskID: domain.TaskID(d.TaskID), Status: d.Status, PreviewCount: d.PreviewCount, Error: d.Error,
	}
	for _, m := range d.Mockups {
		t.Mockups = append(t.Mockups, domain.PreviewMockup{
			VariantID: domain.VariantID(m.VariantID), MockupID: m.MockupID, Src: m.Src,
		})
	}
	return t
}

// ListOptions fetches the personalization fields of a product.
func (g *PersonalizationGateway) ListOptions(ctx context.Context, shop domain.ShopID, product domain.ProductID) ([]domain.PersonalizationOption, error) {
	var dtos []personalizationOptionDTO
	path := fmt.Sprintf("/v1/shops/%d/products/%s/personalization_options.json", shop, product)
	if err := g.client.do(ctx, http.MethodGet, path, nil, nil, &dtos); err != nil {
		return nil, err
	}
	out := make([]domain.PersonalizationOption, 0, len(dtos))
	for _, d := range dtos {
		out = append(out, domain.PersonalizationOption{
			FieldID: d.FieldID, Label: d.Label, Type: d.Type, CharacterLimit: d.CharacterLimit,
		})
	}
	return out, nil
}

// CreateConfig configures personalization for a variant.
func (g *PersonalizationGateway) CreateConfig(ctx context.Context, shop domain.ShopID, product domain.ProductID, cfg domain.CreatePersonalizationConfig) (*domain.PersonalizationConfig, error) {
	items := make([]personalizationConfigItemDTO, 0, len(cfg.Items))
	for _, item := range cfg.Items {
		input := map[string]any{}
		if item.Text != "" {
			input["text"] = item.Text
		}
		if item.ImageURL != "" {
			input["image_url"] = item.ImageURL
		}
		items = append(items, personalizationConfigItemDTO{FieldID: item.FieldID, Type: item.Type, Input: input})
	}
	body := map[string]any{"variant_id": int64(cfg.VariantID), "items": items}
	var dto personalizationConfigDTO
	path := fmt.Sprintf("/v1/shops/%d/products/%s/personalization.json", shop, product)
	if err := g.client.do(ctx, http.MethodPost, path, nil, body, &dto); err != nil {
		return nil, err
	}
	config := &domain.PersonalizationConfig{Strategy: dto.Strategy, Instructions: dto.Instructions}
	return config, nil
}

// CreatePreviewTask queues a preview render.
func (g *PersonalizationGateway) CreatePreviewTask(ctx context.Context, shop domain.ShopID, product domain.ProductID, task domain.CreatePreviewTask) (*domain.PreviewTask, error) {
	variantIDs := make([]int64, 0, len(task.VariantIDs))
	for _, id := range task.VariantIDs {
		variantIDs = append(variantIDs, int64(id))
	}
	body := map[string]any{
		"variant_ids": variantIDs,
		"personalization": map[string]any{
			"personalisation_strategy":     task.Personalisation.Strategy,
			"personalisation_instructions": task.Personalisation.Instructions,
		},
	}
	if task.ExternalID != "" {
		body["external_id"] = task.ExternalID
	}
	var dto previewTaskDTO
	path := fmt.Sprintf("/v1/shops/%d/products/%s/personalization_previews.json", shop, product)
	if err := g.client.do(ctx, http.MethodPost, path, nil, body, &dto); err != nil {
		return nil, err
	}
	t := dto.toDomain()
	return &t, nil
}

// GetPreviewTask returns a preview task's status.
func (g *PersonalizationGateway) GetPreviewTask(ctx context.Context, shop domain.ShopID, product domain.ProductID, id domain.TaskID) (*domain.PreviewTask, error) {
	var dto previewTaskDTO
	path := fmt.Sprintf("/v1/shops/%d/products/%s/personalization_previews/tasks/%s.json", shop, product, id)
	if err := g.client.do(ctx, http.MethodGet, path, nil, nil, &dto); err != nil {
		return nil, err
	}
	t := dto.toDomain()
	return &t, nil
}
