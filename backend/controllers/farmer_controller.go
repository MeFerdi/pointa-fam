package controllers

import (
	"net/http"

	"pointafam/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetFarmers(c *gin.Context) {
	var farmers []models.Farmer
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Find(&farmers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch farmers"})
		return
	}
	c.JSON(http.StatusOK, farmers)
}

func CreateFarmer(c *gin.Context) {
	var farmer models.Farmer
	if err := c.ShouldBindJSON(&farmer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	if err := db.Create(&farmer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create farmer"})
		return
	}

	c.JSON(http.StatusCreated, farmer)
}

func UpdateFarmer(c *gin.Context) {
	var farmer models.Farmer
	id := c.Param("id")

	db := c.MustGet("db").(*gorm.DB)
	if err := db.First(&farmer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Farmer not found"})
		return
	}

	if err := c.ShouldBindJSON(&farmer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Save(&farmer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update farmer"})
		return
	}

	c.JSON(http.StatusOK, farmer)
}

func DeleteFarmer(c *gin.Context) {
	var farmer models.Farmer
	id := c.Param("id")

	db := c.MustGet("db").(*gorm.DB)
	if err := db.First(&farmer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Farmer not found"})
		return
	}

	if err := db.Delete(&farmer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete farmer"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
