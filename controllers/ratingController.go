package controllers

import (
	"ecotrueque/initializers"
	"ecotrueque/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

/*
* ABOUT THIS CONTROLLER
* The rating controller allows adding new rating from a user to other.
* Also it has a func for get the rating value.
 */

func AddRaitingToUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var rating models.Rating

		if err := ctx.ShouldBindJSON(&rating); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Undetermined data",
				"error":   err.Error()})
			return
		}

		//* Validated data input
		if err := validate.Struct(rating); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Check the required fields",
				"error":   err.Error()})
			return
		}
		//* Calculate the RatingScore
		var user models.User
		if err := initializers.DB.First(&user, "id=?", rating.RatedID).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "Error loading user",
				"error":   err.Error()})
		}
		CalculateRatingScore(user)
		//* Save the RatingScore in the user
		if err := initializers.DB.Save(&user).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "Error saving user with his RatingScore",
				"error":   err.Error()})
		}

		//* Save the rating
		if err := initializers.DB.Save(&rating).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "Error saving rating in DB",
				"error":   err.Error()})
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"status":  "success",
			"message": "Rating registred successfully",
			"rating":  rating})

	}
}

func getRating() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		var rating models.Rating
		if err := initializers.DB.First(&rating, "id=?", id).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "Error finding a record",
				"error":   err.Error()})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Rating found",
			"rating":  rating})

	}
}

func getRatings() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		page := 1
		perPage := 10
		var ratings []models.Rating

		if ctx.Query("page") != "" {
			page, _ = strconv.Atoi(ctx.Query("page"))
		}

		if ctx.Query("perPage") != "" {
			perPage, _ = strconv.Atoi(ctx.Query("perPage"))
		}

		offset := (page - 1) * perPage

		if err := initializers.DB.Offset(offset).Limit(10).Find(&ratings).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "Error retrieving ratings",
				"error":   err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Data loaded",
			"ratings": ratings})
	}
}

func UpdateRating() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var payload models.Rating

		if err := ctx.ShouldBindJSON(&payload); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Undetermined data",
				"error":   err.Error()})
			return
		}

		id_rating := ctx.Param("id")
		var rating models.Rating

		if err := initializers.DB.First(&rating, "id=?", id_rating).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "Rating searching error",
				"error":   err.Error()})
			return
		}

		rating.Value = payload.Value
		rating.UpdatedAt = time.Now()

		if err := initializers.DB.Save(&rating).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "Rating searching error",
				"error":   err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":   "success",
			"menssage": "Updates perfomed",
			"rating":   rating})

	}
}

func DeleteRating() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id_rating := ctx.Param("id")

		var rating models.Rating

		if err := initializers.DB.First(&rating, "id=?", id_rating).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "Rating searching error",
				"error":   err.Error()})
			return
		}

		if err := initializers.DB.Delete(&rating).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "Rating deleting error",
				"error":   err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "succes",
			"message": "Rating deleted successfully",
			"rating":  rating})

	}
}

func CalculateRatingScore(u models.User) {
	if len(u.Ratings) == 0 {
		u.RatingScore = 0.0
		return
	}

	var totalScore float64
	for _, rating := range u.Ratings {
		totalScore += float64(rating.Value)
	}

	u.RatingScore = totalScore / float64(len(u.Ratings))

}
