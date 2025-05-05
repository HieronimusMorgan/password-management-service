package repository

import (
	"gorm.io/gorm"
	"password-management-service/internal/models/user"
)

type UserRepository interface {
	GetUserByID(id uint) (*user.Users, error)
	GetUserByPhoneNumber(number string) (*user.Users, error)
	GetUserByClientID(clientID string) (*user.Users, error)
	GetUserByClientAndRole(clientID, roleID uint) (*[]user.Users, error)
}

type userRepository struct {
	db gorm.DB
}

func NewUserRepository(db gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r userRepository) GetUserByID(id uint) (*user.Users, error) {
	var users user.Users
	if err := r.db.Where("user_id = ?", id).First(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

func (r userRepository) GetUserByPhoneNumber(number string) (*user.Users, error) {
	var users user.Users
	if err := r.db.Where("phone_number = ?", number).Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

func (r userRepository) GetUserByClientID(clientID string) (*user.Users, error) {
	var users user.Users
	if err := r.db.Where("client_id = ?", clientID).Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

func (r userRepository) GetUserByClientAndRole(clientID, roleID uint) (*[]user.Users, error) {
	var users []user.Users
	if err := r.db.Where("client_id = ? AND role_id = ?", clientID, roleID).Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}
