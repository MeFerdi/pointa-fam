package controllers

import (
	"net/http"
	"pointafam/backend/models"
	"pointafam/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RetailerController struct {
	Service *services.RetailerService
}

// NewRetailerController initializes a new RetailerController with a RetailerService instance
func NewRetailerController(service *services.RetailerService) *RetailerController {
	return &RetailerController{Service: service}
}

// CreateRetailerHandler handles the creation of a new retailer
func (ctrl *RetailerController) CreateRetailerHandler(c *gin.Context) {
	var retailer models.Retailer
	if err := c.ShouldBindJSON(&retailer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	if err := ctrl.Service.CreateRetailer(&retailer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, retailer)
}

// GetAllRetailersHandler handles retrieving all retailers
func (ctrl *RetailerController) GetAllRetailersHandler(c *gin.Context) {
	retailers, err := ctrl.Service.GetAllRetailers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, retailers)
}

// UpdateRetailerHandler handles updating an existing retailer by ID
func (ctrl *RetailerController) UpdateRetailerHandler(c *gin.Context) {
	id := c.Param("id")
	_, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}
	var updatedRetailer models.Retailer
	if err := c.ShouldBindJSON(&updatedRetailer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	if err := ctrl.Service.UpdateRetailer(id, &updatedRetailer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedRetailer)
}

// DeleteRetailerHandler handles deleting a retailer by ID
func (ctrl *RetailerController) DeleteRetailerHandler(c *gin.Context) {
	id := c.Param("id")
	_, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}
	if err := ctrl.Service.DeleteRetailer(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Retailer deleted successfully"})
}
