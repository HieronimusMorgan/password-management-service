package password

import (
	"gorm.io/gorm"
	"time"
)

type PasswordGroup struct {
	GroupID   uint           `gorm:"primaryKey;column:group_id" json:"group_id,omitempty"`
	UserID    uint           `gorm:"column:user_id" json:"user_id,omitempty"`
	Name      string         `gorm:"column:name" json:"name,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at,omitempty"`
	CreatedBy *string        `gorm:"column:created_by" json:"created_by,omitempty"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at,omitempty"`
	UpdatedBy *string        `gorm:"column:updated_by" json:"updated_by,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	DeletedBy *string        `gorm:"column:deleted_by" json:"deleted_by,omitempty"`
}
