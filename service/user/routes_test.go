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

type mockAuthStore struct{}

func (m *mockAuthStore) CreateToken(secret []byte, userID string, firstName string, lastName string, email string, address string) (string, error) {
	return "mockToken", nil
}

func (m *mockAuthStore) HashPassword(password string) (string, error) {
	return "", nil
}

func (m *mockAuthStore) ComparePassword(hashedPassword, password string) error {
	if hashedPassword == "hashedPassword" && password == "password" {
		return nil
	}
	return fmt.Errorf("invalid password")
}

type mockUserStore struct{}

func (m *mockUserStore) GetUserByEmail(email string) (*domain.User, error) {
	if email == "existing.user@gmail.com" {
		return &domain.User{
			ID:        "1",
			FirstName: "Existing",
			LastName:  "User",
			Email:     "existing.user@gmail.com",
			Password:  "hashedPassword",
			Address:   "Existing User Address",
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

func TestHandleRegister(t *testing.T) {
	userStore := &mockUserStore{}
	authStore := &mockAuthStore{}
	handler := NewHandler(userStore, authStore)

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

func TestHandleLogin(t *testing.T) {
	userStore := &mockUserStore{}
	authStore := &mockAuthStore{}
	handler := NewHandler(userStore, authStore)

	t.Run("should return 400 if the payload is invalid", func(t *testing.T) {
		payload := "invalid payload"

		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte(payload)))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/login", handler.handleLogin)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should return 400 if user not found", func(t *testing.T) {
		payload := domain.LoginUserPayload{
			Email:    "non.existing.user@gmail.com",
			Password: "password",
		}

		marshaled, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(marshaled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/login", handler.handleLogin)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should return 400 if password is incorrect", func(t *testing.T) {
		payload := domain.LoginUserPayload{
			Email:    "existing.user@gmail.com",
			Password: "wrongpassword",
		}

		marshaled, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(marshaled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/login", handler.handleLogin)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should return 200 when login is successful", func(t *testing.T) {
		payload := domain.LoginUserPayload{
			Email:    "existing.user@gmail.com",
			Password: "password",
		}

		marshaled, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(marshaled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/login", handler.handleLogin)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}
