package repository

import (
	"gorm.io/gorm"
	"password-management-service/internal/models/password"
	"password-management-service/internal/utils"
)

type PasswordEntryKeysRepository interface {
	GetPasswordEntryKeyByEntryID(entryID uint) (*password.PasswordEntryKey, error)
}

type passwordEntryKeysRepository struct {
	db gorm.DB
}

func NewPasswordEntryKeysRepository(db gorm.DB) PasswordEntryKeysRepository {
	return &passwordEntryKeysRepository{
		db: db,
	}
}

func (r *passwordEntryKeysRepository) GetPasswordEntryKeyByEntryID(entryID uint) (*password.PasswordEntryKey, error) {
	var passwordEntryKey password.PasswordEntryKey
	if err := r.db.Table(utils.TablePasswordEntryKeyName).Where("entry_id = ?", entryID).First(&passwordEntryKey).Error; err != nil {
		return &passwordEntryKey, err
	}
	return &passwordEntryKey, nil
}
