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
	var cart []models.CartItem
	retailerID := c.Param("retailer_id")
	db.Where("retailer_id = ?", retailerID).Find(&cart)
	c.HTML(http.StatusOK, "cart_list.html", gin.H{"cartItems": cart})
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
