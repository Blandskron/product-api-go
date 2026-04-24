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
		// 1. Get Product (or create if new)
		product, err := u.productRepo.GetByID(purchase.ProductID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Product does not exist, create a new one (assuming purchase includes product details)
				// For simplicity, we'll assume product details are already in the purchase.ProductID
				// In a real app, you'd likely have a separate endpoint to create products or
				// the purchase payload would include more product details.
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
	})
}
