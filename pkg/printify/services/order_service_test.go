package services

import (
	"context"
	"errors"
	"testing"

	"github.com/printify-go/pkg/printify/domain"
	"github.com/printify-go/pkg/printify/ports/driving"
)

type stubOrderGateway struct {
	listFn              func(ctx context.Context, shop domain.ShopID, f driving.OrderFilter) ([]domain.Order, error)
	getFn               func(ctx context.Context, shop domain.ShopID, id domain.OrderID) (*domain.Order, error)
	submitFn            func(ctx context.Context, shop domain.ShopID, o domain.SubmitOrder) (*domain.OrderIDResult, error)
	submitExpressFn     func(ctx context.Context, shop domain.ShopID, o domain.ExpressOrder) ([]map[string]any, error)
	cancelFn            func(ctx context.Context, shop domain.ShopID, id domain.OrderID) (*domain.Order, error)
	sendToProductionFn  func(ctx context.Context, shop domain.ShopID, id domain.OrderID) (*domain.OrderIDResult, error)
	calculateShippingFn func(ctx context.Context, shop domain.ShopID, o domain.SubmitOrder) (*domain.ShippingCosts, error)
	changeAddressFn     func(ctx context.Context, shop domain.ShopID, id domain.OrderID, a domain.AddressChange) (*domain.AddressChangeResult, error)
}

func (s stubOrderGateway) List(ctx context.Context, shop domain.ShopID, f driving.OrderFilter) ([]domain.Order, error) {
	return s.listFn(ctx, shop, f)
}

func (s stubOrderGateway) Get(ctx context.Context, shop domain.ShopID, id domain.OrderID) (*domain.Order, error) {
	return s.getFn(ctx, shop, id)
}

func (s stubOrderGateway) Submit(ctx context.Context, shop domain.ShopID, o domain.SubmitOrder) (*domain.OrderIDResult, error) {
	return s.submitFn(ctx, shop, o)
}

func (s stubOrderGateway) SubmitExpress(ctx context.Context, shop domain.ShopID, o domain.ExpressOrder) ([]map[string]any, error) {
	return s.submitExpressFn(ctx, shop, o)
}

func (s stubOrderGateway) Cancel(ctx context.Context, shop domain.ShopID, id domain.OrderID) (*domain.Order, error) {
	return s.cancelFn(ctx, shop, id)
}

func (s stubOrderGateway) SendToProduction(ctx context.Context, shop domain.ShopID, id domain.OrderID) (*domain.OrderIDResult, error) {
	return s.sendToProductionFn(ctx, shop, id)
}

func (s stubOrderGateway) CalculateShipping(ctx context.Context, shop domain.ShopID, o domain.SubmitOrder) (*domain.ShippingCosts, error) {
	return s.calculateShippingFn(ctx, shop, o)
}

func (s stubOrderGateway) ChangeAddress(ctx context.Context, shop domain.ShopID, id domain.OrderID, a domain.AddressChange) (*domain.AddressChangeResult, error) {
	return s.changeAddressFn(ctx, shop, id, a)
}

func validAddress() domain.Address {
	return domain.Address{
		FirstName: "Jane",
		LastName:  "Doe",
		Address1:  "1 Main St",
		City:      "Arlington",
		Country:   "US",
		Zip:       "22201",
	}
}

func validSubmitOrder() domain.SubmitOrder {
	return domain.SubmitOrder{
		ExternalID: "ext-1",
		LineItems: []domain.LineItemBlueprint{
			{PrintProviderID: 9, BlueprintID: 101, VariantID: 1, Quantity: 2},
		},
		AddressTo: validAddress(),
	}
}

func TestOrderService_List(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		shop    domain.ShopID
		gateway stubOrderGateway
		want    int
		wantErr error
	}{
		{
			name: "rejects non-positive shop without calling gateway",
			shop: 0,
			gateway: stubOrderGateway{listFn: func(context.Context, domain.ShopID, driving.OrderFilter) ([]domain.Order, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "returns orders and passes filter through",
			shop: 1,
			gateway: stubOrderGateway{listFn: func(_ context.Context, shop domain.ShopID, f driving.OrderFilter) ([]domain.Order, error) {
				if shop != 1 {
					t.Fatalf("gateway got shop %d, want 1", shop)
				}
				if f.Status != "created" {
					t.Fatalf("gateway got status %q, want %q", f.Status, "created")
				}
				return []domain.Order{{ID: "ord-1"}}, nil
			}},
			want: 1,
		},
		{
			name: "propagates gateway error",
			shop: 1,
			gateway: stubOrderGateway{listFn: func(context.Context, domain.ShopID, driving.OrderFilter) ([]domain.Order, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewOrderService(tt.gateway)
			got, err := svc.List(context.Background(), tt.shop, driving.OrderFilter{Status: "created"})
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if len(got) != tt.want {
				t.Fatalf("len(orders) = %d, want %d", len(got), tt.want)
			}
		})
	}
}

func TestOrderService_Get(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		shop    domain.ShopID
		id      domain.OrderID
		gateway stubOrderGateway
		wantErr error
	}{
		{
			name: "rejects non-positive shop",
			shop: 0,
			id:   "ord-1",
			gateway: stubOrderGateway{getFn: func(context.Context, domain.ShopID, domain.OrderID) (*domain.Order, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects empty order id",
			shop: 1,
			id:   "",
			gateway: stubOrderGateway{getFn: func(context.Context, domain.ShopID, domain.OrderID) (*domain.Order, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "returns order for valid ids",
			shop: 1,
			id:   "ord-1",
			gateway: stubOrderGateway{getFn: func(_ context.Context, shop domain.ShopID, id domain.OrderID) (*domain.Order, error) {
				if shop != 1 || id != "ord-1" {
					t.Fatalf("gateway got shop=%d id=%q, want 1 ord-1", shop, id)
				}
				return &domain.Order{ID: "ord-1"}, nil
			}},
		},
		{
			name: "propagates gateway error",
			shop: 1,
			id:   "ord-1",
			gateway: stubOrderGateway{getFn: func(context.Context, domain.ShopID, domain.OrderID) (*domain.Order, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewOrderService(tt.gateway)
			got, err := svc.Get(context.Background(), tt.shop, tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got == nil {
				t.Fatal("order = nil, want non-nil")
			}
		})
	}
}

func TestOrderService_Submit(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		shop    domain.ShopID
		payload domain.SubmitOrder
		gateway stubOrderGateway
		wantErr error
	}{
		{
			name:    "rejects non-positive shop",
			shop:    0,
			payload: validSubmitOrder(),
			gateway: stubOrderGateway{submitFn: func(context.Context, domain.ShopID, domain.SubmitOrder) (*domain.OrderIDResult, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects missing line items",
			shop: 1,
			payload: func() domain.SubmitOrder {
				o := validSubmitOrder()
				o.LineItems = nil
				return o
			}(),
			gateway: stubOrderGateway{submitFn: func(context.Context, domain.ShopID, domain.SubmitOrder) (*domain.OrderIDResult, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects zero quantity line item",
			shop: 1,
			payload: func() domain.SubmitOrder {
				o := validSubmitOrder()
				o.LineItems[0].Quantity = 0
				return o
			}(),
			gateway: stubOrderGateway{submitFn: func(context.Context, domain.ShopID, domain.SubmitOrder) (*domain.OrderIDResult, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects missing last name in address",
			shop: 1,
			payload: func() domain.SubmitOrder {
				o := validSubmitOrder()
				o.AddressTo.LastName = ""
				return o
			}(),
			gateway: stubOrderGateway{submitFn: func(context.Context, domain.ShopID, domain.SubmitOrder) (*domain.OrderIDResult, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects missing city in address",
			shop: 1,
			payload: func() domain.SubmitOrder {
				o := validSubmitOrder()
				o.AddressTo.City = ""
				return o
			}(),
			gateway: stubOrderGateway{submitFn: func(context.Context, domain.ShopID, domain.SubmitOrder) (*domain.OrderIDResult, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "submits valid order",
			shop: 1,
			payload: func() domain.SubmitOrder {
				o := validSubmitOrder()
				return o
			}(),
			gateway: stubOrderGateway{submitFn: func(_ context.Context, shop domain.ShopID, o domain.SubmitOrder) (*domain.OrderIDResult, error) {
				if shop != 1 || o.ExternalID != "ext-1" || len(o.LineItems) != 1 {
					t.Fatalf("gateway got unexpected shop=%d external=%q lineItems=%d", shop, o.ExternalID, len(o.LineItems))
				}
				return &domain.OrderIDResult{ID: "ord-9"}, nil
			}},
		},
		{
			name:    "propagates gateway error",
			shop:    1,
			payload: validSubmitOrder(),
			gateway: stubOrderGateway{submitFn: func(context.Context, domain.ShopID, domain.SubmitOrder) (*domain.OrderIDResult, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewOrderService(tt.gateway)
			got, err := svc.Submit(context.Background(), tt.shop, tt.payload)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got == nil {
				t.Fatal("result = nil, want non-nil")
			}
			if tt.wantErr == nil && got != nil && got.ID != "ord-9" {
				t.Fatalf("result.ID = %q, want ord-9", got.ID)
			}
		})
	}
}

func TestOrderService_SubmitExpress(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		shop    domain.ShopID
		payload domain.ExpressOrder
		gateway stubOrderGateway
		wantErr error
	}{
		{
			name:    "rejects non-positive shop",
			shop:    0,
			payload: validExpressOrder(),
			gateway: stubOrderGateway{submitExpressFn: func(context.Context, domain.ShopID, domain.ExpressOrder) ([]map[string]any, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects missing line items",
			shop: 1,
			payload: func() domain.ExpressOrder {
				o := validExpressOrder()
				o.LineItems = nil
				return o
			}(),
			gateway: stubOrderGateway{submitExpressFn: func(context.Context, domain.ShopID, domain.ExpressOrder) ([]map[string]any, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects missing country in address",
			shop: 1,
			payload: func() domain.ExpressOrder {
				o := validExpressOrder()
				o.AddressTo.Country = ""
				return o
			}(),
			gateway: stubOrderGateway{submitExpressFn: func(context.Context, domain.ShopID, domain.ExpressOrder) ([]map[string]any, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "submits valid express order",
			shop: 1,
			payload: func() domain.ExpressOrder {
				o := validExpressOrder()
				return o
			}(),
			gateway: stubOrderGateway{submitExpressFn: func(_ context.Context, shop domain.ShopID, o domain.ExpressOrder) ([]map[string]any, error) {
				if shop != 1 || len(o.LineItems) != 1 {
					t.Fatalf("gateway got unexpected shop=%d lineItems=%d", shop, len(o.LineItems))
				}
				return []map[string]any{{"id": "exp-1"}}, nil
			}},
		},
		{
			name:    "propagates gateway error",
			shop:    1,
			payload: validExpressOrder(),
			gateway: stubOrderGateway{submitExpressFn: func(context.Context, domain.ShopID, domain.ExpressOrder) ([]map[string]any, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewOrderService(tt.gateway)
			got, err := svc.SubmitExpress(context.Background(), tt.shop, tt.payload)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && len(got) != 1 {
				t.Fatalf("len(result) = %d, want 1", len(got))
			}
		})
	}
}

func validExpressOrder() domain.ExpressOrder {
	return domain.ExpressOrder{
		ExternalID: "ext-1",
		LineItems:  []domain.LineItemBlueprint{{PrintProviderID: 9, BlueprintID: 101, VariantID: 1, Quantity: 1}},
		AddressTo:  validAddress(),
	}
}

func TestOrderService_Cancel(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		shop    domain.ShopID
		id      domain.OrderID
		gateway stubOrderGateway
		wantErr error
	}{
		{
			name: "rejects non-positive shop",
			shop: 0,
			id:   "ord-1",
			gateway: stubOrderGateway{cancelFn: func(context.Context, domain.ShopID, domain.OrderID) (*domain.Order, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects empty order id",
			shop: 1,
			id:   "",
			gateway: stubOrderGateway{cancelFn: func(context.Context, domain.ShopID, domain.OrderID) (*domain.Order, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "cancels valid order",
			shop: 1,
			id:   "ord-1",
			gateway: stubOrderGateway{cancelFn: func(_ context.Context, shop domain.ShopID, id domain.OrderID) (*domain.Order, error) {
				if shop != 1 || id != "ord-1" {
					t.Fatalf("gateway got shop=%d id=%q, want 1 ord-1", shop, id)
				}
				return &domain.Order{ID: "ord-1"}, nil
			}},
		},
		{
			name: "propagates gateway error",
			shop: 1,
			id:   "ord-1",
			gateway: stubOrderGateway{cancelFn: func(context.Context, domain.ShopID, domain.OrderID) (*domain.Order, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewOrderService(tt.gateway)
			got, err := svc.Cancel(context.Background(), tt.shop, tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got == nil {
				t.Fatal("order = nil, want non-nil")
			}
		})
	}
}

func TestOrderService_SendToProduction(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		shop    domain.ShopID
		id      domain.OrderID
		gateway stubOrderGateway
		wantErr error
	}{
		{
			name: "rejects non-positive shop",
			shop: 0,
			id:   "ord-1",
			gateway: stubOrderGateway{sendToProductionFn: func(context.Context, domain.ShopID, domain.OrderID) (*domain.OrderIDResult, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects empty order id",
			shop: 1,
			id:   "",
			gateway: stubOrderGateway{sendToProductionFn: func(context.Context, domain.ShopID, domain.OrderID) (*domain.OrderIDResult, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "sends valid order to production",
			shop: 1,
			id:   "ord-1",
			gateway: stubOrderGateway{sendToProductionFn: func(_ context.Context, shop domain.ShopID, id domain.OrderID) (*domain.OrderIDResult, error) {
				if shop != 1 || id != "ord-1" {
					t.Fatalf("gateway got shop=%d id=%q, want 1 ord-1", shop, id)
				}
				return &domain.OrderIDResult{ID: "ord-1"}, nil
			}},
		},
		{
			name: "propagates gateway error",
			shop: 1,
			id:   "ord-1",
			gateway: stubOrderGateway{sendToProductionFn: func(context.Context, domain.ShopID, domain.OrderID) (*domain.OrderIDResult, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewOrderService(tt.gateway)
			got, err := svc.SendToProduction(context.Background(), tt.shop, tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got == nil {
				t.Fatal("result = nil, want non-nil")
			}
		})
	}
}

func TestOrderService_CalculateShipping(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		shop    domain.ShopID
		payload domain.SubmitOrder
		gateway stubOrderGateway
		wantErr error
	}{
		{
			name:    "rejects non-positive shop",
			shop:    0,
			payload: validSubmitOrder(),
			gateway: stubOrderGateway{calculateShippingFn: func(context.Context, domain.ShopID, domain.SubmitOrder) (*domain.ShippingCosts, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects missing line items",
			shop: 1,
			payload: func() domain.SubmitOrder {
				o := validSubmitOrder()
				o.LineItems = nil
				return o
			}(),
			gateway: stubOrderGateway{calculateShippingFn: func(context.Context, domain.ShopID, domain.SubmitOrder) (*domain.ShippingCosts, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "calculates shipping for valid order",
			shop: 1,
			payload: func() domain.SubmitOrder {
				o := validSubmitOrder()
				return o
			}(),
			gateway: stubOrderGateway{calculateShippingFn: func(_ context.Context, shop domain.ShopID, o domain.SubmitOrder) (*domain.ShippingCosts, error) {
				if shop != 1 {
					t.Fatalf("gateway got shop %d, want 1", shop)
				}
				if o.AddressTo.Zip != "22201" {
					t.Fatalf("gateway got zip %q, want 22201", o.AddressTo.Zip)
				}
				return &domain.ShippingCosts{}, nil
			}},
		},
		{
			name:    "propagates gateway error",
			shop:    1,
			payload: validSubmitOrder(),
			gateway: stubOrderGateway{calculateShippingFn: func(context.Context, domain.ShopID, domain.SubmitOrder) (*domain.ShippingCosts, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewOrderService(tt.gateway)
			got, err := svc.CalculateShipping(context.Background(), tt.shop, tt.payload)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got == nil {
				t.Fatal("costs = nil, want non-nil")
			}
		})
	}
}

func TestOrderService_ChangeAddress(t *testing.T) {
	gatewayErr := errors.New("boom")
	validChange := domain.AddressChange{
		FirstName: "New",
		LastName:  "Name",
		Address1:  "2 Main St",
		City:      "Arlington",
		Country:   "US",
		Zip:       "22201",
	}
	tests := []struct {
		name    string
		shop    domain.ShopID
		id      domain.OrderID
		payload domain.AddressChange
		gateway stubOrderGateway
		wantErr error
	}{
		{
			name:    "rejects non-positive shop",
			shop:    0,
			id:      "ord-1",
			payload: validChange,
			gateway: stubOrderGateway{changeAddressFn: func(context.Context, domain.ShopID, domain.OrderID, domain.AddressChange) (*domain.AddressChangeResult, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:    "rejects empty order id",
			shop:    1,
			id:      "",
			payload: validChange,
			gateway: stubOrderGateway{changeAddressFn: func(context.Context, domain.ShopID, domain.OrderID, domain.AddressChange) (*domain.AddressChangeResult, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects missing first name",
			shop: 1,
			id:   "ord-1",
			payload: func() domain.AddressChange {
				a := validChange
				a.FirstName = ""
				return a
			}(),
			gateway: stubOrderGateway{changeAddressFn: func(context.Context, domain.ShopID, domain.OrderID, domain.AddressChange) (*domain.AddressChangeResult, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects missing address1",
			shop: 1,
			id:   "ord-1",
			payload: func() domain.AddressChange {
				a := validChange
				a.Address1 = ""
				return a
			}(),
			gateway: stubOrderGateway{changeAddressFn: func(context.Context, domain.ShopID, domain.OrderID, domain.AddressChange) (*domain.AddressChangeResult, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects missing zip",
			shop: 1,
			id:   "ord-1",
			payload: func() domain.AddressChange {
				a := validChange
				a.Zip = ""
				return a
			}(),
			gateway: stubOrderGateway{changeAddressFn: func(context.Context, domain.ShopID, domain.OrderID, domain.AddressChange) (*domain.AddressChangeResult, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:    "changes address for valid payload",
			shop:    1,
			id:      "ord-1",
			payload: validChange,
			gateway: stubOrderGateway{changeAddressFn: func(_ context.Context, shop domain.ShopID, id domain.OrderID, a domain.AddressChange) (*domain.AddressChangeResult, error) {
				if shop != 1 || id != "ord-1" || a.Address1 != "2 Main St" {
					t.Fatalf("gateway got shop=%d id=%q address1=%q", shop, id, a.Address1)
				}
				return &domain.AddressChangeResult{Resolution: "updated_directly", OrderID: "ord-1"}, nil
			}},
		},
		{
			name:    "propagates gateway error",
			shop:    1,
			id:      "ord-1",
			payload: validChange,
			gateway: stubOrderGateway{changeAddressFn: func(context.Context, domain.ShopID, domain.OrderID, domain.AddressChange) (*domain.AddressChangeResult, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewOrderService(tt.gateway)
			got, err := svc.ChangeAddress(context.Background(), tt.shop, tt.id, tt.payload)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got == nil {
				t.Fatal("result = nil, want non-nil")
			}
		})
	}
}
