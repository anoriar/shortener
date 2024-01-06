package request

// AddURLRequestDto missing godoc.
type AddURLRequestDto struct {
	AuthRequest
	URL string `json:"url"`
}
