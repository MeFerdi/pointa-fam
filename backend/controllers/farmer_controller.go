package controllers

import (
	"net/http"
	"pointafam/backend/models"
	"pointafam/backend/services"
	"pointafam/backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FarmerService instance
var farmerService *services.FarmerService

// SetFarmerService is called to set the FarmerService instance (usually done in main)
func SetFarmerService(service *services.FarmerService) {
	farmerService = service
}

// RegisterFarmer handles farmer registration
func RegisterFarmer(c *gin.Context) {
	var farmer models.Farmer
	if err := c.ShouldBindJSON(&farmer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Hash the password before saving it to the database
	hashedPassword, err := utils.HashPassword(farmer.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}
	farmer.Password = hashedPassword

	if err := farmerService.CreateFarmer(&farmer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create farmer"})
		return
	}

	c.JSON(http.StatusCreated, farmer)
}

// GetFarmers retrieves all farmers from the database
func GetFarmers(c *gin.Context) {
	farmers, err := farmerService.GetAllFarmers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch farmers"})
		return
	}
	c.JSON(http.StatusOK, farmers)
}

// CreateFarmer registers a new farmer
func CreateFarmer(c *gin.Context) {
	var farmer models.Farmer
	if err := c.ShouldBindJSON(&farmer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := farmerService.CreateFarmer(&farmer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create farmer"})
		return
	}

	c.JSON(http.StatusCreated, farmer)
}

// UpdateFarmer updates an existing farmer's details by ID
func UpdateFarmer(c *gin.Context) {
	var updatedFarmer models.Farmer
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32) // Convert string ID to uint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.ShouldBindJSON(&updatedFarmer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := farmerService.UpdateFarmer(uint(id), &updatedFarmer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update farmer"})
		return
	}

	c.JSON(http.StatusOK, updatedFarmer)
}

// DeleteFarmer removes a farmer by ID from the database.
func DeleteFarmer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32) // Convert string ID to uint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := farmerService.DeleteFarmer(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete farmer"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
