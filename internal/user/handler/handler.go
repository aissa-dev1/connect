package userhandler

import (
	userservice "connect/internal/user/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	userservice.CreateTableIfNotExists()

	userGroup := router.Group("/users")

	userGroup.GET("/u/:username", GetUserByUserNameHandler)

	userGroup.POST("/search", SearchUsersHandler)
}
