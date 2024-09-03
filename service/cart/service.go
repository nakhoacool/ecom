package cart

import (
	"ecom/domain"
	"fmt"
)

func getCartItemsIDs(items []domain.CartItem) ([]int, error) {
	productIDs := make([]int, len(items))
	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for product %d", item.ProductID)
		}
		productIDs[i] = item.ProductID
	}
	return productIDs, nil
}

func (h *Handler) createOrder(userID string, items []domain.CartItem, products *[]domain.Product) (int, float64, error) {
	// Create a map for products
	productMap := make(map[int]domain.Product)
	for _, product := range *products {
		productMap[product.ID] = product
	}

	// Check if the products are available in stock
	for _, item := range items {
		product, exists := productMap[item.ProductID]
		if !exists {
			return 0, 0, fmt.Errorf("product %d not found", item.ProductID)
		}
		if item.Quantity > product.Quantity {
			return 0, 0, fmt.Errorf("insufficient stock for product %d", item.ProductID)
		}
	}

	// Calculate the total price
	var totalPrice float64
	for _, item := range items {
		product := productMap[item.ProductID]
		totalPrice += float64(item.Quantity) * product.Price
	}

	// Reduce the stock
	for _, item := range items {
		err := h.productStore.UpdateProductStock(item.ProductID, -item.Quantity)
		if err != nil {
			return 0, 0, err
		}
	}

	// Create the order
	order := domain.Order{
		UserID: userID,
		Total:  totalPrice,
		Status: "pending",
		// Add other necessary fields like Address if needed
	}
	orderID, err := h.orderStore.CreateOrder(order)
	if err != nil {
		return 0, 0, err
	}

	// Create the order items
	for _, item := range items {
		orderItem := domain.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductID].Price,
		}
		err := h.orderStore.CreateOrderItem(orderItem)
		if err != nil {
			return 0, 0, err
		}
	}

	return orderID, totalPrice, nil
}
