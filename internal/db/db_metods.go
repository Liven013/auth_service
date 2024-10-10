package db

import (
	"auth_service/internal/models"
)

func (ps PostgresStorage) GetOne(guid string) models.User {
	var user models.User
	ps.DataBase.First(&user, "guid = ?", guid)
	return user
}

func (ps PostgresStorage) GetAll() ([]models.User, error) {
	var users []models.User
	result := ps.DataBase.Find(&users)
	return users, result.Error
}

func (ps PostgresStorage) Update(userUpdate models.User) error {
	return ps.DataBase.Save(&userUpdate).Error
}

func (ps PostgresStorage) Create(user models.User) error {
	return ps.DataBase.Create(&user).Error
}
