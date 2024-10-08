package main

import (
	"auth_service/internal/db"
	"auth_service/internal/models"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectTestDB() (db.Storage, error) {
	dsn := "user=postgres password=1 dbname=auth port=5432 sslmode=disable"
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return db.PostgresStorage{db: DB}, err
	}

}

var testdb db.Storage = db.ConnectTestDB()
var testUser models.User = models.User{
	GUID:             "00",
	Email:            "e@mail.com",
	Password:         "zzz",
	RefreshTokenHash: "123abc456",
}

func TestMain(m *testing.M) {

	router := getRouter()
	router.Run("localhost:8080")

	testdb.Create(testUser)
}
func TestGetUsers(t *testing.T) {

}
