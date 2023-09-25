package main

import (
	"ecotrueque/controllers"
	"ecotrueque/initializers"
	"ecotrueque/routes"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	fmt.Print("hola")
	location, err := time.LoadLocation("America/Guayaquil")
	if err != nil {
		// Manejar el error
		panic(err)
	}

	// Establecer la zona horaria por defecto
	time.Local = location

	router := gin.New()
	router.Static("/assets", "./assets")
	router.MaxMultipartMemory = 8 << 20 //8MiB
	//router.Use(middleware.CORSMiddleware()) // Todo: configuración de los CORS

	//** Rutas del servicios
	router.GET("/api/seed", controllers.CreateSeed)
	routes.AuthRoutes(router)
	routes.UserRoutes(router)          // Todo: Rutas CRUD users
	routes.PostRoutes(router)          // Todo: Rutas CRUD posts
	routes.CategoryRoutes(router)      // Todo: Rutas CRUD categorias
	routes.ExchangesRoutes(router)     // Todo: Rutas CRUD exchanges
	routes.NotificationsRoutes(router) // Todo: Rutas CRUD notifications
	router.Run()
}
