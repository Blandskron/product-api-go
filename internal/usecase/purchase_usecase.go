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

// ErrProductNotFoundForPurchase is returned when a purchase is attempted for a non-existent product.
var ErrProductNotFoundForPurchase = errors.New("product not found for purchase, please create product first")

type PurchaseUsecase struct {
	productRepo  repository.ProductRepository
	purchaseRepo repository.PurchaseRepository
	db           *gorm.DB // To manage transactions
	validator    *validator.Validate
}

func NewPurchaseUsecase(pr repository.ProductRepository, pur repository.PurchaseRepository, db *gorm.DB) *PurchaseUsecase {
	return &PurchaseUsecase{
		productRepo:  pr,
		purchaseRepo: pur,
		db:           db,
		validator:    validator.New(),
	}
}

func (u *PurchaseUsecase) ProcessPurchase(purchase *domain.Purchase) error {
	if err := u.validator.Struct(purchase); err != nil {
		return ValidationError{Err: err}
	}

	// Start a database transaction
	return u.db.Transaction(func(tx *gorm.DB) error {
		return u.processPurchaseLogic(tx, purchase)
	})
}

// processPurchaseLogic contains the core business logic for processing a purchase.
func (u *PurchaseUsecase) processPurchaseLogic(tx *gorm.DB, purchase *domain.Purchase) error {
	// 1. Get Product
	product, err := u.productRepo.GetByID(purchase.ProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// In a real app, you might create the product here if the payload contained enough info.
			// For now, we enforce that the product must exist.
			return ErrProductNotFoundForPurchase
		}
		return err
	}

	// 2. Add stock
	product.Stock += purchase.Quantity
	if err := u.productRepo.Update(&product); err != nil {
		return err
	}

	// 3. Create Purchase record
	purchase.ID = uuid.New().String()
	purchase.PurchaseDate = time.Now()
	purchase.TotalCost = product.Price * float64(purchase.Quantity) // Assuming purchase cost is same as sale price for simplicity
	return u.purchaseRepo.Create(purchase, tx)
}
