package controllers

import (
	"log"
	"net/http"
	"pointafam/backend/models"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var db *gorm.DB
var jwtKey = []byte("your_secret_key")

func SetDB(database *gorm.DB) {
	db = database
}

func SignUp(c *gin.Context) {
	var SignupInput struct {
		Username        string `json:"username" binding:"required"`
		Email           string `json:"email" binding:"required,email"`
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
		PhoneNumber     string `json:"phoneNumber"`
		Location        string `json:"location"`
		Role            string `json:"role" binding:"required,oneof=farmer retailer"` // Validate role
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
		Username:    SignupInput.Username,
		Email:       SignupInput.Email,
		Password:    string(hashedPassword),
		PhoneNumber: SignupInput.PhoneNumber,
		Location:    SignupInput.Location,
		Role:        SignupInput.Role,
	}

	if err := db.Create(&user).Error; err != nil {
		log.Printf("Could not create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create user"})
		return
	}

	// Create associated Farmer or Retailer record based on role
	if SignupInput.Role == "farmer" {
		farmer := models.Farmer{
			ID: user.ID,

			Name:        user.Username,
			PhoneNumber: SignupInput.PhoneNumber,
			Location:    SignupInput.Location,
		}
		if err := db.Create(&farmer).Error; err != nil {
			log.Printf("Could not create farmer: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create farmer"})
			return
		}
	} else if SignupInput.Role == "retailer" {
		retailer := models.Retailer{
			ID:          user.ID,
			Name:        user.Username,
			PhoneNumber: SignupInput.PhoneNumber,
			Location:    SignupInput.Location,
		}
		if err := db.Create(&retailer).Error; err != nil {
			log.Printf("Could not create retailer: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create retailer"})
			return
		}
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
func GetUserProfile(c *gin.Context) {
	var user models.User
	if err := db.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}
