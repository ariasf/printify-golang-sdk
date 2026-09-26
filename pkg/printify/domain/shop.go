// Package domain holds the SDK's entities and value objects.
// It has no dependencies on other layers.
package domain

// Shop is a sales channel connected to a Printify account.
type Shop struct {
	ID      ShopID
	Title   string
	Channel string
}
