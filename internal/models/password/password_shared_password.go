package password

import "time"

type SharedPassword struct {
	ShareID               uint      `gorm:"primaryKey;column:share_id"`
	EntryID               uint      `gorm:"column:entry_id"`
	FromUserID            uint      `gorm:"column:from_user_id"`
	ToUserID              uint      `gorm:"column:to_user_id"`
	EncryptedSymmetricKey string    `gorm:"column:encrypted_symmetric_key"`
	SharedAt              time.Time `gorm:"column:shared_at"`
}
