package main

import (
	authhandler "connect/internal/auth/handler"
	blockhandler "connect/internal/block/handler"
	"connect/internal/db"
	friendshiphandler "connect/internal/friendship/handler"
	messagehandler "connect/internal/message/handler"
	"connect/internal/middleware"
	"connect/internal/pkg/hasher"
	profilehandler "connect/internal/profile/handler"
	userhandler "connect/internal/user/handler"
	wshandler "connect/internal/ws/handler"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	envErr := godotenv.Load()

	if envErr != nil {
		log.Fatalf("Unable to load .env file %v\n", envErr)
	}

	hasher.SetGlobalHash(hasher.NewHasher(hasher.NewBcrypt()))

	router := gin.Default()
	router.Use(middleware.AllowCors())

	pool := db.Connect()
	defer pool.Close()

	authhandler.RegisterRoutes(router)
	userhandler.RegisterRoutes(router)
	profilehandler.RegisterRoutes(router)
	friendshiphandler.RegisterRoutes(router)
	blockhandler.RegisterRoutes(router)
	messagehandler.RegisterRoutes(router)
	wshandler.RegisterRoutes(router)

	router.Run(":8080")
}
