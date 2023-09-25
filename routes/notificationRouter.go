package routes

import (
	"ecotrueque/controllers"
	"ecotrueque/middleware"

	"github.com/gin-gonic/gin"
)

func NotificationsRoutes(router *gin.Engine) {
	router.POST("/api/notifications", middleware.Authenticate(), controllers.CreateNotification())
	router.GET("/api/notifications/:id", controllers.ReadNotification())
	router.GET("/api/notifications", controllers.ReadNotifications())
	router.PUT("/api/notifications/:id", middleware.Authenticate(), controllers.UpdateNotification())
	router.DELETE("/api/notifications/:id", middleware.Authenticate(), controllers.DeleteNotification())
}
