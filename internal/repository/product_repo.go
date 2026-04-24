package repository

import (
	"product-api-go/internal/domain"

	"gorm.io/gorm"
)

type ProductRepository interface {
	GetAll() ([]domain.Product, error)
	GetByID(id string) (domain.Product, error)
	Create(product *domain.Product) error
	Update(product *domain.Product) error
	Delete(id string) error
}

type PostgresProductRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) ProductRepository {
	return &PostgresProductRepo{db: db}
}

func (r *PostgresProductRepo) GetAll() ([]domain.Product, error) {
	var products []domain.Product
	err := r.db.Find(&products).Error
	return products, err
}

func (r *PostgresProductRepo) GetByID(id string) (domain.Product, error) {
	var product domain.Product
	err := r.db.First(&product, "id = ?", id).Error
	return product, err
}

func (r *PostgresProductRepo) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

func (r *PostgresProductRepo) Update(product *domain.Product) error {
	// Usamos .Model y .Where para asegurar que solo actualizamos un registro existente.
	result := r.db.Model(&domain.Product{}).Where("id = ?", product.ID).Updates(product)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // No se actualizó ninguna fila, así que no existía.
	}
	return nil
}

func (r *PostgresProductRepo) Delete(id string) error {
	result := r.db.Delete(&domain.Product{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
