package services

import (
	"context"
	"fmt"

	"github.com/printify-go/pkg/printify/domain"
	"github.com/printify-go/pkg/printify/ports/driven"
	"github.com/printify-go/pkg/printify/ports/driving"
)

// OrderService implements driving.OrderService.
type OrderService struct {
	orders driven.OrderGateway
}

// NewOrderService builds an OrderService backed by the given gateway.
func NewOrderService(orders driven.OrderGateway) *OrderService {
	return &OrderService{orders: orders}
}

func validOrder(id domain.OrderID) error {
	if !id.Valid() {
		return fmt.Errorf("order id %q: %w", id, domain.ErrInvalidInput)
	}
	return nil
}

// List returns a page of orders for a shop.
func (s *OrderService) List(ctx context.Context, shop domain.ShopID, f driving.OrderFilter) ([]domain.Order, error) {
	if err := validShop(shop); err != nil {
		return nil, err
	}
	orders, err := s.orders.List(ctx, shop, f)
	if err != nil {
		return nil, fmt.Errorf("listing orders in shop %d: %w", shop, err)
	}
	return orders, nil
}

// Get returns one order.
func (s *OrderService) Get(ctx context.Context, shop domain.ShopID, id domain.OrderID) (*domain.Order, error) {
	if err := validShop(shop); err != nil {
		return nil, err
	}
	if err := validOrder(id); err != nil {
		return nil, err
	}
	o, err := s.orders.Get(ctx, shop, id)
	if err != nil {
		return nil, fmt.Errorf("getting order %s in shop %d: %w", id, shop, err)
	}
	return o, nil
}

// Submit creates a new order.
func (s *OrderService) Submit(ctx context.Context, shop domain.ShopID, o domain.SubmitOrder) (*domain.OrderIDResult, error) {
	if err := validShop(shop); err != nil {
		return nil, err
	}
	if err := validateSubmitOrder(o); err != nil {
		return nil, err
	}
	created, err := s.orders.Submit(ctx, shop, o)
	if err != nil {
		return nil, fmt.Errorf("submitting order in shop %d: %w", shop, err)
	}
	return created, nil
}

// SubmitExpress creates a Printify Express order.
func (s *OrderService) SubmitExpress(ctx context.Context, shop domain.ShopID, o domain.ExpressOrder) ([]map[string]any, error) {
	if err := validShop(shop); err != nil {
		return nil, err
	}
	if err := validateExpressOrder(o); err != nil {
		return nil, err
	}
	data, err := s.orders.SubmitExpress(ctx, shop, o)
	if err != nil {
		return nil, fmt.Errorf("submitting express order in shop %d: %w", shop, err)
	}
	return data, nil
}

// Cancel cancels an order.
func (s *OrderService) Cancel(ctx context.Context, shop domain.ShopID, id domain.OrderID) (*domain.Order, error) {
	if err := validShop(shop); err != nil {
		return nil, err
	}
	if err := validOrder(id); err != nil {
		return nil, err
	}
	o, err := s.orders.Cancel(ctx, shop, id)
	if err != nil {
		return nil, fmt.Errorf("cancelling order %s in shop %d: %w", id, shop, err)
	}
	return o, nil
}

// SendToProduction sends an order to production.
func (s *OrderService) SendToProduction(ctx context.Context, shop domain.ShopID, id domain.OrderID) (*domain.OrderIDResult, error) {
	if err := validShop(shop); err != nil {
		return nil, err
	}
	if err := validOrder(id); err != nil {
		return nil, err
	}
	r, err := s.orders.SendToProduction(ctx, shop, id)
	if err != nil {
		return nil, fmt.Errorf("sending order %s to production in shop %d: %w", id, shop, err)
	}
	return r, nil
}

// CalculateShipping returns shipping costs for an order.
func (s *OrderService) CalculateShipping(ctx context.Context, shop domain.ShopID, o domain.SubmitOrder) (*domain.ShippingCosts, error) {
	if err := validShop(shop); err != nil {
		return nil, err
	}
	if err := validateSubmitOrder(o); err != nil {
		return nil, err
	}
	costs, err := s.orders.CalculateShipping(ctx, shop, o)
	if err != nil {
		return nil, fmt.Errorf("calculating shipping in shop %d: %w", shop, err)
	}
	return costs, nil
}

// ChangeAddress attempts to change an order's shipping address.
func (s *OrderService) ChangeAddress(ctx context.Context, shop domain.ShopID, id domain.OrderID, a domain.AddressChange) (*domain.AddressChangeResult, error) {
	if err := validShop(shop); err != nil {
		return nil, err
	}
	if err := validOrder(id); err != nil {
		return nil, err
	}
	if err := validateAddressChange(a); err != nil {
		return nil, err
	}
	result, err := s.orders.ChangeAddress(ctx, shop, id, a)
	if err != nil {
		return nil, fmt.Errorf("changing address for order %s in shop %d: %w", id, shop, err)
	}
	return result, nil
}

func validateSubmitOrder(o domain.SubmitOrder) error {
	if len(o.LineItems) == 0 {
		return fmt.Errorf("line items: at least one required, %w", domain.ErrInvalidInput)
	}
	for i, li := range o.LineItems {
		if li.Quantity < 1 {
			return fmt.Errorf("line item %d: quantity must be >= 1, %w", i, domain.ErrInvalidInput)
		}
	}
	if err := validateAddress(o.AddressTo); err != nil {
		return err
	}
	return nil
}

func validateExpressOrder(o domain.ExpressOrder) error {
	if len(o.LineItems) == 0 {
		return fmt.Errorf("line items: at least one required, %w", domain.ErrInvalidInput)
	}
	return validateAddress(o.AddressTo)
}

func validateAddress(a domain.Address) error {
	if a.FirstName == "" || a.LastName == "" {
		return fmt.Errorf("address: first and last name required, %w", domain.ErrInvalidInput)
	}
	for _, v := range []struct {
		name, value string
	}{
		{"address1", a.Address1},
		{"city", a.City},
		{"country", a.Country},
		{"zip", a.Zip},
	} {
		if v.value == "" {
			return fmt.Errorf("address: %s required, %w", v.name, domain.ErrInvalidInput)
		}
	}
	return nil
}

func validateAddressChange(a domain.AddressChange) error {
	if a.FirstName == "" || a.LastName == "" {
		return fmt.Errorf("address: first and last name required, %w", domain.ErrInvalidInput)
	}
	for _, v := range []struct {
		name, value string
	}{
		{"address1", a.Address1},
		{"city", a.City},
		{"country", a.Country},
		{"zip", a.Zip},
	} {
		if v.value == "" {
			return fmt.Errorf("address: %s required, %w", v.name, domain.ErrInvalidInput)
		}
	}
	return nil
}
