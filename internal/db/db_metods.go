package db

import (
	"auth_service/internal/models"
	"fmt"
)

func (ps PostgresStorage) GetOne(guid string) (models.User, error) {
	var user models.User
	result := ps.db.First(&user, "guid = ?", guid).Error
	return user, fmt.Errorf(result.Error())
}

func (ps PostgresStorage) GetAll() ([]models.User, error) {
	var users []models.User
	result := ps.db.Find(&users)
	return users, result.Error
}

func (ps PostgresStorage) Update(userUpdate models.User) error {
	return ps.db.Save(&userUpdate).Error
}
