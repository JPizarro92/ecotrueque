package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func DeletePostImageInTemp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		/*
			1 recibir nombre de imagen
			2 verificar en carpeta
			3 eliminar imagen en carpeta
			4 devolver true or false
		*/
	}
}

func UploadPostImage2Temp() gin.HandlerFunc {
	return func(c *gin.Context) {
		//? Verifica que se envíe el archivo y la propiedad "image"
		file, err := c.FormFile("imagen")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"code":    http.StatusBadRequest,
				"message": "No se proporcionó ningún archivo.",
				"error":   err.Error(),
			})
			return
		}

		//? Obtener User contexto validado por el token
		userGet, _ := c.Get("uid")
		user, ok := userGet.(string)
		if !ok {
			if userStr, ok := userGet.(fmt.Stringer); ok {
				user = userStr.String()
			} else {
				fmt.Println("Error al transformar usuario. Módulo Post")
			}
		}

		//? Crear Carpetas y obtener el path
		path := MakeFolderUser(user)

		//?Guardar el archivo en el sistema de archivos
		err = c.SaveUploadedFile(file, filepath.Join(path, "/temp/", file.Filename))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error al guardar el archivo",
			})
			return
		}

	}
}

func DeleteImages() {

}

func UpdateImage() {

}
