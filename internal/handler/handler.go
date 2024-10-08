package handler

import (
	"auth_service/internal/db"
	"auth_service/internal/models"
	"auth_service/internal/services"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	var users []models.User
	users, err := db.DB.GetAll()
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("not found DB %s", err)})
	}
	c.IndentedJSON(http.StatusOK, users)
}

func Login(c *gin.Context) {
	ip := c.Param("ip")
	id := c.Param("id")

	var user models.LoginInfo
	if err := c.ShouldBindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//поиск пользователя с данным guid
	existingUser := db.DB.GetOne(id)

	if existingUser.GUID == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("user id:%s does not exist", id)})
		return
	}

	//проверка email

	if existingUser.Email != user.Email {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
		return
	}

	// проверка password, но лучше хранить в БД хеш паролей
	if existingUser.Password != user.Password {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	//генерация пары токенов
	accessToken, refreshToken, err := services.GeneratePairToken(existingUser, ip)
	if err != nil {
		c.IndentedJSON(http.StatusMethodNotAllowed, gin.H{"error": "tokens not create"})
		return
	}

	//сохраняем refresh токен в куки
	c.SetCookie("RefreshToken", refreshToken, int(time.Now().Add(services.RefreshTokenExpiry).Unix()), "/", "localhost", false, true)

	//выдача access токена пользователю
	c.IndentedJSON(http.StatusOK, gin.H{"success": fmt.Sprintf("user %s logged in.       access token: %s", user.Email, accessToken)})
}

func RefreshTokens(c *gin.Context) {
	ip := c.Param("ip")

	accessToken := c.GetHeader("Authorization")
	if !strings.HasPrefix(accessToken, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Access token not provided"})
		return
	}
	accessToken = strings.TrimPrefix(accessToken, "Bearer ")

	refreshToken, err := c.Cookie("RefreshToken")
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "not found cookies"})
		return
	}

	newAccessToken, newRefreshToken, err := services.RefreshTokens(accessToken, refreshToken, ip)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "no refresh"})
		return
	}
	c.SetCookie("RefreshToken", newRefreshToken, int(time.Now().Add(services.RefreshTokenExpiry).Unix()), "/", "localhost", false, true)

	c.IndentedJSON(http.StatusOK, gin.H{"success": "refresh tokens complete: \n Refrsh: " + newRefreshToken + "\n Access" + newAccessToken})
}
