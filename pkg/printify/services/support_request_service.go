package services

import (
	"context"
	"fmt"

	"github.com/printify-go/pkg/printify/domain"
	"github.com/printify-go/pkg/printify/ports/driven"
)

// MaxDescriptionLength bounds support request descriptions.
const MaxDescriptionLength = 400

// MaxImageURLs bounds the number of images per support request.
const MaxImageURLs = 5

// SupportRequestService implements driving.SupportRequestService.
type SupportRequestService struct {
	supportRequests driven.SupportRequestGateway
}

// NewSupportRequestService builds a SupportRequestService backed by the given gateway.
func NewSupportRequestService(gw driven.SupportRequestGateway) *SupportRequestService {
	return &SupportRequestService{supportRequests: gw}
}

func validSupportRequestID(id domain.SupportRequestID) error {
	if !id.Valid() {
		return fmt.Errorf("support request id %q: %w", id, domain.ErrInvalidInput)
	}
	return nil
}

// List returns the support requests for an order.
func (s *SupportRequestService) List(ctx context.Context, shop domain.ShopID, order domain.OrderID) ([]domain.SupportRequest, error) {
	if err := validShop(shop); err != nil {
		return nil, err
	}
	if err := validOrder(order); err != nil {
		return nil, err
	}
	reqs, err := s.supportRequests.List(ctx, shop, order)
	if err != nil {
		return nil, fmt.Errorf("listing support requests for order %s in shop %d: %w", order, shop, err)
	}
	return reqs, nil
}

// Get returns one support request for an order.
func (s *SupportRequestService) Get(ctx context.Context, shop domain.ShopID, order domain.OrderID, id domain.SupportRequestID) (*domain.SupportRequest, error) {
	if err := validShop(shop); err != nil {
		return nil, err
	}
	if err := validOrder(order); err != nil {
		return nil, err
	}
	if err := validSupportRequestID(id); err != nil {
		return nil, err
	}
	req, err := s.supportRequests.Get(ctx, shop, order, id)
	if err != nil {
		return nil, fmt.Errorf("getting support request %s for order %s in shop %d: %w", id, order, shop, err)
	}
	return req, nil
}

// RequestReprint creates a reprint support request.
func (s *SupportRequestService) RequestReprint(ctx context.Context, shop domain.ShopID, order domain.OrderID, r domain.ReprintRequest) (*domain.SupportRequest, error) {
	if err := validShop(shop); err != nil {
		return nil, err
	}
	if err := validOrder(order); err != nil {
		return nil, err
	}
	if err := validateSupportRequestPayload(r.Reason, r.Description, r.ImageURLs, r.LineItems); err != nil {
		return nil, err
	}
	req, err := s.supportRequests.Reprint(ctx, shop, order, r)
	if err != nil {
		return nil, fmt.Errorf("requesting reprint for order %s in shop %d: %w", order, shop, err)
	}
	return req, nil
}

// RequestRefund creates a refund support request.
func (s *SupportRequestService) RequestRefund(ctx context.Context, shop domain.ShopID, order domain.OrderID, r domain.RefundRequest) (*domain.SupportRequest, error) {
	if err := validShop(shop); err != nil {
		return nil, err
	}
	if err := validOrder(order); err != nil {
		return nil, err
	}
	if err := validateSupportRequestPayload(r.Reason, r.Description, r.ImageURLs, r.LineItems); err != nil {
		return nil, err
	}
	req, err := s.supportRequests.Refund(ctx, shop, order, r)
	if err != nil {
		return nil, fmt.Errorf("requesting refund for order %s in shop %d: %w", order, shop, err)
	}
	return req, nil
}

func validateSupportRequestPayload(reason, description string, imageURLs []string, lineItems []domain.SupportRequestLineItem) error {
	if reason == "" {
		return fmt.Errorf("reason: %w", domain.ErrInvalidInput)
	}
	if description == "" {
		return fmt.Errorf("description: %w", domain.ErrInvalidInput)
	}
	if len(description) > MaxDescriptionLength {
		return fmt.Errorf("description: max %d characters, %w", MaxDescriptionLength, domain.ErrInvalidInput)
	}
	if len(imageURLs) > MaxImageURLs {
		return fmt.Errorf("image urls: max %d, %w", MaxImageURLs, domain.ErrInvalidInput)
	}
	if len(lineItems) == 0 {
		return fmt.Errorf("line items: at least one required, %w", domain.ErrInvalidInput)
	}
	for i, li := range lineItems {
		if li.LineItemID == "" {
			return fmt.Errorf("line item %d: line item id required, %w", i, domain.ErrInvalidInput)
		}
		if li.Quantity < 1 {
			return fmt.Errorf("line item %d: quantity must be >= 1, %w", i, domain.ErrInvalidInput)
		}
	}
	return nil
}
