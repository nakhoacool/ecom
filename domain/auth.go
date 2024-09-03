package domain

type AuthService interface {
	CreateToken(secret []byte, userID string) (string, error)
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
}
