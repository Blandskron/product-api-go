package usecase

import (
	"errors"
	"testing"

	"product-api-go/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockSaleRepository is a mock for the SaleRepository interface.
type MockSaleRepository struct {
	mock.Mock
}

func (m *MockSaleRepository) Create(sale *domain.Sale, tx *gorm.DB) error {
	args := m.Called(sale, tx)
	return args.Error(0)
}

func TestProcessSale_Success(t *testing.T) {
	// Arrange
	mockProductRepo := new(MockProductRepository) // Re-using from product_usecase_test
	mockSaleRepo := new(MockSaleRepository)
	// For unit tests, we pass nil for the db. We will test the logic method directly.
	usecase := NewSaleUsecase(mockProductRepo, mockSaleRepo, nil)

	product := domain.Product{ID: "prod-123", Name: "Test Product", Price: 20.0, Stock: 50}
	sale := &domain.Sale{ProductID: "prod-123", Quantity: 5}

	// Set up expectations
	mockProductRepo.On("GetByID", "prod-123").Return(product, nil)
	// We expect Update to be called with the product having its stock reduced.
	mockProductRepo.On("Update", mock.MatchedBy(func(p *domain.Product) bool {
		return p.ID == "prod-123" && p.Stock == 45
	})).Return(nil)
	mockSaleRepo.On("Create", mock.AnythingOfType("*domain.Sale"), (*gorm.DB)(nil)).Return(nil)

	// Act
	// We call the testable logic method directly, passing nil for the transaction object.
	// The mocks don't use the tx object, so this is safe.
	err := usecase.processSaleLogic(nil, sale)

	// Assert
	assert.NoError(t, err)
	mockProductRepo.AssertExpectations(t)
	mockSaleRepo.AssertExpectations(t)
}

func TestProcessSale_InsufficientStock(t *testing.T) {
	// Arrange
	mockProductRepo := new(MockProductRepository)
	mockSaleRepo := new(MockSaleRepository)
	usecase := NewSaleUsecase(mockProductRepo, mockSaleRepo, nil)

	// This product has less stock than the sale requires.
	product := domain.Product{ID: "prod-123", Name: "Test Product", Price: 20.0, Stock: 2}
	sale := &domain.Sale{ProductID: "prod-123", Quantity: 5}

	// Set up expectations
	mockProductRepo.On("GetByID", "prod-123").Return(product, nil)

	// Act
	err := usecase.processSaleLogic(nil, sale)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrInsufficientStock, err)
	// Verify that Update and Create were never called because the logic should have failed early.
	mockProductRepo.AssertNotCalled(t, "Update")
	mockSaleRepo.AssertNotCalled(t, "Create")
}

func TestProcessSale_ProductNotFound(t *testing.T) {
	// Arrange
	mockProductRepo := new(MockProductRepository)
	mockSaleRepo := new(MockSaleRepository)
	usecase := NewSaleUsecase(mockProductRepo, mockSaleRepo, nil)

	sale := &domain.Sale{ProductID: "prod-that-does-not-exist", Quantity: 5}

	// Set up expectations: GetByID will return a "record not found" error.
	mockProductRepo.On("GetByID", "prod-that-does-not-exist").Return(domain.Product{}, gorm.ErrRecordNotFound)

	// Act
	err := usecase.processSaleLogic(nil, sale)

	// Assert
	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}
