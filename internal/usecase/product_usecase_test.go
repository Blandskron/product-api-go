package usecase

import (
	"testing"

	"product-api-go/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductRepository is a mock implementation of the ProductRepository interface.
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) GetAll() ([]domain.Product, error) {
	args := m.Called()
	return args.Get(0).([]domain.Product), args.Error(1)
}
func (m *MockProductRepository) GetByID(id string) (domain.Product, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Product), args.Error(1)
}
func (m *MockProductRepository) Create(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}
func (m *MockProductRepository) Update(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}
func (m *MockProductRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateProduct_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductRepository)
	usecase := NewProductUsecase(mockRepo)
	product := &domain.Product{
		Name:  "Valid Product",
		Price: 10.0,
		Stock: 100,
	}

	// Set up the expectation.
	// We expect the Create method to be called once with any product object.
	// We configure it to return no error (nil).
	mockRepo.On("Create", mock.AnythingOfType("*domain.Product")).Return(nil)

	// Act
	err := usecase.CreateProduct(product)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, product.ID, "Product ID should be set after creation")
	mockRepo.AssertExpectations(t) // Verify that the mock was called as expected.
}

func TestCreateProduct_ValidationError(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductRepository)
	usecase := NewProductUsecase(mockRepo)
	// Product with an empty name, which should fail validation.
	product := &domain.Product{
		Name:  "",
		Price: 10.0,
		Stock: 100,
	}

	// Act
	err := usecase.CreateProduct(product)

	// Assert
	assert.Error(t, err)
	_, ok := err.(ValidationError)
	assert.True(t, ok, "Error should be of type ValidationError")
	// Ensure the repository's Create method was never called.
	mockRepo.AssertNotCalled(t, "Create")
}
