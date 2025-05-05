package password

import "time"

type PasswordHistory struct {
	HistoryID         uint      `gorm:"primaryKey;column:history_id"`
	EntryID           uint      `gorm:"column:entry_id"`
	EncryptedPassword string    `gorm:"column:encrypted_password"`
	ChangedAt         time.Time `gorm:"column:changed_at"`
	ChangedBy         *string   `gorm:"column:changed_by"`
}
