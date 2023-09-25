package controllers

import (
	"net/http"

	"ecotrueque/initializers"
	"ecotrueque/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func CreateNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
		var notification models.Notification
		if err := c.ShouldBind(&notification); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validate := validator.New()
		if err := validate.Struct(notification); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		role_authenticated, _ := c.Get("role")

		if role_authenticated != "ADMIN" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized to created"})
			return
		}

		if err := initializers.DB.Create(&notification).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusCreated, notification)

	}
}

func ReadNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
		notification_id := c.Param("id")

		var notification models.Notification
		if err := initializers.DB.First(&notification, notification_id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}

		c.JSON(http.StatusOK, notification)

	}
}

func ReadNotifications() gin.HandlerFunc {
	return func(c *gin.Context) {
		var notifications []models.Notification

		if err := initializers.DB.Find(&notifications).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"notifications": notifications})
	}
}

func UpdateNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
		notification_id := c.Param("id")
		var notification models.Notification
		var notification_updated models.Notification
		if err := initializers.DB.First(&notification, notification_id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}

		if err := c.ShouldBind(&notification_updated); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "notification received is empty"})
			return
		}

		notification.Name = notification_updated.Name
		notification.Description = notification_updated.Description
		notification.Status = notification_updated.Status
		notification.Image = notification_updated.Image

		if err := initializers.DB.Save(&notification).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}

		c.JSON(http.StatusOK, notification)

	}
}

func DeleteNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
		notification_id := c.Param("id")
		var notification models.Notification

		if err := initializers.DB.First(&notification, notification_id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}

		if err := initializers.DB.Delete(&notification).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": true})

	}
}
