package repository

import (
	"gorm.io/gorm"
	"password-management-service/internal/dto/out"
	"password-management-service/internal/models/password"
	"password-management-service/internal/utils"
)

type PasswordEntryRepository interface {
	AddPasswordEntry(passwordEntry *password.PasswordEntry, passwordEntryKey *password.PasswordEntryKey) error
	UpdatePasswordEntry(passwordEntry *password.PasswordEntry) error
	UpdatePasswordEntryAndEntryKey(passwordEntry password.PasswordEntry, passwordEntryKey password.PasswordEntryKey) error
	DeletePasswordEntry(entryID uint) error
	GetListPasswordEntryResponse(userID uint) ([]out.PasswordEntryListResponse, error)
	GetPasswordEntryByEntryIDAndUserID(entryID, userID uint) (*password.PasswordEntry, error)
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

func (r *passwordEntryRepository) UpdatePasswordEntry(passwordEntry *password.PasswordEntry) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(utils.TablePasswordEntryName).Where("entry_id = ?", passwordEntry.EntryID).Updates(&passwordEntry).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *passwordEntryRepository) UpdatePasswordEntryAndEntryKey(passwordEntry password.PasswordEntry, passwordEntryKey password.PasswordEntryKey) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(utils.TablePasswordEntryName).Where("entry_id = ?", passwordEntry.EntryID).Updates(&passwordEntry).Error; err != nil {
			return err
		}
		if err := tx.Table(utils.TablePasswordEntryKeyName).Where("entry_id = ?", passwordEntry.EntryID).Updates(&passwordEntryKey).Error; err != nil {
			return err
		}
		return nil

	})
}

func (r *passwordEntryRepository) DeletePasswordEntry(entryID uint) error {
	if err := r.db.Unscoped().Table(utils.TablePasswordEntryName).Delete(&password.PasswordEntry{}, entryID).Error; err != nil {
		return err
	}
	return nil
}

func (r *passwordEntryRepository) GetListPasswordEntryResponse(userID uint) ([]out.PasswordEntryListResponse, error) {
	var passwordEntry []out.PasswordEntryListResponse

	err := r.db.Raw(`
		SELECT 
			pe.entry_id, 
			pe.title, 
			pe.username, 
			pe.url, 
			pe.tags, 
			pg.name AS group_name
		FROM password_entries pe
		LEFT JOIN password_groups pg ON pg.group_id = pe.group_id
		WHERE pe.user_id = ? AND pe.deleted_at IS NULL
		ORDER BY pe.entry_id ASC
	`, userID).Scan(&passwordEntry).Error

	if err != nil {
		return nil, err
	}
	return passwordEntry, nil
}

func (r *passwordEntryRepository) GetPasswordEntryByEntryIDAndUserID(entryID, userID uint) (*password.PasswordEntry, error) {
	var passwordEntry password.PasswordEntry
	if err := r.db.Where("entry_id = ? AND user_id = ?", entryID, userID).First(&passwordEntry).Error; err != nil {
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
