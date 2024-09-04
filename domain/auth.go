package domain

type AuthService interface {
	CreateToken(secret []byte, userID string, firstName string, lastName string, email string, address string) (string, error)
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) error
}
