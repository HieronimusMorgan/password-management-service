package password

import "time"

type SharedPassword struct {
	ShareID               uint      `gorm:"primaryKey;column:share_id" json:"share_id,omitempty"`
	EntryID               uint      `gorm:"column:entry_id" json:"entry_id,omitempty"`
	FromUserID            uint      `gorm:"column:from_user_id" json:"from_user_id,omitempty"`
	ToUserID              uint      `gorm:"column:to_user_id" json:"to_user_id,omitempty"`
	EncryptedSymmetricKey string    `gorm:"column:encrypted_symmetric_key" json:"encrypted_symmetric_key,omitempty"`
	SharedAt              time.Time `gorm:"column:shared_at" json:"shared_at,omitempty"`
}
