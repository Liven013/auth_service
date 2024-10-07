package models

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	GUID             string `json:"guid" gorm:"primaryKey"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	RefreshTokenHash string `json:"rth"`
}

type LoginInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AccessClaims struct {
	jwt.RegisteredClaims
	GUID string `json:"sub"`
	EXP  int64  `json:"exp"`
	Time int64  `json:"time"`
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	GUID   string `json:"sub"`
	UserIP string `json:"ip"`
	EXP    int64  `json:"exp"`
	Time   int64  `json:"time"`
}
