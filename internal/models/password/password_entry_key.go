package password

type PasswordEntryKey struct {
	EntryID               uint   `gorm:"primaryKey;column:entry_id" json:"entry_id,omitempty"`
	EncryptedSymmetricKey string `gorm:"column:encrypted_symmetric_key" json:"encrypted_symmetric_key,omitempty"`
}
