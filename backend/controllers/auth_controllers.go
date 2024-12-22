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
)

var jwtKey = []byte("your_secret_key")

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("Registration error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Password = string(hashedPassword)
	if err := db.Create(&user).Error; err != nil {
		log.Printf("Failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	if user.Role == "farmer" {
		var farmer models.Farmer
		farmer.Name = user.FirstName + " " + user.LastName
		farmer.Location = user.Location
		farmer.ContactInfo = user.ContactInfo
		farmer.Password = user.Password
		if err := db.Create(&farmer).Error; err != nil {
			log.Printf("Failed to create farmer: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create farmer"})
			return
		}
	} else if user.Role == "retailer" {
		var retailer models.Retailer
		retailer.Name = user.FirstName + " " + user.LastName
		retailer.Location = user.Location
		retailer.ContactInfo = user.ContactInfo
		retailer.Password = user.Password
		if err := db.Create(&retailer).Error; err != nil {
			log.Printf("Failed to create retailer: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create retailer"})
			return
		}
	}

	log.Printf("User registered successfully: %s", user.Email)
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("Login error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dbUser models.User
	if err := db.Where("email = ?", user.Email).First(&dbUser).Error; err != nil {
		log.Printf("Invalid login attempt for email: %s - %v", user.Email, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		log.Printf("Invalid password for email: %s - %v", user.Email, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(int(dbUser.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Failed to generate token for email: %s - %v", user.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	role := dbUser.Role
	log.Printf("User logged in successfully: %s with role: %s", user.Email, role)

	if role == "farmer" {
		c.JSON(http.StatusOK, gin.H{"token": tokenString, "redirect": "/static/farmer_dashboard.html"})
	} else if role == "retailer" {
		c.JSON(http.StatusOK, gin.H{"token": tokenString, "redirect": "/static/retailer_dashboard.html"})
	}
}
