package domain

import "time"

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ProductStock struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type ProductPayload struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Image       string  `json:"image" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required"`
}

type ProductRepository interface {
	GetProducts() (*[]Product, error)
	CreateProduct(product Product) error
	GetProductByID(id int) (*Product, error)
	GetProductByIDs(ids []int) (*[]Product, error)
	UpdateProduct(product Product) error

	GetProductStock(productID int) (*ProductStock, error)
	UpdateProductStock(productID, quantity int) error
}
