package services

import (
	"context"
	"errors"
	"testing"

	"github.com/printify-go/pkg/printify/domain"
)

type stubPersonalizationGateway struct {
	listOptionsFn   func(ctx context.Context, shop domain.ShopID, product domain.ProductID) ([]domain.PersonalizationOption, error)
	createConfigFn  func(ctx context.Context, shop domain.ShopID, product domain.ProductID, cfg domain.CreatePersonalizationConfig) (*domain.PersonalizationConfig, error)
	createPreviewFn func(ctx context.Context, shop domain.ShopID, product domain.ProductID, task domain.CreatePreviewTask) (*domain.PreviewTask, error)
	getPreviewFn    func(ctx context.Context, shop domain.ShopID, product domain.ProductID, id domain.TaskID) (*domain.PreviewTask, error)
}

func (s stubPersonalizationGateway) ListOptions(ctx context.Context, shop domain.ShopID, product domain.ProductID) ([]domain.PersonalizationOption, error) {
	return s.listOptionsFn(ctx, shop, product)
}

func (s stubPersonalizationGateway) CreateConfig(ctx context.Context, shop domain.ShopID, product domain.ProductID, cfg domain.CreatePersonalizationConfig) (*domain.PersonalizationConfig, error) {
	return s.createConfigFn(ctx, shop, product, cfg)
}

func (s stubPersonalizationGateway) CreatePreviewTask(ctx context.Context, shop domain.ShopID, product domain.ProductID, task domain.CreatePreviewTask) (*domain.PreviewTask, error) {
	return s.createPreviewFn(ctx, shop, product, task)
}

func (s stubPersonalizationGateway) GetPreviewTask(ctx context.Context, shop domain.ShopID, product domain.ProductID, id domain.TaskID) (*domain.PreviewTask, error) {
	return s.getPreviewFn(ctx, shop, product, id)
}

func TestPersonalizationService_ListOptions(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		shop    domain.ShopID
		product domain.ProductID
		gateway stubPersonalizationGateway
		want    int
		wantErr error
	}{
		{
			name:    "rejects non-positive shop",
			shop:    0,
			product: "abc",
			gateway: stubPersonalizationGateway{listOptionsFn: func(context.Context, domain.ShopID, domain.ProductID) ([]domain.PersonalizationOption, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:    "rejects empty product id",
			shop:    1,
			product: "",
			gateway: stubPersonalizationGateway{listOptionsFn: func(context.Context, domain.ShopID, domain.ProductID) ([]domain.PersonalizationOption, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:    "returns options",
			shop:    1,
			product: "abc",
			gateway: stubPersonalizationGateway{listOptionsFn: func(_ context.Context, shop domain.ShopID, product domain.ProductID) ([]domain.PersonalizationOption, error) {
				if shop != 1 || product != "abc" {
					t.Fatalf("gateway got shop=%d product=%q, want 1 abc", shop, product)
				}
				return []domain.PersonalizationOption{{FieldID: "opt-1"}}, nil
			}},
			want: 1,
		},
		{
			name:    "propagates gateway error",
			shop:    1,
			product: "abc",
			gateway: stubPersonalizationGateway{listOptionsFn: func(context.Context, domain.ShopID, domain.ProductID) ([]domain.PersonalizationOption, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewPersonalizationService(tt.gateway)
			got, err := svc.ListOptions(context.Background(), tt.shop, tt.product)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if len(got) != tt.want {
				t.Fatalf("len(options) = %d, want %d", len(got), tt.want)
			}
		})
	}
}

func TestPersonalizationService_CreateConfig(t *testing.T) {
	gatewayErr := errors.New("boom")
	validConfig := func() domain.CreatePersonalizationConfig {
		return domain.CreatePersonalizationConfig{
			VariantID: 1,
			Items:     []domain.PersonalizationConfigItem{{FieldID: "text-1"}},
		}
	}
	tests := []struct {
		name    string
		shop    domain.ShopID
		product domain.ProductID
		cfg     domain.CreatePersonalizationConfig
		gateway stubPersonalizationGateway
		wantErr error
	}{
		{
			name:    "rejects non-positive shop",
			shop:    0,
			product: "abc",
			cfg:     validConfig(),
			gateway: stubPersonalizationGateway{createConfigFn: func(context.Context, domain.ShopID, domain.ProductID, domain.CreatePersonalizationConfig) (*domain.PersonalizationConfig, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:    "rejects empty product id",
			shop:    1,
			product: "",
			cfg:     validConfig(),
			gateway: stubPersonalizationGateway{createConfigFn: func(context.Context, domain.ShopID, domain.ProductID, domain.CreatePersonalizationConfig) (*domain.PersonalizationConfig, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:    "rejects invalid variant id",
			shop:    1,
			product: "abc",
			cfg: func() domain.CreatePersonalizationConfig {
				c := validConfig()
				c.VariantID = 0
				return c
			}(),
			gateway: stubPersonalizationGateway{createConfigFn: func(context.Context, domain.ShopID, domain.ProductID, domain.CreatePersonalizationConfig) (*domain.PersonalizationConfig, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:    "rejects missing items",
			shop:    1,
			product: "abc",
			cfg: func() domain.CreatePersonalizationConfig {
				c := validConfig()
				c.Items = nil
				return c
			}(),
			gateway: stubPersonalizationGateway{createConfigFn: func(context.Context, domain.ShopID, domain.ProductID, domain.CreatePersonalizationConfig) (*domain.PersonalizationConfig, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:    "rejects item without field id",
			shop:    1,
			product: "abc",
			cfg: func() domain.CreatePersonalizationConfig {
				c := validConfig()
				c.Items[0].FieldID = ""
				return c
			}(),
			gateway: stubPersonalizationGateway{createConfigFn: func(context.Context, domain.ShopID, domain.ProductID, domain.CreatePersonalizationConfig) (*domain.PersonalizationConfig, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:    "creates config for valid payload",
			shop:    1,
			product: "abc",
			cfg:     validConfig(),
			gateway: stubPersonalizationGateway{createConfigFn: func(_ context.Context, shop domain.ShopID, product domain.ProductID, cfg domain.CreatePersonalizationConfig) (*domain.PersonalizationConfig, error) {
				if shop != 1 || product != "abc" || cfg.VariantID != 1 || cfg.Items[0].FieldID != "text-1" {
					t.Fatalf("gateway got unexpected shop=%d product=%q cfg=%+v", shop, product, cfg)
				}
				return &domain.PersonalizationConfig{Strategy: "custom", Instructions: "print it"}, nil
			}},
		},
		{
			name:    "propagates gateway error",
			shop:    1,
			product: "abc",
			cfg:     validConfig(),
			gateway: stubPersonalizationGateway{createConfigFn: func(context.Context, domain.ShopID, domain.ProductID, domain.CreatePersonalizationConfig) (*domain.PersonalizationConfig, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewPersonalizationService(tt.gateway)
			got, err := svc.CreateConfig(context.Background(), tt.shop, tt.product, tt.cfg)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got == nil {
				t.Fatal("config = nil, want non-nil")
			}
		})
	}
}

func TestPersonalizationService_CreatePreviewTask(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		shop    domain.ShopID
		product domain.ProductID
		task    domain.CreatePreviewTask
		gateway stubPersonalizationGateway
		wantErr error
	}{
		{
			name:    "rejects non-positive shop",
			shop:    0,
			product: "abc",
			task:    domain.CreatePreviewTask{VariantIDs: []domain.VariantID{1}},
			gateway: stubPersonalizationGateway{createPreviewFn: func(context.Context, domain.ShopID, domain.ProductID, domain.CreatePreviewTask) (*domain.PreviewTask, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:    "rejects empty product id",
			shop:    1,
			product: "",
			task:    domain.CreatePreviewTask{VariantIDs: []domain.VariantID{1}},
			gateway: stubPersonalizationGateway{createPreviewFn: func(context.Context, domain.ShopID, domain.ProductID, domain.CreatePreviewTask) (*domain.PreviewTask, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:    "rejects missing variant ids",
			shop:    1,
			product: "abc",
			task:    domain.CreatePreviewTask{VariantIDs: nil},
			gateway: stubPersonalizationGateway{createPreviewFn: func(context.Context, domain.ShopID, domain.ProductID, domain.CreatePreviewTask) (*domain.PreviewTask, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:    "creates preview task for valid payload",
			shop:    1,
			product: "abc",
			task:    domain.CreatePreviewTask{VariantIDs: []domain.VariantID{1, 2}},
			gateway: stubPersonalizationGateway{createPreviewFn: func(_ context.Context, shop domain.ShopID, product domain.ProductID, task domain.CreatePreviewTask) (*domain.PreviewTask, error) {
				if shop != 1 || product != "abc" || len(task.VariantIDs) != 2 {
					t.Fatalf("gateway got unexpected shop=%d product=%q variants=%d", shop, product, len(task.VariantIDs))
				}
				return &domain.PreviewTask{TaskID: "task-1"}, nil
			}},
		},
		{
			name:    "propagates gateway error",
			shop:    1,
			product: "abc",
			task:    domain.CreatePreviewTask{VariantIDs: []domain.VariantID{1}},
			gateway: stubPersonalizationGateway{createPreviewFn: func(context.Context, domain.ShopID, domain.ProductID, domain.CreatePreviewTask) (*domain.PreviewTask, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewPersonalizationService(tt.gateway)
			got, err := svc.CreatePreviewTask(context.Background(), tt.shop, tt.product, tt.task)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got == nil {
				t.Fatal("task = nil, want non-nil")
			}
		})
	}
}

func TestPersonalizationService_GetPreviewTask(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		shop    domain.ShopID
		product domain.ProductID
		id      domain.TaskID
		gateway stubPersonalizationGateway
		wantErr error
	}{
		{
			name:    "rejects non-positive shop",
			shop:    0,
			product: "abc",
			id:      "task-1",
			gateway: stubPersonalizationGateway{getPreviewFn: func(context.Context, domain.ShopID, domain.ProductID, domain.TaskID) (*domain.PreviewTask, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:    "rejects empty product id",
			shop:    1,
			product: "",
			id:      "task-1",
			gateway: stubPersonalizationGateway{getPreviewFn: func(context.Context, domain.ShopID, domain.ProductID, domain.TaskID) (*domain.PreviewTask, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:    "rejects empty task id",
			shop:    1,
			product: "abc",
			id:      "",
			gateway: stubPersonalizationGateway{getPreviewFn: func(context.Context, domain.ShopID, domain.ProductID, domain.TaskID) (*domain.PreviewTask, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name:    "returns preview task for valid ids",
			shop:    1,
			product: "abc",
			id:      "task-1",
			gateway: stubPersonalizationGateway{getPreviewFn: func(_ context.Context, shop domain.ShopID, product domain.ProductID, id domain.TaskID) (*domain.PreviewTask, error) {
				if shop != 1 || product != "abc" || id != "task-1" {
					t.Fatalf("gateway got shop=%d product=%q id=%q, want 1 abc task-1", shop, product, id)
				}
				return &domain.PreviewTask{TaskID: "task-1"}, nil
			}},
		},
		{
			name:    "propagates gateway error",
			shop:    1,
			product: "abc",
			id:      "task-1",
			gateway: stubPersonalizationGateway{getPreviewFn: func(context.Context, domain.ShopID, domain.ProductID, domain.TaskID) (*domain.PreviewTask, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewPersonalizationService(tt.gateway)
			got, err := svc.GetPreviewTask(context.Background(), tt.shop, tt.product, tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got == nil {
				t.Fatal("task = nil, want non-nil")
			}
		})
	}
}
