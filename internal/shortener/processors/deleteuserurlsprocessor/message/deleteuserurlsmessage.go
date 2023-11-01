package message

type DeleteUserURLsMessage struct {
	UserID    string
	ShortURLs []string
}
