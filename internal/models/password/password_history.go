package password

import "time"

type PasswordHistory struct {
	HistoryID         uint      `gorm:"primaryKey;column:history_id" json:"history_id,omitempty"`
	EntryID           uint      `gorm:"column:entry_id" json:"entry_id,omitempty"`
	EncryptedPassword string    `gorm:"column:encrypted_password" json:"encrypted_password,omitempty"`
	ChangedAt         time.Time `gorm:"column:changed_at" json:"changed_at,omitempty"`
	ChangedBy         *string   `gorm:"column:changed_by" json:"changed_by,omitempty"`
}
