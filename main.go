package main

import (
	"auth_service/internal/handler"

	"github.com/gin-gonic/gin"
)

func getRouter() *gin.Engine {
	router := gin.Default()
	//router.Use(handler.RefreshTokenMiddleware())
	router.POST("/auth/login/:ip/:id", handler.Login)
	router.GET("/auth/refresh/:ip", handler.RefreshTokens)

	router.GET("/users", handler.GetUsers)
	return router
}

func main() {
	router := getRouter()

	router.Run("localhost:8080")
}
