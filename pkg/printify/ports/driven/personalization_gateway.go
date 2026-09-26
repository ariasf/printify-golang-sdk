package driven

import (
	"context"

	"github.com/printify-go/pkg/printify/domain"
)

// PersonalizationGateway abstracts the Printify personalization API.
type PersonalizationGateway interface {
	ListOptions(ctx context.Context, shop domain.ShopID, product domain.ProductID) ([]domain.PersonalizationOption, error)
	CreateConfig(ctx context.Context, shop domain.ShopID, product domain.ProductID, cfg domain.CreatePersonalizationConfig) (*domain.PersonalizationConfig, error)
	CreatePreviewTask(ctx context.Context, shop domain.ShopID, product domain.ProductID, task domain.CreatePreviewTask) (*domain.PreviewTask, error)
	GetPreviewTask(ctx context.Context, shop domain.ShopID, product domain.ProductID, id domain.TaskID) (*domain.PreviewTask, error)
}
