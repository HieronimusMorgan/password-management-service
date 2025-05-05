package password

import (
	"gorm.io/gorm"
	"time"
)

type PasswordGroup struct {
	GroupID   uint           `gorm:"primaryKey;column:group_id"`
	UserID    uint           `gorm:"column:user_id"`
	Name      string         `gorm:"column:name"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	CreatedBy *string        `gorm:"column:created_by"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	UpdatedBy *string        `gorm:"column:updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
	DeletedBy *string        `gorm:"column:deleted_by"`
}
