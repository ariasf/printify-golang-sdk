package domain

// ShippingMethod is a v2 catalog shipping method.
type ShippingMethod string

// Available v2 catalog shipping methods.
const (
	ShippingStandard ShippingMethod = "standard"
	ShippingPriority ShippingMethod = "priority"
	ShippingExpress  ShippingMethod = "express"
	ShippingEconomy  ShippingMethod = "economy"
)

// Location is a physical address of a print provider.
type Location struct {
	Address1 string
	Address2 string
	City     string
	Country  string
	Region   string
	Zip      string
}

// Blueprint is a product design in the Printify catalog.
type Blueprint struct {
	ID               BlueprintID
	Title            string
	Description      *string
	Brand            string
	Model            string
	Images           []string
	Tags             []string
	Features         []string
	CareInstructions []string
}

// PrintProvider is a manufacturer that fulfills orders.
type PrintProvider struct {
	ID                PrintProviderID
	Title             string
	Location          *Location
	Blueprints        []Blueprint
	DecorationMethods []string
}

// PrintProviderRef is an id-only reference to a print provider.
type PrintProviderRef struct {
	ID                PrintProviderID
	Title             string
	DecorationMethods []string
}

// VariantOptions captures the option values of a variant.
type VariantOptions struct {
	Color string
	Size  string
}

// VariantPlaceholder is a print position on a blueprint.
type VariantPlaceholder struct {
	Position         string
	DecorationMethod string
	Height           int
	Width            int
}

// Variant is a purchasable option combination of a blueprint.
type Variant struct {
	ID                VariantID
	Title             string
	Options           *VariantOptions
	Placeholders      []VariantPlaceholder
	DecorationMethods []string
}

// Variants is the variants response for a blueprint and print provider.
type Variants struct {
	BlueprintID BlueprintID
	Title       string
	Variants    []Variant
}

// HandlingTime is how long a provider takes before shipping.
type HandlingTime struct {
	Value int
	Unit  string
}

// ShippingCost is a cost in minor currency units.
type ShippingCost struct {
	Cost     int
	Currency string
}

// ShippingProfile groups variants with the same shipping terms.
type ShippingProfile struct {
	VariantIDs      []VariantID
	FirstItem       ShippingCost
	AdditionalItems ShippingCost
	Countries       []string
}

// BlueprintShipping is the shipping information for a blueprint.
type BlueprintShipping struct {
	HandlingTime HandlingTime
	Profiles     []ShippingProfile
}

// SizeRange is a measurement range for one size.
type SizeRange struct {
	From int
	To   int
}

// SizeType is a measurement type (e.g. width) with one range per size.
type SizeType struct {
	Name   string
	Unit   string
	Values []SizeRange
}

// SizeGuide lists available sizes and measurement ranges for a blueprint.
type SizeGuide struct {
	Sizes []string
	Types []SizeType
}
