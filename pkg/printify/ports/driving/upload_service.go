package driving

import (
	"context"

	"github.com/printify-go/pkg/printify/domain"
)

// UploadFilter narrows the upload list.
type UploadFilter struct {
	Limit *int
	Page  *int
}

// UploadService exposes image uploads to SDK consumers.
type UploadService interface {
	List(ctx context.Context, f UploadFilter) ([]domain.Upload, error)
	Get(ctx context.Context, id domain.UploadID) (*domain.Upload, error)
	UploadImage(ctx context.Context, img domain.UploadImage) (*domain.Upload, error)
	Archive(ctx context.Context, id domain.UploadID) error
}
