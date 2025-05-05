package repository

import (
	"gorm.io/gorm"
	"password-management-service/internal/models/password"
	"password-management-service/internal/utils"
)

type PasswordEntryRepository interface {
	AddPasswordEntry(passwordEntry *password.PasswordEntry, passwordEntryKey *password.PasswordEntryKey) error
	UpdatePasswordEntry(passwordEntry password.PasswordEntry) error
	DeletePasswordEntry(entryID uint) error
	GetPasswordEntryByEntryID(entryID uint) (*password.PasswordEntry, error)
	GetPasswordEntryByUserID(userID string) ([]password.PasswordEntry, error)
	GetPasswordEntryByGroupID(groupID uint) ([]password.PasswordEntry, error)
	GetPasswordEntryByGroupIDAndUserID(groupID uint, userID string) ([]password.PasswordEntry, error)
	GetPasswordEntryByGroupIDAndEntryID(groupID uint, entryID uint) (*password.PasswordEntry, error)
	GetPasswordEntryByGroupIDAndUserIDAndEntryID(groupID uint, userID string, entryID uint) (*password.PasswordEntry, error)
}

type passwordEntryRepository struct {
	db gorm.DB
}

func NewPasswordEntryRepository(db gorm.DB) PasswordEntryRepository {
	return &passwordEntryRepository{
		db: db,
	}
}

func (r *passwordEntryRepository) AddPasswordEntry(passwordEntry *password.PasswordEntry, passwordEntryKey *password.PasswordEntryKey) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(utils.TablePasswordEntryName).Create(&passwordEntry).Error; err != nil {
			return err
		}
		passwordEntryKey.EntryID = passwordEntry.EntryID
		if err := tx.Table(utils.TablePasswordEntryKeyName).Create(&passwordEntryKey).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *passwordEntryRepository) UpdatePasswordEntry(passwordEntry password.PasswordEntry) error {
	if err := r.db.Save(&passwordEntry).Error; err != nil {
		return err
	}
	return nil
}

func (r *passwordEntryRepository) DeletePasswordEntry(entryID uint) error {
	if err := r.db.Delete(&password.PasswordEntry{}, entryID).Error; err != nil {
		return err
	}
	return nil
}

func (r *passwordEntryRepository) GetPasswordEntryByEntryID(entryID uint) (*password.PasswordEntry, error) {
	var passwordEntry password.PasswordEntry
	if err := r.db.Where("entry_id = ?", entryID).First(&passwordEntry).Error; err != nil {
		return nil, err
	}
	return &passwordEntry, nil
}

func (r *passwordEntryRepository) GetPasswordEntryByUserID(userID string) ([]password.PasswordEntry, error) {
	var passwordEntry []password.PasswordEntry
	if err := r.db.Where("user_id = ?", userID).Find(&passwordEntry).Error; err != nil {
		return nil, err
	}
	return passwordEntry, nil
}

func (r *passwordEntryRepository) GetPasswordEntryByGroupID(groupID uint) ([]password.PasswordEntry, error) {
	var passwordEntry []password.PasswordEntry
	if err := r.db.Where("group_id = ?", groupID).Find(&passwordEntry).Error; err != nil {
		return nil, err
	}
	return passwordEntry, nil
}

func (r *passwordEntryRepository) GetPasswordEntryByGroupIDAndUserID(groupID uint, userID string) ([]password.PasswordEntry, error) {
	var passwordEntry []password.PasswordEntry
	if err := r.db.Where("group_id = ? AND user_id = ?", groupID, userID).Find(&passwordEntry).Error; err != nil {
		return nil, err
	}
	return passwordEntry, nil
}

func (r *passwordEntryRepository) GetPasswordEntryByGroupIDAndEntryID(groupID uint, entryID uint) (*password.PasswordEntry, error) {
	var passwordEntry password.PasswordEntry
	if err := r.db.Where("group_id = ? AND id = ?", groupID, entryID).First(&passwordEntry).Error; err != nil {
		return nil, err
	}
	return &passwordEntry, nil
}

func (r *passwordEntryRepository) GetPasswordEntryByGroupIDAndUserIDAndEntryID(groupID uint, userID string, entryID uint) (*password.PasswordEntry, error) {
	var passwordEntry password.PasswordEntry
	if err := r.db.Where("group_id = ? AND user_id = ? AND id = ?", groupID, userID, entryID).First(&passwordEntry).Error; err != nil {
		return nil, err
	}
	return &passwordEntry, nil
}
