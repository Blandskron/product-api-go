package domain

import "time"

type Sale struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	ProductID  string    `json:"product_id" validate:"required"`
	Quantity   int       `json:"quantity" validate:"required,gt=0"`
	SaleDate   time.Time `json:"sale_date"`
	TotalPrice float64   `json:"total_price"` // Calculated based on product price and quantity
}
