package services

import (
	"context"
	"errors"
	"testing"

	"github.com/printify-go/pkg/printify/domain"
)

type stubSupportRequestGateway struct {
	listFn    func(ctx context.Context, shop domain.ShopID, order domain.OrderID) ([]domain.SupportRequest, error)
	getFn     func(ctx context.Context, shop domain.ShopID, order domain.OrderID, id domain.SupportRequestID) (*domain.SupportRequest, error)
	reprintFn func(ctx context.Context, shop domain.ShopID, order domain.OrderID, r domain.ReprintRequest) (*domain.SupportRequest, error)
	refundFn  func(ctx context.Context, shop domain.ShopID, order domain.OrderID, r domain.RefundRequest) (*domain.SupportRequest, error)
}

func (s stubSupportRequestGateway) List(ctx context.Context, shop domain.ShopID, order domain.OrderID) ([]domain.SupportRequest, error) {
	return s.listFn(ctx, shop, order)
}

func (s stubSupportRequestGateway) Get(ctx context.Context, shop domain.ShopID, order domain.OrderID, id domain.SupportRequestID) (*domain.SupportRequest, error) {
	return s.getFn(ctx, shop, order, id)
}

func (s stubSupportRequestGateway) Reprint(ctx context.Context, shop domain.ShopID, order domain.OrderID, r domain.ReprintRequest) (*domain.SupportRequest, error) {
	return s.reprintFn(ctx, shop, order, r)
}

func (s stubSupportRequestGateway) Refund(ctx context.Context, shop domain.ShopID, order domain.OrderID, r domain.RefundRequest) (*domain.SupportRequest, error) {
	return s.refundFn(ctx, shop, order, r)
}

func TestSupportRequestService_List(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		shop    domain.ShopID
		order   domain.OrderID
		gateway stubSupportRequestGateway
		want    int
		wantErr error
	}{
		{
			name:  "rejects non-positive shop",
			shop:  0,
			order: "ord-1",
			gateway: stubSupportRequestGateway{listFn: func(context.Context, domain.ShopID, domain.OrderID) ([]domain.SupportRequest, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:  "rejects empty order id",
			shop:  1,
			order: "",
			gateway: stubSupportRequestGateway{listFn: func(context.Context, domain.ShopID, domain.OrderID) ([]domain.SupportRequest, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:  "returns requests for valid ids",
			shop:  1,
			order: "ord-1",
			gateway: stubSupportRequestGateway{listFn: func(_ context.Context, shop domain.ShopID, order domain.OrderID) ([]domain.SupportRequest, error) {
				if shop != 1 || order != "ord-1" {
					t.Fatalf("gateway got shop=%d order=%q, want 1 ord-1", shop, order)
				}
				return []domain.SupportRequest{{ID: "req-1"}}, nil
			}},
			want: 1,
		},
		{
			name:  "propagates gateway error",
			shop:  1,
			order: "ord-1",
			gateway: stubSupportRequestGateway{listFn: func(context.Context, domain.ShopID, domain.OrderID) ([]domain.SupportRequest, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewSupportRequestService(tt.gateway)
			got, err := svc.List(context.Background(), tt.shop, tt.order)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if len(got) != tt.want {
				t.Fatalf("len(requests) = %d, want %d", len(got), tt.want)
			}
		})
	}
}

func TestSupportRequestService_Get(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		shop    domain.ShopID
		order   domain.OrderID
		id      domain.SupportRequestID
		gateway stubSupportRequestGateway
		wantErr error
	}{
		{
			name:  "rejects non-positive shop",
			shop:  0,
			order: "ord-1",
			id:    "req-1",
			gateway: stubSupportRequestGateway{getFn: func(context.Context, domain.ShopID, domain.OrderID, domain.SupportRequestID) (*domain.SupportRequest, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:  "rejects empty order id",
			shop:  1,
			order: "",
			id:    "req-1",
			gateway: stubSupportRequestGateway{getFn: func(context.Context, domain.ShopID, domain.OrderID, domain.SupportRequestID) (*domain.SupportRequest, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:  "rejects empty request id",
			shop:  1,
			order: "ord-1",
			id:    "",
			gateway: stubSupportRequestGateway{getFn: func(context.Context, domain.ShopID, domain.OrderID, domain.SupportRequestID) (*domain.SupportRequest, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:  "returns request for valid ids",
			shop:  1,
			order: "ord-1",
			id:    "req-1",
			gateway: stubSupportRequestGateway{getFn: func(_ context.Context, shop domain.ShopID, order domain.OrderID, id domain.SupportRequestID) (*domain.SupportRequest, error) {
				if shop != 1 || order != "ord-1" || id != "req-1" {
					t.Fatalf("gateway got shop=%d order=%q id=%q, want 1 ord-1 req-1", shop, order, id)
				}
				return &domain.SupportRequest{ID: "req-1"}, nil
			}},
		},
		{
			name:  "propagates gateway error",
			shop:  1,
			order: "ord-1",
			id:    "req-1",
			gateway: stubSupportRequestGateway{getFn: func(context.Context, domain.ShopID, domain.OrderID, domain.SupportRequestID) (*domain.SupportRequest, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewSupportRequestService(tt.gateway)
			got, err := svc.Get(context.Background(), tt.shop, tt.order, tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got == nil {
				t.Fatal("request = nil, want non-nil")
			}
		})
	}
}

func TestSupportRequestService_RequestReprint(t *testing.T) {
	gatewayErr := errors.New("boom")
	validReq := func() domain.ReprintRequest {
		return domain.ReprintRequest{
			Reason:      "misprinted",
			Description: "wrong color",
			ImageURLs:   []string{"https://example.com/img.png"},
			LineItems:   []domain.SupportRequestLineItem{{LineItemID: "li-1", Quantity: 1}},
		}
	}
	tests := []struct {
		name    string
		shop    domain.ShopID
		order   domain.OrderID
		payload domain.ReprintRequest
		gateway stubSupportRequestGateway
		wantErr error
	}{
		{
			name:    "rejects non-positive shop",
			shop:    0,
			order:   "ord-1",
			payload: validReq(),
			gateway: stubSupportRequestGateway{reprintFn: func(context.Context, domain.ShopID, domain.OrderID, domain.ReprintRequest) (*domain.SupportRequest, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:    "rejects empty order id",
			shop:    1,
			order:   "",
			payload: validReq(),
			gateway: stubSupportRequestGateway{reprintFn: func(context.Context, domain.ShopID, domain.OrderID, domain.ReprintRequest) (*domain.SupportRequest, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:  "rejects empty reason",
			shop:  1,
			order: "ord-1",
			payload: func() domain.ReprintRequest {
				r := validReq()
				r.Reason = ""
				return r
			}(),
			gateway: stubSupportRequestGateway{reprintFn: func(context.Context, domain.ShopID, domain.OrderID, domain.ReprintRequest) (*domain.SupportRequest, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:  "rejects empty description",
			shop:  1,
			order: "ord-1",
			payload: func() domain.ReprintRequest {
				r := validReq()
				r.Description = ""
				return r
			}(),
			gateway: stubSupportRequestGateway{reprintFn: func(context.Context, domain.ShopID, domain.OrderID, domain.ReprintRequest) (*domain.SupportRequest, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:  "rejects description over 400 chars",
			shop:  1,
			order: "ord-1",
			payload: func() domain.ReprintRequest {
				r := validReq()
				long := make([]byte, 401)
				for i := range long {
					long[i] = 'a'
				}
				r.Description = string(long)
				return r
			}(),
			gateway: stubSupportRequestGateway{reprintFn: func(context.Context, domain.ShopID, domain.OrderID, domain.ReprintRequest) (*domain.SupportRequest, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:  "rejects more than 5 image URLs",
			shop:  1,
			order: "ord-1",
			payload: func() domain.ReprintRequest {
				r := validReq()
				r.ImageURLs = []string{"a", "b", "c", "d", "e", "f"}
				return r
			}(),
			gateway: stubSupportRequestGateway{reprintFn: func(context.Context, domain.ShopID, domain.OrderID, domain.ReprintRequest) (*domain.SupportRequest, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:  "rejects missing line items",
			shop:  1,
			order: "ord-1",
			payload: func() domain.ReprintRequest {
				r := validReq()
				r.LineItems = nil
				return r
			}(),
			gateway: stubSupportRequestGateway{reprintFn: func(context.Context, domain.ShopID, domain.OrderID, domain.ReprintRequest) (*domain.SupportRequest, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:  "rejects line item without id",
			shop:  1,
			order: "ord-1",
			payload: func() domain.ReprintRequest {
				r := validReq()
				r.LineItems[0].LineItemID = ""
				return r
			}(),
			gateway: stubSupportRequestGateway{reprintFn: func(context.Context, domain.ShopID, domain.OrderID, domain.ReprintRequest) (*domain.SupportRequest, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:  "rejects line item with zero quantity",
			shop:  1,
			order: "ord-1",
			payload: func() domain.ReprintRequest {
				r := validReq()
				r.LineItems[0].Quantity = 0
				return r
			}(),
			gateway: stubSupportRequestGateway{reprintFn: func(context.Context, domain.ShopID, domain.OrderID, domain.ReprintRequest) (*domain.SupportRequest, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:    "creates reprint request for valid payload",
			shop:    1,
			order:   "ord-1",
			payload: validReq(),
			gateway: stubSupportRequestGateway{reprintFn: func(_ context.Context, shop domain.ShopID, order domain.OrderID, r domain.ReprintRequest) (*domain.SupportRequest, error) {
				if shop != 1 || order != "ord-1" || r.Reason != "misprinted" {
					t.Fatalf("gateway got unexpected shop=%d order=%q reason=%q", shop, order, r.Reason)
				}
				return &domain.SupportRequest{ID: "req-9", Type: "reprint"}, nil
			}},
		},
		{
			name:    "propagates gateway error",
			shop:    1,
			order:   "ord-1",
			payload: validReq(),
			gateway: stubSupportRequestGateway{reprintFn: func(context.Context, domain.ShopID, domain.OrderID, domain.ReprintRequest) (*domain.SupportRequest, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewSupportRequestService(tt.gateway)
			got, err := svc.RequestReprint(context.Background(), tt.shop, tt.order, tt.payload)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got == nil {
				t.Fatal("request = nil, want non-nil")
			}
		})
	}
}

func TestSupportRequestService_RequestRefund(t *testing.T) {
	gatewayErr := errors.New("boom")
	validReq := func() domain.RefundRequest {
		return domain.RefundRequest{
			Reason:      "damaged",
			Description: "arrived broken",
			LineItems:   []domain.SupportRequestLineItem{{LineItemID: "li-1", Quantity: 2}},
		}
	}
	tests := []struct {
		name    string
		shop    domain.ShopID
		order   domain.OrderID
		payload domain.RefundRequest
		gateway stubSupportRequestGateway
		wantErr error
	}{
		{
			name:    "rejects non-positive shop",
			shop:    0,
			order:   "ord-1",
			payload: validReq(),
			gateway: stubSupportRequestGateway{refundFn: func(context.Context, domain.ShopID, domain.OrderID, domain.RefundRequest) (*domain.SupportRequest, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:    "rejects empty order id",
			shop:    1,
			order:   "",
			payload: validReq(),
			gateway: stubSupportRequestGateway{refundFn: func(context.Context, domain.ShopID, domain.OrderID, domain.RefundRequest) (*domain.SupportRequest, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:  "rejects empty reason",
			shop:  1,
			order: "ord-1",
			payload: func() domain.RefundRequest {
				r := validReq()
				r.Reason = ""
				return r
			}(),
			gateway: stubSupportRequestGateway{refundFn: func(context.Context, domain.ShopID, domain.OrderID, domain.RefundRequest) (*domain.SupportRequest, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:  "rejects missing line items",
			shop:  1,
			order: "ord-1",
			payload: func() domain.RefundRequest {
				r := validReq()
				r.LineItems = nil
				return r
			}(),
			gateway: stubSupportRequestGateway{refundFn: func(context.Context, domain.ShopID, domain.OrderID, domain.RefundRequest) (*domain.SupportRequest, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:  "rejects line item with zero quantity",
			shop:  1,
			order: "ord-1",
			payload: func() domain.RefundRequest {
				r := validReq()
				r.LineItems[0].Quantity = 0
				return r
			}(),
			gateway: stubSupportRequestGateway{refundFn: func(context.Context, domain.ShopID, domain.OrderID, domain.RefundRequest) (*domain.SupportRequest, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:    "creates refund request for valid payload",
			shop:    1,
			order:   "ord-1",
			payload: validReq(),
			gateway: stubSupportRequestGateway{refundFn: func(_ context.Context, shop domain.ShopID, order domain.OrderID, r domain.RefundRequest) (*domain.SupportRequest, error) {
				if shop != 1 || order != "ord-1" || r.Reason != "damaged" {
					t.Fatalf("gateway got unexpected shop=%d order=%q reason=%q", shop, order, r.Reason)
				}
				return &domain.SupportRequest{ID: "req-10", Type: "refund"}, nil
			}},
		},
		{
			name:    "propagates gateway error",
			shop:    1,
			order:   "ord-1",
			payload: validReq(),
			gateway: stubSupportRequestGateway{refundFn: func(context.Context, domain.ShopID, domain.OrderID, domain.RefundRequest) (*domain.SupportRequest, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewSupportRequestService(tt.gateway)
			got, err := svc.RequestRefund(context.Background(), tt.shop, tt.order, tt.payload)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got == nil {
				t.Fatal("request = nil, want non-nil")
			}
		})
	}
}
