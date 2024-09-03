package user

import (
	"bytes"
	"ecom/domain"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockUserStore struct{}

func (m *mockUserStore) GetUserByEmail(email string) (*domain.User, error) {
	if email == "existing.user@gmail.com" {
		return &domain.User{
			FirstName: "Existing",
			LastName:  "User",
			Email:     "existing.user@gmail.com",
			Password:  "12345678",
		}, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (m *mockUserStore) GetUserByID(id int) (*domain.User, error) {
	return nil, nil
}

func (m *mockUserStore) CreateUser(user domain.User) error {
	return nil
}

func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	t.Run("should return 400 if the payload is invalid", func(t *testing.T) {
		payload := domain.RegisterUserPayload{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe.com",
			Password:  "12345678",
		}

		marshaled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshaled))

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/register", handler.handleRegister)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should return 400 if user already exists", func(t *testing.T) {
		payload := domain.RegisterUserPayload{
			FirstName: "Existing",
			LastName:  "User",
			Email:     "existing.user@gmail.com",
			Password:  "12345678",
		}

		marshaled, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshaled))

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/register", handler.handleRegister)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should return 201 when user created successfully", func(t *testing.T) {
		payload := domain.RegisterUserPayload{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@gmail.com",
			Password:  "12345678",
		}

		marshaled, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshaled))

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/register", handler.handleRegister)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}
