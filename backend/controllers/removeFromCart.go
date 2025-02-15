package controllers

import (
	"net/http"
	"pointafam/backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RemoveFromCart(c *gin.Context) {
	itemIDStr := c.Param("item_id")
	itemID, err := strconv.ParseUint(itemIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	if err := models.DeleteCartItem(db, uint(itemID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove from cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart successfully"})
}
