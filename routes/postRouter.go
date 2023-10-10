package routes

import (
	"ecotrueque/controllers"
	"ecotrueque/middleware"

	"github.com/gin-gonic/gin"
)

func PostRoutes(router *gin.Engine) {
	//Todo: rutas para gestión de posts
	router.POST("/api/posts/upload-post-image", middleware.Authenticate(), controllers.UploadPostImage2Temp())
	router.DELETE("/api/posts/delete-post-image/:id", middleware.Authenticate(), controllers.DeleteImageInPostFolder())

	router.POST("/api/posts", middleware.Authenticate(), controllers.CreatePost())
	router.GET("/api/posts/:id", controllers.ReadPost())
	router.GET("/api/posts/user/:id", controllers.ReadUserPosts())
	router.GET("/api/posts/user", middleware.Authenticate(), controllers.ReadUserPostsSignIn())
	router.GET("/api/posts/category/:id", controllers.ReadPostsByCategory())
	router.PUT("/api/posts/:id", middleware.Authenticate(), controllers.UpdatePost())
	router.DELETE("/api/posts/:id", middleware.Authenticate(), controllers.DeletePost())
}
