package controllers

import (
	"net/http"

	"ecotrueque/helpers"
	"ecotrueque/initializers"
	"ecotrueque/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
 */

func CreateExchange() gin.HandlerFunc {
	return func(c *gin.Context) {
		var exchange models.Exchange
		if err := c.ShouldBindJSON(&exchange); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validate.Struct(exchange); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if err := initializers.DB.Create(&exchange).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error al crear trueque"})
			return
		}
		c.JSON(http.StatusCreated, exchange)
	}
}

/*
Todo: Read types
* 1. Lectura de 1 item
* 2. Lectura de trueques:
*		- por usuarios
*		- enviados
*		- recibidos
*		- aceptados
*/
func ReadExchange() gin.HandlerFunc {
	return func(c *gin.Context) {
		id_exchange := c.Param("id")

		var exchange models.Exchange

		if err := initializers.DB.First(&exchange, id_exchange).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, exchange)
	}
}

func ReadSendExchangesByUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		uid, _ := uuid.Parse(userID)

		var exchanges []models.Exchange

		if err := initializers.DB.Preload("ProposedUser").Preload("ProposedPost").Preload("ProposedPost.Images").Preload("PostUser").Preload("PostShared").Preload("PostShared.Images").Limit(10).Find(&exchanges, "proposed_user_id=?", uid).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving exchanges"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"exchanges": exchanges})
	}
}

func ReadReceivedExchangesByUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		uid, _ := uuid.Parse(userID)

		var exchanges []models.Exchange

		if err := initializers.DB.Preload("ProposedUser").Preload("ProposedPost").Preload("ProposedPost.Images").Preload("PostUser").Preload("PostShared").Preload("PostShared.Images").Limit(10).Find(&exchanges, "post_user_id=?", uid).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving exchanges"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"exchanges": exchanges})
	}
}

func UpdateExchange() gin.HandlerFunc {
	return func(c *gin.Context) {
		exchange_id := c.Param("id")
		var exchange models.Exchange
		var exchange_updated models.Exchange
		if err := c.ShouldBind(&exchange_updated); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "exchange received is empty"})
			return
		}

		if err := validate.Struct(exchange_updated); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if err := initializers.DB.First(&exchange, exchange_id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "exchange not found"})
			return
		}

		exchange.Message = exchange_updated.Message
		exchange.Observations = exchange_updated.Observations
		exchange.Status = exchange_updated.Status
		exchange.ProposedPostID = exchange_updated.ProposedPostID

		if err := initializers.DB.Save(&exchange).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "fail save updated"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": true})
	}
}

func DeleteExchange() gin.HandlerFunc {
	return func(c *gin.Context) {
		exchange_id := c.Param("id")
		uid_string := helpers.GetUidString(c)
		var exchange models.Exchange

		if err := initializers.DB.First(&exchange, exchange_id).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if uid_string != exchange.ProposedUserID.String() {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization"})
			return
		}

		if err := initializers.DB.Delete(&exchange, exchange_id).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "exchange deleted"})
	}
}
