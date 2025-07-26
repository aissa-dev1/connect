package authhandler

import (
	authmiddleware "connect/internal/auth/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	authGroup := router.Group("/auth")
	authRequiredGroup := authGroup.Group("/")

	authRequiredGroup.Use(authmiddleware.AuthRequired())

	authGroup.POST("/sign-up", SignUpHandler)
	authGroup.POST("/sign-in", SignInHandler)
	authRequiredGroup.POST("/sign-out", SignOutHandler)
}
