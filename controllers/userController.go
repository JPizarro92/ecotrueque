package controllers

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"ecotrueque/helpers"
	"ecotrueque/initializers"
	"ecotrueque/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New()

// todo: this function only activated for user admin
func CreateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//* Get User Login
		var uid_string = helpers.GetUidString(ctx)
		var userLogin models.User
		if err := initializers.DB.First(&userLogin, "id=?", uid_string).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "error to charge user login",
				"error":   err.Error()})
			return
		}

		//* Verify User Role
		if err := helpers.CheckUserType(ctx, userLogin.Role); err != nil || userLogin.Role != "ADMIN" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": err,
				"error":   ""})
			return
		}

		//* Get New User Data
		var payload models.User
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Undetermined data",
				"error":   err.Error()})
			return
		}

		//* Validated data input
		if err := validate.Struct(payload); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": "Check the required fields",
				"error":   err.Error()})
			return
		}

		user := models.User{
			Name:     payload.Name,
			Surname:  payload.Surname,
			Email:    strings.ToLower(payload.Email),
			Password: helpers.HashPassword(payload.Password),
			Verified: false,
			Role:     payload.Role,
			Avatar:   "avatar.png",
			Status:   "ACTIVE",
		}

		if err := initializers.DB.Create(&user).Error; err != nil && strings.Contains(err.Error(), "duplicate key value violates unique") {
			ctx.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": "User with that email already exists",
				"error":   err.Error()})
			return
		} else if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{
				"status":  "error",
				"message": "Something bad happened",
				"error":   ""})
			return
		}

		token, refreshToken, _ := helpers.GenerateAllTokens(user.Email, user.Name, user.Surname, user.Role, user.ID.String())
		user.Token = token
		user.Refresh_token = refreshToken

		if err := initializers.DB.Save(&user).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": "Error saving user.",
				"error":   err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"status":  "success",
			"message": "User created successfully",
			"user":    user})

	}
}

// Todo: function for read the user with his post and exchanges
func ReadUserByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//** Get id from header
		id := ctx.Param("id")
		//** Find the user by id
		var user models.User
		if err := initializers.DB.Preload("Posts").Preload("ExchangesSend").Preload("ExchangesReceived").First(&user, "id=?", id).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "User not found",
				"error":   err.Error()})
			return
		}

		// ** Validated User
		if err := helpers.MatchUserRoleToUid(ctx, user.ID.String()); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{
				"status":  "error",
				"message": err.Error(),
				"error":   ""})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "User charged successfully",
			"user":    user})
	}
}

// Todo: function enable for User Admin, read all users
func ReadUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := helpers.CheckUserType(ctx, "ADMIN"); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
				"error":   ""})
			return
		}
		page := 1     //* Página predeterminada
		perPage := 10 //* Cantidad de posts por página
		var users []models.User

		//* Obtener los parámetros de consulta de la URL
		if ctx.Query("page") != "" {
			page, _ = strconv.Atoi(ctx.Query("page"))
		}
		if ctx.Query("perPage") != "" {
			perPage, _ = strconv.Atoi(ctx.Query("perPage"))
		}

		//* Calcular el desplazamiento (offset) según el número de página y la cantidad de posts por página
		offset := (page - 1) * perPage

		//* Obtener Usuarios
		if err := initializers.DB.Offset(offset).Limit(10).Find(&users).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "Error retrieving users",
				"error":   err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Data loaded",
			"users":   users})
	}
}

// Todo: function for Update a user
func UpdateUserData() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id := ctx.Param("id")
		uid, _ := uuid.Parse(id)

		var userLog models.User

		if err := initializers.DB.First(&userLog, uid).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "Userl loading error",
				"error":   err.Error()})
			return
		}

		//** Validated User Role
		if err := helpers.MatchUserRoleToUid(ctx, userLog.ID.String()); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{
				"status":  "error",
				"message": err.Error(),
				"error":   ""})
			return
		}

		//** Get User Update
		var input models.User
		if err := ctx.ShouldBindJSON(&input); err != nil { // Obtener los datos de usuario actualizados del cuerpo de la solicitud
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "You must enter all required field",
				"error":   err.Error()})
			return
		}
		//* Get User
		var user models.User
		if err := initializers.DB.First(&user, uid).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "User not found",
				"error":   err.Error()})
			return
		}

		//* Update it
		user.Name = input.Name
		user.Surname = input.Surname
		user.Phone = input.Phone
		user.BirthDate = input.BirthDate
		user.Age = input.Age

		//* Save user
		if err := initializers.DB.Save(&user).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "User saving error",
				"error":   err.Error()})
			return
		}

		//* Respond
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "User updated",
			"user":    user})

	}
}

type Password struct {
	Password string `json:"password" binding:"required,min=8"`
}

func UpdateUserPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//* Recibir la contraseña nueva
		id := ctx.Param("id")

		var payload Password

		if err := ctx.ShouldBindJSON(&payload); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Data not received",
				"error":   err.Error()})
			return
		}

		var user models.User
		if err := initializers.DB.First(&user, "id=?", id).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "User not found error",
				"error":   err.Error()})
			return
		}

		if !helpers.VerifySameUid(ctx, id) && user.Role == "USER" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Unauthorized to access to this source.",
				"error":   ""})
			return
		}

		user.Password = helpers.HashPassword(payload.Password)
		if err := initializers.DB.Save(&user).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "Password saving error",
				"error":   err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Password Updated"})
	}
}

func DeleteUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id := ctx.Param("id")
		var user models.User

		//* Buscar el usuario en la base de datos por ID
		err := initializers.DB.First(&user, "id=?", id).Error
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "fail",
				"message": "User not found",
				"error.":  err.Error()})
			return
		}

		role, _ := ctx.Get("role")

		if role != "ADMIN" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Unauthorized to access to this source",
				"error":   ""})
			return
		}

		err = initializers.DB.Delete(&user).Error
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "fail",
				"message": "User deleting error",
				"error":   err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "User deleted successfully",
		})

	}
}

/*
! Sobre almacenamiento de imagenes investigar acerca del manejo de memoría en disco cuando
! el servicio esta dockerizado. Aquí se tiene las siguientes opciones:
! Mantener cargando imagen en disco en el caso de ser dinámico, sino ver otro servidor de
! almacenamiento. Caso contrario almecenar en la DB en forma de texto (verificar si existe binarizado).
*/
func UploadAvatar() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//* Get image from context with name "avatar"
		file, err := ctx.FormFile("avatar")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "No file was provided",
				"error":   err.Error()})
			return
		}

		//* The avatar must be save with the user id
		uid_string := helpers.GetUidString(ctx)
		imgName := uid_string + filepath.Ext(file.Filename)
		savePath := filepath.Join("assets/img/avatars/", imgName)

		err = ctx.SaveUploadedFile(file, savePath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "File saving error",
				"error":   err.Error(),
			})
			return
		}

		var user models.User
		err = initializers.DB.First(&user, "id=?", uid_string).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "User not found error",
				"error":   err.Error()})
			return
		}

		//* Save name avatar (add the time with the id?)
		user.Avatar = imgName
		err = initializers.DB.Save(&user).Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Avatar saving error",
				"error":   err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Successfully saved file",
			"avatar":  savePath,
		})
	}
}
