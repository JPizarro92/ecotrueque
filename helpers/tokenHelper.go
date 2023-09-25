package helpers

import (
	"fmt"
	"log"
	"os"
	"time"

	"ecotrueque/initializers"
	"ecotrueque/models"

	"github.com/golang-jwt/jwt/v5"
)

var SECRET_KEY = os.Getenv("SECRET")

func GenerateAllTokens(email string, name string, surname string, role string, userId string) (signedToken string, signedRefreshToken string, err error) {

	fmt.Println("TH id: ", userId)
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
	}).SignedString([]byte(SECRET_KEY))

	tokenRefresh, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
	}).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, tokenRefresh, err
}

func ValidateToken(signedToken string) (user models.User, msg string) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	token, _ := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(SECRET_KEY), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			msg = fmt.Sprintf("the token is expired")
			return
		}

		// Find the user with token sub
		//! err := initializers.DB.First(&user, "user_id = ?", claims["sub"]).Error
		//userID := claims["sub"].(string)
		err := initializers.DB.First(&user, "id=?", claims["sub"]).Error
		if err != nil {
			msg = fmt.Sprintf("user token is invalid")
			msg = fmt.Sprintf(err.Error())
			return
		}

	} else {
		msg = "token is invalid"
		return
	}

	return user, msg
}
