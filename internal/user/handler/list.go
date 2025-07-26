package userhandler

import (
	"connect/internal/pkg/response"
	userconstants "connect/internal/user/constants"
	usermodel "connect/internal/user/model"
	userservice "connect/internal/user/service"

	"github.com/gin-gonic/gin"
)

func GetUserByUserNameHandler(ctx *gin.Context) {
	userNameParam := ctx.Param("username")

	user, userErr := userservice.GetMinimalUserByUsername(userNameParam)

	if userErr != nil {
		response.RespondInternalError(ctx, userErr.Error())
		return
	}
	if user == nil {
		response.RespondNotFoundError(ctx, userconstants.ErrUserNotFound.Error())
		return
	}

	response.RespondData[usermodel.MinimalUser](ctx, *user)
}
