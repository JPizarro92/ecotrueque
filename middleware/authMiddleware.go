package middleware

import (
	"fmt"
	"net/http"

	"ecotrueque/helpers"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provied")})
			c.Abort()
			return
		}

		user, err := helpers.ValidateToken(clientToken)

		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", user.Email)
		c.Set("name", user.Name)
		c.Set("surname", user.Surname)
		c.Set("uid", user.ID)
		c.Set("role", user.Role)
		c.Next()

	}
}
