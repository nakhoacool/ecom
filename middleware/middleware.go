package middleware

import (
	"context"
	"ecom/config"
	"ecom/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	userIDKey  contextKey = "userID"
	addressKey contextKey = "address"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("missing authorization header"))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid authorization header format"))
			return
		}

		secret := []byte(config.ENV.JWTSecret)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid token"))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid token claims"))
			return
		}

		userID, ok := claims["userID"].(string)
		if !ok {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid token claims"))
			return
		}

		address, ok := claims["address"].(string)
		if !ok {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid token claims"))
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		ctx = context.WithValue(ctx, addressKey, address)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok {
		return "", fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

func GetAddressFromContext(ctx context.Context) (string, error) {
	address, ok := ctx.Value(addressKey).(string)
	if !ok {
		return "", fmt.Errorf("address not found in context")
	}
	return address, nil
}
