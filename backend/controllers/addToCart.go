package controllers

import (
	"net/http"
	"pointafam/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddToCart(c *gin.Context) {
	var request struct {
		ProductID uint `json:"productId"`
		UserID    uint `json:"userID"`
		Quantity  int  `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	var product models.Product
	if err := db.First(&product, request.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Get or create the cart for the retailer
	cart, err := models.GetCart(db, request.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get or create cart"})
		return
	}

	// Add the item to the cart
	cartItem := models.CartItem{
		CartID:    cart.ID,
		ProductID: request.ProductID,
		Quantity:  request.Quantity,
		UserID:    request.UserID,
	}

	if err := db.Create(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to add to cart"})
		return
	}

	c.JSON(http.StatusOK, cartItem)
}
