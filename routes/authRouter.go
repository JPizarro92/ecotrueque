package routes

import (
	"ecotrueque/controllers"
	"ecotrueque/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	// ** Rutas Login - SignUp
	router.POST("/api/login", controllers.SignInUser()) //? Login de Usuario
	router.POST("/api/signup", controllers.SignUpUser())
	router.GET("/api/validate", middleware.Authenticate(), controllers.ValidateLogin())
}
