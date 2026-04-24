package usecase

import (
	"errors"
	"product-api-go/internal/domain"
	"product-api-go/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ErrInsufficientStock is returned when a sale cannot be processed due to lack of stock.
var ErrInsufficientStock = errors.New("insufficient stock")

type SaleUsecase struct {
	productRepo repository.ProductRepository
	saleRepo    repository.SaleRepository
	db          *gorm.DB // To manage transactions
	validator   *validator.Validate
}

func NewSaleUsecase(pr repository.ProductRepository, sr repository.SaleRepository, db *gorm.DB) *SaleUsecase {
	return &SaleUsecase{
		productRepo: pr,
		saleRepo:    sr,
		db:          db,
		validator:   validator.New(),
	}
}

func (u *SaleUsecase) ProcessSale(sale *domain.Sale) error {
	if err := u.validator.Struct(sale); err != nil {
		return ValidationError{Err: err}
	}

	// Start a database transaction
	return u.db.Transaction(func(tx *gorm.DB) error {
		// The business logic is extracted to a separate, testable method.
		return u.processSaleLogic(tx, sale)
	})
}

// processSaleLogic contains the core business logic for processing a sale.
func (u *SaleUsecase) processSaleLogic(tx *gorm.DB, sale *domain.Sale) error {
	// 1. Get Product and check stock
	product, err := u.productRepo.GetByID(sale.ProductID)
	if err != nil {
		return err // Will be gorm.ErrRecordNotFound if not found
	}
	if product.Stock < sale.Quantity {
		return ErrInsufficientStock
	}

	// 2. Deduct stock
	product.Stock -= sale.Quantity
	if err := u.productRepo.Update(&product); err != nil {
		return err
	}

	// 3. Create Sale record
	sale.ID = uuid.New().String()
	sale.SaleDate = time.Now()
	sale.TotalPrice = product.Price * float64(sale.Quantity)
	return u.saleRepo.Create(sale, tx)
}
