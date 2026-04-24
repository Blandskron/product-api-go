package usecase

import (
	"testing"

	"product-api-go/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockPurchaseRepository is a mock for the PurchaseRepository interface.
type MockPurchaseRepository struct {
	mock.Mock
}

func (m *MockPurchaseRepository) Create(purchase *domain.Purchase, tx *gorm.DB) error {
	args := m.Called(purchase, tx)
	return args.Error(0)
}

func TestProcessPurchase_Success(t *testing.T) {
	// Arrange
	mockProductRepo := new(MockProductRepository)
	mockPurchaseRepo := new(MockPurchaseRepository)
	usecase := NewPurchaseUsecase(mockProductRepo, mockPurchaseRepo, nil)

	product := domain.Product{ID: "prod-123", Name: "Test Product", Price: 20.0, Stock: 50}
	purchase := &domain.Purchase{ProductID: "prod-123", Quantity: 10}

	// Set up expectations
	mockProductRepo.On("GetByID", "prod-123").Return(product, nil)
	// We expect Update to be called with the product having its stock increased.
	mockProductRepo.On("Update", mock.MatchedBy(func(p *domain.Product) bool {
		return p.ID == "prod-123" && p.Stock == 60
	})).Return(nil)
	mockPurchaseRepo.On("Create", mock.AnythingOfType("*domain.Purchase"), (*gorm.DB)(nil)).Return(nil)

	// Act
	err := usecase.processPurchaseLogic(nil, purchase)

	// Assert
	assert.NoError(t, err)
	mockProductRepo.AssertExpectations(t)
	mockPurchaseRepo.AssertExpectations(t)
}

func TestProcessPurchase_ProductNotFound(t *testing.T) {
	// Arrange
	mockProductRepo := new(MockProductRepository)
	mockPurchaseRepo := new(MockPurchaseRepository)
	usecase := NewPurchaseUsecase(mockProductRepo, mockPurchaseRepo, nil)

	purchase := &domain.Purchase{ProductID: "prod-does-not-exist", Quantity: 10}

	// Set up expectations
	mockProductRepo.On("GetByID", "prod-does-not-exist").Return(domain.Product{}, gorm.ErrRecordNotFound)

	// Act
	err := usecase.processPurchaseLogic(nil, purchase)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrProductNotFoundForPurchase, err)
	mockProductRepo.AssertExpectations(t)
	mockPurchaseRepo.AssertNotCalled(t, "Create")
}

func TestProcessPurchase_ValidationError(t *testing.T) {
	// Arrange
	mockProductRepo := new(MockProductRepository)
	mockPurchaseRepo := new(MockPurchaseRepository)
	usecase := NewPurchaseUsecase(mockProductRepo, mockPurchaseRepo, nil)

	// This purchase has an invalid quantity, which should be caught by the validator.
	purchase := &domain.Purchase{ProductID: "prod-123", Quantity: 0}

	// Act
	err := usecase.ProcessPurchase(purchase)

	// Assert
	assert.Error(t, err)
	_, ok := err.(ValidationError)
	assert.True(t, ok, "Error should be of type ValidationError")
	mockProductRepo.AssertNotCalled(t, "GetByID")
	mockPurchaseRepo.AssertNotCalled(t, "Create")
}
