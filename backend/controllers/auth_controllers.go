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
		FirstName       string `json:"firstName" binding:"required"`
		LastName        string `json:"lastName" binding:"required"`
		Email           string `json:"email" binding:"required,email"`
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
		PhoneNumber     string `json:"phoneNumber"`
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
		Email:       SignupInput.Email,
		Password:    string(hashedPassword),
		PhoneNumber: SignupInput.PhoneNumber,
		FirstName:   SignupInput.FirstName,
		LastName:    SignupInput.LastName,
		Role:        SignupInput.Role,
	}

	if err := db.Create(&user).Error; err != nil {
		log.Printf("Could not create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create user"})
		return
	}

	if user.Role == "farmer" {
		var farmer models.Farmer
		farmer.Name = user.FirstName + " " + user.LastName
		farmer.PhoneNumber = user.PhoneNumber // Assuming you want to store this as well

		if err := db.Create(&farmer).Error; err != nil {
			log.Printf("Failed to create farmer: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create farmer"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"redirect": "/farmer/dashboard"})
	} else if user.Role == "retailer" {
		var retailer models.Retailer
		retailer.Name = user.FirstName + " " + user.LastName
		retailer.PhoneNumber = user.PhoneNumber // Assuming you want to store this as well

		if err := db.Create(&retailer).Error; err != nil {
			log.Printf("Failed to create retailer: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create retailer"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"redirect": "/retailer/dashboard"})
	}
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
		"role":    user.Role,
	})
}
