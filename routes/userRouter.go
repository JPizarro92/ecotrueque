package routes

import (
	"ecotrueque/controllers"
	"ecotrueque/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	// ** Rutas de Usuario
	router.POST("/api/user", middleware.Authenticate(), controllers.CreateUser())
	router.GET("/api/user/:id", middleware.Authenticate(), controllers.ReadUserByID())
	router.GET("/api/user", middleware.Authenticate(), controllers.ReadUsers())
	router.PUT("/api/user/data/:id", middleware.Authenticate(), controllers.UpdateUserData())
	router.PUT("/api/user/password/:id", middleware.Authenticate(), controllers.UpdateUserPassword())
	router.DELETE("/api/user/:id", middleware.Authenticate(), controllers.DeleteUser())
	router.POST("/api/user/avatar", middleware.Authenticate(), controllers.UploadAvatar())
}
