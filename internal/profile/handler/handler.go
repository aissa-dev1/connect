package profilehandler

import (
	authmiddleware "connect/internal/auth/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	profileGroup := router.Group("/profile")

	profileGroup.Use(authmiddleware.AuthRequired())

	profileGroup.GET("", GetMyProfileHandler)
}
