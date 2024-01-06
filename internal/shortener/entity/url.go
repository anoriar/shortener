package entity

// URL missing godoc.
type URL struct {
	UUID        string `json:"uuid" db:"uuid"`
	ShortURL    string `json:"short_url" db:"short_url"`
	OriginalURL string `json:"original_url" db:"original_url"`
	IsDeleted   bool   `json:"is_deleted" db:"is_deleted"`
}
