package services

import (
	"errors"
	"strconv"

	"github.com/dgrijalva/jwt-go"
)

// ValidateToken validates the JWT token and returns the userID and role
func ValidateToken(tokenString string) (uint, string, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte("your-secret-key"), nil // Replace with your actual secret key
	})
	if err != nil {
		return 0, "", err
	}

	// Extract the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr := claims["userID"].(string)
		role := claims["role"].(string)

		// Convert userID to uint
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			return 0, "", errors.New("invalid userID in token")
		}

		return uint(userID), role, nil
	}

	return 0, "", errors.New("invalid token")
}
