package friendshiphandler

import (
	authmiddleware "connect/internal/auth/middleware"
	friendshipservice "connect/internal/friendship/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	friendshipservice.CreateTableIfNotExists()

	friendshipGroup := router.Group("/friendships")

	friendshipGroup.Use(authmiddleware.AuthRequired())

	friendshipGroup.GET("/receiver/:id", GetFriendRequestHandler)
	friendshipGroup.GET("/receiver/:id/all", GetAllFriendRequestsHandler)
	friendshipGroup.GET("/requests", GetRequestsHandler)
	friendshipGroup.GET("/friends", GetFriendsHandler)

	friendshipGroup.POST("/:receiverId", SendFriendRequestHandler)

	friendshipGroup.PATCH("/accept/:requesterId", AcceptFriendRequestHandler)

	friendshipGroup.DELETE("/:receiverId", DeleteFriendRequestHandler)
}
