package controllers

import (
	"net/http"
	"pointafam/backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeleteUser removes a user by ID from the database.
func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32) // Convert string ID to uint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := models.DeleteUser(db, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete user"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func UpdateUserProfile(c *gin.Context) {
    var user models.User
    if err := db.First(&user, c.Param("id")).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
        return
    }

    var input struct {
        FirstName   string `json:"first_name"`
        LastName    string `json:"last_name"`
        Email       string `json:"email"`
        PhoneNumber string `json:"phone_number"`
        Location    string `json:"location"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
        return
    }

    user.FirstName = input.FirstName
    user.LastName = input.LastName
    user.Email = input.Email
    user.PhoneNumber = input.PhoneNumber
    user.Location = input.Location

    if err := db.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user"})
        return
    }

    c.JSON(http.StatusOK, user)
}
