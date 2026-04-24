package domain

type Product struct {
	ID    string  `json:"id" gorm:"primaryKey" validate:"-"`
	Name  string  `json:"name" validate:"required,min=3,max=100"`
	Price float64 `json:"price" validate:"required,gt=0"` // Price should always be positive
	Stock int     `json:"stock" validate:"gte=0"`         // Stock can be 0, but not negative
}
