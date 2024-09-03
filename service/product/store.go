package product

import (
	"database/sql"
	"ecom/domain"
	"errors"
	"fmt"
	"sync"
)

type Store struct {
	db *sql.DB
	mu sync.Mutex
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetProducts() (*[]domain.Product, error) {
	rows, err := s.db.Query(`
        SELECT p.id, p.name, p.description, p.image, p.price, p.createdAt, ps.quantity
        FROM products p
        LEFT JOIN product_stock ps ON p.id = ps.product_id
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]domain.Product, 0)
	for rows.Next() {
		var product domain.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Image,
			&product.Price,
			&product.CreatedAt,
			&product.Quantity,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return &products, nil
}

func (s *Store) GetProductByID(id int) (*domain.Product, error) {
	row := s.db.QueryRow(`
		SELECT p.id, p.name, p.description, p.image, p.price, p.createdAt, ps.quantity
		FROM products p
		LEFT JOIN product_stock ps ON p.id = ps.product_id
		WHERE p.id = ?
	`, id)

	product := new(domain.Product)
	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Image,
		&product.Price,
		&product.CreatedAt,
		&product.Quantity,
	)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *Store) GetProductByIDs(ids []int) (*[]domain.Product, error) {
	rows, err := s.db.Query(`
		SELECT p.id, p.name, p.description, p.image, p.price, p.createdAt, ps.quantity
		FROM products p
		LEFT JOIN product_stock ps ON p.id = ps.product_id
		WHERE p.id IN (?)
	`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]domain.Product, 0)
	for rows.Next() {
		var product domain.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Image,
			&product.Price,
			&product.CreatedAt,
			&product.Quantity,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return &products, nil
}

func (s *Store) CreateProduct(product domain.Product) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO products (name, description, image, price)
		VALUES (?, ?, ?, ?)
	`, product.Name, product.Description, product.Image, product.Price)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO product_stock (product_id, quantity)
		VALUES (LAST_INSERT_ID(), ?)
	`, product.Quantity)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *Store) UpdateProduct(product domain.Product) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE products
		SET name = ?, description = ?, image = ?, price = ?
		WHERE id = ?
	`, product.Name, product.Description, product.Image, product.Price, product.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		UPDATE product_stock
		SET quantity = ?
		WHERE product_id = ?
	`, product.Quantity, product.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *Store) GetProductStock(productID int) (*domain.ProductStock, error) {
	row := s.db.QueryRow("SELECT product_id, quantity FROM product_stock WHERE product_id = ?", productID)

	stock := new(domain.ProductStock)
	err := row.Scan(&stock.ProductID, &stock.Quantity)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("product stock not found")
	} else if err != nil {
		return nil, err
	}

	return stock, nil
}

func (s *Store) UpdateProductStock(productID, quantity int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	var currentQuantity int
	err = tx.QueryRow("SELECT quantity FROM product_stock WHERE product_id = ?", productID).Scan(&currentQuantity)
	if err != nil {
		return err
	}

	newQuantity := currentQuantity + quantity
	if newQuantity < 0 {
		return fmt.Errorf("insufficient stock")
	}

	_, err = tx.Exec("UPDATE product_stock SET quantity = ? WHERE product_id = ?", newQuantity, productID)
	return err
}
