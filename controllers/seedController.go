package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"ecotrueque/helpers"
	"ecotrueque/initializers"
	"ecotrueque/models"

	"github.com/gin-gonic/gin"
)

var users []models.User
var categories []models.Category
var posts []models.Post

func CategorySeed() {
	//* Leer el archivo JSON
	fileBytes, err := ioutil.ReadFile("./controllers/data_categories.json")
	if err != nil {
		log.Fatal(err)
	}
	//* decodificar JSON a array de estructura
	err = json.Unmarshal(fileBytes, &categories)
	if err != nil {
		log.Fatal(err)
	}
	//* Agregar en DB
	for _, category := range categories {
		err = initializers.DB.Create(&category).Error
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}
}

func UserSeed() {
	//* Leer el archivo JSON
	fileBytes, err := ioutil.ReadFile("./controllers/data_users.json")
	if err != nil {
		log.Fatal(err)
	}
	//* decodificar JSON a array de estructura
	err = json.Unmarshal(fileBytes, &users)
	if err != nil {
		log.Fatal(err)
	}
	//* Agregar en DB
	for _, user := range users {
		user.Password = helpers.HashPassword(user.Password)
		user.Avatar = "av-1.png"
		err = initializers.DB.Create(&user).Error
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}
}

func PostSeed() {
	//* Leer el archivo JSON
	fileBytes, err := ioutil.ReadFile("./controllers/data_posts.json")
	if err != nil {
		log.Fatal(err)
	}
	//* decodificar JSON a array de estructura
	err = json.Unmarshal(fileBytes, &posts)
	if err != nil {
		log.Fatal(err)
	}
	//* Agregar en DB
	for _, post := range posts {

		var user models.User
		var category models.Category

		err = initializers.DB.First(&user, post.UserID).Error
		if err != nil {
			fmt.Println("Error User: ", err.Error())
		}

		err = initializers.DB.First(&category, post.CategoryID).Error
		if err != nil {
			fmt.Println("Error Category: ", err.Error())
		}

		post.User = &user
		post.Category = &category

		err = initializers.DB.Create(&post).Error
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}
}

func CreateSeed(c *gin.Context) {
	CategorySeed()
	UserSeed()
	//PostSeed()
	c.JSON(http.StatusAccepted, gin.H{
		"users":      users,
		"categories": categories,
		"posts":      posts,
	})

}
