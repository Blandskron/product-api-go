package domain

import "time"

type Purchase struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	ProductID    string    `json:"product_id" validate:"required"`
	Quantity     int       `json:"quantity" validate:"required,gt=0"`
	PurchaseDate time.Time `json:"purchase_date"`
	TotalCost    float64   `json:"total_cost"` // Calculated based on product cost and quantity
}
