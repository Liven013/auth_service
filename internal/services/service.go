package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"time"

	db "auth_service/internal/db"
	"auth_service/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenExpiry  = 2 * time.Minute
	RefreshTokenExpiry = 24 * time.Hour
	secretKey          = "secret_key" // безопасное хранение ключа
)

func FindUserByGUID(userGUID string) (models.User, error) {
	return db.DB.GetOne(userGUID)
}

func GenerateAccessToken(userGUID string, timeCreated int64) (string, error) {
	claims := models.AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{},
		GUID:             userGUID,
		EXP:              time.Now().Add(RefreshTokenExpiry).Unix(),
		Time:             timeCreated,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(secretKey))
}

func GenerateRefreshToken(userGUID, userIP string, timeCreatedAccess int64) (string, error) {
	claims := models.RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{},
		GUID:             userGUID,
		UserIP:           userIP,
		EXP:              time.Now().Add(RefreshTokenExpiry).Unix(),
		Time:             timeCreatedAccess,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(secretKey))
}

func GeneratePairToken(user models.User, ip string) (string, string, error) {
	timeCreated := time.Now().Unix()

	accessToken, err := GenerateAccessToken(user.GUID, timeCreated)
	if err != nil {
		return "", "", fmt.Errorf("accessToken error")
	}

	refreshToken, err := GenerateRefreshToken(user.GUID, ip, timeCreated)
	if err != nil {
		return "", "", fmt.Errorf("refreshToken error")
	}

	err = StoreRefreshToken(user.GUID, refreshToken)
	if err != nil {
		fmt.Println("ошибка записи в БД")
		return "", "", fmt.Errorf("store error")
	}

	return accessToken, refreshToken, nil
}

func StoreRefreshToken(userGUID string, refreshToken string) error {
	hash := hashToken(refreshToken)
	user, err := db.DB.GetOne(userGUID)
	if err != nil {
		return err
	}

	user.RefreshTokenHash = string(hash)
	return db.DB.Update(user)
}

func hashToken(token string) string {
	hasher := sha256.New()
	hasher.Write([]byte(token))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetTokenAccessCreatedTime(token string) (int64, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS512 {
			return nil, fmt.Errorf("no actual alg: %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	}

	claims := &models.AccessClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, keyFunc)
	if err != nil {
		return 0, fmt.Errorf("parse error")
	}

	if !parsedToken.Valid {
		return 0, fmt.Errorf("unvalid token")
	}

	return claims.Time, nil
}

func GetTokenRefreshCreatedTime(token string) (int64, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS512 {
			return nil, fmt.Errorf("no actual alg: %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	}

	claims := &models.RefreshClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, keyFunc)
	if err != nil {
		return 0, fmt.Errorf("parse error")
	}

	if !parsedToken.Valid {
		return 0, fmt.Errorf("unvalid token")
	}

	return claims.Time, nil
}

func GetGUIDByAccessToken(token string) (string, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS512 {
			return nil, fmt.Errorf("no actual alg: %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	}

	claims := &models.AccessClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, keyFunc)
	if err != nil {
		fmt.Println("ошибка тут", err)
		return "", fmt.Errorf("parse error ")
	}

	if !parsedToken.Valid {
		return "", fmt.Errorf("unvalid token")
	}
	return claims.GUID, nil
}

func CheckLink(RefreshToken, AccessToken string) (bool, error) {
	timeCreatedInRefreshToken, err := GetTokenRefreshCreatedTime(RefreshToken)
	if err != nil {
		return false, err
	}
	timeCreatedAccesToken, err := GetTokenAccessCreatedTime(AccessToken)
	if err != nil {
		return false, err
	}

	if timeCreatedAccesToken == timeCreatedInRefreshToken {
		return true, nil
	}

	return false, nil
}

func GetIPByRefreshToken(token string) (string, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS512 {
			return nil, fmt.Errorf("no actual alg: %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	}

	claims := &models.RefreshClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, keyFunc)
	if err != nil {
		return "", fmt.Errorf("parse error")
	}

	if !parsedToken.Valid {
		return "", fmt.Errorf("unvalid token")
	}

	return claims.UserIP, nil
}

func ChekIP(token, actualIP string) (bool, error) {
	tokensIP, err := GetIPByRefreshToken(token)
	if err != nil {
		return false, fmt.Errorf("parse id error")
	}
	if tokensIP == actualIP {
		return true, nil
	}
	return false, nil
}

func RefreshTokens(providedRefreshToken, providedAcessToken, clientIP string) (string, string, error) {
	guid, err := GetGUIDByAccessToken(providedAcessToken)
	if err != nil {
		return "", "", err
	}

	user, err := db.DB.GetOne(guid)
	if err != nil {
		return "", "", err
	}

	providedRefreshTokenHash := hashToken(providedRefreshToken)
	if providedRefreshTokenHash == user.RefreshTokenHash {
		return "", "", fmt.Errorf("wrong Refresh token")
	}

	linkFlag, err := CheckLink(providedRefreshToken, providedAcessToken)
	if err != nil {
		return "", "", err
	}
	if !linkFlag {
		return "", "", fmt.Errorf("unlinked tokens")
	}

	valid, err := ChekIP(providedRefreshToken, clientIP)
	if err != nil {
		return "", "", err
	}

	if !valid {
		fmt.Println("email massage")
	}

	newAccessToken, newRefreshToken, err := GeneratePairToken(user, clientIP)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}
