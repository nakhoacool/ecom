package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func CreateToken(secret []byte, userID string) (string, error) {
	expiration := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    expiration.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
