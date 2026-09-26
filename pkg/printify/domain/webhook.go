package domain

// Webhook is a subscription to Printify events.
type Webhook struct {
	ID     WebhookID
	Topic  string
	URL    string
	ShopID ShopID
}

// CreateWebhook is the payload to subscribe a URL to a topic.
type CreateWebhook struct {
	Topic string
	URL   string
}

// WebhookTopic is an event topic a webhook can subscribe to.
type WebhookTopic string

// Known webhook topics.
const (
	TopicShopDisconnected                WebhookTopic = "shop:disconnected"
	TopicProductDeleted                  WebhookTopic = "product:deleted"
	TopicProductCreated                  WebhookTopic = "product:created"
	TopicProductUpdated                  WebhookTopic = "product:updated"
	TopicProductPublishStarted           WebhookTopic = "product:publish:started"
	TopicOrderCreated                    WebhookTopic = "order:created"
	TopicOrderUpdated                    WebhookTopic = "order:updated"
	TopicOrderShipmentCreated            WebhookTopic = "order:shipment:created"
	TopicOrderShipmentDelivered          WebhookTopic = "order:shipment:delivered"
	TopicOrderSentToProduction           WebhookTopic = "order:sent-to-production"
	TopicPersonalizationPreviewProcessed WebhookTopic = "personalization-preview-task:processed"
	TopicSupportRequestUpdated           WebhookTopic = "support-request:updated"
)
