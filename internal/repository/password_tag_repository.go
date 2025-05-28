package repository

import (
	"errors"
	"gorm.io/gorm"
	"password-management-service/internal/models/password"
)

type PasswordTagRepository interface {
	AddPasswordTag(tag *password.PasswordTag) error
	UpdatePasswordTag(tag *password.PasswordTag) error
	GetPasswordTagByIDAndUserID(id uint, userID uint) (*password.PasswordTag, error)
	GetListPasswordTag(userID uint, index int, size int) (*[]password.PasswordTag, error)
	GetPasswordTagsByEntryID(entryID uint) ([]*password.PasswordTag, error)
	DeletePasswordTag(tag *password.PasswordTag) error
	FindOrCreate(name string, createdBy string) (*password.PasswordTag, error)
	LinkTagToEntry(entryID uint, tagID uint) error
	GetCountPasswordTag(userID uint) (int64, error)
}

type passwordTagRepository struct {
	db gorm.DB
}

func NewPasswordTagRepository(db gorm.DB) PasswordTagRepository {
	return &passwordTagRepository{
		db: db,
	}
}

func (r *passwordTagRepository) AddPasswordTag(tag *password.PasswordTag) error {
	if tag == nil {
		return errors.New("tag cannot be nil")
	}
	if tag.Name == "" {
		return errors.New("tag name cannot be empty")
	}
	if tag.CreatedBy == nil || *tag.CreatedBy == "" {
		return errors.New("created by cannot be empty")
	}
	if tag.UpdatedBy == nil || *tag.UpdatedBy == "" {
		return errors.New("updated by cannot be empty")
	}

	return r.db.Create(tag).Error
}

func (r *passwordTagRepository) UpdatePasswordTag(tag *password.PasswordTag) error {
	if tag == nil {
		return errors.New("tag cannot be nil")
	}
	if tag.Name == "" {
		return errors.New("tag name cannot be empty")
	}
	if tag.CreatedBy == nil || *tag.CreatedBy == "" {
		return errors.New("created by cannot be empty")
	}
	if tag.UpdatedBy == nil || *tag.UpdatedBy == "" {
		return errors.New("updated by cannot be empty")
	}

	return r.db.Save(tag).Error
}

func (r *passwordTagRepository) GetPasswordTagByIDAndUserID(tagID uint, userID uint) (*password.PasswordTag, error) {
	var tag password.PasswordTag
	err := r.db.Where("tag_id = ? AND user_id = ?", tagID, userID).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *passwordTagRepository) GetListPasswordTag(userID uint, index int, size int) (*[]password.PasswordTag, error) {
	var tags []password.PasswordTag
	err := r.db.
		Where("deleted_at IS NULL AND user_id = ?", userID).
		Order("tag_id ASC").
		Limit(size).
		Offset((index - 1) * size).
		Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return &tags, nil
}

func (r *passwordTagRepository) GetPasswordTagsByEntryID(entryID uint) ([]*password.PasswordTag, error) {
	var tags []*password.PasswordTag
	err := r.db.Table("password_entry_tags").
		Select("password_tags.*").
		Joins("JOIN password_tags ON password_entry_tags.tag_id = password_tags.id").
		Where("password_entry_tags.entry_id = ?", entryID).
		Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *passwordTagRepository) DeletePasswordTag(tag *password.PasswordTag) error {
	if tag == nil {
		return errors.New("tag cannot be nil")
	}
	if tag.DeletedBy == nil || *tag.DeletedBy == "" {
		return errors.New("deleted by cannot be empty")
	}

	var count int64
	err := r.db.Table("password_entry_tags").Where("tag_id = ?", tag.TagID).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("cannot delete tag as it is linked to password entries")
	}

	return r.db.Model(tag).
		Update("deleted_by", tag.DeletedBy).
		Delete(tag).Error
}

func (r *passwordTagRepository) FindOrCreate(name string, createdBy string) (*password.PasswordTag, error) {
	var tag password.PasswordTag
	err := r.db.Where("name = ?", name).First(&tag).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		tag = password.PasswordTag{Name: name, CreatedBy: &createdBy, UpdatedBy: &createdBy}
		if err := r.db.Create(&tag).Error; err != nil {
			return nil, err
		}
		return &tag, nil
	}
	return &tag, err
}

func (r *passwordTagRepository) LinkTagToEntry(entryID uint, tagID uint) error {
	return r.db.Exec(`INSERT INTO password_entry_tags (entry_id, tag_id) VALUES (?, ?) ON CONFLICT DO NOTHING`, entryID, tagID).Error
}

func (r *passwordTagRepository) GetCountPasswordTag(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&password.PasswordTag{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
