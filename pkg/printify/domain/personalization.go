package domain

// PersonalizationOption is a configurable personalization field.
type PersonalizationOption struct {
	FieldID        string
	Label          string
	Type           string
	CharacterLimit *int
}

// PersonalizationConfigItem is one field's input.
type PersonalizationConfigItem struct {
	FieldID string
	Type    string
	// Text is used when Type is a text field.
	Text string
	// ImageURL is used when Type is an image field.
	ImageURL string
}

// CreatePersonalizationConfig is the payload to configure personalization.
type CreatePersonalizationConfig struct {
	VariantID VariantID
	Items     []PersonalizationConfigItem
}

// PersonalizationConfig is the opaque configuration result to pass to an order.
type PersonalizationConfig struct {
	Strategy     string
	Instructions string
}

// PreviewMockup is one rendered preview image.
type PreviewMockup struct {
	VariantID VariantID
	MockupID  string
	Src       string
}

// PreviewTask is the state of an async personalization preview task.
type PreviewTask struct {
	TaskID       TaskID
	Status       string
	PreviewCount int
	Mockups      []PreviewMockup
	Error        *string
}

// CreatePreviewTask is the payload to request a personalization preview.
type CreatePreviewTask struct {
	ExternalID      string
	VariantIDs      []VariantID
	Personalisation Personalisation
}
