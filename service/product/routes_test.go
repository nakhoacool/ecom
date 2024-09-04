package product

import (
	"bytes"
	"ecom/domain"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockProductStore struct {
	products []domain.Product
}

func (m *mockProductStore) GetProducts() (*[]domain.Product, error) {
	return &m.products, nil
}

func (m *mockProductStore) CreateProduct(product domain.Product) error {
	m.products = append(m.products, product)
	return nil
}

func (m *mockProductStore) GetProductByID(id int) (*domain.Product, error) {
	for _, product := range m.products {
		if product.ID == id {
			return &product, nil
		}
	}
	return nil, errors.New("product not found")
}

func (m *mockProductStore) GetProductByIDs(ids []int) (*[]domain.Product, error) {
	return nil, nil
}

func (m *mockProductStore) UpdateProduct(product domain.Product) error {
	return nil
}

func (m *mockProductStore) GetProductStock(productID int) (*domain.ProductStock, error) {
	return nil, nil
}

func (m *mockProductStore) UpdateProductStock(productID, quantity int) error {
	return nil
}

func TestHandleGetProducts(t *testing.T) {
	store := &mockProductStore{
		products: []domain.Product{
			{ID: 1, Name: "Product 1"},
			{ID: 2, Name: "Product 2"},
		},
	}
	handler := NewHandler(store)

	req, err := http.NewRequest(http.MethodGet, "/products", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/products", handler.handleGetProducts)

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var products []domain.Product
	if err := json.NewDecoder(rr.Body).Decode(&products); err != nil {
		t.Fatal(err)
	}

	if len(products) != 2 {
		t.Errorf("expected 2 products, got %d", len(products))
	}
}

func TestHandleCreateProduct(t *testing.T) {
	store := &mockProductStore{}
	handler := NewHandler(store)

	payload := domain.ProductPayload{
		Name:        "New Product",
		Description: "New Product Description",
		Image:       "image.png",
		Price:       10.0,
		Quantity:    5,
	}

	marshaled, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(marshaled))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/products", handler.handleCreateProduct)

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
	}

	if len(store.products) != 1 {
		t.Errorf("expected 1 product, got %d", len(store.products))
	}
}

func TestHandleGetProductByID(t *testing.T) {
	store := &mockProductStore{
		products: []domain.Product{
			{ID: 1, Name: "Product 1"},
		},
	}
	handler := NewHandler(store)

	req, err := http.NewRequest(http.MethodGet, "/products/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/products/{id}", handler.handleGetProductByID)

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var product domain.Product
	if err := json.NewDecoder(rr.Body).Decode(&product); err != nil {
		t.Fatal(err)
	}

	if product.ID != 1 {
		t.Errorf("expected product ID 1, got %d", product.ID)
	}
}
