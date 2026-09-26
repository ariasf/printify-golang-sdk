package domain

import "time"

// Address is a shipping address.
type Address struct {
	FirstName string
	LastName  string
	Region    string
	Address1  string
	Address2  string
	City      string
	Zip       string
	Email     string
	Phone     string
	Country   string
	Company   string
}

// LineItemMetadata carries buyer-side line item metadata.
type LineItemMetadata struct {
	Title        string
	Price        int
	VariantLabel string
	SKU          string
	Country      string
	ExternalID   string
}

// LineItem is one item on an order.
type LineItem struct {
	ProductID          ProductID
	Quantity           int
	VariantID          VariantID
	PrintProviderID    PrintProviderID
	Cost               int
	ShippingCost       int
	Status             string
	Metadata           *LineItemMetadata
	SentToProductionAt *time.Time
	FulfilledAt        *time.Time
}

// OrderMetadata describes how the order was created.
type OrderMetadata struct {
	OrderType              string
	ShopOrderID            int
	ShopOrderLabel         string
	ShopFulfilledAt        *time.Time
	IsReprint              bool
	ReprintedOrderIDs      []string
	ChildReprintedOrderIDs []string
}

// Shipment is a carrier shipment for an order.
type Shipment struct {
	Carrier     string
	Number      string
	URL         string
	DeliveredAt *time.Time
}

// Order is a Printify order.
type Order struct {
	ID                 OrderID
	AppOrderID         string
	AddressTo          Address
	LineItems          []LineItem
	Metadata           *OrderMetadata
	TotalPrice         int
	TotalShipping      int
	TotalTax           int
	Status             string
	ShippingMethod     int
	IsPrintifyExpress  bool
	IsEconomyShipping  bool
	Shipments          []Shipment
	CreatedAt          time.Time
	SentToProductionAt *time.Time
	FulfilledAt        *time.Time
}

// OrderIDResult is the id of a newly created order.
type OrderIDResult struct {
	ID OrderID
}

// ShippingCosts are per-method shipping costs in minor currency units.
type ShippingCosts struct {
	Standard        *int
	Express         *int
	Priority        *int
	PrintifyExpress *int
	Economy         *int
}

// Personalisation attaches a personalization to a line item.
type Personalisation struct {
	Strategy     string
	Instructions string
	BuyerRequest string
}

// LineItemBlueprint references a blueprint directly in an order.
type LineItemBlueprint struct {
	PrintProviderID PrintProviderID
	BlueprintID     BlueprintID
	VariantID       VariantID
	PrintAreas      map[string]string
	Quantity        int
	ExternalID      string
	Personalisation *Personalisation
}

// SubmitOrder is the payload to create an order.
type SubmitOrder struct {
	ExternalID               string
	Label                    string
	LineItems                []LineItemBlueprint
	ShippingMethod           *int
	IsPrintifyExpress        *bool
	IsEconomyShipping        *bool
	SendShippingNotification *bool
	AddressTo                Address
}

// ExpressOrder is the payload to create a Printify Express order.
type ExpressOrder struct {
	ExternalID               string
	Label                    string
	LineItems                []LineItemBlueprint
	ShippingMethod           *int
	SendShippingNotification *bool
	AddressTo                Address
}

// AddressChangeResult is the outcome of an address change request.
type AddressChangeResult struct {
	// Resolution is "updated_directly" or "support_request_created".
	Resolution     string
	OrderID        OrderID
	SupportRequest *SupportRequest
	AddressTo      *Address
}
