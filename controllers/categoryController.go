package controllers

import (
	"net/http"

	"ecotrueque/initializers"
	"ecotrueque/models"

	"github.com/gin-gonic/gin"
)

func CreateCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		//* Obtener category
		var category models.Category

		if err := c.ShouldBind(&category); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//** Validated data input
		validationErr := validate.Struct(category)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		//* Obtención de datos respecto al usuario quien realiza la actividad
		userId, _ := c.Get("uid")
		var user models.User
		err := initializers.DB.First(&user, userId).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if user.Role != "ADMIN" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No authorization"})
			return
		}

		if err = initializers.DB.Create(&category).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusAccepted, category)

	}
}
func ReadCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		//* Obtener datos relacionados respecto a la categoría
		categoryID := c.Param("id")
		var category models.Category

		if err := initializers.DB.Preload("Posts").First(&category, categoryID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, category)

	}
}

func ReadCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		var categories []models.Category
		if err := initializers.DB.Find(&categories).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error retrieving categories",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"categories": categories})
	}
}

func UpdateCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryID := c.Param("id")
		var category models.Category

		if err := initializers.DB.First(&category, categoryID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if category.ID == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "category not found"})
			return
		}
		//* Obtención de datos respecto al usuario quien realiza la actividad
		userId, _ := c.Get("uid")
		var user models.User

		if err := initializers.DB.First(&user, userId).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if user.Role != "ADMIN" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No authorization"})
			return
		}

		//* Get category updated
		var input models.Category
		if err := c.ShouldBind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		category.Name = input.Name

		if err := initializers.DB.Save(&category).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, category)
	}
}

func DeleteCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryID := c.Param("id")
		var category models.Category

		if err := initializers.DB.First(&category, categoryID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if category.ID == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "category not found"})
			return
		}
		//* Obtención de datos respecto al usuario quien realiza la actividad
		userId, _ := c.Get("uid")
		var user models.User

		if err := initializers.DB.First(&user, userId).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if user.Role != "ADMIN" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No authorization"})
			return
		}

		//* Delete Category
		if err := initializers.DB.Delete(&category).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": true})
	}
}
