package product

import (
	"database/sql"
	"ecom/types"
	"sync"
)

type Store struct {
	db *sql.DB
	mu sync.Mutex
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

//func (s *Store) GetProductStock(productID int) (*types.ProductStock, error) {
//	row := s.db.QueryRow("SELECT product_id, quantity FROM product_stock WHERE product_id = ?", productID)
//
//	stock := new(types.ProductStock)
//	err := row.Scan(&stock.ProductID, &stock.Quantity)
//	if errors.Is(err, sql.ErrNoRows) {
//		return nil, fmt.Errorf("product stock not found")
//	} else if err != nil {
//		return nil, err
//	}
//
//	return stock, nil
//}

func (s *Store) GetProducts() (*[]types.Product, error) {
	rows, err := s.db.Query(`
        SELECT p.id, p.name, p.description, p.image, p.price, p.createdAt, ps.quantity
        FROM products p
        LEFT JOIN product_stock ps ON p.id = ps.product_id
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]types.Product, 0)
	for rows.Next() {
		var product types.Product
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

func (s *Store) CreateProduct(product types.Product) error {
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
