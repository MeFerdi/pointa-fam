package controllers

import (
	"net/http"
	"pointafam/backend/models"

	"github.com/gin-gonic/gin"
)

func GetFarmers(c *gin.Context) {
	var farmers []models.Farmer
	db.Find(&farmers)
	c.HTML(http.StatusOK, "farmers_list.html", gin.H{"farmers": farmers})
}

func CreateFarmer(c *gin.Context) {
	var farmer models.Farmer
	if err := c.ShouldBind(&farmer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&farmer)
	GetFarmers(c)
}

func UpdateFarmer(c *gin.Context) {
	var farmer models.Farmer
	id := c.Param("id")
	if err := db.Where("id = ?", id).First(&farmer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Farmer not found"})
		return
	}
	if err := c.ShouldBindJSON(&farmer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	idInt := 0
	farmer.ID = uint(idInt)
	if err := db.Save(&farmer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Farmer updated successfully"})
}

func DeleteFarmer(c *gin.Context) {
	var farmer models.Farmer
	id := c.Param("id")
	if err := db.Where("id = ?", id).First(&farmer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Farmer not found"})
		return
	}
	db.Delete(&farmer)
	GetFarmers(c)
}
