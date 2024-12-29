package controllers

import (
	"net/http"
	"pointafam/backend/models"

	"github.com/gin-gonic/gin"
)

func AddToCart(c *gin.Context) {
	var cartItem models.CartItem
	if err := c.ShouldBindJSON(&cartItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&cartItem)
	ViewCart(c)
}

func ViewCart(c *gin.Context) {
	var cartItems []models.CartItem
	retailerID := c.Param("retailer_id")
	db.Joins("JOIN carts ON carts.id = cart_items.cart_id").Where("carts.retailer_id = ?", retailerID).Find(&cartItems)
	c.JSON(http.StatusOK, gin.H{"cartItems": cartItems})
}

func DeleteFromCart(c *gin.Context) {
	var cartItem models.CartItem
	itemID := c.Param("item_id")
	if err := db.Where("id = ?", itemID).First(&cartItem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}
	db.Delete(&cartItem)
	ViewCart(c)
}
