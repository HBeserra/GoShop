package ProductCatalog_test

import (
	"context"
	"errors"
	"github.com/HBeserra/GoShop/ProductCatalog"
	"github.com/HBeserra/GoShop/pkg/currency"
	"go.uber.org/mock/gomock"
	"testing"
	"time"

	"github.com/HBeserra/GoShop/domain"
	"github.com/google/uuid"
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
		userID          uuid.UUID
		product         *domain.Product
		setup           func(t setupParams)
		expectedError   error
		expectedProduct func() *domain.Product
	}

	tests := []productTest{
		{
			name:   "permission denied",
			userID: uuid.New(),
			product: &domain.Product{
				Title: "Test Product",
			},
			setup: func(t setupParams) {
				t.authService.EXPECT().HasPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
			},
			expectedError: domain.ErrUnauthorized,
		},
		{
			name:   "auth service error",
			userID: uuid.New(),
			product: &domain.Product{
				Title: "Test Product",
			},
			setup: func(t setupParams) {
				t.authService.EXPECT().HasPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, errors.New("any error"))
			},
			expectedError: domain.ErrUnauthorized,
		},
		{
			name:   "invalid product title - too short",
			userID: uuid.New(),
			product: &domain.Product{
				Title: "Short",
			},
			setup: func(t setupParams) {
				t.authService.EXPECT().HasPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
			},
			expectedError: domain.ErrInvalidProductTitle,
		},
		{
			name:   "invalid product price",
			userID: uuid.New(),
			product: &domain.Product{
				Title: "Valid Product Title",
				Price: currency.NewFromFloat(0.5),
			},
			setup: func(t setupParams) {
				t.authService.EXPECT().HasPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
			},
			expectedError: domain.ErrInvalidProductPrice,
		},
		{
			name:   "invalid variant price",
			userID: uuid.New(),
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
				t.authService.EXPECT().HasPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
			},
			expectedError: domain.ErrInvalidProductPrice,
		},
		{
			name:   "too large price difference in variants",
			userID: uuid.New(),
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
				t.authService.EXPECT().HasPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
			},
			expectedError: domain.ErrInvalidProductPrice,
		},
		{
			name:   "repository error",
			userID: uuid.New(),
			product: &domain.Product{
				Title:  "Valid Product Title",
				Price:  currency.NewFromFloat(50),
				Status: domain.ProductStatusOutOfStock,
			},
			setup: func(t setupParams) {
				t.authService.EXPECT().HasPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				t.repoService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("repository error"))
			},
			expectedError: domain.ErrFailedToCreateProduct,
		},
		{
			name:   "event bus publish error",
			userID: uuid.New(),
			product: &domain.Product{
				Title:  "Valid Product Title",
				Price:  currency.NewFromFloat(50),
				Status: domain.ProductStatusAvailable,
			},
			setup: func(t setupParams) {
				t.authService.EXPECT().HasPermission(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				t.repoService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				t.busService.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("event publish error"))
			},
			expectedError: nil,
		},
		{
			name:   "valid product creation",
			userID: uuid.New(),
			product: &domain.Product{
				Title:  "Valid Product Title",
				Price:  currency.NewFromFloat(50),
				Status: domain.ProductStatusAvailable,
			},
			setup: func(t setupParams) {
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

			authCtrl := gomock.NewController(t)
			defer authCtrl.Finish()
			repoCtrl := gomock.NewController(t)
			defer repoCtrl.Finish()
			busCtrl := gomock.NewController(t)
			defer busCtrl.Finish()

			mockAuth := NewMockAuthService(authCtrl)
			mockRepo := NewMockProductRepository(repoCtrl)
			mockBus := NewMockEventBus(busCtrl)
			service := ProductCatalog.NewProductService(mockRepo, mockBus, mockAuth)

			if tt.setup != nil {
				tt.setup(setupParams{
					busService:  mockBus,
					repoService: mockRepo,
					authService: mockAuth,
				})
			}

			prod := tt.product
			err := service.CreateProduct(context.Background(), tt.userID, prod)

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
