package request

type AddURLRequestDto struct {
	Url string `json:"url" valid:"url"`
}
