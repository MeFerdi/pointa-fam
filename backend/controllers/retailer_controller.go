package controllers

import (
	"net/http"
	"pointafam/backend/models"
	"pointafam/backend/services"
	"pointafam/backend/utils"

	"github.com/gin-gonic/gin"
)

// RetailerService instance
var retailerService *services.RetailerService

// SetRetailerService is called to set the RetailerService instance (usually done in main)
func SetRetailerService(service *services.RetailerService) {
	retailerService = service
}

// RegisterRetailer handles retailer registration
func RegisterRetailer(c *gin.Context) {
	var retailer models.Retailer
	if err := c.ShouldBindJSON(&retailer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	hashedPassword, err := utils.HashPassword(retailer.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}
	retailer.Password = hashedPassword

	if err := retailerService.CreateRetailer(&retailer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create retailer"})
		return
	}

	c.JSON(http.StatusCreated, retailer)
}

// GetRetailers retrieves all retailers from the database
func GetRetailers(c *gin.Context) {
	retailers, err := retailerService.GetAllRetailers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch retailers"})
		return
	}

	c.JSON(http.StatusOK, retailers)
}

// CreateRetailer registers a new retailer
func CreateRetailer(c *gin.Context) {
	var retailer models.Retailer
	if err := c.ShouldBindJSON(&retailer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := retailerService.CreateRetailer(&retailer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create retailer"})
		return
	}

	c.JSON(http.StatusCreated, retailer)
}

// UpdateRetailer updates an existing retailer's details
func UpdateRetailer(c *gin.Context) {
	var retailer models.Retailer
	id := c.Param("id")

	if err := c.ShouldBindJSON(&retailer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := retailerService.UpdateRetailer(id, &retailer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update retailer"})
		return
	}

	c.JSON(http.StatusOK, retailer)
}

// DeleteRetailer removes a retailer by ID from the database.
func DeleteRetailer(c *gin.Context) {
	id := c.Param("id")

	if err := retailerService.DeleteRetailer(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete retailer"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
