package out

import "github.com/lib/pq"

type PasswordEntryListResponse struct {
	EntryID   uint            `json:"entry_id"`
	Title     string          `json:"title"`
	GroupName *string         `json:"group_name,omitempty"`
	URL       *string         `json:"url,omitempty"`
	Tags      *pq.StringArray `gorm:"type:text[]" json:"tags,omitempty"`
}
