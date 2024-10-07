package db

import (
	"log"

	"auth_service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	dsn := "user=postgres password=1 dbname=auth port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	// Автоматическая миграция
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Ошибка миграции базы данных: %v", err)
	}
}
