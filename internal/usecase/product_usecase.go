package usecase

import (
	"product-api-go/internal/domain"
	"product-api-go/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// ValidationError wraps errors from the validator library.
type ValidationError struct {
	Err error // Use a named field to hold the underlying error, typically validator.ValidationErrors
}

func (v ValidationError) Error() string {
	return v.Err.Error()
}

type ProductUsecase struct {
	repo      repository.ProductRepository
	validator *validator.Validate
}

func NewProductUsecase(r repository.ProductRepository) *ProductUsecase {
	return &ProductUsecase{
		repo:      r,
		validator: validator.New(),
	}
}

func (u *ProductUsecase) GetProducts() ([]domain.Product, error) {
	return u.repo.GetAll()
}

func (u *ProductUsecase) GetProduct(id string) (domain.Product, error) {
	return u.repo.GetByID(id)
}

func (u *ProductUsecase) CreateProduct(p *domain.Product) error {
	if err := u.validator.Struct(p); err != nil {
		return ValidationError{Err: err}
	}
	// Generate a new UUID for the product
	p.ID = uuid.New().String()
	return u.repo.Create(p)
}

func (u *ProductUsecase) UpdateProduct(p *domain.Product) error {
	if err := u.validator.Struct(p); err != nil {
		return ValidationError{Err: err}
	}
	return u.repo.Update(p)
}

func (u *ProductUsecase) DeleteProduct(id string) error {
	return u.repo.Delete(id)
}
