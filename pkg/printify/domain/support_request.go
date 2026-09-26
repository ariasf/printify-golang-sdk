package domain

import "time"

// SupportRequest is an order support case (refund, reprint, address change).
type SupportRequest struct {
	ID            SupportRequestID
	OrderID       OrderID
	Type          string
	Status        string
	StatusDetails string
	Outcome       string
	LineItemIDs   []string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// SupportRequestLineItem targets one order line item.
type SupportRequestLineItem struct {
	LineItemID string
	Quantity   int
}

// ReprintRequest is the payload for a reprint support request.
type ReprintRequest struct {
	Reason      string
	Description string
	ImageURLs   []string
	LineItems   []SupportRequestLineItem
}

// RefundRequest is the payload for a refund support request.
type RefundRequest struct {
	Reason      string
	Description string
	ImageURLs   []string
	LineItems   []SupportRequestLineItem
}

// AddressChange is the payload to change an order's shipping address.
type AddressChange struct {
	FirstName string
	LastName  string
	Address1  string
	Address2  string
	City      string
	Region    string
	Country   string
	Zip       string
	Email     string
	Phone     string
}
