package services

import (
	"context"
	"errors"
	"testing"

	"github.com/printify-go/pkg/printify/domain"
	"github.com/printify-go/pkg/printify/ports/driving"
)

type stubUploadGateway struct {
	listFn    func(ctx context.Context, f driving.UploadFilter) ([]domain.Upload, error)
	getFn     func(ctx context.Context, id domain.UploadID) (*domain.Upload, error)
	uploadFn  func(ctx context.Context, img domain.UploadImage) (*domain.Upload, error)
	archiveFn func(ctx context.Context, id domain.UploadID) error
}

func (s stubUploadGateway) List(ctx context.Context, f driving.UploadFilter) ([]domain.Upload, error) {
	return s.listFn(ctx, f)
}

func (s stubUploadGateway) Get(ctx context.Context, id domain.UploadID) (*domain.Upload, error) {
	return s.getFn(ctx, id)
}

func (s stubUploadGateway) UploadImage(ctx context.Context, img domain.UploadImage) (*domain.Upload, error) {
	return s.uploadFn(ctx, img)
}

func (s stubUploadGateway) Archive(ctx context.Context, id domain.UploadID) error {
	return s.archiveFn(ctx, id)
}

func TestUploadService_List(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		gateway stubUploadGateway
		want    int
		wantErr error
	}{
		{
			name: "returns uploads and passes filter through",
			gateway: stubUploadGateway{listFn: func(_ context.Context, f driving.UploadFilter) ([]domain.Upload, error) {
				if f.Limit == nil || *f.Limit != 10 {
					t.Fatalf("gateway got limit %+v, want 10", f.Limit)
				}
				return []domain.Upload{{ID: "up-1"}}, nil
			}},
			want: 1,
		},
		{
			name: "propagates gateway error",
			gateway: stubUploadGateway{listFn: func(context.Context, driving.UploadFilter) ([]domain.Upload, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewUploadService(tt.gateway)
			limit := 10
			got, err := svc.List(context.Background(), driving.UploadFilter{Limit: &limit})
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if len(got) != tt.want {
				t.Fatalf("len(uploads) = %d, want %d", len(got), tt.want)
			}
		})
	}
}

func TestUploadService_Get(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		id      domain.UploadID
		gateway stubUploadGateway
		wantErr error
	}{
		{
			name: "rejects empty upload id",
			id:   "",
			gateway: stubUploadGateway{getFn: func(context.Context, domain.UploadID) (*domain.Upload, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "returns upload for valid id",
			id:   "up-1",
			gateway: stubUploadGateway{getFn: func(_ context.Context, id domain.UploadID) (*domain.Upload, error) {
				if id != "up-1" {
					t.Fatalf("gateway got id %q, want up-1", id)
				}
				return &domain.Upload{ID: "up-1"}, nil
			}},
		},
		{
			name: "propagates gateway error",
			id:   "up-1",
			gateway: stubUploadGateway{getFn: func(context.Context, domain.UploadID) (*domain.Upload, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewUploadService(tt.gateway)
			got, err := svc.Get(context.Background(), tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got == nil {
				t.Fatal("upload = nil, want non-nil")
			}
		})
	}
}

func TestUploadService_UploadImage(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		payload domain.UploadImage
		gateway stubUploadGateway
		wantErr error
	}{
		{
			name: "rejects empty file name",
			payload: domain.UploadImage{
				URL: "https://example.com/img.png",
			},
			gateway: stubUploadGateway{uploadFn: func(context.Context, domain.UploadImage) (*domain.Upload, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects both url and contents empty",
			payload: domain.UploadImage{
				FileName: "img.png",
			},
			gateway: stubUploadGateway{uploadFn: func(context.Context, domain.UploadImage) (*domain.Upload, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "rejects both url and contents set",
			payload: domain.UploadImage{
				FileName: "img.png",
				URL:      "https://example.com/img.png",
				Contents: "base64data",
			},
			gateway: stubUploadGateway{uploadFn: func(context.Context, domain.UploadImage) (*domain.Upload, error) {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "uploads image via url",
			payload: domain.UploadImage{
				FileName: "img.png",
				URL:      "https://example.com/img.png",
			},
			gateway: stubUploadGateway{uploadFn: func(_ context.Context, img domain.UploadImage) (*domain.Upload, error) {
				if img.FileName != "img.png" || img.URL != "https://example.com/img.png" {
					t.Fatalf("gateway got unexpected image %+v", img)
				}
				return &domain.Upload{ID: "up-2", FileName: "img.png"}, nil
			}},
		},
		{
			name: "uploads image via contents",
			payload: domain.UploadImage{
				FileName: "img.png",
				Contents: "base64data",
			},
			gateway: stubUploadGateway{uploadFn: func(_ context.Context, img domain.UploadImage) (*domain.Upload, error) {
				if img.FileName != "img.png" || img.Contents != "base64data" {
					t.Fatalf("gateway got unexpected image %+v", img)
				}
				return &domain.Upload{ID: "up-3", FileName: "img.png"}, nil
			}},
		},
		{
			name:    "propagates gateway error",
			payload: domain.UploadImage{FileName: "img.png", URL: "https://example.com/img.png"},
			gateway: stubUploadGateway{uploadFn: func(context.Context, domain.UploadImage) (*domain.Upload, error) {
				return nil, gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewUploadService(tt.gateway)
			got, err := svc.UploadImage(context.Background(), tt.payload)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got == nil {
				t.Fatal("upload = nil, want non-nil")
			}
		})
	}
}

func TestUploadService_Archive(t *testing.T) {
	gatewayErr := errors.New("boom")
	tests := []struct {
		name    string
		id      domain.UploadID
		gateway stubUploadGateway
		wantErr error
	}{
		{
			name: "rejects empty upload id",
			id:   "",
			gateway: stubUploadGateway{archiveFn: func(context.Context, domain.UploadID) error {
				panic("gateway must not be called")
			}},
			wantErr: domain.ErrInvalidInput,
		},
		{
			name: "archives upload for valid id",
			id:   "up-1",
			gateway: stubUploadGateway{archiveFn: func(_ context.Context, id domain.UploadID) error {
				if id != "up-1" {
					t.Fatalf("gateway got id %q, want up-1", id)
				}
				return nil
			}},
		},
		{
			name: "propagates gateway error",
			id:   "up-1",
			gateway: stubUploadGateway{archiveFn: func(context.Context, domain.UploadID) error {
				return gatewayErr
			}},
			wantErr: gatewayErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewUploadService(tt.gateway)
			err := svc.Archive(context.Background(), tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
