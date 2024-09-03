package domain

import "time"

type Order struct {
	ID        int       `json:"id"`
	UserID    string    `json:"userId"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
}

type OrderItem struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"orderId"`
	ProductID int     `json:"productId"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type OrderRepository interface {
	CreateOrder(order Order) (int, error)
	CreateOrderItem(orderItem OrderItem) error
}
