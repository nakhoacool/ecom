package user

import (
	"ecom/service/auth"
	"ecom/types"
	"ecom/utils"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	store types.UserStoreInterface
}

func NewHandler(store types.UserStoreInterface) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) UserRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	//get JSON payload
	var payload types.RegisterUserPayload
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
	// check if the user exists
	_, err = h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user already exists with email %s", payload.Email))
		return
	}

	// hash the password
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// create the user
	err = h.store.CreateUser(types.User{
		ID:        uuid.New().String(),
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "user created"})
}
