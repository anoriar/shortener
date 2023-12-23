package entity

// UserURLCollection missing godoc.
type UserURLCollection []string

// FindShortURL missing godoc.
func (uc UserURLCollection) FindShortURL(shortURL string) (string, bool) {
	for _, item := range uc {
		if item == shortURL {
			return item, true
		}
	}
	return "", false
}

// User missing godoc.
type User struct {
	UUID        string
	SavedURLIDs UserURLCollection
}
