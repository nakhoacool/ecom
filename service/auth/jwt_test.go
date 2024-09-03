package auth

import "testing"

func TestCreateToken(t *testing.T) {
	secret := []byte("secret")
	userID := "123"
	token, err := CreateToken(secret, userID)
	if err != nil {
		t.Error(err)
	}
	if token == "" {
		t.Error("token is empty")
	}
}
