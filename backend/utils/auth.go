package utils

import (
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a hashed password with its plain-text version
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidateStruct validates a struct using the go-playground validator
func ValidateStruct(s interface{}) error {
	validate := validator.New()
	return validate.Struct(s)
}

// UserRegistration struct for validating user registration data
type UserRegistration struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// ValidateUserRegistration validates user registration data
func ValidateUserRegistration(user UserRegistration) error {
	if err := ValidateStruct(user); err != nil {
		return err
	}
	return nil
}
