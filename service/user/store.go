package user

import (
	"database/sql"
	"ecom/domain"
	"errors"
	"fmt"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(email string) (*domain.User, error) {
	row := s.db.QueryRow("SELECT * FROM users WHERE email = ?", email)

	user := new(domain.User)
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Address,
		&user.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) GetUserByID(id int) (*domain.User, error) {
	row := s.db.QueryRow("SELECT * FROM users WHERE id = ?", id)

	user := new(domain.User)
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Address,
		&user.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) CreateUser(user domain.User) error {
	_, err := s.db.Exec(
		"INSERT INTO users (id,firstName, lastName, email, password, address) VALUES (?,?, ?, ?, ?, ?)",
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.Address,
	)

	return err
}
