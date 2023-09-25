package routes

import (
	"ecotrueque/controllers"
	"ecotrueque/middleware"

	"github.com/gin-gonic/gin"
)

func ExchangesRoutes(router *gin.Engine) {
	router.POST("/api/exchanges", middleware.Authenticate(), controllers.CreateExchange())
	router.GET("/api/exchanges/:id", middleware.Authenticate(), controllers.ReadExchange())
	router.GET("/api/exchanges/s/user/:id", middleware.Authenticate(), controllers.ReadSendExchangesByUser())
	router.GET("/api/exchanges/r/user/:id", middleware.Authenticate(), controllers.ReadReceivedExchangesByUser())
	router.PUT("/api/exchanges/:id", middleware.Authenticate(), controllers.UpdateExchange())
	router.DELETE("/api/exchanges/:id", middleware.Authenticate(), controllers.DeleteExchange())
}
