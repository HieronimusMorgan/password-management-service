package repository

import (
	"errors"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"password-management-service/internal/dto/out"
	"password-management-service/internal/models/password"
	"password-management-service/internal/utils"
	"strings"
)

type PasswordEntryRepository interface {
	AddPasswordEntry(passwordEntry *password.PasswordEntry, passwordEntryKey *password.PasswordEntryKey, tags pq.StringArray, userID uint) error
	UpdatePasswordEntry(passwordEntry *password.PasswordEntry) error
	UpdatePasswordEntryAndEntryKey(passwordEntry password.PasswordEntry, passwordEntryKey password.PasswordEntryKey) error
	DeletePasswordEntry(entryID uint) error
	GetListPasswordEntryResponse(userID uint, tags string, index int, size int) ([]out.PasswordEntryListResponse, error)
	GetListPasswordEntryResponseByTags(userID uint, tags []string, index int, size int) ([]out.PasswordEntryListResponse, error)
	GetPasswordEntryByEntryIDAndUserID(entryID, userID uint) (*password.PasswordEntry, error)
	GetPasswordEntryByUserID(userID string) ([]password.PasswordEntry, error)
	GetPasswordEntryByGroupID(groupID uint) ([]password.PasswordEntry, error)
	GetPasswordEntryByGroupIDAndUserID(groupID uint, userID string) ([]password.PasswordEntry, error)
	GetPasswordEntryByGroupIDAndEntryID(groupID uint, entryID uint) (*password.PasswordEntry, error)
	GetPasswordEntryByGroupIDAndUserIDAndEntryID(groupID uint, userID string, entryID uint) (*password.PasswordEntry, error)
	GetCountPasswordEntriesByUserID(userID uint) (int64, error)
	GetCountPasswordEntriesByTags(id uint, tags []string) (int64, error)
}

type passwordEntryRepository struct {
	db gorm.DB
}

func NewPasswordEntryRepository(db gorm.DB) PasswordEntryRepository {
	return &passwordEntryRepository{
		db: db,
	}
}

func (r *passwordEntryRepository) AddPasswordEntry(passwordEntry *password.PasswordEntry, passwordEntryKey *password.PasswordEntryKey, tags pq.StringArray, userID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(utils.TablePasswordEntryName).Create(passwordEntry).Error; err != nil {
			return err
		}

		passwordEntryKey.EntryID = passwordEntry.EntryID
		if err := tx.Table(utils.TablePasswordEntryKeyName).Create(passwordEntryKey).Error; err != nil {
			return err
		}

		for _, tagName := range tags {
			var tag password.PasswordTag
			if err := r.db.Where("name = ? AND user_id = ?", tagName, userID).First(&tag).Error; errors.Is(err, gorm.ErrRecordNotFound) {
				tag = password.PasswordTag{UserID: userID, Name: tagName, CreatedBy: passwordEntry.CreatedBy, UpdatedBy: passwordEntry.CreatedBy}
				if err := r.db.Create(&tag).Error; err != nil {
					return err
				}
			}

			if err := tx.Table(utils.TablePasswordEntryTagName).Create(&password.PasswordEntryTag{
				EntryID: passwordEntry.EntryID,
				TagID:   tag.TagID,
			}).Error; err != nil {
				return err
			}
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

func (r *passwordEntryRepository) GetListPasswordEntryResponse(userID uint, tagParams string, index int, size int) ([]out.PasswordEntryListResponse, error) {
	var passwordEntries []out.PasswordEntryListResponse

	// Build base query
	query := `
		SELECT DISTINCT
			pe.entry_id,
			pe.title,
			pe.url,
			pg.name AS group_name
		FROM password_entries pe
		LEFT JOIN password_groups pg ON pg.group_id = pe.group_id
		LEFT JOIN password_entry_tags pet ON pet.entry_id = pe.entry_id
		LEFT JOIN password_tags pt ON pt.tag_id = pet.tag_id
		WHERE pe.user_id = ? AND pe.deleted_at IS NULL
	`
	var args []interface{}
	args = append(args, userID)

	// Filter by tags if provided
	if tagParams != "" {
		tags := strings.Split(tagParams, ",")
		query += ` AND pt.name = ANY(?)`
		args = append(args, pq.Array(tags))
	}

	query += ` ORDER BY pe.entry_id ASC LIMIT ? OFFSET ?`
	args = append(args, size, (index-1)*size)

	// Execute main query
	if err := r.db.Raw(query, args...).Scan(&passwordEntries).Error; err != nil {
		return nil, err
	}

	// Fetch tags for all entries in a single query for efficiency
	entryIDs := make([]uint, 0, len(passwordEntries))
	entryIndex := make(map[uint]int)
	for i, entry := range passwordEntries {
		entryIDs = append(entryIDs, entry.EntryID)
		entryIndex[entry.EntryID] = i
	}

	if len(entryIDs) > 0 {
		type tagResult struct {
			EntryID uint
			Name    string
		}
		var tagResults []tagResult
		if err := r.db.Table(utils.TablePasswordEntryTagName).
			Select("password_entry_tags.entry_id, pt.name").
			Joins("JOIN password_tags pt ON pt.tag_id = password_entry_tags.tag_id").
			Where("password_entry_tags.entry_id IN ?", entryIDs).
			Scan(&tagResults).Error; err != nil {
			return nil, err
		}

		// Map tags to entries
		tagMap := make(map[uint][]string)
		for _, tr := range tagResults {
			tagMap[tr.EntryID] = append(tagMap[tr.EntryID], tr.Name)
		}
		for entryID, tags := range tagMap {
			strArray := pq.StringArray(tags)
			passwordEntries[entryIndex[entryID]].Tags = &strArray
		}
	}

	return passwordEntries, nil
}

func (r *passwordEntryRepository) GetListPasswordEntryResponseByTags(userID uint, tags []string, index int, size int) ([]out.PasswordEntryListResponse, error) {
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
		WHERE pe.user_id = ? AND pe.deleted_at IS NULL AND pe.tags && ?::text[]
		ORDER BY pe.entry_id ASC
		LIMIT ? OFFSET ?
	`, userID, pq.Array(tags), size, (index-1)*size).Scan(&passwordEntry).Error

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

func (r *passwordEntryRepository) GetCountPasswordEntriesByUserID(userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&password.PasswordEntry{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *passwordEntryRepository) GetCountPasswordEntriesByTags(id uint, tags []string) (int64, error) {
	var count int64
	if err := r.db.Model(&password.PasswordEntry{}).
		Where("user_id = ? AND tags && ?::text[]", id, pq.Array(tags)).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
