package prodcatalog_test

import (
	"context"
	"errors"
	"github.com/HBeserra/GoShop/domain"
	"github.com/HBeserra/GoShop/internal/prodcatalog"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestUpdateProduct(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepo := NewMockProductRepository(mockCtrl)
	mockAuth := NewMockAuthService(mockCtrl)
	mockBus := NewMockEventBus(mockCtrl)
	mockMedia := NewMockMediaCtrl(mockCtrl)

	service, _ := prodcatalog.NewProductService(mockRepo, mockBus, mockAuth, mockMedia)

	validProduct := &domain.Product{
		ID:     uuid.New(),
		Title:  "Test Product",
		Price:  500,
		Stock:  50,
		Status: domain.ProductStatusOutOfStock,
	}

	tests := []struct {
		name          string
		setupMocks    func()
		product       *domain.Product
		expectedError error
	}{
		{
			name: "successful update",
			setupMocks: func() {
				mockAuth.EXPECT().GetUserID(gomock.Any()).Return(uuid.New(), nil)
				mockAuth.EXPECT().CheckPermissions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				mockRepo.EXPECT().Update(gomock.Any(), validProduct).Return(nil)
			},
			product:       validProduct,
			expectedError: nil,
		},
		{
			name: "user unauthorized",
			setupMocks: func() {
				mockAuth.EXPECT().GetUserID(gomock.Any()).Return(uuid.New(), nil)
				mockAuth.EXPECT().CheckPermissions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
			},
			product:       validProduct,
			expectedError: domain.ErrUnauthorized,
		},
		{
			name: "auth service error",
			setupMocks: func() {
				mockAuth.EXPECT().GetUserID(gomock.Any()).Return(uuid.Nil, errors.New("auth error"))
			},
			product:       validProduct,
			expectedError: domain.ErrUnauthorized,
		},
		{
			name: "repository update error",
			setupMocks: func() {
				mockAuth.EXPECT().GetUserID(gomock.Any()).Return(uuid.New(), nil)
				mockAuth.EXPECT().CheckPermissions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				mockRepo.EXPECT().Update(gomock.Any(), validProduct).Return(domain.ErrProductNotFound)
			},
			product:       validProduct,
			expectedError: domain.ErrProductNotFound,
		},
		{
			name: "validation error",
			setupMocks: func() {
				mockAuth.EXPECT().GetUserID(gomock.Any()).Return(uuid.New(), nil)
				mockAuth.EXPECT().CheckPermissions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
			},
			product:       &domain.Product{ID: uuid.New(), Title: "Produto 12", Price: -100}, // invalid price
			expectedError: domain.ErrInvalidProductPrice,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := service.UpdateProduct(context.Background(), "", tt.product)

			assert.ErrorIs(t, err, tt.expectedError)
		})
	}
}
