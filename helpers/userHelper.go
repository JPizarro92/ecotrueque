package helpers

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetUidString(ctx *gin.Context) (uid string) {
	uidGet, _ := ctx.Get("uid")
	uid, ok := uidGet.(string)
	if !ok {
		if uidStr, ok := uidGet.(fmt.Stringer); ok {
			uid = uidStr.String()
		}
	}
	return uid
}

// Todo: function for convert the password string to a hash
func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err.Error())
	}
	return string(hash)
}

// Todo: function for validate the login password with the user hash
func VerifyPassword(hashPassword string, password string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	check := true
	msg := ""
	if err != nil {
		fmt.Println(err.Error())
		msg = fmt.Sprintf("email or password is incorrect")
		check = false
	}
	return check, msg
}

// Todo: function for validate the same UID
func VerifySameUid(ctx *gin.Context, uid_user string) bool {
	uid_string := GetUidString(ctx)
	if uid_string != uid_user {
		return false
	}
	return true
}
