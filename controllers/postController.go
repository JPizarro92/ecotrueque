package controllers

import (
	"errors"
	"fmt"

	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"ecotrueque/helpers"
	"ecotrueque/initializers"
	"ecotrueque/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var path_img = os.Getenv("PATH_IMG")

func CreatePost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var post models.Post
		if err := ctx.ShouldBindJSON(&post); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Undetermined data",
				"error":   err.Error()})
			return
		}
		validate := validator.New()
		if err := validate.Struct(post); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Check all required fields",
				"error":   err.Error()})
			return
		}
		uid := helpers.GetUidString(ctx)
		if uid != post.User.ID.String() {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Unauthorized to access this resource",
				"error":   "inconsistent user is trying to register a post"})
			return
		}

		var user models.User
		if err := initializers.DB.First(&user, uid).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "User not found",
				"error":   err.Error()})
			return
		}

		//* Crear los PostImage
		var images []models.PostImage

		// Todo: for external storage we must get the all imagen urls
		filestxt, _ := LoadPostImagesToFolder(uid)
		for _, name := range filestxt {
			var image models.PostImage
			image.Img = name
			images = append(images, image)
		}

		post.Images = images

		//* Crear el Post
		if err := initializers.DB.Create(&post).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error ...": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"status": "success", "post": post})
	}
}

func ReadPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var post models.Post
		err := initializers.DB.Preload("User").Preload("Category").Preload("Images").First(&post, id).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if post.ID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "post no found"})
			return
		}

		c.JSON(http.StatusOK, post)
	}
}

func UpdatePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		//* Obtenemos el ID del post a actualizar
		postID := c.Param("id")

		var post models.Post

		if err := initializers.DB.First(&post, postID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		//* Obtenemos el ID del User autenticado
		userID := helpers.GetUidString(c)
		var user models.User

		if err := initializers.DB.First(&user, "id=?", userID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("Aqui :vd")
		if user.ID != post.UserID {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No authorization"})
			return
		}

		//* Obtener el post updated
		var input models.Post
		if err := c.ShouldBind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		post.Title = input.Title
		post.Price = input.Price
		post.ShortDescription = input.ShortDescription
		post.LongDescription = input.LongDescription
		post.ExchangeRate = input.ExchangeRate
		post.Tags = input.Tags
		post.ProductStatus = input.ProductStatus
		post.PostStatus = input.PostStatus
		post.CategoryID = input.CategoryID
		post.Category = input.Category

		//** Save post updated
		if err := initializers.DB.Save(&post).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, post)

	}
}

func DeletePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		//* Obtener el ID del User Auth
		userID := helpers.GetUidString(c)
		var user models.User
		if err := initializers.DB.First(&user, "id=?", userID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		//* post id
		postID := c.Param("id")
		var post models.Post
		if err := initializers.DB.First(&post, postID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if user.ID != post.UserID {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No authorization"})
			return
		}

		if err := initializers.DB.Delete(&post).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": true})

	}
}

func ReadUserPostsSignIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id, _ := c.Get("uid")
		userID, _ := user_id.(string)
		uid, err := uuid.Parse(userID)
		if err != nil {
			fmt.Println(uid)
		}

		page := 1     //* Página predeterminada
		perPage := 10 //* Cantidad de posts por página
		var posts []models.Post
		//* Obtener los parámetros de consulta de la URL
		if c.Query("page") != "" {
			page, _ = strconv.Atoi(c.Query("page"))
		}
		if c.Query("perPage") != "" {
			perPage, _ = strconv.Atoi(c.Query("perPage"))
		}

		//* Calcular el desplazamiento (offset) según el número de página y la cantidad de posts por página
		offset := (page - 1) * perPage

		if err := initializers.DB.Preload("User").Preload("Category").Preload("Images").Offset(offset).Limit(10).Find(&posts, "user_id=?", uid).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error retrieving posts",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"posts": posts})
	}
}

func ReadUserPosts() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		uid, err := uuid.Parse(userID)

		if err != nil {
			fmt.Println(uid)
		}

		page := 1     //* Página predeterminada
		perPage := 10 //* Cantidad de posts por página
		var posts []models.Post
		//* Obtener los parámetros de consulta de la URL
		if c.Query("page") != "" {
			page, _ = strconv.Atoi(c.Query("page"))
		}
		if c.Query("perPage") != "" {
			perPage, _ = strconv.Atoi(c.Query("perPage"))
		}

		//* Calcular el desplazamiento (offset) según el número de página y la cantidad de posts por página
		offset := (page - 1) * perPage

		if err := initializers.DB.Preload("User").Preload("Category").Preload("Images").Offset(offset).Limit(10).Find(&posts, "user_id=?", uid).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error retrieving posts",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"posts": posts})
	}
}

func ReadPostsByCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryId := c.Param("id")
		page := 1     //* Página predeterminada
		perPage := 10 //* Cantidad de posts por página
		var posts []models.Post

		//* Obtener los parámetros de consulta de la URL
		if c.Query("page") != "" {
			page, _ = strconv.Atoi(c.Query("page"))
		}
		if c.Query("perPage") != "" {
			perPage, _ = strconv.Atoi(c.Query("perPage"))
		}

		//* Calcular el desplazamiento (offset) según el número de página y la cantidad de posts por página
		offset := (page - 1) * perPage
		fmt.Println("Codigo categoria: ", categoryId)
		fmt.Println("Pagina :", page)
		if categoryId == "0" {
			if err := initializers.DB.Preload("User").Preload("Category").Preload("Images").Offset(offset).Limit(10).Find(&posts).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Error retrieving posts",
				})
				return
			}
		} else {
			if err := initializers.DB.Preload("User").Preload("Category").Preload("Images").Offset(offset).Limit(10).Find(&posts, "category_id=?", categoryId).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Error retrieving posts",
				})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"posts": posts})
	}
}

// Todo: function for validated the image extension
func ValidatedExtension(filename string) bool {
	//Validar la extensión del archivo
	allowedExtensions := []string{".jpg", ".jpeg", ".png"}
	extension := filepath.Ext(filename)
	isAllowed := false
	for _, ext := range allowedExtensions {
		if ext == extension {
			isAllowed = true
			break
		}
	}
	return isAllowed
}

func MakeFolderUser(userID string) string {
	pathUser := fmt.Sprintf("%s%s", path_img, userID)
	pathUserTemp := fmt.Sprintf("%s%s%s", path_img, userID, "/temp")
	//? Verifica si existe la carpeta del usuario
	exist, _ := ValidateFolderExist(pathUser)
	//? Si no existe se crea las carpetas de usuario y temp
	if !exist {
		os.Mkdir(pathUser, 0755)
		os.Mkdir(pathUserTemp, 0755)
	}
	return pathUser
}

func LoadPostImagesToFolder(userID string) ([]string, error) {

	//* Move images from temp to posts
	pathTemp := fmt.Sprintf("%s%s%s", path_img, userID, "/temp")
	pathPost := fmt.Sprintf("%s%s%s", path_img, userID, "/posts")

	//* Verificar si existen las carpetas
	existPost, _ := ValidateFolderExist(pathPost)
	existTemp, _ := ValidateFolderExist(pathTemp)

	if !existPost && !existTemp {
		return nil, errors.New("error while trying to create folders")
	}

	//* Leer archivos y cargar los nombres
	var filestxt []string

	if !existTemp {
	} else {
		files, err := os.ReadDir(pathTemp)
		if err != nil {
			return nil, errors.New("error while reading the temporal folder")
		} else {
			for _, file := range files {
				if !file.IsDir() {
					filestxt = append(filestxt, file.Name())
					srcPath := filepath.Join(pathTemp, file.Name())
					destPath := filepath.Join(pathPost, file.Name())
					err := os.Rename(srcPath, destPath)
					if err != nil {
						fmt.Printf("Error al mover el archivo %s: %v", file.Name(), err)
					} else {
						fmt.Printf("Archivo %s movido éxitosamente\n", file.Name())
					}
				}
			}
		}
	}
	return filestxt, nil
}

// Todo: this function validates if the folder exist, otherwise the system creates its
func ValidateFolderExist(folderPath string) (bool, error) {
	_, err := os.Stat(folderPath)
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(folderPath, 0755)
			return true, nil
		}
		return false, err
	}
	return true, nil
}
