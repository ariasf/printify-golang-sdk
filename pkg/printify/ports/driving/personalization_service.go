package driving

import (
	"context"

	"github.com/printify-go/pkg/printify/domain"
)

// PersonalizationService exposes personalization configuration to SDK consumers.
type PersonalizationService interface {
	ListOptions(ctx context.Context, shop domain.ShopID, product domain.ProductID) ([]domain.PersonalizationOption, error)
	CreateConfig(ctx context.Context, shop domain.ShopID, product domain.ProductID, cfg domain.CreatePersonalizationConfig) (*domain.PersonalizationConfig, error)
	CreatePreviewTask(ctx context.Context, shop domain.ShopID, product domain.ProductID, task domain.CreatePreviewTask) (*domain.PreviewTask, error)
	GetPreviewTask(ctx context.Context, shop domain.ShopID, product domain.ProductID, id domain.TaskID) (*domain.PreviewTask, error)
}
