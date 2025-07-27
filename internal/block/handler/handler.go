package blockhandler

import (
	authmiddleware "connect/internal/auth/middleware"
	blockservice "connect/internal/block/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	blockservice.CreateTableIfNotExists()

	blockGroup := router.Group("/blocks")

	blockGroup.Use(authmiddleware.AuthRequired())

	blockGroup.GET("/between/:blockerId/:blockedId", GetBlocksBetweenUsersHandler)
	blockGroup.GET("/all", GetBlockedUsersHandler)

	blockGroup.POST("/:blockedId", BlockHandler)

	blockGroup.DELETE("/:blockedId", UnblockHandler)
}
