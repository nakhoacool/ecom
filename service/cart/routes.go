package cart

import (
	"ecom/domain"
	"ecom/utils"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	orderStore   domain.OrderRepository
	productStore domain.ProductRepository
}

func NewHandler(orderStore domain.OrderRepository, productStore domain.ProductRepository) *Handler {
	return &Handler{
		orderStore:   orderStore,
		productStore: productStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", h.handleCheckout).Methods(http.MethodPost)
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var payload domain.CartCheckoutPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	err := utils.Validate.Struct(payload)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", validationErrors))
		return
	}

	// get the products from the store
	productIDs, err := getCartItemsIDs(payload.Items)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	products, err := h.productStore.GetProductByIDs(productIDs)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// create the order
	// TODO: implement JWT authentication and get the user ID from the token
	orderID, totalPrice, err := h.createOrder("1", payload.Items, products)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"orderID":    orderID,
		"totalPrice": totalPrice,
	})
}
