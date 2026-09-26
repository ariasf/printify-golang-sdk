package domain

// Typed identifiers for Printify resources.
type (
	ShopID           int64
	BlueprintID      int64
	PrintProviderID  int64
	VariantID        int64
	ProductID        string
	OrderID          string
	UploadID         string
	WebhookID        string
	SupportRequestID string
	TaskID           string
)

func (id ShopID) Valid() bool           { return id > 0 }
func (id BlueprintID) Valid() bool      { return id > 0 }
func (id PrintProviderID) Valid() bool  { return id > 0 }
func (id VariantID) Valid() bool        { return id > 0 }
func (id ProductID) Valid() bool        { return id != "" }
func (id OrderID) Valid() bool          { return id != "" }
func (id UploadID) Valid() bool         { return id != "" }
func (id WebhookID) Valid() bool        { return id != "" }
func (id SupportRequestID) Valid() bool { return id != "" }
func (id TaskID) Valid() bool           { return id != "" }
