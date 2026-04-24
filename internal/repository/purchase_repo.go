package repository

import (
	"product-api-go/internal/domain"

	"gorm.io/gorm"
)

type PurchaseRepository interface {
	Create(purchase *domain.Purchase, tx *gorm.DB) error // Pass transaction explicitly
}

type PostgresPurchaseRepo struct {
	db *gorm.DB
}

func NewPurchaseRepo(db *gorm.DB) *PostgresPurchaseRepo {
	return &PostgresPurchaseRepo{db: db}
}

func (r *PostgresPurchaseRepo) Create(purchase *domain.Purchase, tx *gorm.DB) error {
	if tx == nil {
		tx = r.db // Use the main DB if no transaction is provided (shouldn't happen in usecase)
	}
	return tx.Create(purchase).Error
}
