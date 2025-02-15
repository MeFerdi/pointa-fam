package controllers

import (
	"net/http"
	"pointafam/backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ViewCart(c *gin.Context) {
	retailerIDStr := c.Param("retailer_id")
	retailerID, err := strconv.ParseUint(retailerIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid retailer ID"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	cart, err := models.GetCart(db, uint(retailerID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart"})
		return
	}

	c.JSON(http.StatusOK, cart)
}
