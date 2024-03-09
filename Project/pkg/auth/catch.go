package auth

import (
	"fmt"
	"net/http"
	"os"

	"strings"

	"github.com/golang-jwt/jwt"
)

func JwtPayloadFromRequest(w http.ResponseWriter, r *http.Request) (jwt.MapClaims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header is required", http.StatusUnauthorized)
		return nil, fmt.Errorf("authorization header is required")
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		http.Error(w, "Invalid token format", http.StatusUnauthorized)
		return nil, fmt.Errorf("invalid token format")
	}

	tokenString := authHeader[len(bearerPrefix):]
	jwtSecretKey := os.Getenv("secretKey")
	if jwtSecretKey == "" {
		http.Error(w, "JWT secret key is not set", http.StatusInternalServerError)
		return nil, fmt.Errorf("JWT secret key is not set")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	http.Error(w, "Invalid token claims", http.StatusUnauthorized)
	return nil, fmt.Errorf("invalid token claims")
}
