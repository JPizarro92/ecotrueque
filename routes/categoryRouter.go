package routes

import (
	"ecotrueque/controllers"
	"ecotrueque/middleware"

	"github.com/gin-gonic/gin"
)

func CategoryRoutes(router *gin.Engine) {
	router.POST("/api/category", middleware.Authenticate(), controllers.CreateCategory())
	router.GET("/api/category/:id", controllers.ReadCategory())
	router.GET("/api/category", controllers.ReadCategories())
	router.PUT("/api/category/:id", middleware.Authenticate(), controllers.UpdateCategory())
}
