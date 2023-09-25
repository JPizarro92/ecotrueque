package controllers

import (
	"net/http"
	"strings"
	"time"

	"ecotrueque/helpers"
	"ecotrueque/initializers"
	"ecotrueque/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

// Todo: SignIn or Login of users
func SignInUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//* Get data from context
		var payload models.SignInInput
		if err := ctx.BindJSON(&payload); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
			return
		}

		//* Search user in DB
		var user models.User
		if err := initializers.DB.First(&user, "email = ?", payload.Email).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "email or password is incorrect"})
			return
		}

		//* Validate password with hash
		passwordIsValid, msg := helpers.VerifyPassword(user.Password, payload.Password)
		if !passwordIsValid {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": msg})
			return
		}

		//* Generate the token
		token, refreshToken, _ := helpers.GenerateAllTokens(user.Email, user.Name, user.Surname, user.Role, user.ID.String())
		user.Token = token
		user.Refresh_token = refreshToken
		user.UpdatedAt = time.Now()

		if err := initializers.DB.Save(&user).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "success", "user": user})
	}
}

// Todo: SignUp or Register of users
func SignUpUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload models.SignUpInput
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		user := models.User{
			Name:     payload.Name,
			Surname:  payload.Surname,
			Email:    strings.ToLower(payload.Email),
			Password: helpers.HashPassword(payload.Password),
			Verified: false,
			Role:     "USER",
			Avatar:   "avatar.png",
			Status:   "ACTIVE",
		}

		if err := initializers.DB.Create(&user).Error; err != nil && strings.Contains(err.Error(), "duplicate key value violates unique") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email already exists"})
			return
		} else if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Something bad happened"})
			return
		}

		// ! Agregar código para verificar por código de correo electrónico ó teléfono
		// ! Validación de dos pasos
		token, refreshToken, _ := helpers.GenerateAllTokens(user.Email, user.Name, user.Surname, user.Role, user.ID.String())

		user.Token = token
		user.Refresh_token = refreshToken

		if err := initializers.DB.Save(&user).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Error generating token."})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"status": "success", "user": user})

	}
}

// Todo: Validate if user is login or not
func ValidateLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid_string := helpers.GetUidString(ctx)
		var user models.User
		if err := initializers.DB.First(&user, "id=?", uid_string).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "user": user})
	}
}

func IsDuplicateKeyError(err error) bool {
	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		return false
	}
	return pgErr.Code == "23505" // Todo: código de error específico de PostgreSQL para duplicados de clave única
}
