package password

type PasswordEntryKey struct {
	EntryID               uint   `gorm:"primaryKey;column:entry_id"`
	EncryptedSymmetricKey string `gorm:"column:encrypted_symmetric_key"`
}
