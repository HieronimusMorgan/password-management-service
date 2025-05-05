package password

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

type PasswordEntry struct {
	EntryID           uint            `gorm:"primaryKey;column:entry_id;not null" json:"entry_id,omitempty"`
	UserID            uint            `gorm:"column:user_id;not null" json:"user_id,omitempty"`
	GroupID           *uint           `gorm:"column:group_id" json:"group_id,omitempty"`
	Title             string          `gorm:"column:title;not null" json:"title,omitempty"`
	Username          string          `gorm:"column:username;not null" json:"username,omitempty"`
	EncryptedPassword string          `gorm:"column:encrypted_password;not null" json:"encrypted_password,omitempty"`
	EncryptedNotes    *string         `gorm:"column:encrypted_notes" json:"encrypted_notes,omitempty"`
	URL               *string         `gorm:"column:url" json:"url,omitempty"`
	Tags              *pq.StringArray `gorm:"type:text[];column:tags" json:"tags,omitempty"`
	ExpiresAt         *time.Time      `gorm:"column:expires_at" json:"expires_at,omitempty"`
	LastAccessedAt    *time.Time      `gorm:"column:last_accessed_at" json:"last_accessed_at,omitempty"`
	CreatedAt         time.Time       `gorm:"column:created_at" json:"created_at,omitempty"`
	CreatedBy         *string         `gorm:"column:created_by" json:"created_by,omitempty"`
	UpdatedAt         time.Time       `gorm:"column:updated_at" json:"updated_at,omitempty"`
	UpdatedBy         *string         `gorm:"column:updated_by" json:"updated_by,omitempty"`
	DeletedAt         gorm.DeletedAt  `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	DeletedBy         *string         `gorm:"column:deleted_by" json:"deleted_by,omitempty"`
}
