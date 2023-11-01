package request

type DeleteUserURLsRequestDto struct {
	AuthRequest
	ShortURLs []string
}
