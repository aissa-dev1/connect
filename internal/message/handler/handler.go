package messagehandler

import (
	authmiddleware "connect/internal/auth/middleware"
	messageservice "connect/internal/message/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	messageservice.CreateTableIfNotExists()

	messageGroup := router.Group("/messages")

	messageGroup.Use(authmiddleware.AuthRequired())

	messageGroup.GET("/between/:senderId/:receiverId", GetMessagesBetweenUsersHandler)
}
