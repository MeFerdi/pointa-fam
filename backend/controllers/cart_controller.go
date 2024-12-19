package controllers

import (
	"net/http"
	"pointafam/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func SetDB(database *gorm.DB) {
	db = database
}

func AddToCart(c *gin.Context) {
	var cartItem models.CartItem
	if err := c.ShouldBindJSON(&cartItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&cartItem)
	c.JSON(http.StatusOK, cartItem)
}

func ViewCart(c *gin.Context) {
	var cart models.Cart
	retailerID := c.Param("retailer_id")
	db.Preload("Items").Where("retailer_id = ?", retailerID).First(&cart)
	c.JSON(http.StatusOK, cart)
}

func DeleteFromCart(c *gin.Context) {
	var cartItem models.CartItem
	itemID := c.Param("item_id")
	db.Where("id = ?", itemID).Delete(&cartItem)
	c.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
}
