package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {
	store := NewStore()
	secret := []byte("mysecret")
	userID := "123"
	firstName := "John"
	lastName := "Doe"
	email := "john.doe@example.com"
	address := "123 Main St"

	tokenString, err := store.CreateToken(secret, userID, firstName, lastName, email, address)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	assert.NoError(t, err)
	assert.NotNil(t, token)

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, userID, claims["userID"])
	assert.Equal(t, firstName, claims["firstName"])
	assert.Equal(t, lastName, claims["lastName"])
	assert.Equal(t, email, claims["email"])
	assert.Equal(t, address, claims["address"])
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), time.Unix(int64(claims["exp"].(float64)), 0), time.Minute)
}

func TestHashPassword(t *testing.T) {
	store := NewStore()
	password := "password"
	hash, err := store.HashPassword(password)
	if err != nil {
		t.Error(err)
	}
	if hash == "" {
		t.Error("hash is empty")
	}
	if hash == password {
		t.Error("hash is not hashed")
	}
}

func TestComparePassword(t *testing.T) {
	store := NewStore()
	password := "password"
	hash, err := store.HashPassword(password)
	if err != nil {
		t.Error(err)
	}
	err = store.ComparePassword(hash, password)
	if err != nil {
		t.Error(err)
	}
	err = store.ComparePassword(hash, "wrongpassword")
	if err == nil {
		t.Error("wrong password should not match")
	}
}
