package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"pointafam/backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UploadProfilePicture(c *gin.Context) {
	userID := c.Param("id")

	file, err := c.FormFile("profile_picture")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}

	// Create the uploads directory if it doesn't exist
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", os.ModePerm)
	}

	// Save the file to the uploads directory
	filePath := filepath.Join("uploads", filepath.Base(file.Filename))
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Update the user's profile picture URL in the database
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.ProfilePictureUrl = "/" + filePath
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile picture"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"profilePictureUrl": user.ProfilePictureUrl})
}

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
