// Package printify is the public entrypoint of the SDK. It only wires
// adapters into services; all logic lives in the services layer.
package printify

import (
	"errors"
	"net/http"

	"github.com/printify-go/pkg/printify/adapters/rest"
	"github.com/printify-go/pkg/printify/ports/driving"
	"github.com/printify-go/pkg/printify/services"
)

// Filter types exposed at the facade level so users can write
// `printify.OrderFilter{...}` instead of `driving.OrderFilter{...}`.
// They are plain aliases: the underlying service methods accept the
// identical struct by value.
type (
	OrderFilter   = driving.OrderFilter
	ProductFilter = driving.ProductFilter
	UploadFilter  = driving.UploadFilter
	WebhookFilter = driving.WebhookFilter
)

// Config configures the SDK client.
type Config struct {
	// APIToken is the Printify personal access token. Required.
	APIToken string
	// BaseURL overrides the API endpoint. Defaults to the production API.
	BaseURL string
	// HTTPClient overrides the underlying HTTP client. Defaults to http.DefaultClient.
	HTTPClient *http.Client
}

// SDK exposes the Printify API through driving ports.
type SDK struct {
	Shops           driving.ShopService
	Catalog         driving.CatalogService
	Products        driving.ProductService
	Orders          driving.OrderService
	Personalization driving.PersonalizationService
	SupportRequests driving.SupportRequestService
	Uploads         driving.UploadService
	Webhooks        driving.WebhookService
}

// New wires the REST adapters into the services and returns a ready SDK.
func New(cfg Config) (*SDK, error) {
	if cfg.APIToken == "" {
		return nil, errors.New("printify: APIToken is required")
	}
	var doer rest.Doer
	if cfg.HTTPClient != nil {
		doer = cfg.HTTPClient
	}
	client := rest.NewClient(cfg.BaseURL, cfg.APIToken, doer)

	return &SDK{
		Shops:           services.NewShopService(rest.NewShopGateway(client)),
		Catalog:         services.NewCatalogService(rest.NewCatalogGateway(client)),
		Products:        services.NewProductService(rest.NewProductGateway(client)),
		Orders:          services.NewOrderService(rest.NewOrderGateway(client)),
		Personalization: services.NewPersonalizationService(rest.NewPersonalizationGateway(client)),
		SupportRequests: services.NewSupportRequestService(rest.NewSupportRequestGateway(client)),
		Uploads:         services.NewUploadService(rest.NewUploadGateway(client)),
		Webhooks:        services.NewWebhookService(rest.NewWebhookGateway(client)),
	}, nil
}
