package password

type PasswordEntryTag struct {
	EntryID uint `json:"entry_id" binding:"required"`
	TagID   uint `json:"tag_id" binding:"required"`
}
