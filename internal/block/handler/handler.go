package blockhandler

import (
	authmiddleware "connect/internal/auth/middleware"
	blockservice "connect/internal/block/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	blockservice.CreateTableIfNotExists()

	blockGroup := router.Group("/blocks")
	blockAuthGroup := blockGroup.Group("/")

	blockAuthGroup.Use(authmiddleware.AuthRequired())

	blockAuthGroup.GET("/between/:blockerId/:blockedId", GetBlocksBetweenUsersHandler)
	blockAuthGroup.GET("/all", GetBlockedUsersHandler)

	blockAuthGroup.POST("/:blockedId", BlockHandler)

	blockAuthGroup.DELETE("/:blockedId", UnblockHandler)
}
