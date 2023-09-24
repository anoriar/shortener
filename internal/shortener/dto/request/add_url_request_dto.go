package request

type AddURLRequestDto struct {
	URL string `json:"url" valid:"url"`
}
