package repository

import (
	"product-api-go/internal/domain"

	"gorm.io/gorm"
)

type SaleRepository interface {
	Create(sale *domain.Sale, tx *gorm.DB) error // Pass transaction explicitly
}

type PostgresSaleRepo struct {
	db *gorm.DB
}

func NewSaleRepo(db *gorm.DB) *PostgresSaleRepo {
	return &PostgresSaleRepo{db: db}
}

func (r *PostgresSaleRepo) Create(sale *domain.Sale, tx *gorm.DB) error {
	if tx == nil {
		tx = r.db // Use the main DB if no transaction is provided (shouldn't happen in usecase)
	}
	return tx.Create(sale).Error
}
