package in

type PasswordGroupRequest struct {
	Name string `json:"name,omitempty"`
}

type PasswordGroupEntryRequest struct {
	EntryID uint `json:"entry_id,omitempty"`
}
