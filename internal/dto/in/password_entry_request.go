package in

import (
	"github.com/lib/pq"
)

type PasswordEntryRequest struct {
	Title    string          `json:"title"`
	Username string          `json:"username"`
	Password string          `json:"password"`
	Notes    *string         `json:"notes"`
	URL      *string         `json:"url"`
	Tags     *pq.StringArray `json:"tags" gorm:"type:text[]"`
}
