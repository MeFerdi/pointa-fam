package controllers

import (
	"net/http"
	"pointafam/backend/models"
	"pointafam/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FarmerController struct {
	Service *services.FarmerService
}

// NewFarmerController initializes a new FarmerController with a FarmerService instance
func NewFarmerController(service *services.FarmerService) *FarmerController {
	return &FarmerController{Service: service}
}

// CreateFarmerHandler handles the creation of a new farmer
func (ctrl *FarmerController) CreateFarmerHandler(c *gin.Context) {
	var farmer models.Farmer
	if err := c.ShouldBindJSON(&farmer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	if err := ctrl.Service.CreateFarmer(&farmer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, farmer)
}

// GetAllFarmersHandler handles retrieving all farmers
func (ctrl *FarmerController) GetAllFarmersHandler(c *gin.Context) {
	farmers, err := ctrl.Service.GetAllFarmers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, farmers)
}

// UpdateFarmerHandler handles updating an existing farmer by ID
func (ctrl *FarmerController) UpdateFarmerHandler(c *gin.Context) {
	id := c.Param("id")
	uintID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}
	var updatedFarmer models.Farmer
	if err := c.ShouldBindJSON(&updatedFarmer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	if err := ctrl.Service.UpdateFarmer(uint(uintID), &updatedFarmer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedFarmer)
}

// DeleteFarmerHandler handles deleting a farmer by ID
func (ctrl *FarmerController) DeleteFarmerHandler(c *gin.Context) {
	id := c.Param("id")
	uintID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}
	if err := ctrl.Service.DeleteFarmer(uint(uintID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Farmer deleted successfully"})
}
