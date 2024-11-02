package repository

import (
	"fmt"

	"github.com/Ion-Stefan/saas-go-fiber/database"
	"github.com/Ion-Stefan/saas-go-fiber/internal/model"
	"gorm.io/gorm"
)

// CreateUser adds a new user to the database.
func CreateUser(user *model.User) error {
	result := database.DB.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetUserByID retrieves a user by their ID. Returns an error if not found or if the ID is invalid.
func GetUserByID(id uint) (*model.User, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid user ID: %d", id)
	}
	var user model.User
	result := database.DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with ID %d not found", id)
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by their email. Returns an error if not found.
func GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with email '%s' not found", email)
		}
		return nil, result.Error
	}
	return &user, nil
}

// UpdateUser updates user details by their ID. Returns an error if the update fails.
func UpdateUser(id uint, user *model.User) error {
	result := database.DB.Model(&model.User{}).Where("id = ?", id).Updates(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteUser removes a user by their ID. Returns an error if deletion fails.
func DeleteUser(id uint) error {
	result := database.DB.Unscoped().Where("id = ?", id).Delete(&model.User{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
