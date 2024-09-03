package product

import (
	"ecom/types"
	"ecom/utils"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	store types.ProductStoreInterface
}

func NewHandler(store types.ProductStoreInterface) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) ProductRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.handleGetProducts).Methods("GET")
	router.HandleFunc("/products", h.handleCreateProduct).Methods("POST")
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, products)
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var payload types.ProductPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the payload
	err := utils.Validate.Struct(payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// create the product
	err = h.store.CreateProduct(types.Product{
		Name:        payload.Name,
		Description: payload.Description,
		Image:       payload.Image,
		Price:       payload.Price,
		Quantity:    payload.Quantity,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "product created"})
}
