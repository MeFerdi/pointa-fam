package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"pointafam/backend/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	jwtKey = []byte("your_secret_key")
)

func SetDB(database *gorm.DB) {
	db = database
}

func SignUp(c *gin.Context) {
	var SignupInput struct {
		FirstName       string `json:"first_name" binding:"required"`
		LastName        string `json:"last_name" binding:"required"`
		Email           string `json:"email" binding:"required,email"`
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
		PhoneNumber     string `json:"phone_number"`
		Role            string `json:"role" binding:"required,oneof=farmer retailer"`
	}

	if err := c.ShouldBindJSON(&SignupInput); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	if SignupInput.Password != SignupInput.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Passwords do not match"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(SignupInput.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create user"})
		return
	}

	user := models.User{
		FirstName:   SignupInput.FirstName,
		LastName:    SignupInput.LastName,
		Email:       SignupInput.Email,
		Password:    string(hashedPassword),
		PhoneNumber: SignupInput.PhoneNumber,
		Role:        SignupInput.Role,
	}

	if err := db.Create(&user).Error; err != nil {
		log.Printf("Could not create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create user"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Failed to generate token for email: %s - %v", SignupInput.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Registration successful",
		"token":   tokenString,
		"userID":  user.ID,
		"role":    user.Role,
	})
}

func Login(c *gin.Context) {
	var LoginInput struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&LoginInput); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	var user models.User
	if err := db.Where("email = ?", LoginInput.Email).First(&user).Error; err != nil {
		log.Printf("User not found: %s - %v", LoginInput.Email, err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(LoginInput.Password)); err != nil {
		log.Printf("Invalid password for email: %s - %v", LoginInput.Email, err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Failed to generate token for email: %s - %v", LoginInput.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
		"userID":  user.ID,
		"role":    user.Role,
	})
}

// User Profile
func GetUserProfile(c *gin.Context) {
	var user models.User
	if err := db.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}
