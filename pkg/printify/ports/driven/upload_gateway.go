package driven

import (
	"context"

	"github.com/printify-go/pkg/printify/domain"
	"github.com/printify-go/pkg/printify/ports/driving"
)

// UploadGateway abstracts the Printify uploads API.
type UploadGateway interface {
	List(ctx context.Context, f driving.UploadFilter) ([]domain.Upload, error)
	Get(ctx context.Context, id domain.UploadID) (*domain.Upload, error)
	UploadImage(ctx context.Context, img domain.UploadImage) (*domain.Upload, error)
	Archive(ctx context.Context, id domain.UploadID) error
}
