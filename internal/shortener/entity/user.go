package entity

type UserURLCollection []string

func (uc UserURLCollection) FindShortURL(shortURL string) (string, bool) {
	for _, item := range uc {
		if item == shortURL {
			return item, true
		}
	}
	return "", false
}

type User struct {
	UUID        string
	SavedURLIDs UserURLCollection
}
