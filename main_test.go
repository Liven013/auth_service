package main

import (
	"auth_service/internal/db"
	"auth_service/internal/models"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectTestDB() db.Storage {
	dsn := "user=postgres password=1 dbname=auth port=5432 sslmode=disable"
	DB, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db.PostgresStorage{DataBase: DB}
}

var testdb db.Storage = ConnectTestDB()

var testUser models.User = models.User{
	GUID:             "00",
	Email:            "e@mail.com",
	Password:         "zzz",
	RefreshTokenHash: "123abc456",
}

func TestMain(m *testing.M) {
	router := getRouter()
	router.Run("localhost:8080")
}
func TestGetUsers(t *testing.T) {
	testdb.Create(testUser)
}
