package domain

import "time"

// OptionValue is a single value of a product option.
type OptionValue struct {
	ID    int
	Title string
}

// ProductOption is a configurable option (e.g. size) on a product.
type ProductOption struct {
	Name   string
	Type   string
	Values []OptionValue
}

// ProductVariant is a sellable variant of a product.
type ProductVariant struct {
	ID              VariantID
	SKU             string
	Cost            int
	Price           int
	Title           string
	Grams           int
	Enabled         bool
	IsDefault       bool
	Available       bool
	ExpressEligible bool
	Options         []int
}

// ProductImage is an image attached to product variants.
type ProductImage struct {
	Src       string
	VariantID []VariantID
	Position  string
	IsDefault bool
}

// Pattern describes repeating placement of an image.
type Pattern struct {
	SpacingX int
	SpacingY int
	Scale    int
	Offset   int
}

// PlaceholderImage is an image placed on a print placeholder.
type PlaceholderImage struct {
	ID         string
	Src        *string
	Name       *string
	Type       string
	Height     int
	Width      int
	X          *float64
	Y          *float64
	Scale      *float64
	Angle      int
	FontFamily *string
	FontSize   *int
	FontWeight *int
	FontColor  *string
	FontStyle  *string
	InputText  *string
	TextAlign  *string
	Pattern    *Pattern
}

// Placeholder groups placeholder images by print position.
type Placeholder struct {
	Position string
	Images   []PlaceholderImage
}

// PrintArea defines the print region for a set of variants.
type PrintArea struct {
	VariantIDs   []VariantID
	Placeholders []Placeholder
	Background   string
}

// ViewFile is a file (e.g. SVG) for a product view.
type ViewFile struct {
	Src       string
	VariantID []VariantID
}

// View is a named side of a product (e.g. front).
type View struct {
	ID       int
	Label    string
	Position string
	Files    []ViewFile
}

// GpsrInfo is a GPSR compliance entry for a product.
type GpsrInfo struct {
	Title string
	Text  string
}

// Product is a shop product based on a blueprint.
type Product struct {
	ID                ProductID
	Title             string
	Description       *string
	SafetyInformation *string
	Tags              []string
	Options           []ProductOption
	Variants          []ProductVariant
	Images            []ProductImage
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Visible           bool
	Locked            bool
	ExpressEligible   bool
	ExpressEnabled    bool
	EconomyEligible   bool
	EconomyEnabled    bool
	BlueprintID       BlueprintID
	UserID            int
	ShopID            ShopID
	PrintProviderID   PrintProviderID
	PrintAreas        []PrintArea
	Views             []View
	SalesChannelProps []any
}

// CreateProduct describes the payload to create a product.
type CreateProduct struct {
	Title             string
	Description       *string
	SafetyInformation *string
	BlueprintID       BlueprintID
	PrintProviderID   PrintProviderID
	Variants          []ProductVariant
	PrintAreas        []PrintArea
}

// UpdateProduct describes the payload to update a product.
type UpdateProduct struct {
	Title             string
	Description       *string
	SafetyInformation *string
	BlueprintID       BlueprintID
	PrintProviderID   PrintProviderID
	Variants          []ProductVariant
	PrintAreas        []PrintArea
}

// PublishProduct selects which fields are pushed to the sales channel.
type PublishProduct struct {
	Title            bool
	Description      bool
	Images           bool
	Variants         bool
	Tags             bool
	KeyFeatures      bool
	ShippingTemplate bool
}

// PublishSucceeded carries the externally published product reference.
type PublishSucceeded struct {
	ExternalID     string
	ExternalHandle string
}

// PublishFailed carries the reason publishing failed.
type PublishFailed struct {
	Reason string
}
