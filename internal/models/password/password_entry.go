package password

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

type PasswordEntry struct {
	EntryID           uint            `gorm:"primaryKey;column:entry_id"`
	UserID            uint            `gorm:"column:user_id"`
	GroupID           *uint           `gorm:"column:group_id"`
	Title             string          `gorm:"column:title"`
	Username          string          `gorm:"column:username"`
	EncryptedPassword string          `gorm:"column:encrypted_password"`
	EncryptedNotes    *string         `gorm:"column:encrypted_notes"`
	URL               *string         `gorm:"column:url"`
	Tags              *pq.StringArray `gorm:"type:text[];column:tags"`
	ExpiresAt         *time.Time      `gorm:"column:expires_at"`
	LastAccessedAt    *time.Time      `gorm:"column:last_accessed_at"`
	CreatedAt         time.Time       `gorm:"column:created_at"`
	CreatedBy         *string         `gorm:"column:created_by"`
	UpdatedAt         time.Time       `gorm:"column:updated_at"`
	UpdatedBy         *string         `gorm:"column:updated_by"`
	DeletedAt         gorm.DeletedAt  `gorm:"column:deleted_at"`
	DeletedBy         *string         `gorm:"column:deleted_by"`
}
