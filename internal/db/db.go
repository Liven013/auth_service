package db

import (
	"log"

	"auth_service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage interface {
	GetOne(string) models.User
	GetAll() ([]models.User, error)
	Update(models.User) error
	Create(models.User) error
}

var DB = ConnectDB()

type PostgresStorage struct {
	DataBase *gorm.DB
}

func ConnectDB() PostgresStorage {
	dsn := "user=postgres password=1 dbname=auth port=5432 sslmode=disable"
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	// Автоматическая миграция
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Ошибка миграции базы данных: %v", err)
	}
	return PostgresStorage{DataBase: DB}
}
