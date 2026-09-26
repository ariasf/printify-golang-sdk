package rest

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/printify-go/pkg/printify/domain"
	"github.com/printify-go/pkg/printify/ports/driving"
)

// UploadGateway implements driven.UploadGateway over the Printify REST API.
type UploadGateway struct {
	client *Client
}

// NewUploadGateway builds an UploadGateway using the shared transport client.
func NewUploadGateway(client *Client) *UploadGateway {
	return &UploadGateway{client: client}
}

type uploadDTO struct {
	ID         string `json:"id"`
	FileName   string `json:"file_name"`
	Height     int    `json:"height"`
	Width      int    `json:"width"`
	Size       int    `json:"size"`
	MimeType   string `json:"mime_type"`
	PreviewURL string `json:"preview_url"`
	UploadTime string `json:"upload_time"`
}

func (d uploadDTO) toDomain() domain.Upload {
	return domain.Upload{
		ID: domain.UploadID(d.ID), FileName: d.FileName, Height: d.Height, Width: d.Width,
		Size: d.Size, MimeType: d.MimeType, PreviewURL: d.PreviewURL, UploadTime: parseTime(d.UploadTime),
	}
}

func uploadListQuery(f driving.UploadFilter) url.Values {
	q := url.Values{}
	if f.Limit != nil {
		q.Set("limit", fmt.Sprint(*f.Limit))
	}
	if f.Page != nil {
		q.Set("page", fmt.Sprint(*f.Page))
	}
	return q
}

// List fetches a page of uploads.
func (g *UploadGateway) List(ctx context.Context, f driving.UploadFilter) ([]domain.Upload, error) {
	var dtos []uploadDTO
	if err := g.client.do(ctx, http.MethodGet, "/v1/uploads.json", uploadListQuery(f), nil, &dtos); err != nil {
		return nil, err
	}
	out := make([]domain.Upload, 0, len(dtos))
	for _, d := range dtos {
		out = append(out, d.toDomain())
	}
	return out, nil
}

// Get fetches one upload.
func (g *UploadGateway) Get(ctx context.Context, id domain.UploadID) (*domain.Upload, error) {
	var dto uploadDTO
	path := fmt.Sprintf("/v1/uploads/%s.json", id)
	if err := g.client.do(ctx, http.MethodGet, path, nil, nil, &dto); err != nil {
		return nil, err
	}
	u := dto.toDomain()
	return &u, nil
}

// UploadImage uploads an image from a URL or base64 contents.
func (g *UploadGateway) UploadImage(ctx context.Context, img domain.UploadImage) (*domain.Upload, error) {
	body := map[string]any{"file_name": img.FileName}
	if img.URL != "" {
		body["url"] = img.URL
	}
	if img.Contents != "" {
		body["contents"] = img.Contents
	}
	var dto uploadDTO
	if err := g.client.do(ctx, http.MethodPost, "/v1/uploads/images.json", nil, body, &dto); err != nil {
		return nil, err
	}
	u := dto.toDomain()
	return &u, nil
}

// Archive removes an uploaded image.
func (g *UploadGateway) Archive(ctx context.Context, id domain.UploadID) error {
	path := fmt.Sprintf("/v1/uploads/%s/archive.json", id)
	return g.client.do(ctx, http.MethodPost, path, nil, nil, nil)
}
