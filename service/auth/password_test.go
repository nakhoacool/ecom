package auth

import "testing"

func TestHashPassword(t *testing.T) {
	password := "password"
	hash, err := HashPassword(password)
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
	password := "password"
	hash, err := HashPassword(password)
	if err != nil {
		t.Error(err)
	}
	err = ComparePassword(hash, password)
	if err != nil {
		t.Error(err)
	}
	err = ComparePassword(hash, "wrongpassword")
	if err == nil {
		t.Error("wrong password should not match")
	}
}
