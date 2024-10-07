package main

import (
	db "auth_service/internal/database"
	"auth_service/internal/handler"

	"github.com/gin-gonic/gin"
)

func getRouter() *gin.Engine {
	router := gin.Default()
	//router.Use(handler.RefreshTokenMiddleware())
	router.POST("/auth/login/:ip/:id", handler.Login)
	router.GET("/auth/refresh/:ip", handler.RefreshTokens)

	router.GET("/users", handler.GetUsers)
	router.GET("/user/cookie", handler.Cookies)
	return router
}

func main() {
	router := getRouter()

	db.ConnectDB()
	router.Run("localhost:8080")
}
