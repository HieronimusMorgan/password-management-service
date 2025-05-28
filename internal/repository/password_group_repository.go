package repository

import (
	"gorm.io/gorm"
	"password-management-service/internal/dto/out"
	"password-management-service/internal/models/password"
	"password-management-service/internal/utils"
	"time"
)

type PasswordGroupRepository interface {
	AddPasswordGroup(group *password.PasswordGroup) error
	UpdatePasswordGroup(group *password.PasswordGroup) error
	GetPasswordGroupByID(groupID uint) (*password.PasswordGroup, error)
	GetPasswordGroupByUserID(userID uint) ([]password.PasswordGroup, error)
	GetPasswordGroupByUserIDAndGroupID(userID, groupID uint) (*password.PasswordGroup, error)
	GetItemListPasswordGroup(groupID uint, userID uint) (interface{}, error)
	DeletePasswordGroupByID(groupID uint, clientID string) error
}

type passwordGroupRepository struct {
	db gorm.DB
}

func NewPasswordGroupRepository(db gorm.DB) PasswordGroupRepository {
	return &passwordGroupRepository{
		db: db,
	}
}

func (r *passwordGroupRepository) AddPasswordGroup(group *password.PasswordGroup) error {
	if err := r.db.Table(utils.TablePasswordGroupName).Create(&group).Error; err != nil {
		return err
	}
	return nil
}

func (r *passwordGroupRepository) UpdatePasswordGroup(group *password.PasswordGroup) error {
	if err := r.db.Table(utils.TablePasswordGroupName).Where("group_id = ?", group.GroupID).Updates(group).Error; err != nil {
		return err
	}
	return nil
}

func (r *passwordGroupRepository) GetPasswordGroupByID(groupID uint) (*password.PasswordGroup, error) {
	var passwordGroup password.PasswordGroup
	if err := r.db.Table(utils.TablePasswordGroupName).Where("group_id = ?", groupID).First(&passwordGroup).Error; err != nil {
		return &passwordGroup, err
	}
	return &passwordGroup, nil
}

func (r *passwordGroupRepository) GetPasswordGroupByUserID(userID uint) ([]password.PasswordGroup, error) {
	var passwordGroups []password.PasswordGroup
	if err := r.db.Table(utils.TablePasswordGroupName).Where("user_id = ?", userID).Find(&passwordGroups).Error; err != nil {
		return passwordGroups, err
	}
	return passwordGroups, nil
}

func (r *passwordGroupRepository) GetPasswordGroupByUserIDAndGroupID(userID, groupID uint) (*password.PasswordGroup, error) {
	var passwordGroup password.PasswordGroup
	if err := r.db.Table(utils.TablePasswordGroupName).Where("user_id = ? AND group_id = ?", userID, groupID).First(&passwordGroup).Error; err != nil {
		return &passwordGroup, err
	}
	return &passwordGroup, nil
}

func (r *passwordGroupRepository) GetItemListPasswordGroup(groupID uint, userID uint) (interface{}, error) {
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
		WHERE pe.user_id = ? AND pe.group_id = ? AND pe.deleted_at IS NULL
		ORDER BY pe.entry_id ASC
	`, userID, groupID).Scan(&passwordEntry).Error

	if err != nil {
		return nil, err
	}

	return passwordEntry, nil
}

func (r *passwordGroupRepository) DeletePasswordGroupByID(groupID uint, clientID string) error {
	if err := r.db.Table(utils.TablePasswordGroupName).Where("group_id = ?", groupID).
		Updates(map[string]interface{}{
			"deleted_by": clientID,
			"deleted_at": time.Now(),
		}).Error; err != nil {
		return err
	}
	return nil
}
