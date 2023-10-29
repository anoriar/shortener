package request

type AddURLRequestDto struct {
	AuthRequest
	URL string `json:"url"`
}
