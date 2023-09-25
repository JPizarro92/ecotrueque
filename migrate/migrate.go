package main

import (
	"ecotrueque/helpers"
	"ecotrueque/initializers"
	"ecotrueque/models"
	"fmt"
	"log"
	"os"

	"gorm.io/gorm"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	production := os.Getenv("PRODUCTION")

	if production == "false" {
		initializers.DB.Migrator().DropTable(&models.Exchange{})
		initializers.DB.Migrator().DropTable(&models.Rating{})
		initializers.DB.Migrator().DropTable(&models.Category{})
		initializers.DB.Migrator().DropTable(&models.PostImage{})
		initializers.DB.Migrator().DropTable(&models.Post{})
		initializers.DB.Migrator().DropTable(&models.User{})
	}

	initializers.DB.AutoMigrate(&models.User{})
	initializers.DB.AutoMigrate(&models.Rating{})
	initializers.DB.AutoMigrate(&models.Post{})
	initializers.DB.AutoMigrate(&models.PostImage{})
	initializers.DB.AutoMigrate(&models.Category{})
	initializers.DB.AutoMigrate(&models.Exchange{})

	createAdminUser()

}

func createAdminUser() {

	var admin models.User

	result := initializers.DB.First(&admin, "role=?", "ADMIN")

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.Fatal(result.Error)
	}

	if result.Error == gorm.ErrRecordNotFound {
		user_admin := models.User{
			Name:     "admin",
			Surname:  "admin",
			Email:    "admin@ecotrueque.ups.com",
			Password: helpers.HashPassword("admin_ecotrueque"),
			Verified: true,
			Role:     "ADMIN",
			Avatar:   "avatar.png",
			Status:   "ACTIVE",
		}
		if err := initializers.DB.Create(&user_admin).Error; err != nil {
			log.Fatal(err)
		}
		fmt.Println("Usuario 'ADMIN' Creado.")
		fmt.Println(user_admin)
	} else {
		fmt.Println("Usuario 'ADMIN' Existente.")
	}

}
