package controllers

import (
	"net/http"
	"path/filepath"
	"pointafam/backend/models"

	"github.com/gin-gonic/gin"
)

func UpdateUserProfile(c *gin.Context) {
	var user models.User
	if err := db.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	var updateUserInput struct {
		Username    string `json:"username"`
		PhoneNumber string `json:"phoneNumber"`
		Location    string `json:"location"`
	}

	if err := c.ShouldBindJSON(&updateUserInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	user.Username = updateUserInput.Username
	user.PhoneNumber = updateUserInput.PhoneNumber
	user.Location = updateUserInput.Location

	// Handle profile picture upload
	file, err := c.FormFile("profile_picture")
	if err == nil {
		// Save the file to the server
		filename := filepath.Base(file.Filename)
		filepath := filepath.Join("uploads", filename)
		if err := c.SaveUploadedFile(file, filepath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to upload profile picture"})
			return
		}
		user.ProfilePictureUrl = "/" + filepath
	} else {
		// Set default avatar if no profile picture is uploaded
		defaultAvatarUrl := "/uploads/default_avatar.png"
		user.ProfilePictureUrl = defaultAvatarUrl
	}

	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}
