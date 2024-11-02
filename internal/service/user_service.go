package service

import (
	"github.com/Ion-Stefan/saas-go-fiber/internal/model"
	"github.com/Ion-Stefan/saas-go-fiber/internal/repository"
)

func CreateUser(user *model.User) error {
	return repository.CreateUser(user)
}

func UpdateUser(id uint, user *model.User) error {
	return repository.UpdateUser(id, user)
}

func GetUserByID(id uint) (*model.User, error) {
	return repository.GetUserByID(id)
}

func GetUserByEmail(email string) (*model.User, error) {
	return repository.GetUserByEmail(email)
}

func DeleteUser(id uint) error {
	return repository.DeleteUser(id)
}
