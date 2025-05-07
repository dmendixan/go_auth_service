package handlers

import (
	"auth-service/config"
	"auth-service/internal/models"
	"gorm.io/gorm"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Profile(c *gin.Context, db *gorm.DB) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
	})
}
