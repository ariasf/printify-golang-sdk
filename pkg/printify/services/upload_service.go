package services

import (
	"context"
	"fmt"

	"github.com/printify-go/pkg/printify/domain"
	"github.com/printify-go/pkg/printify/ports/driven"
	"github.com/printify-go/pkg/printify/ports/driving"
)

// UploadService implements driving.UploadService.
type UploadService struct {
	uploads driven.UploadGateway
}

// NewUploadService builds an UploadService backed by the given gateway.
func NewUploadService(uploads driven.UploadGateway) *UploadService {
	return &UploadService{uploads: uploads}
}

func validUpload(id domain.UploadID) error {
	if !id.Valid() {
		return fmt.Errorf("upload id %q: %w", id, domain.ErrInvalidInput)
	}
	return nil
}

// List returns a page of uploaded images.
func (s *UploadService) List(ctx context.Context, f driving.UploadFilter) ([]domain.Upload, error) {
	uploads, err := s.uploads.List(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("listing uploads: %w", err)
	}
	return uploads, nil
}

// Get returns one uploaded image.
func (s *UploadService) Get(ctx context.Context, id domain.UploadID) (*domain.Upload, error) {
	if err := validUpload(id); err != nil {
		return nil, err
	}
	u, err := s.uploads.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting upload %s: %w", id, err)
	}
	return u, nil
}

// UploadImage uploads an image from a URL or base64 contents.
func (s *UploadService) UploadImage(ctx context.Context, img domain.UploadImage) (*domain.Upload, error) {
	if err := validateUploadImage(img); err != nil {
		return nil, err
	}
	u, err := s.uploads.UploadImage(ctx, img)
	if err != nil {
		return nil, fmt.Errorf("uploading image %q: %w", img.FileName, err)
	}
	return u, nil
}

// Archive removes an uploaded image.
func (s *UploadService) Archive(ctx context.Context, id domain.UploadID) error {
	if err := validUpload(id); err != nil {
		return err
	}
	if err := s.uploads.Archive(ctx, id); err != nil {
		return fmt.Errorf("archiving upload %s: %w", id, err)
	}
	return nil
}

func validateUploadImage(img domain.UploadImage) error {
	if img.FileName == "" {
		return fmt.Errorf("file name: %w", domain.ErrInvalidInput)
	}
	if (img.URL == "") == (img.Contents == "") {
		return fmt.Errorf("exactly one of url or contents required, %w", domain.ErrInvalidInput)
	}
	return nil
}
