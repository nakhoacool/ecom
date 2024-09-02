package user

import (
	"database/sql"
	"ecom/types"
	"errors"
	"fmt"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	row := s.db.QueryRow("SELECT * FROM users WHERE email = ?", email)

	user := new(types.User)
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	row := s.db.QueryRow("SELECT * FROM users WHERE id = ?", id)

	user := new(types.User)
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec(
		"INSERT INTO users (id,firstName, lastName, email, password) VALUES (?,?, ?, ?, ?)",
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
	)

	return err
}
