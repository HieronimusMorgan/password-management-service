package out

import "github.com/lib/pq"

type PasswordEntryListResponse struct {
	EntryID   uint            `json:"entry_id"`
	Title     string          `json:"title"`
	GroupName *string         `json:"group_name"`
	URL       *string         `json:"url"`
	Tags      *pq.StringArray `json:"tags" gorm:"type:text[]"`
}
