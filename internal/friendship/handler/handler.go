package friendshiphandler

import (
	authmiddleware "connect/internal/auth/middleware"
	friendshipservice "connect/internal/friendship/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	friendshipservice.CreateTableIfNotExists()

	friendshipGroup := router.Group("/friendships")
	friendshipAuthGroup := friendshipGroup.Group("/")

	friendshipAuthGroup.Use(authmiddleware.AuthRequired())

	friendshipAuthGroup.GET("/receiver/:id", GetFriendRequestHandler)
	friendshipAuthGroup.GET("/receiver/:id/all", GetAllFriendRequestsHandler)
	friendshipAuthGroup.GET("/requests", GetRequestsHandler)
	friendshipAuthGroup.GET("/friends", GetFriendsHandler)

	friendshipAuthGroup.POST("/:receiverId", SendFriendRequestHandler)

	friendshipAuthGroup.PATCH("/accept/:requesterId", AcceptFriendRequestHandler)

	friendshipAuthGroup.DELETE("/:receiverId", DeleteFriendRequestHandler)
}
