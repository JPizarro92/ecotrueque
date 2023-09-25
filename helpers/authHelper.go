package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("role")
	err = nil
	if userType != role {
		err = errors.New("Unauthorized to access this resource")
	}
	return err
}

func MatchUserRoleToUid(ctx *gin.Context, userId string) (err error) {
	userType := ctx.GetString("role")
	err = nil
	if userType == "USER" && GetUidString(ctx) != userId {
		err = errors.New("Unauthorized to access this resource")
		return err
	}
	err = CheckUserType(ctx, userType)
	return err
}
