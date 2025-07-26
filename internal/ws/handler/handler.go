package wshandler

import (
	authmiddleware "connect/internal/auth/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	wsGroup := router.Group("/ws")
	wsAuthGroup := wsGroup

	wsAuthGroup.Use(authmiddleware.AuthRequired())

	wsAuthGroup.GET("/chat", HandleChat)
}
