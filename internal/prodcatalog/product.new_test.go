package prodcatalog_test

import (
	"context"
	"errors"
	"github.com/HBeserra/GoShop/internal/prodcatalog"
	"github.com/HBeserra/GoShop/pkg/currency"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"testing"
	"time"

	"github.com/HBeserra/GoShop/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {

	type setupParams struct {
		busService  *MockEventBus
		repoService *MockProductRepository
		authService *MockAuthService
	}

	type productTest struct {
		name            string
		product         *domain.Product
		setup           func(t setupParams)
		expectedError   error
		expectedProduct func() *domain.Product
	}

	tests := []productTest{
		{
			name: "no user id",
			product: &domain.Product{
				Title: "Test Product",
			},
			setup: func(t setupParams) {
				t.authService.EXPECT().GetUserID(gomock.Any()).Return(uuid.Nil, nil)
			},
			expectedError: domain.ErrUnauthorized,
		},
		{
			name: "permission denied",
			product: &domain.Product{
				Title: "Test Product",
			},
			setup: func(t setupParams) {
				t.authService.EXPECT().GetUserID(gomock.Any()).Return(uuid.New(), nil)
				t.authService.EXPECT().HasPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
			},
			expectedError: domain.ErrUnauthorized,
		},
		{
			name: "auth service error",
			product: &domain.Product{
				Title: "Test Product",
			},
			setup: func(t setupParams) {
				t.authService.EXPECT().GetUserID(gomock.Any()).Return(uuid.New(), nil)
				t.authService.EXPECT().HasPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, errors.New("any error"))
			},
			expectedError: domain.ErrUnauthorized,
		},
		{
			name: "invalid product title - too short",
			product: &domain.Product{
				Title: "Short",
			},
			setup: func(t setupParams) {
				t.authService.EXPECT().GetUserID(gomock.Any()).Return(uuid.New(), nil)
				t.authService.EXPECT().HasPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
			},
			expectedError: domain.ErrInvalidProductTitle,
		},
		{
			name: "invalid product price",
			product: &domain.Product{
				Title: "Valid Product Title",
				Price: currency.NewFromFloat(0.5),
			},
			setup: func(t setupParams) {
				t.authService.EXPECT().GetUserID(gomock.Any()).Return(uuid.New(), nil)
				t.authService.EXPECT().HasPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
			},
			expectedError: domain.ErrInvalidProductPrice,
		},
		{
			name: "invalid variant price",
			product: &domain.Product{
				Title: "Valid Product Title",
				Price: currency.NewFromFloat(500),
				Variants: []domain.ProductVariant{
					{
						Title: "Variant 1",
						Price: currency.NewFromFloat(5),
					},
				},
			},
			setup: func(t setupParams) {
				t.authService.EXPECT().GetUserID(gomock.Any()).Return(uuid.New(), nil)
				t.authService.EXPECT().HasPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
			},
			expectedError: domain.ErrInvalidProductPrice,
		},
		{
			name: "too large price difference in variants",
			product: &domain.Product{
				Title: "Valid Product Title",
				Price: currency.NewFromFloat(50),
				Variants: []domain.ProductVariant{
					{
						Title: "Variant 1",
						Price: currency.NewFromFloat(100),
					},
					{
						Title: "Variant 2",
						Price: currency.NewFromFloat(10),
					},
				},
			},
			setup: func(t setupParams) {
				t.authService.EXPECT().GetUserID(gomock.Any()).Return(uuid.New(), nil)
				t.authService.EXPECT().HasPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
			},
			expectedError: domain.ErrInvalidProductPrice,
		},
		{
			name: "repository error",
			product: &domain.Product{
				Title:  "Valid Product Title",
				Price:  currency.NewFromFloat(50),
				Status: domain.ProductStatusOutOfStock,
			},
			setup: func(t setupParams) {
				t.authService.EXPECT().GetUserID(gomock.Any()).Return(uuid.New(), nil)
				t.authService.EXPECT().HasPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				t.repoService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("repository error"))
			},
			expectedError: domain.ErrFailedToCreateProduct,
		},
		{
			name: "event bus publish error",
			product: &domain.Product{
				Title:  "Valid Product Title",
				Price:  currency.NewFromFloat(50),
				Status: domain.ProductStatusAvailable,
			},
			setup: func(t setupParams) {
				t.authService.EXPECT().GetUserID(gomock.Any()).Return(uuid.New(), nil)
				t.authService.EXPECT().HasPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				t.repoService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				t.busService.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("event publish error"))
			},
			expectedError: nil,
		},
		{
			name: "valid product creation",
			product: &domain.Product{
				Title:  "Valid Product Title",
				Price:  currency.NewFromFloat(50),
				Status: domain.ProductStatusAvailable,
			},
			setup: func(t setupParams) {
				t.authService.EXPECT().GetUserID(gomock.Any()).Return(uuid.New(), nil)
				t.authService.EXPECT().HasPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				t.repoService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				t.busService.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
			expectedProduct: func() *domain.Product {
				return &domain.Product{
					Title:  "Valid Product Title",
					Price:  currency.NewFromFloat(50),
					Status: domain.ProductStatusAvailable,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuth := NewMockAuthService(ctrl)
			mockRepo := NewMockProductRepository(ctrl)
			mockBus := NewMockEventBus(ctrl)
			mockMedia := NewMockMediaCtrl(ctrl)
			service, err := prodcatalog.NewProductService(mockRepo, mockBus, mockAuth, mockMedia)
			assert.NoError(t, err)

			if tt.setup != nil {
				tt.setup(setupParams{
					busService:  mockBus,
					repoService: mockRepo,
					authService: mockAuth,
				})
			}

			prod := tt.product
			err = service.CreateProduct(context.Background(), prod)

			assert.ErrorIs(t, err, tt.expectedError)
			if tt.expectedError == nil && tt.expectedProduct != nil {
				expectedProd := tt.expectedProduct()
				assert.NotNil(t, prod.ID)
				assert.WithinDuration(t, time.Now(), prod.CreatedAt, time.Second)
				assert.WithinDuration(t, time.Now(), prod.UpdatedAt, time.Second)
				assert.Equal(t, expectedProd.Title, prod.Title)
				assert.Equal(t, expectedProd.Price, prod.Price)
				assert.Equal(t, expectedProd.Status, prod.Status)
			}
		})
	}
}
